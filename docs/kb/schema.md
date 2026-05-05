# Knowledge Base Schema

## Purpose

This file tells an AI agent how to maintain the workspace knowledge base.

The knowledge base is not a dumping ground for raw article text. It is a compiled Markdown wiki built from curated source notes.

## Directory rules

- `sources/YYYY/*.md`: one file per ingested source
- `ja/sources/YYYY/*.md`: Japanese mirror source notes with the same relative path shape when translations exist
- rendered `Sources` year landing pages are generated from `sources/YYYY/*.md` during build/check and are not hand-maintained in git
- `wiki/index.md`: top-level navigation page
- `wiki/log.md`: ingest and maintenance log
- `wiki/projects/*.md`: project-specific compiled pages
- `wiki/topics/*.md`: broad conceptual pages
- `wiki/tools/*.md`: tool and framework pages
- `wiki/patterns/*.md`: reusable methods and heuristics
- `ja/wiki/index.md`: Japanese top-level navigation page
- `ja/wiki/log.md`: Japanese ingest and maintenance log
- `ja/wiki/**`: Japanese mirror wiki pages using the same relative path layout where translations exist

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
source_file:
source_type: article|post|docs|video|video_backed_article|repo|demo
original_language: en|ja|...
ingested_on:
status: active|watch|superseded
tags: []
related_pages: []
conversion_method:
conversion_notes: []
---
```

`original_language` records the language of the original source material, not the locale of the KB note file. A translated Japanese mirror note therefore keeps the same `original_language` as its canonical English counterpart.

`source_file` is optional and should be used when the original source was a local file or another non-durable file identity instead of a stable public URL. Keep it descriptive enough for later human review without depending on a machine-specific absolute path when that path is not durable.

`conversion_method` and `conversion_notes` are optional and should be used when a temporary preprocessing step such as `markitdown` materially affected how the durable source note was drafted.

Video-oriented source notes SHOULD also preserve durable retrieval anchors when available:

```yaml
video_url:
video_platform:
channel:
published_on:
duration:
time_anchors:
selected_screenshots: []
named_entities: []
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
- English pages at `docs/kb/**` are canonical for AI retrieval; Japanese mirror pages under `docs/kb/ja/**` are excluded from QA/retrieval by default
- Japanese mirror pages should keep the same relative layout as the English tree so `/foo/` can map cleanly to `/ja/foo/` when a translation exists

## During ingest

For each new source:
1. Create the source note.
2. Decide which existing wiki pages should be updated.
3. Create new wiki pages only when the concept does not already have a natural home.
4. Update `wiki/index.md` if a new section or notable page was added.
5. Append a dated item to `wiki/log.md`.

When practical, also update the Japanese mirror under `ja/wiki/` and `ja/sources/` with the same relative path layout.

Do not hand-edit generated `Sources` navigation or yearly source landing pages.

When the drafting input came from temporary converted Markdown:
- preserve the original source URL when one exists
- otherwise preserve a stable file identity in `source_file`
- record the converter used and any important cleanup caveats
- do not commit the raw converted Markdown itself

When compressing a source, prefer this order of preservation:
1. concrete product or document names
2. selection criteria or trade-offs
3. workspace-specific recommendation

## Video-specific ingest rules

- Prefer `source_type: video` for direct video URLs and `source_type: video_backed_article` for thin articles whose useful content lives in a referenced video.
- Preserve canonical video metadata even when the original ingest target was the article wrapper.
- Keep segment summaries concise and attach concrete time anchors so the human can jump back into the source quickly.
- Preserve named entities, tool names, library names, commands, and UI labels found through transcript and OCR when they are useful retrieval anchors.
- Store raw transcripts, OCR dumps, extracted frame sets, and other bulky intermediates outside `docs/kb/`.
- Store durable screenshots only when they materially improve comprehension and only under `docs/kb/assets/source-images/<year>/<source-slug>/`.
- Durable screenshots should be referenced from the source note with enough context to explain why the image was kept.
- Do not keep more than a small curated set of screenshots in git. Bulk candidate frames remain temporary job artifacts.
- For video-backed sources, add a visible `## Source` section in the source note body with the source URL, video URL when present, subtitle/transcript availability, and enough retrieval guidance for a human to revisit the source from the rendered wiki page.
- If subtitles or transcripts are retrieved, preserve the durable result by segmenting notes under the same time anchors used in frontmatter. Do not commit full third-party verbatim transcripts unless the source license or user-provided text makes that permitted; keep detailed segment notes and short retrieval anchors instead.
- When screenshots are useful for human skim, keep at least one selected screenshot per important time anchor unless the frame is blank, duplicative, or legally/sensitively unsuitable. In the source-note body, place each screenshot with its matching anchor notes rather than collecting all screenshots in a separate section.

## Document-conversion fallback rules

- Treat converted Markdown as temporary drafting input, not as a durable KB artifact.
- Prefer `source_type: docs` for document-like sources unless another existing type is clearly more accurate.
- Preserve the original file or document identity in frontmatter and, when useful, in a visible `## Source` section in the note body.
- Record conversion caveats when the helper lost ordering, tables, slide structure, or attachments that materially affect interpretation.
- If the base converter produces low-quality output, stop and escalate rather than laundering the damaged output into a durable KB note.

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
