# Spec: Workflow Linter

## Goal

Mechanically enforce AI-Centered Development workflow rules declared in `AI_WORKFLOW.md`. Runs as a pre-push git hook (local) and as a CI check. All checks are **warnings only** (exit 0) to provide visibility without blocking.

## Components

### `tools/workflow-lint.sh`

Main linter script.

**Interface:**
```
tools/workflow-lint.sh --mode=pre-push
tools/workflow-lint.sh --mode=ci [--pr-title=TITLE] [--pr-body=BODY]
```

**Behavior:**
- Computes changed files via `git diff --name-only --diff-filter=ADMR origin/main...HEAD`
- Runs checks based on mode
- Outputs colored warnings to stderr
- Always exits 0

**Checks:**

| # | Check | Mode | Description | Rule Source |
|---|-------|------|-------------|-------------|
| 1 | Issue lifecycle | pre-push, ci | Files removed from `docs/issues/` must appear in `docs/issues/done/` (moved, not deleted) | AI_WORKFLOW.md Step 3: "Issue Resolution" |
| 2 | Docs-change hint | ci only | If code files changed but no `docs/` files changed, and PR title/body does not contain `[trivial]`, emit warning | AI_WORKFLOW.md: "Spec-Code Parity" principle |

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
- Checks out with full history (`fetch-depth: 0`) so `origin/main...HEAD` diff works
- Passes PR title and body from GitHub event context to `--pr-title` / `--pr-body`

## Non-Goals

- Spec-sync checking (context-dependent, left to human review)
- Trivial change detection (human declaration via `[trivial]`)
- Project-specific lint (managed per-project)
- Branch naming enforcement (future plan)
- Exec-plan existence/completion enforcement (future plan)
