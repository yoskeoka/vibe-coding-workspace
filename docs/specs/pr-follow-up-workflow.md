# PR Follow-up Workflow

## Purpose

Define the workflow contract after a pull request is created or updated. PR creation is not the end of review preparation: agents must watch the new PR head long enough to collect mechanical verification results and configured automatic AI review signals.

## Ownership

The `review-task` skill owns PR preparation and initial post-PR follow-up for every workflow step that routes through it.

This ownership applies when:

- a new PR is created
- an existing PR is updated
- a later push changes the PR branch head SHA

Each newer head SHA restarts the bounded follow-up cycle.

## Monitoring Loop

For each PR head SHA, `review-task` must:

1. Record the current PR number and head SHA.
2. Wait 30 seconds after PR creation or a later push so CI/checks and review automation have time to start.
3. Poll the PR with the compact follow-up helper when it exists:

   ```sh
   skills/review-task/scripts/gh-pr-followup poll <owner/repo> <pr-number>
   ```

   If the helper is missing or fails, report the failure reason and stop automatic follow-up for this PR head SHA. Do not automatically fall back to raw `gh` reads.
4. Inspect required CI/check status from the helper output and continue the CI failure loop below when checks fail or expose actionable logs.
5. Inspect the PR issue timeline from the helper output.
6. Inspect review summaries and inline review comments from the helper output.
7. Detect automatic reviewer activity from timeline events, including:
   - `copilot_work_started`
   - `review_requested` events where `requested_reviewer.login` or `requested_team.name` identifies Copilot, Claude, `gh aw`, agent workflow, or another configured bot/agent reviewer
   - timeline events from bot or agent actors that indicate review work has started
8. If automatic reviewer activity has started but no final review/comments are visible yet, wait for submitted review/comments using the bounded bot-review cadence below.
9. If no automatic reviewer activity is present, do not spend the bot-review wait budget; record that no automatic review start was observed.
10. Before handoff, use the latest helper output for review summaries and inline review comments. If the helper failed, hand off the failure reason and tell the user that PR follow-up should be checked later rather than spending extra context on raw GitHub API output.
11. Treat actionable bot/agent review comments like normal review feedback: resolve them in scope or document why they are not being acted on before calling the PR ready for handoff.
12. Before handoff, verify the PR head SHA did not change during monitoring. If it changed, restart from step 1 for the new head SHA.

Bounded bot-review wait cadence:

- First wait: 5 minutes.
- Second wait: 1 minute.
- Third wait: 1 minute.
- Total wait budget: 7 minutes across 3 polling turns.
- After each wait turn, poll with the compact helper. If review/comments were submitted, stop waiting and triage them. If the helper fails, stop the automatic wait loop and report the failure reason.

If bot review has started but no review/comments have been submitted after the 7-minute budget, treat the bot-review wait as timed out and document the state in the handoff.

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

If the helper is missing or fails during automatic PR follow-up, agents must not switch to raw `gh api` polling by default. The low-token behavior is to report the helper failure, stop the automatic CI/Copilot monitoring loop for that head SHA, and leave later PR inspection to a human or a separately requested follow-up. Raw `gh` inspection remains available only when the user explicitly asks for it or when a targeted, bounded command is needed to diagnose the helper itself.

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

## Copilot Review Policy

GitHub Copilot review comments are advisory review input, not mechanical verification failures.

Agents must not silently auto-apply Copilot suggestions. For each substantive Copilot comment, the handoff should summarize:

- the comment or requested change
- whether action is recommended
- whether the suggestion can be applied as-is or needs adaptation
- whether the item should be deferred into a separate issue
- suggested response or implementation options

After the agent has already implemented a change, it should evaluate Copilot comments from the implementation context and present response options in the current session. Do not post that triage back to the PR unless the user explicitly asks for a PR comment.

Explicit implementation work may follow if the user asks for it or if the workflow has already authorized addressing review feedback. Larger comment-driven rewrites still count as pushes that restart the monitoring loop.

## Stop Conditions

The bounded follow-up cycle may stop when:

- all required checks pass and no bot review-start activity appears after the 30-second startup wait
- checks pass and available Copilot comments have been summarized for human review
- CI fails but cannot be fixed automatically within scope, and the blocker is documented
- Copilot review remains pending after bot review-start activity and the 7-minute wait budget, and the timeout state is documented
- the compact follow-up helper is missing or fails, and the failure reason is documented
- the user explicitly asks to stop waiting

The handoff must say which condition was reached.
