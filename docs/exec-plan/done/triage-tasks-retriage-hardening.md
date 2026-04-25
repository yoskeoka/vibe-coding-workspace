# Triage Tasks Re-Triage Hardening

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Harden the `triage-tasks` full re-triage workflow so it updates the GitHub
Project accurately, cheaply, and predictably when repo state has drifted.

This plan captures the improvements identified during the latest full
re-triage: cheaper read-only exploration, update-before-add reconciliation,
stale source handling, stricter subagent output, clearer mutation planning, and
a consistent final briefing.

## Background

The current skill can refresh the Project cache and collect repo-by-repo state,
but the recent full re-triage exposed several workflow gaps:

- read-only repo exploration may use a larger default model unless the operator
  explicitly chooses a small one
- existing board items can point at stale PR numbers or local files that have
  moved to `done/`
- the population step says to create Project items, which biases agents toward
  duplicate `pj add` calls instead of updating existing items
- old item bodies may remain as `Source:` only, even after newer body templates
  exist
- project-plan phase gaps, local issues, exec plans, PRs, and GitHub Issues can
  be mixed into one undifferentiated suggested-item list
- missing local checkouts are not covered clearly
- workflow-sync PR items need repo-local replacement behavior because new sync
  PRs supersede old ones
- many `pj update` commands are slow unless the agent first builds a mutation
  plan and applies updates deliberately
- subagent summaries vary in schema, repo naming, priority casing, and `next`
  field shape
- the final briefing should consistently report what changed and what remains
  actionable

Past decisions and constraints:

- GitHub Projects is the canonical workspace triage board.
- `.local/pj/cache.json` is derived data, not source of truth.
- Full re-triage is expensive; day-to-day triage should prefer `pj list` and
  cached `Priority`.
- Human review over token burn means the workflow should use compact, structured
  artifacts and low-cost models for routine read-only collection.

## Trade-offs

### Option A: Keep full re-triage as loose guidance

- Flexible for unusual sessions.
- Leaves repeated errors to operator memory: duplicate adds, stale PR items, and
  model choice mistakes remain likely.

### Option B: Add a strict reconciliation and mutation workflow to the skill

- Gives agents a repeatable order: collect, reconcile, plan mutations, apply,
  sync, brief.
- Reduces duplicate Project items and stale board state.
- Slightly lengthens the skill, but keeps the complexity where the workflow
  decision actually happens.
- Recommended for this plan.

### Option C: Move most re-triage logic into `tools/pj`

- Could eventually make reconciliation more mechanical.
- Requires product design for source identity, stale detection, and task
  categories inside the CLI.
- Too large for this documentation/skill hardening slice; the current plan
  should improve agent behavior first.

## Spec Changes

### `docs/specs/triage-tasks.md`

Update the `triage-tasks` integration contract with these requirements:

- Read-only repo exploration should explicitly choose the cheapest currently
  available small model when delegation is used, unless the repo question
  requires deeper reasoning. Keep the contract model-name agnostic.
- Full re-triage must reconcile existing Project items before adding new ones:
  - compare by `Source` URL/path
  - compare by normalized title when source is missing or old
  - prefer updating an equivalent item over creating a new item
  - only add when no equivalent item exists
- Source reconciliation must handle:
  - local source still in `docs/exec-plan/todo/` or `docs/issues/`
  - local source moved to `docs/exec-plan/done/` or `docs/issues/done/`
  - local source missing
  - PR URL open, closed, merged, or superseded
  - GitHub Issue open or closed
- Existing item corrections should normalize old bodies to the durable handoff
  body format when title, source, priority, repo, or kind is already being
  updated.
- Full re-triage should classify findings before mutation:
  - direct backlog item: exec plan, local issue, open PR, open GitHub Issue
  - roadmap gap: unchecked project-plan phase or unmet requirement
  - do not add yet: vague or duplicate roadmap gap better represented by an
    existing direct item
- Missing local repo checkouts must be reported explicitly. GitHub PR/issue
  state may still be collected, but local docs inspection must be marked
  unavailable instead of inferred.
- Workflow-sync PRs should maintain at most one active item per repo for the
  current sync PR. When a newer open sync PR supersedes an old one, update the
  existing item rather than adding another.
- Full re-triage should build a short mutation plan before changing the Project:
  - Done updates
  - PR source/title/body updates
  - priority/kind/repo/body normalization
  - new adds
  - final sync/list
- When `tools/pj` batch update is available, the skill may use it for existing
  item updates; otherwise it should use single-item `pj update`.
- Subagent output should follow a strict schema with stable field names,
  workspace repo basenames, and canonical priority casing.
- The final briefing should include:
  - Project URL
  - synced item count
  - number of items marked Done
  - number of items updated
  - number of items added
  - high-priority Todo shortlist
  - caveats such as missing local checkouts or GitHub-only inspection
  - the standard numbered next choices

### `docs/specs/github-projects-task-cli.md`

- No direct CLI behavior change is required by this plan.
- If execution references future batch updates in the triage spec, keep wording
  conditional so this plan does not depend on the batch command being complete.

## Code Changes

### `skills/triage-tasks/SKILL.md`

- Strengthen the subagent delegation instruction so routine read-only
  exploration explicitly selects an available low-cost small model.
- Add an "Existing Board Reconciliation" step before Project population.
- Replace "Create a Project item" wording with "update-or-add Project items".
- Add stale source rules for local paths, PRs, and GitHub Issues.
- Require normalization of old `Source:`-only item bodies when touching the
  item for another reason.
- Add the direct backlog / roadmap gap / do-not-add-yet classification.
- Add missing-checkout handling.
- Add workflow-sync PR replacement rules.
- Add mutation-plan ordering and note future batch update usage as an
  optimization.
- Replace loose subagent output guidance with the strict schema.
- Add the final full re-triage briefing format.
- Keep the current user-confirmation boundary: the skill must not auto-execute a
  selected task after triage.

### Optional supporting docs

- Update `docs/specs/triage-tasks.md` and `skills/triage-tasks/SKILL.md`
  together so the spec and skill stay in sync.
- Do not edit `AGENTS.md` unless execution finds top-level workspace guidance is
  actively misleading.

## Design Decisions

No ADR update is expected. This plan tightens an existing workflow contract
rather than changing workspace architecture.

Past decision: the workspace optimizes for AI-readable context and low-cost
review artifacts. Apply the same reasoning here: full re-triage should produce
structured summaries and deliberate mutation plans instead of raw broad output
or repeated remote updates.

## Sub-tasks

- [ ] Review the current `docs/specs/triage-tasks.md` and
      `skills/triage-tasks/SKILL.md` to identify exact insertion points.
- [ ] [parallel] Update the spec with the full re-triage reconciliation,
      classification, stale source, missing checkout, mutation-plan, and final
      briefing rules.
- [ ] [parallel] Draft the strict subagent output schema and model-selection
      wording in the skill.
- [ ] [depends on: spec update] Update the skill flow so Step 3 uses
      update-or-add instead of add-first population.
- [ ] [depends on: spec update] Add workflow-sync PR replacement and stale
      local-source handling to the skill.
- [ ] [depends on: skill update] Verify the skill still preserves the standard
      Step 1 choices and Step 4 fresh-session handoff.
- [ ] [depends on: verification] Run documentation/workflow lint checks that
      apply to changed Markdown and skill files.
- [ ] [depends on: verification] Move this plan to
      `docs/exec-plan/done/triage-tasks-retriage-hardening.md`.

## Parallelism

- Spec wording and skill schema drafting can happen in parallel after the
  insertion points are reviewed.
- Stale source rules and workflow-sync replacement rules are independent within
  the skill update, but both should be reconciled before final mutation-plan
  wording is checked.
- Verification depends on both the spec and skill updates being complete.

## Verification

- Confirm `docs/specs/triage-tasks.md` and `skills/triage-tasks/SKILL.md`
  agree on update-before-add reconciliation.
- Confirm read-only subagent guidance requires explicit low-cost small-model
  selection without naming a fixed model version.
- Confirm stale local sources, stale PRs, closed GitHub Issues, and missing
  checkouts have defined handling.
- Confirm old item body normalization is required when an item is already being
  updated.
- Confirm full re-triage output includes mutation counts, caveats, high-priority
  shortlist, Project URL, and the standard next choices.
- Confirm the skill does not require the future `pj` batch update command to
  exist before it can still perform triage.

## Expected Outcome

- Full re-triage updates the existing board instead of creating avoidable
  duplicates.
- Stale Project items are closed or corrected consistently.
- Routine repo exploration uses cheaper read-only subagents by default.
- The final briefing is predictable enough for a human to understand what
  changed and choose the next task quickly.
