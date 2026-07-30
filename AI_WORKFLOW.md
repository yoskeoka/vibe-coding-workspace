# AI-Centered Development Workflow

This is the canonical lifecycle reference. Read the section for the current
phase; task-specific skills own their procedures. See
[the context contract](docs/specs/workflow-context-contract.md) for document
ownership and minimum reads.

## Shared rules

- Every change uses a clean branch from current `main` and GitHub PR review.
  Start normally with `ww create <type>/<name>` and `cd "$(ww cd <type>/<name>)"`.
  From the workspace root for a child repo, use `ww create --repo <repo> ...`.
- Branches are `<type>/<description>` using `plan`, `feat`, `fix`, `chore`, or
  `docs`; descriptions are kebab-case. `feat` and `fix` names map to a matching
  numbered plan suffix. Active plans and local issues are `<sequence>-<name>.md`.
- Raw Git branch creation is only for documented `ww` bootstrap/recovery or
  unreleased `ww` work. Record a normal-workflow `ww` failure with command, cwd,
  target, expected/actual result, fallback, and impact.
- `docs/specs/` describes observable product and operated-service behavior, not
  local harness or CI mechanics. Update it before code.
- Keep unrelated findings in `docs/issues/`; delete resolved local issues with
  their completed execution plan after verification and PR preparation. Retrieve
  completed tracker details from the implementation PR or Git history. Read the
  [design-decision index](docs/design-decisions/README.md) and core beliefs before
  making a design decision.

## Project planning

Update `docs/project-plan.md` when goals, significance, or requirements change.
Use a `plan/<name>` worktree and `plan-project`; submit the project-plan PR
through [PR and follow-up](#pr-and-follow-up).

## Execution planning

Use `plan/<name>` and create a numbered plan under `docs/exec-plan/todo/`.
Each executable plan has an external completion boundary, concrete existing
references (paths, symbols, ranges), a `(NEW)/(MODIFY)/(DELETE)` change map,
black-box spec changes, applicable `Addresses:`, dependencies/parallelism, and
the `/execute-task` instruction directly below its title. A high-level parent
uses `N/A - detail required before execution` and is split before execution.

Plan PRs use [PR and follow-up](#pr-and-follow-up). External issues listed in a
plan are also listed under the plan PR's Issues section.

## Execution

Use `feat/<name>` or `fix/<name>` only after the matching plan is merged. Follow
this order: spec update, implementation, scoped issue logging, verification,
then delete the completed plan. Delete local issues named on `Addresses:` in the
same branch; retrieve them later through the PR or Git history. For external
issues, the implementation PR includes `Closes #<n>` (same repo) or `Closes
<URL>`, unless it explains why the issue stays open. Invoke `post-task-review`
for significant work, then
[PR and follow-up](#pr-and-follow-up).

The workspace-local Codex Stop reminder is a guardrail that asks an agent to
confirm this completion boundary. It does not replace this caller's workflow
responsibility, a user-directed stop, or `review-task` ownership of the
latest-head follow-up.

## PR and follow-up

Run applicable lint, tests, builds, and required visual/manual checks before PR
creation. `review-task` verifies scope, title, applicable template fields, and
the latest-head landing loop. Before every push to a PR branch, confirm the PR
is still open. After create/update: wait 30 seconds, inspect checks, timeline,
review summaries, and inline comments. For execution PRs with pending required
checks, make up to two additional 30-second polls. If advisory review starts,
poll at 3m, 2m, 1m, and 1m; inspect every review body and inline comment.

Fix actionable, in-scope CI failures and advisory findings in the PR; defer a
clearly separate larger item to a plan or issue. The handoff groups advisory
findings by source and gives the PR URL, worktree path, and a short
fresh-session prompt. Wait for human approval before merging.

## Post-task review

For significant completion, use `post-task-review` to capture unrecorded intent,
surface concrete follow-ups for human approval, and maintain active exceptions.
Do not create speculative issues or generic lessons.

## Workflow setup

Use `manage-workflow` to bootstrap a repository. It creates the docs structure,
keeps `AGENTS.md` canonical with `CLAUDE.md` as a symlink, and installs the
workflow assets. Do not overwrite project-specific guidance.

## Task triage

Use `triage-tasks` for workspace priorities. GitHub Projects is canonical;
`.local/pj/` is derived cache only. Execution plans and local issues remain the
implementation trackers after task selection.
