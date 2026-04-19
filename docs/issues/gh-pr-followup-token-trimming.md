# Reduce PR follow-up token usage

## Summary

`review-task` currently instructs agents to inspect PR checks, issue timeline events, PR reviews, and inline comments with raw `gh pr view` / `gh api` calls. In practice this can return large JSON payloads that are mostly irrelevant to the handoff:

- issue timeline events include full cross-referenced issue and repository objects
- inline PR comments include full `diff_hunk`, full user objects, reactions, and link objects
- repeated polling returns the same previously-seen timeline and comment entries
- status check rollups are needed, but only a small subset of each check object matters

This wastes context during the bounded post-PR follow-up loop. A recent `ww` PR follow-up spent tens of thousands of tokens mostly on raw timeline/comment JSON rather than on decisions.

## Proposed Solution

Keep the first iteration as a shell wrapper that can live inside the workflow/skill boundary and be called by `review-task` instead of raw `gh api` commands.

Add a script such as:

```sh
tools/gh-pr-followup poll <owner/repo> <pr-number>
```

or, if skill-local tooling is preferred:

```sh
skills/review-task/scripts/gh-pr-followup poll <owner/repo> <pr-number>
```

The script should:

1. Store local non-canonical polling state under `.local/gh-pr-followup/`, keyed by owner, repo, and PR number.
2. Record at least:
   - `head_sha`
   - `last_timeline_event_id`
   - `last_review_comment_id`
   - `last_checked_at`
3. Use `gh ... --jq` to emit only AI-relevant fields.
4. Filter out timeline and review-comment entries that are older than or equal to the stored last-seen marker before returning output to the agent.
5. Reset the markers when the PR head SHA changes, so review comments and check status are interpreted against the current head.
6. Return compact JSON or line-oriented summaries that are safe to paste directly into agent context.

Suggested timeline query shape:

```sh
gh api "repos/$repo/issues/$pr/timeline" --paginate \
  --jq '.[] | select(
    .event == "copilot_work_started" or
    .event == "review_requested" or
    .event == "reviewed" or
    .event == "committed"
  ) | {
    id,
    event,
    created_at,
    actor: .actor.login,
    reviewer: .requested_reviewer.login,
    app: .performed_via_github_app.slug,
    commit_id,
    review_state: .state
  }'
```

Suggested inline comment query shape:

```sh
gh api "repos/$repo/pulls/$pr/comments" --paginate \
  --jq '.[] | {
    id,
    path,
    line,
    user: .user.login,
    body,
    created_at,
    updated_at,
    commit_id,
    html_url
  }'
```

Suggested status query shape:

```sh
gh pr view "$pr" --repo "$repo" \
  --json headRefOid,statusCheckRollup,reviewDecision \
  --jq '{
    head_sha: .headRefOid,
    review_decision: .reviewDecision,
    checks: [
      .statusCheckRollup[] | {
        name,
        workflow: .workflowName,
        status,
        conclusion,
        details_url: .detailsUrl
      }
    ]
  }'
```

The `review-task` skill should then be updated to prefer this wrapper for post-PR follow-up, and to fall back to raw `gh` commands only when the wrapper is missing or fails.

## Priority

Medium. This does not change product behavior, but it directly improves agent reliability by reducing context churn during PR follow-up. It also makes repeated CI/Copilot polling cheaper and less error-prone because the agent sees only new events and actionable review comments.
