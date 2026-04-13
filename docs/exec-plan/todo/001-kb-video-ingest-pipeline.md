# Knowledge-Base Video Skim and Ingest Pipeline

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Extend the workspace `knowledge-base` skill so it can process video-backed sources, especially YouTube-style technical walkthroughs, without dumping raw transcripts into `docs/kb/`.

The pipeline should:

- support a `skim` mode for quick human review before deciding whether to ingest
- support an `ingest` mode that produces KB-ready source-note drafts and wiki update guidance
- keep heavy processing inside the `knowledge-base` skill directory for portability
- minimize token cost by using Python for mechanical steps and AI only at explicit checkpoints

Addresses:

- article pages whose useful content actually lives in the linked video
- poor human skimmability when transcripts alone miss diagrams, UI states, or slide-like visuals
- unnecessary token spend from sending full transcripts and too many screenshots to an AI model

## Background

Recent KB ingest experiments surfaced a recurring case: the article HTML is thin, while the actual explanation is in an embedded or linked video. For these sources, the current URL ingest flow is insufficient because it lacks:

- transcript fallback handling
- frame extraction and OCR
- representative screenshot selection for fast human review
- a clear boundary between temporary processing artifacts and durable KB content

The intended design is a hybrid pipeline:

1. Python orchestrates download, transcript normalization, frame extraction, OCR, dedupe, and segmenting.
2. AI summarizes segments, selects representative screenshots, and decides what is durable enough for KB inclusion.
3. Only compressed outputs reach `docs/kb/`; raw and intermediate artifacts stay outside the KB tree.

## Relevant Prior Decisions

- `docs/design-decisions/adr.md` records that the knowledge base lives under `docs/kb/` and uses a dedicated `knowledge-base` skill rather than execution-workflow skills.
- `docs/design-decisions/core-beliefs.md` favors AI-first context structures and spec-code parity. This plan follows that by defining durable artifacts up front and keeping human-readable KB outputs compact.

Apply the same reasoning here: keep KB ingest conceptually separate, preserve AI-friendly artifacts, and avoid turning the repo into a raw media archive.

## Spec Changes

### `docs/specs/knowledge-base.md`

- Add video-backed ingest as an explicit supported source shape.
- Define `skim` and `ingest` as separate operating modes.
- Define where temporary job artifacts live and which outputs are durable.
- Define the AI checkpoints:
  - segment summarization
  - representative frame selection
  - KB compile/review
- Define screenshot usage rules for human-skimmable outputs.

### `docs/kb/schema.md`

- Extend source-note guidance for `source_type: video` and video-backed articles.
- Define the minimal durable fields worth preserving from a video:
  - metadata
  - segment summaries
  - named entities
  - time anchors
  - selected screenshots when they materially improve comprehension
- Clarify that raw transcripts and bulk OCR output are not stored in `docs/kb/`.

### `docs/kb/ingest.md`

- Add a video ingest flow that explains:
  - subtitles-first transcript strategy
  - OCR-assisted representative frame selection
  - `skim` before `ingest` when KB value is uncertain

## Code Changes

### `skills/knowledge-base/SKILL.md`

- Expand the skill instructions to cover video-backed ingest.
- Document the new `skim` and `ingest` flows.
- Explain when to keep artifacts in OS temp storage versus a repo-local working directory.

### `skills/knowledge-base/scripts/`

- Add a Python orchestrator for the video pipeline.
- Keep the implementation self-contained under the skill for portability.
- Use existing tools for heavy lifting rather than reimplementing them:
  - `yt-dlp` for metadata, subtitles, and media fetch
  - transcript fallback via an external transcription backend when subtitles are missing
  - `ffmpeg` for candidate frame extraction
  - `PaddleOCR` for text extraction from frames

Planned modules:

- `kb_video_ingest.py` as the CLI entrypoint
- job artifact schema helpers
- fetch/transcript/frame/OCR/dedupe/segment/compile helpers
- prompt templates for AI checkpoints

### `skills/knowledge-base/prompts/`

- Add structured prompt contracts for:
  - segment summarization
  - representative frame selection
  - KB compile/review

### Optional repo-local scratch path

- Support an explicit repo-local working directory such as `.local/kb-ingest/` for resumable jobs.
- Default to OS temp storage when persistence is unnecessary.

## Artifact Model

Temporary job artifacts should be written outside `docs/kb/`, with a schema that supports resume/debug/review. The initial job model should cover:

- job metadata
- normalized transcript segments
- extracted candidate frames
- OCR outputs
- deduped frame candidates by segment
- AI-written segment summaries
- AI-selected representative frames
- source-note and wiki-update drafts

Durable KB outputs remain limited to:

- source notes in `docs/kb/sources/<year>/`
- wiki updates in `docs/kb/wiki/`
- log updates in `docs/kb/wiki/log.md`

## Design Decisions

- Keep the pipeline implementation inside `skills/knowledge-base/` rather than exposing a repo-global helper command for now.
- Prefer subtitles over transcription when available to reduce cost and latency.
- Use OCR plus AI, not AI-only vision, for first-pass frame filtering and dedupe.
- Treat representative screenshots as optional, but allow up to two per segment when they materially improve human skimmability.
- Preserve `skim` as a first-class mode so users can inspect a video before deciding whether it belongs in the KB.

If implementation reveals stable cross-skill utility, extraction into a repo-level helper can be evaluated in a later plan.

## Sub-tasks

- [ ] Update KB specs and ingest/schema docs for video-backed sources, temporary artifacts, and `skim`/`ingest` modes
- [ ] Define the job artifact schema and prompt contracts under `skills/knowledge-base/`
- [ ] [parallel] Scaffold the skill-local script layout and CLI entrypoint
- [ ] [parallel] Implement fetch and transcript normalization with subtitles-first fallback behavior
- [ ] [parallel] Implement frame extraction, OCR, and candidate dedupe
- [ ] [depends on: job schema, fetch/transcript, frame/OCR] Implement segment building and AI segment summarization
- [ ] [depends on: segment summaries, frame/OCR] Implement AI representative-frame selection
- [ ] [depends on: segment summaries, frame selection] Implement `skim` output generation
- [ ] [depends on: skim output] Implement `ingest` draft generation for source notes and wiki update targets
- [ ] [depends on: all above] Verify the pipeline on at least one video-backed source and confirm the outputs stay human-skimmable

## Verification

- Confirm the pipeline can ingest a video-backed source without storing raw transcript/OCR dumps in `docs/kb/`
- Confirm `skim` output includes segment summaries and only the most useful representative screenshots
- Confirm AI checkpoints receive reduced, structured context rather than entire transcripts or all frames
- Confirm `ingest` produces KB-ready draft outputs with durable retrieval anchors
- Confirm the implementation remains self-contained under `skills/knowledge-base/` apart from optional scratch storage

## Expected Outcome

- The workspace gains a practical way to inspect and ingest video-heavy technical references
- Human review becomes faster because important segments can include one or two AI-selected screenshots
- Token spend drops because Python handles mechanical preprocessing and AI only sees narrowed, meaningful context
- The `knowledge-base` skill becomes more portable because the pipeline logic lives inside the skill itself
