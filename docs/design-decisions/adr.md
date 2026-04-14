# Architectural Decision Records (ADR)

## 2026-04-07 Skip child-repo workflow PRs when `skills/` is unchanged

### Context
`sync-workflow-to-child-repos.yml` currently opens PRs in every child repository whenever `vibe-coding-workspace` receives a push to `main`. However, `setup-workspace.sh` only symlinks skills from the workflow submodule, so changes outside `skills/` do not affect child repos.

### Decision
The sync workflow will compare the child repo's recorded workflow submodule commit against the pushed workflow commit and only create a PR when the diff contains changes under `skills/`.

### Consequences
- Child repos stop receiving no-op workflow PRs for unrelated repo changes.
- The workflow becomes slightly more selective, so future changes that should sync to child repos must live under `skills/` or update the sync rule explicitly.

---

## 2026-04-11 Store the workspace knowledge base in-repo under `docs/kb/`

### Context
The workspace needs a durable place to collect useful references for AI-centered personal development: articles, demo links, tool notes, and cross-project patterns. GitHub Wiki was considered, but the workspace already treats git-tracked Markdown under `docs/` as the canonical AI context. The knowledge base also needs a human-readable publishing path.

### Decision
The knowledge base will live in this repository under `docs/kb/` as Markdown files. It will follow a three-layer structure inspired by Karpathy's "LLM Knowledge Base" pattern:

- `docs/kb/sources/` for immutable source notes
- `docs/kb/wiki/` for compiled wiki pages
- `docs/kb/schema.md` and `docs/kb/ingest.md` for maintenance rules

AI ingestion will use a dedicated `knowledge-base` skill instead of the execution workflow skills. Human-readable rendering will use MkDocs Material and publish to GitHub Pages from the same Markdown sources.

### Consequences
- The same files remain easy for agents to edit, grep, diff, review, and publish.
- Knowledge-base ingest stays conceptually separate from `execute-task` and `exec-plan`.
- The repository gains a second documentation surface (raw Markdown plus rendered Pages), so the structure and schema must remain disciplined.

---

## 2026-04-14 Dogfood the released global `ww` binary in normal workflow startup

### Context
The workspace workflow currently teaches raw git branch creation (`git fetch origin && git switch -c ...`) at most task entry points. That bypasses the product this workspace is already building for parallel multi-repo worktree flows: `ww`.

At the same time, `ww`'s current product direction already favors the patterns this workspace needs:
- workspace-aware repo targeting
- centralized `.worktrees/` layout in workspace mode
- explicit path-oriented shell interfaces such as `ww cd`
- git-native behavior instead of a parallel branch-management abstraction

If the workspace keeps using raw git for normal planning and execution, `ww` loses the best possible dogfood loop and branch/worktree contention remains a manual concern in the primary checkout.

### Decision
Use the globally installed released `ww` binary as the default operator path for normal workflow startup in the workspace and its child repos.

Concretely:
- planning and execution flows should start with `ww create ...`
- operators should enter worktrees with `ww cd` (or equivalent `ww`-based navigation)
- raw git branch creation is no longer the default documented path for ordinary work
- the latest in-repo `ww` build is reserved for work inside `ww/` that changes or verifies unreleased `ww` behavior
- `ww` failures discovered during normal workflow use are first-class outputs that must be recorded back to `ww`, not silent friction to step around

### Consequences
- The workspace continuously exercises the real released `ww` UX during day-to-day work.
- Parallel task execution becomes the default mental model instead of an optional advanced path.
- Workflow docs, skills, and handoff prompts must be updated together so they stop emitting raw git branch commands.
- Some temporary fallback paths still remain necessary when `ww` itself is the thing being debugged or when a documented `ww` failure blocks unrelated work.

---
