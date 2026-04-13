# Project Field Provisioning for `tools/pj`

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Remove the remaining manual GitHub Project schema setup from the workspace triage flow by teaching `tools/pj` to provision the required custom ProjectV2 fields during bootstrap.

After this plan is executed, `pj init` should leave the canonical `Workspace Task Triage` board ready for normal `sync`, `list`, `add`, and `move` usage without requiring the operator to create `Repo`, `Kind`, or `Priority` by hand first.

This supports the project-plan requirement that workspace task coordination stays low-friction and GitHub Projects-backed, with `.local/pj/cache.json` remaining derived data rather than a second source of truth.

## Background

Current behavior creates or resolves the canonical board and writes cache metadata, but then fails if the board is missing custom fields. That leaves bootstrap half-finished on a fresh board and keeps a manual GitHub UI step in the critical path.

Past decision: the original CLI spike explicitly treated custom-field auto-provisioning as out of scope. The new issue in [docs/issues/pj-project-field-provisioning.md](/docs/issues/pj-project-field-provisioning.md) narrows the next slice: extend bootstrap to provision the required field model without changing the existing cache-first, single-board workflow.

Relevant existing constraints:

- `tools/pj` already treats field metadata in `.local/pj/cache.json` as derived data populated from GitHub, so provisioning must refresh and persist remote IDs after mutation rather than inventing local IDs.
- `add` and `move` currently assume single-select fields and option-name-to-option-ID lookup, so the provisioned schema should stay aligned with that model.
- Owner-scope configuration and repository linking are being clarified in separate plans; this plan should not expand into those areas beyond keeping `pj init` compatible with them.

## Spec Changes

### `docs/specs/triage-tasks.md`

- Update the minimum-field contract so `Repo`, `Kind`, and `Priority` are expected single-select workflow fields, not manually pre-created prerequisites.
- Define `pj init` as responsible for provisioning missing workflow fields during bootstrap.
- Require provisioning to be idempotent for an already-compatible board.

### `docs/specs/github-projects-task-cli.md`

- Change `pj init` from "fail if custom fields are missing" to "create or reconcile the required workflow fields, then validate compatibility".
- Define the failure mode when an existing field has the wrong type or lacks required option metadata.
- Keep schema management scoped to the required workspace fields only; arbitrary field management stays out of scope.

## Code Changes

### `tools/pj/internal/pj/github.go`

- Add GraphQL mutation support for creating ProjectV2 single-select fields on the canonical board.
- Add read-after-write sync logic so provisioned field IDs and option IDs are immediately persisted into cache.
- Decide how to handle an existing field with the right name but incompatible shape:
  - detect and fail clearly for the first implementation
  - avoid silent destructive mutation or duplicate-field creation

### `tools/pj/internal/pj/app.go`

- Fold field provisioning into `runInit` so bootstrap creates the board and required field model in one idempotent flow.
- Keep `sync` read-only; it should report schema problems but not mutate them.
- Keep `add` and `move` unchanged except for benefiting from the stronger `init` contract.

### `tools/pj/internal/pj/types.go`

- Add a canonical workflow-field schema definition that code and tests can share:
  - required field names
  - required field type (`single-select`)
  - initial option sets needed by the workspace triage flow

### Tests

- Extend GraphQL-focused tests to cover:
  - existing compatible board: no new fields created
  - freshly created board: required fields are provisioned
  - partially configured board: only missing fields are created
  - incompatible existing field: `pj init` fails with a compatibility error
  - cache refresh after provisioning includes the new field and option IDs

## Design Decisions

- Provisioning belongs in `pj init`, not a separate required bootstrap command. `init` already owns "make the canonical board usable" and splitting field creation into a second mandatory command would preserve the current footgun.
- `sync` remains read-only. Mutating schema during sync would make routine refreshes surprising and would blur the line between bootstrap and day-to-day operations.
- The first implementation should provision only the required workspace fields and canonical options needed by current commands. Broader schema reconciliation can remain future work.
- Existing incompatible fields should fail fast instead of being rewritten automatically. This matches the current preference for correctness over convenience when remote state is ambiguous.

## Sub-tasks

- [ ] Finalize the canonical option sets for provisioned `Repo`, `Kind`, and `Priority` fields so they cover current workspace triage usage without overfitting future workflows
- [ ] Update `docs/specs/triage-tasks.md` and `docs/specs/github-projects-task-cli.md` to define bootstrap provisioning and compatibility rules
- [ ] [parallel] Review the current `runInit` and cache-refresh flow to identify the cleanest insertion point for provisioning without duplicating sync logic
- [ ] [parallel] Verify the GitHub GraphQL mutation shape and response data needed to create ProjectV2 single-select fields and re-read option metadata
- [ ] [depends on: init flow review, GraphQL mutation verification] Add an internal canonical field-schema definition shared by init logic and tests
- [ ] [depends on: canonical field schema, GraphQL mutation verification] Implement idempotent provisioning inside `pj init`
- [ ] [depends on: provisioning implementation] Add compatibility validation for wrong-type or under-specified existing fields
- [ ] [depends on: provisioning implementation] Extend unit tests for no-op, create, partial-create, and incompatible-schema paths
- [ ] [depends on: tests] Run a controlled manual verification against a disposable or otherwise non-primary ProjectV2 target before using the flow on the long-lived workspace board
- [ ] [depends on: manual verification] Re-run `pj init` against the same target to confirm idempotent no-op behavior and stable cache metadata

## Verification

- Confirm the updated specs describe `pj init` as provisioning required custom fields instead of requiring manual pre-creation
- Confirm `pj init` creates missing `Repo`, `Kind`, and `Priority` fields on a board that only has GitHub's default schema
- Confirm a second `pj init` run against an already-compatible board does not create duplicate fields
- Confirm cache refresh after provisioning contains the created field IDs and single-select option IDs used by `add` and `move`
- Confirm `pj init` fails clearly when a required field name exists with an incompatible type or missing required options
- Confirm `pj sync` still behaves as a read-only cache refresh and surfaces incompatibilities without mutating remote schema

## Manual Verification Strategy

- Prefer a disposable ProjectV2 target owned by the same GitHub account or a dedicated sandbox owner for the first real mutation test; do not make the long-lived workspace board the first place new schema mutations are exercised.
- Use unit tests with `httptest` as the primary safety net for mutation ordering and error handling, because they can cover partial-field and incompatible-schema scenarios that are awkward to stage safely in the live board.
- If a disposable remote target is not available, treat manual verification on the canonical board as a one-time additive migration and only after mocked tests prove the no-op and compatibility paths.

## Expected Outcome

- `pj init` becomes a complete bootstrap step for the workspace board instead of a partial setup command that still depends on GitHub UI work
- The workspace cache continues to mirror remote Project metadata after bootstrap, including field and option IDs needed for later mutations
- Operators can create or recover the canonical workspace triage board with one command, while incompatible remote schema still fails loudly instead of being silently rewritten
