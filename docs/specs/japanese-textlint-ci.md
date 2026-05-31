# Spec: Japanese textlint CI

## Goal

Run `textlint` in CI for changed Japanese Markdown in pull requests to `main`.
Show findings as GitHub Actions warnings.
Keep one PR comment up to date with the latest findings and replacement hints.

## Scope

- Repository: `vibe-coding-workspace`
- Trigger: GitHub Actions `pull_request` events targeting `main`
- Files considered: changed `*.md` files anywhere in the repository
- Non-goal: linting every historical Markdown file on every pull request
- Non-goal: linting non-Markdown files
- Non-goal: high-context rewrite suggestions beyond deterministic rule output

## Changed Japanese Markdown Definition

A file counts as changed Japanese Markdown when all of these are true:

1. The file appears in `git diff --name-only --diff-filter=AMR <base-ref>...HEAD`.
2. The file still exists in the worktree after checkout.
3. The path ends with `.md`.
4. The file content contains at least one character from the helper's Japanese-writing ranges:
   - Hiragana and fullwidth Katakana (`U+3040`-`U+30FF`)
   - Katakana Phonetic Extensions (`U+31F0`-`U+31FF`)
   - CJK Unified Ideographs Extension A (`U+3400`-`U+4DBF`)
   - CJK Unified Ideographs (`U+4E00`-`U+9FFF`)
   - CJK Compatibility Ideographs (`U+F900`-`U+FAFF`)
   - Halfwidth Katakana (`U+FF65`-`U+FF9F`)

## Helper Script

`tools/list-changed-japanese-markdown.sh` is the source of truth for local and CI detection.

### Interface

```sh
tools/list-changed-japanese-markdown.sh [base-ref]
```

### Behavior

- defaults `base-ref` to `origin/main`
- compares `${base_ref}...HEAD`
- emits one matching path per line
- ignores deleted files
- limits eligibility to paths ending with `.md`
- includes changed Markdown anywhere in the repository when the content contains at least one character from the Japanese-writing ranges

## textlint Runtime

- `package.json` pins `textlint` and `textlint-rule-preset-ai-writing` as devDependencies
- `.textlintrc.json` enables `preset-ai-writing`
- the workflow and local verification commands load one repo-local custom rule with `--rulesdir ./tools/textlint-rules`
- the repo-local custom rule reads replacement terms from `config/textlint/terms.jsonl`

### Replacement Dictionary Contract

`config/textlint/terms.jsonl` uses one JSON object per line.

Each line must include:

- `pattern`: JavaScript regular expression source text
- `replacement`: preferred replacement text shown in findings

Example:

```json
{"pattern":"\\btaxonomy\\b","replacement":"分類"}
```

The workflow must fail for malformed JSONL or invalid regular expressions.

## GitHub Actions Workflow

`.github/workflows/japanese-textlint.yml` is the CI entrypoint.

### Workflow Behavior

- triggers on pull requests targeting `main`
- starts for changed Markdown and the workflow/helper/config/spec/runtime files that control this CI
- uses the repository-standard `actions/checkout` reference managed through `pinact`
- provisions a pinned Node runtime before installing dependencies
- checks out with full history (`fetch-depth: 0`) so `<base-ref>...HEAD` resolves cleanly
- explicitly fetches the PR base branch before diffing
- installs dependencies with `npm ci`
- calls `tools/list-changed-japanese-markdown.sh` with the fetched base ref
- exits successfully with a clear message when no changed Japanese Markdown files are found
- runs one `textlint --rulesdir ./tools/textlint-rules --format json` command over the full changed-file set
- converts `textlint` findings into GitHub Actions warning annotations
- writes a step summary with the file count, finding count, and a compact findings table
- upserts one PR comment with a stable marker instead of posting a new comment on every rerun or push
- includes the latest file count, finding count, findings table, and replacement hints in that PR comment
- paginates PR comment reads before deciding whether to create or update the marker comment
- uses explicit timeouts for both the `textlint` subprocess and GitHub API requests
- requests `contents: read`, `issues: write`, and `pull-requests: write` permissions for PR comment upsert
- keeps the job green when `textlint` reports prose findings, but fails the job if dependency install, config loading, dictionary parsing, or tool output parsing fails

## Reporting Contract

Each `textlint` finding should appear as a GitHub Actions warning annotation with:

- file path
- line
- column
- `ruleId` as the warning title
- lint message text
- replacement hint when the custom dictionary rule provides one

The step summary should include:

- number of checked files
- number of findings
- per-finding table when findings exist

The workflow should also keep one PR comment per pull request:

- comment body includes a stable hidden marker so reruns can find and update it
- comment body includes the latest file count, finding count, and findings table
- reruns after new pushes update the existing comment instead of creating another one

## Local Verification

Use the helper directly:

```sh
tools/list-changed-japanese-markdown.sh origin/main
```

For a focused local CI-style check, run:

```bash
npm ci
mapfile -t files < <(tools/list-changed-japanese-markdown.sh origin/main)
[ "${#files[@]}" -eq 0 ] || npx textlint --rulesdir ./tools/textlint-rules --format json "${files[@]}"
```
