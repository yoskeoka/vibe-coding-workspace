# Historical workflow artifact retention

## Summary

Resolved execution plans and local issues remain in the checked-out repository
under `done/`. They are not normally useful to implementation agents, but they
are returned by ordinary repository searches and inflate task context.

## Desired outcome

Keep only active plans and unresolved local issues in the working tree. Treat
the plan PR, implementation PR, and Git history as the audit trail for a
resolved artifact. Preserve durable decisions in ADRs, specs, code comments, or
the project plan rather than in completed task trackers.

## Scope

- Replace the move-to-`done/` lifecycle with deletion after successful
  implementation.
- Update the workflow contract, templates, skills, and linter together.
- Remove existing completed artifacts from `docs/exec-plan/done/` and
  `docs/issues/done/`.

## Non-goals

- Put full plan text into commit messages.
- Delete active plans or unresolved issues.
- Replace GitHub PRs or Git history as the completion/audit trail.

