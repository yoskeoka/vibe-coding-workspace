# Workflow Item Sequence Prefix

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Introduce a durable sequence-prefix rule for active workflow plan and issue files so backlog age and creation order remain visible without opening Git history or PR lists.

This plan updates the workspace workflow contract first, then rolls the same rule into child repositories that use the workflow.

## Scope

In scope:

- workspace workflow docs, specs, templates, and skills
- `tools/workflow-lint.sh` checks for active plan / issue naming
- renaming active workspace plan and issue files to the new numbered format
- applying the same workflow rule and active-file renames to child repositories

Out of scope:

- renaming historical files already under `docs/exec-plan/done/` or `docs/issues/done/`
- changing branch naming to include sequence numbers
- reconstructing historical ordering beyond what existing `done/` counts and file creation dates can support

## Required Rule

### Active plan / issue filenames

- Active execution plans under `docs/exec-plan/todo/` MUST use `<sequence>-<name>.md`.
- Active local issues under `docs/issues/` MUST use `<sequence>-<name>.md`.
- Sequence numbers use four-digit zero padding from `0001` through `9999`.
- Sequence numbers at `10000` or above MUST be written without zero padding.
- Mixed states such as `9999-foo.md` and `12345-bar.md` are valid.
- `README.md` files are exempt.

### Historical files

- Files already in `docs/exec-plan/done/` and `docs/issues/done/` remain unchanged even if they use older naming.
- Historical navigation for completed work is expected to rely primarily on PR history rather than retroactive file renames.

### Branch mapping

- Branch naming remains `plan/<name>`, `feat/<name>`, `fix/<name>`, `chore/<name>`, or `docs/<name>`.
- For `feat/<name>` and `fix/<name>`, the matching exec-plan is the single active or completed plan whose basename suffix is `-<name>.md`.
- Workflow docs and lint rules must describe this as `docs/exec-plan/{todo,done}/<sequence>-<name>.md`, not as a branch name that repeats the sequence.

### `docs/lessons.md`

- Add one durable rule near the top of `docs/lessons.md` that new lessons are appended at the end of the file.
- Reinforce the same append-at-end rule in workflow docs and relevant skills.

## Renaming Rule For Existing Active Items

For this repository and each child repository:

1. Count existing completed files under `docs/exec-plan/done/` and `docs/issues/done/`, excluding `README.md`.
2. Collect only active files under `docs/exec-plan/todo/` and `docs/issues/`, excluding `README.md`.
3. Sort active files by creation date from Git history (`git log --diff-filter=A --follow`).
4. If two files share the same creation timestamp, break the tie lexicographically by current path.
5. Assign new sequence numbers starting from `done_count + 1`, oldest active file first.
6. Rename only active files; never rename files already in `done/`.

Workspace current baseline at planning time:

- exec-plan done count: `37`, so active plan numbering starts at `0038`
- issue done count: `18`, so active issue numbering starts at `0019`

## Relevant Prior Decisions

Past decision: the workflow is a first-class product and child projects are part of the proving ground. Apply the same reasoning here by treating this naming rule as a workspace-wide workflow contract, not as a workspace-root-only preference.

Past decision: normal workflow startup should dogfood the released global `ww` binary. Apply the same discipline here by keeping branch UX stable and moving the ordering signal into file naming and lint rules instead of forcing sequence numbers into branch names.

## Code Changes

### Workspace root

- Update `tools/workflow-lint.sh` so it can:
  - validate active plan naming under `docs/exec-plan/todo/`
  - validate active issue naming under `docs/issues/`
  - resolve `feat/*` / `fix/*` exec-plan mapping through the numbered filename suffix
  - keep `README.md` exempt
  - tolerate historical non-numbered files in `done/`
- Update any helper logic or path-resolution code that assumes active filenames equal the branch description exactly.

### Child repositories

- Apply the same workflow-lint behavior and active-file renames in each child repo:
  - `ai-arena`
  - `reversi-adventure`
  - `vim-learning-game`
  - `ww`
  - `envdiff`
- Use repo-local PRs for child rollouts where local docs/tools must change; do not assume the existing workflow-sync path can propagate these changes automatically.

## Spec Changes

### Workspace specs / docs

- Update `AI_WORKFLOW.md` and `AGENTS.md` so plan / issue numbering, branch mapping, and `docs/lessons.md` append behavior are explicit.
- Update `docs/specs/workflow-linter.md` with the new mechanically enforced checks and suffix-based exec-plan resolution.
- Update `docs/exec-plan/todo/README.md`, `docs/issues/README.md`, and workflow-facing template files so newly bootstrapped repos inherit the numbered rule.
- Update `docs/lessons.md` with the durable append-at-end rule near the top.

### Skills

- Update workflow-facing skills that mention plan / issue creation or lesson updates:
  - `skills/plan-execution/SKILL.md`
  - `skills/execute-task/SKILL.md`
  - `skills/post-task-review/SKILL.md`
  - `skills/triage-tasks/SKILL.md`
  - `skills/manage-workflow/SKILL.md`
- Ensure examples and instructions consistently use numbered active plan / issue filenames and `docs/lessons.md` append-at-end wording.

## Execution Strategy

1. Update the workspace workflow contract and linter first.
2. Rename active workspace plan / issue files according to the documented numbering rule.
3. Verify workspace workflow docs, skills, and linter behavior locally.
4. Roll the same rule into each child repo, including active-file renames and repo-local workflow docs / tooling updates.
5. Open or update child-repo PRs so the workspace family converges on the same workflow contract.

## Sub-tasks

- [ ] Update workspace workflow docs and specs for numbered active plan / issue files
- [ ] Add the durable `docs/lessons.md` append-at-end rule and mirror it in relevant skills
- [ ] [parallel] Extend `tools/workflow-lint.sh` to enforce active plan / issue numbering and suffix-based exec-plan lookup
- [ ] [parallel] Update workflow templates and bootstrapping docs so new repos inherit the numbered rule
- [ ] Rename active workspace plan files using the done-count plus creation-date rule
- [ ] Rename active workspace issue files using the done-count plus creation-date rule
- [ ] Verify workspace docs and lint behavior locally
- [ ] Roll the same workflow rule and active-file renames into `ai-arena`
- [ ] Roll the same workflow rule and active-file renames into `reversi-adventure`
- [ ] Roll the same workflow rule and active-file renames into `vim-learning-game`
- [ ] Roll the same workflow rule and active-file renames into `ww`
- [ ] Roll the same workflow rule and active-file renames into `envdiff`

## Verification

- `tools/workflow-lint.sh --mode=pre-push` warns on malformed active plan / issue filenames and passes on compliant ones
- `tools/workflow-lint.sh` finds the matching exec-plan for `feat/<name>` and `fix/<name>` through numbered filename suffixes
- Workspace active items are renamed exactly from the computed numbering order and `done/` files remain untouched
- `docs/lessons.md`, workflow docs, and skills all agree that lessons append at the end
- Each child repo receives the same workflow rule, lint behavior, and active-file numbering migration
