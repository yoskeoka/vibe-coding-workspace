---
title: Godot
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-18-classmethod-godot-codex-three-themes.md
---

# Godot

## 現在のシグナル

Godot は runtime engine としてだけでなく、gameplay code と editor extension を横断して AI 補助の検証を行える共通環境として有用である。

## メモ

- 今のところ最も強い KB signal は、`Godot 4.6.1` と `Codex` を 2D action、小規模 multiplayer、`EditorPlugin` tooling に適用した Classmethod の比較記事である。
- この比較における Godot の価値は、操作感が重要な runtime task と editor UX task の両方を、ひとつの toolchain で試せる点にある。
- AI が生成した first pass は動く土台として十分強いが、難所は依然として人間が確認する側に残る。movement feel、synchronization correctness、editor usability の粗さがそれに当たる。
- `Multiplayer Bomber Demo`、`EditorPlugin`、`GraphEdit` は今後の Godot 実験でも useful retrieval anchor になる。

## 関連ページ

- [ai-game-dev](../topics/ai-game-dev.md)
