---
title: reversi-adventure
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-04-11-openai-gpt-5-4-phaser-tactical-rpg.md
  - ../../sources/2026/2026-04-22-unity-interactive-writing-esoteric-ebb.md
  - ../../sources/2021/2021-12-30-mametter-6x6-reversi-god.md
  - ../../sources/2023/2023-12-10-y-tetsu-variant-board-reversi-strong-solve.md
---

# reversi-adventure

## 関連性

この project はワークスペース内で進行中の game-development testbed なので、engine stack が違っても AI 補助ゲーム制作のパターンを将来的に移植できる可能性がある。

## 現時点の要点

- Phaser tactical RPG の参照は直接再利用できるわけではないが、agent-driven な playtest loop や、より厚みのある game prototype への道筋を補強している。
- `Esoteric Ebb` の narrative-authoring ノートは、adventure layer が authored choice、quest state、player-specific callback を持つようになったとき relevant になる。転用したいのは CRPG 規模ではなく、state naming、choice と後の feedback の接続、分岐 content の事前 validation という規律である。
- まめめもの 6x6 リバーシ記事は、browser 上で deterministic な board-game core を成立させる最良の参照である。compact な bitboard、段階分割した exact search、圧縮した opening data、Rust/WASM solver と front-end 表示の分離という形が具体的に残っている。
- y-tetsu の変形盤面リバーシ記事は、solver を作り直さずに非矩形 board や authored board を試すための最良の参照である。移植したい核は、実績ある engine に対する薄い `hole` mask patch と、盤面 authoring / 棋譜 replay を別 tooling に逃がす分離である。

## 関連ページ

- [ai-game-dev](../topics/ai-game-dev.md)
- [branching-narrative-authoring](../patterns/branching-narrative-authoring.md)
- [browser-perfect-play-games](../patterns/browser-perfect-play-games.md)
