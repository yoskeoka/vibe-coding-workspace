# KB Bilingual Rendering

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Extend the workspace knowledge base so the published wiki supports English at `/` and Japanese at `/ja/` while keeping English as the canonical AI-facing corpus for QA and retrieval.

This plan covers:

- adding bilingual KB source trees with matching path structure under `docs/kb/` and `docs/kb/ja/`
- publishing both locales through MkDocs with a shared navigation model and language switcher
- keeping `docs/kb/ja/**` out of QA and retrieval flows
- recording original source language in source-note front matter

## Design Decisions

Past belief: `AI-First` favors structures that are easy for LLMs to search and exclude deliberately. Apply the same reasoning here by isolating Japanese content under `docs/kb/ja/` so retrieval can ignore that subtree without heuristic filtering.

Past decision: the KB lives in-repo under `docs/kb/` and publishes from the same Markdown source of truth. Apply the same reasoning here by keeping bilingual Markdown in-repo and extending the existing derived MkDocs build rather than introducing a separate renderer stack.

Considered options:

- `docs/kb/` English only, Japanese wiki pages only: simpler duplication story, but source-note browsing becomes inconsistent between locales and the published Japanese wiki lacks the source-oriented layer.
- `docs/kb/` canonical English plus full `docs/kb/ja/` mirror for wiki and source notes: slightly more ingest work, but consistent rendered IA, simple path mapping, and straightforward retrieval exclusion.
- move the human-facing site to a different framework such as Astro: more control over navigation and performance, but it increases stack complexity before the KB publishing model is settled.

Recommended option confirmed in discussion: keep canonical English under `docs/kb/`, add a Japanese mirror under `docs/kb/ja/` for wiki and source notes, continue using MkDocs Material, and disable `navigation.instant` to stay compatible with static i18n routing.

No ADR update is expected unless implementation reveals a materially different renderer architecture or retrieval contract.

## Code Changes

### `tools/kb`

- Add multilingual build support on top of the existing generated-config flow.
- Install and invoke the static i18n plugin through the KB build path.
- Keep `check`, `build`, and `serve` working against invocation-owned generated roots.

### `tools/kb_generate.py`

- Generate publish inputs for both the canonical English tree and the Japanese mirror tree.
- Preserve identical relative paths between locales so `/foo/` maps cleanly to `/ja/foo/`.
- Generate locale-aware `Sources` navigation and yearly index pages for both languages.
- Keep frontmatter-derived visible sections working within each locale tree.

### `mkdocs.kb.template.yml`

- Add the static i18n plugin and locale configuration for English and Japanese.
- Remove `navigation.instant` because it conflicts with Material's locale switching flow in this setup.
- Add any Material language-switcher configuration needed for `/` and `/ja/`.

### `docs/kb/`

- Keep English wiki and source notes as the canonical QA-facing corpus.
- Add or update frontmatter conventions so source notes record `original_language`.

### `docs/kb/ja/`

- Introduce a Japanese mirror layout for:
  - `wiki/`
  - `sources/YYYY/`
- Treat this tree as human-facing published content and an ingest-time translation target.

### QA / retrieval integration

- Update the relevant KB query/retrieval path so `docs/kb/ja/**` is excluded from AI-facing search and question answering.
- Keep the English canonical corpus as the only retrieval input unless a future spec explicitly changes that contract.

## Spec Changes

### `docs/specs/knowledge-base.md`

- Define the bilingual KB layout and the role split between canonical English content and Japanese published mirror content.
- Define that source notes and wiki pages may both exist under `docs/kb/ja/` with matching path structure.
- Define that QA and retrieval operate on canonical English content and exclude `docs/kb/ja/**`.
- Define the MkDocs i18n publishing contract for `/` and `/ja/`.
- Record that `navigation.instant` is disabled for compatibility with the chosen i18n approach.
- Define `original_language` as required source-note metadata.

### `docs/kb/schema.md`

- Extend the directory rules and frontmatter conventions to cover `docs/kb/ja/`.
- Document `original_language` in source-note front matter.
- Clarify that Japanese source notes and wiki pages mirror the English path layout for publishing, while English remains canonical for AI retrieval.

### `docs/kb/ingest.md`

- Define bilingual ingest expectations: English canonical note/page creation first, Japanese mirror creation as a best-effort companion output.
- Clarify that ingest should update both locale trees when practical, while retrieval-facing workflows continue to rely on English only.

## Sub-tasks

- [ ] Update KB specs and schema for bilingual layout, retrieval scoping, and `original_language`
- [ ] Implement MkDocs static i18n configuration and remove `navigation.instant`
- [ ] [parallel] Extend KB input generation for locale-aware nav, yearly source indexes, and per-locale visible relationship sections
- [ ] [parallel] Add `docs/kb/ja/` structure and seed mirror content needed to validate the bilingual publishing path
- [ ] [depends on: spec updates, MkDocs static i18n configuration, locale-aware generation] Wire KB build and serve commands to render `/` and `/ja/` correctly
- [ ] [depends on: spec updates] Update ingest/retrieval flows so Japanese files are generated as a companion output but excluded from QA/retrieval inputs
- [ ] [depends on: all implementation tasks] Verify with `tools/kb check` and `tools/kb build`

## Verification

- `tools/kb check`
- `tools/kb build`
- Inspect the generated site to confirm:
  - English pages render at `/`
  - Japanese pages render at `/ja/`
  - locale switching preserves path shape where both locales exist
  - `docs/kb/ja/**` is excluded from QA/retrieval inputs
