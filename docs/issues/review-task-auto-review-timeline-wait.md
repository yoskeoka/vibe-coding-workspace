# review-task should wait for automatic reviewer timeline completion

## Summary

The `review-task` PR gate currently checks PR readiness through regular PR
state, status checks, and review/comment data, but it does not explicitly check
the GitHub issue timeline for automatic reviewer lifecycle events.

That missed an important case on PR #73:

- GitHub UI showed Copilot review-requested and Copilot started-reviewing events.
- `gh pr view <pr> --json reviewRequests` returned an empty list.
- `gh pr view <pr> --json reviews` only showed Copilot after the review was
  submitted.
- The actionable inline comments were discoverable only after the Copilot review
  completed.

The missing check was:

```bash
gh api repos/<owner>/<repo>/issues/<pr-number>/timeline --paginate
```

That timeline included:

- `review_requested` with `requested_reviewer.login=Copilot`
- `copilot_work_started`
- the later `reviewed` event once Copilot submitted its review

## Proposed Solution

Update `skills/review-task/SKILL.md` so that after CI/status checks are green,
the PR gate also inspects the issue timeline API before final handoff.

The skill should:

1. Run:

   ```bash
   gh api repos/<owner>/<repo>/issues/<pr-number>/timeline --paginate
   ```

2. Detect automatic review activity from systems such as:
   - Copilot
   - agent workflow / `gh aw`
   - Claude
   - other configured bot or agent reviewers

3. If an automatic reviewer has started but has not yet submitted its final
   review/comments, wait rather than handing off immediately.

4. After the automatic review completes, fetch both review summaries and inline
   review comments:

   ```bash
   gh pr view <pr-number> --json reviews,comments,statusCheckRollup
   gh api repos/<owner>/<repo>/pulls/<pr-number>/comments
   ```

5. Treat actionable bot/agent comments like normal review feedback: resolve or
   explicitly document why they are not being acted on before calling the PR
   ready.

## Priority

High. This is a PR-quality gate issue: failing to wait for automatic review
comments can leave known review feedback unaddressed while the agent reports
the PR as ready.
