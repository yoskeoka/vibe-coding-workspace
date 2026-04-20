# `tools/pj` Pagination Follow-up

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Teach `tools/pj` to fetch complete GitHub Projects data with cursor-based pagination instead of failing when a board exceeds the current single-page spike limits.

This resolves the follow-up tracked in `docs/issues/pj-pagination-follow-up.md` and supports the project-plan requirement that workspace task coordination uses a GitHub Projects-backed CLI reusable enough to evaluate for later standalone use.

## Background

`tools/pj` currently queries a fixed amount of GitHub Projects data during sync:

- project fields: `first: 50`
- project items: `first: 100`
- item field values: `first: 20`

The current implementation correctly fails when project fields or items report `hasNextPage`, which protects the cache from silent truncation. The next slice should replace that conservative failure with full cursor traversal while preserving the same cache model and clear failure behavior for genuinely unsupported cases.

Past decision: GitHub Projects remains the remote source of truth and `.local/pj/cache.json` remains derived data optimized for AI reads. The implementation should therefore improve remote reads without introducing a second local state model.

Past decision: tool output and remote payloads should be trimmed at the source. Pagination should continue requesting only fields needed by the current cache and commands instead of expanding the GraphQL selection set.

## Design Options

### Option A: Keep one large project query and add `after` variables

- Fetch fields and items in the same query, then repeat while either connection has another page.
- Practical issue: independently paginating two top-level connections in one query makes request control awkward, because field and item cursors can advance on different schedules.
- This keeps the visible query shape close to today's implementation, but risks making truncation and test cases harder to reason about.

### Option B: Split sync into small paginated loaders

- Resolve the project shell once, then load fields, items, and item field values through focused helper queries.
- Each helper owns one cursor loop and one response shape.
- This creates a few more HTTP requests but keeps correctness and tests straightforward.

Recommended option: Option B. It aligns with the workspace preference for correctness over speed, keeps each GraphQL response compact, and makes it easy to preserve clear errors if GitHub returns an unexpected page shape.

## Spec Changes

### `docs/specs/github-projects-task-cli.md`

- Change the `pj sync` contract from failing on field/item single-page overflow to requiring cursor-based pagination for project fields and project items.
- Define item field value pagination behavior:
  - either implement pagination for `fieldValues` in the same execution slice, or
  - keep a clear per-item failure if GitHub reports more field values than the implementation supports.
- Remove or revise the current error-handling requirement that query results exceeding single-page limits fail because pagination is not implemented.
- Keep `.local/pj/cache.json` as derived data and keep the normalized item model unchanged.

## Code Changes

### `tools/pj/internal/pj/github.go`

- Refactor `syncProject` so project identity, fields, items, and field values can be loaded through cursor-aware helpers.
- Add GraphQL queries with `pageInfo { hasNextPage endCursor }` and `after` variables for:
  - project fields
  - project items
  - item field values, if included in this slice
- Preserve existing field parsing for single-select options and item normalization for draft issues, issues, and pull requests.
- Remove the field/item `hasNextPage` hard failures after full pagination exists.
- Keep owner project discovery pagination out of this plan unless it becomes a blocker; this issue is specifically about sync data completeness for the configured board.

### `tools/pj/internal/pj/github_test.go`

- Add tests that serve multiple pages of project fields and assert the final cache contains fields from every page.
- Add tests that serve multiple pages of project items and assert the final cache contains items from every page.
- Add tests for item field value pagination if implemented in this slice.
- Keep or adapt existing tests around project lookup, provisioning, and mutation behavior so the pagination refactor does not weaken bootstrap coverage.

### `tools/pj/internal/pj/types.go`

- Update shared response structs only if the pagination refactor needs reusable typed page containers. Avoid broad type churn if local helper structs stay clearer.

## Design Decisions

- No ADR update is expected for this plan. The architectural direction remains the existing GitHub Projects-backed CLI with derived local cache; pagination is an implementation hardening step.
- Prefer a small helper-query refactor over a large GraphQL abstraction. The CLI is still a workspace-local spike, and a general GraphQL pagination framework would be premature unless the implementation naturally repeats enough code to justify it.
- Treat item field value pagination as part of execution discovery: implement it in this slice if the helper-query refactor makes it small; otherwise preserve a clear failure for oversized item field values and log a narrower follow-up issue.

## Sub-tasks

- [x] Update `docs/specs/github-projects-task-cli.md` to require cursor-based pagination for fields and items and to clarify item field value handling.
- [x] [parallel] Review the current `syncProject` response parsing and identify the smallest helper boundaries for field, item, and field-value loading.
- [x] [parallel] Add failing `httptest` coverage for multi-page field and item responses.
- [x] [depends on: helper boundary review, multi-page tests] Refactor `syncProject` into cursor-aware loaders while preserving the existing cache shape.
- [x] [depends on: cursor-aware loaders] Decide and implement item field value pagination or retain a clear oversized-field-values error with a follow-up issue.
- [x] [depends on: pagination implementation] Run `go -C tools/pj test ./...`.
- [x] [depends on: tests] Run `go -C tools/pj run ./cmd/pj sync` against the configured workspace board when credentials and project config are available.
- [x] [depends on: verification] Move `docs/issues/pj-pagination-follow-up.md` to `docs/issues/done/` during execution if the issue is fully resolved.

## Parallelism

- The spec update, helper-boundary review, and test design can proceed independently.
- The implementation depends on agreeing on helper boundaries and should happen after the failing test cases define expected pagination behavior.
- Live `pj sync` verification depends on tests passing and on local GitHub credentials/config being available.

## Verification

- PASS: `go -C tools/pj test ./...`
- PASS: `go -C tools/pj run ./cmd/pj sync --owner yoskeoka --owner-type user --project 1 --config /tmp/pj-pagination-config.json --cache /tmp/pj-pagination-cache.json`
- Inspect `.local/pj/cache.json` after sync only as derived verification output; do not commit it.
- Confirm the spec no longer promises failure for ordinary field/item pagination overflow once implementation supports it.

## Expected Outcome

- `pj sync` can cache boards larger than the current one-page field and item limits.
- Cache contents remain normalized exactly as before, so `pj list`, `pj add`, and `pj update` keep their current contracts.
- Oversized item field values are either fully handled or fail with a narrower, explicit message and a recorded follow-up issue.
