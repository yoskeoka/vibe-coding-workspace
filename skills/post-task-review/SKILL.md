---
name: post-task-review
description: Capture durable follow-up knowledge after significant completed work.
---

# Post-task Review

Use after significant execution, investigation, merge, review correction, or
workaround; skip trivial edits. See [post-task review](../../AI_WORKFLOW.md#post-task-review).

## Procedure

1. Identify human intent, corrected assumptions, and concrete findings not
   already recorded in the plan, specs, ADRs, or issues.
2. Present prioritized findings and ask before creating new local issues. Do not
   create speculative issues or duplicate an external canonical tracker.
3. Promote a recurring rule to its owning spec, skill, linter, or ADR whenever
   possible. Keep `docs/lessons.md` only for unresolved recurring exceptions;
   each entry needs its remediation link and review trigger.
4. Propose `AGENTS.md` or skill changes only when the knowledge is universal and
   not already canonical. Preserve project-specific guidance and the compact
   read contract.
