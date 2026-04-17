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
2. Wait for required CI/check runs to settle, using the bounded wait cadence below.
3. Inspect PR timeline/review data for GitHub Copilot auto-review activity.
4. If Copilot activity is present, wait for submitted review comments using the same bounded wait cadence.
5. If no Copilot activity or submitted review appears within the bounded interval, stop waiting for Copilot and continue the handoff.
6. Before handoff, verify the PR head SHA did not change during monitoring. If it changed, restart from step 1 for the new head SHA.

Bounded wait cadence:

- First wait: 5 minutes.
- Second wait: 1 minute.
- Third wait: 1 minute.
- Total wait budget: 7 minutes across 3 polling turns.

If required checks are still pending, have not started, or Copilot has started but not submitted a review after the 7-minute budget, treat the wait as blocked or timed out and document the state in the handoff.

Polling-style waiting should be delegated to a low-cost subagent only when the platform supports delegation and the current session explicitly authorizes subagent use. The main agent remains responsible for deciding what to fix, what to defer, and what to report.

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

Explicit implementation work may follow if the user asks for it or if the workflow has already authorized addressing review feedback. Larger comment-driven rewrites still count as pushes that restart the monitoring loop.

## Stop Conditions

The bounded follow-up cycle may stop when:

- all required checks pass and no Copilot activity appears within the 7-minute wait budget
- checks pass and available Copilot comments have been summarized for human review
- CI fails but cannot be fixed automatically within scope, and the blocker is documented
- CI/checks or Copilot review remain pending after the 7-minute wait budget, and the timeout state is documented
- the user explicitly asks to stop waiting

The handoff must say which condition was reached.
