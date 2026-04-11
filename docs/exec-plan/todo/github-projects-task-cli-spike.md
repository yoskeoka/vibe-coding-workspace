# GitHub Projects Task CLI Spike

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Replace the current beads-based workspace task management flow with a lighter GitHub Projects-based approach that keeps AI-friendly structured data locally while using GitHub as the remote source of truth.

The immediate goal is a workspace-local spike, not a publishable product. The CLI should prove that the workspace can triage, list, add, and move tasks without Dolt, beads hooks, or a separate issue database.

## Background

The current workspace uses beads (`bd`) for cross-project task tracking. In practice this has created avoidable overhead:

- Dolt server and remote configuration add operational complexity for a single-user hobby workspace
- `bd dolt push` is still referenced in workspace instructions even though no beads remote is configured
- `.beads/` is much heavier than the current task-tracking needs

Recent evaluation compared three directions:

1. keep beads and simplify the workflow
2. invent a lightweight local structured task tracker
3. use GitHub Projects as the remote system and keep local structured cache/views for AI workflows

The chosen direction is option 3. GitHub Projects already provides the human-facing board and remote persistence, while a small Go CLI can provide an AI-friendly local interface and later serve as the basis for a reusable personal tool.

## Spec Changes

### `docs/specs/triage-tasks.md`

- Replace the volatile `.local/priority.md` / beads-centered assumptions with a GitHub Projects-centered workflow
- Define GitHub Projects as the canonical remote source for workspace triage state
- Define the local cache/view artifact used by the CLI for AI-friendly access
- Document the minimum project fields used by the workflow:
  - `Status`
  - `Repo`
  - `Kind`
  - `Priority`
- Define the initial operational commands supported by the CLI:
  - `sync`
  - `list`
  - `add`
  - `move`

### `AGENTS.md`

- Remove or revise the beads-specific requirement that every session must run `bd dolt push`
- Update task-tracking guidance so workspace-level triage refers to the GitHub Projects flow instead of beads commands

## Code Changes

### New workspace-local Go CLI

- Add a small Go command under the workspace (path to be chosen during execution) for GitHub Projects task operations
- Implement GitHub authentication by reading a token from `gh auth token`
- Use GitHub Projects GraphQL APIs rather than shelling out to `gh project` for item operations
- Keep the implementation scoped to a workspace-local spike, with thin internal boundaries so a later extraction remains possible

### Local cache / config

- Store a local structured cache for project items to support AI-friendly reads and fast local listing
- Keep GitHub Projects as the source of truth; the cache is derived data
- Cache project metadata needed to map field names and single-select options to GitHub IDs

### Workflow cleanup

- Remove workspace instructions and helper flows that assume beads remains the task tracker
- Decide whether existing `.beads/` artifacts should be deleted immediately, ignored operationally, or cleaned up in a follow-up task

## Design Decisions

- The spike remains inside this workspace for now; do not create a new repository in this plan
- The first implementation should optimize for low-friction local use, not public distribution polish
- `gh auth` remains the operator-facing authentication mechanism even though the CLI uses direct GitHub API calls internally
- GitHub Projects is the only remote backend in scope for the spike
- Do not add dependency graphs, task trees, or Dolt-like local history in this plan

## Sub-tasks

- [ ] Update `docs/specs/triage-tasks.md` for the GitHub Projects-based workflow and local cache model
- [ ] Update `AGENTS.md` to remove beads-specific remote-sync requirements and align task tracking guidance with the new workflow
- [ ] Add a focused execution spec for the CLI structure, cache location, and required project field mapping
- [ ] Implement the initial Go CLI with:
  - token acquisition via `gh auth token`
  - project metadata lookup
  - item listing
  - draft item creation
  - status movement through single-select field updates
- [ ] Add a sync command that refreshes the local cache from GitHub Projects
- [ ] Verify the CLI against the workspace GitHub Project using real `gh` authentication
- [ ] Decide whether `.beads/` retirement belongs in the same execution PR or a follow-up cleanup plan
- [ ] Record standalone-tool graduation criteria so a later project can be evaluated explicitly

## Verification

- Confirm the workspace GitHub Project can be queried and updated through the CLI using `gh`-managed authentication
- Confirm the local cache can be regenerated from GitHub Projects without manual edits
- Confirm the minimal field set (`Status`, `Repo`, `Kind`, `Priority`) is enough to support workspace triage
- Confirm updated workflow docs no longer instruct agents to use beads-specific remote sync steps for workspace task tracking

## Expected Outcome

- Workspace task tracking no longer depends on Dolt-backed beads workflows
- A small Go CLI supports the core operations needed for workspace triage
- The workspace gains a realistic spike that can later be evaluated for extraction into a standalone personal/work task tool
