---
title: LLM Knowledge Bases
last_reviewed: 2026-04-11
status: active
sources:
  - ../../sources/2026/2026-04-11-karpathy-llm-knowledge-bases-x-post.md
  - ../../sources/2026/2026-04-11-karpathy-llm-wiki-gist.md
  - ../../sources/2026/2026-04-11-classmethod-llm-knowledge-base.md
---

# LLM Knowledge Bases

## Summary

The key pattern is to let an LLM maintain a compiled Markdown wiki from curated sources instead of repeatedly rediscovering the same information from bookmarks, uploads, or ad hoc search.

## Why it fits this workspace

- The workspace already uses git-tracked Markdown as AI context.
- Knowledge should accumulate across many hobby projects, not stay trapped in chat history.
- A compiled wiki is easier to skim than scattered source links.

## Working model

- The origin point is Karpathy's X post: spend more LLM effort on manipulating knowledge, not only code.
- `sources/` holds source-oriented notes with provenance for articles, papers, repos, images, and posts.
- `schema.md` defines maintenance rules for the agent, similar to how `CLAUDE.md` defines behavior elsewhere.
- `wiki/` holds durable synthesized pages and should keep retrieval anchors like `index.md` and `log.md`.
- `ingest`, `query`, and `lint` form the operating loop.

## Concrete anchors worth preserving

- Viewer and local knowledge UX: `Obsidian`
- Web capture example: `Obsidian Web Clipper`
- Search or memory layer mentioned in the practical write-up: `Mem0`, `pgvector`
- Workspace command example from the article: `/kb-compile`
- Core files called out repeatedly: `schema.md`, `index.md`, `log.md`

## Trade-off framing

- This pattern is not anti-RAG; it is anti-rebuilding everything from scratch on every query.
- The compiled wiki is the durable layer. A search layer can still exist beside it.
- The approach appears especially attractive at small-to-medium scale, before a corpus forces heavier retrieval infrastructure.

## Related pages

- [source-ingestion](../patterns/source-ingestion.md)
- [ww](../projects/ww.md)
