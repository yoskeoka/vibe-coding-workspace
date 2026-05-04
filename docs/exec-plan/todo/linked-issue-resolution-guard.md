# linked-issue-resolution-guard
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Make resolved local issues harder to forget by turning the issue-to-plan relationship into explicit workflow data and a mechanical review signal.

Addresses: `docs/issues/linked-issue-resolution-guard.md`

## Current State

- `AI_WORKFLOW.md` and `AGENTS.md` already say resolved local issues should move from `docs/issues/` to `docs/issues/done/`.
- Some completed plans already use `Addresses:`, but it is not yet a required workflow contract for execution plans.
- `.github/PULL_REQUEST_TEMPLATE.md` links plans and issues, but it does not explicitly require that linked resolved local issues be moved to `docs/issues/done/` or justified.
- `tools/workflow-lint.sh` checks issue lifecycle only after an issue file has already been changed; it does not use plan metadata to detect "this issue should have been moved in this PR".

## Options Considered

### Option A: skill-only reminder

- Update workflow skills to remind the agent again about issue moves.
- This is low cost, but it repeats a rule that already exists and still depends on the operator remembering the final manual step.

### Option B: PR-template/process reminder only

- Add stronger checklist wording to the PR template and maybe AI workflow docs.
- This improves reviewer visibility, but it still relies on humans noticing the mismatch and does not give local mechanical feedback before push.

### Option C: explicit plan linkage plus linter/template enforcement

- Require plan-level `Addresses:` metadata for local issues that the work resolves.
- Make the PR template ask whether linked local issues were moved to `docs/issues/done/` or intentionally left open with justification.
- Extend `workflow-lint` so execution work that claims a local issue but forgets the matching done move gets a `fixable` warning.

### Recommended Approach

Adopt **Option C**.

It keeps the canonical rule in repo-visible workflow artifacts instead of only in skills, gives reviewers a concrete checklist hook, and adds a machine-readable signal that can be checked before push and in CI.

## Spec Changes

### `AI_WORKFLOW.md`

- Define that execution plans resolving tracked local issues should declare them in an `Addresses:` section.
- Clarify that an implementation PR resolving a linked local issue should move that issue to `docs/issues/done/` in the same execution branch unless the PR body explains why the issue remains open.

### `docs/specs/workflow-linter.md`

- Document the `Addresses:` convention as workflow-lint input for execution work.
- Add a new `fixable` check for linked local issue resolution:
  - when execution work declares `Addresses: docs/issues/<name>.md`
  - and the PR/branch diff appears to implement that work
  - but the issue file is not moved to `docs/issues/done/<name>.md`
- Clarify scope boundaries so this remains a workflow rule, not a generalized semantic code-analysis rule.

### `.github/PULL_REQUEST_TEMPLATE.md`

- Add a checklist item confirming that resolved linked local issues were moved to `docs/issues/done/`, or the PR body explains why they remain open.
- Tighten the `Plan / Issues` guidance so linked local issue paths are expected to match the plan's `Addresses:` section when applicable.

## Expected Code Changes

### `tools/workflow-lint.sh`

- Parse `Addresses:` entries from the matching execution plan file for `feat/*` and `fix/*` branches.
- Recognize local issue paths under `docs/issues/`.
- Emit a `fixable` warning when the branch appears to resolve a linked local issue but does not move that file to `docs/issues/done/`.
- Keep the check narrow and explainable so false positives remain rare and understandable.

## Design Decisions

Past decisions:

- Workflow-lint warnings are intentionally non-blocking but classified, with `fixable` meaning "normally resolve before push/PR".
- The repo already prefers machine-visible workflow contracts in docs/specs/template/linter over relying on agent memory alone.

Apply the same reasoning here: the missing safeguard is not another reminder sentence, but a clearer contract and a linter-visible signal.

## Sub-tasks

- [ ] Define the `Addresses:` contract for execution plans and linked local issue closure in `AI_WORKFLOW.md`
- [ ] Update `docs/specs/workflow-linter.md` with the new linked-issue-resolution check and its guardrails
- [ ] Update `.github/PULL_REQUEST_TEMPLATE.md` so reviewers must confirm the linked-issue closure state
- [ ] Implement the new `tools/workflow-lint.sh` check for linked local issue resolution
- [ ] Verify the new check with at least one positive and one negative repo-state example

## Parallelism

- `AI_WORKFLOW.md`, `docs/specs/workflow-linter.md`, and `.github/PULL_REQUEST_TEMPLATE.md` can be drafted in parallel once the contract wording is settled.
- `tools/workflow-lint.sh` depends on the contract wording because the check should enforce the documented behavior, not invent new policy.
- Verification depends on the linter implementation.

## Risks and Mitigations

- Risk: the linter overreaches and guesses issue resolution from unrelated diffs.
  - Mitigation: scope the new check to linked local issues declared in the matching execution plan instead of trying to infer issue closure from arbitrary branch content.
- Risk: some linked issues intentionally stay open across multiple PRs.
  - Mitigation: let the workflow and PR template permit an explicit justification when a linked issue remains open.
- Risk: `Addresses:` becomes required in situations where no local issue file exists.
  - Mitigation: keep the contract conditional: use `Addresses:` when the work resolves tracked local issues, not for every plan unconditionally.
