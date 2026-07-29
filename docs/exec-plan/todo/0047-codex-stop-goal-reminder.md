# Codex Stop Goal Reminder

> **Execution**: Use `/execute-task` to implement this plan. After implementation is complete, use `/review-task` to prepare and create the PR.

## Objective and completion boundary

Prevent a Codex agent from treating a progress update as completion while it is carrying out a workspace workflow task. On the first `Stop` event of a turn in a workspace `plan/*`, `feat/*`, or `fix/*` worktree, Codex must receive one concise continuation prompt that asks it to check the user-requested completion boundary. When an active matching execution plan is available, that prompt must name its objective or completion boundary. The second `Stop` event for that same turn must be allowed so a genuinely complete task can hand off normally.

Completion is a workspace-local configuration, script, specification, and focused-test change that is ready for a reviewable PR. It does not change child-repository hooks or make a hook an unconditional blocker.

Addresses: N/A

## Current evidence

- `.codex/config.toml:10-12` enables trusted project-scoped Codex hooks for the workspace.
- `.codex/hooks.json:16-27` already registers a `Stop` command that dispatches child-repository quality gates through `dispatch_hook.py`; Codex launches matching command hooks concurrently, so the new reminder must not depend on that dispatcher running first.
- `.codex/hooks/dispatch_hook.py:13-76` discovers and forwards `post-tool-use` and `stop` payloads to `ai-arena`, including when the session cwd is below a child repository.
- `skills/execute-task/SKILL.md:11-25` and `docs/specs/pr-follow-up-workflow.md:7-18,212-227` make review-task's latest-head stop condition the current execution completion boundary.
- The current Codex Hooks manual documents `Stop` input fields `stop_hook_active` and `last_assistant_message`, and documents that `decision: "block"` continues the turn with the hook reason as a new user prompt. It also documents that the transcript format is not a stable hook interface, so this change must not depend on transcript parsing.

## Change map

### Observable contract and operator guidance

- (NEW) `docs/specs/codex-stop-goal-reminder.md`: define workspace-only applicability, eligible workflow branches, one-continuation-per-turn behavior, active-plan summary behavior, non-blocking exits, and the expected reminder content for plan versus execution work.
- (MODIFY) `docs/specs/README.md`: list the new Codex stop-goal-reminder contract.
- (MODIFY) `AI_WORKFLOW.md`: identify the workspace-local Stop reminder as a guardrail that asks for completion-boundary confirmation but does not replace the documented `review-task` ownership or user-directed stops.

### Codex hook implementation

- (MODIFY) `.codex/hooks.json`: register the reminder as a second `Stop` command with a short timeout and a status message. Preserve the existing child quality-gate dispatcher unchanged; do not assume handler ordering.
- (NEW) `.codex/hooks/stop_goal_reminder.py`: read the JSON payload from stdin; prove that the session cwd belongs to this workspace worktree rather than a child repository; inspect the current branch and matching active plan without reading the unstable transcript; and emit valid Stop-hook JSON. On an eligible first Stop, emit `decision: "block"` with a concise goal/completion reminder. On `stop_hook_active`, a non-workspace cwd, an ineligible branch, malformed/absent context, or an intentional non-workflow state, exit successfully without a continuation decision.
- (NEW) `.codex/hooks/test_stop_goal_reminder.py`: use standard-library temporary Git repositories or injected command seams to cover workspace versus child/worktree scope, branch eligibility, matching-plan extraction, missing-plan fallback text, first-stop continuation output, and second-stop pass-through behavior.

## Black-box contract changes

1. A trusted workspace Codex session running from the workspace root or one of its worktrees on `plan/*`, `feat/*`, or `fix/*` receives at most one automatic completion check per turn.
2. The completion check tells the agent to continue work unless it can affirm the user-requested goal is complete or identify a genuine blocker or required user decision. It reminds planning work that its reviewable-plan PR and initial latest-head follow-up are the boundary; it reminds execution work that verification, PR creation, and `review-task`'s latest-head stop condition are the boundary.
3. When a matching active plan exists under `docs/exec-plan/todo/`, the check includes a bounded summary from that plan instead of guessing from the session transcript.
4. A session rooted in a managed child repository, a non-workflow branch, or a second Stop event from the same turn receives no continuation from this hook. Existing Stop handlers continue to operate independently.

## Execution steps

1. Write the black-box contract first. State that this is a lightweight workflow guardrail, not a source of truth for CI or PR state, and align `AI_WORKFLOW.md` with the existing caller/`review-task` completion ownership.
2. Implement the focused Python Stop-hook command. Resolve the workspace Git top-level from the hook location and compare it with the session cwd's Git top-level so a child-repository session cannot trigger it. Read the current branch only after scope succeeds; accept only `plan/`, `feat/`, and `fix/` names.
3. Locate the active plan whose filename suffix matches the eligible branch description. Extract only a short, deterministic H1/objective-or-completion-boundary summary with a hard output limit. Never parse `transcript_path`, never serialize user content beyond the bounded plan text, and use a generic workflow reminder when no matching plan is available.
4. Make first Stop output a valid JSON `decision: "block"` and a concrete continuation reason. If `stop_hook_active` is true, emit no decision so the continued agent may complete its handoff. Treat malformed input and Git lookup failures as fail-open with a concise diagnostic only where Codex can safely surface it.
5. Add focused tests before registering the hook, then add the second independent handler to `.codex/hooks.json`. Verify the existing dispatcher still receives its payload and that a child cwd does not produce the new reminder.
6. Run the focused hook tests, configuration parsing/JSON validation, workspace workflow checks, and a manual stdin fixture for both first and second Stop payloads. Then use `review-task` for the plan PR and its Step 2 latest-head follow-up.

## Dependencies and parallelism

- Step 1 precedes implementation because the observable scope and stop semantics must be agreed before a hook can alter the agent loop.
- Steps 2 and 3 are one implementation unit because scope, branch mapping, and bounded plan extraction jointly determine whether continuation is safe.
- The test module can be written alongside Steps 2-4, but hook registration waits until the tested command behavior is stable.
- Existing child dispatch remains independent and requires a regression check after registration, not a code merge with this feature.

## Verification

- Run `python3 .codex/hooks/test_stop_goal_reminder.py` (or the focused standard-library test command selected by the implementation) and show coverage of all scope, branch, plan-summary, and loop-prevention cases.
- Pipe representative JSON payloads into `.codex/hooks/stop_goal_reminder.py`; validate that the first eligible payload returns parseable JSON with `decision: "block"`, while the repeated Stop payload returns success without a decision.
- Run the existing dispatcher with an `ai-arena` cwd fixture and confirm its behavior is unchanged; run the reminder with the same child cwd and confirm it is silent.
- Run `python3 -m json.tool .codex/hooks.json`, `./tools/test-workflow-context-contract.sh`, `./tools/workflow-lint.sh --mode=pre-push`, and `git diff --check`.
- In the plan PR, complete the `review-task` Step 2 initial latest-head compact poll and state the reached stop condition.
