# Follow Up: Prevent `tools/kb` Races on Shared Generated Paths

## Summary

`tools/kb check` and `tools/kb build` both write to the shared generated workspace under `.local/kb-generated`.

When they were run in parallel, one process removed or replaced generated inputs while the other was still using them, which produced failures such as:

- `OSError: [Errno 39] Directory not empty: '.local/kb-generated'`
- MkDocs config errors because `.local/kb-generated/docs` no longer existed at the moment of use

Running the commands sequentially succeeds.

## Why It Matters

- The current behavior is easy to trigger from parallel agent/tool execution.
- The failures are noisy and look like content or MkDocs problems even though the underlying issue is a shared scratch-path race.
- This makes KB verification less robust in automated or agentic workflows.

## Proposed Solution

- Make generated paths invocation-specific instead of globally shared, or add locking around generation/build steps.
- Ensure `tools/kb` reports a clearer error when a concurrent run is detected.
- Decide whether `.local/kb-generated` should be a cacheable shared workspace or a per-run temporary directory.
