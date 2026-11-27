# 0048 変更範囲に絞った Rust lint の運用

> **Execution**: Use `/execute-task` to implement this plan. After implementation is complete, use `/review-task` to prepare and create the PR.

## 目的と完了条件

Rust を変更した各エージェントが、担当した crate とターゲットを引き渡し前に lint する。PR の CI はワークスペース全体を、単独 crate のリポジトリではその crate 全体を検査する。変更側で必要なテストや WASM などの対象別検証は引き続き各リポジトリの契約に従う。

対象は Reversi Adventure、Reversi AI Arena、AI Arena の Rust 製 WASM フィクスチャ。調査時点の管理対象では、ほかのプロジェクトに追跡対象の Cargo マニフェストはなかった。各リポジトリの通常のエージェント指示から対象別コマンドを見つけられ、CI に委ねた範囲が PR で実際に実行されることを確認した時点で完了とする。

## 共通契約と Cargo 設定

- 各エージェントは担当差分と所属 crate を確認し、Rust ソース、マニフェスト、ビルドスクリプト、または crate が取り込む入力を変更した場合に、その crate の lint を実行する。共有入力の変更は影響先の crate も選ぶ。共有依存やツールチェーン設定の影響を絞れない場合は、ローカルで確認した代表範囲と CI に委ねた範囲を明記する。Rust に影響しない差分では Rust lint を起動しない。
- 通常の Rust ソースは選択した crate の lib／bin を対象とする。テスト、example、bench、および対応する条件付きコードを変更した場合は、その crate の全ターゲットを対象にする。WASM を対象とするコードには実際のターゲット指定を使う。非デフォルト feature のコードを変更したときは必要な feature を明示して確認する。
- 選択式のリポジトリでは `make rust-clippy-local CLIPPY_PACKAGES="crate-a crate-b"` を用意する。テスト等には `CLIPPY_ALL_TARGETS=1` を使う。選択が空、完全なパッケージ名以外の指定、または不明な crate の場合は失敗を返す。Cargo へ crate ごとの `-p` と `-- --no-deps -D warnings` を渡す。参照する依存のコンパイルは Cargo に任せる。必要な feature を指定する場合の Cargo 直接実行例も案内する。
- ローカル lint コマンドの `CARGO_BUILD_JOBS` は未指定時に `2` とし、呼び出し側の指定を尊重する。設定はローカル入口に適用する。通常の worktree のビルドディレクトリとキャッシュを再利用し、同じディレクトリへの Cargo 実行を直列にする。全リポジトリ共通の `target` やグローバルな Cargo 設定は導入しない。
- 各エージェントは引き渡し時に対象 crate、必要ならターゲットと feature、実行結果を記録する。親エージェントはその記録を PR に集約し、担当差分が増えた場合や報告後に関連箇所を修正した場合に該当範囲を確認する。CI の結果は PR の最新コミットで確認する。

## リポジトリ別の変更

### ワークスペース

- (MODIFY) `docs/specs/workflow-context-contract.md` の Preserved enforcement と `AI_WORKFLOW.md` の Execution／PR and follow-up。担当エージェントの対象確認、PR での確認結果の集約、CI での全体チェックを定義する。品質ゲートの具体的なコマンドは各リポジトリの `AGENTS.md` から辿れるようにする。

### Reversi Adventure

- (NEW) `docs/specs/rust-quality-gates.md`、(MODIFY) `docs/specs/README.md` と `AGENTS.md` の Build & Lint Commands（現行 105–110 行）。対象別 lint、CI の全体 lint、既存のエンジン・AI テストと GDExtension ビルドの責務を記す。
- (MODIFY) `Makefile` の `CARGO` 定義周辺に `rust-clippy-local` を追加する。パッケージ例は `reversi-engine`、`reversi-ai`、`reversi-godot`。前二者の選択時に Godot を Cargo のコンパイル対象に含めないことを確認する。関係箇所はルート `Cargo.toml` の `workspace.members` と `rust/reversi-godot/Cargo.toml` の依存定義。
- (MODIFY) `.github/workflows/ci.yml` の `Run clippy`（21–22 行）。現行の `--workspace` に `--all-targets` を加えてテスト用ターゲットも CI で lint する。既存のテスト、GDExtension、Oracle のジョブを維持する。

### Reversi AI Arena

- (MODIFY) `docs/specs/quality-gates.md` の Required Local Commands、Minimum Pre-PR Verification、CI Contract と、`docs/specs/language-toolchain-policy.md` の Verification Expectations。全体の `make verify-rust` を CI のゲートとして位置づけ、担当範囲のローカル lint と既存の WASM 要件を記す。エージェントの入口を `AGENTS.md` に追加し、`README.md` の検証コマンド案内を揃える。
- (MODIFY) `Makefile` の `rust-clippy` 周辺と `tools/rust-ci.sh` の `clippy`／`verify` ディスパッチ（14–26、48–53 行）。`rust-clippy-local` から選択式モードを呼び出す。例は `reversi-game`、`aiarena-protocol`、`reversi-rust-reference-player`。WASM 側を変更した crate の lint には `wasm32-wasip1` を指定できるようにする。既存の `rust-clippy`、`verify-rust` は全体を検査する入口として維持する。
- (MODIFY) `.github/workflows/rust-ci.yml` の PR／push の対象パスと `detect` の `full`。追跡されている Rust workspace member を列挙し、現行で漏れている `internal/aiarena-protocol/**` と、新しい lint 設定を追加した場合はその設定ファイルも網羅する。共有 crate だけを変えた差分で `workspace-verify` が選ばれることを確認する。

### AI Arena

- (MODIFY) `AGENTS.md` と `docs/specs/janken-game.md` の Phase 4 Multi-language Evaluation Lane（252 行以降）。既存の Rust フィクスチャを変更したエージェントの lint と、WASM CI での lint および動作検証を記す。
- (MODIFY) `Makefile` の Rust WASM フィクスチャターゲット周辺（現行 237–240 行）。`make rust-clippy-local` は `testdata/ai/janken/janken-rust-wasm-ai/Cargo.toml` と `wasm32-wasip1` を固定して呼ぶ。単独 crate のため package 指定は不要。テスト等の変更時は `CLIPPY_ALL_TARGETS=1` を使い、CI 用の入口 `make rust-clippy` は常に全ターゲットを検査する。
- (MODIFY) `.github/workflows/wasm-verification.yml` の `rust-wasm`。Clippy コンポーネントを利用可能にし、既存の Rust-WASM ビルド・E2E の前にフィクスチャの全ターゲット lint を実行する。Janken フィクスチャと Makefile の既存パス条件を利用する。

## 実施順序

1. この計画の承認後、ワークスペースの契約と先行する Adventure の子リポジトリ計画を用意する。子リポジトリ側では個別の番号付き計画をレビューし、承認後に仕様、コマンド、エージェント指示、CI の順で実装する。
2. Adventure で選択式コマンドの引数、ジョブ上限、Godot を含まない実行を確認する。その結果を反映して Arena 2件の個別計画を提出する。Arena 2件の実装は独立して進められる。
3. 各実装 PR の最新コミットで CI の対象と結果を確認する。全リポジトリの結果をワークスペースの実装 PR に集約し、完了時にこの計画を削除する。

## 検証

- Cargo の呼び出しを記録できる代替コマンドで、crate の複数指定、空指定、未知の crate、全ターゲット指定、ジョブ数の既定値と上書き、Cargo の失敗伝播を検証する。WASM 対象のコマンドは manifest・target の引数も照合する。
- Adventure は `reversi-ai`、Reversi AI Arena は `aiarena-protocol`、AI Arena は Janken のフィクスチャで実際のローカル Clippy を各1回確認する。Adventure の通常の選択式実行では Cargo の詳細ログで Godot のコンパイルが走らないことを確認する。エージェントは自分の差分に応じた lint 結果を PR に残す。
- 実装 PR の CI で各ワークスペースまたはフィクスチャの全ターゲット lint が通り、既存の関連ゲートが維持されることを確認する。Reversi AI Arena は共有 crate 単独の変更パスを使ったフィルター検証も行う。
- 計画 PR ではドキュメントの差分、構造、workflow lint を確認する。実装 PR では各リポジトリの workflow lint と必要な CI 定義の検証を実行する。
