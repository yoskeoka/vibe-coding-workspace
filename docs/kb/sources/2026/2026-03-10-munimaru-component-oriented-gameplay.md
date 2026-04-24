---
title: なぜ私はゲーム開発で「疎結合」と「コンポーネント指向」に異常なほどこだわるのか
source_url: https://zenn.dev/munimaru62o/articles/c6ed730c6e4c61
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - ai-game-dev
  - architecture
  - unreal-engine
  - component-design
published_on: 2026-03-10
named_entities:
  - Unreal Engine
  - UE5
  - GameCoreFramework
  - GameplayTag
  - DataAsset
  - ScriptableObject
related_pages:
  - ../../wiki/topics/ai-game-dev.md
  - ../../wiki/patterns/component-oriented-gameplay.md
---

# Munimaru on component-oriented gameplay architecture

## Why it matters here

This article is a strong design reference for game-system architecture that aims at emergent play rather than scripted one-to-one interactions. It is useful even outside Unreal because the core framing is engine-agnostic.

## Summary

- The article argues that strict domain separation and loose coupling are not just code-quality goals; they are the basis for emergent gameplay where systems can combine in unexpected ways.
- It frames the central problem as escaping fixed input-output patterns that only replay the developer's pre-authored solutions.
- Three architectural supports are emphasized: a safety net that does not break when systems are recombined, a common intent protocol between components, and a data-driven workflow that lets non-programmers explore combinations.
- Concrete Unreal-flavored anchors include separating persistent and temporary state, leaning on `GameplayTag`-style intent passing, and using `DataAsset` / `ScriptableObject`-like assets for designer-facing composition.
- The article is explicit about trade-offs: balance complexity grows fast, cognitive load increases, and debugging and performance can become harder.

## Workspace takeaways

- When a game plan claims to support combinatorial mechanics, it should say what the safety net is, what the intent protocol is, and which parts are data-driven.
- "Loose coupling" is not enough as a slogan. The durable retrieval anchors are safety, protocol, and authoring workflow.
- This is a useful comparison point for any future engine-specific framework work in the workspace.

