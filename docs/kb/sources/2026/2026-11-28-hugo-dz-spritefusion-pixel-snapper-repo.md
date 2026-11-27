---
title: Sprite Fusion Pixel Snapper repository
source_url: https://github.com/Hugo-Dz/spritefusion-pixel-snapper
source_type: repo
original_language: en
ingested_on: 2026-11-28
status: active
tags:
  - pixel-art
  - image-processing
  - tooling
related_pages:
  - ../../wiki/topics/pixel-art-ui.md
---

# Sprite Fusion Pixel Snapper repository

## Summary

`Hugo-Dz/spritefusion-pixel-snapper` is the open-source Rust implementation behind Sprite Fusion Pixel Snapper, a post-processing tool for AI-generated and procedural pixel art. It snaps inconsistent pixels to a regular grid, maintains a consistent grid resolution, and quantizes colors to a limited palette while aiming to preserve detail such as dithering.

## Concrete interfaces

- The hosted browser tool is available at https://spritefusion.com/pixel-snapper; the repository also links a desktop edition.
- Install a native CLI with Homebrew (`brew install Hugo-Dz/tap/spritefusion-pixel-snapper`) or Cargo (`cargo install spritefusion-pixel-snapper`).
- CLI form: `spritefusion-pixel-snapper <INPUT> <OUTPUT> [COLOR_COUNT] [OPTIONS]`. Inputs and outputs can be PNG/JPEG files or directories for batch processing; color count defaults to 16.
- `--pixel-size <PIXELS>` overrides automatic pixel-size detection. `--palette <HEX,...>` constrains output to specified six-digit hex colors.
- The Rust library can also be built for the web with `wasm-pack build --target web --out-dir pkg --release`. JavaScript calls `process_image(inputBytes, kColors?, pixelSizeOverride?, paletteHex?)`.
- The repository is MIT licensed.

## Workspace relevance

This makes the existing Pixel Snapper note actionable beyond the browser UI: local batch processing and a WASM integration path can fit asset pipelines. Choose the tool when assets have recognizable details worth retaining but need consistent pixel scale and palette; inspect the output because snapping and color quantization still transform the source art.

## Source

- Repository and README: https://github.com/Hugo-Dz/spritefusion-pixel-snapper
- Online tool: https://spritefusion.com/pixel-snapper
- License: MIT
