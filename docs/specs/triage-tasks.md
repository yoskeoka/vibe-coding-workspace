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
- The workspace MUST use a dedicated ProjectV2 named `Workspace Task Triage` for cross-project task coordination.
- This board is reserved for workspace triage data; unrelated personal/work boards MUST NOT be reused as the canonical workspace tracker.
- Workspace triage MUST begin with an explicit bootstrap step, `pj init --owner <owner> --owner-type user|org`, which resolves the canonical board by name and creates it when absent.
- If `Workspace Task Triage` does not exist yet, `pj init` MUST create it before later `pj` commands manage items on it.
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

#### `init`
- Authenticates via `gh auth token`
- Resolves the canonical `Workspace Task Triage` board by owner and title
- Creates the canonical board when absent
- Writes the resolved project identity into the local cache so later commands can reuse it
- Fails clearly if the minimum workflow fields are still missing after bootstrap

#### `sync`
- Authenticates via `gh auth token`
- Queries the configured GitHub Project through the GraphQL API
- Refreshes the local cache from remote state
- Reuses the cached project identity after a successful `init`

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
- The skill MUST bootstrap or refresh the cache with `pj init` and/or `pj sync` before relying on local task state.
- The skill MUST treat `pj list` output plus the cached `Priority` field as the day-to-day "what next?" view; the current spike does not provide a separate `ready` command.
- The skill MUST create new triage items with `pj add` and claim or complete them by changing `Status` with `pj move`.
- After `pj init` or `pj sync`, the skill SHOULD include the canonical GitHub Project URL in the briefing when the owner scope and project number are known.
- When `Priority` is missing or incomplete, the skill MUST still rank a small shortlist using explicit heuristics such as: active execution plans over vague future ideas, broken/failing workflow items over aspirational enhancements, and tasks in the currently active repo over distant backlog items.
- The skill SHOULD present the top-priority shortlist before dumping the full board so the user can choose quickly, while still making the Project URL or full list available.
- After presenting the current task list, the skill MUST offer the next step as explicit numbered choices:
  `1. Pick a task`, `2. Update the list`, `3. Full re-triage`.
- After offering those choices, the skill MUST wait for user confirmation instead of auto-choosing one of them.
- During full re-triage, the skill MAY delegate repo-by-repo read-only exploration to subagents in order to keep the main context smaller.
- When delegating that exploration, the skill SHOULD prefer an available low-cost small model rather than a frontier-sized model, but MUST avoid hard-coding a specific model version name in the workflow contract.
- The main agent MUST keep responsibility for final prioritization, Project mutations, and the handoff prompt even when read-only exploration is delegated.
- After the user picks a task, the default handoff SHOULD be a fresh-session prompt rather than immediately starting implementation in the same session, because triage often leaves broad cross-repo context in the conversation.
- That handoff prompt SHOULD include enough context to start the next workflow step cleanly: target repo path, suggested branch command, files to read first, the goal, deliverables, and key constraints.
- The skill MUST emit the handoff prompt in the same language the user is currently using in the chat, rather than always generating multiple language variants.
- The skill MUST NOT instruct agents to use beads-only concepts such as Dolt backup restore, dependency edges, or `bd update --claim`.
- Because the current spike only mutates `Status` on existing items, the skill SHOULD treat edits to `Repo`, `Kind`, or `Priority` as a manual GitHub Project follow-up unless the item is recreated.

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
