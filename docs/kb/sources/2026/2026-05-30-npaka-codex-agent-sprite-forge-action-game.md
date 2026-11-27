---
title: "Trying action game generation with Codex and Agent Sprite Forge"
source_url: https://note.com/npaka/n/ne0c0daafbfa3
source_type: article
original_language: ja
ingested_on: 2026-11-28
status: active
tags:
  - agent-sprite-forge
  - codex
  - phaser
  - game-development
related_pages:
  - ../../wiki/tools/agent-sprite-forge.md
  - ../../wiki/topics/ai-game-dev.md
---

# Trying action game generation with Codex and Agent Sprite Forge

## Summary

npaka documents an iterative workflow for building a side-scrolling action game with Codex, Agent Sprite Forge, Vite, TypeScript, and Phaser. The initial request uses two reference images: a target game screenshot and a character/stage overview. Codex creates a 16:9 game with idle and run animations, then Agent Sprite Forge is used to revise player sprites from a character reference.

The prototype is refined in focused passes: tile spacing and jump reach, damage and restart behavior, enemy direction and stomp interactions, stage-clear/game-over presentation, UI scale and typography, and a map editor with save/load behavior. The author also adjusts a looping background in Photoshop with Generative Fill and adds separately sourced BGM and sound effects to a project `sounds` folder.

## Reusable workflow

- Provide concrete visual references and name each file's role in the prompt.
- Specify the runtime stack and presentation constraints up front: `Vite + TypeScript + Phaser` and 16:9 in this example.
- Iterate through coherent game concerns instead of one broad revision: movement and level geometry, enemy behavior, player animation, UI, editor, and audio.
- Treat generated output as a prototype that needs playtesting. The article reports issues such as impossible first gaps, inconsistent editor/game tiles, and controls breaking after saving.
- For cyclic walk animation, inspect repeated limb poses. The author observes that generated horizontal walks often repeat the same leading foot, and suggests using a hand-corrected motion image as a reference or asking Agent Sprite Forge to perform post-processing only.
- Keep asset production connected to runtime use: adjust sprite sheets, background loops, editor tiles, and audio in the context of the game.

## Workspace relevance

This is a concrete end-to-end example for the existing Agent Sprite Forge and AI game-development notes. Its main useful anchor is the iteration loop from image-guided Phaser generation through gameplay correction and asset cleanup. The map editor section is also a reminder to verify that editor output matches the runtime representation and remains usable after saving.

## Source

- Original Japanese article: https://note.com/npaka/n/ne0c0daafbfa3
- Author: npaka
- Published: 2026-05-30
