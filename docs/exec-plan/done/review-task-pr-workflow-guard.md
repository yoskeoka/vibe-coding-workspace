# Review-Task PR Workflow Guard

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Make `review-task` the explicit PR-preparation gate for workflow work so lightweight changes do not skip branch-type, scope, and PR-title checks. Reduce duplicated PR instructions across skills, and add a small `AGENTS.md` reminder for trivial/no-plan changes.

Addresses:
- Recent workflow miss where a lightweight project-rules change went straight to PR creation without first classifying the branch as `chore/*`
- Ongoing duplication and drift risk between `review-task`, `plan-execution`, and `execute-task`

## Current State

- `review-task` contains shared PR workflow guidance, but it still teaches outdated raw-git branch setup and numbered branch examples.
- `plan-execution` and `execute-task` each restate their own PR creation flow instead of routing operators through `review-task`.
- `AGENTS.md` says trivial changes may skip an execution plan, but it does not explicitly remind the agent to re-check `AI_WORKFLOW.md` for the correct branch type and PR title before creating a PR.
- This leaves a gap for lightweight `docs/*` or `chore/*` changes: the agent can correctly decide "no exec-plan required" but still miss workflow-specific branch/type/title discipline.

## Code Changes

### `skills/review-task/SKILL.md`

- Update the branch guidance to match current `ww`-first workflow and current branch naming rules in `AI_WORKFLOW.md`.
- Add an explicit first step that classifies the current work as `plan`, `feat`, `fix`, `chore`, or `docs` before any PR creation work.
- Require a pre-PR check for:
  - branch name matches the classified change type
  - exec-plan requirement is satisfied or explicitly exempt
  - PR title matches the branch/change classification
  - over-scoping and obvious out-of-scope changes are called out before proceeding
- Make the skill's terminal path: if checks pass and no PR exists yet, create the PR; otherwise verify/update the existing PR so it is ready for review. This keeps PR creation as the completion path only for new PRs, not a separate implied step.

### `skills/plan-execution/SKILL.md`

- Replace duplicated detailed PR creation instructions with a shorter handoff that explicitly says to invoke `review-task` before creating the planning PR.
- Keep only step-specific requirements that are unique to execution-plan PRs.

### `skills/execute-task/SKILL.md`

- Replace duplicated detailed PR workflow instructions with an explicit handoff to `review-task` after implementation/spec/verification are complete.
- Keep only execution-specific requirements such as spec-first ordering, plan move to `done/`, and verification artifact expectations.

### `AGENTS.md`

- Add one explicit rule for lightweight/no-plan changes:
  - even when skipping an exec-plan, check `AI_WORKFLOW.md` before PR creation and align the branch type and PR title with the change scope (`docs`, `chore`, etc.)

## Spec Changes

None expected in `docs/specs/`. This work updates workflow/skill behavior documents rather than product/runtime behavior.

## Design Decisions

Past decisions:
- `AI_WORKFLOW.md` already establishes `ww` as the default startup path and makes `chore` / `docs` branches exempt from exec-plan mapping.
- `review-task` already exists as the shared PR workflow skill, so the right direction is to strengthen it rather than introduce a second overlapping PR skill.

Apply the same reasoning here: one shared PR gate is better than duplicated PR instructions spread across multiple skills.

## Sub-tasks

- [ ] [parallel] Update `skills/review-task/SKILL.md` so it becomes the explicit PR-preparation gate, including branch/type/title/scope checks
- [ ] [parallel] Update `skills/plan-execution/SKILL.md` and `skills/execute-task/SKILL.md` to route PR preparation through `review-task` and remove duplicated PR instructions
- [ ] [parallel] Add the lightweight-change guard sentence to `AGENTS.md`
- [ ] [depends on: all above] Verify the updated workflow reads coherently end-to-end and prepare the execution PR
