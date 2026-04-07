# Workflow Sync to Child Repos

## Purpose

The `Sync Workflow to Child Repos` GitHub Actions workflow keeps the `.claude/vendor/workflow` submodule in child repositories aligned with workflow changes that matter to them.

`setup-workspace.sh` only exposes workflow content through the `skills/` directory, so child repositories only need PRs when workflow changes touch `skills/`.

## Sync Rule

For each child repository:

1. Read the current workflow submodule commit recorded in the child repo.
2. Compare that commit against the pushed commit in `vibe-coding-workspace`.
3. If the commit range contains no changes under `skills/`, skip PR creation for that child repo.
4. If the commit range includes at least one change under `skills/`, continue with the existing sync flow:
   - close stale `workflow-sync` PRs
   - create a new PR that updates `.claude/vendor/workflow`

## PR Expectations

- PRs are created only for workflow commits that change `skills/`.
- PR titles and bodies should describe the update as a `skills` sync, not a generic full-submodule refresh.
- Skipped repos should produce a clear log message explaining that no `skills/` changes were found.
