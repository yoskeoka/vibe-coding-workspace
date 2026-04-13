# Spec: Triage Tasks — GitHub Projects Workflow

## Goal
Provide a low-friction workspace triage flow where GitHub Projects is the canonical remote task board and a local CLI keeps an AI-friendly structured cache for fast reads and session handoff.

## Scope
- Workspace-level task tracking for this meta-repo
- Synchronization between a single GitHub Project and a local derived cache
- Minimal commands needed to triage and update tasks
- Integration points for the `triage-tasks` skill and `AGENTS.md`

## Requirements

### 1. Session start hook (AGENTS.md)
A short entry in `AGENTS.md` MUST propose triage at the start of every new session:
- "At the start of a new session, if the user has not given you a specific task yet, suggest running `triage-tasks` to review priorities across all managed projects."
- Keep to 1–2 lines. Details live in the skill.

### 2. Canonical state
- The canonical remote source of workspace triage state is a GitHub Project (ProjectV2).
- The local cache is derived data only. Deleting it must not lose task state.
- `docs/exec-plan/todo/` remains the canonical tracker for implementation plans once a task is selected.

### 3. Minimum project fields
The workflow depends on these GitHub Project fields:
- `Status`
- `Repo`
- `Kind`
- `Priority`

`Status` MUST be a single-select field. `Repo`, `Kind`, and `Priority` MAY be single-select fields for the initial spike; no other custom fields are required for the workflow to function.

### 4. Local cache
- The CLI MUST store a structured cache under `.local/pj/`.
- The primary cache artifact MUST be JSON so agents can read it without re-querying GitHub.
- The cache MUST include:
  - sync timestamp
  - project identity (`owner`, `owner_type`, `project_number`, `project_id`)
  - field metadata needed to map names to GitHub IDs
  - the current item list with normalized field values
- The cache directory is non-Git and must be ignored by `.gitignore`.

### 5. Workspace CLI
A workspace-local Go CLI provides the task operations. The initial command set is:

#### `sync`
- Authenticates via `gh auth token`
- Queries the configured GitHub Project through the GraphQL API
- Refreshes the local cache from remote state

#### `list`
- Reads the local cache and renders the current task list
- MAY support simple filtering, but full triage scoring is not required in the spike

#### `add`
- Creates a draft project item in GitHub Projects
- Supports setting the minimum field set when values are provided
- Refreshes the local cache after mutation succeeds

#### `move`
- Updates the `Status` field for an existing project item
- Refreshes the local cache after mutation succeeds

### 6. `triage-tasks` skill integration
The `triage-tasks` skill may use the local cache and/or the CLI output as its workspace-task source instead of `bd` or `.local/priority.md`.

### 7. Placement and structure
- The spike CLI lives under `tools/pj/` as an independent Go module.
- The binary entrypoint lives under `tools/pj/cmd/pj/`.
- Internal packages MAY live under `tools/pj/internal/`.
- The placement is intentionally workspace-local and utility-oriented; extraction to a standalone repository remains a later decision.

## Non-Goals
- Replacing GitHub Issues, child-repo issue trackers, or `docs/exec-plan/todo/`
- Implementing dependency graphs, recursive task trees, or Dolt-like history
- Publishing the CLI as a reusable public tool in this spike
- Full automation of task selection without user confirmation
