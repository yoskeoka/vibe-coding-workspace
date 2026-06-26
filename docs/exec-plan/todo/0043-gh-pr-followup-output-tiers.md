# GH PR Follow-up Output Tiers
> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Reduce token churn during repeated PR follow-up polling by splitting
`skills/review-task/scripts/gh-pr-followup poll` into a default compact tier and
an opt-in verbose tier.

The existing helper already removed raw GitHub API payloads, but the default
output still returns more review/check/timeline detail than the common
`review-task` landing loop usually needs. This plan keeps the correctness
contract of post-PR follow-up while making the default polling surface smaller.

This supports the workspace goal of keeping hobby workflow cost low and matches
`core-beliefs.md` items 4 and 5: prefer compact review artifacts and trim tool
output at the source.

## Current State

- `docs/specs/pr-follow-up-workflow.md` defines `gh-pr-followup poll` as the
  preferred repeated polling path for `review-task`.
- `skills/review-task/scripts/gh-pr-followup` already compacts GitHub reads and
  deduplicates old timeline/comment entries through `.local/gh-pr-followup/`
  markers.
- The current helper still returns full review bodies, full compact check lists,
  and all new timeline/comment entries on every normal poll, even when the
  caller only needs head SHA, required-check summary, and whether advisory
  review started.
- `ai-arena/docs/issues/0034-reduce-token-heavy-verification-command-surfaces.md`
  identified this helper as a separate concern from repo-owned verification
  command surfaces, so the workflow-side change should live in its own plan/PR.

## Spec Changes

### `docs/specs/pr-follow-up-workflow.md`

- Define output tiers for `gh-pr-followup poll`.
- Require the default poll surface to include only the fields needed for the
  standard landing decision:
  - `repo`
  - `pr`
  - `head_sha`
  - `review_decision`
  - compact required-check summary
  - advisory-review-start detection data
  - newly observed inline comments only when present
- Define an opt-in verbose tier for detailed review bodies, expanded timeline
  context, and full compact check entries.
- Clarify when `review-task` should use default mode versus verbose mode.

### `skills/review-task/SKILL.md`

- Update the follow-up loop to call the default compact tier for ordinary
  landing polls.
- State when the operator should escalate to verbose mode, such as helper
  diagnosis or substantive advisory review triage.

## Code Changes

### `skills/review-task/scripts/gh-pr-followup`

- Add an output-tier interface without changing the helper's ownership or local
  state model.
- Keep the current default command path short for callers already using
  `gh-pr-followup poll <owner/repo> <pr-number>`.
- Add an opt-in verbose path that emits the richer review/timeline/check detail
  currently returned by default.
- Ensure marker updates and head-SHA reset behavior remain identical across
  tiers.
- Keep failure behavior unchanged: helper errors still stop automatic follow-up
  instead of falling back to raw `gh`.

### `docs/issues/` and related docs

- No new workspace issue file is required unless execution uncovers a separate
  workflow gap beyond the output-tier split itself.

## Design Decisions

Past decisions reviewed before planning:

- `core-beliefs.md` says human review should prefer compact artifacts over token
  burn and that large tool output should be trimmed at the source.
- `docs/exec-plan/done/gh-pr-followup-token-trimming.md` established the helper
  as the workflow-owned compact polling path instead of repeated raw GitHub API
  reads.
- `docs/specs/pr-follow-up-workflow.md` and prior review-task plans keep
  `review-task` as the sole owner of the bounded PR follow-up loop.

Apply the same reasoning here:

- Keep `review-task` as the owner; do not create a second polling helper.
- Preserve the existing `poll` entrypoint so current callers do not need a broad
  migration.
- Move richer output behind an explicit verbose path instead of returning it on
  every landing poll.

Viable interface options:

1. Add `poll --verbose` while keeping plain `poll` as the compact default.
2. Add `poll --mode=compact|verbose` and require every caller to choose.

Recommended option: option 1. It keeps the common path shortest, preserves the
existing command shape for current callers, and matches the goal of making the
default surface cheaper rather than merely configurable.

No ADR update is expected unless execution reveals a broader workflow-wide
pattern for tiered helper output beyond `gh-pr-followup`.

## Sub-tasks

- [ ] [parallel] Update `docs/specs/pr-follow-up-workflow.md` with the default
      compact tier and verbose-tier contract.
- [ ] [parallel] Update `skills/review-task/SKILL.md` so ordinary PR landing
      polls stay on the default tier and verbose mode is reserved for explicit
      deeper inspection.
- [ ] [depends on: spec update] Adjust
      `skills/review-task/scripts/gh-pr-followup` to support tiered output while
      preserving marker/state behavior.
- [ ] [depends on: helper update] Add focused verification for default versus
      verbose output shape and unchanged marker-reset behavior.
- [ ] [depends on: all above] Run workflow lint and relevant shell/script checks,
      then prepare the planning PR through `review-task`.

## Parallelism

- Spec wording and `review-task` skill wording can proceed in parallel.
- Helper implementation depends on the chosen output-tier contract.
- Verification depends on the final helper behavior.

## Verification

- `./tools/workflow-lint.sh --mode=pre-push`
- shell syntax check for `skills/review-task/scripts/gh-pr-followup`
- targeted smoke test with fake `GH_BIN` or equivalent fixture to verify:
  - plain `poll` omits verbose-only fields
  - verbose mode returns the richer review/timeline/check detail
  - marker reset and deduplication still behave the same when head SHA changes
- manual spec/skill parity check against helper behavior

## Expected Outcome

- Ordinary `review-task` landing polls paste a smaller, decision-oriented JSON
  payload into agent context.
- Detailed review bodies and expanded timeline/check context remain available,
  but only through an explicit verbose path.
- The helper stays the canonical workflow-owned polling surface without
  weakening CI/advisory-review follow-up correctness.
