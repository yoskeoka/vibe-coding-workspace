# Follow Up: Add Pagination Support to `tools/pj`

## Summary

`tools/pj` currently queries GitHub Projects with fixed page sizes:

- fields: 50
- items: 100
- item field values: 20

The current spike now detects truncation at the project field and item list level and fails clearly, but it still does not implement full cursor-based pagination.

## Why It Matters

- Larger boards would produce incomplete caches without pagination
- Silent truncation would be dangerous for triage, so the current spike must remain conservative
- A reusable task CLI should handle normal board growth without hard operational limits

## Proposed Solution

- Implement cursor-based pagination using `pageInfo { hasNextPage endCursor }`
- Continue fetching until all project fields and items are loaded
- Decide whether item field values also need pagination support in the same pass
- Add verification against a board shape that exceeds the current single-page limits
