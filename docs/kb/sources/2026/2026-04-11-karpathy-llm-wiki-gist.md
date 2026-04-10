---
title: llm-wiki
source_url: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f
source_type: article
ingested_on: 2026-04-11
status: active
tags:
  - knowledge-base
  - llm
  - wiki
related_pages:
  - ../../wiki/topics/llm-knowledge-bases.md
  - ../../wiki/patterns/source-ingestion.md
---

# llm-wiki

## Why it matters here

This is the seed idea for turning `vibe-coding-workspace` into a durable AI-maintained knowledge base instead of a pile of bookmarks and open tabs.

## Summary

- Proposes a three-layer model: raw sources, schema, and compiled wiki.
- Treats ingest, query, and lint as the core maintenance loop.
- Emphasizes that the exact file layout should stay domain-specific rather than universal.
- Suggests keeping index and log files so both humans and agents can navigate the wiki.
- Positions the wiki as a persistent knowledge layer that grows through use, not just through one-time summarization.

## Workspace takeaways

- Keep the knowledge base in Markdown that agents can edit directly.
- Use a schema file to make the agent a disciplined maintainer.
- Separate source-oriented notes from durable concept pages.

## Follow-up

- Keep the initial structure small and editable.
- Add a light lint routine once enough pages accumulate.
