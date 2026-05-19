---
title: エージェント型コーディングワークフロー
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-05-19-agent-quality-controls-slopless.md
  - ../../sources/2026/2026-05-11-zenn-claude-code-local-operations-guide.md
---

# エージェント型コーディングワークフロー

## 現在のシグナル

有用な workflow note は、もはや prompt の書き方だけではない。runtime の置き場所、privacy boundary、deterministic guardrail、そして agent をどう repeatable な loop に入れるかが中心になっている。

## このワークスペースで重要な理由

- このワークスペースは agent-first な開発を既に行っているため、運用資料の再利用価値が高い。
- 同じ repo でも、ある task では local privacy boundary が重要で、別の task では prose quality enforcement が重要になる。
- これらの source は、generic な感想ではなく、具体 command、environment variable、integration loop を残している点が強い。

## 再利用できるパターン

- agent workflow は systems problem として設計する。
- model をどこで動かすか
- どの output に deterministic review をかけるか
- 繰り返しやすい invocation contract をどう保つか
- privacy-sensitive な作業では、hosted default を前提にせず、local の `Messages API` 互換 endpoint へ向きを変える。
- prose や reviewable Markdown には、open-ended rewrite ではなく `slopless` のような deterministic lint feedback を使う選択肢を持つ。
- JSON findings のような raw artifact は、post-run inspection や CI integration に役立つなら保存する。

## 残すべき具体アンカー

- `Claude Code`
- `LM Studio`
- `Qwen3-Coder-30B-A3B-Instruct`
- `ANTHROPIC_BASE_URL`
- `slopless`
- `textlint`
- `.slopless/findings/`

## 未解決の問い

- このワークスペースのどの repo なら、capability 低下を受け入れてでも local-only routing を正当化できるか。
- deterministic な writing check は pre-commit、CI、agent loop 内のどこで回すのがよいか。
- `slopless` を入れるだけの価値がある English-doc subset はどこか。

## 関連ページ

- [llm-knowledge-bases](./llm-knowledge-bases.md)
- [slopless](../tools/slopless.md)
