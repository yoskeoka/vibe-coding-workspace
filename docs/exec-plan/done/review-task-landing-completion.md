# Review-Task Landing Completion

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Strengthen the workflow contract so `plan-execution` and `execute-task` do not stop at local implementation or pre-commit narration when the user expects full workflow completion. After either skill hands work to `review-task`, the documented completion point should remain `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll and inspection for the latest PR head SHA`, with any blocker or timeout explicitly reported.

## Current State

- `AI_WORKFLOW.md`, `AGENTS.md`, and `review-task` already describe bounded post-PR follow-up ownership.
- `plan-execution` and `execute-task` both mention routing into `review-task`, but they do not strongly state that the caller must not report completion before the shared PR gate finishes the initial follow-up loop.
- The repo already records several PR follow-up lessons, but there is not yet a workflow-facing rule aimed at the failure mode where the agent narrates progress, updates `docs/lessons.md`, and then stops before `commit -> push -> PR`.
- This leaves room for an implementation to satisfy the local coding portion while still violating the user's expected landing steps.

## Spec Changes

### `AI_WORKFLOW.md`

- Make the ownership boundary more explicit: when a workflow step routes into `review-task`, the step is not complete until the initial post-PR follow-up stop condition is reached for the latest pushed head SHA.
- Clarify that this applies equally to planning PRs and execution PRs.

### `docs/specs/pr-follow-up-workflow.md`

- Add an explicit handoff rule that caller skills must not report task completion before `review-task` reaches a documented stop condition for the current PR head SHA.
- Name the minimum landing path expected for non-blocked PR creation flows as one initial `gh-pr-followup` poll after `commit -> push -> PR create/update -> 30-second wait`, including inspection of checks, timeline events, review summaries, and inline comments.

### `skills/plan-execution/SKILL.md`

- Strengthen the PR workflow section so the planning step is only considered complete after `review-task` finishes its initial PR creation/update and follow-up loop.
- Make the "do not stop at local plan file creation" expectation explicit.

### `skills/execute-task/SKILL.md`

- Strengthen the PR workflow section so execution completion is only reported after `review-task` finishes its initial PR creation/update and follow-up loop.
- Make the "do not stop at local implementation, commit drafting, or pre-push narration" expectation explicit.

### `docs/lessons.md`

- Add a durable lesson that execution-oriented requests requiring branch lifecycle completion must not stop before `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll`, even if the agent already updated lessons or finished local verification.

## Code Changes

- Update the relevant workflow docs and skill instructions to align on the same completion boundary.
- Keep the change documentation-only unless the wording audit reveals a helper/script mismatch that materially contradicts the intended PR follow-up behavior.

## Design Decisions

Past decisions:
- `review-task` is already the shared PR gate, so strengthening caller-skill completion semantics is preferable to adding another landing or handoff skill.
- `core-beliefs.md` favors correctness and compact, review-ready artifacts, which supports a strict completion boundary around actual PR state rather than narration about intended next steps.

Apply the same reasoning here:
- preserve `review-task` as the single owner of PR creation/update and initial follow-up
- make caller skills explicitly subordinate to that owner for completion reporting
- prefer minimal wording changes over new workflow machinery

No ADR update is expected.

## Sub-tasks

- [ ] [parallel] Audit `AI_WORKFLOW.md`, `docs/specs/pr-follow-up-workflow.md`, `skills/plan-execution/SKILL.md`, and `skills/execute-task/SKILL.md` for wording that still permits early stop points
- [ ] [parallel] Capture the new durable lesson in `docs/lessons.md`
- [ ] [depends on: wording audit] Update the workflow and spec docs with the explicit completion boundary
- [ ] [depends on: wording audit] Update `plan-execution` and `execute-task` so their `review-task` handoff language forbids early completion reporting
- [ ] [depends on: all above] Verify the final wording is consistent across top-level workflow docs, specs, and caller skills

## Parallelism

- The wording audit and lesson capture can start independently.
- Doc and skill edits depend on the wording audit.
- Final verification depends on all wording updates landing.

## Verification

- Confirm `AI_WORKFLOW.md`, `docs/specs/pr-follow-up-workflow.md`, `skills/plan-execution/SKILL.md`, and `skills/execute-task/SKILL.md` all describe the same completion boundary.
- Confirm the workflow now says callers must not report completion before `review-task` reaches a documented stop condition for the latest PR head SHA.
- Confirm the minimum non-blocked landing path is named explicitly as `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll`, with checks, timeline, reviews, and inline comments inspected from that poll.
- Confirm the lesson entry states the same rule in operational language.

## Expected Outcome

- Planning and execution workflow steps no longer leave room to treat local implementation or local plan creation as the end of the task once `review-task` is required.
- The repo documents one explicit, shared answer to "when is this task complete enough to report back?" for both planning and execution PRs.
