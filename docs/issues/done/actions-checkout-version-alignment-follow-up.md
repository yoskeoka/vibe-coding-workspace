# Align actions/checkout versions across workflows

## Summary

Resolved by standardizing repository GitHub Actions workflows on `actions/checkout@v6`.

Before this fix, `.github/workflows/kb-pages.yml` used checkout major 4 while
`.github/workflows/workflow-lint.yml` and `.github/workflows/shellcheck.yml`
used checkout major 6.

This was not a functional bug in the ShellCheck CI task, but it created
avoidable maintenance ambiguity for future workflow updates and review comments.

## Resolution

The repository now uses `actions/checkout@v6` in:

- `.github/workflows/kb-pages.yml`
- `.github/workflows/workflow-lint.yml`
- `.github/workflows/shellcheck.yml`

The chosen standard is documented in the relevant workflow specs.

## Priority

Medium. This is maintenance work rather than a blocker, but leaving mixed major versions in place invites repeated review churn and can hide real workflow compatibility differences.
