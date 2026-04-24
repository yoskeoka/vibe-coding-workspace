---
title: お父さん、ゲーム作れるの？
source_url: https://zenn.dev/sonicmoov/articles/76a5098905c978
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - ai-game-dev
  - next2d
  - webgl
  - webgpu
  - engine-dev
published_on: 2026-04-10
named_entities:
  - Next2D
  - swf2js
  - OffscreenCanvas
  - Transferable
  - ArrayBuffer
  - WebGL
  - WebGPU
  - MVVM
  - Clean Architecture
  - Atomic Design
  - MCP
related_pages:
  - ../../wiki/topics/ai-game-dev.md
  - ../../wiki/tools/next2d.md
---

# Next2D origin story and engine/tooling lessons

## Why it matters here

This article is both a motivation note and a technical anchor for owning more of the game-production stack. It combines engine design, performance engineering, production tooling, and AI-era workflow integration.

## Summary

- The article frames `Next2D` as a response to a simple personal question: whether the author could truly make games directly, not only inside a large-team production context.
- `Next2D` is described as three parts: the `Player` runtime, a `Framework` grounded in `MVVM + Clean Architecture + Atomic Design`, and a browser-based `AnimationTool`.
- A major technical focus is maintaining stable `60fps` by pushing rendering to a worker with `OffscreenCanvas`, keeping logic on the main thread, and using `Transferable` `ArrayBuffer` ownership transfer instead of repeated copies.
- The article stresses memory discipline, buffer reuse, object pooling, and the practical difficulty of keeping asynchronous main-thread and sub-thread work smooth without GC spikes.
- `WebGPU` support is described not as a clean rewrite but as a careful compatibility effort around an existing `WebGL`-shaped command stream.
- It also notes a `Next2D` `MCP` server released in 2026, tying the engine/toolchain story directly into AI-assisted development.

## Workspace takeaways

- Engine ownership only pays off if tooling, rendering architecture, and content-production flow are all designed together.
- `OffscreenCanvas`, `Transferable`, buffer reuse, and object pools are the durable retrieval anchors here; they describe the real performance shape of the system.
- The `MCP` reference makes this relevant not only as engine engineering but as a context-rich environment for future AI-assisted workflows.

