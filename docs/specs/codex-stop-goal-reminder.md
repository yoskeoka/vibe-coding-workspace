# Codex Stop Goal Reminder

## Purpose

The workspace provides a lightweight Codex Stop-hook guardrail for workflow
work. It prompts an agent to check the user-requested completion boundary; it
is not a source of truth for CI, pull-request, or review state.

## Scope and behavior

- The guardrail applies only when the hook payload `cwd` resolves to the Git
  top-level of this workspace root or its current `ww` worktree. A managed
  child repository is out of scope.
- It applies only on `plan/*`, `feat/*`, and `fix/*` branches.
- On the first Stop event of a turn, it returns one `decision: "block"` with a
  concise continuation prompt. When `stop_hook_active` is true, it returns no
  decision so the second Stop event can hand off normally.
- A malformed payload, unavailable Git context, non-workflow branch, missing
  cwd, and other intentional non-workflow state fail open without a
  continuation. Other Stop handlers remain independent.

## Reminder content

When exactly one active plan in `docs/exec-plan/todo/` has a filename suffix
matching a non-empty, path-safe current branch description, the reminder
includes a bounded heading and objective or completion-boundary summary from
that plan. Ambiguous or unsafe branch descriptions use the generic reminder.
It never reads the unstable transcript interface.

For `plan/*`, the prompt asks the agent to continue unless it can affirm that
the reviewable plan PR and its initial latest-head follow-up are complete, or
it identifies a genuine blocker or required user decision. For `feat/*` and
`fix/*`, it instead names verification, PR creation, and `review-task`'s
latest-head stop condition. If no matching plan is available, the prompt uses
the same branch-appropriate generic completion check.
