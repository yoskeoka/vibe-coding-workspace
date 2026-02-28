---
description: When implementing code, writing code for a planned task, updating specs before coding, executing an existing plan, coding a feature or fix, updating docs/specs, logging unrelated issues found during implementation, or moving a completed plan from todo to done.
---

# Execute Task (Workflow Step 3)

**Position in workflow**: This is **Step 3** of the AI-Centered Development cycle. An execution plan (Step 2) must exist in `docs/exec-plan/todo/` before starting. After execution, proceed to Step 4 (Review) to create a PR.

## What to Do

Implement the changes described in the active execution plan, following strict ordering rules.

### Rules — Execution Order

1. **Spec First**: Update `docs/specs/` to reflect the intended changes *before* modifying any code.
2. **Implement**: Write the code to match the updated spec exactly.
3. **Log Issues**: If unrelated problems are found during implementation, log them in `docs/issues/<name>.md`. Do **not** fix them within the current plan unless they are blockers.
4. **Completion**: When all work in the plan is done, move the plan file from `docs/exec-plan/todo/` to `docs/exec-plan/done/`.

### Spec-Code Parity

- `docs/specs/` must **strictly match** the actual code at all times.
- Never write code without a corresponding specification update.
- If the implementation reveals that the spec needs adjustment, update the spec first, then adjust the code.

### Issue Logging

When encountering unrelated issues during execution:

1. Create `docs/issues/<descriptive-name>.md`.
2. Document the issue clearly.
3. Continue with the current plan — do not get sidetracked.
4. These issues can become future execution plans.

### Completing the Plan

When all tasks in the plan are done:

1. Verify spec-code parity: `docs/specs/` matches the implementation.
2. Move the plan file: `docs/exec-plan/todo/<NNN>-name.md` → `docs/exec-plan/done/<NNN>-name.md`.
3. Proceed to Review (Step 4).

## Next Step

After execution is complete and the plan is moved to `done/`, proceed to **Review** (Step 4): create a PR with all artifacts.
