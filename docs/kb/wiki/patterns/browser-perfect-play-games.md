---
title: Browser-Playable Perfect-Play Board Games
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2021/2021-12-30-mametter-6x6-reversi-god.md
  - ../../sources/2023/2023-12-10-y-tetsu-variant-board-reversi-strong-solve.md
---

# Browser-Playable Perfect-Play Board Games

## Pattern

If a small perfect-information board game should feel immediate in the browser, keep the solver compact, phase-aware, and isolated from the rendering layer. The browser runtime should spend time on interaction, not on recomputing known opening structure or paying avoidable representation costs.

## Concrete anchors

- Represent the board in a form that copies cheaply and supports bitwise move generation, such as two `u64` bitboards for a 6x6 Reversi board.
- Split the solve path by game phase: exact late-game search, heuristic-assisted midgame ordering, and a precomputed opening tree for the earliest plies.
- Use heuristics to accelerate exact search, not to replace it, when the product claim depends on perfect play.
- Compress precomputed opening data aggressively enough to ship it with the client, for example with `LOUDS` plus direct binary embedding.
- Run the solver in `WebAssembly` and a `Web worker`, then keep UI technology decisions separate from solver-language decisions.
- Treat integer width and target architecture as part of the product contract. A native 64-bit assumption can fail silently when the browser target is effectively 32-bit.
- When exploring board-shape variants, first check whether the engine can absorb a board-mask extension at the legal-move boundary instead of rewriting the move generator.
- Keep board editing and replay tooling outside the core solver so custom shapes and solved records can iterate faster than the search engine itself.

## Trade-offs

- Exact or near-exact search pushes a lot of complexity into data layout, search ordering, and precomputation.
- The fastest internal data structures may justify `unsafe` or similarly sharp tools, which raises maintenance cost.
- Browser packaging can expose a different class of bugs than native execution, especially around integer width, worker boundaries, and bundler configuration.
- Reusing a native engine accelerates experimentation, but it also drags in that engine's build assumptions, distribution format, and platform-specific setup.
- A perfect-play experience is compelling, but it narrows design freedom compared with fuzzier, more style-driven opponents.

## Related pages

- [reversi-adventure](../projects/reversi-adventure.md)
