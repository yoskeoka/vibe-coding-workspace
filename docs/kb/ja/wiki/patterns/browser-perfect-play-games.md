---
title: ブラウザで遊べる完全プレイ盤上ゲーム
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2021/2021-12-30-mametter-6x6-reversi-god.md
  - ../../sources/2023/2023-12-10-y-tetsu-variant-board-reversi-strong-solve.md
---

# ブラウザで遊べる完全プレイ盤上ゲーム

## パターン

小さな完全情報盤上ゲームを browser で即応的に遊ばせたいなら、solver は compact で phase-aware にし、rendering layer から分離する。browser 実行時間は、既知の序盤計算や高コスト表現のやり直しではなく、interaction に使うべきである。

## 具体アンカー

- 6x6 リバーシのような盤面は、`u64` 2 本の bitboard のように cheap copy と bitwise move generation がしやすい形で持つ。
- solve path は局面段階で分ける。終盤は exact search、中盤は heuristic-assisted ordering、序盤は事前計算 tree を使う。
- product claim が perfect play に依存するなら、heuristic は exact search の代替ではなく加速用途に限定する。
- 事前計算した opening data は、`LOUDS` のような方法で十分に圧縮し、client と一緒に配布できる形にする。
- solver は `WebAssembly` と `Web worker` に載せ、UI 技術の選定は solver 実装言語から切り離す。
- 整数幅と target architecture を product contract の一部として扱う。native 64-bit 前提は browser 側の実質 32-bit target で壊れうる。
- 盤面形状 variant を試す段階では、move generator を全面的に書き換える前に、legal-move boundary に board-mask 拡張を差し込めないか確認する。
- board editing や棋譜 replay の tooling は core solver から切り離し、custom shape や solved record を速く回せるようにする。

## トレードオフ

- exact / near-exact search は data layout、search ordering、precompute に多くの複雑さを押し込む。
- 最速の内部表現は `unsafe` のような鋭い手段を正当化しうるが、保守コストは上がる。
- browser packaging では、native 実行とは別種の bug が出る。整数幅、worker 境界、bundler 設定が典型である。
- native engine 再利用は実験速度を上げるが、build 前提、配布形式、platform 固有 setup も一緒に背負う。
- perfect-play 体験は魅力的だが、曖昧でスタイル重視の opponent より設計自由度は狭くなる。

## 関連ページ

- [reversi-adventure](../projects/reversi-adventure.md)
