# ww Unknown Status For Missing Default Base

**Execution**: Use `/execute-task` to implement this plan.

## Objective

Make `ww list` and `ww clean` robust when a detected repository does not expose a confirmed default base via `origin/HEAD`. Instead of failing the entire command, `ww` should keep listing worktrees, surface degraded status information, and exclude uncertain worktrees from cleanup.

This plan resolves [ww-status-default-base-unknown.md](/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/docs/issues/ww-status-default-base-unknown.md).

## Background

The current status pipeline assumes that every repository can resolve a base branch from either:

- `default_base`
- `refs/remotes/origin/HEAD`

That assumption is too strict for workspace mode. A repository may lack `origin/HEAD` while still having enough information to infer likely defaults such as `origin/main`. Today that condition hard-fails `ww list`, which is not acceptable for a best-effort visibility command.

## Scope

Project: `ww/`

Primary goals:

- Add `unknown` as a first-class status
- Add `status_detail` so degraded states are machine-readable
- Render degraded states in text output as `unknown(<detail>)`
- Teach base resolution a best-effort inference sequence
- Ensure `--cleanable` and `ww clean` never act on `unknown`

## Design Decisions

### Status model

Status values become:

- `active`
- `merged`
- `stale`
- `unknown`

`status_detail` is added for degraded cases. Initial values:

- `no-origin-head`
- `heuristic-base`
- `no-remote`
- `remote-query-failed`
- `base-detect-failed`

### Default-base resolution order

When `default_base` is unset, resolve in this order:

1. `origin/HEAD`
2. local `main` tracking `origin/main`
3. multiple local branches converging on upstream `origin/main`
4. existence of remote ref `origin/main`
5. otherwise unresolved

If step 2-4 succeeds, treat the inferred base as heuristic. It may be useful for diagnostics, but worktree status remains `unknown(heuristic-base)` unless the base came from explicit config or `origin/HEAD`.

### Cleanup safety

- `ww list --cleanable` remains limited to `merged` and `stale`
- `ww clean` remains limited to `merged` and `stale`
- `unknown` is never treated as cleanable

## Spec Changes

Update `ww/docs/specs/cli-commands.md`:

- add `unknown` to the `ww list` status table
- document `status_detail` in JSON output
- show degraded text output as `unknown(<detail>)`
- clarify that `--cleanable` filters only `merged` and `stale`
- clarify that `ww clean` skips `unknown`

Update `ww/docs/specs/git-operations.md`:

- document default-base inference order
- document which sources are authoritative vs heuristic
- document that missing `origin/HEAD` no longer hard-fails listing

Update `ww/docs/specs/configuration.md`:

- clarify that explicit `default_base` remains authoritative

## Code Changes

### `ww/git`

- replace the single-step default-branch lookup with a richer resolver that can report:
  - resolved base ref
  - resolution source
  - failure/detail classification
- add helpers for:
  - reading local branch upstreams
  - testing whether `origin/main` exists
  - distinguishing `origin/HEAD` absence from other failures

### `ww/worktree`

- extend `WorktreeInfo` with `StatusDetail`
- change status computation to degrade to `unknown` instead of returning an error for:
  - missing `origin/HEAD`
  - heuristic-only base
  - missing branch remote
  - remote query failure
- keep `merged` / `stale` logic only for authoritative base resolution

### `ww/cmd/ww`

- update `list` text output to render `unknown(<detail>)`
- keep JSON output machine-readable with `status` plus `status_detail`
- keep cleanable filtering restricted to `merged` and `stale`

## Tests

- Add unit tests for default-base resolution order and source classification
- Add list tests covering:
  - missing `origin/HEAD` but `origin/main` exists
  - branch with no tracking remote
  - remote query failure
  - authoritative `default_base`
- Add integration coverage proving:
  - `ww list` succeeds and emits `unknown(...)` instead of failing
  - `ww list --cleanable` excludes `unknown`
  - `ww clean` leaves `unknown` worktrees untouched

## Sub-tasks

- [ ] [parallel] Update specs for `unknown` / `status_detail` behavior
- [ ] [parallel] Implement default-base resolver with authoritative vs heuristic sources
- [ ] [depends on: default-base resolver] Update worktree status classification and `WorktreeInfo`
- [ ] [depends on: worktree status classification] Update CLI text and JSON rendering
- [ ] [depends on: status classification] Add unit and integration tests for degraded status behavior
- [ ] [depends on: all above] Verify `ww list`, `ww list --cleanable`, and `ww clean` behavior in workspace mode

## Verification

- `go test ./...` in `ww/`
- Manual reproduction of the current failure case from the workspace root
- Confirm that the affected repo now shows `unknown(<detail>)` and that no `unknown` worktree is cleaned
