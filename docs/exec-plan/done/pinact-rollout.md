# Pinact Rollout

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Adopt `pinact` as the standard way to pin and update GitHub Actions references across `vibe-coding-workspace` and its managed child projects, using a deliberately small durable policy surface: add a one-line rule to each relevant repository's `AGENTS.md`, then run `pinact` in every repository that actually uses GitHub Actions and open PRs with the resulting diffs.

This plan covers:

- `vibe-coding-workspace`
- `ww`
- `ai-arena`
- `reversi-adventure`
- `vim-learning-game`
- `envdiff`

## Current State

- The workspace already treats workflow policy as a first-class deliverable, and high-visibility operational rules belong in `AGENTS.md`.
- The workspace root, `ai-arena`, and `reversi-adventure` currently use tag-based GitHub Action references such as `actions/checkout@v6`, `actions/cache@v4`, and `dtolnay/rust-toolchain@stable`.
- `ww` contains both ordinary workflows and `gh aw` generated review workflows, but the `gh aw` family is already pinned and is explicitly out of scope for this rollout.
- `pinact` should be treated as already installed for this rollout; installation/distribution is not part of the plan.
- `vim-learning-game` and `envdiff` are listed in `setup.sh` as managed repos, but they are not present in this local workspace snapshot, so their workflow inventory must be completed during execution.

## Options Considered

### Option A: add only a one-line reminder in an execution plan

- This is the lightest change, but execution plans are task-local artifacts and are a poor place for a durable cross-repo workflow rule.
- It would be easy for later workflow edits in other repos to miss the rule entirely.

### Option B: add a one-line rule to each relevant repository's `AGENTS.md`, use the `pinact` CLI manually, and keep enforcement human-driven for now

- This gives the highest visibility at the lowest operational cost.
- It matches the current preference for low-friction, low-token workflow rules before adding more automation.
- It avoids GitHub App setup and avoids introducing bot-driven workflow mutations before the team knows the repo-specific edge cases.

### Option C: adopt `pinact-action` immediately across repositories

- This can automate validation or auto-fix in CI, but it adds token/secret management and workflow-chaining concerns immediately.
- The upstream docs recommend a GitHub App installation token; using only `GITHUB_TOKEN` cannot update workflows in the fully automated way this action is designed for.
- This is more appropriate as a second phase after the CLI rollout proves stable.

### Recommended Approach

Adopt **Option B** first.

That means:

- make `pinact` the documented operator path for GitHub Actions updates
- add a short durable rule in each relevant repository's `AGENTS.md`
- run `pinact` repo by repo during the rollout
- defer `pinact-action` until the manual CLI flow has proven worthwhile

This keeps setup cost low, aligns with the workspace's "workflow as product" goal, and avoids overcommitting to automation before repo-specific exceptions are understood.

## Scope Inventory

### Workspace root

- `.github/workflows/check-pj.yml`
- `.github/workflows/kb-pages.yml`
- `.github/workflows/shellcheck.yml`
- `.github/workflows/sync-workflow-to-child-repos.yml`
- `.github/workflows/workflow-lint.yml`

### `ai-arena`

- `.github/workflows/go-ci.yml`

### `reversi-adventure`

- `.github/workflows/ci.yml`

### `ww`

- ordinary workflows only:
  - `.github/workflows/ci.yml`
  - `.github/workflows/copilot-setup-steps.yml`
  - `.github/workflows/release.yml`
- explicit non-scope:
  - `.github/workflows/plan-review.md`
  - `.github/workflows/impl-review.md`
  - `.github/workflows/spec-code-sync.md`
  - corresponding `.lock.yml`
  - `.github/aw/actions-lock.json`

### Repositories requiring inventory during execution

- `vim-learning-game`
- `envdiff`

## Documentation Changes

### `AGENTS.md`

- Add one high-visibility rule in every repository that this rollout edits:
  - when editing GitHub Actions workflows or composite actions, use `pinact` to pin or update `uses:` references rather than hand-editing version tags
- Keep the wording short and operational.
- Do not add installation instructions or environment setup guidance.

## Expected Code and Config Changes

### Workspace root

- Update `AGENTS.md` with the one-line `pinact` rule
- Run `pinact` on workspace workflows and review the diff
- Add `.pinact.yaml` only if default targeting is insufficient or exclusions are needed

### Child repositories

- Update each edited repository's `AGENTS.md` with the same minimal operator guidance
- Run `pinact` and commit resulting workflow updates repo by repo
- Add repo-local `.pinact.yaml` only when the default target set is wrong or exclusions are required

### `ww` special handling

- Limit the rollout to ordinary workflow YAML files.
- Do not touch `gh aw` source `.md`, generated `.lock.yml`, or `.github/aw/actions-lock.json` in this plan.

## Sub-tasks

- [ ] Add the one-line `pinact` rule to `AGENTS.md` in each repository that this rollout edits
- [ ] Inventory each target repository's workflow files and separate in-scope ordinary YAML workflows from explicit exclusions
- [ ] Add `.pinact.yaml` in the workspace root if exclusions or custom file patterns are needed
- [ ] Run `pinact` in the workspace root and review/fix the resulting workflow diffs
- [ ] Roll the same process out to `ai-arena`
- [ ] Roll the same process out to `reversi-adventure`
- [ ] Roll the same process out to `ww` for ordinary workflow YAML files only, excluding `gh aw` sources and generated lock files
- [ ] Clone or otherwise inspect `vim-learning-game` and `envdiff`, then apply the same rollout if workflows exist
- [ ] Open or update PRs for each repository with the `AGENTS.md` line and `pinact` diffs together
- [ ] Decide after the manual rollout whether a second follow-up plan is justified for `pinact-action` validation-only CI

## Parallelism

- The workspace `AGENTS.md` update can be drafted in parallel with repository workflow inventory.
- `ai-arena` and `reversi-adventure` are independent rollout targets once the shared policy wording is settled.
- `ww` should remain separate because its in-scope and out-of-scope workflow files differ.
- `vim-learning-game` and `envdiff` depend first on obtaining the missing local repository context.

## Risks and Mitigations

- Risk: a one-line `AGENTS.md` rule is visible but easy to miss in child repos if only the workspace root is updated.
  - Mitigation: update every edited repository's own `AGENTS.md`, not just the workspace root.
- Risk: `pinact` changes workflow behavior unexpectedly when tags and SHAs do not line up cleanly.
  - Mitigation: review diffs repo by repo, prefer conservative first-pass pinning, and use config-based exclusions for intentional exceptions.
- Risk: `ww` review workflow assets get touched accidentally even though they are already pinned and out of scope.
  - Mitigation: exclude the `gh aw` source and generated files explicitly from inventory, execution, and review.
- Risk: immediate `pinact-action` adoption adds credentials and CI recursion complexity before the manual path is understood.
  - Mitigation: defer CI automation to a separate follow-up decision after the CLI rollout.
- Risk: missing local clones hide additional workflow files in `vim-learning-game` or `envdiff`.
  - Mitigation: keep those repos explicitly in scope and block completion until their inventory is confirmed.
