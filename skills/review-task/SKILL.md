---
name: review-task
description: When creating or updating a pull request, preparing a PR for review, generating verification artifacts, collecting test results or screenshots for review, monitoring initial CI/Copilot follow-up, submitting changes for human review, or checking that all PR requirements (code, specs, plan, verification, follow-up) are met.
metadata:
  author: yoskeoka
  version: '2.0.0'
---

# Review Task (PR Preparation and Follow-up Gate)

**Position in workflow**: PR review is **not a standalone step** — it is embedded into every step of the AI-Centered Development cycle. Steps 1 (Project Plan), 2 (Execution Plan), and 3 (Execution) can route PR preparation through this skill so branch/type/title/scope discipline, verification evidence, PR readiness, and bounded post-PR follow-up are checked before human review.

## Step 1: Classify the Change

Before creating a `ww` worktree/branch for new work, classify the change into one of the workflow change types from `AI_WORKFLOW.md`. If work is already in progress on a branch, use this classification to verify the branch still matches the intended scope before any PR creation or PR update work:

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

If the existing branch does **not** match the classification, fix that mismatch before proceeding to PR preparation. Do not carry a mis-typed branch into review.

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

Do **NOT** create or update a PR for review until all required local checks are green.

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
- Otherwise, if working in a child repo that vendors this workflow under `.claude/vendor/workflow`, use `.claude/vendor/workflow/.github/PULL_REQUEST_TEMPLATE.md`.
- Otherwise, if working in the workflow repo itself, use `.github/PULL_REQUEST_TEMPLATE.md`.

If no PR exists yet, determine which template path applies in the current repo, then run:

```sh
git push origin <branch-name>
gh pr create --title "<descriptive title>" --body-file <template-path>
```

For example, if working in this workflow repo:

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

After the PR is created or updated, continue into the post-PR follow-up loop below. Do not hand off immediately after PR creation unless the user explicitly asked to stop before monitoring.

## Post-PR Follow-up Loop

For each new PR, updated PR, or later push to the PR branch, monitor the latest PR head SHA before handoff:

1. Record the PR number and current head SHA.
2. Wait for CI/check runs to settle for that SHA.
3. Inspect PR timeline/review data for GitHub Copilot auto-review activity.
4. If Copilot activity is present, wait a short bounded interval for review completion and comments.
5. If no Copilot activity appears within the bounded interval, stop waiting for Copilot and record that no Copilot review was observed.
6. Re-check the PR head SHA before handoff. If it changed, restart this loop for the new SHA.

Prefer delegating polling-style waiting to a low-cost worker when the platform and session allow delegation. Keep final decisions, code changes, review-comment triage, and user handoff in the main agent.

### CI Failure Loop

Treat CI failures as mechanical verification failures:

1. Inspect failing check logs.
2. If the failure is actionable from logs and the fix stays within branch scope, fix it in the same branch.
3. Re-run the relevant local verification.
4. Commit and push the fix.
5. Restart the post-PR follow-up loop for the new PR head SHA.

If the failure is not actionable, outside scope, or caused by external infrastructure, stop and document the blocker with the relevant check/log context.

### Copilot Review Triage

Treat GitHub Copilot comments as advisory review input, not automatic patch instructions. Do **not** silently auto-apply Copilot suggestions.

For each substantive Copilot comment, prepare a concise human-review briefing:

- comment summary
- whether action is recommended
- whether a GitHub suggestion can be applied as-is or needs adaptation
- whether the item should be deferred into `docs/issues/` or a future plan
- suggested response or implementation options

If later user-approved or workflow-authorized changes are pushed in response to Copilot or human comments, restart the post-PR follow-up loop for the new head SHA.

### Stop Conditions

The follow-up loop can stop when one of these is true:

- required checks pass and no Copilot auto-review activity appears within the bounded wait window
- required checks pass and available Copilot comments have been summarized
- CI is blocked or not actionable, and the blocker is documented
- the user explicitly asks to stop waiting

Wait for GitHub PR review approval before merging into `main`.

## Verification First Principle

Human review happens **after** mechanical tests and "visual" verification data are ready. The PR should contain enough evidence that a reviewer can confirm correctness without running the code themselves.

## After Merge

After the PR is merged, return to the appropriate workflow step:

- If the project plan needs updating → Step 1 (Project Plan)
- If the next task needs planning → Step 2 (Execution Plan)
- If a plan is ready for implementation → Step 3 (Execution)

Repeat steps 1–3 until the Project Plan is complete.
