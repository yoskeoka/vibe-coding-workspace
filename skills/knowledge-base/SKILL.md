---
name: knowledge-base
description: Use when the user wants to ingest URLs or notes into the workspace knowledge base, update compiled wiki pages, file durable answers back into the wiki, or lint the knowledge base for staleness and weak structure.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Knowledge Base

This is a workspace-only skill for maintaining `docs/kb/`.

## When to use

- The user says to ingest URLs, articles, notes, or references into the knowledge base
- The user asks a question that should be answered from `docs/kb/` and preserved
- The user asks to lint, review, or reorganize the knowledge base

## Read first

1. `docs/kb/schema.md`
2. `docs/kb/ingest.md`
3. `docs/kb/wiki/index.md`

## Ingest workflow

1. Read the provided URLs or source material.
2. Create one source note per source under `docs/kb/sources/<year>/`.
3. Update the most relevant pages under `docs/kb/wiki/`.
4. Update `docs/kb/wiki/index.md` if navigation should change.
5. Append a dated line to `docs/kb/wiki/log.md`.
6. If the source suggests a concrete experiment, record it in the relevant page.

Prefer updating existing wiki pages over creating new ones. Preserve provenance. Avoid copying long source text.
Keep concrete retrieval anchors such as service names, library names, product names, and document names when they are part of the source's value. Do not collapse `Render vs Cloud Run` into `easy backend vs fast backend` if the original concrete options are useful later.

## Query filing-back workflow

If a user asks a question and the answer is durable:
- create or update a compact page under `docs/kb/wiki/`
- link it from a relevant topic, pattern, tool, or project page
- note the change in `docs/kb/wiki/log.md`

## Lint workflow

Review the knowledge base for:
- orphan pages
- duplicated topics
- stale date-sensitive claims
- weak cross-linking
- source notes that are not reflected in the wiki

## Rendering

- Use `tools/kb check` to validate the structure.
- Use `tools/kb build` to build the Pages site.
- Use `tools/kb serve` for local preview if needed.
- `tools/kb` will prefer `uv` when available and otherwise use the local Python MkDocs installation.
