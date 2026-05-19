---
title: Agent Sprite Forge
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-05-15-npaka-agent-sprite-forge.md
---

# Agent Sprite Forge

## 現在のシグナル

`Agent Sprite Forge` は、2D game 向け asset を raw image ではなく engine で使いやすい形まで持っていく Codex-oriented workflow である。

## メモ

- tool surface は `generate2dsprite` と `generate2dmap` に分かれている。
- transparent sprite sheet、抽出 animation frame、GIF preview、layered map asset、prop pack、collision metadata、engine-ready export が具体 output anchor になる。
- background removal、alignment、slicing、validation、export まで local finishing pass が含まれている点が重要である。
- そのため art ideation だけでなく、「playable prototype を止めない程度の asset set を短時間で作る」という目的にも効く。

## 関連ページ

- [ai-game-dev](../topics/ai-game-dev.md)
- [pixel-art-ui](../topics/pixel-art-ui.md)
