# Align actions/checkout versions across workflows

## Summary

The repository currently mixes `actions/checkout` major versions across GitHub Actions workflows.

- `.github/workflows/kb-pages.yml` uses `actions/checkout@v4`
- `.github/workflows/workflow-lint.yml` uses `actions/checkout@v6`
- `.github/workflows/shellcheck.yml` uses `actions/checkout@v6`

This is not a functional bug in the ShellCheck CI task, but it creates avoidable maintenance ambiguity for future workflow updates and review comments.

## Proposed Solution

Choose a single supported `actions/checkout` major version for this repository and update all workflows to that version in one focused change.

- inventory current workflow requirements that depend on checkout behavior
- confirm whether `v6` is the intended repo standard or whether workflows should stay on `v4` for compatibility
- update all workflow files together so the repository stops drifting by file
- document the chosen standard if the repo intentionally stays mixed for a time

## Priority

Medium. This is maintenance work rather than a blocker, but leaving mixed major versions in place invites repeated review churn and can hide real workflow compatibility differences.
