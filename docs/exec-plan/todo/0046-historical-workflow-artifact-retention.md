# Historical workflow artifact retention

> **Execution**: Use `/execute-task` to implement this plan. After implementation is complete, use `/review-task` to prepare and create the PR.

## Objective

Remove resolved execution plans and local issues from the checked-out
repository so ordinary repository search returns only active implementation
trackers. Preserve reviewability through the existing plan PR, implementation
PR, and Git history; require durable product or workflow decisions to live in
their canonical locations instead of in completed task files.

Completion boundary: after the implementation PR merges, `docs/exec-plan/` has
only `todo/` and `docs/issues/` has only unresolved issues; neither checked-out
`done/` directory nor its generated-template counterpart exists. The linter
accepts and validates the new completion transition without weakening the
pre-implementation plan requirement or external GitHub issue closure check.

Addresses: docs/issues/0028-historical-workflow-artifact-retention.md

## Confirmed decision

- Delete a completed exec-plan and every linked resolved local issue in the
  implementation PR; do not move them to `done/`.
- Use the plan PR, implementation PR, and Git commit history as the audit
  trail. Do not duplicate plan bodies in commit messages: commit messages stay
  concise change summaries and links, while the PR diff preserves the complete
  artifact.
- Existing `docs/exec-plan/done/` and `docs/issues/done/` records are historical
  task trackers, not canonical decisions, and will be removed in this migration.

## Relevant references

- `AI_WORKFLOW.md:21-22` and `AI_WORKFLOW.md:45-51` currently require resolved
  local issues and completed plans to move into `done/`.
- `AGENTS.md:16-18` repeats the same execution-branch completion rule.
- `docs/specs/workflow-linter.md:67-75` and `:87-98` define `done/` as both the
  lifecycle destination and the signal used for linter checks.
- `tools/workflow-lint.sh:99-140` resolves plans across `todo/` and `done/`;
  `:475-505` forbids issue deletion; `:661-712` checks linked issue movement.
- `skills/execute-task/SKILL.md:19-20`, `skills/plan-execution/SKILL.md:19-21`,
  `docs/exec-plan/todo/README.md:1-5`, and `docs/issues/README.md:1-4` state or
  imply the existing retention policy.
- `skills/manage-workflow/templates/docs/exec-plan/todo/README.md` and
  `skills/manage-workflow/templates/docs/issues/README.md` propagate it to
  newly bootstrapped repositories.

## Change map

### Workflow contract and operator guidance

- (MODIFY) `AI_WORKFLOW.md`, `AGENTS.md`, `README.md`, and the active-directory
  README files: describe deletion of resolved local artifacts after verification
  and PR preparation; describe Git/PR history as the retrieval path.
- (MODIFY) `skills/execute-task/SKILL.md` and `skills/plan-execution/SKILL.md`:
  retain active-plan creation and plan/issue linkage, but replace move-to-done
  closeout instructions with deletion after the required evidence is captured.
- (MODIFY) `skills/manage-workflow/templates/docs/exec-plan/todo/README.md` and
  `skills/manage-workflow/templates/docs/issues/README.md`; (DELETE) their
  `done/` template directories and README files.

### Observable workflow-linter contract

- (MODIFY) `docs/specs/workflow-linter.md`: define a completed execution as a
  matching plan deletion in the branch diff, require linked local issue
  deletions in that same diff, and retain `Addresses:`-based external GitHub
  closure metadata validation.
- (MODIFY) `tools/workflow-lint.sh`: resolve an active matching plan normally;
  for closeout validation, recover the deleted matching plan from the merge-base
  side of the diff, parse its `Addresses:` content, and verify every linked
  local issue is deleted or explicitly remains open. Remove all `done/` lookup
  and rename requirements. Keep warnings `fixable` and non-blocking.
- (MODIFY) linter fixtures/tests, at their existing locations: cover active
  execution, compliant plan-and-linked-issue deletion, deleted-plan-only
  closeout, missing linked-issue deletion, and external issue closing metadata.

### Historical cleanup

- (DELETE) all tracked files under `docs/exec-plan/done/` and
  `docs/issues/done/`, then remove the now-empty directories.
- (MODIFY) every remaining tracked reference to those directories, including
  specs, templates, skills, and lint messages. Do not rewrite ADRs or historical
  Git commits; their past paths remain valid in history.

## Black-box contract changes

1. Before coding on a `feat/*` or `fix/*` branch, a matching active plan under
   `docs/exec-plan/todo/` remains mandatory.
2. A successful execution removes its matching plan from the working tree. If
   that plan names local `docs/issues/` paths in `Addresses:`, the same PR
   removes those resolved issue files unless its PR body explicitly says they
   remain open and why.
3. A resolved local tracker is discoverable through the PR/commit history, not
   through an in-tree `done/` directory. Canonical lasting decisions remain in
   ADRs, specs, `docs/project-plan.md`, code, or code comments.
4. External GitHub issues named in a removed plan still require `Closes` metadata
   (or an explicit remaining-open justification) in CI.

## Execution steps

1. Update the black-box workflow-linter specification first, including the
   deleted-plan diff source, linked-local-issue rule, and unchanged external
   issue rule. Update the lifecycle and human-facing documents in the same
   change so no surface directs an agent to `done/`.
2. Refactor `workflow-lint.sh` around the new lifecycle. Keep the active-plan
   check separate from closeout validation, and use a deleted plan's base-side
   content only to evaluate explicit links already declared by the author.
3. Update or add focused linter coverage before relying on the implementation;
   run it against both a normal active branch and closeout-diff fixtures.
4. Remove all existing completed plan and issue files and the template `done/`
   directories. Search tracked active guidance for stale paths and correct only
   current workflow contracts, not immutable ADR/history records.
5. Run the workflow quality gates and manually verify `rg` over the working
   tree no longer finds historical task bodies while `git log --all --
   docs/exec-plan` can retrieve a removed plan.

## Dependencies and parallelism

- Step 1 must precede linter implementation because the spec owns the observable
  rule.
- Documentation/template updates and focused test-fixture updates can proceed
  in parallel after the new contract is fixed.
- Historical cleanup follows the reference sweep and precedes final stale-path
  verification.

## Verification

- Run the focused workflow-linter tests and `./tools/workflow-lint.sh` in both
  pre-push and CI-style fixture modes as supported by the existing harness.
- Run `./tools/test-workflow-context-contract.sh`, `./tools/workflow-lint.sh
  --mode=pre-push`, and `git diff --check`.
- Confirm active `feat/*`/`fix/*` branches still receive a fixable warning when
  no matching `todo/` plan exists.
- Confirm a deleted matching plan with a deleted linked local issue is accepted,
  while omission of that issue remains fixable and an external linked issue
  still needs PR closing metadata.
- Confirm tracked-tree searches find no `docs/exec-plan/done/` or
  `docs/issues/done/` references except intentionally historical documentation
  if any is explicitly justified; confirm a removed artifact remains retrievable
  from `git log`/`git show`.
