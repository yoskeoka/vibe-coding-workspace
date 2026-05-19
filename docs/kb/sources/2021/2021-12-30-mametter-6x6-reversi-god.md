---
title: 6x6リバーシの神
source_url: https://mametter.hatenablog.com/entry/2021/12/30/114111
source_type: post
original_language: ja
ingested_on: 2026-05-19
status: active
tags:
  - reversi
  - rust
  - webassembly
  - browser-games
  - exact-search
published_on: 2021-12-30
named_entities:
  - Rust
  - WebAssembly
  - TypeScript
  - three.js
  - wasm-bindgen
  - wasm-pack
  - Web worker
  - MTD(f)
  - LOUDS
  - Joel Feinstein
related_pages:
  - ../../wiki/projects/reversi-adventure.md
  - ../../wiki/patterns/browser-perfect-play-games.md
---

# God-tier 6x6 Reversi in the browser

## Why it matters here

This is a strong reference for any workspace game that wants exact or near-exact board-play quality in a browser without relying on a server. It is especially relevant to `reversi-adventure` because it shows how to turn a solved small-board game into a playable, intentionally unwinnable experience.

## Summary

- The post describes a 6x6 Reversi AI that plays perfect defense as white, making the human black player unable to win.
- The core board representation is a bitboard implemented as two Rust `u64` values, one for black stones and one for white stones.
- The solver is split by game phase: exact negamax for the endgame, `MTD(f)` with a transposition table and heuristic move ordering for the midgame, and a precomputed opening tree for the first 16 plies.
- The opening tree is compressed with `LOUDS` and embedded directly into the binary with `include_bytes!`, keeping the payload practical for browser delivery.
- The browser build uses `wasm-bindgen`, `wasm-pack`, and a `Web worker`, while the UI is implemented separately in TypeScript and `three.js`.

## Workspace takeaways

- A browser-playable perfect-information game does not need a server if the board model is compact enough and the expensive search is phased carefully.
- The article is a concrete example of mixing exact search with heuristic acceleration without giving up the final exactness claim. The heuristic evaluator only improves move ordering until the solver can switch back to exact search.
- Precomputing the opening is not just a speed hack. It is a delivery strategy that lets the browser spend runtime budget on the interactive part instead of repeating known work.
- The post is also a reminder that platform changes matter: the solver worked natively, then failed in WASM because `usize` assumptions changed under a 32-bit target.
- UI and solver separation is deliberate here. The author kept the search engine in Rust/WASM and chose TypeScript plus `three.js` for the presentation layer instead of forcing everything through one language stack.

## Concrete anchors

- Board representation: two `u64` bitboards and macro-generated flip code specialized for all 36 squares
- Endgame solver: exact negamax once empty squares drop to 11 or fewer
- Midgame search: `MTD(f)` plus move ordering from hand-built features such as `edge2x` patterns and legal-move counts
- Learned evaluation weights: roughly 3000 floating-point values embedded as Rust code
- Opening strategy: about 16 plies of precomputed best-response data encoded with `LOUDS`
- Browser packaging: `wasm-bindgen`, `wasm-pack`, `Web worker`, TypeScript, `three.js`
- Portability bug: native 64-bit success masked a 32-bit WASM `usize` overflow

## Follow-up

- If `reversi-adventure` wants a deterministic strong-opponent mode, compare its current board representation and search split against this `bitboard -> phased solver -> browser worker` shape.
- If a future workspace board game aims for client-only hosting, preserve the target-width assumptions explicitly in tests before shipping to WASM.
