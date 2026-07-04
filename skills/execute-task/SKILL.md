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

Before making any changes, create a fresh worktree from the latest `main` with the globally installed `ww` CLI:

From the target repo root:

```sh
ww create feat/<name>
cd "$(ww cd feat/<name>)"
```

From the workspace root when targeting a child repo:

```sh
ww create --repo <repo> feat/<name>
cd "$(ww cd --repo <repo> feat/<name>)"
```

Use `feat/<name>` or `fix/<name>` to match the execution-plan filename (e.g., `feat/initial-setup`, `fix/bug-x`).

## What to Do

Implement the changes described in the active execution plan, following strict ordering rules.

### Rules — Execution Order

1. Spec First: Update `docs/specs/` to reflect the intended changes _before_ modifying any code.
2. Implement: Write the code to match the updated spec exactly.
3. Log Issues: If unrelated problems are found during implementation, log them in `docs/issues/<sequence>-<name>.md`. Do **not** fix them within the current plan unless they are blockers.
4. Completion: When all work in the plan is done, move the plan file from `docs/exec-plan/todo/` to `docs/exec-plan/done/`.

### Spec-Code Parity

- `docs/specs/` must **strictly match** the actual code at all times.
- Never write code without a corresponding specification update.
- If the implementation reveals that the spec needs adjustment, update the spec first, then adjust the code.

### Issue Logging

When encountering unrelated issues during execution:

1. Create `docs/issues/<sequence>-<descriptive-name>.md`.
2. Document the issue clearly.
3. Continue with the current plan — do not get sidetracked.
4. These issues can become future execution plans.

### Issue Resolution

When an issue in `docs/issues/` is resolved:

1. Move the file from `docs/issues/<sequence>-<name>.md` to `docs/issues/done/<sequence>-<name>.md`.
2. Include the move in the same PR that fixes the issue.
3. Trivial issues (single-line fixes, typos, doc-only) may be fixed directly without a formal execution plan — just branch, fix, and PR.

When the matching execution plan resolves an external GitHub issue:

1. Keep that GitHub issue as the canonical tracker when the feedback came from outside the workspace and already lives in the target repo issue tracker.
2. Add a matching closing keyword to the implementation PR body:
   - same repo: `Closes #<number>`
   - cross repo: `Closes https://github.com/<owner>/<repo>/issues/<number>`
3. If the issue should stay open after this PR, explain that explicitly in the PR body instead of omitting the linkage silently.

### Self-Improvement Loop

After ANY correction from the user:

1. Create or update `docs/lessons.md` with the pattern, appending new lessons at the end of the file. (don't add lessons when the user correction creates issues/plans at the same time)
2. Use this format:
   - Mistake: What went wrong (be specific)
   - Pattern: The underlying cause or anti-pattern
   - Rule: Concrete, actionable rule to prevent recurrence
   - Applied: Where this rule applies (specific files, patterns, situations)
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
2. Move the numbered plan file: `docs/exec-plan/todo/<sequence>-<name>.md` → `docs/exec-plan/done/<sequence>-<name>.md`.
3. Proceed to Verify & PR.

## Verify (Pre-PR Gate)

Run **all** project lint and test commands using non-AI tooling (e.g., `make lint`, `npm run lint`, `go vet`, `pytest`, `npm test`, or whatever the project defines).

- If any check fails:
  1. Fix the issue in the same branch.
  2. Re-run the checks until **all pass**.
- Do **NOT** proceed to `review-task` until lint and tests are green.

## PR Workflow

After all checks pass:

1. Invoke **`review-task`** so the shared PR gate checks branch classification, exec-plan requirements, PR title alignment, scope, verification evidence, PR template completeness, and bounded post-PR follow-up before review.
2. Through that `review-task` flow, create the PR if one does not exist yet, or update the existing PR so it is ready for review and monitored for the latest head SHA.
3. The PR must include:
   - Code changes.
   - Spec updates (`docs/specs/`).
   - The plan file moved to `docs/exec-plan/done/`.
   - `Closes` entries for any external GitHub issues declared in the plan's `Addresses:` line, or an explicit reason they remain open.
   - Verification artifacts (test results, screenshots, logs) for human review.
4. Complete the initial CI/Copilot follow-up cycle through `review-task`, then wait for GitHub PR review approval before merging into `main`.
5. Do not report execution complete after only finishing local implementation, updating `docs/lessons.md`, drafting a commit, or opening the PR. Execution is complete only after `review-task` reaches a documented stop condition for the latest pushed PR head SHA. For a non-blocked PR create/update flow, that minimum landing path is `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll`, and if required CI checks are still pending for that Step 3 PR head SHA, `review-task` must use a bounded CI-settling window with up to two more `30-second wait -> poll` turns before handoff.

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
