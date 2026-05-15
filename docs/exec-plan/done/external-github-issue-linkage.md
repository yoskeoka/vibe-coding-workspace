# External GitHub Issue Linkage

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Make external repo-native feedback as hard to forget as local `docs/issues/` follow-up by carrying the issue target from the execution plan into the implementation PR body.

## Current State

- `docs/issues/` is the workspace-local issue log, but the workflow text does not clearly say that external feedback may stay canonical in GitHub Issues instead.
- Execution plans can already use `Addresses:` for local issue files, but they do not clearly define how to record external GitHub issues.
- The PR template and workflow skills do not explicitly require implementation PRs to include `Closes` entries for external GitHub issues declared in the plan.
- `workflow-lint` can enforce linked local issue moves, but it does not currently check external issue closure metadata in PR bodies.

## Options Considered

### Option A: skill-only reminder

- Add another reminder sentence to the workflow skills.
- This is cheap, but the workspace already has multiple reminder layers and the recent miss happened anyway.

### Option B: docs and template only

- Clarify the rule in `AI_WORKFLOW.md` and the PR template.
- This improves reviewer visibility, but it still relies on humans catching the omission after the fact.

### Option C: explicit plan linkage plus PR/linter enforcement

- Define external GitHub issues as valid `Addresses:` targets using full issue URLs.
- Require plan PRs to repeat those issues under `Issues`.
- Require implementation PRs to add explicit `Closes` entries for those same issues.
- Extend `workflow-lint` to warn in CI when a completed plan links an external GitHub issue but the PR body does not close it.

### Recommended Approach

Adopt **Option C**.

It keeps the canonical tracker where outside contributors already interact, preserves the workspace-local `docs/issues/` role for internal workflow notes, and adds a mechanical guard for the PR-body step that GitHub issue closure actually depends on.

## Spec Changes

### `AI_WORKFLOW.md` and `AGENTS.md`

- Clarify that `docs/issues/` is for workspace-local development follow-up, while external feedback can stay canonical in the target repo's GitHub Issues.
- Extend `Addresses:` guidance so plans may declare external GitHub issues using full issue URLs.
- Require implementation PRs to include explicit closing keywords for those linked external issues unless the PR body explains why they remain open.

### Workflow skills and PR template

- Update `plan-execution` so plan authors record external GitHub issues in `Addresses:` and repeat them in the plan PR body under `Issues`.
- Update `execute-task` and `review-task` so implementation PRs must include matching `Closes` entries.
- Update `.github/PULL_REQUEST_TEMPLATE.md` with an explicit `Closes` section and checklist coverage.

### `docs/specs/workflow-linter.md`

- Define the external-GitHub-issue `Addresses:` convention.
- Add a CI-only `fixable` warning when a completed plan links external GitHub issues but the PR body lacks matching closing keywords.

## Expected Code Changes

### `tools/workflow-lint.sh`

- Parse full GitHub issue URLs from `Addresses:` lines in completed execution plans.
- Detect the current GitHub repo slug so same-repo issues can be satisfied by `Closes #<number>`.
- Emit a `fixable` warning in CI when the PR body neither closes the linked issue nor explicitly justifies leaving it open.

## Sub-tasks

- [x] Clarify the local-vs-external issue tracking rule in workflow docs
- [x] Define the external GitHub issue `Addresses:` convention
- [x] Update workflow skills and PR template with plan-PR and implementation-PR responsibilities
- [x] Extend workflow-linter spec and implementation with CI validation for linked external issue closure metadata
- [x] Verify the new linter behavior with positive and negative examples
