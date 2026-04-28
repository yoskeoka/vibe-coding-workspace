# KB MarkItDown Fallback

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Add an explicit `markitdown`-based fallback path for KB ingest when the source format is not well handled by the current conversational URL flow or the existing video-backed pipeline.

This plan assumes the comparison work in `kb-markitdown-comparison.md` has produced a positive enough result for narrow adoption.

## Background

The current KB workflow is strong when:

- the source is a web page the agent can read directly
- the useful content lives in a video and benefits from transcript/OCR/skimming

The current KB workflow is weak when:

- the useful source is a local file instead of a web page
- the source is a PDF slide deck or paper
- the source is a DOCX or PPTX file
- the source is machine-readable but awkward to convert cleanly by hand

`markitdown` is a plausible fallback because it is intentionally optimized for LLM-facing Markdown output instead of human-publishing fidelity. That matches the KB source-note drafting step better than general-purpose format converters do.

At the same time, the fallback should stay narrow:

- not a replacement for article ingest
- not a replacement for video skim/ingest
- not an excuse to ingest raw third-party documents into `docs/kb/` without curation

## Relevant Prior Decisions

Past decision: the KB uses a dedicated skill and durable Markdown artifacts under `docs/kb/`. Apply the same reasoning here by keeping `markitdown` inside the KB skill boundary and by treating converted Markdown as temporary preprocessing output.

Past belief: `Correctness over Speed` means the ingest contract should be written before helpers are added. Apply the same reasoning here by updating the KB spec and ingest docs before adding a new fallback operator path.

## Options

### Option A: Manual ad hoc fallback only

- Continue to let the operator improvise when the source is PDF/DOCX/PPTX

Pros:

- zero implementation cost

Cons:

- repeated manual work
- inconsistent quality and provenance

### Option B: Explicit `markitdown` fallback for unsupported file-like sources

- Keep current primary flows
- Add a documented fallback path when source acquisition should start from file-to-Markdown conversion

Pros:

- narrow and pragmatic
- extends coverage without destabilizing current ingest

Cons:

- adds dependency management and failure modes
- may need secondary OCR/high-fidelity paths for scanned PDFs

### Option C: Make `markitdown` the default preprocessing layer for all sources

- route most sources through `markitdown` first

Pros:

- more uniform preprocessing path

Cons:

- likely redundant for ordinary articles
- weakens the distinction between curated ingest and generic conversion
- higher risk of over-ingesting low-value raw text

## Recommendation Pending Human Confirmation

Recommended option: Option B only.

Adoption rule:

- use `markitdown` when the source is file-like or otherwise unsupported by the normal KB flow
- do not insert it in front of every web article
- do not replace the video-backed pipeline with it

If later evidence suggests broader use, that should come from a separate plan.

## Code Changes

### `skills/knowledge-base/SKILL.md`

- Add a documented fallback branch in the ingest workflow for file-like sources
- Explain when to choose:
  - normal conversational article ingest
  - video skim/ingest
  - `markitdown` fallback

### `skills/knowledge-base/scripts/`

- Add a helper entrypoint or wrapper that:
  - accepts a local path or supported remote URI
  - invokes `markitdown` with the required dependency profile
  - writes the converted Markdown to a temporary job directory
  - passes only the cleaned, narrowed output into the KB source-note drafting flow

Possible shape:

- `kb_markitdown_ingest.py`
- or an extension of an existing KB helper if that keeps the operator flow simpler

### Temporary artifact handling

- Store converted intermediate Markdown outside `docs/kb/`, ideally under OS temp or `.local/kb-ingest/<job-id>/`
- Keep only the curated source note, wiki updates, and optional durable assets in git

### Dependency handling

- Decide whether the workspace should depend on:
  - base `markitdown`
  - a minimal extras set such as `pdf,docx,pptx,xlsx`
  - optional OCR/plugin paths only when explicitly configured
- Add setup guidance and fail-fast dependency checks

## Spec Changes

### `docs/specs/knowledge-base.md`

- Define a third ingest acquisition mode for unsupported file-like sources:
  - convert to Markdown outside `docs/kb/`
  - draft curated KB outputs from that conversion
- Clarify that converted Markdown is temporary and non-durable
- Define the allowed source shapes for the fallback:
  - local files
  - direct document URLs
  - other file-like inputs explicitly supported by the chosen dependency profile
- Define that video-backed sources still prefer the dedicated video pipeline

### `docs/kb/ingest.md`

- Add a decision tree for choosing article ingest vs video ingest vs `markitdown` fallback
- Add operator steps for unsupported-format ingest
- Document when to stop and escalate to a higher-fidelity path such as:
  - scanned PDFs
  - badly ordered multi-column PDFs
  - OCR-heavy slides where the base converter loses too much structure

### `docs/kb/schema.md`

- Clarify provenance expectations when the source note was drafted from a converted file
- Ensure the durable note still records the original source URL or file identity, date, and any conversion caveats

## Design Decisions

- Keep `markitdown` as a fallback, not a default
- Keep converted Markdown temporary and outside `docs/kb/`
- Preserve source provenance explicitly when the input is a local file or document URL
- Do not bundle optional OCR/plugin paths into the first implementation unless the comparison plan proves they are needed

If the fallback path needs a default OCR/provider decision, record that separately in `docs/design-decisions/adr.md`.

## Sub-tasks

- [ ] Confirm the comparison plan produced a positive narrow-scope recommendation
- [ ] Update KB specs, schema, and ingest docs for the fallback decision tree and temporary-artifact contract
- [ ] Add a `markitdown` dependency/setup strategy with fail-fast checks
- [ ] [parallel] Implement a KB helper entrypoint for file-like source conversion
- [ ] [parallel] Update the KB skill instructions to route unsupported formats through the fallback path
- [ ] [depends on: helper entrypoint, docs] Connect converted Markdown to the source-note drafting flow without storing raw intermediates in git
- [ ] [depends on: all above] Verify on at least one PDF and one Office-format source
- [ ] Record any residual gaps that still need separate OCR or layout-specific handling

## Verification

- Confirm the operator can ingest a PDF or Office file through a documented KB fallback path
- Confirm converted Markdown stays outside `docs/kb/` and only curated outputs become durable
- Confirm article URLs and video-backed sources still take their current primary paths
- Confirm source notes preserve enough provenance and caveat data for later human review
