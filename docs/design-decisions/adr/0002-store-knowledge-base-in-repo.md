# Store the workspace knowledge base in-repo under `docs/kb/`

## Status

Accepted — 2026-04-11

## Context

The workspace needs durable references for AI-centered development. GitHub Wiki
was considered, but git-tracked Markdown under `docs/` is already canonical AI
context and also needs a human-readable publishing path.

## Decision

Keep the knowledge base in `docs/kb/`: immutable source notes in `sources/`,
compiled pages in `wiki/`, and maintenance rules in `schema.md` and `ingest.md`.
Use the `knowledge-base` skill for ingestion and MkDocs Material to publish the
same Markdown.

## Consequences

- Sources stay easy to edit, diff, review, and publish.
- Knowledge-base ingest stays separate from execution-plan workflow skills.
- Raw Markdown and rendered Pages are two documentation surfaces and require
  disciplined structure.
