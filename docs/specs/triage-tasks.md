# Triage tasks in a GitHub Projects workflow

## Goal
Provide a low-friction workspace triage flow where GitHub Projects is the canonical remote task board and a local CLI keeps an AI-friendly structured cache for fast reads and session handoff.

## Scope
- Workspace-level task tracking for this meta-repo
- Synchronization between a single GitHub Project and a local derived cache
- Minimal commands needed to triage, update, and navigate tasks
- Integration points for the `triage-tasks` skill and `AGENTS.md`

## Requirements

### 1. Session start hook (AGENTS.md)
A short entry in `AGENTS.md` MUST propose triage at the start of every new session:
- "At the start of a new session, if the user has not given you a specific task yet, suggest running `triage-tasks` to review priorities across all managed projects."
- Keep to 1–2 lines. Details live in the skill.

### 2. Canonical state
- The canonical remote source of workspace triage state is a GitHub Project (ProjectV2).
- The workspace MUST use a dedicated ProjectV2 named `Workspace Task Triage` for cross-project task coordination.
- The canonical ProjectV2 has one owner scope: `user` or `org`.
- The canonical ProjectV2 SHOULD be linked to this workspace repository so it appears in the repository's Projects tab for discoverability.
- GitHub only supports linking a ProjectV2 to repositories owned by the same user or organization as the Project owner; if the workspace later moves to an organization-owned board, the linked workspace repository must be owned by that same organization.
- This board is reserved for workspace triage data; unrelated personal/work boards MUST NOT be reused as the canonical workspace tracker.
- Workspace triage MUST begin with an explicit bootstrap step, `pj init --owner <owner> --owner-type user|org`, which resolves the canonical board by name and creates it when absent.
- If `Workspace Task Triage` does not exist yet, `pj init` MUST create it before later `pj` commands manage items on it.
- The active owner target for that board MUST be stored explicitly in `.local/pj/config.json`.
- A local workspace uses one owner scope at a time. Switch boards through an explicit configuration change.
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

`pj init` MUST provision the custom `Workspace Repo`, `Kind`, and `Priority` fields when they are missing, deriving the `Workspace Repo` display values from `setup.sh` plus the workspace repository itself and using the workspace's canonical non-repo option sets for the spike.

Provisioning MUST be idempotent for an already-compatible board; later `pj init` runs must reuse existing compatible fields instead of creating duplicates.

The canonical `Workspace Repo` option display values for the current workspace are the basenames of those setup-derived repositories, including `vibe-coding-workspace`, `ww`, `ai-arena`, `reversi-adventure`, `vim-learning-game`, and `envdiff`.

The local cache MUST include ordered enriched repo metadata so operators can resolve `--repo` by display basename, `owner/repo`, `github.com/owner/repo`, unique prefix, or stable integer index without depending on GitHub's display labels alone.

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
- Supports both `--body` and `--body-file`, rejecting both together
- Refreshes the local cache after mutation succeeds

#### `update`
- Updates title, body, `Status`, `Repo`, `Kind`, and `Priority` for an existing project item when values are provided
- Supports both `--body` and `--body-file`, rejecting both together
- Refreshes the local cache after mutation succeeds
- Replaces the former `move` command; status changes use `pj update --item <id> --status <value>`

#### `update-batch`
- Applies a JSON-file mutation plan to multiple existing project items
- Validates the full input before remote mutation when possible
- Refreshes the local cache once after all successful mutations instead of after every item
- Stops on the first remote failure and leaves the local cache stale with a clear instruction to run `pj sync`
- Is an optimization for reconciling existing items; it does not replace `pj add` for genuinely missing tasks

#### `url`
- Prints the canonical GitHub Project URL from cached project metadata

#### `open`
- Opens the canonical GitHub Project URL from cached project metadata in the operator's browser

#### `config`
- Prints, sets, or clears the active owner target
- Must be the explicit mechanism for switching the local workspace from one owner scope to another
- Must clear incompatible cached project identity when the owner scope changes
- `clear` must also remove the cached project snapshot so later commands cannot keep operating on the old board implicitly

#### `repo-link`
- Checks, adds, or removes the repository-level ProjectV2 link for the cached canonical board.
- Supports `status`, `add`, and `remove` subcommands with an explicit `<owner>/<repo>` target.
- Uses `pj repo-link status <owner>/<repo>` to report whether the target repository already exposes the canonical board through its Projects tab.
- Uses `pj repo-link add <owner>/<repo>` to link the canonical board to the target repository.
- Uses `pj repo-link remove <owner>/<repo>` to unlink the canonical board from the target repository when needed.
- Must reject target repositories whose owner does not match the configured Project owner.
- Setting a default repository for the Project is not required by this workspace triage flow because task creation uses Project draft items; repository linking is required only for repository-tab discoverability.

### 6. `triage-tasks` skill integration
The `triage-tasks` skill may use the local cache and/or the CLI output as its workspace-task source instead of historical local-priority-file flows.
- The skill MUST bootstrap or refresh the cache with `pj init` and/or `pj sync` before relying on local task state.
- The skill MUST treat `pj list` output plus the cached `Priority` field as the day-to-day "what next?" view; the current spike does not provide a separate `ready` command.
- The workflow Project MUST include a `Priority` field. Fallback ranking handles an empty, unset, or unknown item value.
- The skill MUST create new triage items with `pj add` and claim or complete them by changing `Status` with `pj update --item <id> --status <value>`.
- The skill SHOULD use `pj add --body-file` for generated startup handoff bodies that are long enough to make inline shell quoting awkward.
- The skill MAY correct existing items through `pj update` instead of forcing manual GitHub edits or item recreation for `Repo`, `Kind`, and `Priority`.
- During full re-triage, the skill MAY build a short mutation plan first and apply existing-item corrections with `pj update-batch` when that command is available.
- The skill MUST continue using `pj add` for genuinely missing items; `pj update-batch` only changes items that already exist.
- If `pj update-batch` is unavailable or fails, the skill MUST remain compatible with one-item-at-a-time `pj update`.
- During full re-triage, every new `pj add` item MUST include a compact, remote-facing startup handoff.
- The Project item body MUST stay concise enough for GitHub Project scanning and SHOULD be a short Markdown block with these minimum fields:
  - `Source`: the local plan, local issue, GitHub PR, GitHub Issue, or discovered source reference
  - `Repo`: the target workspace repo
  - `Next`: the recommended next-step skill when clear, such as `plan-execution` or `execute-task`; otherwise `Manual triage`
  - `Start`: the suggested `ww create` and `ww cd` command when a concrete planning or execution branch is meaningful
  - `Read`: initial files, docs, PRs, or issues to inspect first
  - `Goal`: a one-sentence outcome for the task
- Stored Project item bodies MUST use English as the stable workspace-board language so the remote board remains consistent across sessions. Chat handoff prompts remain user-facing output and MUST follow the current session language rule below.
- For execution-plan items, the body SHOULD recommend `execute-task`, include the plan path, and suggest a `feat/<plan-name>` or `fix/<plan-name>` startup command that matches the plan's expected branch type.
- For local issue follow-up items, the body SHOULD recommend `plan-execution` when the issue is non-trivial, include the issue path, and suggest a `plan/<issue-name>` startup command.
- For open PR review or follow-up items, the body SHOULD avoid a misleading implementation startup when the action is review, approval, or post-review response; it SHOULD use `Manual triage` or a specific review skill only when clear, include the PR URL, and name the expected review/follow-up goal.
- For open GitHub Issue items, the body SHOULD use `Manual triage` unless the issue clearly maps to planning or execution, include the issue URL, and avoid inventing local plan paths that do not exist yet.
- When a concrete startup prompt is unavailable, the body MUST include source, repo, initial reading context, and a short goal. It SHOULD set `Next: Manual triage` and `Start: Not yet specified`.
- After `pj init` or `pj sync`, the skill SHOULD include the canonical GitHub Project URL in the briefing by using `pj url` when the cache is available.
- When an item's `Priority` value is empty, unset, incomplete, or displayed as unknown (for example `-` in `pj list`), the skill MUST still rank a small shortlist using explicit heuristics such as: active execution plans over vague future ideas, broken/failing workflow items over aspirational enhancements, and tasks in the currently active repo over distant backlog items.
- The skill SHOULD present the top-priority shortlist before dumping the full board so the user can choose quickly, while still making the Project URL from `pj url` or full list available.
- After presenting the current task list, the skill MUST offer the next step as explicit numbered choices:
  `1. Pick a task`, `2. Update the list`, `3. Full re-triage`.
- After offering those choices, the skill MUST wait for user confirmation instead of auto-choosing one of them.
- During full re-triage, the skill MAY delegate repo-by-repo read-only exploration to subagents to keep the main context smaller.
- When delegating that exploration, the skill SHOULD explicitly choose an available low-cost small model for routine read-only collection unless the repo question needs deeper reasoning, but MUST avoid hard-coding a specific model version name in the workflow contract.
- The main agent MUST keep responsibility for final prioritization, Project mutations, and the handoff prompt even when read-only exploration is delegated.
- Full re-triage MUST classify collected findings before mutating the Project:
  - direct backlog item: pending exec plan, local issue, open PR, or open GitHub Issue
  - roadmap gap: an unchecked project-plan phase or unmet requirement that is not already represented by a more concrete item
  - do not add yet: vague, duplicate, or superseded findings that are better represented by an existing direct item
- Full re-triage MUST reconcile existing Project items before creating new ones.
- Reconciliation MUST compare a candidate against existing items by `Source` first, then by a normalized title fallback when `Source` is missing, stale, or still in an older body format.
- Reconciliation MUST prefer updating an equivalent existing item over creating a duplicate item.
- Full re-triage MUST only create a new Project item when no equivalent existing item can be reconciled.
- Reconciliation MUST handle local sources in all of these states:
  - still present in `docs/exec-plan/todo/` or `docs/issues/`
  - deleted by a completed execution and retrievable through its PR or Git
    history
  - missing locally because the repo checkout is unavailable
  - missing locally because the file no longer exists
- Reconciliation MUST handle remote sources in all of these states:
  - PR URL still open
  - PR URL closed or merged
  - PR URL superseded by a newer open PR for the same workflow-sync task in the same repo
  - GitHub Issue still open
  - GitHub Issue closed
- Missing local repo checkouts MUST be reported explicitly during full re-triage.
- When a local checkout is missing, the skill MAY still inspect GitHub PRs and issues for that repo, but MUST mark local project-plan, exec-plan, and issue inspection as unavailable instead of inferring local state.
- When a local checkout is missing, delegated summaries MUST still emit the standard schema shape, using empty arrays for unavailable local `project_plan_gaps`, `exec_plans`, and `issues`, plus a caveat that only GitHub PR/issue inspection was possible.
- During full re-triage, a workflow-sync PR MUST be identified primarily by the GitHub label `workflow-sync`.
- If label data is unavailable for a candidate PR or existing Project item, the skill MAY fall back to the normalized title prefix `chore: update workflow skills to `, but the label remains the canonical identifier when present.
- During full re-triage, workflow-sync PRs MUST maintain at most one active Project item per repo for the current sync PR.
- When a newer open workflow-sync PR supersedes an older one in the same repo, the skill MUST update the existing item to the newer PR instead of creating another Project item.
- When an existing item is already being updated for title, source, repo, kind, priority, or status reasons, the skill SHOULD also normalize an old `Source:`-only or otherwise outdated body to the current durable handoff body format in the same mutation.
- Full re-triage SHOULD build a short mutation plan before changing the Project.
- That mutation plan SHOULD be ordered as:
  - mark stale items `Done`
  - update PR source/title/body replacements
  - normalize `Priority`, `Kind`, `Repo`, and body fields on existing items
  - add genuinely missing items
  - run a final `pj sync`/`pj list`
- During full re-triage, the skill MAY apply existing-item corrections with `pj update-batch` when that command is available.
- If `pj update-batch` is unavailable, the skill MUST remain compatible with one-item-at-a-time `pj update`.
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
- The skill SHOULD use `pj update` for corrections to `Repo`, `Kind`, or `Priority` on existing items.
- During full re-triage, delegated repo summaries SHOULD use a strict schema with stable field names, workspace repo basenames, canonical priority casing (`High`, `Medium`, `Low`), and a stable `next` string rather than free-form shapes.
- The final full re-triage briefing SHOULD include:
  - the Project URL
  - the synced item count
  - how many items were marked `Done`
  - how many existing items were updated
  - how many new items were added
  - a high-priority Todo shortlist
  - caveats such as missing local checkouts or GitHub-only inspection
  - the standard numbered next choices

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
