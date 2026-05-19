---
title: ピクセルアート UI
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-04-11-phaser-pixui.md
  - ../../sources/2026/2026-05-15-npaka-agent-sprite-forge.md
  - ../../sources/2026/2026-05-16-gigazine-sprite-fusion-pixel-snapper.md
---

# ピクセルアート UI

## 現在のシグナル

pixel-art に特化した UI tooling は追っておく価値がある。汎用 UI stack は retro-style game の視覚制約と相性が悪いことが多い。

## このワークスペースで重要な理由

- 今後の Phaser prototype で、専用 widget を毎回作らずに crisp な pixel UI を出したい場面がありうる。
- hobby game project では、UI readability が polish 工程の遅いボトルネックになりやすい。
- AI 補助の asset generation が増えると、「pixel-art 風だが pixel art としては破綻している」出力をどう整えるかも別の論点になる。

## 再利用できるアンカー

- `PixUI` は引き続き Phaser 側の widget work の主 anchor である。
- `Agent Sprite Forge` は、sprite、map、prop pack を素早く用意して prototype を前に進めるときの anchor になる。
- `Sprite Fusion Pixel Snapper` は、AI 出力が惜しいが grid / palette の規律を満たしていないときの post-processing anchor になる。

## 実務ルール

生成した pixel art は prompt の結果ではなく pipeline として扱う。

- candidate asset を生成する
- transparency、alignment、slicing を整える
- 必要なら real pixel grid に snap し、palette noise を減らす
- その上で、UI や prototype art として十分かを評価する

## 関連ページ

- [phaser](../tools/phaser.md)
- [agent-sprite-forge](../tools/agent-sprite-forge.md)
