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

# Next2D の起源とエンジン / ツール設計の学び

## ここで重要な理由

この記事は、ゲーム制作スタックをより多く自前で持つことの動機と技術的足場を同時に示しています。エンジン設計、性能設計、制作ツール、AI 時代のワークフロー統合がひとつにつながっています。

## 要約

- `Next2D` は、「自分は本当にゲームを作れるのか」という個人的な問いへの応答として位置づけられている。
- 構成は `Player` ランタイム、`MVVM + Clean Architecture + Atomic Design` に基づく `Framework`、そしてブラウザベースの `AnimationTool` の 3 要素として説明されている。
- 技術面では、描画を worker 側の `OffscreenCanvas` に寄せ、ロジックは main thread に残し、コピーではなく `Transferable` な `ArrayBuffer` 所有権移譲を使うことで `60fps` を維持する設計が中心になっている。
- メモリ節約、buffer 再利用、object pool、非同期な main/sub thread の滑らかさ維持と GC spike 回避の難しさも強調している。
- `WebGPU` 対応も、きれいな書き直しではなく、既存の `WebGL` 風コマンドストリームとの互換性を意識した慎重な拡張として語られている。
- 2026 年公開の `Next2D` `MCP` サーバーにも触れており、エンジン / ツールチェーンの話をそのまま AI 補助開発文脈につないでいる。

## ワークスペースでの含意

- エンジンを自前で持つ価値は、ツール、描画構成、コンテンツ制作フローをまとめて設計できるときに初めて出る。
- `OffscreenCanvas`、`Transferable`、buffer 再利用、object pool が、このシステムの実際の性能形状を示す耐久アンカーになる。
- `MCP` 参照があるので、これは単なるエンジン実装資料ではなく、今後の AI 補助ワークフロー向けに文脈豊富な環境という意味でも重要。
