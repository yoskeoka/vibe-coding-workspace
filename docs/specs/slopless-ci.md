# Spec: Slopless CI

## Goal

Run `slopless` in CI for changed Markdown in pull requests to `main`.
Show findings as GitHub Actions warnings.
Keep one PR comment up to date with the latest findings and repair hints.

## Scope

- Repository: `vibe-coding-workspace`
- Trigger: GitHub Actions `pull_request` events targeting `main`
- Files considered: changed `*.md` files added, modified, or renamed in the PR diff
- Non-goal: linting every historical Markdown file on every pull request
- Non-goal: reproducing the local Codex hook behavior inside CI

## Changed Markdown Definition

A file counts as changed Markdown when all of these are true:

1. The file appears in `git diff --name-only --diff-filter=AMR <base-ref>...HEAD`.
2. The file still exists in the worktree after checkout.
3. The path ends with `.md`.

CI does not filter by language.
Any changed `*.md` file is eligible.

## Helper Script

`tools/list-changed-markdown.sh` is the source of truth for local and CI detection.

### Interface

```sh
tools/list-changed-markdown.sh [base-ref]
```

### Behavior

- defaults `base-ref` to `origin/main`
- compares `${base_ref}...HEAD`
- emits one matching path per line
- ignores deleted files
- treats every changed `*.md` path as eligible

## GitHub Actions Workflow

`.github/workflows/slopless.yml` is the CI entrypoint.

### Workflow Behavior

- triggers on pull requests targeting `main`
- starts only for Markdown-related paths and the workflow/helper/spec files
- uses the repository-standard `actions/checkout` reference managed through `pinact`
- checks out with full history (`fetch-depth: 0`) so `<base-ref>...HEAD` resolves cleanly
- explicitly fetches the PR base branch before diffing
- calls `tools/list-changed-markdown.sh` with the fetched base ref
- exits successfully with a clear message when no changed Markdown files are found
- runs `slopless` against every changed Markdown file
- pins `slopless` to version `0.2.10`
- verifies the npm `gitHead` for `slopless@0.2.10` matches commit `c40c40f3127d0c61cbfc1c34cacf0a5f49ed7e26` before linting
- converts `slopless` findings into GitHub Actions warning annotations
- appends rule-specific repair hints to each warning annotation message
- writes a step summary with the file count, finding count, and a compact findings table
- upserts one PR comment with a stable marker instead of posting a new comment on every rerun or push
- includes the latest file count, findings table, and rule-specific repair hints in that PR comment
- keeps the job green when `slopless` reports prose findings, but fails the job if the pinned package cannot be verified or the tool output is malformed

## Reporting Contract

Each `slopless` finding should appear as a GitHub Actions warning annotation with:

- file path
- line
- column
- `ruleId` as the warning title
- lint message text
- a short repair hint derived from the matching readability rule

The step summary should include:

- number of checked files
- number of findings
- per-finding table when findings exist

The workflow should also keep one PR comment per pull request:

- comment body includes a stable hidden marker so reruns can find and update it
- comment body includes the latest file count, finding count, findings table, and per-rule repair hints
- reruns after new pushes update the existing comment instead of creating another one

## Local Verification

Use the helper directly:

```sh
tools/list-changed-markdown.sh origin/main
```

For a focused local CI-style check, run:

```sh
actual_git_head="$(npm view "slopless@0.2.10" gitHead | tr -d '\n')"
test "$actual_git_head" = "c40c40f3127d0c61cbfc1c34cacf0a5f49ed7e26"
tools/list-changed-markdown.sh origin/main | while IFS= read -r path; do
  [ -n "$path" ] || continue
  npx --yes --package=slopless@0.2.10 slopless "$path" || true
done
```
