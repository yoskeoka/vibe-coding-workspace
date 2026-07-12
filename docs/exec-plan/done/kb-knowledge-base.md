# Knowledge Base Ingest and Publish Flow

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Add a repo-native knowledge base for the workspace so useful external references can be ingested by AI, compiled into durable wiki pages, and published as a browsable site for humans.

Addresses:
- Scattered bookmarks and open tabs that are hard to rediscover later
- Lack of a durable, AI-maintainable place for cross-project knowledge
- No human-friendly rendered view of accumulated workspace knowledge

## Code Changes

### `tools/kb`
- Add a dedicated helper command for knowledge-base maintenance.
- Support at least `build`, `serve`, and `check` operations.
- Keep build/publish concerns separate from execution workflow skills.

### `skills/knowledge-base/SKILL.md`
- Add a workspace-only skill for knowledge-base ingest, query filing-back, and linting.
- Define how AI should turn URLs into source notes and wiki updates.

### `.github/workflows/kb-pages.yml`
- Publish the rendered knowledge base to GitHub Pages when relevant files change on `main`.

### `.gitignore`
- Ignore the local MkDocs build output.

## Spec Changes

### `docs/specs/knowledge-base.md`
- Define the directory layout and page conventions.
- Define the ingest, query, lint, and compile flows.
- Define how source notes relate to compiled wiki pages.
- Define the publishing path with MkDocs and GitHub Pages.

### `docs/kb/*`
- Seed the new knowledge base with schema, ingest instructions, index, log, and a few initial source notes and wiki pages.

## Design Decisions

- Record the storage-location decision as an indexed file under `docs/design-decisions/adr/`.
- Keep the knowledge base in-repo rather than using GitHub Wiki.

## Sub-tasks

- [ ] [parallel] Write the spec for the knowledge-base structure and operations
- [ ] [parallel] Add the knowledge-base schema, ingest guide, and seed wiki pages
- [ ] [parallel] Add the dedicated `knowledge-base` skill
- [ ] [parallel] Add the `tools/kb` helper command
- [ ] [depends on: spec, docs/kb, tools/kb] Add MkDocs config and GitHub Pages workflow
- [ ] [depends on: all above] Verify structure and rendering command expectations
