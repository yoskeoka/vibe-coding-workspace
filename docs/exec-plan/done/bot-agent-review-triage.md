# Bot Agent Review Triage

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Update the PR follow-up workflow so advisory bot/agent review feedback is collected and triaged for the human before any comment-driven edits are made.

Addresses: `docs/issues/bot-agent-review-triage-should-not-auto-apply.md`

This supports the project-plan goals of treating the workflow as a first-class product, keeping review manageable for one human across many projects, and keeping agent behavior auditable instead of allowing advisory automation to create unreviewed scope creep.

## Current State

- `docs/specs/pr-follow-up-workflow.md` and `skills/review-task/SKILL.md` already separate CI failures from GitHub Copilot comments.
- Both documents still have Copilot-specific sections and stop conditions, even though the monitoring loop can detect Copilot, Claude, `gh aw`, agent workflow, or other bot/agent reviewer activity.
- The current wording says actionable bot/agent review comments should be resolved in scope or documented before handoff. That can be read as permission to silently apply advisory bot/agent comments during PR follow-up.
- `docs/issues/bot-agent-review-triage-should-not-auto-apply.md` records the desired policy: extract each advisory bot/agent finding, include the implementer's view and a short explanation, recommend whether to fix in the current PR, and wait for explicit human approval before editing or pushing comment-driven changes.
- Later pushes should always inspect normal CI for the latest head SHA, but should not spend the longer bot/agent review wait budget unless new review-start activity appears for that head SHA or the human asks to wait.

## Spec Changes

### `docs/specs/pr-follow-up-workflow.md`

- Rename the Copilot-specific policy to a general advisory bot/agent review policy.
- Define advisory bot/agent review comments as human-review input, not automatic patch instructions.
- Require the handoff to group findings by source reviewer/workflow and include:
  - source reviewer/workflow
  - comment location or link
  - extracted comment summary
  - implementer's view
  - concise 1-2 line explanation
  - recommendation: fix in this PR, defer, or no action
- Clarify that passing or approving advisory bot/agent checks can still contain substantive observations in review bodies, so review summaries need inspection even when the overall state is not blocking.
- Separate the two loops:
  - CI/check loop: always inspect required mechanical checks for every pushed PR head SHA.
  - Advisory review loop: wait for the bounded bot/agent cadence only on initial review startup, when new review-start activity appears for the latest head SHA, or when the human explicitly asks to wait.
- Update stop conditions so the workflow can stop after required checks pass and advisory bot/agent findings have been summarized, without requiring automatic edits.

## Code Changes

### `skills/review-task/SKILL.md`

- Rename `Copilot Review Triage` to advisory bot/agent review triage.
- Replace Copilot-only triage fields with the source-grouped triage fields from the spec.
- Explicitly forbid editing files, applying suggestions, committing, or pushing in response to advisory bot/agent findings unless the human approves or a prior human instruction already authorized that specific review-feedback implementation work.
- Clarify that `review-task` should summarize bot/agent findings in the current session and ask for a decision when the next action would mutate the branch.
- Update later-push guidance so required CI/checks are always inspected for the latest head SHA, while the longer advisory wait budget is skipped unless new bot/agent review-start activity is observed for that SHA or the human asks to wait.
- Update stop conditions so they refer to advisory bot/agent review comments rather than Copilot-only comments.

### `AI_WORKFLOW.md`

- Align the top-level PR follow-up wording with the generalized advisory bot/agent review policy.
- Clarify that CI remediation remains automatic when actionable and in scope, but advisory review remediation requires human approval before mutation.
- Clarify later-push behavior: always check CI for the new head SHA; only repeat the longer advisory wait when new review-start activity appears or the human asks.

### `AGENTS.md`

- Update PR follow-up and session completion instructions so agents triage Copilot, `gh aw`, agent workflow, and other advisory bot/agent review findings before implementation.
- Include the required compact triage output fields so future agents have the policy in their top-level operating instructions.
- Clarify that advisory bot/agent checks are reported by status but do not block handoff by default on second and later pushes unless required CI fails, review-start activity appears for the latest SHA, or the human requests more waiting.

### `docs/issues/`

- Move `docs/issues/bot-agent-review-triage-should-not-auto-apply.md` to `docs/issues/done/bot-agent-review-triage-should-not-auto-apply.md` during execution after the workflow/spec/skill updates are complete.

## Design Decisions

Past decisions reviewed before planning:

- `docs/design-decisions/core-beliefs.md` favors AI-first context retrieval, correctness over speed, human review over token burn, and trimming tool output at the source.
- `docs/design-decisions/adr.md` records that normal workflow startup should dogfood the released global `ww` binary.
- `docs/exec-plan/done/review-task-post-pr-follow-up.md` established `review-task` as the owner of PR readiness and initial post-PR monitoring.
- `docs/exec-plan/done/review-task-pr-workflow-guard.md` strengthened `review-task` as the shared PR gate instead of adding overlapping PR workflow skills.
- `docs/exec-plan/done/gh-pr-followup-token-trimming.md` added the compact PR follow-up helper so repeated polling does not flood the main context.

Trade-offs considered:

- Keep the policy in `review-task` and `pr-follow-up-workflow`: This preserves the existing owner for PR follow-up and keeps the workflow easy for agents to find. Recommended.
- Create a separate bot-review triage skill: This would isolate the policy, but it would split one PR handoff flow across two skills and increase the chance agents skip the boundary.
- Treat all bot/agent findings as normal review feedback to fix immediately: This maximizes automation, but it conflicts with human-review-over-token-burn and with the issue's goal of preventing advisory feedback from becoming unreviewed branch mutation.

Recommendation: update the existing PR follow-up contract in place. No ADR update is expected because this refines an established PR-review policy rather than introducing a new architectural boundary.

## Sub-tasks

- [x] [parallel] Inventory current Copilot-specific wording in `docs/specs/pr-follow-up-workflow.md`, `skills/review-task/SKILL.md`, `AI_WORKFLOW.md`, and `AGENTS.md`.
- [x] [parallel] Define the exact advisory bot/agent triage output shape and later-push wait rules from the issue.
- [x] [depends on: advisory triage output shape] Update `docs/specs/pr-follow-up-workflow.md` with the generalized policy, source-grouped handoff requirements, and separated CI/advisory loops.
- [x] [depends on: spec update] Update `skills/review-task/SKILL.md` so the skill follows the spec and does not authorize silent advisory-comment edits.
- [x] [depends on: spec update] Update `AI_WORKFLOW.md` and `AGENTS.md` with the same policy boundary and later-push behavior.
- [x] [depends on: workflow docs and skill updates] Move the source issue to `docs/issues/done/`.
- [x] [depends on: all above] Run workflow lint and manually verify the policy is consistent across spec, skill, and top-level agent instructions.

## Parallelism

- The wording inventory and triage-output design can happen independently.
- Spec changes should land before skill and top-level instruction updates so the skill mirrors the spec.
- Issue resolution and final verification depend on the documentation updates.

## Verification

- Run `tools/workflow-lint.sh --mode=pre-push`.
- Search for remaining Copilot-only policy wording and confirm either:
  - it was generalized to advisory bot/agent review triage, or
  - it intentionally refers to a Copilot-specific event such as `copilot_work_started`.
- Confirm `docs/specs/pr-follow-up-workflow.md`, `skills/review-task/SKILL.md`, `AI_WORKFLOW.md`, and `AGENTS.md` agree on:
  - CI failures can be fixed automatically when actionable and in scope.
  - Advisory bot/agent comments are summarized for human review before mutation.
  - Each substantive advisory finding gets source, location/link, summary, implementer's view, explanation, and fix/defer/no-action recommendation.
  - Later pushes always inspect CI for the new head SHA.
  - Later pushes only spend the longer advisory wait when new review-start activity appears for that SHA or the human asks to wait.
- Confirm the issue file is moved to `docs/issues/done/` only after execution completes the documented updates.

## Expected Outcome

- Agents stop treating advisory bot/agent review comments as automatic patch instructions.
- Human handoff includes a compact, source-grouped review briefing with the implementer's judgment for each substantive advisory finding.
- CI remains an automatic remediation loop, while advisory review remains a human-controlled decision point.
- Later PR pushes do not repeatedly burn the longer bot/agent review wait budget unless there is new evidence that review automation started for the latest head SHA.
