# `tools/pj` CI Checks

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Add a dedicated GitHub Actions workflow that continuously verifies the workspace-local `tools/pj` Go module. The goal is to keep the GitHub Projects-backed task CLI mechanically checked in PRs instead of relying on ad hoc local test runs.

This supports the project-plan requirements for low-friction workspace task coordination and a reusable task CLI spike by making formatting, linting, vetting, and tests visible in normal PR review.

## Background

`tools/pj` now has Go tests under `tools/pj/internal/pj/`, but the repository's current CI only covers workflow linting and changed Bash scripts. That leaves regressions in `tools/pj` invisible unless an operator remembers to run Go checks locally.

Past decisions and constraints:

- `tools/pj` intentionally remains a workspace-local Go module under `tools/pj/`.
- GitHub Projects is the remote source of truth for workspace triage, and `.local/pj/` remains derived local state.
- Existing CI workflows use small focused jobs rather than a monolithic all-repo check.

## Trade-offs

### Option A: Add a dedicated `check-pj` workflow

- Keeps `tools/pj` checks obvious in the PR checks list.
- Lets the workflow run with `working-directory: tools/pj`, matching the independent module layout.
- Avoids coupling future non-Go workspace checks to the `pj` module.
- Recommended for this plan because the current repo already separates workflow-lint and ShellCheck by responsibility.

### Option B: Fold `tools/pj` checks into the existing workflow-lint workflow

- Reduces the number of workflow files.
- Makes the workflow-lint job name less accurate and mixes unrelated failure modes.
- Less aligned with the existing focused-workflow pattern.

### Option C: Add a root `make test` / `make lint` first, then call it from CI

- Could create a nice local abstraction for humans and agents.
- Adds a second layer before the immediate CI gap is closed.
- Better as a follow-up if more workspace tools need shared commands.

## Spec Changes

### `docs/specs/github-projects-task-cli.md`

- Add a CI/verification section for the `tools/pj` module.
- Define that PRs touching `tools/pj`, its specs, or its workflow should run the `check-pj` GitHub Actions workflow.
- Specify the expected checks:
  - `goimports` formatting check
  - Go lint check
  - `go vet ./...`
  - `go test ./...`
- Document that the workflow runs from the `tools/pj` module directory.

## Code Changes

### `.github/workflows/check-pj.yml`

- Create a new GitHub Actions workflow named `check-pj`.
- Trigger it on pull requests to `main` with a `paths` filter limited to:
  - `tools/pj/**`
  - `docs/specs/github-projects-task-cli.md`
  - `.github/workflows/check-pj.yml`
- Use `actions/checkout@v6`, matching the existing workflow version alignment.
- Set up the Go version required by `tools/pj/go.mod`.
- Install the lint/format tooling used by the workflow:
  - `goimports` for formatting checks
  - a Go lint tool suitable for CI
- Run all checks with `working-directory: tools/pj`.
- Keep commands deterministic and reviewable; avoid mutating files in CI.

### Tooling choice during execution

- Pin formatter and linter install versions in the workflow (for example, `goimports@vX.Y.Z` and linter `@vA.B.C`) so CI runs are deterministic.
- Prefer a commonly maintained Go lint path during implementation. If `golint` is unsuitable or stale, choose a better-supported linter and record the exact pinned command in the spec and PR verification.
- Document the chosen pinned versions and a simple upgrade process in `docs/specs/github-projects-task-cli.md`.

## Design Decisions

- No ADR update is expected. This is CI coverage for an existing module, not a new architecture direction.
- Keep the first version scoped to `tools/pj`; do not create a repository-wide Go module or a shared `Makefile` unless implementation shows it is necessary.
- Do not change `tools/pj` behavior as part of this plan. Formatting or lint fixes needed to satisfy CI are in scope only when they are mechanical and directly required by the new checks.

## Sub-tasks

- [ ] Update `docs/specs/github-projects-task-cli.md` with the `tools/pj` CI contract.
- [ ] [parallel] Review `tools/pj/go.mod` and existing Go files to choose compatible Go setup and lint/format commands.
- [ ] [parallel] Review existing GitHub Actions workflows for trigger, checkout, and naming conventions.
- [ ] [depends on: spec update, workflow review] Add `.github/workflows/check-pj.yml`.
- [ ] [depends on: workflow addition] Run the same checks locally where possible:
  - `go -C tools/pj test ./...`
  - `go -C tools/pj vet ./...`
  - `goimports` check over `tools/pj/**/*.go`
  - chosen Go lint command
- [ ] [depends on: local checks] Fix any formatting or lint issues that are directly required by the new CI checks.
- [ ] [depends on: verification] Move this plan to `docs/exec-plan/done/`.

## Verification

- Confirm `docs/specs/github-projects-task-cli.md` describes the new `check-pj` CI contract.
- Confirm `.github/workflows/check-pj.yml` exists, targets pull requests to `main`, and uses the documented `paths` filter.
- Confirm the workflow checks `tools/pj` with formatting, linting, `go vet ./...`, and `go test ./...`.
- Confirm local verification commands pass before opening the implementation PR.
- After the implementation PR is created, confirm the `check-pj` GitHub Actions check runs for the latest PR head SHA.

## Expected Outcome

- `tools/pj` test coverage and Go quality checks become visible in every relevant PR.
- Future `tools/pj` work has a clear mechanical quality gate instead of relying only on local operator discipline.
- The workflow remains small and focused enough to evolve later if more workspace-local Go tools are added.
