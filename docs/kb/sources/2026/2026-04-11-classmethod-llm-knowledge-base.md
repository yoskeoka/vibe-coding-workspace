---
title: Karpathy 氏が言語化した「LLM Knowledge Base」というパターン
source_url: https://dev.classmethod.jp/articles/karpathy-llm-knowledge-base/
source_type: article
ingested_on: 2026-04-11
status: active
tags:
  - knowledge-base
  - llm
  - workflow
related_pages:
  - ../../wiki/topics/llm-knowledge-bases.md
  - ../../wiki/patterns/source-ingestion.md
---

# Karpathy 氏が言語化した「LLM Knowledge Base」というパターン

## Why it matters here

This article translates the abstract gist into a practical workspace pattern: raw knowledge, schema, and compiled wiki maintained by an agent.

## Summary

- Frames the core move as letting an LLM compile curated documents into structured Markdown.
- Highlights the three layers as raw sources, schema, and wiki.
- Distinguishes the pattern from pure RAG by emphasizing durable compiled knowledge.
- Shows a practical command-driven workflow with targeted compile and lint operations.
- Reinforces that humans should steer curation while the LLM handles routine maintenance.

## Workspace takeaways

- A command such as `kb build` can stay separate from task execution workflows.
- A small command-and-skill setup is enough to make the pattern real.
- Project and topic pages are a better fit than a flat bookmark list.

## Follow-up

- Keep ingest interactive rather than fully automatic.
- Favor source notes over full-page article copies.
