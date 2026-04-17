# Follow Up: Revisit KB Renderer Pinning and Long-Term MkDocs Choice

## Summary

The workspace KB now pins the immediate renderer versions used during this ingest:

- `mkdocs==1.6.1`
- `mkdocs-material==9.7.6`

During `tools/kb build`, the current `mkdocs-material` release emits an upstream warning about the future `MkDocs 2.0` direction. The build still succeeds today, and the rendered wiki is acceptable, but the warning remains a reminder that the KB publishing stack has an unresolved long-term dependency decision.

## Why It Matters

- The current setup is working, and the rendered wiki is acceptable, so there is no immediate user-facing problem to fix.
- Exact pins reduce short-term drift, but they do not answer whether `MkDocs + Material for MkDocs` is the right long-term renderer.
- The upstream warning raises a decision that should be made deliberately instead of by drift: should the workspace continue with `MkDocs + Material for MkDocs`, or switch to a different renderer later if upstream compatibility or governance becomes a concern?

## Proposed Solution

- Document the current rationale for staying with MkDocs Material if the answer is "keep it".
- Decide whether exact pins in `requirements-kb.txt` are enough, or whether `tools/kb` should use a lockfile-based dependency path.
- If the team wants a lower-risk long-term path, evaluate alternatives against the current requirements:
  - git-tracked Markdown as source of truth
  - lightweight local build/serve flow
  - GitHub Pages publishing
  - good sidebar/navigation support for the compiled wiki
- Record an explicit decision in `docs/design-decisions/adr.md` once the trade-off is reviewed.
