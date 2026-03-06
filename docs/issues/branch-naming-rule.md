# Branch Naming Convention Rule

## Summary

Current workflow docs (AI_WORKFLOW.md, AGENTS.md) use branch naming examples (`plan/002-feature-x`, `feat/002-feature-x`, `fix/003-bug-x`) but do not declare a mandatory naming rule.

## Impact

- Workflow linter cannot enforce branch naming (no rule to enforce)
- Exec plan existence/completion checks depend on branch→plan file mapping, which requires a declared naming convention

## Blocked

- Workflow linter check: branch naming (001-workflow-linter.md)
- Workflow linter check: exec plan existence (001-workflow-linter.md)
- Workflow linter check: exec plan completion (001-workflow-linter.md)

## Proposed Rule (draft)

```
Branch names MUST match: <type>/<description>
  type: plan | feat | fix | chore | docs
  description: free-form, but if associated with an exec-plan, SHOULD start with the plan number (e.g., 002-feature-name)
```

## Action Required

Formalize the rule in `AI_WORKFLOW.md` and `AGENTS.md`, then unblock the linter checks.
