# Copilot review triage should not auto-apply comments

## Summary

During PR follow-up, the current workflow/skill guidance can lead an agent to
silently apply Copilot review feedback after extracting comments from a PR.
That is too aggressive for advisory bot review.

The desired behavior is:

- extract each Copilot review comment
- include the implementer's view on the comment
- add a concise 1-2 line explanation
- decide only whether the item should be fixed in the current PR
- avoid editing files or pushing changes unless the human explicitly approves
  that action

This surfaced while handling `yoskeoka/ww#164`: Copilot pointed out ambiguity in
the worktree browser viewport spec, and I clarified the spec immediately instead
of first presenting the comment summary, implementer view, and "fix in this PR?"
decision for human review.

## Proposed Solution

Update the relevant workflow skill guidance, likely `review-task` and any
execution/PR handoff instructions that describe Copilot review handling, so
Copilot review comments are triaged before implementation.

The skill should require a compact table or list with these fields:

- comment location/link
- extracted comment summary
- implementer's view
- 1-2 line explanation
- recommendation: fix in this PR, defer, or no action

Only after that triage should the agent ask for or receive explicit human
approval before applying changes.

## Priority

Medium. This prevents advisory bot feedback from becoming unreviewed scope
creep, preserves human control over PR changes, and makes Copilot review
follow-up more auditable.
