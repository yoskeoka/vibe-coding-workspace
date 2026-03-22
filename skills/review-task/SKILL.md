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

Use the **PR template** when creating pull requests. Template priority:
- If the child project has `.github/PULL_REQUEST_TEMPLATE.md`, use it.
- Otherwise, use the workspace-level `.github/PULL_REQUEST_TEMPLATE.md`.

```sh
git push origin <branch-name>
gh pr create --title "<descriptive title>" --fill
```

The `--fill` flag auto-populates from the template. After creation, edit the PR body to complete all template sections:

1. **Plan / Issues** — Link the exec-plan, issue, or project-plan that triggered this PR.
2. **Type of Change** — Check the applicable box.
3. **Instructions** — Fill in the execution command. Under "Additional Context from Instructing Human", record any human instructions, decisions, or intent NOT already captured in the plan/specs/code. Include the AI's question when the human's answer was brief (e.g., "Yes" or "A") so the context is self-contained.
4. **Verification** — Check off and fill in the commands used.
5. **Checklist** — Confirm all items.
6. **Breaking Changes / Screenshots** — Fill or delete as applicable.

Wait for GitHub PR review approval before merging into `main`.

## Verification First Principle

Human review happens **after** mechanical tests and "visual" verification data are ready. The PR should contain enough evidence that a reviewer can confirm correctness without running the code themselves.

## After Merge

After the PR is merged, return to the appropriate workflow step:
- If the project plan needs updating → Step 1 (Project Plan)
- If the next task needs planning → Step 2 (Execution Plan)
- If a plan is ready for implementation → Step 3 (Execution)

Repeat steps 1–3 until the Project Plan is complete.
