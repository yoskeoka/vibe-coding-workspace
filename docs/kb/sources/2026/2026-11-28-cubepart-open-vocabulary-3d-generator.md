---
title: "CubePart: An Open-Vocabulary Part-Controllable 3D Generator"
source_url: https://cubepart.github.io/
source_type: docs
original_language: en
ingested_on: 2026-11-28
status: active
tags:
  - 3d-generation
  - game-development
  - research
related_pages:
  - ../../wiki/topics/ai-game-dev.md
---

# CubePart: An Open-Vocabulary Part-Controllable 3D Generator

## Summary

CubePart is a research system for generating or decomposing multi-part 3D meshes from a text prompt and a user-defined, open-vocabulary parts schema. It produces a mesh for each named part and aims to keep the pieces structurally complete and assembled as a coherent object, so they can support animation, physics, or scripted behavior.

The project describes a two-stage model: a schema-aware stage generates the holistic shape, then a multi-part decoder produces individual part latents with cross-part attention. Its part-labeled dataset is reported as 462K assets and 2.02M parts. The project page includes a gallery of interactive meshes and behavior videos, plus comparisons against other decomposition approaches.

## Workspace relevance

This is a useful research anchor for game assets that need meaningful editable parts, such as a robot with independent limbs or a chest with a separate lid. It connects text-to-3D asset generation to downstream animation and game logic. Treat the results as research claims and inspect the linked paper, code, and demo before assuming a production-ready asset pipeline.

## Source

- Project page: https://cubepart.github.io/
- Paper and code links are provided on the project page.
- The authors list SIGGRAPH 2026 and a CC BY 4.0 license for the project-page content.
