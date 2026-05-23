# Spec: Workflow Sync to Child Repos

## Purpose

The `Sync Workflow to Child Repos` GitHub Actions workflow keeps the `.claude/vendor/workflow` submodule in child repositories aligned with workflow changes that matter to them.

Child repositories consume workflow changes through two surfaces:

- the vendored workflow submodule itself (`.claude/vendor/workflow`)
- copied runtime assets installed from that submodule into the child repo:
  - `.githooks/pre-push`
  - `tools/workflow-lint.sh`
  - `.github/workflows/workflow-lint.yml`

## Sync Rule

The workflow matrix is the registry of managed child repositories that should receive automated workflow-sync PRs. When the workspace adds a new child repository that uses `setup-workspace.sh`, update `.github/workflows/sync-workflow-to-child-repos.yml` in the same change if that repository should receive this sync behavior.

For each child repository:

1. Read the current workflow submodule commit recorded in the child repo.
2. Compare that commit against the pushed commit in `vibe-coding-workspace`.
3. Evaluate whether the commit range touches any documented child-consumed workflow surface:
   - `skills/`
   - `setup-workspace.sh`
   - `tools/install-hooks.sh`
   - `.githooks/pre-push`
   - `tools/workflow-lint.sh`
   - `.github/workflows/workflow-lint.yml`
4. If the commit range contains none of those paths, skip PR creation for that child repo.
5. If the commit range includes at least one of those paths, continue with the existing sync flow:
   - close stale `workflow-sync` PRs
   - create a new PR that updates `.claude/vendor/workflow`
   - run `tools/install-hooks.sh` from the updated submodule against the child repo so the copied runtime assets stay aligned with the new workflow commit without re-running full submodule/bootstrap setup

## PR Expectations

- PRs are created only for workflow commits that change documented child-consumed workflow surfaces.
- PR titles and bodies should describe the update as a workflow asset sync, not a generic full-submodule refresh.
- PR bodies should include the source commit subject as the human-readable summary of the change.
- PR bodies should include the relevant workflow commit messages from the synced range so reviewers can judge merge priority without opening the source repository history.
- PR bodies should include a `diff stat` section so reviewers can confirm the touched files at a glance.
- Skipped repos should produce a clear log message explaining that no child-consumed workflow changes were found.
- Before creating the PR, the workflow must ensure the `workflow-sync` label exists in the child repository so label assignment cannot fail on a fresh repo.
- Re-running the workflow for the same source commit must be safe even if a previous attempt already pushed `workflow-sync/<short-sha>` but failed before PR creation.
- When reusing an existing `workflow-sync/<short-sha>` branch, the workflow should update it with a lease-protected force push so it does not silently overwrite unexpected remote changes.
