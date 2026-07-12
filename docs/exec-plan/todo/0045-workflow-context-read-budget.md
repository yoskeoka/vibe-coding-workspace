# Workflow context read budget
> **Execution**: Use `/execute-task` to implement this plan. After implementation is complete, use `/review-task` to prepare and create the PR.

Addresses: docs/issues/0026-workflow-context-read-budget.md

## Objective

Make the workspace workflow readable by task phase: agents can find universal
guards, the applicable lifecycle rule, and the owning procedure without loading
duplicated workflow prose, while all existing enforcement remains observable.

## Existing Implementation References

- `AGENTS.md:1-188` — universal guards and workspace routing.
- `AI_WORKFLOW.md:1-177` — lifecycle, branch, execution, and PR rules.
- `.github/PULL_REQUEST_TEMPLATE.md:1-134` — conditional PR metadata.
- `skills/{plan-execution,execute-task,review-task,post-task-review,manage-workflow,plan-project,triage-tasks}/SKILL.md` — phase procedures and duplicate lifecycle text.
- `docs/design-decisions/README.md:1-10`, `docs/design-decisions/adr/*.md`, `docs/lessons.md:1-8` — compact decision and active-exception records to preserve.
- `tools/workflow-lint.sh:626-809` and `docs/specs/workflow-linter.md` — structural validation conventions.

## Code Change Map

- `docs/specs/workflow-context-contract.md` (NEW) — define compact-read behavior and retained enforcement.
- `AGENTS.md`, `AI_WORKFLOW.md`, `.github/PULL_REQUEST_TEMPLATE.md` (MODIFY) — separate universal, lifecycle, and conditional author guidance.
- `skills/*/SKILL.md` (MODIFY) — give each workflow skill its procedure and lifecycle links.
- `docs/design-decisions/README.md`, `docs/design-decisions/adr/*.md` (NEW); `docs/design-decisions/adr.md` (DELETE) — indexed immutable Nygard ADR records.
- `docs/lessons.md` (MODIFY) — active exceptions register.
- `skills/manage-workflow/templates/docs/design-decisions/*` (MODIFY) — scaffold compact ADR layout.
- `tools/workflow-lint.sh`, `tools/test-workflow-context-contract.sh`, `docs/specs/workflow-linter.md` (MODIFY) — add focused structural checks and verification.

## Spec Changes

- Add `docs/specs/workflow-context-contract.md` for the observable workflow
  document ownership, ADR layout, lessons register, and retained lifecycle
  obligations.

## Sub-tasks

- [ ] [parallel] Establish contract, lifecycle links, and compact PR entrypoints.
- [ ] [parallel] Migrate ADRs, lessons, and workflow templates.
- [ ] [depends on: contract, migration] Add and run structural workflow-lint checks.

## Verification

- Run `./tools/test-workflow-context-contract.sh` and
  `./tools/workflow-lint.sh --mode=pre-push`.
- Verify each old ADR is indexed and the planning, execution, and review
  entrypoints identify their minimum reads.
