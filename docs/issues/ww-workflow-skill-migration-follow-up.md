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

## Priority

Medium. This does not block the current workflow migration because the primary startup docs and core execution/planning skills are aligned, but it leaves avoidable ambiguity in adjacent workflows and weakens the dogfooding contract over time.
