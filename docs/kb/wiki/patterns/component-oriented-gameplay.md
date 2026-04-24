---
title: Component-Oriented Gameplay
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-10-munimaru-component-oriented-gameplay.md
---

# Component-Oriented Gameplay

## Pattern

If a game wants emergent combinations instead of fixed authored interactions, it needs more than loosely coupled code. It needs explicit safety boundaries, an intent protocol between systems, and a data-driven way to test combinations without constant engineering intervention.

## Concrete anchors

- Safety net: persistent and temporary state are separated so possession, swapping, or dynamic recombination does not corrupt the underlying player state.
- Intent protocol: systems communicate through tags, interfaces, or equivalent intent messages rather than direct knowledge of each other's concrete types.
- Designer-facing composition: new combinations should be testable through assets and configuration, not by touching engine code every time.

## Trade-offs

- Combinatorial balance gets harder quickly.
- Debugging and performance can become more difficult.
- Cognitive load rises unless the framework makes setup and lifecycle ordering safe by default.

## Related pages

- [ai-game-dev](../topics/ai-game-dev.md)

