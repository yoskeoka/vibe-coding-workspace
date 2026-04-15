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
- The active owner target for that board MUST be stored explicitly in `.local/pj/config.json`.
- A single local workspace operates against one owner scope at a time. Switching from a personal board to an organization board later requires an explicit configuration change, not a different one-off flag on a later command.
- `.local/pj/` is the only supported local workspace-triage state for the current workflow.
- The workspace MUST NOT depend on committed legacy tracker runtime artifacts or local database state for current task coordination.
- The local cache is derived data only. Deleting it must not lose task state.
- `docs/exec-plan/todo/` remains the canonical tracker for implementation plans once a task is selected.

### 3. Minimum project fields
The workflow depends on these GitHub Project fields:
- `Status`
- `Workspace Repo` (exposed by the CLI as `repo`)
- `Kind`
- `Priority`

`Status`, `Workspace Repo`, `Kind`, and `Priority` MUST be available as single-select fields on the canonical workspace board.
`pj init` MUST provision the custom `Workspace Repo`, `Kind`, and `Priority` fields when they are missing, using the workspace's canonical option sets for the spike.
Provisioning MUST be idempotent for an already-compatible board; later `pj init` runs must reuse existing compatible fields instead of creating duplicates.
The canonical `Workspace Repo` options for the spike are `vibe-coding-workspace`, `ww`, `ai-arena`, `reversi-adventure`, `vim-learning-game`, and `envdiff`.
The canonical `Kind` options for the spike are `Feature`, `Bug`, `Chore`, and `Research`.
The canonical `Priority` options for the spike are `High`, `Medium`, and `Low`.
GitHub currently rejects `Repo` as a custom ProjectV2 field name, so the workflow MUST use `Workspace Repo` as the remote field name while keeping `--repo` as the CLI/operator-facing flag.
If one of those required field names already exists with the wrong field type or without the required canonical options, `pj init` MUST fail with a compatibility error instead of rewriting the field automatically.
No other custom fields are required for the workflow to function.

### 4. Local cache
- The CLI MUST store a structured cache under `.local/pj/`.
- The CLI MUST store owner-scope configuration under `.local/pj/config.json`.
- The primary cache artifact MUST be JSON so agents can read it without re-querying GitHub.
- The owner config MUST include `owner` and `owner_type`.
- No committed fallback cache, archive, or mirror under legacy tracker runtime paths is part of the supported workflow contract.
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
- Provisions the minimum custom workflow fields when they are missing
- Writes the active owner target to `.local/pj/config.json`
- Writes the resolved project identity into the local cache so later commands can reuse it
- Fails clearly if the minimum workflow fields cannot be provisioned or are still incompatible after bootstrap

#### `sync`
- Authenticates via `gh auth token`
- Queries the configured GitHub Project through the GraphQL API
- Refreshes the local cache from remote state
- Reuses the configured owner target after a successful `init`
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

#### `config`
- Prints, sets, or clears the active owner target
- Must be the explicit mechanism for switching the local workspace from one owner scope to another
- Must clear incompatible cached project identity when the owner scope changes
- `clear` must also remove the cached project snapshot so later commands cannot keep operating on the old board implicitly

### 6. `triage-tasks` skill integration
The `triage-tasks` skill may use the local cache and/or the CLI output as its workspace-task source instead of historical local-priority-file flows.
- The skill MUST bootstrap or refresh the cache with `pj init` and/or `pj sync` before relying on local task state.
- The skill MUST treat `pj list` output plus the cached `Priority` field as the day-to-day "what next?" view; the current spike does not provide a separate `ready` command.
- The canonical Project used by this workflow MUST include a `Priority` field; fallback ranking applies when an item's `Priority` value is empty, unset, or otherwise unknown, not when the field is absent from the Project schema.
- The skill MUST create new triage items with `pj add` and claim or complete them by changing `Status` with `pj move`.
- After `pj init` or `pj sync`, the skill SHOULD include the canonical GitHub Project URL in the briefing when the owner scope and project number are known.
- When an item's `Priority` value is empty, unset, incomplete, or displayed as unknown (for example `-` in `pj list`), the skill MUST still rank a small shortlist using explicit heuristics such as: active execution plans over vague future ideas, broken/failing workflow items over aspirational enhancements, and tasks in the currently active repo over distant backlog items.
- The skill SHOULD present the top-priority shortlist before dumping the full board so the user can choose quickly, while still making the Project URL or full list available.
- After presenting the current task list, the skill MUST offer the next step as explicit numbered choices:
  `1. Pick a task`, `2. Update the list`, `3. Full re-triage`.
- After offering those choices, the skill MUST wait for user confirmation instead of auto-choosing one of them.
- During full re-triage, the skill MAY delegate repo-by-repo read-only exploration to subagents in order to keep the main context smaller.
- When delegating that exploration, the skill SHOULD prefer an available low-cost small model rather than a frontier-sized model, but MUST avoid hard-coding a specific model version name in the workflow contract.
- The main agent MUST keep responsibility for final prioritization, Project mutations, and the handoff prompt even when read-only exploration is delegated.
- After the user picks a task, the default handoff SHOULD be a fresh-session prompt rather than immediately starting implementation in the same session, because triage often leaves broad cross-repo context in the conversation.
- That handoff prompt SHOULD include enough context to start the next workflow step cleanly: target repo path, suggested `ww` worktree command, files to read first, the goal, deliverables, and key constraints.
- For normal planning/execution handoff, that suggested command SHOULD use the globally installed `ww` binary rather than raw `git switch -c`.
- The handoff contract SHOULD support both startup contexts:
  - from the target repo root: `ww create <type>/<name>` then `cd "$(ww cd <type>/<name>)"`
  - from the workspace root targeting a child repo: `ww create --repo <repo> <type>/<name>` then `cd "$(ww cd --repo <repo> <type>/<name>)"`
- That handoff prompt SHOULD explicitly name the skill to use in the next session when the next workflow step is clear.
- If the selected task is non-trivial and does not already have an execution plan, the handoff prompt SHOULD direct the next session to use `plan-execution`.
- If the selected task already has an execution plan and the next step is implementation, the handoff prompt SHOULD direct the next session to use `execute-task`.
- The skill MUST emit the handoff prompt in the same language the user is currently using in the chat, rather than always generating multiple language variants.
- The skill MUST NOT instruct agents to use legacy tracker-only concepts such as backup restore flows, dependency edges, or claim-style mutations outside the GitHub Project workflow.
- Because the current spike only mutates `Status` on existing items, the skill SHOULD treat edits to `Repo`, `Kind`, or `Priority` as a manual GitHub Project follow-up unless the item is recreated.

### 7. Placement and structure
- The spike CLI lives under `tools/pj/` as an independent Go module.
- The binary entrypoint lives under `tools/pj/cmd/pj/`.
- Internal packages MAY live under `tools/pj/internal/`.
- The placement is intentionally workspace-local and utility-oriented; extraction to a standalone repository remains a later decision.

## Non-Goals
- Replacing GitHub Issues, child-repo issue trackers, or `docs/exec-plan/todo/`
- Implementing dependency graphs, recursive task trees, or database-style local history
- Publishing the CLI as a reusable public tool in this spike
- Full automation of task selection without user confirmation
