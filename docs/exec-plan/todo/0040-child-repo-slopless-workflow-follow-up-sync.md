# Sync child repos to the latest `slopless` workflow

**Execution**: Use `/execute-task` to implement this plan.

Addresses: docs/issues/0024-child-repo-slopless-workflow-follow-up-sync.md

## Objective

Bring every managed child repo onto the current workspace `.github/workflows/slopless.yml` behavior without widening scope beyond the workflow file unless a repo-specific compatibility fix is strictly required.

This plan should end with:

- one workflow-only PR per target child repo
- each PR updating `.github/workflows/slopless.yml` to match the workspace version
- repo-specific exceptions called out explicitly instead of being folded into silent extra changes

## Background

The first child-repo `slopless` rollout landed an earlier workflow shape. The workspace workflow later gained follow-up fixes around:

- explicit Node setup instead of relying on the runner default toolchain
- batching `npx` execution over the changed file set
- scanning PR comments beyond the first page before marker-comment upsert
- narrower workflow permissions
- explicit subprocess and GitHub API timeouts

The issue already constrains this follow-up to a narrow rollout: port the current workspace workflow behavior to each child repo and avoid unrelated fixes.

## Design Decisions

### 1. Use the workspace workflow file as the source of truth

- copy the current workspace `.github/workflows/slopless.yml` shape into each target repo
- do not re-design the workflow independently per repo unless a compatibility break forces a documented exception

### 2. Keep this as a workflow-only rollout

- in normal cases, modify only `.github/workflows/slopless.yml`
- if one repo needs an extra compatibility adjustment, document why that exception is necessary and keep the delta minimal
- do not expand into spec, docs, or pre-existing `slopless` findings remediation

### 3. Prefer direct child-repo rollout PRs over sync-automation changes

- the workspace already contains the desired workflow shape
- changing rollout automation would enlarge scope beyond the issue's narrow follow-up intent
- manual or scripted file-sync into dedicated child-repo PRs is acceptable for this one-shot convergence task

## Spec Changes

No `docs/specs/` updates are planned for this rollout.

The execution branch should stay workflow-only unless implementation reveals that one repo's existing workflow contract is incompatible with the current workspace file. If that happens, record the exception in the relevant PR and, if needed, open a follow-up issue rather than broadening this plan silently.

## Code Changes

### Workspace repo

- no runtime code change expected
- add this execution plan only

### Child repos

For each target repo, update:

- `.github/workflows/slopless.yml`

Target repos:

- `ai-arena`
- `dungeon-game-ai-arena`
- `envdiff`
- `reversi-adventure`
- `reversi-ai-arena`
- `vim-learning-game`
- `ww`

## PR / Rollout Plan

1. Create one implementation branch and PR in each target child repo.
2. Copy the current workspace `slopless` workflow into the child repo.
3. Run lightweight workflow-file validation in that repo.
4. Push and open/update the child PR.
5. Complete the bounded post-PR follow-up loop for the latest head SHA.
6. Repeat until every target repo has a reviewable rollout PR.

## Sub-tasks

- [ ] Capture the current workspace `.github/workflows/slopless.yml` as the rollout source
- [ ] [parallel] Update `ai-arena` to the latest workflow file
- [ ] [parallel] Update `dungeon-game-ai-arena` to the latest workflow file
- [ ] [parallel] Update `envdiff` to the latest workflow file
- [ ] [parallel] Update `reversi-adventure` to the latest workflow file
- [ ] [parallel] Update `reversi-ai-arena` to the latest workflow file
- [ ] [parallel] Update `vim-learning-game` to the latest workflow file
- [ ] [parallel] Update `ww` to the latest workflow file
- [ ] [depends on: child repo updates] Validate workflow syntax and pinned actions in each repo
- [ ] [depends on: validation] Open or update the seven child-repo PRs
- [ ] [depends on: PR creation] Complete initial post-PR follow-up for each latest head SHA

## Parallelism

- The seven child-repo file-sync updates are independent once the workspace source file is fixed.
- Validation and PR creation can also proceed per repo in parallel, subject to serialized Git write operations within each individual worktree.

## Verification

For each child repo:

- run `pinact run .github/workflows/slopless.yml`
- run `python3 -c 'import yaml; yaml.safe_load(open(".github/workflows/slopless.yml"))'`
- run `git diff --check`
- confirm the diff is limited to `.github/workflows/slopless.yml` unless a documented compatibility exception exists

## Expected Outcome

- all seven managed child repos receive PRs that align their `slopless` workflow with the current workspace version
- the rollout stays narrowly scoped to workflow-file parity
- any repo-specific incompatibility is isolated, documented, and not silently generalized into broader workflow changes
