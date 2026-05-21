# Codex Hook Dispatch

`vibe-coding-workspace` は、workspace root から child repo を編集する運用を前提に、project-scoped Codex hook を workspace 側で受ける。workspace 自身の English Markdown prose QA も同じ hook 入口で扱う。

## Scope

- workspace root で開始した Codex session
- child repo ごとの formatter / lint / test hook への dispatch
- workspace 自身の English Markdown に対する prose lint
- repo 固有 hook 実装の配置先の規約

この spec は child repo 固有の formatter / lint / test の中身そのものは定義しない。各 child repo の `docs/specs/*` がその契約を持つ。

## Workspace Hook Entry Points

- `vibe-coding-workspace/.codex/config.toml`
  - `features.hooks = true` を有効にする
- `vibe-coding-workspace/.codex/hooks.json`
  - `PostToolUse`
  - `Stop`
- `vibe-coding-workspace/.codex/hooks/dispatch_hook.py`
  - child repo 判定と委譲だけを担当する
- `vibe-coding-workspace/.codex/hooks/slopless_post_tool_use.py`
  - workspace 自身の English Markdown にだけ `slopless` を適用する

workspace 側 hook は child repo の formatter / test 実装を持たず、repo 固有の重い品質ゲートは child repo script に委譲する。workspace 自身で持つローカル判定は English Markdown 向け prose lint に限る。

## Workspace Markdown Lint

workspace root の `PostToolUse` は、`apply_patch`, `Edit`, `Write` で触れた file path のうち `*.md` にだけ `slopless` を適用してよい。

- 実行タイミングは `Stop` ではなく `PostToolUse` とする
- 理由は、Markdown 編集直後に対象 file 単位で失敗させた方が feedback が早く、session stop 時に repo 全体を再走査する必要もないため
- 実行対象は file content に日本語文字が含まれない Markdown に限る
- `slopless` は `npx` 実行を使うが、package version は固定し、対応する npm `gitHead` が期待する GitHub commit SHA と一致することを hook 側で確認してから起動する
- npm cache は repo 外の一時 directory を使ってよい

この lint は workspace 自身の prose QA 用であり、child repo 側 `Stop` hook のような広い lint/test gate の代替ではない。

## Child Repo Contract

各 child repo が Codex hook を持つ場合、workspace 側 dispatcher が期待する entrypoint は以下とする。

- `<child-repo>/tools/codex-hook-post-tool-use.sh`
- `<child-repo>/tools/codex-hook-stop.sh`

workspace 側 dispatcher は、session cwd と hook payload から対象 child repo を判定し、対象 repo の script に同じ stdin JSON payload をそのまま渡す。

## ai-arena Routing

`ai-arena` については、workspace 側 dispatcher は以下のどちらかを見つけられれば委譲対象として扱ってよい。

- `<workspace-root>/ai-arena`
- `<workspace-root>/.worktrees/...` とは別に存在する canonical workspace checkout 上の `ai-arena`

`PostToolUse` dispatch は `apply_patch`, `Edit`, `Write` を対象にし、hook payload の `tool_input` 全体が `ai-arena/` を含むか、session cwd 自体が `ai-arena` repo 内にある場合にだけ起動する。

`Stop` dispatch は session cwd が `ai-arena` repo 内にある場合、または canonical workspace checkout 上の `ai-arena` が存在する場合に child repo 側 script へ委譲してよい。重い実行を避けるための最終判定は child repo 側 script が持つ。

## Restart Behavior

- `.codex/config.toml` と `.codex/hooks.json` の変更は、以後の Codex session で使う前提とする
- Codex は起動時に project-scoped hook を読む前提なので、設定変更後は Codex restart が必要になることがある
