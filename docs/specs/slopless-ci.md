# Spec: Slopless CI

## Goal

Run `slopless` on eligible changed Markdown in pull requests to `main`.
Show CI findings as GitHub Actions warnings and keep one PR comment up to date.
Use the same plan-file exclusion in the local post-edit hook.

## Scope

- Repositories: `vibe-coding-workspace` and managed child repositories that run Slopless CI
- Trigger: GitHub Actions `pull_request` events targeting `main`
- Common exclusion: Markdown under `docs/exec-plan/` is never linted by CI or the local hook
- Trigger exclusion: a change only under `docs/exec-plan/` does not trigger Slopless CI
- Workspace files considered: changed Markdown under `docs/specs/`, `docs/design-decisions/`, `docs/development/`, `docs/kb/`, and `docs/references/`
- Child project files considered: changed Markdown in each project's current Slopless scope, excluding `docs/exec-plan/`
- Non-goal: linting every historical Markdown file on every pull request
- Non-goal: reproducing the local Codex hook behavior inside CI

## Changed Markdown Definition

A file counts as eligible changed Markdown when all of these are true:

1. The file appears in `git diff --name-only --diff-filter=AMR <base-ref>...HEAD`.
2. The file still exists in the worktree after checkout.
3. The root-relative path is not under `docs/exec-plan/`.
4. The path is within the repository's Slopless scope:
   - In `vibe-coding-workspace`, it is under one of these directories:
     - `docs/specs/`
     - `docs/design-decisions/`
     - `docs/development/`
     - `docs/kb/`
     - `docs/references/`
   - In managed child projects, it matches that project's existing workflow scope.
5. The path ends with `.md`.
6. The file content does not contain characters from the helper's Japanese-writing ranges:
   - Hiragana and fullwidth Katakana (`U+3040`-`U+30FF`)
   - Katakana Phonetic Extensions (`U+31F0`-`U+31FF`)
   - CJK Unified Ideographs Extension A (`U+3400`-`U+4DBF`)
   - CJK Unified Ideographs (`U+4E00`-`U+9FFF`)
   - CJK Compatibility Ideographs (`U+F900`-`U+FAFF`)
   - Halfwidth Katakana (`U+FF65`-`U+FF9F`)

## Helper Script

In `vibe-coding-workspace`, `tools/list-changed-markdown.sh` is the source of
truth for local and CI detection. Child project workflows apply the common
exclusion in both their trigger and candidate selection.

### Interface

```sh
tools/list-changed-markdown.sh [base-ref]
```

### Behavior

- defaults `base-ref` to `origin/main`
- compares `${base_ref}...HEAD`
- emits one matching path per line
- ignores deleted files
- excludes all paths under `docs/exec-plan/`, even if a future scope change adds a broader Markdown pattern
- limits eligibility to Markdown under `docs/specs/`, `docs/design-decisions/`, `docs/development/`, `docs/kb/`, and `docs/references/`
- skips files whose content contains any character from those Japanese-writing ranges, even if they also contain ASCII, punctuation, or emoji

## Local Post-Edit Check

The local post-edit Slopless check follows the same `docs/exec-plan/` exclusion.
An execution-plan file is skipped even when an edit event includes it alongside
other Markdown files. The shared hook change must reach child projects through
the workflow submodule sync.

## GitHub Actions Workflow

`.github/workflows/slopless.yml` is the CI entrypoint.

### Workflow Behavior

- triggers on pull requests targeting `main`
- starts only for the scoped long-lived Markdown paths and the workflow/helper/spec files
- excludes `docs/exec-plan/**` from triggers and candidate sets in managed child workflows, while preserving each project's current path scope
- uses the repository-standard `actions/checkout` reference managed through `pinact`
- provisions a pinned Node runtime before calling `npm view` or `npx`
- checks out with full history (`fetch-depth: 0`) so `<base-ref>...HEAD` resolves cleanly
- explicitly fetches the PR base branch before diffing
- calls `tools/list-changed-markdown.sh` with the fetched base ref
- exits successfully with a clear message when no changed Markdown files are found
- runs one `slopless` command over the full changed Markdown file set
- pins `slopless` to version `0.2.10`
- verifies the npm `gitHead` for `slopless@0.2.10` matches commit `c40c40f3127d0c61cbfc1c34cacf0a5f49ed7e26` before linting
- converts `slopless` findings into GitHub Actions warning annotations
- appends rule-specific repair hints to each warning annotation message
- writes a step summary with the file count, finding count, and a compact findings table
- upserts one PR comment with a stable marker instead of posting a new comment on every rerun or push
- includes the latest file count, findings table, and rule-specific repair hints in that PR comment
- paginates PR comment reads before deciding whether to create or update the marker comment
- uses explicit timeouts for both the `slopless` subprocess and GitHub API requests
- requests `contents: read`, `issues: write`, and `pull-requests: write` permissions for PR comment upsert
- logs the GitHub API response body for HTTP failures so permission errors remain diagnosable
- treats `slopless` exit code `0` with empty stdout as a valid zero-findings result
- treats missing `filePath` fields in tool output defensively so the workflow never annotates the repository root as the target file
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

```bash
actual_git_head="$(npm view "slopless@0.2.10" gitHead | tr -d '\n')"
test "$actual_git_head" = "c40c40f3127d0c61cbfc1c34cacf0a5f49ed7e26"
mapfile -t files < <(tools/list-changed-markdown.sh origin/main)
[ "${#files[@]}" -eq 0 ] || npx --yes --package=slopless@0.2.10 slopless "${files[@]}" || true
```
