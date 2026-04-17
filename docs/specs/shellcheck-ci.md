# Spec: ShellCheck CI

## Goal

Run ShellCheck in CI only for Bash scripts changed by a pull request targeting `main`. This keeps shell lint feedback focused on the scripts a PR actually touches, instead of making unrelated historical shell issues a blocker.

## Scope

- Repository: `vibe-coding-workspace`
- Trigger: GitHub Actions `pull_request` events targeting `main`
- Files considered: changed files added, modified, or renamed in the PR diff
- Non-goal: linting every historical shell script on every pull request

## Changed Bash Script Definition

A file is treated as a changed Bash script when all of the following are true:

1. The file appears in `git diff --name-only --diff-filter=AMR <base-ref>...HEAD`.
2. The file still exists in the worktree after checkout.
3. The file matches at least one of these rules:
   - path ends with `.sh`
   - path ends with `.bash`
   - first line contains a Bash shebang such as `#!/bin/bash` or `#!/usr/bin/env bash`

In this workspace, `.sh` files are treated as Bash automation scripts by policy unless this spec is intentionally narrowed in a future change.

## Helper Script

`tools/list-changed-bash-scripts.sh` is the canonical implementation for local and CI detection.

### Interface

```sh
tools/list-changed-bash-scripts.sh [base-ref]
```

### Behavior

- defaults `base-ref` to `origin/main`
- compares `${base_ref}...HEAD`
- emits one matching path per line
- ignores deleted files
- uses repository-safe shell logic for ordinary path names used in this workspace

## GitHub Actions Workflow

`.github/workflows/shellcheck.yml` defines a separate project lint workflow.

### Workflow Behavior

- checks out the PR head
- checks out with enough history to resolve the merge-base used by `<base-ref>...HEAD`
- explicitly fetches the PR base branch before diffing
- installs ShellCheck on the runner instead of assuming it is preinstalled
- calls `tools/list-changed-bash-scripts.sh` with the fetched base ref
- exits successfully with a clear message when no changed Bash scripts are found
- runs ShellCheck against every matching file
- passes `--shell=bash` for extension-based matches that do not provide their own shell directive or shebang

## Relationship to Workflow Lint

This check is separate from `.github/workflows/workflow-lint.yml`.

- `workflow-lint.yml` enforces AI workflow rules declared in `AI_WORKFLOW.md`
- ShellCheck CI is project-specific lint for shell automation scripts
- project lint remains intentionally out of scope for `tools/workflow-lint.sh`

## Local Verification

Use the helper script directly:

```sh
tools/list-changed-bash-scripts.sh origin/main
```

If `shellcheck` is installed locally, verify the helper and detected files with:

```sh
shellcheck tools/list-changed-bash-scripts.sh
tools/list-changed-bash-scripts.sh origin/main | while IFS= read -r path; do
  [ -n "$path" ] || continue
  case "$path" in
    *.sh|*.bash)
      shellcheck --shell=bash "$path"
      ;;
    *)
      shellcheck "$path"
      ;;
  esac
done
```
