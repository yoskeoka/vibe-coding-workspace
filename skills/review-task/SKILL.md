---
name: review-task
description: When creating a pull request, preparing a PR for review, generating verification artifacts, collecting test results or screenshots for review, submitting changes for human review, or checking that all PR requirements (code, specs, plan, verification) are met.
metadata:
  author: yoskeoka
  version: '2.0.0'
---

# Review Task (PR Workflow Reference)

**Position in workflow**: PR review is **not a standalone step** — it is embedded into every step of the AI-Centered Development cycle. Steps 1 (Project Plan), 2 (Execution Plan), and 3 (Execution) each require their own branch and PR. This skill describes the shared PR requirements and verification standards.

## Branch Rule

Every PR must come from a fresh branch created from the latest `main`:

```sh
git fetch origin
git switch -c <branch-name> origin/main
```

Branch naming conventions:
- Project plan changes: `plan/project-plan-<description>`
- Execution plan changes: `plan/<NNN>-<description>`
- Code execution: `feat/<NNN>-<description>` or `fix/<NNN>-<description>`

## Pre-PR Gate (Verify)

Before creating any PR, run **all** applicable checks:

1. **Lint & Test**: Run project lint and test commands using non-AI tooling (e.g., `make lint`, `npm run lint`, `go vet`, `pytest`, `npm test`).
2. **Fix failures**: If any check fails, fix in the same branch and re-run until all pass.
3. **Doc-only PRs**: Skip lint/test when no tooling covers documentation, but still verify Markdown formatting.

Do **NOT** create a PR until all checks are green.

## PR Must Include

The PR contents depend on which workflow step produced it:

### For Step 1 (Project Plan) PRs
- Updated `docs/project-plan.md`.

### For Step 2 (Execution Plan) PRs
- New plan file in `docs/exec-plan/todo/`.
- Any `docs/design-decisions/` updates if architectural choices were made.

### For Step 3 (Execution) PRs
1. **Code changes**: The implementation.
2. **Spec updates**: The updated `docs/specs/` files that match the code.
3. **Plan file moved to `done/`**: The execution plan in `docs/exec-plan/done/` proving the task was completed through the proper workflow.
4. **Verification artifacts**: Test results, screenshots, logs, or other evidence for human reviewers. Human review happens _after_ mechanical tests and verification data are ready.

## Verification Standards by Task Type

| Task Type   | Minimum Verification             |
| ----------- | -------------------------------- |
| Bug fix     | Reproduce → Fix → Verify fixed   |
| Feature     | Tests pass + manual demo         |
| Refactor    | Behavior unchanged + tests pass  |
| Performance | Before/after metrics             |
| Security    | Specific vulnerability addressed |

## Pre-PR Checklist

Before creating the PR, verify:

- [ ] Branch was created from the latest `origin/main`.
- [ ] `docs/specs/` matches the implementation (Spec-Code Parity) — for Step 3 PRs.
- [ ] Plan file has been moved from `docs/exec-plan/todo/` to `docs/exec-plan/done/` — for Step 3 PRs.
- [ ] All lint and test checks pass (non-AI tooling).
- [ ] Any visual or behavioral changes have screenshots/logs attached.
- [ ] No unresolved blockers remain (non-blockers should be in `docs/issues/`).

## Creating the PR

```sh
git push origin <branch-name>
gh pr create --title "<descriptive title>" --body "<summary of changes>"
```

Wait for GitHub PR review approval before merging into `main`.

## Verification First Principle

Human review happens **after** mechanical tests and "visual" verification data are ready. The PR should contain enough evidence that a reviewer can confirm correctness without running the code themselves.

## After Merge

After the PR is merged, return to the appropriate workflow step:
- If the project plan needs updating → Step 1 (Project Plan)
- If the next task needs planning → Step 2 (Execution Plan)
- If a plan is ready for implementation → Step 3 (Execution)

Repeat steps 1–3 until the Project Plan is complete.
