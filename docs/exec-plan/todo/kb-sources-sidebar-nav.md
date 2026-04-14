# Knowledge Base Sources Sidebar Navigation

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Fix the rendered knowledge-base navigation so ingested source notes do not disappear from the sidebar when new files are added, while keeping the `Sources` section compact enough for humans to browse on the web.

This plan covers three user-facing improvements:

- group source notes by year in the sidebar with labels like `2026 (6)`
- add year landing pages so the sidebar links stay compact even as source count grows
- expose `sources` and `related_pages` metadata as visible links in page content instead of hiding those relationships only in YAML frontmatter

## Code Changes

### `tools/kb`
- Add a nav-generation step or helper that derives the `Sources` navigation from files under `docs/kb/sources/YYYY/`.
- Ensure the generated nav labels include per-year counts and stay deterministic.
- Generate publish-only year index pages and a generated MkDocs config in a temporary or local-only workspace, not as committed source files.
- Keep `build` strict and fail fast if generated KB navigation inputs are stale or missing.

### `mkdocs.kb.template.yml`
- Replace the hand-maintained `mkdocs.kb.yml` with a template file that keeps the stable site configuration and fixed wiki navigation.
- Leave the `Sources` section to be injected by the KB generation step so newly added source notes cannot be omitted accidentally.
- Remove always-expanded source navigation if it prevents the year grouping from collapsing.

### Generated KB build inputs
- Generate one landing page per year that lists source-note dates and titles.
- Keep these year index pages out of git and generate them only for build, check, and publish flows.
- Generate the effective MkDocs config as a derived file such as `.local/mkdocs.kb.generated.yml` and pass it to `mkdocs -f`.

### KB rendering hooks or templates
- Add a small rendering extension so wiki pages show their `sources:` links in page content.
- Add a matching visible section for source notes so `related_pages:` links are visible in the rendered site.
- Keep source-of-truth metadata in frontmatter; the rendered links should be derived from it rather than duplicated manually.

## Spec Changes

### `docs/specs/knowledge-base.md`
- Define that `Sources` navigation is grouped by year rather than a flat manually curated list.
- Define that each `sources/YYYY/` directory has a human-facing yearly index page in the rendered site, but that these pages are derived artifacts rather than committed source files.
- Define that source relationships stored in frontmatter must also be visible in rendered page content for web readers.
- Define that the KB site uses a template config plus generated nav instead of a hand-maintained full MkDocs nav file.

### `docs/kb/schema.md`
- Document the year index page convention as a generated publish artifact.
- Document that `sources` and `related_pages` metadata are rendered into visible link sections in the site output.

### `docs/kb/ingest.md`
- Clarify that yearly index pages and source nav are generated automatically during build/check rather than maintained during ingest.

## Design Decisions

Past decision: the KB is stored in-repo under `docs/kb/` so the same Markdown files remain useful to both AI agents and humans. Apply the same reasoning here: navigation should still be derived from repo-native Markdown instead of introducing a separate external catalog.

## Sub-tasks

- [ ] Add spec updates for grouped `Sources` navigation and visible metadata-derived link sections
- [ ] Implement template-based KB config generation, grouped source nav generation, and publish-only yearly source index pages
- [ ] Implement rendering support for visible `sources` / `related_pages` link sections
- [ ] Update existing KB pages or generated outputs as needed for the new navigation structure
- [ ] Verify with `tools/kb check` and `tools/kb build`
