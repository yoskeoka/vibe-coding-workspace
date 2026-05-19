---
title: 〖リバーシ(オセロ)〗変形ボードの強解決！
source_url: https://qiita.com/y-tetsu/items/6e91ceb110951d799704
source_type: article
original_language: ja
ingested_on: 2026-05-19
status: active
tags:
  - reversi
  - othello
  - edax
  - exact-search
  - board-editors
published_on: 2023-12-10
named_entities:
  - Edax
  - Othello is Solved
  - setboard
  - MinGW
related_pages:
  - ../../wiki/projects/reversi-adventure.md
  - ../../wiki/patterns/browser-perfect-play-games.md
---

# Strong-solving variant-board Reversi with minimal Edax surgery

## Why it matters here

This is a useful bridge between pure solver work and playable board-game tooling. It shows that a strong existing engine can be repurposed for custom board shapes with a very small rules-layer patch, then wrapped with lightweight authoring and replay workflows instead of rebuilding the whole search stack from scratch.

## Summary

- The article adapts `Edax` to strongly solve small Reversi boards with non-standard shapes and initial placements.
- The key representation change is a third 64-bit mask `H` for holes, alongside the usual black and white bitboards.
- The move generator stays mostly intact. The author simply removes `H` positions from legal moves by changing the final empty-square mask from `~(P|O)` to `~(P|O|H)`.
- The engine is configured for full-depth search up to 60 plies and disables the standard 8x8 opening book so custom boards are solved from search rather than inherited book data.
- The workflow adds `setboard` for custom 64-cell board strings, a `display` command for move-record output, and a separate browser board editor / replay site to author shapes and replay solved records.

## Workspace takeaways

- If the goal is to prototype a strong opponent or solve custom board variants quickly, patching a proven engine can beat building a fresh solver.
- A board-mask abstraction is often enough to unlock variant support when the underlying move-generation core is already reliable.
- Separating concerns matters here too: native engine changes stay minimal, while board editing and record replay live in auxiliary tools.
- This is complementary to the existing Rust/WASM 6x6 reference. One note shows how to ship a compact solver in the browser; this one shows how to validate custom-rule or custom-shape ideas earlier with an adapted native engine.
- The seven sample shapes are also a practical reminder that "solvable" and "pleasant to play" are different filters. Some boards finish quickly while move-rich shapes such as the ring still take around two minutes to analyze.

## Concrete anchors

- Engine baseline: `Edax`
- Strong-solve settings: search depth raised from `21` to `60`, book usage disabled
- Variant-board representation: third bitboard `H` for holes
- Legal-move patch: final mask changes to `~(P|O|H)`
- Authoring interface: `setboard "<64 cells><side-to-move>"`, with spaces representing holes
- Replay support: `display` command prints the solved move record
- Sample shapes: `4x4`, `十字`, `風車`, `卍`, `X`, `外側に石`, `リング`

## Follow-up

- If `reversi-adventure` needs to test unusual board layouts or authored puzzle boards, compare a fast native-solver prototyping loop like this against direct browser implementation.
- If a future board-game experiment reuses an external engine, keep the variant support patch as close to the legal-move boundary as possible so solver correctness remains easy to reason about.
