# Project Owner Scope Configuration

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Make `tools/pj` explicitly track which GitHub owner scope it is operating against so the workspace does not silently create or query the canonical `Workspace Task Triage` Project under the wrong target.

The intended outcome is:

- the operator can intentionally choose `user` or `org`
- the chosen owner scope is persisted in local metadata once established
- later commands such as `sync`, `add`, `move`, and `repo-link` operations can reuse that stored configuration without re-specifying it every time

This supports the project-plan requirement that the GitHub Projects-based workspace triage flow be practical enough for repeated daily use without surprising behavior.

## Background

The current CLI accepts `--owner` and `--owner-type`, but the workflow still leaves room for confusion about which owner scope is the canonical target for a given local workspace.

That creates two related risks:

1. the operator expects an existing board under one owner scope but the CLI resolves or creates under another scope, which looks like "the board does not exist"
2. repeated commands become noisy because the operator has to keep restating the same owner metadata even after the workspace has already established a canonical target

The local cache already persists project identity after bootstrap and sync. This plan extends that idea into a first-class owner-scope configuration model so the workspace can intentionally remember its canonical Project owner target.

## Spec Changes

### `docs/specs/github-projects-task-cli.md`

- Clarify that the workspace must operate against one explicitly chosen owner target at a time:
  - `owner`
  - `owner_type` (`user` or `org`)
- Define how the CLI establishes that target:
  - explicit flags on first use
  - optional explicit configuration/init command if that proves clearer than implicit cache seeding
- Define how later commands reuse cached owner metadata when flags are omitted
- Define how the CLI behaves when the operator provides flags that conflict with cached owner metadata
- Define whether changing owner scope requires an explicit reset or override action rather than silent drift

### `docs/specs/triage-tasks.md`

- Clarify that the canonical workspace board is owner-scoped and that the workspace's local metadata remembers which owner scope is active
- Clarify how a single local workspace should behave if the operator wants to switch from a personal board to an organization board later

### `AGENTS.md`

- Update workspace task tracking guidance to reflect that `pj` commands can reuse the locally established owner scope
- Add any operator guidance needed to avoid accidental cross-owner confusion

### `README.md`

- Document the first-run setup flow for selecting the owner scope
- Document how to inspect or override the stored owner target when needed

## Code Changes

### `tools/pj/`

- Add a clear owner-scope configuration model on top of the existing cache metadata
- Decide whether this is exposed as:
  - improved semantics for `pj init`
  - a dedicated config command
  - both
- Ensure `sync`, `add`, `move`, and future repo-link commands reuse the stored owner target when flags are omitted
- Add clear errors or confirmation requirements when a command would operate against a different owner target than the stored one
- Consider whether the cache artifact alone is sufficient or whether a separate config artifact is needed for clearer intent

### Cache / metadata

- Persist enough information to make the active owner target explicit and inspectable
- Avoid ambiguous fallback behavior where partially supplied flags silently mix old and new owner metadata

## Design Decisions

- The workspace should treat the owner target as an explicit configuration decision, not an incidental side effect of whichever flags happened to be used most recently
- Once a canonical owner target is established locally, normal day-to-day commands should not require restating it
- Switching owner targets should require an explicit operator action so the CLI does not silently create or query boards in the wrong place
- The configuration model should remain simple enough that the workspace still feels lighter than the previous beads-based flow

## Sub-tasks

- [ ] Update `docs/specs/github-projects-task-cli.md` with an explicit owner-target configuration model
- [ ] Update `docs/specs/triage-tasks.md`, `AGENTS.md`, and `README.md` so operator guidance matches the owner-target model
- [ ] [parallel] Review the current cache/flag merge behavior in `tools/pj` and identify where mixed old/new owner metadata could cause surprising behavior
- [ ] [parallel] Decide whether owner-target persistence should stay inside `.local/pj/cache.json` or move into a dedicated config artifact
- [ ] [depends on: cache/flag review, config storage decision] Implement explicit owner-target persistence and reuse
- [ ] [depends on: owner-target persistence] Add safe behavior for owner-target overrides or resets
- [ ] [depends on: implementation] Verify first-run bootstrap, repeated sync/add/move usage, and intentional owner-target switching

## Verification

- Confirm first-run setup requires or clearly establishes the intended `owner` and `owner_type`
- Confirm repeated commands can omit owner flags once local metadata is established
- Confirm conflicting owner inputs do not silently drift the workspace to a different owner target
- Confirm the stored owner target is easy to inspect from local metadata and/or CLI output
- Confirm docs/specs describe the owner-target behavior consistently

## Expected Outcome

- `tools/pj` no longer feels ambiguous about whether it is targeting a user-owned or org-owned Project
- A workspace can bootstrap once and then use the same owner target consistently across later commands
- Switching to a different owner target becomes an intentional, explicit operation instead of an accidental side effect
