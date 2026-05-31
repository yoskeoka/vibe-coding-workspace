# Workspace Japanese textlint CI
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Workspace repo に、変更された日本語 Markdown を対象とする `textlint` ベースの PR comment CI を導入する。`slopless` の English-only prose review と並べて、日本語 docs / plan / issue / KB mirror にも deterministic な文言チェックをかけ、後続 spec へ好ましくない表現が流入する前に review できる状態を作る。

## Scope and outcome

- workspace repo で変更された `*.md` を列挙し、日本語を含む file だけを lint 対象にする
- `textlint-rule-preset-ai-writing` を導入する
- repo-local の custom rule で JSONL 置換辞書を読み、禁止/非推奨表現を検知して代替表現を提案する
- 初期辞書として `{"pattern":"\\btaxonomy\\b","replacement":"分類"}` を入れる
- CI は PR ごとに warning annotation と stable marker 付き comment を upsert する
- `README.md` に textlint セクションを追加し、追加辞書の置き場所、JSONL の 1 行 1 ルール形式、追加手順を説明する

## Non-goals

- child repo への同時ロールアウト
- 既存 Markdown 全件の一括修正
- 日本語以外の prose lint を `slopless` 以外へ置き換えること
- 高度な形態素解析や文脈依存の言い換え提案

## Design notes

- changed-file の入口は `slopless` と同じ pull_request CI に寄せるが、filter は逆にして「日本語を含む Markdown のみ」を対象にする
- 対象 path は `docs/specs/` や `docs/kb/ja/` に限定せず、repo 内で変更された Markdown 全般を対象にする。これにより `docs/exec-plan/` や `docs/issues/` の日本語文面も review 対象へ含める
- custom rule は publish 済み package を増やさず repo 内実装とし、辞書メンテを workspace 側で完結させる
- CI failure boundary は `slopless` と同様に「lint findings 自体では job を fail させず、tooling/config の破損だけ fail」とするのが第一候補

## Code changes

### Workflow and helper

- `.github/workflows/` に日本語 textlint CI workflow を追加する
- `tools/` に changed Japanese Markdown helper を追加する
- 必要なら comment upsert / JSON parsing の処理を workflow 内 script または repo-local helper に切り出す

### textlint runtime

- repo root に `package.json` を追加または更新し、`textlint` と `textlint-rule-preset-ai-writing` を devDependency として固定する
- `.textlintrc` 系設定を追加し、preset と repo-local custom rule を有効化する
- repo-local custom rule 実装を `tools/textlint-rules/` 配下に追加する
- 置換辞書を `config/textlint/terms.jsonl` のような repo-local path に追加する

### Documentation

- `README.md` に textlint セクションを追加する
- 必要なら `docs/specs/README.md` に新 spec への導線を追加する

## Spec changes

- `docs/specs/` に日本語 textlint CI の spec を追加する
- spec では以下を明記する
  - trigger と path scope
  - changed Markdown 判定
  - 日本語 file 判定
  - `textlint-rule-preset-ai-writing` の採用
  - JSONL 置換辞書 rule の contract
  - PR annotation / comment の出力形式
  - tooling failure と findings の扱い
- `docs/specs/README.md` に追加 spec を列挙する

## Verification

- helper が changed Japanese Markdown だけを返すことを手元で確認する
- `textlint` が `docs/kb/ja/**` と日本語を含む plan / issue sample に対して動くことを確認する
- 初期辞書 `taxonomy -> 分類` が検知され、comment/summary に代替案が出ることを確認する
- workflow YAML の構文と pin 更新を確認する

## Sub-tasks

- [ ] [parallel] 現行 `slopless` workflow と helper の comment/upsert 方式を流用できる範囲に分解する
- [ ] [parallel] `textlint` config、repo-local custom rule、JSONL 辞書フォーマットを設計する
- [ ] [depends on: workflow helper design, textlint runtime design] 日本語 textlint CI workflow / helper / config を実装する
- [ ] [depends on: textlint runtime design] 初期辞書エントリ `taxonomy -> 分類` を追加する
- [ ] [depends on: CI implementation] `README.md` と関連 spec を更新して運用方法を明記する
- [ ] [depends on: implementation, docs] ローカル verification を実行し、PR comment 出力を確認する

## Success criteria

- workspace PR で変更された日本語 Markdown に対して `textlint` comment が 1 つにまとまって更新される
- `docs/exec-plan/` や `docs/issues/` の日本語表現も path ではなく内容で対象化される
- 辞書 entry の追加方法が `README.md` を読めば分かる
- 初期辞書 rule が repo 内設定だけで再現できる
