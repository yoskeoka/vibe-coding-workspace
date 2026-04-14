# Exec-Plan Filename Convention

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Resolve the remaining inconsistency between the current workflow rule for exec-plan filenames and the older numbered examples that still appear in live docs/templates. Establish a single source of truth for active naming guidance, then migrate active plan files and references so future reviews do not get conflicting signals.

Addresses:
- PR review comment: `https://github.com/yoskeoka/vibe-coding-workspace/pull/57#discussion_r3082505684`
- Existing inconsistency between `AI_WORKFLOW.md` / `skills/plan-execution/SKILL.md` and `docs/exec-plan/todo/README.md`

## Current State

The repository already formalized a non-numeric naming rule:

- `AI_WORKFLOW.md` says exec-plan filenames match the branch description and use free-form kebab-case with no numeric prefixes.
- `skills/plan-execution/SKILL.md` repeats the same rule and explicitly rejects names like `004a-feature-name.md`.
- `docs/exec-plan/todo/README.md` and `skills/manage-workflow/templates/docs/exec-plan/todo/README.md` still say `XXX-description.md` with examples like `001-init.md`.
- Active plan files still include mixed conventions such as `001-kb-video-ingest-pipeline.md`, `003-workspace-project-linking.md`, `004-project-owner-scope-config.md`, and non-numeric files in the same directory.

## Decision

Keep the already-formalized non-numeric convention as the source of truth for active workflow docs and plan filenames.

Rationale:
- This matches the current branch naming and exec-plan mapping rules already declared in `AI_WORKFLOW.md`.
- It keeps branch names and plan filenames aligned without extra numbering semantics.
- Numeric prefixes are already treated as historical baggage in `docs/exec-plan/done/branch-naming-and-mapping.md`.

Historical completed artifacts may continue to mention old numbered examples when they are documenting prior state, but active instructions and active plan files should no longer use them.

## Code Changes

### Live workflow docs and templates

- Update `docs/exec-plan/todo/README.md` to describe the non-numeric kebab-case convention and reference branch-to-plan name matching.
- Update `skills/manage-workflow/templates/docs/exec-plan/todo/README.md` so newly bootstrapped repos inherit the current rule.
- Audit workflow-facing docs/skills for any remaining live references that tell users to create numbered exec-plan filenames, and update them if found.

### Active exec-plan files

- Rename active todo plans that still carry numeric prefixes to non-numeric names that match the intended branch description.
- Update any references that point at the renamed plan files.
- Verify the resulting active set uses one convention consistently.

## Spec Changes

- Update workflow documentation/spec text where needed so live instructions consistently describe non-numeric exec-plan filenames.
- If the execution work changes `docs/specs/workflow-linter.md` expectations or adds guardrails for this convention, update that spec first; otherwise, no linter behavior change is required.

## Design Decisions

Past decision:
- `docs/exec-plan/done/branch-naming-and-mapping.md` already chose non-numeric filenames and explicitly planned to rename numeric exec-plan files.

Apply the same reasoning here rather than introducing a second convention.

## Sub-tasks

- [ ] [parallel] Update live workflow-facing docs/templates to state the non-numeric exec-plan filename convention
- [ ] [parallel] Rename active numeric exec-plan files in `docs/exec-plan/todo/`
- [ ] [depends on: both above] Update any in-repo references affected by the renames and verify there are no remaining live numbered instructions
- [ ] [depends on: all above] Run repo checks relevant to the touched files and prepare the execution PR
