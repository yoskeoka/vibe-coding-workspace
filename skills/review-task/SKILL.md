---
name: review-task
description: When creating a pull request, preparing a PR for review, generating verification artifacts, collecting test results or screenshots for review, submitting changes for human review, or checking that all PR requirements (code, specs, plan, verification) are met.
metadata:
  author: yoskeoka
  version: '2.0.0'
---

# Review Task (PR Preparation Gate)

**Position in workflow**: PR review is **not a standalone step** — it is embedded into every step of the AI-Centered Development cycle. Steps 1 (Project Plan), 2 (Execution Plan), and 3 (Execution) each require their own branch and PR. This skill is the shared pre-PR gate that checks branch/type/title/scope discipline, verification evidence, and PR readiness before human review.

## Step 1: Classify the Change

Before any PR creation or PR update work, classify the current branch into one of the workflow change types from `AI_WORKFLOW.md`:

- `plan`: Project-plan or execution-plan authoring/update work
- `feat`: Implementation work executing an approved plan
- `fix`: Bug-fix work executing an approved plan
- `chore`: Non-functional changes such as CI, tooling, or dependency updates
- `docs`: Documentation-only changes that do not fit the plan/execute flow

Use the classification to drive every later PR decision:

- branch name
- exec-plan requirement or exemption
- PR title
- PR template checkboxes
- reviewer expectations

## Branch Rule

Every PR must come from a fresh branch created from the latest `main` with the globally installed `ww` CLI:

From the target repo root:

```sh
ww create <type>/<description>
cd "$(ww cd <type>/<description>)"
```

From the workspace root when targeting a child repo:

```sh
ww create --repo <repo> <type>/<description>
cd "$(ww cd --repo <repo> <type>/<description>)"
```

Branch naming follows `AI_WORKFLOW.md`:

- Project plan or execution plan changes: `plan/<name>`
- Approved-plan implementation: `feat/<name>` or `fix/<name>`
- Non-plan exempt work: `chore/<name>` or `docs/<name>`

## Pre-PR Gate (Verify)

Before creating any PR, run **all** applicable checks:

1. **Lint & Test**: Run project lint and test commands using non-AI tooling (e.g., `make lint`, `npm run lint`, `go vet`, `pytest`, `npm test`).
2. **Fix failures**: If any check fails, fix in the same branch and re-run until all pass.
3. **Doc-only PRs**: Skip lint/test when no tooling covers documentation, but still verify Markdown formatting.

Do **NOT** create or update a PR for review until all required checks are green.

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

Before creating or updating the PR, verify:

- [ ] The work has been classified as `plan`, `feat`, `fix`, `chore`, or `docs`.
- [ ] Branch was created from the latest `origin/main`.
- [ ] Branch name matches the classified change type and follows `<type>/<description>`.
- [ ] Exec-plan requirement is satisfied (`feat/*` and `fix/*`) or explicitly exempt (`plan/*`, `chore/*`, `docs/*`).
- [ ] PR title matches the classified change type and scope.
- [ ] `docs/specs/` matches the implementation (Spec-Code Parity) — for Step 3 PRs.
- [ ] Plan file has been moved from `docs/exec-plan/todo/` to `docs/exec-plan/done/` — for Step 3 PRs.
- [ ] All lint and test checks pass (non-AI tooling).
- [ ] Any visual or behavioral changes have screenshots/logs attached.
- [ ] The diff is not obviously over-scoped for the branch/plan; any out-of-scope changes are removed or called out before proceeding.
- [ ] No unresolved blockers remain (non-blockers should be in `docs/issues/`).

### PR Title Rule

The PR title should make the workflow classification obvious to reviewers.

- `plan/*`: title describes the plan being added or updated
- `feat/*`: title describes the implemented feature or workflow behavior
- `fix/*`: title describes the bug being fixed
- `chore/*`: title is clearly labeled as maintenance/tooling/CI work
- `docs/*`: title is clearly labeled as documentation-only work

Do not use a title that implies implementation when the branch is `docs/*` or `chore/*`, and do not use a generic docs/chore title for `feat/*` or `fix/*` execution work.

## Completing the Gate

Use the **PR template** when creating or repairing pull requests. Template priority:

- If the child project has `.github/PULL_REQUEST_TEMPLATE.md`, use it.
- Otherwise, use the workspace-level `.github/PULL_REQUEST_TEMPLATE.md`.

If no PR exists yet, determine which template to use (project-level if present, otherwise workspace-level), then run:

```sh
git push origin <branch-name>
gh pr create --title "<descriptive title>" --body-file <template-path>
```

For example, if using the workspace-level template:

```sh
gh pr create --title "<descriptive title>" --body-file .github/PULL_REQUEST_TEMPLATE.md
```

> **Note**: `--fill` populates the title/body from commits, not from the PR template. Use `--body-file` to pre-populate with the correct template content.

If a PR already exists, do **not** treat PR creation as the next mandatory step. Instead:

1. Confirm the branch still matches the intended scope.
2. Confirm the existing PR title/body/checklists still match the current classification and diff.
3. Update the existing PR so it is review-ready rather than creating a duplicate PR.

Whether the PR is new or existing, complete the template/body so these sections are correct:

1. **Plan / Issues** — Link the exec-plan, issue, or project-plan that triggered this PR.
2. **Type of Change** — Check the applicable box.
3. **Instructions** — Fill in the execution command. Under "Additional Context from Instructing Human", record any human instructions, decisions, or intent NOT already captured in the plan/specs/code. Include the AI's question when the human's answer was brief (e.g., "Yes" or "A") so the context is self-contained.
4. **Verification** — Check off and fill in the commands used.
5. **Checklist** — Confirm all items.
6. **Dependencies** — List PRs/issues that block or are blocked by this PR. N/A if none.
7. **Reviewer Notes** — Highlight areas for review focus, known trade-offs, or intentional oddities. N/A if none.
8. **Links** — External references (library docs, design references, discussions). N/A if none.
9. **Breaking Changes / Screenshots** — Fill or delete as applicable.

Wait for GitHub PR review approval before merging into `main`.

## Verification First Principle

Human review happens **after** mechanical tests and "visual" verification data are ready. The PR should contain enough evidence that a reviewer can confirm correctness without running the code themselves.

## After Merge

After the PR is merged, return to the appropriate workflow step:

- If the project plan needs updating → Step 1 (Project Plan)
- If the next task needs planning → Step 2 (Execution Plan)
- If a plan is ready for implementation → Step 3 (Execution)

Repeat steps 1–3 until the Project Plan is complete.
