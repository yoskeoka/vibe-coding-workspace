---
title: Godot と Codex で 3 つの題材を試してみたら、初回出力の強さと人手確認の重さが見えてきた
source_url: https://dev.classmethod.jp/articles/godot-codex-three-theme-verification/
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - ai-game-dev
  - godot
  - codex
  - agentic-dev
published_on: 2026-03-18
named_entities:
  - Godot
  - Codex
  - GPT-5.4
  - EditorPlugin
  - GraphEdit
  - Multiplayer Bomber Demo
related_pages:
  - ../../wiki/topics/ai-game-dev.md
  - ../../wiki/tools/godot.md
---

# Godot and Codex across three game-development tasks

## Why it matters here

This is one of the most concrete comparisons in the KB for what an AI coding agent can do on first pass versus what still needs a human in game development.

## Summary

- The article tests `Codex` with `Godot 4.6.1` across three separate themes: fast 2D action, small-scale multiplayer, and an `EditorPlugin`-based creation tool.
- Its headline conclusion is consistent across all three: the first pass can build a useful base quickly, but detailed tuning, feel, editor usability, and final validation still need human review.
- The author deliberately isolates each theme into its own session and directory so cross-task context does not contaminate the comparison.
- The concrete test targets matter: movement feel, wall actions, coyote time, and camera behavior for action; ownership and sync correctness for multiplayer; and practical editor UX for the plugin path.
- The article positions Godot as especially useful for this kind of experiment because runtime code and editor extensions can both be exercised in the same environment.

## Workspace takeaways

- Separate sessions by theme when comparing agent performance across substantially different problem classes.
- Judge success by the hard-to-fake surfaces: movement feel, collision correctness, sync integrity, and editor workflow roughness.
- For future AI-assisted Godot work, preserve both the initial prompt and the human-found deviations. The delta is the real benchmark.

