# PR Follow-up Workflow

## Purpose

Define the workflow contract after a pull request is created or updated. PR creation is not the end of review preparation: agents must watch the new PR head long enough to collect mechanical verification results and configured advisory bot/agent review signals.

## Ownership

The `review-task` skill owns PR preparation and initial post-PR follow-up for every workflow step that routes through it.

Caller skills such as `plan-execution` and `execute-task` must not report task completion before `review-task` reaches a documented stop condition for the current PR head SHA.

This ownership applies when:

- a new PR is created
- an existing PR is updated
- a later push changes the PR branch head SHA

Each newer head SHA restarts the required CI/check inspection loop. The longer advisory bot/agent review wait budget only runs when review-start activity is observed for the latest head SHA, on initial review startup, or when the human explicitly asks to wait.

## Monitoring Loop

For a non-blocked PR creation or update flow, the minimum landing path is `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll`.

Step-specific CI settling rule:

- Step 2 (`plan-execution`) may stop after that initial poll when no other stop-condition work remains.
- Step 3 (`execute-task`) uses a bounded CI-settling window as part of the landing check. If required checks are still pending after the initial poll for the latest PR head SHA, wait another 30 seconds and poll again. If required checks are still pending after that second poll, wait a third 30-second turn and poll once more.
- Stop the extra Step 3 settling polls early when required checks finish, advisory reviewer activity starts and the advisory-review cadence takes over, the helper fails, the PR head SHA changes, or the user explicitly asks to stop waiting.

For each PR head SHA, `review-task` must:

1. Record the current PR number and head SHA.
2. Wait 30 seconds after PR creation or a later push so CI/checks and review automation have time to start.
3. Poll the PR with the compact follow-up helper when it exists:

   ```sh
   skills/review-task/scripts/gh-pr-followup poll <owner/repo> <pr-number>
   ```

   If the helper is missing or fails, report the failure reason and stop automatic follow-up for this PR head SHA. Do not automatically fall back to raw `gh` reads.
4. Inspect required CI/check status from the helper output and continue the CI failure loop below when checks fail or expose actionable logs. This CI/check inspection happens for every pushed PR head SHA.
5. Inspect the PR issue timeline from the helper output.
6. Inspect review summaries and inline review comments from the helper output.
7. Detect advisory bot/agent reviewer activity from timeline events, including:
   - `copilot_work_started`
   - `review_requested` events where `requested_reviewer.login` or `requested_team.name` identifies Copilot, Claude, `gh aw`, agent workflow, or another configured bot/agent reviewer
   - timeline events from bot or agent actors that indicate review work has started
8. If advisory reviewer activity has started for the latest head SHA but no final review/comments are visible yet, wait for submitted review/comments using the bounded advisory-review cadence below.
9. If no advisory reviewer activity is present, do not spend the advisory-review wait budget; record that no advisory review start was observed.
10. Before handoff, use the latest helper output for review summaries and inline review comments. If the helper failed, hand off the failure reason and tell the user that PR follow-up should be checked later rather than spending extra context on raw GitHub API output.
11. Triage substantive advisory bot/agent findings for human review before any comment-driven branch mutation. Do not edit files, apply suggestions, commit, or push in response to advisory bot/agent findings unless the human explicitly approves that specific follow-up or a prior human instruction already authorized implementing that exact review feedback.
12. Before handoff, verify the PR head SHA did not change during monitoring. If it changed, restart from step 1 for the new head SHA.

For Step 3 execution PRs, if step 4 shows required checks still pending and no higher-priority stop condition has already been reached, perform up to two additional `wait 30 seconds -> poll compact helper` turns before handoff.

Bounded advisory-review wait cadence:

- First wait: 5 minutes.
- Second wait: 1 minute.
- Third wait: 1 minute.
- Total wait budget: 7 minutes across 3 polling turns.
- After each wait turn, poll with the compact helper. If review/comments were submitted, stop waiting and triage them. If the helper fails, stop the automatic wait loop and report the failure reason.

If advisory bot/agent review has started but no review/comments have been submitted after the 7-minute budget, treat the advisory-review wait as timed out and document the state in the handoff.

Polling-style waiting should be delegated to a low-cost subagent only when the platform supports delegation and the current session explicitly authorizes subagent use. The main agent remains responsible for deciding what to fix, what to defer, and what to report.

## Compact Polling Helper

`skills/review-task/scripts/gh-pr-followup` is the preferred polling path for repeated post-PR checks because it emits compact, decision-oriented JSON instead of raw GitHub API objects.

Interface:

```sh
skills/review-task/scripts/gh-pr-followup poll <owner/repo> <pr-number>
```

Behavior:

- Reads PR head, review decision, review summaries, issue timeline, inline review comments, and status checks through GitHub CLI.
- Stores non-canonical local polling markers under `.local/gh-pr-followup/`.
- Keys state by owner, repo, and PR number using a filesystem-safe filename.
- Records at least:
  - `head_sha`
  - `last_timeline_event_id`
  - `last_review_comment_id`
  - `last_checked_at`
- Resets timeline and inline-comment markers when the PR head SHA changes.
- Emits only timeline events and inline comments newer than the stored markers for the current head SHA.
- Updates the state file only after all GitHub reads and JSON compaction succeed.
- Supports `GH_BIN` and `GH_PR_FOLLOWUP_STATE_DIR` environment overrides for local verification.

If the helper is missing or fails during automatic PR follow-up, agents must not switch to raw `gh api` polling by default. The low-token behavior is to report the helper failure, stop the automatic CI/advisory-review monitoring loop for that head SHA, and leave later PR inspection to a human or a separately requested follow-up. Raw `gh` inspection remains available only when the user explicitly asks for it or when a targeted, bounded command is needed to diagnose the helper itself.

The helper output must include:

- `repo`
- `pr`
- `head_sha`
- `review_decision`
- `checks`
- `reviews`
- `timeline_events`
- `inline_comments`
- `state`

Compact check entries include only:

- `name`
- `workflow`
- `status`
- `conclusion`
- `details_url`

Compact review entries include only:

- `id`
- `state`
- `user`
- `body`
- `submitted_at`
- `commit_id`
- `html_url`

Compact timeline event entries include only fields needed to detect review-start and review-complete signals:

- `id`
- `event`
- `created_at`
- `actor`
- `reviewer`
- `team`
- `app`
- `commit_id`
- `review_state`

Compact inline review comment entries include only:

- `id`
- `path`
- `line`
- `user`
- `body`
- `created_at`
- `updated_at`
- `commit_id`
- `html_url`

The helper is not the source of truth for PR state. Its `.local/` markers are derived cache data that can be deleted; deletion only causes the next poll to return currently visible timeline events and inline comments again.

## CI Failure Policy

CI failures are mechanical verification failures. When logs are available and the failure is actionable within the branch scope, the agent should:

1. Inspect the failing check logs.
2. Fix the issue in the same branch.
3. Re-run the relevant local verification.
4. Commit and push the fix.
5. Restart the post-push monitoring loop for the new PR head SHA.

If the failure is not actionable from available logs, is outside scope, or depends on external infrastructure, the handoff must state the blocker and include the relevant check/log context.

## Advisory Bot/Agent Review Policy

Copilot, Claude, `gh aw`, agent workflow reviews, and other configured bot/agent review comments are advisory review input, not mechanical verification failures.

Agents must not silently auto-apply advisory bot/agent suggestions. Passing or approving advisory bot/agent checks can still contain substantive observations in review bodies, so review summaries must be inspected even when the overall state is not blocking.

For each substantive advisory finding, the handoff must group findings by source reviewer/workflow and include:

- source reviewer/workflow
- comment location or link
- extracted comment summary
- implementer's view
- concise 1-2 line explanation
- recommendation: fix in this PR, defer, or no action

After the agent has already implemented a change, it should evaluate advisory bot/agent findings from the implementation context and present response options in the current session. Do not post that triage back to the PR unless the user explicitly asks for a PR comment.

Explicit implementation work may follow only if the human asks for it or a prior human instruction already authorized that specific review-feedback implementation work. Larger comment-driven rewrites still count as pushes that restart required CI/check inspection for the new head SHA and may restart the advisory wait only when new review-start activity appears or the human asks to wait.

## Stop Conditions

The bounded follow-up cycle may stop when:

- all required checks pass and no advisory bot/agent review-start activity appears after the 30-second startup wait
- for Step 2 planning PRs, the initial follow-up poll completed and no other stop-condition work remains
- for Step 3 execution PRs, the bounded CI-settling window ended after the initial follow-up poll plus up to two additional 30-second CI-settling polls, and required checks are still pending with no advisory bot/agent review-start activity
- checks pass and available advisory bot/agent findings have been summarized for human review
- CI fails but cannot be fixed automatically within scope, and the blocker is documented
- advisory bot/agent review remains pending after review-start activity and the 7-minute wait budget, and the timeout state is documented
- the compact follow-up helper is missing or fails, and the failure reason is documented
- the user explicitly asks to stop waiting

The handoff must say which condition was reached.

Caller skills that routed work into `review-task` must treat one of these stop conditions as the completion boundary for the current PR head SHA rather than reporting success immediately after local edits, local verification, commit creation, or PR creation.
