---
title: ww
last_reviewed: 2026-04-11
status: seed
sources:
  - ../../sources/2026/2026-04-11-karpathy-llm-knowledge-bases-x-post.md
  - ../../sources/2026/2026-04-11-karpathy-llm-wiki-gist.md
  - ../../sources/2026/2026-04-11-classmethod-llm-knowledge-base.md
---

# ww

## Relevance

`ww` is the workspace tooling project, so it is the most natural place to absorb future lessons about AI workflow ergonomics, multi-agent orchestration, and durable knowledge support.

## Current knowledge hooks

- LLM knowledge-base ideas likely feed into future workspace tooling rather than the CLI directly.
- If `ww` gains knowledge-aware commands later, this page can track those patterns.
- The concrete command shape in the article, `/kb-compile`, is a useful reference point for future `ww` or workspace tooling decisions.
- The practical stack described around the pattern also suggests an interface boundary: compiled wiki in Markdown first, optional search layer such as `Mem0` and `pgvector` second.

## Related pages

- [llm-knowledge-bases](../topics/llm-knowledge-bases.md)
