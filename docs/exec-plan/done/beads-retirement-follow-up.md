# Beads Retirement Follow-up

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Finish the workspace migration away from beads by removing or archiving the remaining `.beads/` artifacts, deleting stale ignore patterns and workflow remnants, and aligning the surviving docs/specs with the GitHub Projects-based triage flow.

This closes the deferred cleanup captured in `docs/issues/beads-retirement-follow-up.md` and reinforces the project-plan direction that workspace task coordination now uses GitHub Projects plus the local `pj` cache instead of beads/Dolt state.

## Current Context

The GitHub Projects spike already made beads non-canonical, but the workspace still contains:

- a tracked `.beads/` directory with historical backup and Dolt runtime data
- beads-specific `.gitignore` entries that are only useful while those artifacts remain
- historical documentation that still references the old beads-based workflow as active context

The current `docs/specs/triage-tasks.md` already declares the GitHub Projects flow, so this follow-up should finish the operational cleanup without reintroducing ambiguity about the source of truth.

## Spec Changes

### `docs/specs/triage-tasks.md`

- Clarify that `.local/pj/` is the only supported local workspace triage state for the current workflow
- State that the workspace must not depend on committed `.beads/` artifacts or Dolt runtime state for task coordination
- Add any brief operator note needed if historical beads data is intentionally archived elsewhere instead of kept in-repo

### Workflow docs sweep

- Review `AI_WORKFLOW.md`, `AGENTS.md`, and any active skill/spec text that still implies beads artifacts are part of the supported workflow
- Keep historical records in completed plans/issues where useful, but remove any wording that reads like current operational guidance

## Code Changes

### Workspace cleanup

- Decide the concrete retirement path for `.beads/`:
  - delete it entirely if no history needs to be preserved in-repo, or
  - move only the minimum historical notes needed into Markdown under `docs/` before deletion
- Remove obsolete beads-specific `.gitignore` rules once `.beads/` is no longer a supported workspace artifact
- Remove any leftover helper/config files that only existed to support beads runtime behavior, if they remain after the directory cleanup

### Reference cleanup

- Sweep the repo for stale runtime-facing references to `bd`, `.beads/`, and Dolt backup behavior
- Update or trim those references so only historical context remains in done plans/issues, not active workflow guidance

## Design Decisions

- Follow the existing `core-beliefs.md` guidance: prefer the simplest cleanup that preserves correctness and AI-readable context
- Unless execution uncovers a concrete need for long-term beads history, do not keep bulky runtime artifacts in-repo just for nostalgia
- If execution decides to preserve any history, prefer compact Markdown documentation under `docs/` over retaining the raw `.beads/` runtime directory
- No ADR update is expected unless execution reveals a broader archival policy that should govern future tool retirements

## Sub-tasks

- [ ] [parallel] Inventory the current `.beads/` contents and confirm what historical value, if any, must be preserved before deletion
- [ ] [parallel] Update `docs/specs/triage-tasks.md` and any active workflow docs to state the post-beads contract explicitly
- [ ] [depends on: inventory] Remove `.beads/` or replace it with a compact archived note, following the chosen preservation path
- [ ] [depends on: inventory] Remove obsolete beads-specific `.gitignore` entries and any leftover runtime-only helper files
- [ ] [depends on: docs update, cleanup] Sweep for remaining active beads references and either delete them or rewrite them as historical context
- [ ] [depends on: all above] Verify the repo still communicates GitHub Projects + `pj` as the only supported workspace triage flow, then move `docs/issues/beads-retirement-follow-up.md` to `docs/issues/done/`

## Parallelism

- The inventory/preservation decision and the doc/spec update can proceed independently at first
- File deletion and `.gitignore` cleanup depend on the inventory decision
- Final reference sweep should happen after the chosen cleanup path is implemented so docs match the resulting repository state

## Verification

- Confirm `.beads/` is either removed or intentionally replaced by a documented lightweight archive path
- Confirm `.gitignore` no longer contains beads-specific rules that are obsolete after cleanup
- Confirm active workflow docs/specs no longer present beads/Dolt artifacts as part of the supported workspace flow
- Confirm any remaining `bd`/`.beads/` mentions are limited to historical records such as completed plans, issues, or retrospective notes

## Expected Outcome

- Workspace task coordination is unambiguously GitHub Projects plus `.local/pj/`
- Bulky legacy beads runtime artifacts are gone from the active repo state
- Workflow docs stop mixing current GitHub Projects guidance with stale beads-era cleanup leftovers
