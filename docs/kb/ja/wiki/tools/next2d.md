---
title: Next2D
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-04-10-sonicmoov-next2d-origin-story.md
---

# Next2D

## 現在のシグナル

`Next2D` は、runtime、framework の規約、animation tooling、AI 向け context integration まで含めて browser-game stack をより広く自前化する参照先として強い。

## メモ

- このシステムは engine runtime 単体ではなく、`Player + Framework + AnimationTool` として説明されている。
- runtime 側では、`OffscreenCanvas`、worker ベース rendering、`Transferable` `ArrayBuffer` handoff、buffer 再利用、object pool によって framerate を守る設計が強調されている。
- framework 側は `MVVM + Clean Architecture + Atomic Design` に支えられており、構造化された game UI / application architecture の retrieval anchor としても relevant。
- 記事は `Next2D` の `MCP` server にも触れており、この stack が rendering performance だけでなく、将来の AI 補助開発フローにも関係することを示している。

## 関連ページ

- [ai-game-dev](../topics/ai-game-dev.md)
