# Codex stop goal reminder

## Purpose

This workspace hook checks a workflow task before Codex stops. It asks the
agent to check the user's goal. It is not a source of truth for CI, PRs, or
review state.

## Scope and behavior

- The hook runs only when the payload `cwd` has this workspace Git root. It
  also runs from a `ww` worktree. It does not run in a child repository.
- It runs only on `plan/*`, `feat/*`, and `fix/*` branches.
- The first Stop event returns one `decision: "block"` and a short prompt. Any
  truthy `stop_hook_active` value returns no decision. This lets the second
  Stop event hand off as usual.
- Bad input, missing Git data, a non-workflow branch, or a missing cwd returns
  no decision. Other Stop hooks still run on their own.

## Reminder content

- Add a plan summary only when one file in `docs/exec-plan/todo` has a safe
  branch suffix.
- Use the generic prompt for an empty, unsafe, or matching multiple names.
- Read only a short heading and objective or completion boundary. Never read
  the transcript.

On `plan/*`, the prompt asks for the reviewable plan PR and its first
latest-head check. On `feat/*` and `fix/*`, it asks for verification, PR
creation, and `review-task`'s latest-head stop condition. In either case, the
agent continues unless the goal is complete, blocked, or needs a user choice.
