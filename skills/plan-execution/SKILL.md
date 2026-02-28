---
description: When planning a new task, creating an execution plan, starting a new feature, planning a bug fix, breaking down work into sub-tasks, creating a plan file in exec-plan/todo, or organizing what code and spec changes are needed before implementation.
---

# Plan Execution (Workflow Step 2)

**Position in workflow**: This is **Step 2** of the AI-Centered Development cycle. The project plan (Step 1) must exist before creating execution plans. After planning, proceed to Step 3 (Execution) to implement.

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
