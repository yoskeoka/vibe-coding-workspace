# Sync child repos to the latest `slopless` Japanese-filter behavior

**Execution**: Use `/execute-task` to implement this plan.

## Objective

Bring every managed child repo onto the current workspace `slopless` changed-Markdown filtering behavior so Japanese-main Markdown is excluded from `slopless` review there as well.

This plan should end with:

- one workflow-only PR per target child repo
- each PR updating `.github/workflows/slopless.yml`
- the child workflows excluding Markdown files that contain Japanese-writing characters before invoking `slopless`
- repo-specific path-scope differences preserved instead of being overwritten by workspace-only assumptions

## Background

The desired child-repo behavior is to exclude changed Markdown files that contain any character from these Japanese-writing ranges before invoking `slopless`:

- Hiragana and fullwidth Katakana (`U+3040`-`U+30FF`)
- Katakana Phonetic Extensions (`U+31F0`-`U+31FF`)
- CJK Unified Ideographs Extension A (`U+3400`-`U+4DBF`)
- CJK Unified Ideographs (`U+4E00`-`U+9FFF`)
- CJK Compatibility Ideographs (`U+F900`-`U+FAFF`)
- Halfwidth Katakana (`U+FF65`-`U+FF9F`)

The child `slopless` workflows still inline their changed-Markdown detection inside `.github/workflows/slopless.yml` instead of calling a shared helper. A child rollout therefore needs workflow-file edits even when the intended behavior matches the latest workspace-side proposal.

The rollout remains narrow:

- copy the current behavior, not a broader redesign
- keep fixes limited to `.github/workflows/slopless.yml` unless a repo needs a minimal compatibility adjustment
- do not remediate existing `slopless` findings as part of this task

## Design Decisions

### 1. Keep the rollout behavior explicit in this plan, then port it into child workflow-local detection

- this plan itself defines the Japanese-writing ranges to mirror, so execution does not depend on an unmerged workspace branch
- child repos do not ship `tools/list-changed-markdown.sh`, so the rollout should update each workflow's inline detection block rather than adding a new shared script in this task
- preserve each child repo's current path scope such as `docs/development/` where present

### 2. Keep this as a workflow-only rollout

- in normal cases, modify only `.github/workflows/slopless.yml`
- if one repo needs an extra compatibility adjustment, document the exception and keep it minimal
- do not widen into child `docs/specs/`, helper-script distribution, or prose cleanup

### 3. Roll out to every managed child repo for parity

- the noise report came from Japanese-main docs in repos such as `ai-arena`, but parity matters more than repo-by-repo divergence
- apply the same filtering behavior across all managed child repos that already carry the copied `slopless` workflow

## Spec Changes

No child-repo spec updates are planned for this rollout.

This plan updates copied workflow behavior only. If implementation shows that one child repo documents `slopless` behavior locally and that documentation would drift, log the discrepancy and decide separately rather than broadening this rollout silently.

## Code Changes

### Workspace repo

- add this execution plan only

### Child repos

For each target repo, update:

- `.github/workflows/slopless.yml`

Expected workflow edits per repo:

- add the same Japanese-writing codepoint detection now used by the workspace helper
- exclude changed Markdown files containing those characters before the `slopless` invocation
- keep the repo's existing Markdown path scope intact unless the current workspace copy has already intentionally changed it for that repo

Target repos:

- `ai-arena`
- `dungeon-game-ai-arena`
- `envdiff`
- `reversi-adventure`
- `reversi-ai-arena`
- `vim-learning-game`
- `ww`

## PR / Rollout Plan

1. Capture the current child workflow shape and apply the Japanese-writing ranges defined in this plan.
2. Inspect each child repo's current `slopless` workflow and preserve any repo-local path-scope differences.
3. Update the inline changed-Markdown detection in each child workflow.
4. Run lightweight workflow validation in each child repo.
5. Push and open or update one PR per child repo.
6. Complete the bounded post-PR follow-up loop for each latest head SHA.

## Sub-tasks

- [ ] Capture the current child repos' inline detector shape and map the plan's Japanese-writing ranges into it
- [ ] [parallel] Update `ai-arena` to exclude Japanese-writing Markdown before `slopless`
- [ ] [parallel] Update `dungeon-game-ai-arena` to exclude Japanese-writing Markdown before `slopless`
- [ ] [parallel] Update `envdiff` to exclude Japanese-writing Markdown before `slopless`
- [ ] [parallel] Update `reversi-adventure` to exclude Japanese-writing Markdown before `slopless`
- [ ] [parallel] Update `reversi-ai-arena` to exclude Japanese-writing Markdown before `slopless`
- [ ] [parallel] Update `vim-learning-game` to exclude Japanese-writing Markdown before `slopless`
- [ ] [parallel] Update `ww` to exclude Japanese-writing Markdown before `slopless`
- [ ] [depends on: child repo updates] Validate workflow syntax and scope preservation in each repo
- [ ] [depends on: validation] Open or update the seven child-repo PRs
- [ ] [depends on: PR creation] Complete initial post-PR follow-up for each latest head SHA

## Parallelism

- The seven child-repo workflow edits are independent once the workspace behavior is fixed as the rollout source.
- Validation and PR creation can also proceed per repo in parallel, subject to serialized Git write operations within each repo worktree.

## Verification

For each child repo:

- run `pinact run .github/workflows/slopless.yml`
- run `python3 -c 'import yaml; yaml.safe_load(open(".github/workflows/slopless.yml"))'`
- run `git diff --check`
- confirm the diff is limited to `.github/workflows/slopless.yml` unless a documented compatibility exception exists
- when practical, inspect the inline detector block to confirm it preserves the repo's existing path scope while adding the Japanese-writing character exclusion

## Expected Outcome

- all seven managed child repos receive PRs aligning their `slopless` workflow behavior with the current workspace Japanese-filter logic
- Japanese-main Markdown stops generating avoidable `slopless` warnings in child repos without relying on path naming alone
- the rollout stays narrowly scoped to copied workflow parity rather than becoming a broader workflow/tooling redesign
