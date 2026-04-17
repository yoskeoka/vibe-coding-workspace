## Summary

The workspace now documents the globally installed `ww` CLI as the default startup path in `AI_WORKFLOW.md`, `AGENTS.md`, `README.md`, and the core planning/execution skills, but some workflow-adjacent skills still instruct raw git startup.

Current mismatches:
- `skills/plan-project/SKILL.md`
- `skills/review-task/SKILL.md`
- `skills/manage-workflow/SKILL.md`

This leaves an inconsistent operator experience: depending on which workflow entrypoint an agent reads first, it may still choose `git fetch origin && git switch -c ...` instead of dogfooding global `ww`.

## Proposed Solution

Create a follow-up execution plan that:
- reviews whether each remaining skill should use global `ww` by default
- updates the applicable skills to `ww create` / `ww cd`
- explicitly documents any justified exceptions where raw git remains the correct startup path
- extends workflow docs and lint coverage if those skills become part of the dogfooding contract
- removes the temporary follow-up note in `AI_WORKFLOW.md` once the remaining migration work is complete or the remaining exceptions are documented in steady-state form

## Priority

Medium. This does not block the current workflow migration because the primary startup docs and core execution/planning skills are aligned, but it leaves avoidable ambiguity in adjacent workflows and weakens the dogfooding contract over time.

## Resolution

Resolved in the `docs/ww-workflow-skill-migration-follow-up` branch:

- Confirmed the remaining adjacent workflow skills now use the global `ww` startup path.
- Extended `docs/specs/ww-dogfooding-workflow.md` so `plan-project`, `review-task`, and `manage-workflow` are part of the documented dogfooding touchpoint set.
- Extended `tools/workflow-lint.sh` and `docs/specs/workflow-linter.md` so raw-git startup wording in those skills is covered by the fixable workflow-startup warning.
- Replaced the temporary follow-up note in `AI_WORKFLOW.md` with steady-state wording that points to the spec-defined migrated touchpoints.
