# Spec: Dogfood Global `ww` in the Workspace Workflow

## Goal
Make the workspace workflow use the globally installed `ww` CLI as the default operator path for branch/worktree creation and task startup so the workspace continuously dogfoods the real released tool.

## Scope
- Workspace-level workflow docs and skill contracts
- Day-to-day planning and execution entry points for the workspace repo and child repos
- Failure handling and issue filing when `ww` behaves unexpectedly during normal use

## Requirements

### 1. Default operator path
- Normal workflow entry points MUST use the globally installed `ww` binary, not `git switch -c`, to start plan or execution work.
- From the target repo root, the default startup flow is:
  - `ww create <type>/<name>`
  - `cd "$(ww cd <type>/<name>)"`
- From the workspace root when targeting a child repo, the default startup flow is:
  - `ww create --repo <repo> <type>/<name>`
  - `cd "$(ww cd --repo <repo> <type>/<name>)"`
- Workflow docs and handoff prompts MUST prefer the stable released global binary. They MUST NOT default to `go run ./cmd/ww`, a repo-local dev build, or raw git branch creation.

### 2. Workflow touchpoints
The workflow must eventually describe `ww` usage consistently anywhere it currently creates a branch or starts task execution.

This planning change defines the required operator-facing touchpoints to align now:
- `AI_WORKFLOW.md` branch setup plus Step 1/2/3 examples
- `AGENTS.md` branch rules and "start new feature / fix a bug" recipes
- `docs/specs/triage-tasks.md` handoff prompt contract
- `README.md` operator-facing workflow summary

The following skill contracts are part of the migration scope but remain pending execution-plan work in this repository:
- `skills/plan-execution/SKILL.md`
- `skills/execute-task/SKILL.md`
- `skills/triage-tasks/SKILL.md`

### 3. Parallel task operator experience
- Each active task SHOULD get its own `ww` worktree. Operators SHOULD avoid reusing the primary checkout as a task branch sandbox.
- `ww list` is the canonical view of active worktrees across the workspace.
- `ww cd` or `ww i` SHOULD be used to navigate existing worktrees instead of manual path guessing.
- If `ww create` reports that the branch or target path already exists, the operator MUST treat that as a branch/worktree contention signal:
  - inspect with `ww list`
  - resume the existing worktree if it is the intended task
  - otherwise create a new branch name rather than repurposing the existing one silently
- Cleanup flows SHOULD prefer `ww list --cleanable` and `ww clean` so post-merge cleanup also dogfoods `ww`.

### 4. Boundary between stable dogfooding and `ww` development
- For normal workspace and child-repo work, use the globally installed released `ww`.
- When the task is to change `ww` itself, the released global `ww` remains the default tool for creating the worktree that will hold the `ww` change.
- Use the latest in-repo `ww` build only when:
  - developing or testing unreleased `ww` behavior inside `ww/`
  - reproducing, debugging, or verifying a `ww` bug or regression
- Other workflow docs MUST NOT require unreleased `ww` behavior to start ordinary planning or execution in non-`ww` repos.

### 5. Failure handling
- If `ww` is missing, outdated, or fails before a task can start, treat that as a workflow problem to resolve explicitly, not as incidental friction to ignore.
- Stop and capture the problem when `ww` blocks required worktree lifecycle operations such as create, navigate, remove, or clean and there is no already-documented safe workaround.
- Raw git fallback is allowed only when one of these is true:
  - the task is explicitly about debugging `ww`
  - continuing without fallback would strand unrelated high-priority work
  - the failure is already captured and the fallback is documented in the handoff, plan, or issue
- When raw git fallback is used, the workflow MUST still record that `ww` failed and why the fallback was necessary.

### 6. `ww` issue filing policy
- Unexpected `ww` bugs, workflow friction, or confusing behavior discovered during normal workflow use are first-class outputs and SHOULD be filed back to `ww`.
- The filing target may be a GitHub issue, a `ww` execution plan, or a `ww` local issue file, depending on the next action, but the finding MUST be recorded somewhere durable.
- A recorded `ww` issue MUST include:
  - `ww` version and install path if known
  - exact command
  - current working directory and target repo
  - expected behavior
  - actual behavior
  - relevant stderr/stdout
  - whether a raw git fallback was used
  - impact on the blocked workspace task

## Non-Goals
- Replacing git for all low-level operations inside `ww`
- Requiring unreleased `ww` features before the workflow can function
- Eliminating every emergency fallback path on day one
