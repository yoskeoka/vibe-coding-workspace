---
title: Sprite Fusion Pixel Snapper のリポジトリ
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

# Sprite Fusion Pixel Snapper のリポジトリ

## 要約

`Hugo-Dz/spritefusion-pixel-snapper` は、AI生成や手続き生成のピクセルアートを後処理する Sprite Fusion Pixel Snapper の Rust 実装です。不揃いなピクセルを規則的なグリッドに合わせ、解像度を一定に保ち、ディザリングなどの細部をなるべく残しながら色を量子化します。

## 具体的な利用方法

- ブラウザー版: https://spritefusion.com/pixel-snapper。リポジトリからデスクトップ版にも移動できます。
- Homebrew (`brew install Hugo-Dz/tap/spritefusion-pixel-snapper`) または Cargo (`cargo install spritefusion-pixel-snapper`) で native CLI を導入できます。
- CLI 形式: `spritefusion-pixel-snapper <INPUT> <OUTPUT> [COLOR_COUNT] [OPTIONS]`。PNG/JPEG ファイルまたは directory を入出力にでき、directory は batch 処理します。色数の既定値は 16 です。
- `--pixel-size <PIXELS>` で自動検出されたピクセルサイズを上書きできます。`--palette <HEX,...>` で6桁hex色の palette を指定できます。
- Web 向けには `wasm-pack build --target web --out-dir pkg --release` でビルドし、JavaScript から `process_image(inputBytes, kColors?, pixelSizeOverride?, paletteHex?)` を呼べます。
- ライセンスは MIT です。

## ワークスペースでの意義

既存の Pixel Snapper の記事ノートを、ブラウザー利用だけでなく、ローカルの batch 処理や WebAssembly 組み込みまで拡張する情報です。形状の細部を残しながらグリッドと色数を整えたい素材に向きますが、変換結果を目視で確認してから採用します。

## 出典

- Repository / README: https://github.com/Hugo-Dz/spritefusion-pixel-snapper
- Online tool: https://spritefusion.com/pixel-snapper
- License: MIT
