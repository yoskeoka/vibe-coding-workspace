# Bot/agent review triage should not auto-apply comments

## Summary

During PR follow-up, the current workflow/skill guidance can lead an agent to
silently apply Copilot, `gh aw`, agent workflow, or other advisory bot review
feedback after extracting comments from a PR. That is too aggressive for
advisory review signals.

The desired behavior is:

- extract each bot/agent review comment
- include the implementer's view on the comment
- add a concise 1-2 line explanation
- decide only whether the item should be fixed in the current PR
- avoid editing files or pushing changes unless the human explicitly approves
  that action

This surfaced while handling `yoskeoka/ww#164`: Copilot pointed out ambiguity in
the worktree browser viewport spec, and I clarified the spec immediately instead
of first presenting the comment summary, implementer view, and "fix in this PR?"
decision for human review.

The same issue applies to non-Copilot automation. For example, `ww` PRs may run
agent workflow reviews in addition to GitHub Copilot review, and those comments
should be triaged with the same human-review-first flow.

Agentic workflow reviews may appear as GitHub Actions checks and still submit
review bodies that say "Approve" overall while containing minor observations,
notes, or non-blocking suggestions that are worth fixing. Those observations
must be extracted and summarized for the human just like Copilot inline review
comments. A passing check or approving review state is not enough by itself;
the review body still needs to be inspected for substantive advisory findings.

There is also a separate follow-up nuance for repeated pushes:

- after the initial PR push, it is useful to wait for configured bot/agent
  review startup and collect review comments
- after later pushes, the agent should always check normal CI/required checks
  for the new head SHA
- after later pushes, the agent should not spend the longer bot/agent review
  wait budget unless new review-start activity appears for the latest head SHA
  or the human explicitly asks to wait
- advisory bot/agent checks that are not required CI should be reported as
  pending/skipped/pass/fail, but they should not block handoff by default on
  second and later pushes

## Proposed Solution

Update the relevant workflow skill guidance, likely `review-task` and any
execution/PR handoff instructions that describe bot/agent review handling, so
all advisory bot/agent review comments are triaged before implementation.

The skill should require a compact table or list with these fields:

- source reviewer/workflow, such as Copilot or `Spec/Code Sync Check`
- comment location/link
- extracted comment summary
- implementer's view
- 1-2 line explanation
- recommendation: fix in this PR, defer, or no action

The handoff should group findings by source, for example:

- Copilot found three items:
  - issue 1
  - issue 2
  - issue 3
- Agentic workflow `Spec/Code Sync Check` found one item, even though its
  overall review decision was approve:
  - issue 4

Only after that triage should the agent ask for or receive explicit human
approval before applying changes.

The skill should also distinguish two loops:

- CI/check loop: required mechanical checks for the latest head SHA, always
  inspected after every push
- advisory review loop: bot/agent review comments, waited for on initial PR
  review startup or when new review-start activity appears for the latest head
  SHA, but not waited on indefinitely after every subsequent push

## Priority

Medium. This prevents advisory bot feedback from becoming unreviewed scope
creep, preserves human control over PR changes, and makes bot/agent review
follow-up more auditable.
