---
title: llm-wiki
source_url: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f
source_type: docs
original_language: en
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

This is the most implementation-relevant source for the workspace KB. It turns the X post into an explicit operating model that can be adapted into repo-native Markdown and agent skills.

## Summary

- Proposes a three-layer model: raw sources, `schema` instructions, and a compiled Markdown `wiki`.
- Treats `ingest`, `query`, and `lint` as the core maintenance loop instead of a one-shot summarize-everything workflow.
- Keeps the exact file layout intentionally abstract, but explicitly calls for durable navigation files such as `index.md` and `log.md`.
- Positions the wiki as something the LLM maintains over time through summaries, entity or concept pages, backlinks, and incremental updates.
- Argues that at roughly `~100` articles and `~400K` words, a maintained wiki can be useful before reaching for "fancy RAG".

## Workspace takeaways

- Keep the knowledge base in Markdown that agents can edit directly and humans can diff in git.
- Use a schema file to make the agent a disciplined maintainer instead of an unconstrained summarizer.
- Separate source-oriented notes from durable concept pages and keep `index.md` / `log.md` as retrieval anchors.

## Follow-up

- Keep the initial structure small and editable instead of copying the gist literally.
- Add a light lint routine once enough pages accumulate.
