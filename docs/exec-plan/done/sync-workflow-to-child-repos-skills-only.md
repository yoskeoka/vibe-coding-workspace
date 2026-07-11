**Execution**: Use `/execute-task` to implement this plan.

# Skip no-op workflow sync PRs for child repos

## Objective

Stop creating child-repo workflow sync PRs when the pushed `vibe-coding-workspace` commit does not change anything under `skills/`.

This matches how `setup-workspace.sh` consumes the workflow submodule: child repositories only symlink skills from `.claude/vendor/workflow/skills/`, so changes outside `skills/` do not need to fan out.

## Background

The current sync workflow compares the child repo's submodule commit directly against `github.sha` and opens PRs whenever they differ. That produces PRs for unrelated workflow repo changes, even though child repos do not consume those paths.

## Design Decisions

- Use a path-limited diff on `skills/` between the child repo's current workflow submodule commit and the pushed workflow commit.
- Keep the existing stale-PR closing and PR creation flow unchanged for commits that do touch `skills/`.
- Record the decision as an indexed file under `docs/design-decisions/adr/`.

## Code Changes

- Update `.github/workflows/sync-workflow-to-child-repos.yml` so each matrix job skips PR creation when the diff contains no `skills/` changes.
- Update the PR title/body text to describe the sync as a `skills` update.

## Spec Changes

- Add `docs/specs/workflow-sync-to-child-repos.md` describing the sync rule and skip behavior.

## Sub-tasks

- [x] [parallel] Add/update the workflow sync spec and ADR entry.
- [x] [depends on: spec] Update `.github/workflows/sync-workflow-to-child-repos.yml` to gate PR creation on `skills/` changes.
- [x] [depends on: workflow] Verify the workflow logic locally with a dry-run style inspection or script-level reasoning.
- [x] [depends on: verification] Move this plan file to `docs/exec-plan/done/`.
