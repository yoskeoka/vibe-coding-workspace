# pj Update Batch Command Design

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Add a batch update capability to `tools/pj` so workspace triage can apply many
Project item mutations without paying the current one-command, one-cache-refresh
cost for every item.

This plan supports the project-plan requirements that workspace task
coordination uses a low-friction GitHub Projects-backed flow and that the local
CLI spike stays practical enough to evaluate for broader reuse.

The execution phase MUST design and document the batch command surface before
implementing code. The implementation PR is not allowed to start from an
implicit command shape.

## Background

Full re-triage commonly needs a short mutation plan:

- mark stale items `Done`
- update old PR items to current PR URLs and titles
- normalize item bodies to the durable handoff format
- adjust `Repo`, `Kind`, or `Priority`
- add only genuinely missing items

Today `pj update` refreshes the cache after every successful mutation. That is
correct for single-item operator use, but it makes bulk triage slow and noisy.
The current workaround is to run many `pj update` commands sequentially and
wait for each refresh.

Past decisions and constraints:

- GitHub Projects remains the canonical remote source of truth.
- `.local/pj/cache.json` remains derived data and must not become a second
  source of truth.
- `pj update` is already the single existing-item mutation path.
- Remote-backed commands should obtain GitHub clients through injected command
  dependencies so tests can stub network behavior.
- The CLI is still a workspace-local spike, so simple, explicit behavior is
  preferred over a large general-purpose automation language.

## Trade-offs

### Option A: Add `pj update-batch --file <path>`

- Keeps batch updates as a sibling of the existing single-item `update`
  command.
- Makes help output easy to scan because single and batch paths are distinct.
- Avoids overloading `pj update` flags with two modes.
- Recommended starting point for command design unless execution finds a
  clearer local convention.

### Option B: Add `pj update --batch-file <path>`

- Keeps all existing-item mutation under the `update` verb.
- Adds modal behavior: `--item` flags and batch-file input would need careful
  mutual-exclusion rules.
- Slightly more compact, but easier for operators and tests to misuse.

### Option C: Add a broader `pj batch --file <path>`

- Leaves room for future mixed operations such as add, update, and close in one
  file.
- Expands the design beyond the immediate pain point and risks inventing a
  mini workflow language before the spike needs one.
- Better as a later evolution if update-only batching proves too narrow.

Recommendation for execution: design a small update-only batch command first.
The command design must still be reviewed in the spec before code is written,
including input format, validation, cache-refresh semantics, error reporting,
and whether the final CLI name is `update-batch` or another explicit shape.

## Spec Changes

### `docs/specs/github-projects-task-cli.md`

- Add a command-design section for batch updates before implementation details.
- Define the final operator command syntax chosen during execution.
- Define the batch input format. The design must cover:
  - item ID
  - optional `title`
  - optional `body` or body-file equivalent
  - optional `status`
  - optional `repo`
  - optional `kind`
  - optional `priority`
- Define validation rules:
  - every operation must identify one Project item
  - unknown fields fail before remote mutation begins
  - empty operations fail clearly
  - invalid field values fail with the same resolver behavior as `pj update`
  - body inline input and body-file input remain mutually exclusive per item
- Define execution semantics:
  - apply operations in input order
  - refresh the cache once after all successful mutations, not after every item
  - fail clearly if any operation fails
  - document whether failure stops immediately or reports all preflight
    validation errors before contacting GitHub
  - never write a partially invented cache; cache state after partial remote
    mutation must come from a real sync or be left stale with a clear warning
- Define output semantics:
  - show a compact per-item result summary
  - print the final cache refresh status
  - make output suitable for agent handoff without dumping raw GraphQL payloads
- Update the command list and help contract.

### `docs/specs/triage-tasks.md`

- Describe that full re-triage should build a mutation plan first and may use
  the batch update command once available.
- Clarify that batch update is an optimization for updating existing items; it
  does not replace `pj add` for genuinely missing items.
- Keep the workflow compatible with single-item `pj update` as a fallback when
  batch update is unavailable or fails.

## Code Changes

### `tools/pj/internal/pj/app.go`

- Add the final batch command parser and help text selected by the command
  design.
- Reuse the existing single-item update parsing and validation logic where
  practical instead of creating a parallel resolver.
- Add a batch execution path that:
  - loads config/cache once
  - validates the full input before remote mutation when possible
  - applies remote item updates in order
  - refreshes cache once at the end after success
  - reports partial mutation states explicitly if a later item fails

### `tools/pj/internal/pj/`

- Factor shared update request types if needed so single and batch updates use
  the same value-resolution and GitHub mutation behavior.
- Keep GitHub client interfaces testable with stubs.
- Avoid adding general-purpose scripting or dependency semantics to the batch
  file format.

### Tests

- Add command tests for:
  - accepted batch input
  - unknown fields
  - missing item ID
  - no-op operation
  - invalid body/body-file combinations
  - invalid enum values
  - successful multi-item mutation with one cache refresh
  - remote failure with clear partial-state reporting
- Add focused tests that prove single-item `pj update` behavior is unchanged.

## Design Decisions

No ADR update is expected unless execution chooses a generalized `pj batch`
command that intentionally expands the CLI product direction beyond update-only
batching.

Past decision: keep the CLI small and workspace-local while using GitHub
Projects as canonical remote state. Apply the same reasoning here: batch update
should reduce command overhead without turning `.local/pj/` into an editable
source of truth or creating a broad automation language.

## Sub-tasks

- [ ] Review current `pj update` parsing, value resolution, GitHub mutation, and
      cache refresh flow.
- [ ] [parallel] Draft the batch command design in
      `docs/specs/github-projects-task-cli.md`, including syntax, input format,
      validation, output, and cache semantics.
- [ ] [parallel] Update `docs/specs/triage-tasks.md` to describe mutation-plan
      batching as an optional optimization after reconciliation.
- [ ] [depends on: command design] Implement the selected batch command in
      `tools/pj/internal/pj/app.go` and supporting package files.
- [ ] [depends on: command design] Add unit tests for parsing, validation,
      success, failure, and cache-refresh behavior.
- [ ] [depends on: implementation] Verify single-item `pj update` still works
      and still refreshes cache after a successful mutation.
- [ ] [depends on: tests] Run `go -C tools/pj test ./...` and any `tools/pj`
      lint/vet checks available in the repository.
- [ ] [depends on: verification] Move this plan to
      `docs/exec-plan/done/pj-update-batch-command-design.md`.

## Parallelism

- The command design and triage spec guidance can be drafted in parallel after
  the current `pj update` behavior is reviewed.
- Tests for input validation can be drafted once the command design is settled,
  while implementation starts on shared update request types.
- Remote failure tests depend on the final execution semantics for partial
  mutation reporting.

## Verification

- Confirm the spec records the final command design before code changes are
  reviewed.
- Confirm batch updates can update at least two items while refreshing cache
  once at the end.
- Confirm malformed batch input fails before remote mutation when validation can
  be performed locally.
- Confirm remote failure output makes partial mutation state explicit.
- Confirm `pj update --item ...` behavior remains compatible with the existing
  spec.
- Confirm `triage-tasks` guidance treats batch update as an optimization and not
  as a required dependency for normal triage.

## Expected Outcome

- Full re-triage can apply many existing-item corrections faster and with less
  command noise.
- The command surface is explicitly designed before implementation, avoiding an
  accidental file format or ambiguous `update` mode.
- The cache remains derived remote state, refreshed once from GitHub after
  successful batch mutation.
