# pj client factory injection follow-up

## Summary

`tools/pj/internal/pj/app_test.go` currently swaps the package-level `newProjectClient` variable to inject test doubles into `runInit` and `runSync`.

That makes the tests depend on global mutable state and prevents safe parallel execution. The current PR mitigates the immediate flake risk by removing `t.Parallel()` from the affected tests, but the underlying design remains.

Relevant files:
- `tools/pj/internal/pj/app.go`
- `tools/pj/internal/pj/app_test.go`

## Proposed Solution

Refactor the command runners so they can receive a client factory or app dependency struct without mutating package globals.

Possible directions:
- pass a client factory into `runInit`, `runSync`, `runAdd`, and `runMove`
- introduce an `App` struct that owns the client factory and IO dependencies

The goal is to make tests hermetic and allow safe parallel execution without package-level mutation.

## Priority

Medium. The immediate correctness risk is mitigated by serializing the tests, but the current design still makes test isolation weaker and future extensions more error-prone.
