---
name: review-task
description: Prepare a workflow PR, verify its conditional metadata, and perform bounded latest-head follow-up.
---

# Review Task

Use for PR creation, a later branch push, or PR/CI/advisory-review follow-up.
Read [PR and follow-up](../../AI_WORKFLOW.md#pr-and-follow-up) for lifecycle
rules and [the PR template](../../.github/PULL_REQUEST_TEMPLATE.md) for author fields.

## Prepare

1. Confirm branch type and plan mapping; ensure all applicable local quality
   gates have passed. Resolve workflow-linter `fixable` warnings or record the
   explicit exception in the PR body.
2. Use the current repository template; otherwise use the workspace, vendored,
   then workflow template. Fill only applicable fields: plan/issues, conditional
   external `Closes`, type, concise summary, uncaptured intent, verification,
   and relevant checklist/reviewer evidence.
3. Before pushing an existing PR branch, confirm that PR is open. Commit, push,
   and create or update the PR.

## Latest-head follow-up

1. Wait 30 seconds and run `scripts/gh-pr-followup poll` for the latest SHA.
   The helper's compact output is canonical; if it fails, report the failure and
   stop automatic raw-API fallback unless the human requests targeted diagnosis.
2. Inspect checks, timeline events, review summaries, and inline comments. For
   an execution PR with pending required checks, poll up to two more times after
   30 seconds each. Fix actionable in-scope CI failures, reverify, push, and
   restart from step 1 for the new head.
3. When advisory review starts for this head, poll after 3m, 2m, 1m, and 1m.
   Inspect approvals/passes as well as comment states. Triage substantive
   findings by source, location, summary, implementer's view, explanation, and
   recommendation; fix reasonable in-scope work, otherwise defer to a plan or
   issue.
4. Hand off the PR URL, worktree path, latest-head status, and grouped advisory
   triage. The human owns PR review, approval, and merge. Omit prompts asking a
   new session to review or merge the PR.

## Post-merge plan prompt

Give a short, copyable prompt for a concrete next task. Say to start after
the current PR merges, and use the actual project repo path and exact plan
file path in the prompt.

- For a plan PR that adds an executable plan, use
  `$execute-task <project-repo> <exec-plan-file-path>` for that plan.
- For an implementation PR with a scheduled executable successor, use
  `$execute-task <project-repo> <exec-plan-file-path>` for the successor.
- If the implemented plan filled a gap in a parent or originating plan that
  needs reconsideration, use
  `$plan-execution <project-repo> <originating-plan-file-path>` and specify
  what to revisit. For a removed plan, include its former path and a PR or
  Git reference that retrieves it.

Omit the prompt when there is no concrete related next task with a known plan
path. A parent marked as requiring more detail is ineligible for `$execute-task`.
