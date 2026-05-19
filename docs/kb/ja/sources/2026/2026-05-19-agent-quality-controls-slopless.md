---
title: slopless
source_url: https://github.com/agent-quality-controls/slopless
source_type: repo
original_language: en
ingested_on: 2026-05-19
status: active
tags:
  - writing
  - lint
  - ai-coding
  - markdown
related_pages:
  - ../../wiki/topics/agentic-coding-workflows.md
  - ../../wiki/tools/slopless.md
---

# slopless

## ここで重要な理由

これは AI 補助の writing loop に deterministic な品質管理を差し込むための、かなり具体的な tool である。

## 要約

- `slopless` は LLM 呼び出しではなく deterministic な `textlint` rule 群で Markdown prose を検査する。
- repo の説明では、AI 由来と human 由来の両方の prose slop を検出し、50 以上の rule と JSON-only の出力を持つ。
- 想定 loop は明確で、package を入れ、Codex か Claude 用 skill を導入し、新しい agent session を始めて、findings JSON が空になるまで rewrite を繰り返す。
- CLI 契約も具体的に覚えておく価値がある。
- `npx slopless "docs/**/*.md"` で file / glob を検査する
- exit `0` は clean
- exit `1` は findings あり
- exit `2` は failure または invocation error
- 生の findings JSON は `.slopless/findings/` に保存する運用が推奨されており、agent loop の外でも review しやすい。
- repo wiki の `Philosophy`、`Comparison`、`Rules`、`Behavior`、`Ignore-Rules` は後続調査の retrieval anchor になる。

## ワークスペースでの含意

- English docs を repeatable に prose cleanup したいときの guardrail 候補として保持する価値がある。
- deterministic JSON 契約があるので、CI や scripted review loop に載せやすい。
- ただし scope は English-only なので、このワークスペースの日本語 internal docs 全体にそのまま適用する前提にはできない。

## フォローアップ

- English KB page や PR 向け doc の prose QA が必要なら、`slopless` を既存の `Vale`、`textlint`、review-only 運用と比較する。
