# Spec: Workflow Sync to Child Repos

## Purpose

The `Sync Workflow to Child Repos` GitHub Actions workflow keeps the `.claude/vendor/workflow` submodule in child repositories aligned with workflow changes that matter to them.

`setup-workspace.sh` only exposes workflow content through the `skills/` directory, so child repositories only need PRs when workflow changes touch `skills/`.

## Sync Rule

The workflow matrix is the registry of managed child repositories that should receive automated workflow-sync PRs. When the workspace adds a new child repository that uses `setup-workspace.sh`, update `.github/workflows/sync-workflow-to-child-repos.yml` in the same change if that repository should receive this sync behavior.

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
- PR bodies should include the source commit subject as the human-readable summary of the change.
- PR bodies should include the relevant workflow commit messages from the synced range so reviewers can judge merge priority without opening the source repository history.
- PR bodies should include a `diff stat` section so reviewers can confirm the touched files at a glance.
- Skipped repos should produce a clear log message explaining that no `skills/` changes were found.
- Before creating the PR, the workflow must ensure the `workflow-sync` label exists in the child repository so label assignment cannot fail on a fresh repo.
- Re-running the workflow for the same source commit must be safe even if a previous attempt already pushed `workflow-sync/<short-sha>` but failed before PR creation.
- When reusing an existing `workflow-sync/<short-sha>` branch, the workflow should update it with a lease-protected force push so it does not silently overwrite unexpected remote changes.
