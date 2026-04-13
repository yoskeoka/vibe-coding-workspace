# Workspace Task Triage Project Bootstrap

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Extend the workspace-local `pj` spike so the GitHub Projects workflow can bootstrap its own canonical board instead of assuming it already exists. This completes the missing setup step for the workspace-wide task triage flow and reduces manual preconfiguration friction.

This plan supports the project-plan requirements that workspace task coordination uses a lightweight GitHub Projects-backed flow and that the local CLI is practical enough to evaluate for longer-term use.

## Background

The current `pj` spike can sync, list, add, and move items only after a dedicated ProjectV2 already exists. The specs now define that the workspace must use a dedicated board named `Workspace Task Triage`, but creation of that board is still a manual prerequisite.

To make the workflow self-bootstrapping, the CLI needs a way to ensure the canonical board exists and can be initialized with the expected field model before normal task operations begin.

## Spec Changes

### `docs/specs/github-projects-task-cli.md`

- Add a board-bootstrap command and behavior for the canonical `Workspace Task Triage` ProjectV2
- Define how the CLI resolves an existing board by name before creating a new one
- Define the initial setup responsibilities after creation:
  - create or verify the dedicated ProjectV2
  - persist its identity into the local cache
  - ensure the minimum workflow fields exist or fail clearly if the API cannot provision them in this spike
- Clarify whether bootstrap is operator-invoked (`pj init`, `pj ensure-project`, etc.) or folded into `sync`

### `docs/specs/triage-tasks.md`

- Document that workspace triage begins by ensuring the canonical `Workspace Task Triage` board exists
- Clarify whether board creation is automatic or an explicit bootstrap step
- Keep the single-board workspace model explicit so the skill and CLI share the same assumption

### `AGENTS.md`

- Update the workspace task tracking section to reference the bootstrap flow once the command name is finalized

## Code Changes

### `tools/pj/`

- Add a board-bootstrap path to the CLI, likely as a dedicated command rather than overloading `sync`
- Implement GitHub GraphQL queries/mutations to:
  - find an existing ProjectV2 by owner and name
  - create the canonical board when absent
  - record the resulting project identity in cache
- Decide whether minimum field provisioning belongs in the same command or in a follow-up task

### Cache and configuration

- Store enough metadata to remember that the canonical board has been resolved or created
- Keep owner / owner_type configuration explicit so a later org-scoped deployment is still possible

## Design Decisions

- The canonical board name is fixed to `Workspace Task Triage`
- This repository remains a workspace utility repo; keep the CLI under `tools/pj/`
- Prefer a dedicated bootstrap command over hidden side effects in `sync` unless GitHub API limitations make explicit bootstrap awkward
- If GitHub's Projects API makes custom-field creation too heavy for this task, split field provisioning into a follow-up plan instead of forcing a partial implicit workflow

## Sub-tasks

- [ ] Update `docs/specs/github-projects-task-cli.md` with the canonical board bootstrap behavior and command contract
- [ ] Update `docs/specs/triage-tasks.md` and `AGENTS.md` to describe the bootstrap step consistently
- [ ] [parallel] Investigate the GraphQL path for discovering ProjectV2 boards by owner and name
- [ ] [parallel] Investigate the GraphQL path for creating a new ProjectV2 under a user or org owner
- [ ] [depends on: GraphQL discovery, GraphQL creation] Implement the CLI bootstrap command in `tools/pj/`
- [ ] [depends on: CLI bootstrap command] Decide whether minimum field provisioning is in scope now or must become a follow-up issue/plan
- [ ] [depends on: CLI bootstrap command] Verify that the canonical board can be created or resolved using `gh auth token`-derived credentials

## Verification

- Confirm the CLI can detect an existing `Workspace Task Triage` ProjectV2 for the configured owner
- Confirm the CLI can create the board when it does not already exist
- Confirm the resulting board identity is written to `.local/pj/cache.json`
- Confirm post-bootstrap commands (`sync`, then `add`/`list`) can target the resolved board without manual project-number lookup
- Confirm docs/specs no longer assume a manually pre-created board

## Expected Outcome

- The workspace task flow has an explicit bootstrap step for the canonical ProjectV2
- The operator no longer needs to manually create the triage board before using `pj`
- The GitHub Projects-based workspace workflow becomes closer to a repeatable end-to-end setup
