# Review-Task Post-PR Follow-up

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Extend the workflow so `review-task` does not stop at PR creation. After opening or updating a PR, the agent should wait for CI completion, detect whether GitHub Copilot code review is active, automatically pursue CI-failure fixes, and produce a human-review briefing for Copilot comments instead of auto-applying them.

This captures the actual operator workflow already used after PR creation:

- wait for CI
- wait briefly for Copilot review when the repo/PR is configured for it
- auto-fix mechanical CI failures when possible
- read Copilot comments and summarize whether to ignore, apply as-is, adapt, or defer to a separate issue

## Current State

- `skills/review-task/SKILL.md` currently stops after the PR is created or updated and ready for review.
- The workflow does not define who owns post-PR waiting, CI retry/fix loops, or AI-review triage after the PR exists.
- In practice, this leaves repetitive manual follow-up work outside the documented skill contract.
- Copilot review needs a different policy boundary from CI:
  - CI failures should be investigated and fixed automatically when feasible.
  - Copilot comments should be read and analyzed, but not auto-resolved without explicit implementation work.

## Spec Changes

### `AI_WORKFLOW.md`

- Clarify that the PR workflow includes post-creation monitoring until initial CI and any configured automatic AI review have settled.
- State the policy split:
  - CI failures are mechanical verification failures and should be fixed in-branch when feasible before handing off.
  - Copilot review comments are advisory review input and require explicit analysis rather than automatic acceptance.

### `AGENTS.md`

- Add a concise session-end rule for PR follow-up:
  - after PR creation, wait for CI to finish
  - if Copilot auto-review appears, wait a short bounded interval for it to post comments
  - treat CI failures as auto-fix candidates
  - treat Copilot comments as human-review prep, not auto-merge criteria and not auto-fix instructions

### `skills/review-task/SKILL.md`

- Expand the skill from a pure PR-creation gate into a PR-preparation-and-initial-follow-up gate.
- Define a bounded post-PR monitoring loop:
  - create/update PR
  - prefer delegating polling to a low-cost subagent when available so the main agent does not burn expensive tokens on wait loops
  - wait for CI checks
  - if the PR timeline shows Copilot auto-review starting, wait for review completion/comments
  - stop waiting after a short bounded timeout when no Copilot review appears
- Define the CI failure loop:
  - inspect failing checks
  - attempt fixes automatically when the failure is actionable from logs and within scope
  - re-run verification and push updates
  - repeat until green or blocked
- Define the Copilot review triage output:
  - comment summary
  - whether action is recommended
  - options for response or implementation
  - whether GitHub suggestion can be applied as-is or needs adaptation
  - whether the comment is better deferred into a separate issue
- Explicitly forbid silently auto-applying Copilot review suggestions.

## Code Changes

### Workflow skill updates

- Update `skills/review-task/SKILL.md` to describe the new post-PR responsibilities, stop conditions, and policy boundaries.
- Prefer a low-cost small subagent for polling-style wait loops when the platform supports delegation; on GPT-family models this should bias toward a `mini`-class worker instead of the main model.
- Update any skill that routes into `review-task` if its wording assumes the skill ends at PR creation.

### Supporting workflow docs

- Update `AI_WORKFLOW.md` and `AGENTS.md` so the broader workflow matches the strengthened `review-task` behavior.
- Update or add a workflow-facing spec if needed so the post-PR contract is documented outside the skill body as well.

## Design Decisions

Past decisions:
- `review-task` is already the shared PR gate, so extending it is preferable to inventing a second overlapping PR-follow-up skill.
- `core-beliefs.md` favors correctness over speed, which supports waiting for real CI/review signals instead of treating PR creation as the end of the job.

Apply the same reasoning here:
- keep one owner skill for PR readiness and initial post-PR follow-up
- separate mechanical verification failures from advisory review feedback
- preserve human judgment for review comments even when automation is available
- use cheaper delegated polling for idle wait loops when available, while keeping final decisions and implementation in the main agent

No ADR update is expected unless execution reveals a broader policy for third-party AI review handling.

## Sub-tasks

- [ ] [parallel] Inventory the current `review-task` wording and identify exactly which steps stop at PR creation
- [ ] [parallel] Decide the minimum workflow contract for CI waiting, retry/fix loops, Copilot detection windows, and low-cost polling delegation
- [ ] [depends on: workflow contract] Update `skills/review-task/SKILL.md` with the post-PR monitoring and triage flow
- [ ] [depends on: workflow contract] Update `AI_WORKFLOW.md` and `AGENTS.md` so the top-level workflow matches the skill behavior
- [ ] [depends on: review-task update] Update any calling skills whose wording now understates `review-task` responsibilities
- [ ] [depends on: all above] Verify the resulting workflow clearly separates automatic CI remediation from manual Copilot-comment decision-making

## Parallelism

- Current-state inventory and workflow-contract definition can start independently.
- Skill/doc updates depend on agreeing on the post-PR contract.
- Caller-skill cleanup depends on the final `review-task` wording.

## Verification

- Confirm the workflow docs no longer imply that `review-task` ends at PR creation.
- Confirm the new contract distinguishes CI auto-fix behavior from Copilot comment analysis.
- Confirm the documented Copilot flow stops at analysis/options and does not authorize automatic suggestion application.
- Confirm the waiting behavior is bounded so the agent does not block indefinitely on absent Copilot review.
- Confirm polling guidance prefers a low-cost delegated worker when available instead of burning the main model on repeated status checks.

## Expected Outcome

- Agents treat PR creation as the start of initial review monitoring, not the terminal step.
- Mechanical CI failures are automatically driven toward resolution when feasible.
- Copilot review comments are collected and summarized into actionable human-review options rather than being auto-applied.
