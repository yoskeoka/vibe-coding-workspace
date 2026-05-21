# Make workflow-lint CI warnings visible and roll the change into child repos

**Execution**: Use `/execute-task` to implement this plan.

Addresses: docs/issues/0022-workflow-lint-ci-fixable-warnings-are-non-blocking.md

## Objective

Keep `tools/workflow-lint.sh` non-blocking while making `fixable` warnings visible in GitHub Actions review surfaces, then distribute the updated workflow-lint runtime assets to the workspace repo and all managed child repos through reviewable PRs.

This plan should produce:

- one workspace implementation PR that updates the warning contract, CI presentation, and workflow-asset distribution rules
- follow-up child-repo PRs for every managed repo that consumes the copied workflow-lint assets

## Background

The current workflow-linter contract keeps all findings non-blocking, including CI mode, to avoid false-positive frustration. That part is still sound.

The problem is visibility:

- `.github/workflows/workflow-lint.yml` only runs `./tools/workflow-lint.sh --mode=ci`
- `tools/workflow-lint.sh` exits `0` even when `fixable` warnings exist
- reviewers must open raw job logs to notice those warnings

We considered making `--mode=ci` exit non-zero and absorbing it in GitHub Actions with `continue-on-error`, but that mixes two concerns:

- linter verdict semantics
- CI presentation semantics

The chosen direction is to keep the script non-blocking and move reviewer-facing warning presentation into the workflow layer.

There is also a rollout constraint: child repos do not consume workflow-lint entirely through the vendored `skills/` subtree. They keep repo-local copies of at least:

- `tools/workflow-lint.sh`
- `.githooks/pre-push`
- `.github/workflows/workflow-lint.yml`

So the current `skills/`-only workflow-sync contract is too narrow for this rollout.

## Design Decisions

### 1. Keep `tools/workflow-lint.sh` as a warning emitter, not a hard gate

- `--mode=pre-push` and `--mode=ci` stay `exit 0`
- `fixable` versus `advisory` classification remains the linter's contract
- the linter should grow machine-readable output or another structured surface that CI can transform into annotations, summaries, or PR comments

### 2. Move warning visibility into `.github/workflows/workflow-lint.yml`

- GitHub Actions should surface `fixable` warnings without turning the check red by default
- acceptable reviewer surfaces include job annotations, step summary, and a PR comment if needed
- the workflow should distinguish `fixable` versus `advisory` findings clearly enough that reviewers can act without opening raw logs

### 3. Treat workflow-lint runtime assets as copied distribution artifacts

For this rollout, the copied child-repo asset set must include at least:

- `tools/workflow-lint.sh`
- `.githooks/pre-push`
- `.github/workflows/workflow-lint.yml`

If implementation shows another repo-local helper is required to make the rollout reproducible, add it to the same documented asset set rather than relying on unstated manual steps.

### 4. Expand workflow-sync beyond `skills/` when runtime assets change

The current child-repo sync rule only opens PRs for `skills/` changes because that was the only consumed submodule surface. This plan changes that assumption.

The updated sync contract should open child PRs when the pushed workflow commit changes any documented child-consumed workflow asset set, not just `skills/`.

## Spec Changes

### `docs/specs/workflow-linter.md`

- keep the non-blocking exit-code contract explicit
- document the new CI presentation contract for `fixable` and `advisory` warnings
- define the machine-readable or otherwise structured output that the CI workflow will consume
- state which warning surfaces are expected in GitHub Actions

### `docs/specs/workflow-sync-to-child-repos.md`

- replace the current `skills/`-only sync rule with a documented child-consumed asset set
- define how the sync workflow detects relevant changes for copied runtime assets as well as `skills/`
- document that child sync PRs may include copied workflow runtime files, not just the submodule bump

### `docs/specs/github-actions-pinning.md`

- update only if the chosen CI implementation adds or changes `uses:` references

### Child-repo workflow-maintenance docs

- update any repo-local workflow-linter maintenance doc that still claims the rollout surface is only `tools/workflow-lint.sh` plus `.githooks/pre-push`
- at minimum, align the statement that `.github/workflows/workflow-lint.yml` is also a copied workflow asset

## Code Changes

### Workspace repo

- update `tools/workflow-lint.sh` to emit structured warning data that CI can consume without scraping human-only prose
- update `.github/workflows/workflow-lint.yml` to transform workflow-lint findings into reviewer-visible warning surfaces while keeping the check non-blocking by default
- add or update any helper script needed to sync copied workflow assets into child repos
- update `setup-workspace.sh` and/or `tools/install-hooks.sh` if the documented child bootstrap path must also copy the workflow-lint CI file or other runtime assets
- update `.github/workflows/sync-workflow-to-child-repos.yml` so generated child PRs include the copied runtime-asset changes, not just a submodule pointer move

### Child repos

For each managed child repo that consumes copied workflow-lint assets:

- update `tools/workflow-lint.sh`
- update `.githooks/pre-push` if needed
- update `.github/workflows/workflow-lint.yml`
- update any repo-local workflow maintenance doc that names the copied asset set

Managed repos in scope:

- `ai-arena`
- `dungeon-game-ai-arena`
- `reversi-adventure`
- `reversi-ai-arena`
- `vim-learning-game`
- `ww`
- `envdiff`

## PR / Rollout Plan

1. Create the workspace implementation PR with the spec, CI, and sync-contract changes.
2. Merge that PR to `main`.
3. Let the updated workspace sync flow create child-repo rollout PRs, or if one repo needs a documented exception, create that child PR manually and record why automation was insufficient.
4. Verify that every in-scope child repo has an open PR carrying the updated copied workflow assets.

## Sub-tasks

- [ ] Update the workflow-linter spec to describe non-blocking CI warning presentation
- [ ] Update the workflow-sync spec to describe the broader child-consumed asset set
- [ ] Decide the structured output shape from `tools/workflow-lint.sh` to CI
- [ ] [parallel] Design the GitHub Actions warning surfaces to use for `fixable` and `advisory` findings
- [ ] [parallel] Audit bootstrap and sync paths for copied workflow runtime assets
- [ ] [depends on: spec decisions] Implement workspace-side workflow-lint structured output and CI presentation
- [ ] [depends on: sync audit] Implement child-asset sync changes in `setup-workspace.sh`, related helper scripts, and `.github/workflows/sync-workflow-to-child-repos.yml`
- [ ] [depends on: workspace implementation] Open and verify the workspace PR
- [ ] [depends on: workspace PR merged] Open or confirm child rollout PRs for all managed repos
- [ ] [depends on: child rollout PRs] Review each child PR for copied-file parity and obvious repo-specific breakage

## Parallelism

- `[parallel]` CI warning-surface design and copied-asset-path audit can proceed independently after the spec direction is fixed.
- Child-repo PR verification can run in parallel across repositories after the rollout PRs exist.

## Verification

### Workspace verification

- run `./tools/workflow-lint.sh --mode=pre-push` and confirm exit code remains `0`
- run `./tools/workflow-lint.sh --mode=ci --pr-title="..." --pr-body="..."` against a state that emits at least one `fixable` warning and confirm the structured output is consumable by CI
- verify the workflow-lint GitHub Actions job shows reviewer-visible warning output without a hard failure when only warnings are present
- verify any newly added or modified workflow `uses:` references are pinned through `pinact`

### Rollout verification

- confirm the updated sync workflow detects runtime-asset changes, not just `skills/` changes
- confirm each in-scope child repo receives a PR that includes the copied workflow-lint assets
- confirm each child PR keeps the repo-local files aligned:
  - `tools/workflow-lint.sh`
  - `.githooks/pre-push`
  - `.github/workflows/workflow-lint.yml`
- confirm any repo-local maintenance doc changes match the copied-asset contract

## Expected Outcome

- `fixable` workflow-lint findings become visible in GitHub Actions review surfaces without converting the job into a hard blocker
- the workflow contract stays robust against false positives because the linter remains non-blocking
- workspace-owned automation can distribute workflow-lint runtime asset changes to child repos
- every managed child repo ends the rollout with an open PR carrying the updated copied workflow assets
