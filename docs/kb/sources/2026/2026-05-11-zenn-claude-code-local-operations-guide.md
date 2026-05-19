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

## Why it matters here

This is a practical privacy-preserving agent-runtime note: keep the Claude Code interaction model, but redirect the model endpoint to a local stack when code cannot leave the machine.

## Summary

- The article's core claim is operational, not ideological: if `Claude Code` talks to an `Anthropic Messages API` endpoint, then a local compatible server can preserve the workflow without sending source code to Anthropic.
- The concrete stack is `Claude Code CLI -> LM Studio (lms server) -> local Qwen3-Coder GGUF model`.
- The runtime anchors worth preserving are:
  - `LM Studio 0.4.1+` for `/v1/messages` compatibility
  - `Qwen3-Coder-30B-A3B-Instruct` as the example local model
  - `ANTHROPIC_BASE_URL=http://localhost:1234`
  - a shell wrapper such as `claude-local()` to switch models cleanly
- The article also preserves realistic expectations. The cited `SWE-bench Verified` gap between frontier Claude models and the local `Qwen3-Coder` setup is large enough that privacy may be the reason to adopt it, not raw capability.
- Context sizing is called out as a practical requirement: agent workflows need large context windows, so VRAM and context-length settings are first-order concerns rather than tuning details.

## Workspace takeaways

- Keep this as the reference when sensitive code cannot be sent to hosted models but the agent-style CLI workflow is still desirable.
- The operational question is not only "does a local model run" but "does it preserve tool-calling, context, and tolerable latency for real repo work."
- This belongs beside review guardrails such as deterministic linting, because both are about making agent workflows usable in constrained environments.

## Follow-up

- If a privacy-sensitive repo needs this pattern, capture the exact wrapper scripts, model identifiers, and hardware assumptions in that repo's own docs rather than relying only on this source note.
