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
3. Inspect required CI/check status and continue the CI failure loop below when checks fail or expose actionable logs.
4. Inspect the PR issue timeline with `gh api repos/<owner>/<repo>/issues/<pr-number>/timeline --paginate`.
5. Detect automatic reviewer activity from timeline events, including:
   - `copilot_work_started`
   - `review_requested` events where `requested_reviewer.login` or `requested_team.name` identifies Copilot, Claude, `gh aw`, agent workflow, or another configured bot/agent reviewer
   - timeline events from bot or agent actors that indicate review work has started
6. If automatic reviewer activity has started but no final review/comments are visible yet, wait for submitted review/comments using the bounded bot-review cadence below.
7. If no automatic reviewer activity is present, do not spend the bot-review wait budget; record that no automatic review start was observed.
8. Before handoff, fetch review summaries and inline review comments with:
   - `gh pr view <pr-number> --json reviews,comments,statusCheckRollup`
   - `gh api repos/<owner>/<repo>/pulls/<pr-number>/comments --paginate`
9. Treat actionable bot/agent review comments like normal review feedback: resolve them in scope or document why they are not being acted on before calling the PR ready for handoff.
10. Before handoff, verify the PR head SHA did not change during monitoring. If it changed, restart from step 1 for the new head SHA.

Bounded bot-review wait cadence:

- First wait: 5 minutes.
- Second wait: 1 minute.
- Third wait: 1 minute.
- Total wait budget: 7 minutes across 3 polling turns.
- After each wait turn, fetch PR reviews and inline comments using the review-summary and pull-comment commands above. If review/comments were submitted, stop waiting and triage them.

If bot review has started but no review/comments have been submitted after the 7-minute budget, treat the bot-review wait as timed out and document the state in the handoff.

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

After the agent has already implemented a change, it should evaluate Copilot comments from the implementation context and present response options in the current session. Do not post that triage back to the PR unless the user explicitly asks for a PR comment.

Explicit implementation work may follow if the user asks for it or if the workflow has already authorized addressing review feedback. Larger comment-driven rewrites still count as pushes that restart the monitoring loop.

## Stop Conditions

The bounded follow-up cycle may stop when:

- all required checks pass and no bot review-start activity appears after the 30-second startup wait
- checks pass and available Copilot comments have been summarized for human review
- CI fails but cannot be fixed automatically within scope, and the blocker is documented
- Copilot review remains pending after bot review-start activity and the 7-minute wait budget, and the timeout state is documented
- the user explicitly asks to stop waiting

The handoff must say which condition was reached.
