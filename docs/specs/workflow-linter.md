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
tools/workflow-lint.sh --mode=ci [--pr-title=TITLE] [--pr-body=BODY] [--report-file=PATH]
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
- When `--report-file` is provided, also writes machine-readable JSON Lines to that path for CI consumption
- Each warning block includes:
  - warning class (`fixable` or `advisory`)
  - one primary finding message
  - rationale line prefixed with `WHY:`
  - remediation line prefixed with `FIX:` for `fixable` warnings
- Each JSON Lines warning record includes:
  - `type: "warning"`
  - `class`
  - `finding`
  - `why`
  - `fix`
- Counts warnings by finding block, not by output line
- Prints summary totals by warning class
- Prints a final reminder when any `fixable` warnings remain
- Appends one JSON Lines summary record when `--report-file` is provided:
  - `type: "summary"`
  - `total`
  - `fixable`
  - `advisory`
- Always exits 0

**Checks:**

| # | Check | Class | Mode | Description | Rule Source |
|---|-------|-------|------|-------------|-------------|
| 1 | Issue lifecycle | `fixable` | pre-push, ci | Files removed from `docs/issues/` must appear in `docs/issues/done/` (moved, not deleted) | AI_WORKFLOW.md Step 3: "Issue Resolution" |
| 2 | Docs-change hint | `advisory` | ci only | If code files changed but no `docs/` files changed, and PR title/body does not contain `[trivial]`, emit warning | AI_WORKFLOW.md: "Spec-Code Parity" principle |
| 3 | Branch naming | `fixable` | pre-push, ci | Branch name must match `<type>/<description>` where type is `plan\|feat\|fix\|chore\|docs` and description is non-empty kebab-case. `main` is exempt. | AI_WORKFLOW.md: "Branch Naming Convention" |
| 4 | Exec-plan existence | `fixable` | pre-push, ci | For `feat/*` and `fix/*` branches, a matching exec-plan must exist whose basename suffix is `-<name>.md`, where `<name>` is the branch description. Historical completed plans may still use `docs/exec-plan/done/<name>.md`. `plan/*`, `chore/*`, `docs/*` branches are exempt. | AI_WORKFLOW.md: "Exec-Plan Mapping" |
| 4a | Ambiguous exec-plan mapping | `fixable` | pre-push, ci | For `feat/*` and `fix/*` branches, if multiple active or completed exec-plans share the same `-<name>.md` suffix in one directory, emit a warning because suffix-based mapping is ambiguous and must be cleaned up. | AI_WORKFLOW.md: "Exec-Plan Mapping" |
| 5 | Workflow startup wording | `fixable` | pre-push, ci | If changed migrated workflow-facing docs or skills reintroduce raw startup snippets like `git fetch origin` or `git switch -c`, emit a warning to keep global `ww` as the default operator path. Covered skills include `plan-execution`, `execute-task`, `triage-tasks`, `plan-project`, `review-task`, and `manage-workflow`. | docs/specs/ww-dogfooding-workflow.md: "Workflow lint guard" |
| 6 | Linked local issue resolution | `fixable` | pre-push, ci | For `feat/*` and `fix/*` branches whose matching exec-plan is completed in `docs/exec-plan/done/` and whose `Addresses:` line names local issue paths under `docs/issues/`, each linked issue must also be moved to `docs/issues/done/` in the same branch unless the PR body explicitly justifies leaving it open (for example `remains open: <reason>`). | AI_WORKFLOW.md Step 3: "Issue Resolution" |
| 6a | Linked external GitHub issue closure metadata | `fixable` | ci only | For `feat/*` and `fix/*` branches whose matching exec-plan is completed in `docs/exec-plan/done/` and whose `Addresses:` line names external GitHub issue URLs, the PR body must include matching closing keywords (`Closes #123` for same-repo issues or `Closes <full-url>` for cross-repo issues) unless the PR body explicitly justifies leaving the issue open. | AI_WORKFLOW.md Step 3: "External GitHub Issue Resolution" |
| 7 | Active plan / issue naming | `fixable` | pre-push, ci | Active files under `docs/exec-plan/todo/` and `docs/issues/` must use `<sequence>-<name>.md`, with four-digit zero padding through `9999` and unpadded numbers at `10000+`. `README.md` is exempt. | AI_WORKFLOW.md: "Active Plan / Issue Naming" |

The missing-base-ref and diff-failure advisories are environment-sensitive guardrail behavior, not repo-policy checks. They exist to keep shallow, partially fetched, or otherwise unusual repository states from producing a misleading "no changes" result while still preserving non-diff checks.

**Operating rule:**
- `fixable` warnings should normally be resolved before push/PR
- `fixable` warnings may be skipped only for explicit human instruction or a clear false positive
- skipped `fixable` warnings must be justified in the PR body
- `advisory` warnings remain non-blocking judgment calls
- exit behavior stays non-blocking (`exit 0`)

**Exec-plan filename convention:**
- Active exec-plan filenames use `<sequence>-<name>.md`.
- For execution branches, the exec-plan basename suffix must match the branch description. For example, `feat/workflow-linter` maps to `docs/exec-plan/todo/0042-workflow-linter.md` during execution and `docs/exec-plan/done/0042-workflow-linter.md` after completion.
- Historical completed artifacts may still use older non-numbered filenames in `docs/exec-plan/done/`, and workflow-lint must keep tolerating them.
- Active local issue filenames use the same `<sequence>-<name>.md` rule under `docs/issues/`.

**`Addresses:` convention for local issues:**
- Execution plans that expect to resolve tracked local issues should include a single `Addresses:` line.
- That line lists one or more local issue paths under `docs/issues/`, for example `Addresses: docs/issues/0019-bug-name.md`.
- `workflow-lint` treats those paths as explicit completion metadata only for the matching `feat/<name>` or `fix/<name>` execution branch after the plan has moved to `docs/exec-plan/done/<sequence>-<name>.md` or a tolerated historical `docs/exec-plan/done/<name>.md`.
- The parser must accept both same-line `Addresses: docs/issues/...` entries and the common multi-line form where `Addresses:` is followed by bullet lines.
- The linked-local-issue check stays narrow on purpose: it does not try to infer issue closure from arbitrary code changes or unrelated branches.

**`Addresses:` convention for external GitHub issues:**
- When the canonical tracker is an external GitHub issue in the target repository, list the issue on the same `Addresses:` line using the full issue URL, for example `Addresses: https://github.com/yoskeoka/ww/issues/227`.
- Use full URLs in plans so the closure target stays unambiguous even when the plan is reviewed from the workspace root or copied across repos.
- `workflow-lint` only checks external issue closure metadata in CI mode because the PR body is required input.
- The CI check is narrow on purpose: it only verifies that explicitly linked issues from the completed plan appear in PR-body closing keywords or are intentionally left open with justification.

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
- Copies `.github/workflows/workflow-lint.yml` to the target repo's `.github/workflows/` directory
- Sets `git config core.hooksPath .githooks` in the target repo
- Idempotent (safe to run multiple times)

### `setup-workspace.sh` (renamed from `setup-skills.sh`)

Adds a workflow-asset installation step after skill symlink setup:
- Calls `tools/install-hooks.sh` for the child repo
- All existing functionality preserved

### `.github/workflows/workflow-lint.yml`

GitHub Actions workflow that runs the linter on PRs targeting `main`.

- Triggers on `pull_request` to `main`
- Uses the repository-standard `actions/checkout` reference managed through `pinact`
- Checks out with full history (`fetch-depth: 0`) so the resolved base ref exists locally and diff checks can run
- Passes PR title and body from GitHub event context to `--pr-title` / `--pr-body`
- Passes `--report-file` to capture JSON Lines warning output from the linter
- Emits GitHub Actions warning annotations for each finding so reviewers can see `fixable` and `advisory` warnings without opening raw logs
- Writes a step summary that includes warning counts and a short finding list
- Keeps the job green when the linter only finds warnings; visibility changes, blocking behavior does not

## Non-Goals

- Spec-sync checking (context-dependent, left to human review)
- Trivial change detection (human declaration via `[trivial]`)
- Project-specific lint (managed per-project)
- ~~Branch naming enforcement~~ (implemented in Check 3)
- ~~Exec-plan existence/completion enforcement~~ (implemented in Check 4)
