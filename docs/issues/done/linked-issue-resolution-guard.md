# Guard Resolved Local Issue Moves with Explicit Plan Linkage

## Summary

The workflow already says resolved local issues must move from `docs/issues/` to `docs/issues/done/`, but it does not consistently make the relationship between a concrete execution plan and the issue it resolves machine-visible.

That leaves a gap:

- the implementation PR can clearly resolve the underlying problem
- the plan/PR can mention the issue informally
- but the final `docs/issues/... -> docs/issues/done/...` move can still be forgotten because it is not enforced as task-local completion criteria

The missed `ai-arena/docs/issues/arena-runner-log-and-persist-split.md` move after PR #120 is the concrete example that motivated this follow-up.

## Why It Matters

- Workflow rules that live only in global guidance are easy to miss at the end of a task.
- The current linter can detect a bad delete-vs-move after someone edits issue files, but it cannot detect "this plan resolved a known issue and forgot to move it".
- Reviewers and agents need a lightweight way to see which local issue files are expected to close with a given execution PR.

## Proposed Direction

- Require executable plans that resolve tracked local issues to declare them in an `Addresses:` section.
- Teach the workflow and PR template that resolved linked issues must move to `docs/issues/done/`, or the PR must explain why not.
- Extend `tools/workflow-lint.sh` with a fixable warning for execution work that declares linked local issues but does not move them to `docs/issues/done/` when the issue appears resolved by the diff.

## Non-Goal

- Do not rely on workflow-skill wording alone as the primary safeguard. Skills should follow the contract, but the enforcement should live in repo-visible workflow docs, specs, template, and linter behavior.
