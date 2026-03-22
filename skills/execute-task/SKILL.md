---
name: execute-task
description: When implementing code, writing code for a planned task, updating specs before coding, executing an existing plan, coding a feature or fix, updating docs/specs, logging unrelated issues found during implementation, or moving a completed plan from todo to done.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Execute Task (Workflow Step 3)

**Position in workflow**: This is **Step 3** of the AI-Centered Development cycle. An execution plan (Step 2) must be merged into `main` before starting. This step requires its own branch, and all lint/tests must pass before creating a PR.

## Branch Setup

Before making any changes, create a fresh branch from the latest `main`:

```sh
git fetch origin
git switch -c feat/<NNN>-<description> origin/main
```

Use a descriptive branch name (e.g., `feat/001-initial-setup`, `fix/003-bug-x`).

## What to Do

Implement the changes described in the active execution plan, following strict ordering rules.

### Rules — Execution Order

1. **Spec First**: Update `docs/specs/` to reflect the intended changes _before_ modifying any code.
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

### Issue Resolution

When an issue in `docs/issues/` is resolved:

1. Move the file from `docs/issues/<name>.md` to `docs/issues/done/<name>.md`.
2. Include the move in the same PR that fixes the issue.
3. Trivial issues (single-line fixes, typos, doc-only) may be fixed directly without a formal execution plan — just branch, fix, and PR.

### Self-Improvement Loop

After ANY correction from the user:

1. Create or update `docs/lessons.md` with the pattern.
2. Use this format:
   - **Mistake**: What went wrong (be specific)
   - **Pattern**: The underlying cause or anti-pattern
   - **Rule**: Concrete, actionable rule to prevent recurrence
   - **Applied**: Where this rule applies (specific files, patterns, situations)
3. Review `docs/lessons.md` at session start for relevant learnings.

> "Be more careful" is not a rule. Rules must be specific and testable.

### Elegance Check

For non-trivial changes, pause and ask: **"Is there a more elegant way?"**

- If a fix feels hacky → reconsider the approach before committing
- If similar code exists elsewhere → look for reuse or abstraction
- If the change touches 5+ files → verify it's the minimal approach

**Skip this for**: single-line fixes, obvious corrections, established patterns in the codebase. Don't over-engineer.

### Completing the Plan

When all tasks in the plan are done:

1. Verify spec-code parity: `docs/specs/` matches the implementation.
2. Move the plan file: `docs/exec-plan/todo/<NNN>-name.md` → `docs/exec-plan/done/<NNN>-name.md`.
3. Proceed to Verify & PR.

## Verify (Pre-PR Gate)

Run **all** project lint and test commands using non-AI tooling (e.g., `make lint`, `npm run lint`, `go vet`, `pytest`, `npm test`, or whatever the project defines).

- If any check fails:
  1. Fix the issue in the same branch.
  2. Re-run the checks until **all pass**.
- Do **NOT** proceed to PR creation until lint and tests are green.

## PR Workflow

After all checks pass:

1. Push the branch and create a PR via `gh pr create`, using the **PR template**.
   - Template priority: **project-level** `.github/PULL_REQUEST_TEMPLATE.md` > **workspace-level** template.
2. Fill in all template sections. Pay special attention to:
   - **Plan / Issues**: Link the exec-plan and any issues that triggered this work.
   - **Instructions → Additional Context from Instructing Human**: Record any human decisions, preferences, or directives that are NOT captured in the plan, specs, or code diff. When the human's answer was brief ("Yes", "A"), include the AI's question/proposal for context.
3. The PR must include:
   - Code changes.
   - Spec updates (`docs/specs/`).
   - The plan file moved to `docs/exec-plan/done/`.
   - Verification artifacts (test results, screenshots, logs) for human review.
4. Wait for GitHub PR review approval before merging into `main`.

### Verification Standards by Task Type

| Task Type   | Minimum Verification             |
| ----------- | -------------------------------- |
| Bug fix     | Reproduce → Fix → Verify fixed   |
| Feature     | Tests pass + manual demo         |
| Refactor    | Behavior unchanged + tests pass  |
| Performance | Before/after metrics             |
| Security    | Specific vulnerability addressed |

## Next Step

After the PR is merged, repeat from **Step 1** (if the project plan needs updating) or **Step 2** (Execution Plan) for the next task. Repeat until the Project Plan is complete.
