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

Ingest SHOULD:
- preserve the user's framing about why the source matters
- prefer updating an existing wiki page over creating duplicates
- identify which workspace projects, tools, topics, and patterns are affected

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

Human-readable publishing MUST use MkDocs with a dedicated config file `mkdocs.kb.yml`.

The local helper command `tools/kb` MUST provide:
- `build` to render the site
- `serve` to preview the site locally
- `check` to validate required files and directories

`build` SHOULD run in strict mode so broken links and missing pages fail fast.
The helper SHOULD prefer `uv` when available and fall back to `python3 -m mkdocs` otherwise.

## Rendering and Publishing

- The rendered site MUST be generated from the same Markdown files stored in git.
- Pull requests that change the knowledge base publishing inputs MUST run a strict MkDocs build in CI before merge.
- GitHub Pages MUST publish the rendered site from GitHub Actions.
- Local build output MUST be ignored by git.

## Seed Content

The initial knowledge base SHOULD include:
- the LLM Knowledge Base references that motivated the feature
- at least one deployment-related source note
- at least one Phaser/game-development source note
- starter wiki pages for topics, tools, and workspace projects
