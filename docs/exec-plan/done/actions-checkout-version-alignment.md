# Actions Checkout Version Alignment
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Resolve `docs/issues/actions-checkout-version-alignment-follow-up.md` by choosing one `actions/checkout` major version for this repository and aligning every GitHub Actions workflow to that standard.

This supports the project-plan tooling goal that workflow infrastructure remains reviewable, mechanically verifiable, and low-maintenance across workspace work.

## Context

Current workflow inventory:

- `.github/workflows/kb-pages.yml` uses `actions/checkout@v4`
- `.github/workflows/workflow-lint.yml` uses `actions/checkout@v6`
- `.github/workflows/shellcheck.yml` uses `actions/checkout@v6`

Relevant prior guidance:

- `docs/project-plan.md` treats the workflow itself as a first-class deliverable and calls for reusable, mechanically enforced workflow tooling.
- `docs/design-decisions/core-beliefs.md` favors spec/code parity and avoiding unnecessary refactors.
- `docs/design-decisions/adr/0003-dogfood-released-global-ww.md` records a dogfooding decision to keep workflow docs, skills, and tooling aligned rather than letting process definitions drift.

External compatibility note:

- `actions/checkout` v6 is the latest major line observed during planning, with v6 release notes calling out Node.js 24 support details and runner compatibility requirements.

## Design Options

### Option A: Standardize on `actions/checkout@v6` (recommended)

Update `.github/workflows/kb-pages.yml` from `v4` to `v6`, leaving the two existing `v6` workflows unchanged.

Why this fits:

- Smallest repository diff.
- Matches the current majority in this repo.
- Aligns with the latest major version currently published for `actions/checkout`.
- Avoids making workflow lint and ShellCheck CI move backward without a demonstrated compatibility need.

Risk:

- `v6` has runner compatibility requirements. Execution must confirm GitHub-hosted `ubuntu-latest` satisfies those requirements or document a reason to choose another standard.

### Option B: Standardize on `actions/checkout@v4`

Downgrade `.github/workflows/workflow-lint.yml` and `.github/workflows/shellcheck.yml` to `v4`.

Why this could fit:

- Conservative if runner compatibility becomes uncertain.

Why this is not recommended:

- Larger diff.
- Moves two already-working workflows away from the current major line.
- Does not match the repository's existing direction.

### Option C: Keep the mixed versions and document why

Leave YAML unchanged and add documentation explaining why `kb-pages` remains on `v4`.

Why this could fit:

- Appropriate only if `kb-pages` has a real compatibility constraint.

Why this is not recommended:

- The issue exists because the mixed state currently has no documented reason.
- No workflow-specific checkout behavior has been identified that requires `v4`.

## Spec Changes

- [x] Update `docs/specs/knowledge-base.md` to state that the `kb-pages` workflow uses the repository-standard checkout major.
- [x] Update `docs/specs/workflow-linter.md` to state that `workflow-lint.yml` uses the repository-standard checkout major with `fetch-depth: 0`.
- [x] Update `docs/specs/shellcheck-ci.md` to state that `shellcheck.yml` uses the repository-standard checkout major with `fetch-depth: 0`.
- [x] Add a small shared note in the relevant spec text that the current repository standard is `actions/checkout@v6`, unless execution discovers a compatibility reason to choose a different major.

## Code Changes

- [x] [parallel] Update `.github/workflows/kb-pages.yml` from `actions/checkout@v4` to the chosen standard, expected to be `actions/checkout@v6`.
- [x] [parallel] Confirm `.github/workflows/workflow-lint.yml` already matches the chosen standard.
- [x] [parallel] Confirm `.github/workflows/shellcheck.yml` already matches the chosen standard.
- [x] Move `docs/issues/actions-checkout-version-alignment-follow-up.md` to `docs/issues/done/actions-checkout-version-alignment-follow-up.md` after the workflow and spec updates are complete.

## Sub-tasks

- [x] Confirm the final `actions/checkout` major against official release notes and runner requirements.
- [x] Update specs first.
- [x] Update workflow YAML.
- [x] Verify no remaining mixed `actions/checkout` major versions with `rg -n "actions/checkout@" .github/workflows docs/specs docs/issues`.
- [x] Run workflow quality gates:
  - `./tools/workflow-lint.sh --mode=pre-push`
  - `shellcheck tools/workflow-lint.sh tools/list-changed-bash-scripts.sh tools/install-hooks.sh setup-workspace.sh` if `shellcheck` is locally available
  - `tools/kb check` if the knowledge-base spec or workflow changes need KB build validation
- [x] Move this plan from `docs/exec-plan/todo/` to `docs/exec-plan/done/` during execution.

## Parallelism

- The three workflow inventory checks are independent and can be done in parallel.
- Spec updates should happen before YAML edits to preserve spec-first workflow.
- Issue closure and plan move depend on successful spec/code verification.

## Design Decisions

No ADR update is expected. This is a maintenance alignment of existing GitHub Actions usage, not a new architectural direction.

If execution finds a real compatibility constraint that prevents `v6`, record the chosen exception in the specs and consider an ADR entry only if the exception becomes a durable repository policy.
