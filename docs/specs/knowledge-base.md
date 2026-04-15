# Spec: Workspace Knowledge Base

## Goal

Provide a repo-native knowledge base for AI-centered personal development work. The knowledge base must let the user hand URLs directly to an agent for ingest, keep durable Markdown notes in git, and publish a human-readable wiki via GitHub Pages.

## Scope

- Knowledge capture for external references relevant to the workspace
- AI-maintained source notes and compiled wiki pages
- Human-readable rendering from the same Markdown files
- Lightweight command support for build, serve, and structural checks

## Non-Goals

- Full-content archival of third-party articles
- Replacing `docs/specs/`, `docs/project-plan.md`, or `docs/exec-plan/`
- Replacing search or RAG systems with a large retrieval stack
- Fully automated crawling without human curation

## Directory Layout

The knowledge base MUST live under `docs/kb/`.

```text
docs/kb/
  README.md
  schema.md
  ingest.md
  sources/
    YYYY/
      YYYY-MM-DD-slug.md
  wiki/
    index.md
    log.md
    projects/
    topics/
    tools/
    patterns/

Generated build inputs MAY be created under `.local/` during build, check, or preview, but the git-tracked source of truth remains `docs/kb/`.

Video ingest jobs MAY also use:

```text
.local/kb-ingest/
  <job-id>/
```

for resumable repo-local scratch data. When the user does not need persistence, the video pipeline SHOULD default to OS temp storage instead.
```

## Content Model

### 1. Source Notes

Each ingested reference MUST create one Markdown file in `docs/kb/sources/<year>/`.

A source note MUST include:
- title
- source URL
- source type
- ingested date
- status
- tags
- workspace relevance
- concise summary bullets
- candidate follow-up actions
- links to related wiki pages

Video-backed sources MAY additionally include:
- canonical video URL when the source URL is a thin article that points to a video
- video duration
- channel or author metadata
- segment summaries with time anchors
- named entities or tool names recovered from transcript/OCR
- selected screenshot references when screenshots materially improve later human review

Source notes are source-oriented and append-only in spirit. They MAY be corrected, but they should not become general wiki pages.

### 2. Wiki Pages

Compiled knowledge MUST live under `docs/kb/wiki/`.

Wiki pages MUST:
- focus on durable concepts rather than one-off links
- synthesize across multiple sources when available
- link back to source notes
- include `last reviewed` metadata

Recommended wiki groupings:
- `projects/` for workspace project relevance
- `topics/` for broad themes
- `tools/` for tool- or framework-specific notes
- `patterns/` for reusable practices

### 3. Navigation Files

The knowledge base MUST contain:
- `docs/kb/wiki/index.md` as the primary human and AI entry point
- `docs/kb/wiki/log.md` as an append-only-ish ingest and maintenance log

The rendered site MUST expose a `Sources` navigation section grouped by year.

Each year entry SHOULD use a compact label such as `2026 (6)` and point to a generated yearly landing page in the rendered site.

Those yearly landing pages are publish artifacts derived from `docs/kb/sources/YYYY/*.md` and MUST NOT become hand-maintained source files in git.

## Operations

### 1. Ingest

The primary ingest UX is conversational:
- The user gives one or more URLs and asks the AI to ingest them.
- The AI uses the `knowledge-base` skill.

Ingest MUST:
- create or update source notes
- update the relevant wiki pages
- update `docs/kb/wiki/index.md` if navigation changes
- append a short entry to `docs/kb/wiki/log.md`

The knowledge-base ingest flow MUST support two video-oriented operating modes:
- `skim`: produce a compact review artifact for deciding whether the source belongs in the KB
- `ingest`: produce KB-ready source-note and wiki-update drafts from the already skimmed job data

Ingest SHOULD:
- preserve the user's framing about why the source matters
- prefer updating an existing wiki page over creating duplicates
- identify which workspace projects, tools, topics, and patterns are affected
- prefer subtitles over fresh transcription when the source provides usable subtitles
- narrow AI context to structured segment summaries and representative frame candidates instead of passing full transcripts or all frames

Ingest does NOT need to hand-maintain rendered `Sources` nav entries or yearly source index pages. Those are derived during build/check.

For video-backed sources, ingest MUST keep raw transcripts, OCR dumps, extracted frame sets, and other bulky intermediates outside `docs/kb/`.

### 1a. Skill discovery

The `knowledge-base` skill is workspace-only and MUST NOT be distributed to child repos via `setup-workspace.sh`.

The workspace repo itself MUST expose the skill through `.claude/skills/knowledge-base` so local agents can invoke it. Agent directories that mirror `.claude/skills/` (for example `.agents/skills/`) MUST therefore resolve the same skill as well.

The skill MUST keep heavy video-processing helpers, prompt contracts, and Python runtime files inside `skills/knowledge-base/` for portability.

## Video Pipeline

Video-backed ingest MUST support direct video URLs and thin articles whose substantive content lives in a referenced video.

The skill-local entrypoint MUST provide:
- `skim` to fetch/normalize video context and generate a human-skimmable review packet
- `ingest` to turn a prepared or resumed job into KB-ready drafts
- `--check-deps` to fail fast with actionable dependency guidance before long-running work starts

The implementation MUST orchestrate these stages:
1. metadata fetch
2. subtitles-first transcript acquisition
3. candidate frame extraction
4. OCR over candidate frames
5. candidate dedupe and segment building
6. AI checkpoint preparation for segment summarization
7. AI checkpoint preparation for representative-frame selection
8. `skim` or `ingest` output compilation

### Dependency model

The video pipeline MUST depend on:
- `yt-dlp` for metadata, subtitles, and media fetch
- `ffmpeg` for frame extraction
- a skill-local Python runtime with `paddleocr`, a CPU `paddlepaddle` variant, and image-processing helpers

The pipeline MUST check for these dependencies at startup and exit with actionable setup instructions if they are missing.

The skill docs MUST include install guidance for:
- macOS/Homebrew
- Debian/Ubuntu
- environment-matched Paddle runtime selection for CPU vs GPU hosts

The CLI SHOULD inspect the runtime environment before heavy work starts and choose safe defaults for:
- OCR batch size
- frame extraction cadence
- GPU enablement when supported by both the host and the installed Paddle runtime

The initial tuning heuristics MAY use signals such as:
- host OS or WSL detection
- available system memory
- CUDA / GPU availability

### AI checkpoints

The implementation MUST define structured prompt contracts for:
- segment summarization
- representative-frame selection
- KB compile/review

The initial implementation MAY emit prompt-ready markdown or JSON payloads instead of invoking an LLM directly, as long as the job artifacts cleanly separate the three checkpoints.

### Temporary and durable artifacts

Temporary job artifacts MUST stay outside `docs/kb/` and SHOULD be resumable. A job directory MUST be able to store:
- job metadata
- normalized transcript segments
- extracted candidate frames
- OCR outputs
- deduped frame candidates by segment
- AI-written or AI-ready segment summary payloads
- AI-selected or AI-ready representative-frame payloads
- source-note and wiki-update drafts

Durable KB outputs remain limited to:
- source notes in `docs/kb/sources/<year>/`
- wiki updates in `docs/kb/wiki/`
- `docs/kb/wiki/log.md`
- selected representative screenshots under `docs/kb/assets/source-images/<year>/<source-slug>/` when screenshots materially improve human skimmability

Raw transcripts, bulk OCR output, and bulk extracted frames MUST NOT be committed under `docs/kb/`.

### Screenshot policy

Representative screenshots are optional. The pipeline SHOULD keep them only when they add meaning beyond transcript text alone, such as:
- UI state changes
- diagrams or slide content
- code or terminal output that anchors a segment summary

The durable screenshot budget SHOULD stay small. The initial implementation SHOULD target at most two selected screenshots per segment and MUST avoid persisting bulk candidates in `docs/kb/`.

### 2. Query

Questions answered from the knowledge base SHOULD cite the relevant source notes or wiki pages. Durable outputs from those questions MAY be filed back into `docs/kb/wiki/`.

### 3. Lint

The knowledge base SHOULD support periodic AI maintenance to detect:
- orphan pages
- stale claims
- duplicate topics
- missing backlinks
- pages that should be split or merged

### 4. Compile

Human-readable publishing MUST use MkDocs with:
- a git-tracked template config file `mkdocs.kb.template.yml`
- a generated effective config file produced during build/check/serve
- a generated docs tree that may include publish-only files such as yearly source index pages

The local helper command `tools/kb` MUST provide:
- `build` to render the site
- `serve` to preview the site locally
- `check` to validate required files and directories

`build` SHOULD run in strict mode so broken links and missing pages fail fast.
The helper SHOULD prefer `uv` when available and fall back to `python3 -m mkdocs` otherwise.

The helper MUST derive the rendered `Sources` nav from `docs/kb/sources/YYYY/*.md` so newly added source notes cannot be omitted accidentally.

## Rendering and Publishing

- The git-tracked Markdown under `docs/kb/` MUST remain the source of truth for the rendered site, but the published site MAY be built from a derived docs tree that includes publish-only artifacts and sections derived from frontmatter.
- The published top-level navigation MUST treat wiki content as the primary browsing surface:
  - `wiki/index.md` MUST be the leading top-level page in the rendered nav.
  - the current wiki groupings (`Log`, `Topics`, `Tools`, `Patterns`, `Projects`) SHOULD remain top-level rendered entries instead of being nested under a `Wiki` parent.
  - `README.md`, `schema.md`, and `ingest.md` MAY remain published pages, but they MUST NOT lead the rendered top-level nav ahead of the wiki entry points.
- The published `Sources` navigation MUST stay grouped by year and default to a compact initial state:
  - the top-level `Sources` section MUST start expanded on first render.
  - generated yearly groups inside `Sources` MUST start expanded on first render.
  - individual source-note entries inside each year group MUST remain collapsed until opened through normal navigation.
  - any KB-specific navigation script used for this default state MUST only adjust the initial render and MUST NOT override later manual expand/collapse actions.
- Frontmatter relationships that drive human navigation MUST also be visible in rendered page content:
  - wiki page `sources:` entries MUST render as a visible `## Sources` section
  - source-note `related_pages:` entries MUST render as a visible `## Related pages` section
- Pull requests that change the knowledge base publishing inputs MUST run a strict MkDocs build in CI before merge.
- GitHub Pages MUST publish the rendered site from GitHub Actions.
- Local build output MUST be ignored by git.

## Seed Content

The initial knowledge base SHOULD include:
- the LLM Knowledge Base references that motivated the feature
- at least one deployment-related source note
- at least one Phaser/game-development source note
- starter wiki pages for topics, tools, and workspace projects
