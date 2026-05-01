---
name: plan-execution
description: When planning a new task, creating an execution plan, starting a new feature, planning a bug fix, breaking down work into sub-tasks, creating a plan file in exec-plan/todo, or organizing what code and spec changes are needed before implementation.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Plan Execution (Workflow Step 2)

**Position in workflow**: This is **Step 2** of the AI-Centered Development cycle. The project plan (Step 1) must exist and be merged before creating execution plans. This step requires its own branch and PR. After the plan is merged, proceed to Step 3 (Execution).

## Do I Need an Execution Plan?

```
Task received
    │
    ├─ Will docs/specs/ be created or updated? → YES, create exec plan
    │
    ├─ 3+ steps or architectural decision? → YES, create exec plan
    │
    ├─ Multiple files will be modified? → YES, create exec plan
    │
    ├─ Uncertain about approach? → YES, create exec plan
    │
    ├─ Single-line fix / typo / trivial change? → NO, execute directly
    │
    └─ Otherwise → Proceed directly
```

**When in doubt, create a plan.** The cost of an unnecessary plan is low; the cost of a botched unplanned change is high.

## Branch Setup

Before making any changes, create a fresh worktree from the latest `main` with the globally installed `ww` CLI:

From the target repo root:

```sh
ww create plan/<name>
cd "$(ww cd plan/<name>)"
```

From the workspace root when targeting a child repo:

```sh
ww create --repo <repo> plan/<name>
cd "$(ww cd --repo <repo> plan/<name>)"
```

Use a descriptive kebab-case name that matches the exec-plan filename (e.g., `plan/initial-setup`, `plan/feature-name`).

## What to Do

Create a new plan file in `docs/exec-plan/todo/` that details the work to be done.

### Rules

1. **Read planning context before designing**: Read `docs/project-plan.md`, `docs/design-decisions/core-beliefs.md`, and `docs/design-decisions/adr.md` before drafting the plan or choosing an approach.
2. **Discuss meaningful trade-offs before deciding**: If the plan has multiple viable approaches, list the options even when one option is common or conventional. Compare them against `docs/project-plan.md`, `docs/design-decisions/core-beliefs.md`, and `docs/design-decisions/adr.md`, identify the recommended option for the current project, then discuss that recommendation with the human before locking it into the plan.
3. **One plan per logical unit of work**: Each plan file corresponds to a coherent, reviewable chunk of changes.
4. **Naming convention**: For a single-file plan, name the file `<name>.md` so it matches the branch name `plan/<name>`, e.g., `initial-setup.md`, `feature-name.md`, `fix-bug-x.md`. If a large goal is split into multiple child plans, the branch MUST still use the shared parent/base name `plan/<major-item>`, and each child plan filename MUST derive from that same `<major-item>` base. When execution order matters across the child plans, use `<major-item>-<two-digit-sequence>-<minor-item>.md`, such as branch `plan/platform-phase2` with child plans `platform-phase2-01-foundation.md`, `platform-phase2-02-fixture-e2e.md`. If a plan is split into multiple child files with no meaningful order, continue using non-numeric, descriptive names derived from the same branch/base name, such as branch `plan/feature-name` with `feature-name-backend.md` and `feature-name-ui.md`; do not use unrelated numeric prefixes like `004a-feature-name.md`.
5. **Split-parent handling**: If one existing parent plan is being decomposed into multiple child plans, do not overwrite or rewrite that parent into another executable plan. Verify that all spec/code/verification scope has been moved into the child plans without gaps, then move the parent plan from `docs/exec-plan/todo/` to `docs/exec-plan/done/` as-is. The completion action is the move itself, not a content rewrite of the parent plan.
6. **Check `docs/exec-plan/done/`** for completed plans to understand prior work and avoid duplication.

### Plan File Content

Each plan file in `docs/exec-plan/todo/` must detail:

- **Objective**: What this plan accomplishes (linked to project requirements).
- **Code changes**: What files/modules will be created or modified.
- **Spec changes**: How `docs/specs/` will be updated to reflect the changes.
- **Sub-tasks**: Break large tasks into smaller steps if needed.
- **Design decisions**: If architectural choices are being made, note them for `docs/design-decisions/adr.md`.
- **Parallelism**: Identify which sub-tasks are independent (see below).
- **Execution instruction**: Place the following line immediately below the H1 title in every executable plan file:
  > **Execution**: Use `/execute-task` to implement this plan.
- **Split-parent exception**: When a parent plan is only being retired because its scope was split into child plans, do not edit the parent plan body to turn it into a new record document. Verify migration completeness separately, then move the parent plan to `docs/exec-plan/done/` without adding a new execution instruction or rewriting its contents.

### Parallel Execution Planning

Design plans for maximum parallel execution:

1. **Identify independent tasks**: Mark sub-tasks that have no dependencies on each other with `[parallel]`.
2. **Explicit dependency notation**: Use `depends on: <task>` for tasks that must wait for others.
3. **Split plans when appropriate**: If a large plan contains 2+ fully independent streams of work, or if one large goal is clearer as ordered child plans, create separate plan files. For independent streams, use descriptive names derived from the work itself, such as branch `plan/feature-name` with `feature-name-backend.md` and `feature-name-ui.md`, or names like `api-endpoints.md` and `ui-components.md` when they share an obvious base in context. For ordered phases toward one larger goal, use the shared parent/base branch name plus `<major-item>-<two-digit-sequence>-<minor-item>.md`, and make the order explicit, e.g., branch `plan/platform-phase2` with `platform-phase2-01-foundation.md`, `platform-phase2-02-fixture-e2e.md`. For large goals with ordered child plans, this naming is required rather than optional.
4. **When replacing one large executable plan with child plans, audit for migration completeness**: explicitly verify that every spec change, code change, verification item, and design note from the original parent plan has either moved into an appropriate child plan or been intentionally dropped with justification. After that audit, retire the parent by moving it from `todo/` to `done/` without editing the file contents.
5. **When child plans reference a parent plan, use a path that survives completion**: parent plans move from `docs/exec-plan/todo/` to `docs/exec-plan/done/` after execution or split completion, so do not hardcode a `todo/`-only path in child-plan metadata. Prefer wording that references the filename and notes it may live at `docs/exec-plan/todo/<name>.md` or `docs/exec-plan/done/<name>.md`, without assuming the current directory.

Example:

```markdown
## Sub-tasks

- [ ] [parallel] Create API schema types
- [ ] [parallel] Set up database migration
- [ ] [depends on: API schema, DB migration] Implement API handlers
- [ ] [depends on: API handlers] Write integration tests
```

**Why**: Independent tasks can be delegated to subagents for parallel execution, reducing total time and keeping the main context clean.

### For a New Feature

1. Read `docs/project-plan.md`.
2. Create `docs/exec-plan/todo/<feature-name>.md`.
3. Outline the spec and code changes.
4. Wait for user confirmation or proceed if authorized.

### For a Bug Fix

1. Create `docs/exec-plan/todo/<fix-bug-name>.md`.
2. Include reproduction steps in the plan.
3. Outline the fix approach, spec updates, and verification steps.

### Architectural Decisions

If the plan involves architectural choices, review and update `docs/design-decisions/`:

- `adr.md`: Append the decision with context and rationale.
- `core-beliefs.md`: Verify the decision aligns with guiding principles.

When trade-offs exist:

1. List the viable options and the practical consequences of each.
2. Explain how prior decisions in `docs/design-decisions/adr.md` and current goals in `docs/project-plan.md` affect the choice.
3. Present a recommended option as a recommendation, not a final decision.
4. Ask the human to confirm, reject, or refine the recommendation before recording it in the plan.

## PR Workflow

After the plan file is created:

1. Commit the plan file and any related `docs/design-decisions/` updates on the branch.
2. Invoke **`review-task`** before any PR creation work so the branch classification, branch name, PR title, scope, template/body, and bounded post-PR follow-up are checked through the shared PR gate.
3. Create or update the planning PR through that `review-task` flow, complete the initial CI/Copilot follow-up cycle for the latest PR head SHA, then wait for GitHub PR review approval before merging into `main`.

## Next Step

After the execution plan PR is merged into `main`, invoke **`/execute-task`** to start implementation. Do not proceed without invoking the skill — it defines the execution workflow (spec-first, branch setup, PR gate, etc.).
