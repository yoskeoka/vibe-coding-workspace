# KB MarkItDown Comparison

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Evaluate whether `markitdown` is a worthwhile addition to the workspace KB ingest stack by comparing it against the current ingest paths on a small, representative source set.

The goal is not to replace the current conversational article ingest or the video-first pipeline blindly. The goal is to answer, with small enough cost, where `markitdown` improves source coverage or operator effort for this workspace.

## Background

Current KB ingest has two strong paths:

- conversational URL/article ingest into curated source notes and wiki pages
- a dedicated video-backed pipeline for YouTube-style sources with skim-first review

Current KB ingest does not yet define a general file-conversion path for formats such as:

- PDF slide decks
- PDF papers
- Word documents
- PowerPoint files outside the video-backed flow

Light evaluation notes captured before planning:

- the current KB spec and skill focus on URLs plus video-backed sources, not generic local/remote document conversion
- `markitdown` 0.1.5 is current on PyPI as of 2026-04-29 and runs as a simple CLI/Python package
- official docs position it as LLM-oriented Markdown conversion for PDF, Office files, HTML, ZIP, EPUB, audio, images, and YouTube, with optional extras and plugin-based OCR support
- official docs also make clear that high-fidelity OCR and scanned-document recovery rely on optional paths such as Azure Document Intelligence or the `markitdown-ocr` plugin rather than the base install

## Relevant Prior Decisions

Past decision: the KB lives in-repo under `docs/kb/` and uses a dedicated `knowledge-base` skill instead of the normal execution workflow. Apply the same reasoning here by evaluating `markitdown` as a KB-ingest helper, not as a general repo-wide document system.

Past belief: `AI-First` and `Trim Tool Output at the Source` both favor narrowing inputs before they reach the model. Apply the same reasoning here by evaluating whether `markitdown` gives the agent cleaner Markdown with lower manual cleanup cost than the current ad hoc fallback behavior.

## Comparison Scope

Run a small comparison matrix that covers:

- one ordinary web article already well-served by current ingest
- one HTML-heavy page with obvious boilerplate such as nav, footer, share UI, or embedded widgets
- one video-backed article or direct YouTube source already well-served by the current video pipeline
- one PDF slide deck
- one text-heavy PDF paper or note
- one Word or PowerPoint document if a suitable sample is available

Measure at least:

- source acquisition friction
- dependency/setup cost
- extraction completeness
- preservation of headings, lists, tables, links, and obvious retrieval anchors
- cleanup effort before a source note can be written
- fit with the existing KB skill and durable artifact rules

For HTML sources specifically, also measure:

- raw HTML size versus `markitdown` output size before any LLM step
- estimated token reduction from passing normalized Markdown instead of raw page HTML
- whether retrieval anchors needed for KB source notes survive the conversion:
  - title
  - headings
  - links
  - concrete tool, product, API, and document names
- boilerplate retention rate:
  - navigation
  - footer/legal text
  - cookie/banner text
  - social/share UI labels
- whether `markitdown` removes noise without removing meaningful article structure or outbound references

## Options To Compare

### Option A: Keep current KB ingest only

- Keep article ingest and video-backed ingest exactly as they are
- Continue handling unsupported formats ad hoc outside the documented KB flow

Pros:

- no new dependency or maintenance surface
- avoids introducing another partially overlapping ingest path

Cons:

- unsupported document formats stay operator-specific and inconsistent
- no shared fallback path for file-centric sources

### Option B: Evaluate `markitdown` as a supplementary converter

- Keep current article and video flows as primary
- Use `markitdown` only as a comparison target and possible later helper for file-like sources

Pros:

- broadens source coverage without forcing a full ingest redesign
- preserves the current high-value video-specific flow

Cons:

- introduces dependency and version-management questions
- may still need separate OCR/high-fidelity fallbacks for scanned or complex PDFs

### Option C: Replace current ingest with `markitdown` where possible

- Prefer `markitdown` for most sources, including some URLs

Pros:

- a single conversion layer is simpler in theory

Cons:

- conflicts with the existing curated conversational ingest design
- likely weakens the specialized video skim/review workflow
- too large a jump for a first adoption step

## Recommendation Pending Human Confirmation

Recommended option: pursue Option B only.

That means:

- compare `markitdown` against current KB ingest on a narrow corpus
- keep current article and video flows as the baseline
- treat replacement ideas as out of scope unless the comparison results are unusually strong

This recommendation should be confirmed before any plan that makes `markitdown` a default path.

## Code Changes

### `skills/knowledge-base/scripts/`

- Add a small repeatable comparison helper or manifest-driven runner for the selected evaluation corpus
- Capture per-source outputs and a normalized scorecard without baking the result into production ingest yet

Possible shape:

- `kb_markitdown_compare.py` or equivalent helper
- a small fixture manifest under `docs/references/` or `skills/knowledge-base/fixtures/`

### `skills/knowledge-base/SKILL.md`

- Document the temporary evaluation workflow so the comparison can be repeated intentionally rather than by memory
- Clarify that this is an evaluation-only path, not the new default ingest flow

### Evaluation artifacts

- Add a durable comparison report under `docs/references/` or a KB-adjacent planning note that records:
  - sources tested
  - environment and dependency choices
  - quality findings
  - recommendation

## Spec Changes

### `docs/specs/knowledge-base.md`

- Add a short note that the KB ingest stack may include a non-default document-conversion helper for unsupported formats, but only after explicit evaluation
- If the comparison results justify it, record the accepted evaluation criteria and the decision boundary for using `markitdown`

### `docs/kb/ingest.md`

- Add a temporary operator note describing how document-format evaluation should be run during the experiment
- Keep the normal article and video ingest flow as the documented primary path

If the comparison concludes "no adoption", the spec update may remain minimal and decision-focused rather than operational.

## Design Decisions

- Do not treat `markitdown` as a replacement for the current video-backed ingest pipeline in this plan
- Do not treat `markitdown` output as durable KB content by itself; it remains a preprocessing artifact
- If the evaluation shows value only for file-like sources, preserve that narrower scope instead of broadening by drift

## Sub-tasks

- [ ] Select a small representative corpus covering current strengths and current gaps
- [ ] Define the comparison rubric and output format for findings
- [ ] Define the HTML-specific comparison rubric for token reduction, boilerplate removal, and retrieval-anchor preservation
- [ ] [parallel] Add a repeatable evaluation helper or manifest-driven command path
- [ ] [parallel] Document the temporary evaluation workflow in the KB skill/docs
- [ ] Run the current KB ingest flow on the in-scope sources where applicable and record operator cost
- [ ] Run `markitdown` on the same corpus with the minimal useful dependency set and record output quality
- [ ] Compare structure preservation, cleanup effort, and missing-content cases
- [ ] Decide whether `markitdown` earns a follow-on integration path, and record the recommendation in durable docs

## Verification

- Confirm the comparison can be rerun from documented commands and sample inputs
- Confirm the result captures both quality and operator-effort trade-offs, not just file-format coverage
- Confirm the comparison distinguishes clearly between:
  - sources already well-served by current ingest
  - sources that are currently awkward or unsupported
