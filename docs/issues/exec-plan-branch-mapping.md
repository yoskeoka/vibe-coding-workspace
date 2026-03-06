# Exec Plan ↔ Branch Mapping Convention

## Summary

The workflow assumes a mapping between branch names and exec-plan files (e.g., `feat/002-feature-x` → `docs/exec-plan/todo/002-feature-x.md`), but this mapping is implicit — never declared as a rule.

## Impact

- Workflow linter cannot check exec-plan existence for a branch
- Workflow linter cannot verify exec-plan completion (todo→done move) for feat/fix branches
- Agents may create branches and plans with non-matching names without detection

## Depends On

- `branch-naming-rule.md` — the naming convention must be formalized first

## Blocked

- Workflow linter check: exec plan existence (001-workflow-linter.md)
- Workflow linter check: exec plan completion (001-workflow-linter.md)

## Action Required

After branch naming rule is formalized, declare the branch→plan mapping rule in `AI_WORKFLOW.md` (e.g., "branches with a plan number MUST have a corresponding file in `docs/exec-plan/`").
