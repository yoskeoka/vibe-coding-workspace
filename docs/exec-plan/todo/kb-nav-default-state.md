# KB Navigation Default State

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Refine the published KB navigation so the top-level information architecture matches how the rendered site is actually used by humans:

- make the wiki pages the effective top-level entry points instead of showing `Home`, `Schema`, and `Ingest` first
- keep `Sources` compact by opening year groups but keeping individual source-note entries collapsed by default

This is a follow-up to `kb-sources-sidebar-nav.md`, which fixed missing source-note nav entries but left the overall default open/closed state too blunt.

## Code Changes

### `mkdocs.kb.template.yml`
- Move the current `Wiki` children up one level in the nav so the top page becomes `wiki/index.md`.
- Keep `Sources` as a top-level nav group below the wiki-oriented entries.
- Add any theme hooks, extra JS, or extra CSS references needed for a KB-specific initial sidebar state.

### KB theme override assets
- Add a small KB-only customization that adjusts the initial sidebar state after Material renders navigation.
- Open `Sources` year groups by default while leaving individual source-note entries collapsed.
- Preserve normal navigation behavior after initial load so the user can still manually expand or collapse anything.

### `tools/kb_generate.py`
- Adjust generated `Sources` nav structure if needed so year groups and individual source-note leaves are distinguishable by the customization layer.
- Keep generated source indexes and visible frontmatter-derived sections working with the new top-level nav layout.

### `tools/kb`
- No major behavior change expected beyond continuing to build the generated config and any new theme assets correctly.

## Spec Changes

### `docs/specs/knowledge-base.md`
- Define that the published KB should use wiki content as the primary top-level browsing surface.
- Define the desired default navigation behavior for `Sources`: year groups visible, individual source notes collapsed until opened.
- Document any KB-specific theme customization used to achieve that behavior.

### `docs/kb/schema.md`
- Clarify that `README.md`, `schema.md`, and `ingest.md` remain source files but are not required to lead the published top-level nav.

## Design Decisions

Past belief: `AI-First` favors repo-native Markdown structure over deep human-oriented hierarchies. Apply the same reasoning here by keeping the on-disk files unchanged while adjusting only the rendered site navigation for human browsing.

Past decision: the KB lives in-repo and publishes from the same Markdown source of truth. Apply the same reasoning here by keeping navigation refinements inside the KB publishing layer instead of introducing a separate curated site tree in git.

## Sub-tasks

- [ ] Update KB specs for the new top-level nav layout and default `Sources` open/closed behavior
- [ ] Restructure the published top-level nav so wiki pages move up one level and `Sources` stays below them
- [ ] Add KB-specific navigation customization so `Sources` year groups start open while individual source-note entries start collapsed
- [ ] Verify with `tools/kb check` and `tools/kb build`
