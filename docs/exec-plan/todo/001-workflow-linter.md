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
- **Branch naming convention**: Current workflow docs use examples (`plan/002-feature-x`, `feat/002-feature-x`) but do not declare a mandatory naming rule. Linter enforcement requires a declared rule first. Out of scope until the rule is formalized.
- **Exec plan existence / completion**: Checking whether a branch has a corresponding exec-plan file requires a branch-name-to-plan-file mapping convention (e.g., `feat/002-*` → `002-*.md`). This convention is implicit, not declared. Out of scope until the naming rule is formalized.

### Pre-commit hooks

**Decision: Not used.** Pre-commit operates at commit granularity, but workflow rules operate at branch/PR granularity. Valid multi-commit patterns (spec commit → code commit, fix commit → issue-done-move commit) would be blocked. Project-specific pre-commit lint (formatters, etc.) is a separate design topic.

### Hook placement

- **Local**: pre-push only (sees full branch diff via `origin/main..HEAD`)
- **CI**: Same checks as pre-push + PR metadata checks (has access to PR title/body)

## Checks

Only rules that are **explicitly declared** in `AI_WORKFLOW.md` and **mechanically verifiable** are included.

### Pre-push checks (low false-positive, warning-only)

| # | Check | Description | Rule source | Severity |
|---|-------|-------------|-------------|----------|
| 1 | Issue lifecycle | Files removed from `docs/issues/` must appear in `docs/issues/done/` (moved, not deleted) | AI_WORKFLOW.md Step 3: "Issue Resolution" | warning |

### CI-only checks (may have false positives, uses PR metadata)

| # | Check | Description | Rule source | Severity |
|---|-------|-------------|-------------|----------|
| 2 | Docs-change hint | If code files changed but no `docs/` files changed, and PR title/body does not contain `[trivial]`, emit warning | AI_WORKFLOW.md: "Spec-Code Parity" principle | warning |

All checks are **warnings** (non-blocking). The goal is visibility, not gatekeeping.

### Future checks (blocked on rule formalization)

These checks are desirable but require workflow rule updates before implementation:

- **Branch naming convention**: Needs a declared rule in `AI_WORKFLOW.md` (e.g., "Branch names MUST match `<type>/<NNN>-<description>` where type is plan|feat|fix|chore|docs").
- **Exec plan existence**: Needs the branch naming rule above to establish the branch→plan mapping.
- **Exec plan completion (todo→done)**: Same dependency on branch→plan mapping.

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
