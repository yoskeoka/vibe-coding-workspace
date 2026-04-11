---
title: Source Ingestion
last_reviewed: 2026-04-11
status: active
sources:
  - ../../sources/2026/2026-04-11-karpathy-llm-wiki-gist.md
  - ../../sources/2026/2026-04-11-classmethod-llm-knowledge-base.md
---

# Source Ingestion

## Pattern

The user gives URLs plus a small amount of framing. The agent converts them into source notes and then updates durable wiki pages rather than leaving them as isolated bookmarks.

## Good ingest behavior

- keep one source note per URL
- summarize for workspace usefulness rather than article completeness
- update existing concept pages before creating new ones
- add short follow-up ideas when the source suggests an experiment

## Anti-patterns

- copying long source text into the repository
- creating a brand-new wiki page for every link
- losing provenance when updating a synthesized page
