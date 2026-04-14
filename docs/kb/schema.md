# Knowledge Base Schema

## Purpose

This file tells an AI agent how to maintain the workspace knowledge base.

The knowledge base is not a dumping ground for raw article text. It is a compiled Markdown wiki built from curated source notes.

## Directory rules

- `sources/YYYY/*.md`: one file per ingested source
- rendered `Sources` year landing pages are generated from `sources/YYYY/*.md` during build/check and are not hand-maintained in git
- `wiki/index.md`: top-level navigation page
- `wiki/log.md`: ingest and maintenance log
- `wiki/projects/*.md`: project-specific compiled pages
- `wiki/topics/*.md`: broad conceptual pages
- `wiki/tools/*.md`: tool and framework pages
- `wiki/patterns/*.md`: reusable methods and heuristics

## Maintenance rules

- Prefer updating an existing wiki page over creating a near-duplicate.
- Preserve source provenance. Every wiki page should point to at least one source note when a concrete claim depends on it.
- Summaries should be concise and biased toward workspace usefulness.
- Do not summarize away concrete retrieval anchors. If a source introduces specific services, tools, libraries, frameworks, APIs, or documents that would help future search or comparison, keep those names in the source note and in the relevant wiki page.
- Human-facing pages should explain "why this matters here", not just "what the article said".
- Avoid copying long passages from third-party sources. Summarize instead.
- If a source looks time-sensitive, record the relevant date in the source note.
- Treat `wiki/log.md` as an audit trail, not a narrative article.

## Naming

- Source notes: `YYYY-MM-DD-short-slug.md`
- Wiki pages: kebab-case by concept or project

## Frontmatter conventions

Source note frontmatter:

```yaml
---
title:
source_url:
source_type: article|post|docs|video|repo|demo
ingested_on:
status: active|watch|superseded
tags: []
related_pages: []
---
```

Wiki page frontmatter:

```yaml
---
title:
last_reviewed:
status: active|watch|seed
sources: []
---
```

In the rendered site:
- wiki pages are the primary top-level browsing surface
- `README.md`, `schema.md`, and `ingest.md` remain source files and published pages, but they do not need to lead the rendered top-level nav
- `Sources` opens to yearly groups by default, while individual source-note leaves remain collapsed until opened through normal navigation
- wiki page `sources:` frontmatter is converted into a visible `## Sources` section
- source-note `related_pages:` frontmatter is converted into a visible `## Related pages` section
- these visible sections are derived from frontmatter during build/check and should not be duplicated manually unless the prose itself needs additional context

## During ingest

For each new source:
1. Create the source note.
2. Decide which existing wiki pages should be updated.
3. Create new wiki pages only when the concept does not already have a natural home.
4. Update `wiki/index.md` if a new section or notable page was added.
5. Append a dated item to `wiki/log.md`.

Do not hand-edit generated `Sources` navigation or yearly source landing pages.

When compressing a source, prefer this order of preservation:
1. concrete product or document names
2. selection criteria or trade-offs
3. workspace-specific recommendation

## During query filing-back

- Only file back durable outputs.
- Prefer short synthesis pages over dumping chat transcripts.
- Link the new page from an existing topic or project page.

## During lint

Look for:
- pages with no inbound links
- stale notes that need reconfirmation
- duplicated topic pages
- sources not connected to any wiki page
- project pages missing obvious relevant sources
