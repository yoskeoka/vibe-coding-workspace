---
name: review-task
description: When creating or updating a pull request, preparing a PR for review, generating verification artifacts, collecting test results or screenshots for review, monitoring initial CI/advisory bot/agent follow-up, submitting changes for human review, or checking that all PR requirements (code, specs, plan, verification, follow-up) are met.
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
- If the plan's `Addresses:` line lists external GitHub issues, the PR body links the same issues under `Issues`.

### For Step 3 (Execution) PRs

1. **Code changes**: The implementation.
2. **Spec updates**: The updated `docs/specs/` files that match the code.
3. **Plan file moved to `done/`**: The execution plan in `docs/exec-plan/done/` proving the task was completed through the proper workflow.
4. **Verification artifacts**: Test results, screenshots, logs, or other evidence for human reviewers. Human review happens _after_ mechanical tests and verification data are ready.
5. **External issue closure metadata**: If the matching execution plan lists external GitHub issues in `Addresses:`, the PR body contains matching `Closes` entries or an explicit reason those issues remain open.

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
- [ ] External GitHub issues declared in the plan's `Addresses:` line are linked under `Issues` for Step 2 PRs and closed with explicit `Closes` entries for Step 3 PRs, or the PR body explains why they remain open.
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

- If the current repo has `.github/PULL_REQUEST_TEMPLATE.md`, use it.
- Otherwise, if working in a child repo inside this workspace, use the workspace root repo template at `<workspace-root>/.github/PULL_REQUEST_TEMPLATE.md`.
- Otherwise, if working in a child repo that vendors this workflow under `.claude/vendor/workflow`, use `.claude/vendor/workflow/.github/PULL_REQUEST_TEMPLATE.md`.
- Otherwise, use the workflow repo's `.github/PULL_REQUEST_TEMPLATE.md`.

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

1. **Plan / Issues** — Link the exec-plan, issue, or project-plan that triggered this PR. If the matching plan `Addresses:` line lists external GitHub issues, repeat those issues here.
2. **Closes** — For Step 3 PRs that resolve external GitHub issues from the plan's `Addresses:` line, add one explicit closing keyword per issue. Use `Closes #123` only when the issue is in the same repo as the PR; otherwise use the full issue URL.
3. **Type of Change** — Check the applicable box.
4. **Instructions** — Fill in the execution command. Under "Additional Context from Instructing Human", record any human instructions, decisions, or intent NOT already captured in the plan/specs/code. Include the AI's question when the human's answer was brief (e.g., "Yes" or "A") so the context is self-contained.
5. **Verification** — Check off and fill in the commands used.
6. **Checklist** — Confirm all items.
7. **Dependencies** — List PRs/issues that block or are blocked by this PR. N/A if none.
8. **Reviewer Notes** — Highlight areas for review focus, known trade-offs, or intentional oddities. N/A if none.
9. **Links** — External references (library docs, design references, discussions). N/A if none.
10. **Breaking Changes / Screenshots** — Fill or delete as applicable.

After the PR is created or updated, continue into the post-PR follow-up loop below. Do not hand off immediately after PR creation unless the user explicitly asked to stop before monitoring.

## Post-PR Follow-up Loop

For each new PR, updated PR, or later push to the PR branch, monitor the latest PR head SHA before handoff:

1. Record the PR number and current head SHA.
2. Wait 30 seconds after PR creation or a later push so CI/checks and review automation have time to start.
3. Poll compact PR follow-up state with:

   ```sh
   skills/review-task/scripts/gh-pr-followup poll <owner/repo> <pr-number>
   ```

   The helper returns the current head SHA, review decision, compact check rollup, review summaries, new timeline events, and new inline review comments. Repeated polls use `.local/gh-pr-followup/` markers so old timeline events and inline comments are not pasted back into the main context.
4. Inspect CI/check status for that SHA from the helper output and continue the CI failure loop below when checks fail or expose actionable logs. This check inspection happens for every pushed PR head SHA.
5. Inspect new timeline events from the helper output to detect advisory bot/agent reviewer activity.
6. If the helper is missing or fails, report the failure reason and stop automatic follow-up for this PR head SHA. Do **not** automatically fall back to raw GitHub reads; spending large context on raw timeline/comment JSON is not appropriate for routine hobby-project PR monitoring. Tell the user the PR can be checked later, or run a targeted raw `gh` command only if the user explicitly asks or the helper itself needs diagnosis.
7. Detect advisory bot/agent reviewer activity from timeline events, including:
   - `copilot_work_started`
   - `review_requested` events where `requested_reviewer.login` or `requested_team.name` identifies Copilot, Claude, `gh aw`, agent workflow, or another configured bot/agent reviewer
   - timeline events from bot or agent actors that indicate review work has started
8. If advisory reviewer activity has started for the latest head SHA but no final review/comments are visible yet, wait for review completion/comments using the bounded advisory-review cadence below.
9. If no advisory reviewer activity is present, do not spend the advisory-review wait budget; record that no advisory review start was observed.
10. Before handoff, use the latest helper output for review summaries and inline review comments. If the helper failed, hand off the failure reason instead of fetching raw review bodies or already-seen comments.
11. Triage substantive advisory bot/agent findings from the implementation context and decide whether each item should be fixed in the current PR, deferred, or treated as no action.
12. Re-check the PR head SHA before handoff. If it changed, restart this loop for the new SHA.

CI settling cadence by workflow step:

- Step 2 (`plan-execution`) keeps the existing minimum path: after the first 30-second wait and poll, it may stop when no other stop-condition work remains.
- Step 3 (`execute-task`) uses a bounded CI-settling window as part of the landing check. After the first 30-second wait and poll, if required checks for the latest head SHA are still pending, wait another 30 seconds and poll again. If required checks are still pending after that second poll, wait a third 30-second turn and poll once more.
- Stop the extra Step 3 settling waits early when required checks finish, advisory reviewer activity starts and moves the flow into the advisory-review cadence, the helper fails, the head SHA changes, or the user explicitly asks to stop waiting.
- This Step 3 CI-settling extension adds at most two extra 30-second wait turns beyond the startup wait, for a maximum of three polls separated by 30-second waits before handoff when required checks are merely still running.

Bounded advisory-review wait cadence:

- First wait: 3 minutes.
- Second wait: 2 minutes.
- Third wait: 1 minute.
- Fourth wait: 1 minute.
- Total wait budget: 7 minutes across 4 polling turns.
- After each wait turn, poll with `skills/review-task/scripts/gh-pr-followup`. If review/comments were submitted, stop waiting and triage them. If the helper fails, stop the automatic wait loop and report the failure reason.

If advisory bot/agent review has started but no review/comments have been submitted after the 7-minute budget, treat the advisory-review wait as timed out and document the state in the handoff.

Prefer delegating polling-style waiting to a low-cost subagent only when the platform supports delegation and the current session explicitly authorizes subagent use. Keep final decisions, code changes, review-comment triage, and user handoff in the main agent.

### CI Failure Loop

Treat CI failures as mechanical verification failures:

1. Inspect failing check logs.
2. If the failure is actionable from logs and the fix stays within branch scope, fix it in the same branch.
3. Re-run the relevant local verification.
4. Commit and push the fix.
5. Restart the post-PR follow-up loop for the new PR head SHA.

If the failure is not actionable, outside scope, or caused by external infrastructure, stop and document the blocker with the relevant check/log context.

### Advisory Bot/Agent Review Triage

Treat Copilot, Claude, `gh aw`, agent workflow reviews, and other configured bot/agent comments as advisory review input, not automatic patch instructions.

Passing or approving advisory bot/agent checks can still contain substantive observations in review bodies. Inspect review summaries and inline comments even when the overall state is not blocking.

For each substantive advisory finding, prepare a concise human-review briefing grouped by source reviewer/workflow:

- source reviewer/workflow
- comment location or link
- extracted comment summary
- implementer's view
- concise 1-2 line explanation
- recommendation: fix in this PR, defer, or no action

After implementing changes, evaluate advisory bot/agent findings from the implementation context and present response options in the current session. Do not post that triage back to the PR unless the user explicitly asks for a PR comment.

When the implementer's view is `fix in this PR` and the change remains reasonably scoped to the current PR, make that follow-up change before handoff rather than asking for another human approval turn.

When the implementer's view is `defer`, route the follow-up based on scope clarity:

- if the separate larger change has a known direction, create a new `plan-execution` task
- if the separate larger change is real but the solution is still unclear, create a `docs/issues/` item

Do **not** silently auto-apply every suggestion. Use implementer judgment to keep branch mutations in scope, and treat clearly separate or large rewrites as deferred work instead of stretching the current PR.

If later workflow-authorized changes are pushed in response to advisory bot/agent or human comments, restart required CI/check inspection for the new head SHA. Skip the longer advisory-review wait unless new review-start activity appears for that SHA or the human asks to wait.

### Stop Conditions

The follow-up loop can stop when one of these is true:

- required checks pass and no advisory bot/agent review-start activity appears after the 30-second startup wait
- for Step 2 planning PRs, the initial follow-up poll completed and no other stop-condition work remains
- for Step 3 execution PRs, the bounded CI-settling window ended after the initial follow-up poll plus up to two additional 30-second CI-settling polls, and required checks are still pending with no advisory bot/agent review-start activity
- required checks pass and any available in-scope advisory bot/agent findings were fixed in the current PR while larger follow-ups were deferred through the documented plan-or-issue split
- CI is blocked or not actionable, and the blocker is documented
- advisory bot/agent review remains pending after review-start activity and the 7-minute wait budget, and the timeout state is documented
- the compact follow-up helper is missing or fails, and the failure reason is documented
- the user explicitly asks to stop waiting

Wait for GitHub PR review approval before merging into `main`.

When the PR reaches a human-review handoff, end with a compact separate-session prompt that future follow-up for the same PR should start in a new session. Keep it to about 3 lines and include the PR number, branch name, and requested follow-up scope.

## Verification First Principle

Human review happens **after** mechanical tests and "visual" verification data are ready. The PR should contain enough evidence that a reviewer can confirm correctness without running the code themselves.

## After Merge

After the PR is merged, return to the appropriate workflow step:

- If the project plan needs updating → Step 1 (Project Plan)
- If the next task needs planning → Step 2 (Execution Plan)
- If a plan is ready for implementation → Step 3 (Execution)

Repeat steps 1–3 until the Project Plan is complete.
