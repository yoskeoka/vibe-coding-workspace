# Spec: Workflow Linter

## Goal

Mechanically enforce AI-Centered Development workflow rules declared in `AI_WORKFLOW.md`. Runs as a pre-push git hook (local) and as a CI check. All checks are **warnings only** (exit 0) to provide visibility without blocking.

## Warning Classes

Warnings are classified to distinguish findings that should normally be fixed before push/PR from findings that are informative only.

| Class | Meaning | Expected Action |
|---|---|---|
| `fixable` | The repo state can usually be corrected immediately by moving files, renaming a branch, creating a missing plan, or making another straightforward change | Resolve before push/PR unless an explicit human instruction conflicts or the warning is a clear false positive |
| `advisory` | Useful workflow signal, but not something that should automatically trigger a repo mutation | Review and use judgment; no default repo mutation required |

When a `fixable` warning is intentionally skipped, the PR body must explain why. Allowed reasons are:
- an explicit user/human instruction takes precedence
- the warning is judged to be a clear false positive

## Components

### `tools/workflow-lint.sh`

Main linter script.

**Interface:**
```
tools/workflow-lint.sh --mode=pre-push
tools/workflow-lint.sh --mode=ci [--pr-title=TITLE] [--pr-body=BODY]
```

**Behavior:**
- Resolves the diff base as `origin/${GITHUB_BASE_REF}` when `GITHUB_BASE_REF` is set, otherwise `origin/main`
- Verifies the resolved base ref exists locally before running diff-based checks
- Computes changed files via `git diff --name-only --diff-filter=ADMR <base-ref>...HEAD`
- If the base ref is missing locally, emits an advisory warning and skips only diff-based checks instead of silently treating that state as "no changes"
- If diff computation fails for another repository-state reason, emits an advisory warning and skips only diff-based checks
- Runs branch / exec-plan checks even when diff-based checks are skipped
- Runs checks based on mode
- Outputs normalized warning blocks to stderr
- Each warning block includes:
  - warning class (`fixable` or `advisory`)
  - one primary finding message
  - rationale line prefixed with `WHY:`
  - remediation line prefixed with `FIX:` for `fixable` warnings
- Counts warnings by finding block, not by output line
- Prints summary totals by warning class
- Prints a final reminder when any `fixable` warnings remain
- Always exits 0

**Checks:**

| # | Check | Class | Mode | Description | Rule Source |
|---|-------|-------|------|-------------|-------------|
| 1 | Issue lifecycle | `fixable` | pre-push, ci | Files removed from `docs/issues/` must appear in `docs/issues/done/` (moved, not deleted) | AI_WORKFLOW.md Step 3: "Issue Resolution" |
| 2 | Docs-change hint | `advisory` | ci only | If code files changed but no `docs/` files changed, and PR title/body does not contain `[trivial]`, emit warning | AI_WORKFLOW.md: "Spec-Code Parity" principle |
| 3 | Branch naming | `fixable` | pre-push, ci | Branch name must match `<type>/<description>` where type is `plan\|feat\|fix\|chore\|docs` and description is non-empty kebab-case. `main` is exempt. | AI_WORKFLOW.md: "Branch Naming Convention" |
| 4 | Exec-plan existence | `fixable` | pre-push, ci | For `feat/*` and `fix/*` branches, `docs/exec-plan/todo/<name>.md` or `docs/exec-plan/done/<name>.md` must exist, where `<name>` is the branch description. `plan/*`, `chore/*`, `docs/*` branches are exempt. | AI_WORKFLOW.md: "Exec-Plan Mapping" |
| 5 | Workflow startup wording | `fixable` | pre-push, ci | If changed migrated workflow-facing docs or skills reintroduce raw startup snippets like `git fetch origin` or `git switch -c`, emit a warning to keep global `ww` as the default operator path. Covered skills include `plan-execution`, `execute-task`, `triage-tasks`, `plan-project`, `review-task`, and `manage-workflow`. | docs/specs/ww-dogfooding-workflow.md: "Workflow lint guard" |

The missing-base-ref and diff-failure advisories are environment-sensitive guardrail behavior, not repo-policy checks. They exist to keep shallow, partially fetched, or otherwise unusual repository states from producing a misleading "no changes" result while still preserving non-diff checks.

**Operating rule:**
- `fixable` warnings should normally be resolved before push/PR
- `fixable` warnings may be skipped only for explicit human instruction or a clear false positive
- skipped `fixable` warnings must be justified in the PR body
- `advisory` warnings remain non-blocking judgment calls
- exit behavior stays non-blocking (`exit 0`)

**Exec-plan filename convention:**
- Active exec-plan filenames use descriptive kebab-case without numeric prefixes.
- For execution branches, the filename stem must match the branch description exactly. For example, `feat/workflow-linter` maps to `docs/exec-plan/todo/workflow-linter.md` during execution and `docs/exec-plan/done/workflow-linter.md` after completion.
- Numeric examples may appear in historical completed artifacts when documenting prior workflow state, but live workflow docs, templates, and active todo plans must use the non-numeric convention.

**Exit codes:**
- Always 0 (warnings only)

### `.githooks/pre-push`

Thin wrapper that invokes `tools/workflow-lint.sh --mode=pre-push`.

- Resolves the repo root via `git rev-parse --show-toplevel`
- Calls the linter script
- Passes through exit code (always 0)

### `tools/install-hooks.sh`

Sets up git hooks for a repository.

**Interface:**
```
tools/install-hooks.sh [repo-path]
```

**Behavior:**
- Defaults to current directory if no path given
- Copies `.githooks/pre-push` to the target repo's `.githooks/` directory
- Copies `tools/workflow-lint.sh` to the target repo's `tools/` directory
- Sets `git config core.hooksPath .githooks` in the target repo
- Idempotent (safe to run multiple times)

### `setup-workspace.sh` (renamed from `setup-skills.sh`)

Adds a hook installation step after skill symlink setup:
- Calls `tools/install-hooks.sh` for the child repo
- All existing functionality preserved

### `.github/workflows/workflow-lint.yml`

GitHub Actions workflow that runs the linter on PRs targeting `main`.

- Triggers on `pull_request` to `main`
- Uses the repository-standard `actions/checkout@v6`
- Checks out with full history (`fetch-depth: 0`) so the resolved base ref exists locally and diff checks can run
- Passes PR title and body from GitHub event context to `--pr-title` / `--pr-body`

## Non-Goals

- Spec-sync checking (context-dependent, left to human review)
- Trivial change detection (human declaration via `[trivial]`)
- Project-specific lint (managed per-project)
- ~~Branch naming enforcement~~ (implemented in Check 3)
- ~~Exec-plan existence/completion enforcement~~ (implemented in Check 4)
