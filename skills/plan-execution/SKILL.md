---
name: plan-execution
description: When planning a new task, creating an execution plan, starting a new feature, planning a bug fix, breaking down work into sub-tasks, creating a plan file in exec-plan/todo, or organizing what code and spec changes are needed before implementation.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Plan Execution (Workflow Step 2)

**Position in workflow**: This is **Step 2** of the AI-Centered Development cycle. The project plan (Step 1) must exist before creating execution plans. After planning, proceed to Step 3 (Execution) to implement.

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

## What to Do

Create a new plan file in `docs/exec-plan/todo/` that details the work to be done.

### Rules

1. **Read `docs/project-plan.md` first**: Understand the project's goals before planning tasks.
2. **One plan per logical unit of work**: Each plan file corresponds to a coherent, reviewable chunk of changes.
3. **Naming convention**: Use sequential numbering, e.g., `001-initial-setup.md`, `002-feature-name.md`, `003-fix-bug-x.md`.
4. **Check `docs/exec-plan/done/`** for completed plans to understand prior work and avoid duplication.

### Plan File Content

Each plan file in `docs/exec-plan/todo/` must detail:

- **Objective**: What this plan accomplishes (linked to project requirements).
- **Code changes**: What files/modules will be created or modified.
- **Spec changes**: How `docs/specs/` will be updated to reflect the changes.
- **Sub-tasks**: Break large tasks into smaller steps if needed.
- **Design decisions**: If architectural choices are being made, note them for `docs/design-decisions/adr.md`.
- **Parallelism**: Identify which sub-tasks are independent (see below).

### Parallel Execution Planning

Design plans for maximum parallel execution:

1. **Identify independent tasks**: Mark sub-tasks that have no dependencies on each other with `[parallel]`.
2. **Explicit dependency notation**: Use `depends on: <task>` for tasks that must wait for others.
3. **Split plans when appropriate**: If a large plan contains 2+ fully independent streams of work, create separate plan files (e.g., `004a-api-endpoints.md`, `004b-ui-components.md`).

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
2. Create `docs/exec-plan/todo/<NNN>-feature-name.md`.
3. Outline the spec and code changes.
4. Wait for user confirmation or proceed if authorized.

### For a Bug Fix

1. Create `docs/exec-plan/todo/<NNN>-fix-bug-name.md`.
2. Include reproduction steps in the plan.
3. Outline the fix approach, spec updates, and verification steps.

### Architectural Decisions

If the plan involves architectural choices, review and update `docs/design-decisions/`:

- `adr.md`: Append the decision with context and rationale.
- `core-beliefs.md`: Verify the decision aligns with guiding principles.

## Next Step

After the plan is created and confirmed, proceed to **Execution** (Step 3): update specs first, then implement code.
