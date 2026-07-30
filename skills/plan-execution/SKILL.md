---
name: plan-execution
description: Plan a non-trivial feature, bug fix, or workflow change in docs/exec-plan/todo.
---

# Plan Execution

Use this for work that changes specs, spans multiple files/steps, has an
architectural choice, or is uncertain. A typo or isolated trivial fix may skip
the plan. Follow [execution planning](../../AI_WORKFLOW.md#execution-planning)
and [PR follow-up](../../AI_WORKFLOW.md#pr-and-follow-up) for shared lifecycle rules.

## Procedure

1. Create `plan/<name>` with `ww`, then read `docs/project-plan.md`, relevant
   specs, `docs/design-decisions/core-beliefs.md`, and
   `docs/design-decisions/README.md`. Discuss meaningful alternatives with the
   human before recording a new decision.
2. Create the next numbered `docs/exec-plan/todo/<sequence>-<name>.md`. Use the
   same `<name>` as the branch. A split parent keeps a shared base name until
   its scope is fully migrated; completed plans are removed and remain
   retrievable through their PR and Git history.
3. Directly below the H1, add:
   `> **Execution**: Use \`/execute-task\` to implement this plan. After implementation is complete, use \`/review-task\` to prepare and create the PR.`
4. Include: user-visible objective/completion boundary; exact references with
   paths, symbols, and ranges; `(NEW)/(MODIFY)/(DELETE)` map; black-box spec
   changes; `Addresses:` for relevant local or full external issue URLs;
   concrete sub-tasks/dependencies/parallelism; and verification. Use
   `N/A - detail required before execution` only for an intentional parent.
5. Send the plan PR through `review-task`. Put any external `Addresses:` issue
   under Issues in the plan PR body.

Do not implement the plan in this branch. After merge, start a matching
`feat/<name>` or `fix/<name>` branch with `execute-task`.
