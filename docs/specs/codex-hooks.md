# Codex Hook Dispatch

`vibe-coding-workspace` は、workspace root から child repo を編集する運用を前提に、project-scoped Codex hook を workspace 側で受ける。

## Scope

- workspace root で開始した Codex session
- child repo ごとの formatter / lint / test hook への dispatch
- repo 固有 hook 実装の配置先の規約

この spec は child repo 固有の formatter / lint / test の中身そのものは定義しない。各 child repo の `docs/specs/*` がその契約を持つ。

## Workspace Hook Entry Points

- `vibe-coding-workspace/.codex/config.toml`
  - `features.codex_hooks = true` を有効にする
- `vibe-coding-workspace/.codex/hooks.json`
  - `PostToolUse`
  - `Stop`
- `vibe-coding-workspace/.codex/hooks/dispatch_hook.py`
  - child repo 判定と委譲だけを担当する

workspace 側 hook は repo 固有の formatter / test 実装を持たず、child repo の script を呼び出すだけに留める。

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
