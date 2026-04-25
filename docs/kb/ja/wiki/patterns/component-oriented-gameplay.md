---
title: コンポーネント指向ゲームプレイ
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-10-munimaru-component-oriented-gameplay.md
---

# コンポーネント指向ゲームプレイ

## パターン

ゲームが固定された authored interaction ではなく創発的な組み合わせを狙うなら、単に loosely coupled な code だけでは足りない。明示的な safety boundary、system 間の intent protocol、そして継続的なエンジニア介入なしに組み合わせを試せる data-driven な検証手段が必要になる。

## 具体アンカー

- Safety net: 永続状態と一時状態を分け、possession、swap、動的再構成が起きても元の player state を壊さない。
- Intent protocol: system 同士は相手の concrete type を直接知るのではなく、tag、interface、または同等の intent message で通信する。
- Designer-facing composition: 新しい組み合わせは engine code を毎回触らず、asset と設定で試せる状態にする。

## トレードオフ

- 組み合わせ爆発による balance 調整は急速に難しくなる。
- debug と performance の難易度も上がる。
- framework 側が setup と lifecycle ordering を安全にしてくれないと、認知負荷が高くなる。

## 関連ページ

- [ai-game-dev](../topics/ai-game-dev.md)
