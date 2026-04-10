---
title: LLM Knowledge Bases
last_reviewed: 2026-04-11
status: active
sources:
  - ../../sources/2026/2026-04-11-karpathy-llm-wiki-gist.md
  - ../../sources/2026/2026-04-11-classmethod-llm-knowledge-base.md
---

# LLM Knowledge Bases

## Summary

The key pattern is to let an LLM maintain a compiled Markdown wiki from curated sources instead of repeatedly rediscovering the same information from bookmarks or ad hoc search.

## Why it fits this workspace

- The workspace already uses git-tracked Markdown as AI context.
- Knowledge should accumulate across many hobby projects, not stay trapped in chat history.
- A compiled wiki is easier to skim than scattered source links.

## Working model

- `sources/` holds source-oriented notes with provenance.
- `schema.md` defines maintenance rules for the agent.
- `wiki/` holds durable synthesized pages.
- ingest, query, and lint form the operating loop.

## Related pages

- [source-ingestion](../patterns/source-ingestion.md)
- [ww](../projects/ww.md)
