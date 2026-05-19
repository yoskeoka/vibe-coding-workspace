---
title: 「このコード、Claudeに見せていいの？」を解決する — Claude Codeローカル運用ガイド
source_url: https://zenn.dev/shintaroamaike/articles/c7e7e6b27509cc
source_type: article
original_language: ja
ingested_on: 2026-05-19
status: active
tags:
  - ai-coding
  - claude-code
  - local-llm
  - privacy
related_pages:
  - ../../wiki/topics/agentic-coding-workflows.md
---

# 「このコード、Claudeに見せていいの？」を解決する — Claude Codeローカル運用ガイド

## ここで重要な理由

これは privacy 制約下でも agent-style coding workflow を維持するための、かなり実務的な runtime note である。

## 要約

- 記事の核心は思想論ではなく運用論で、`Claude Code` が `Anthropic Messages API` を叩くなら、互換 endpoint を local に立てれば source code を Anthropic に送らずに同じ操作感を保てるという話である。
- 具体 stack は `Claude Code CLI -> LM Studio (lms server) -> local Qwen3-Coder GGUF model`。
- 残すべき runtime anchor も明確である。
- `/v1/messages` 互換のための `LM Studio 0.4.1+`
- 例示 model の `Qwen3-Coder-30B-A3B-Instruct`
- `ANTHROPIC_BASE_URL=http://localhost:1234`
- 切り替え用 shell wrapper の `claude-local()`
- 記事は期待値調整もしていて、frontier Claude と local `Qwen3-Coder` の `SWE-bench Verified` 差は大きい。つまり採用理由は raw capability より privacy / data boundary 側にある。
- context size も重要条件として扱われており、agent workflow では VRAM と context-length が tuning detail ではなく成立条件になる。

## ワークスペースでの含意

- 機密 code を hosted model に送れないが agent CLI workflow は維持したい、という状況の基準資料として残せる。
- 問うべきは「local model が動くか」ではなく、「tool calling、context、latency を含めて repo work に耐えるか」である。
- deterministic linting のような review guardrail と並べて保持すると、制約環境下の agent workflow 設計資料として使いやすい。

## フォローアップ

- 実際に privacy-sensitive repo で使うなら、wrapper script、model identifier、hardware 前提はその repo 側 docs に具体化しておく。
