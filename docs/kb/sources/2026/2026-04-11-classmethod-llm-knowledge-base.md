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

This article translates the abstract gist into a concrete workspace pattern with recognizable tools and commands. It is currently the best bridge between Karpathy's concept and how this workspace can actually operate.

## Summary

- Frames the core move as letting an LLM compile curated documents into structured Markdown instead of only answering over retrieved snippets.
- Highlights the three layers as raw sources, schema, and wiki, and gives concrete examples such as articles, papers, repositories, and images.
- Uses concrete tool names that are worth preserving: `Obsidian`, `Obsidian Web Clipper`, `CLAUDE.md`, `Mem0`, `pgvector`, and the custom command `/kb-compile`.
- Distinguishes the pattern from pure RAG by emphasizing durable compiled knowledge while still allowing a separate search layer.
- Shows concrete maintenance operations such as project-scoped compile, all-project compile, and lint-like checks for contradictions, broken links, and stale pages.

## Workspace takeaways

- A command such as `tools/kb build` should stay separate from task execution workflows.
- A small command-and-skill setup is enough to make the pattern real without introducing a full product or platform first.
- Project pages, topic pages, `index.md`, and `log.md` are better retrieval anchors than a flat bookmark list.

## Follow-up

- Keep ingest interactive rather than fully automatic.
- Favor source notes over full-page article copies.
- Revisit whether the workspace later wants a separate search layer in addition to the compiled wiki.
