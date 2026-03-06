# 001: Workflow Linter

## Objective

Implement a bash-based workflow linter that mechanically enforces AI-Centered Development workflow rules. The linter runs as a pre-push git hook locally and as a CI check with additional PR-metadata-aware checks.

## Background

AI agents repeatedly violated workflow rules during development (e.g., misusing GitHub Issues, skipping spec updates, pushing to merged PRs). Human review alone is insufficient — mechanical enforcement is needed for rules that can be checked programmatically.

## Design Decisions

### What the linter does NOT check

- **Spec-sync**: Whether code changes require spec updates is context-dependent (bug fixes, refactors may not need spec changes). This is fundamentally undecidable by a shell script. Left to human review.
- **Trivial detection**: Whether a change is "trivial" is a human declaration (via PR title/body `[trivial]`), not mechanically detectable.
- **Project-specific lint**: `cargo clippy`, `eslint`, etc. are separate concerns managed per-project. Not in scope.

### Pre-commit hooks

**Decision: Not used.** Pre-commit operates at commit granularity, but workflow rules operate at branch/PR granularity. Valid multi-commit patterns (spec commit → code commit, fix commit → issue-done-move commit) would be blocked. Project-specific pre-commit lint (formatters, etc.) is a separate design topic.

### Hook placement

- **Local**: pre-push only (sees full branch diff via `origin/main..HEAD`)
- **CI**: Same checks as pre-push + PR metadata checks (has access to PR title/body)

## Checks

### Pre-push checks (low false-positive, warning-only)

| # | Check | Description | Severity |
|---|-------|-------------|----------|
| 1 | Branch naming | Branch must match `(plan\|feat\|fix\|chore\|docs)/<NNN>-*` or `(plan\|feat\|fix\|chore\|docs)/*` pattern | warning |
| 2 | Exec plan exists | If branch name contains `<NNN>`, a matching file must exist in `docs/exec-plan/todo/` or `docs/exec-plan/done/` | warning |
| 3 | Issue lifecycle | Files removed from `docs/issues/` must appear in `docs/issues/done/` (moved, not deleted) | warning |

### CI-only checks (may have false positives, uses PR metadata)

| # | Check | Description | Severity |
|---|-------|-------------|----------|
| 4 | Docs-change hint | If code files changed but no `docs/` files changed, and PR title/body does not contain `[trivial]`, emit warning | warning |
| 5 | Exec plan completion | If branch is `feat/` or `fix/`, check that the corresponding plan file is moved to `done/` (only if the plan file existed in `todo/`) | warning |

All checks are **warnings** (non-blocking). The goal is visibility, not gatekeeping.

## File Structure

```
tools/
  workflow-lint.sh       # Main linter script
  install-hooks.sh       # Installs pre-push hook in a repo
.githooks/
  pre-push               # Hook script that calls workflow-lint.sh
```

Distribution to child repos: via `setup-skills.sh` (rename to `setup-workspace.sh`) which also handles skill symlinks.

## Code Changes

### `tools/workflow-lint.sh`
- Bash script, `set -euo pipefail`
- Accepts mode flag: `--mode=pre-push` or `--mode=ci`
- CI mode accepts additional flags: `--pr-title=...` `--pr-body=...`
- Reads branch name from `git branch --show-current`
- Reads diff from `git diff --name-only --diff-filter=ADMR origin/main...HEAD`
- Exit code 0 always (warnings only)
- Output: colored warnings to stderr

### `.githooks/pre-push`
- Thin wrapper that calls `tools/workflow-lint.sh --mode=pre-push`

### `tools/install-hooks.sh`
- Sets `git config core.hooksPath` to point to `.githooks/`
- Called by `setup-workspace.sh`

### `setup-skills.sh` → `setup-workspace.sh` rename
- Add hook installation step
- Add tools/ symlink or copy to child repos

## Spec Changes

- Create `docs/specs/workflow-linter.md` documenting the linter's checks, configuration, and usage.

## Sub-tasks

- [ ] [parallel] Create `docs/specs/workflow-linter.md`
- [ ] [parallel] Implement `tools/workflow-lint.sh`
- [ ] [parallel] Create `.githooks/pre-push`
- [ ] [depends on: workflow-lint.sh] Create `tools/install-hooks.sh`
- [ ] [depends on: all above] Rename `setup-skills.sh` → `setup-workspace.sh` and integrate
- [ ] [depends on: all above] Test in reversi-adventure (first child repo)
- [ ] [depends on: test] Update child repo setup documentation
