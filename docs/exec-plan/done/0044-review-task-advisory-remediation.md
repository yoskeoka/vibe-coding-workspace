# Review-Task Advisory Remediation

> **Execution**: Use `/execute-task` to implement this plan. After implementation is complete, use `/review-task` to prepare and create the PR.

## Objective

Update the workflow contract so PR follow-up no longer stops at advisory bot/agent triage when the implementer's judgment is "fix in this PR" and the change is still reasonably in scope. After the workflow change lands, `review-task` should drive the branch to a human-review-ready state by default: CI green, substantive in-scope bot findings fixed, truly large follow-up work deferred into a plan or issue, and a compact separate-session handoff prompt included when asking the human to review.

This plan also changes the advisory-review waiting cadence so no single idle wait exceeds 3 minutes while preserving the current total budget, which keeps prompt-cache reuse viable during the automated landing run.

This supports the workspace goals of keeping hobby token cost low, making the workflow itself a first-class product, and reducing human overhead before the branch is actually ready for review.

## Existing Implementation References

- `AI_WORKFLOW.md`
  - `PR Workflow`, lines 93-112
  - `Execution Plan`, lines 122-140
  - `Execution`, lines 142-160
- `AGENTS.md`
  - `Branch & PR Rules`, lines 42-59
  - `Landing the Plane (Session Completion)`, lines 162-187
- `docs/specs/pr-follow-up-workflow.md`
  - `Ownership`, lines 7-20
  - `Monitoring Loop`, lines 21-67
  - `Advisory Bot/Agent Review Policy`, lines 164-181
  - `Stop Conditions`, lines 183-198
- `skills/review-task/SKILL.md`
  - `Post-PR Follow-up Loop`, lines 173-216
  - `Advisory Bot/Agent Review Triage`, lines 230-249
  - `Stop Conditions`, lines 251-264
- `skills/plan-execution/SKILL.md`
  - `PR Workflow`, lines 174-183
- `skills/execute-task/SKILL.md`
  - `PR Workflow`, lines 116-129
- `docs/exec-plan/done/review-task-post-pr-follow-up.md`
  - `Objective` through `Expected Outcome`, lines 1-110
- `docs/exec-plan/done/bot-agent-review-triage.md`
  - `Objective` through `Expected Outcome`, lines 1-110
- `docs/exec-plan/done/review-task-landing-completion.md`
  - `Objective` through `Expected Outcome`, lines 1-95
- `docs/exec-plan/todo/0043-gh-pr-followup-output-tiers.md`
  - `Objective` through `Expected Outcome`, lines 1-123

## Code Change Map

- `docs/specs/pr-follow-up-workflow.md` (MODIFY)
  - `Monitoring Loop`, `Advisory Bot/Agent Review Policy`, `Stop Conditions` - redefine the default handling for advisory findings so implementer-approved in-scope fixes are applied in the current PR, truly large follow-ups become a plan or issue, the advisory wait cadence is split into <=3-minute chunks without shrinking total wait time, and the human-review handoff requires a short separate-session prompt.
- `skills/review-task/SKILL.md` (MODIFY)
  - `Post-PR Follow-up Loop`, `Advisory Bot/Agent Review Triage`, `Stop Conditions` - mirror the new contract, including "fix unless no-action or genuinely separate large work", explicit defer routing to `plan-execution` or `docs/issues/`, the split wait cadence, and the final 3-line handoff prompt when the PR is ready for human review.
- `AI_WORKFLOW.md` (MODIFY)
  - `PR Workflow`, `Post-PR Follow-up`, `Execution` - align the top-level workflow with the new review-task behavior and the prompt-cache-aware waiting cadence.
- `AGENTS.md` (MODIFY)
  - `Branch & PR Rules`, `Landing the Plane` - update the top-level operating instructions so agents keep running through in-scope advisory fixes before asking for human review and include the separate-session handoff prompt.
- `skills/plan-execution/SKILL.md` (MODIFY)
  - `PR Workflow` - replace the older "CI/Copilot follow-up" wording with the updated advisory-remediation completion boundary so plan PRs and execution PRs describe the same contract.
- `skills/execute-task/SKILL.md` (MODIFY)
  - `PR Workflow` - align Step 3 execution wording with the new remediation-by-default policy and final human-review handoff requirement.
- `docs/lessons.md` (MODIFY)
  - in the top of this file, add one instruction to suggest human to cleanup lessons when the count of lessons is more than 10.
  - also in the top of this file, add the annotation for agent that "don't add a lesson already issue tracked/planned which will be fixed in the future".

## Spec Changes

- `docs/specs/pr-follow-up-workflow.md`
  - Document that advisory bot/agent findings are no longer only a handoff summary surface: when the implementer's view is not `no action` and the fix remains reasonably scoped to the current PR, `review-task` should implement the fix before handing off.
  - Document the stricter defer boundary:
    - truly large, clearly separate work with a known direction becomes a follow-up `plan-execution` task
    - truly large work without a settled solution becomes a `docs/issues/` entry
  - Replace any single 5-minute advisory wait with a sequence of waits no longer than 3 minutes each while preserving the same total advisory wait budget.
  - Require the human-review handoff to end with a compact separate-session prompt template that asks future human-driven fix requests for this PR to continue in a new session.

## Design Decisions

Past decisions reviewed before planning:

- `docs/project-plan.md` says the workflow is itself a product and should keep token cost and human time low.
- `docs/design-decisions/core-beliefs.md` favors human-review efficiency and trimming tool output/cost at the source.
- `docs/exec-plan/done/review-task-post-pr-follow-up.md` made `review-task` the owner of PR follow-up rather than leaving the landing loop fragmented.
- `docs/exec-plan/done/bot-agent-review-triage.md` deliberately kept advisory bot feedback behind a human approval boundary to avoid silent scope creep.

Trade-offs considered:

1. Keep the current explicit human approval gate for every advisory mutation.
   - Safer against unintended scope expansion, but it leaves many clearly in-scope PR cleanups unfinished and pushes avoidable review load onto the human.
2. Shift to a default "fix in this PR unless no-action or truly separate large work" policy.
   - Better matches the user's desired landing shape and keeps human review focused on a PR that is already mechanically and advisory-clean enough to inspect.
3. Auto-fix every advisory suggestion without an implementer judgment filter.
   - Fastest in theory, but too broad: it weakens scope control and would turn questionable bot comments into silent branch mutations.

Recommended option: option 2. The user has already specified this direction, so execution should treat that choice as confirmed rather than reopening it.

No ADR update is expected unless execution reveals a broader change to workflow ownership beyond `review-task`.

## Sub-tasks

- [ ] [parallel] Audit the current advisory-remediation wording in `docs/specs/pr-follow-up-workflow.md`, `skills/review-task/SKILL.md`, `AI_WORKFLOW.md`, `AGENTS.md`, `skills/plan-execution/SKILL.md`, and `skills/execute-task/SKILL.md`.
- [ ] [parallel] Define the exact remediation boundary for advisory findings:
  - fix in the current PR when the implementer's view is not `no action` and the work is still reasonably scoped
  - create a follow-up exec plan when the larger change direction is known
  - create a `docs/issues/` item when the larger solution is still unsettled
- [ ] [parallel] Define the replacement advisory wait cadence so each individual wait stays within 3 minutes while preserving the total budget and keeping the final review-start timeout semantics clear.
- [ ] [depends on: remediation boundary, wait cadence] Update `docs/specs/pr-follow-up-workflow.md` with the new external contract and the required separate-session handoff prompt.
- [ ] [depends on: spec update] Update `skills/review-task/SKILL.md` so the operational flow matches the spec, including the default-fix policy, defer routing, split wait cadence, and final handoff prompt.
- [ ] [depends on: spec update] Update `AI_WORKFLOW.md`, `AGENTS.md`, `skills/plan-execution/SKILL.md`, and `skills/execute-task/SKILL.md` so all caller and top-level docs describe the same completion boundary.
- [ ] [depends on: workflow and skill updates] Append the durable lesson to `docs/lessons.md`.
- [ ] [depends on: all above] Run workflow verification and confirm no conflicting older wording remains.

## Parallelism

- The wording audit, remediation-boundary definition, and wait-cadence design can start independently.
- The spec update depends on the agreed remediation and cadence rules.
- Skill and top-level workflow updates depend on the spec wording.
- The lesson entry and final verification depend on the completed wording changes.

## Verification

- Run `./tools/workflow-lint.sh --mode=pre-push`.
- Search for remaining wording that still says advisory fixes always require a fresh human decision before branch mutation, and confirm any remaining instances are intentional.
- Search for remaining single-wait `5 minutes` wording in the updated workflow/spec/skill surfaces and confirm it has been replaced with sub-3-minute waits while preserving the same total budget.
- Confirm `docs/specs/pr-follow-up-workflow.md`, `skills/review-task/SKILL.md`, `AI_WORKFLOW.md`, `AGENTS.md`, `skills/plan-execution/SKILL.md`, and `skills/execute-task/SKILL.md` all agree on:
  - in-scope advisory findings are fixed in the current PR by default unless the implementer's view is `no action`
  - truly large follow-ups become either a new exec plan or a `docs/issues/` item depending on solution clarity
  - the advisory wait cadence contains no single wait longer than 3 minutes
  - the final human-review request includes a compact separate-session handoff prompt
- Manually inspect the final handoff wording to ensure the new-session prompt stays within about 3 lines and is specific enough to reuse.

## Expected Outcome

- `review-task` drives PRs past advisory triage and through reasonable in-PR cleanup before asking for human review.
- Only truly large follow-up work is deferred, with a clear split between "known direction -> new plan" and "solution unclear -> issue".
- Long idle waits no longer cross the 5-minute prompt-cache boundary in one chunk, while total advisory wait time remains unchanged.
- Human review requests end with a short reusable prompt that steers later human-driven follow-up for the same PR into a fresh session.
