# pj Client Factory Injection Follow-up
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Remove the remaining package-level client factory mutation from `tools/pj` command tests so the GitHub Projects task CLI can keep growing without weaker test isolation or serialized tests around shared global state.

This closes the design debt captured in `docs/issues/pj-client-factory-injection-follow-up.md` and supports the project-plan requirements that the GitHub Projects-backed workspace triage CLI remain practical, low-friction, and reusable enough to evaluate beyond this workspace.

## Context

Current `tools/pj/internal/pj/app.go` defines a package-level `newProjectClient` variable. Tests in `tools/pj/internal/pj/app_test.go` temporarily replace that variable to inject `stubProjectClient` into `runInit`, `runSync`, `runAdd`, and `runUpdate`.

That pattern works only while affected tests stay serialized. It also makes future command tests easier to write incorrectly because the dependency injection point is shared mutable process state.

Past decision: `tools/pj` is an intentionally small Go CLI spike under `tools/pj/`, backed by GitHub Projects as source of truth and `.local/pj/` as derived cache/config. This plan should improve internal testability without changing the external command contract or expanding the spike into a broader framework.

## Design Options

### Option A: pass a client factory to each `runX` function

- Change `runInit`, `runSync`, `runAdd`, and `runUpdate` to accept a `func() (projectClient, error)`.
- Keep `Run` as the only place that supplies the production factory.
- Tests call each runner with a local factory.

This is minimal, but it spreads the same dependency parameter across every remote-backed command and makes future command dependencies likely to be added one argument at a time.

### Option B: introduce an internal `app` dependency struct

- Add an unexported struct, for example `type app struct { newProjectClient func() (projectClient, error) }`.
- Make exported `Run` construct a default app with `newGitHubClient`.
- Move command dispatch and remote-backed runners onto methods such as `a.runInit`, `a.runSync`, `a.runAdd`, and `a.runUpdate`.
- Tests instantiate an app with a local stub factory instead of mutating package globals.

This keeps the production CLI surface unchanged while centralizing command dependencies in one internal type. It also leaves room for future IO or clock dependencies without changing every runner signature repeatedly.

### Recommended Direction

Use Option B. It aligns with the core belief to avoid broad refactors for aesthetics because the change is scoped to the actual test isolation problem, while still giving the CLI a clear place to own command dependencies.

No ADR update is expected unless execution uncovers a broader architectural change beyond internal dependency injection.

## Spec Changes

Update `docs/specs/github-projects-task-cli.md` with a compact implementation/testability note:

- `Run(args, stdout, stderr)` remains the public command entrypoint.
- Remote-backed commands must obtain GitHub Project clients through an app-owned dependency, not through a mutable package-level variable.
- Tests should inject stub clients through that app dependency so command tests can run independently and safely regain `t.Parallel()` where file/cache fixtures are isolated.

No user-facing CLI behavior, cache schema, command flags, or GitHub Project semantics should change. `docs/specs/triage-tasks.md` should not need changes unless execution finds user-facing workflow guidance affected by the refactor.

## Code Changes

### `tools/pj/internal/pj/app.go`

- Replace package-level `newProjectClient` mutation with a default app construction path.
- Keep `func Run(args []string, stdout, stderr io.Writer) error` as the exported entrypoint used by `cmd/pj`.
- Move dispatch and remote-backed runner methods behind the internal app dependency.
- Ensure local-cache-only commands such as `config`, `list`, `url`, and `open` keep their behavior unchanged.
- Ensure `init`, `sync`, `add`, and `update` still create clients lazily only after argument/config/cache validation reaches the point where remote access is needed.

### `tools/pj/internal/pj/app_test.go`

- Replace save/restore mutations of `newProjectClient` with per-test app construction.
- Restore `t.Parallel()` to tests that were serialized only because they mutated the package-level client factory.
- Keep temp-dir config/cache fixtures per test.
- Add or adjust assertions to confirm each remote-backed command uses its injected stub client.

### `docs/issues/`

- Move `docs/issues/pj-client-factory-injection-follow-up.md` to `docs/issues/done/` during execution after the implementation and verification pass.

## Sub-tasks

- [ ] [parallel] Update `docs/specs/github-projects-task-cli.md` with the internal dependency-injection/testability contract.
- [ ] [parallel] Inspect current serialized tests and mark which can regain `t.Parallel()` once global mutation is removed.
- [ ] [depends on: spec update] Introduce the internal app dependency struct and default production constructor in `tools/pj/internal/pj/app.go`.
- [ ] [depends on: internal app dependency] Convert remote-backed command runners to app methods while preserving exported `Run` behavior.
- [ ] [depends on: internal app dependency] Convert tests to instantiate an app with a stub client factory instead of mutating a package global.
- [ ] [depends on: tests converted] Re-enable safe parallel tests and verify no test depends on shared mutable package state.
- [ ] [depends on: implementation verified] Move the follow-up issue file to `docs/issues/done/`.

## Verification

- Run `go -C tools/pj test ./...`.
- Run `make lint` if the workspace linter covers the touched docs/workflow files.
- Confirm `tools/pj/internal/pj/app_test.go` has no save/restore mutation of a package-level `newProjectClient`.
- Confirm `pj` command help and command output remain unchanged for representative local-cache commands:
  - `go -C tools/pj run ./cmd/pj help`
  - `go -C tools/pj run ./cmd/pj list --cache <temp-cache>`

## Non-goals

- Do not change GitHub GraphQL query or mutation behavior.
- Do not change `.local/pj/config.json` or `.local/pj/cache.json` schema.
- Do not add a new public Go API for `tools/pj`; the dependency struct is internal test plumbing.
- Do not fold unrelated follow-ups such as pagination or repository linking into this plan.

## Review Notes

The key review question is whether the internal dependency struct stays small. If execution starts adding abstractions unrelated to command dependency injection, split that work into a separate issue or plan instead of expanding this one.
