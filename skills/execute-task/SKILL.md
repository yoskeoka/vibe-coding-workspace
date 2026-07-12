---
name: execute-task
description: Implement an approved execution plan with spec-first ordering and PR handoff.
---

# Execute Task

Use only after the matching execution plan has merged. The shared lifecycle is
[execution](../../AI_WORKFLOW.md#execution) and
[PR follow-up](../../AI_WORKFLOW.md#pr-and-follow-up).

## Procedure

1. Create matching `feat/<name>` or `fix/<name>` worktree with `ww`; read the
   plan, relevant specs, and referenced implementation.
2. Update the applicable black-box spec first, then implement exactly that
   contract. Log non-blocking unrelated findings as a numbered local issue.
3. Run all applicable project quality gates and capture reviewable evidence.
4. Move the completed plan to `docs/exec-plan/done/`; move linked resolved local
   issues to `docs/issues/done/`. Preserve full external issue URLs for the PR's
   conditional closure metadata.
5. For significant work, run `post-task-review`; then invoke `review-task`.

Do not treat local edits or an opened PR as completion: `review-task` owns the
latest-head follow-up stop condition.
