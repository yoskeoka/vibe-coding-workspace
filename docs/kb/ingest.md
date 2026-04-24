# Ingest Flow

The intended UX is conversational:

> Ingest these URLs into the knowledge base.

For video-heavy sources, the intended UX also includes:

> Skim this video-backed source first, then ingest it if it looks KB-worthy.

## Expected agent behavior

1. Read the given URLs or the provided source material.
2. Capture the user's framing:
   - why the links matter
   - which projects or themes they seem relevant to
3. Create one source note per URL in `sources/<year>/`.
4. Update the most relevant compiled wiki pages in `wiki/`.
5. Add or refresh Japanese mirror files under `ja/sources/<year>/` and `ja/wiki/` when practical.
6. Add a short entry to `wiki/log.md`.
7. If the source suggests a concrete experiment, record it as a follow-up bullet.

Do not hand-maintain rendered `Sources` navigation or yearly source index pages during ingest. `tools/kb check`, `tools/kb build`, and `tools/kb serve` derive those artifacts automatically.

English remains the canonical AI-facing corpus. Japanese mirror files improve the published site for humans, but QA/retrieval flows should continue to read `docs/kb/**` while excluding `docs/kb/ja/**` unless the user explicitly asks for Japanese mirror content.

## Video-backed flow

When the useful content lives in a video:

1. Run the skill-local video pipeline in `skim` mode first when KB value is uncertain.
2. Prefer subtitles over fresh transcription to reduce cost and latency.
3. Extract candidate frames at a conservative interval and run OCR on them.
4. Dedupe candidate frames before any AI checkpoint.
5. Summarize by segment with time anchors instead of treating the full transcript as one blob.
6. Keep only the most useful representative screenshots for human review.
7. Run `ingest` only after the skim output indicates the source is worth keeping.

`skim` should produce a compact review artifact that includes:
- normalized metadata
- segment summaries
- representative screenshot candidates
- suggested KB relevance and tags

`ingest` should produce:
- a source-note draft
- wiki update guidance or draft content
- a short log entry draft

For video-backed source notes, the durable body should also include:
- a visible `## Source` section so rendered wiki browsing exposes the article URL, video URL, and subtitle/transcript retrieval status instead of hiding them only in YAML frontmatter
- segment notes under the same time anchors used in frontmatter
- selected screenshots placed with their matching anchor notes, normally at least one per important time anchor when the frame improves human skim

Full third-party verbatim transcripts should not be copied into `docs/kb/` unless the source license or user-provided material permits it. When captions are available but full copying is not appropriate, keep detailed segment notes, concrete names, commands, UI labels, and source links so the original can be revisited without repeating the whole ingest process.

## Temporary vs durable outputs

Temporary processing artifacts belong in OS temp storage or `.local/kb-ingest/<job-id>/` when resume/debug value matters.

Durable KB outputs remain limited to:
- `docs/kb/sources/<year>/`
- `docs/kb/ja/sources/<year>/`
- `docs/kb/wiki/`
- `docs/kb/ja/wiki/`
- `docs/kb/wiki/log.md`
- `docs/kb/assets/source-images/<year>/<source-slug>/` for selected screenshots only

Raw transcripts, bulk OCR output, and full frame dumps must not be copied into `docs/kb/`.

## Classification hints

Ask these questions during ingest:
- Which workspace project does this help, if any?
- Is this primarily a tool note, topic note, or reusable pattern?
- Is it a durable insight or just a temporary watch item?
- Should it modify an existing page or create a new one?
- If this is video-backed, which segments are actually durable enough to keep?
- Does any screenshot add meaning that the text summary alone would miss?

## Human review expectations

- The human can skim source notes for accuracy.
- The human can browse the compiled wiki or the rendered Pages site.
- The human can redirect emphasis in the next ingest request instead of editing the wiki manually.
- The rendered site should expose source relationships visibly, even when the source-of-truth links live in frontmatter.
- For video-backed sources, the human should be able to review the skim packet without opening a raw transcript dump.
