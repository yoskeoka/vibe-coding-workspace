> **Execution**: Use `/execute-task` to implement this plan.

# Adopt Global `ww` into Workspace Workflow for Dogfooding

## Objective

Change the workspace development workflow so normal plan creation and task execution use the globally installed `ww` CLI by default instead of raw git branch setup or the latest in-repo `ww` build.

This plan supports the workspace project goal that the workflow itself is a first-class product and turns the workspace into a continuous dogfooding loop for `ww`.

## Background

Current workspace docs still tell operators to start work with `git fetch origin && git switch -c ...`. That undermines three things the `ww` product is already designed to support:

- parallel task execution through dedicated worktrees
- workspace-aware repo targeting across child repositories
- shell-friendly navigation via explicit path output

The relevant `ww` design history already points in the same direction:

- `ww` core belief: AI-first, spec-first workflow support
- `ww` ADR: centralized `.worktrees/` in workspace mode
- `ww` ADR: bounded nearest-workspace detection
- `ww` ADR: explicit path-oriented shell integration via `ww cd`

The workspace should apply that reasoning directly instead of continuing to document raw git as the normal operator path.

## Workflow Touchpoints to Migrate

These touchpoints currently create branches or start task execution and must be migrated to `ww`:

- `AI_WORKFLOW.md`
  - branch setup for every workflow step
  - Step 1/2/3 examples that currently imply raw branch creation
- `AGENTS.md`
  - branch and PR rules
  - "Start a new feature"
  - "Fix a bug"
- `docs/specs/triage-tasks.md`
  - fresh-session handoff contract and prompt template
- `README.md`
  - operator-facing workflow summary
- `skills/plan-execution/SKILL.md`
  - planning branch setup
- `skills/execute-task/SKILL.md`
  - execution branch setup
- `skills/triage-tasks/SKILL.md`
  - generated handoff prompt and "do now" transition

## Spec Changes

### `docs/specs/ww-dogfooding-workflow.md`

Create a workspace spec that defines:

- global `ww` as the default operator path
- the expected operator flow for workspace-root vs child-repo startup
- multi-task/worktree contention handling
- the boundary between stable global `ww` dogfooding and unreleased `ww` development
- failure handling and issue-filing policy for `ww`

### `docs/specs/triage-tasks.md`

Update the handoff contract so fresh-session prompts emit `ww` worktree setup commands instead of raw git branch commands.

### `AI_WORKFLOW.md`, `AGENTS.md`, `README.md`

Align top-level workflow docs around:

- using global `ww` for normal startup
- reserving raw git fallback for documented exceptional cases
- treating `ww` failures as workflow outputs that must be captured

## Code / Skill Changes

### Workspace repo

- Update `skills/plan-execution/SKILL.md` to use `ww create` / `ww cd` for planning startup
- Update `skills/execute-task/SKILL.md` to use `ww create` / `ww cd` for execution startup
- Update `skills/triage-tasks/SKILL.md` so handoff prompts recommend `ww`
- Evaluate whether `tools/workflow-lint.sh` should warn when workflow docs/skills emit raw git startup commands after the migration

### `ww` repo follow-up slices

If dogfooding exposes product gaps, capture them as explicit `ww` work:

- Slice A: improve diagnostics for startup failures (clearer collision, repo-target, or workspace-detection errors)
- Slice B: improve operator ergonomics for resuming or inspecting parallel worktrees if `ww list` / `ww cd` guidance proves insufficient
- Slice C: document or implement any missing stable-release behavior the workspace workflow now depends on

These `ww` slices should be filed only when a real dogfooding gap is observed; they must not be hand-waved away as incidental friction.

## Design Decisions

- Normal workflow startup should dogfood the released global `ww` binary, not bypass it with raw git
- A task branch should normally live in its own `ww` worktree rather than in the primary checkout
- The workspace must keep a clear boundary between "use stable global `ww` to do work" and "use the latest in-repo `ww` only when developing or debugging `ww` itself"
- `ww` failures are first-class workflow outputs and should feed back into `ww` planning/issues with concrete reproduction data

## Sub-tasks

- [ ] Update workspace specs/docs (`docs/specs/ww-dogfooding-workflow.md`, `docs/specs/triage-tasks.md`, `AI_WORKFLOW.md`, `AGENTS.md`, `README.md`) so they all describe `ww` as the default startup path
- [ ] [parallel] Update `skills/plan-execution/SKILL.md` and `skills/execute-task/SKILL.md` to replace raw git startup with `ww create` / `ww cd`
- [ ] [parallel] Update `skills/triage-tasks/SKILL.md` so fresh-session prompts recommend `ww` commands and explain exceptional fallback handling
- [ ] [parallel] Define the `ww` issue-filing checklist and ensure it is referenced consistently anywhere the workflow talks about unexpected tooling failures
- [ ] [depends on: spec/doc updates, skill updates] Dry-run the updated workflow from:
  - the workspace repo root
  - a child repo root
  - the workspace root targeting a child repo with `--repo`
- [ ] [depends on: dry-run verification] File explicit `ww` follow-up issues/plans for any product gaps exposed by dogfooding
- [ ] [depends on: settled workflow wording] Decide whether `tools/workflow-lint.sh` should gain a warning for workflow-facing raw git startup commands

## Verification

- Confirm every documented workflow startup path now prefers global `ww`
- Confirm fresh-session handoff prompts no longer default to raw `git switch -c`
- Confirm the docs distinguish stable global `ww` dogfooding from unreleased `ww` development inside `ww/`
- Confirm the `ww` issue-filing policy says when to stop, when fallback is allowed, and what evidence must be recorded
- Confirm the plan leaves room for corresponding `ww` product work instead of assuming the workflow can bypass all `ww` problems

## Expected Outcome

- The workspace workflow continuously dogfoods the released `ww` binary during normal planning and execution
- Parallel worktree-based development becomes the documented default rather than an optional side path
- `ww` bugs and confusing behavior discovered during everyday use get captured back into `ww` as durable work items
