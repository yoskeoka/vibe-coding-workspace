# Pinact Rollout

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Adopt `pinact` as the standard way to pin and update GitHub Actions references across `vibe-coding-workspace` and its managed child projects, while keeping the initial rollout low-friction and compatible with each repository's existing workflow conventions.

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
- `ww` is a mixed case:
  - ordinary workflows such as `ci.yml` and `release.yml` use tag-based references
  - agentic review workflows are source-controlled as `.github/workflows/*.md` and compiled to `.lock.yml`
  - `ww/AGENTS.md` already says to edit `.md` sources, not `.lock.yml`
- `pinact` is not currently installed in this workspace checkout.
- `vim-learning-game` and `envdiff` are listed in `setup.sh` as managed repos, but they are not present in this local workspace snapshot, so their workflow inventory must be completed during execution.

## Options Considered

### Option A: add only a one-line reminder in an execution plan

- This is the lightest change, but execution plans are task-local artifacts and are a poor place for a durable cross-repo workflow rule.
- It would be easy for later workflow edits in other repos to miss the rule entirely.

### Option B: add a one-line rule to `AGENTS.md`, use the `pinact` CLI manually during workflow edits, and keep enforcement human-driven for now

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
- add a short durable rule in `AGENTS.md` and the workflow-facing spec
- introduce a shared `.pinact.yaml` where it helps
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

- ordinary workflows:
  - `.github/workflows/ci.yml`
  - `.github/workflows/copilot-setup-steps.yml`
  - `.github/workflows/release.yml`
- generated workflow family:
  - source: `.github/workflows/plan-review.md`, `impl-review.md`, `spec-code-sync.md`
  - generated outputs: corresponding `.lock.yml`
  - related lock data: `.github/aw/actions-lock.json`

### Repositories requiring inventory during execution

- `vim-learning-game`
- `envdiff`

## Spec Changes

### `AGENTS.md`

- Add one high-visibility rule in the workspace root:
  - when editing GitHub Actions workflows or composite actions, use `pinact` to pin or update `uses:` references rather than hand-editing version tags
- Clarify that repo-specific source-of-truth rules still apply, especially in `ww` where `.md` workflow sources must be edited instead of generated `.lock.yml` files

### `AI_WORKFLOW.md`

- Add a short workflow rule that GitHub Actions updates should pass through `pinact` as the standard pin/update step
- Keep the wording narrow so it remains an operator rule, not a blanket requirement for unrelated YAML edits

### Workflow-facing specs

- Update the workspace spec that governs workflow maintenance so the rule is durable outside `AGENTS.md`
- If execution shows the rule belongs in a more specific spec, prefer:
  - `docs/specs/workflow-linter.md` for workspace workflow discipline
  - repo-local workflow specs in child repos when they already define CI/workflow contracts

## Expected Code and Config Changes

### Workspace root

- Install or document the chosen `pinact` installation path for local operators
- Add `.pinact.yaml` if the default file targeting is insufficient or if exclusions are needed
- Run `pinact` on workspace workflows and review the diff
- Keep any intentionally unpinned references explicit in config or comments rather than silently skipping them

### Child repositories

- Add the same minimal operator guidance where durable and repo-appropriate
- Add repo-local `.pinact.yaml` only when the default target set is wrong or exclusions are required
- Run `pinact` and commit resulting workflow updates repo by repo

### `ww` special handling

- Determine the authoritative edit path before mutation:
  - edit workflow source `.md` files when they generate `.lock.yml`
  - regenerate `.lock.yml` and any related lock artifacts with the repo's existing tooling
- Do not run blind replacements directly on generated lock files unless the repo's documented flow requires it

## Sub-tasks

- [ ] Confirm the preferred `pinact` installation path for this workspace and document it if needed
- [ ] Add the durable operator rule to workspace `AGENTS.md` and `AI_WORKFLOW.md`
- [ ] Update the relevant workflow-maintenance spec in the workspace so the rule is not stored only in agent instructions
- [ ] Inventory each target repository's workflows and classify normal files versus generated sources/outputs
- [ ] Add `.pinact.yaml` in the workspace root if exclusions or custom file patterns are needed
- [ ] Run `pinact` in the workspace root and review/fix the resulting workflow diffs
- [ ] Roll the same process out to `ai-arena`
- [ ] Roll the same process out to `reversi-adventure`
- [ ] Roll the same process out to `ww` using its source-first workflow generation rules
- [ ] Clone or otherwise inspect `vim-learning-game` and `envdiff`, then apply the same rollout if workflows exist
- [ ] Decide after the manual rollout whether a second follow-up plan is justified for `pinact-action` validation-only CI

## Parallelism

- The workspace documentation updates can be drafted in parallel with repository workflow inventory.
- `ai-arena` and `reversi-adventure` are independent rollout targets once the shared policy wording is settled.
- `ww` should remain separate because its generated workflow path needs repo-specific handling.
- `vim-learning-game` and `envdiff` depend first on obtaining the missing local repository context.

## Risks and Mitigations

- Risk: a one-line `AGENTS.md` rule is visible but too weak to stay durable on its own.
  - Mitigation: record the same rule in a workflow-facing spec during the same rollout.
- Risk: `pinact` changes workflow behavior unexpectedly when tags and SHAs do not line up cleanly.
  - Mitigation: review diffs repo by repo, prefer conservative first-pass pinning, and use config-based exclusions for intentional exceptions.
- Risk: `ww` generated workflows drift if the wrong files are edited.
  - Mitigation: treat `.md` sources as authoritative and regenerate lock outputs through the documented repo flow.
- Risk: immediate `pinact-action` adoption adds credentials and CI recursion complexity before the manual path is understood.
  - Mitigation: defer CI automation to a separate follow-up decision after the CLI rollout.
- Risk: missing local clones hide additional workflow files in `vim-learning-game` or `envdiff`.
  - Mitigation: keep those repos explicitly in scope and block completion until their inventory is confirmed.
