---
title: slopless
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-05-19-agent-quality-controls-slopless.md
---

# slopless

## 現在のシグナル

`slopless` は、English Markdown に対して vague / padded / repetitive な prose を deterministic に検査する review tool である。

## メモ

- `Codex` と `Claude Code` 向けの skill install 導線を持つ CLI として提供される。
- 運用上の最大の anchor は JSON-only output で、script から扱いやすく review もしやすい。
- 想定 loop は agent-oriented で、skill を入れ、新しい session を始め、checker を回し、rewrite し、findings がなくなるまで繰り返す。
- repo wiki には rule inventory、philosophy、comparison、behavior、ignore syntax など、後続比較に使える広い surface がある。
- ただし English-only なので、workspace 全体の docs linter にする前に language boundary を切る必要がある。

## 関連ページ

- [agentic-coding-workflows](../topics/agentic-coding-workflows.md)
