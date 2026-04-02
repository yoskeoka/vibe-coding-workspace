# Workflow Linter Warning Discipline

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Strengthen adherence to workflow-linter findings without turning the linter into a hard blocker. The outcome should make fixable warnings harder to ignore, preserve room to skip on explicit human instruction or clear false positives, and document that behavior in the linter spec, agent instructions, and PR review flow.

## Background

The current workflow linter emits warnings only and already includes ad hoc `WHY` / `FIX` guidance, but it does not distinguish between warnings that should normally be resolved before push/PR and warnings that are merely advisory. In practice this makes it too easy for agents to leave behind warnings that were mechanically fixable.

We want a smaller intervention than upgrading warnings to errors:

- classify warnings into `fixable` and `advisory`
- make fixable warnings carry explicit remediation hints
- require agents to resolve fixable warnings unless the human's instruction conflicts or the warning is a false positive
- reflect that expectation in the PR checklist

## Spec Changes

### `docs/specs/workflow-linter.md`

- Define warning classes:
  - `fixable`: the repo state can usually be corrected before push/PR by renaming, moving files, or other straightforward changes
  - `advisory`: useful signal, but not something that should automatically trigger repo mutation
- Update the linter behavior section to require normalized output that includes:
  - warning class
  - rationale (`WHY`)
  - remediation hint (`FIX`) for fixable warnings
  - summary counts by class
- Document the operating rule:
  - fixable warnings should normally be resolved before push/PR
  - they may be skipped only when user instruction takes precedence or the warning is judged to be a false positive
  - skipped fixable warnings must be justified in the PR
- Keep exit behavior non-blocking (`exit 0`)

### `AGENTS.md`

- Add a project-level rule that agents must not ignore workflow-linter findings
- State explicitly that fixable warnings are expected to be resolved before push/PR
- Allow skipping only for:
  - user instruction that conflicts with the warning
  - clear false positives
- Require the skip reason to be recorded in the PR body

### `.github/PULL_REQUEST_TEMPLATE.md`

- Add a checklist item confirming that workflow-linter warnings were reviewed and all fixable warnings were either resolved or explicitly justified

## Code Changes

### `tools/workflow-lint.sh`

- Introduce warning classes in the shell implementation without changing the exit code contract
- Refactor warning emission so each check can emit a structured class (`fixable` or `advisory`) together with message, `WHY`, and optional `FIX`
- Add a summary block at the end showing counts by warning class
- Add a final reminder when fixable warnings remain, instructing the user/agent to resolve them before push/PR unless overridden by human instruction or false-positive judgment
- Preserve compatibility with the existing pre-push and CI entry points

## Design Decisions

- Do not add helper commands or auto-fixers in this plan
- Do not convert warnings into errors in this plan
- Do not add new workflow-linter checks in this plan unless needed to support the warning classification system

## Sub-tasks

- [ ] Update `docs/specs/workflow-linter.md` to define warning classes, normalized output, and skip rules
- [ ] [parallel] Update `AGENTS.md` with the stronger workflow-linter handling rule
- [ ] [parallel] Update `.github/PULL_REQUEST_TEMPLATE.md` with a checklist item for fixable-warning resolution/justification
- [ ] [depends on: spec update] Refactor `tools/workflow-lint.sh` to emit warning classes, `WHY`/`FIX`, and summary counts
- [ ] [depends on: linter update] Manually verify representative outputs for both fixable and advisory warnings

## Verification

- Run `tools/workflow-lint.sh --mode=pre-push` in a repo state that produces at least one fixable warning and confirm:
  - warning class is shown
  - `WHY` is shown
  - `FIX` is shown
  - summary reports at least one fixable warning
- Run `tools/workflow-lint.sh --mode=ci --pr-title="..." --pr-body="..."` in a repo state that produces an advisory warning and confirm:
  - advisory class is shown
  - summary count reflects advisory warnings
- Confirm exit code remains `0`
- Review the PR template and `AGENTS.md` diff to ensure the expected discipline is documented

## Expected Outcome

- Workflow-linter output distinguishes between fixable and advisory warnings
- Agents have a documented rule to resolve fixable warnings before push/PR unless an allowed exception applies
- PR authors are prompted to confirm that fixable warnings were resolved or justified
- The linter remains non-blocking, reducing the risk of false-positive frustration
