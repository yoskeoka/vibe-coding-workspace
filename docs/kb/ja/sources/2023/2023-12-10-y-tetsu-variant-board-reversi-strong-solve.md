---
title: 〖リバーシ(オセロ)〗変形ボードの強解決！
source_url: https://qiita.com/y-tetsu/items/6e91ceb110951d799704
source_type: article
original_language: ja
ingested_on: 2026-05-19
status: active
tags:
  - reversi
  - othello
  - edax
  - exact-search
  - board-editors
published_on: 2023-12-10
named_entities:
  - Edax
  - Othello is Solved
  - setboard
  - MinGW
related_pages:
  - ../../wiki/projects/reversi-adventure.md
  - ../../wiki/patterns/browser-perfect-play-games.md
---

# 最小限の Edax 改造で変形盤面リバーシを強解決する

## ここで重要な理由

これは、solver 単体の話と playable な board-game tooling の間をつなぐ参照である。既存の強い engine を、ルール境界の小さな patch だけで custom board shape に流用し、その上に軽量な盤面 authoring と棋譜 replay を載せるやり方が具体的に残っている。

## 要約

- この記事は、非標準形状や非標準初期配置を持つ小さなリバーシ盤面を `Edax` 改造で強解決している。
- 表現上の中心変更は、黒石・白石 bitboard に加えて、穴を表す 3 本目の 64-bit mask `H` を導入することにある。
- move generator の中核はほぼそのままで、最終的な合法手 mask を `~(P|O)` から `~(P|O|H)` に変えて穴位置を除外している。
- engine 設定は 60 手先まで読む full-depth 探索に寄せ、標準 8x8 用 opening book は無効化している。
- workflow として `setboard` による custom board 文字列入力、`display` による棋譜表示、さらに browser 上の盤面 editor / replay site を用意している。

## ワークスペース向けの示唆

- custom board variant を素早く試したいなら、新しい solver を最初から作るより、既存の実績ある engine に薄い patch を当てる方が早い場合がある。
- variant 対応は、信頼済みの move-generation core があるなら、board-mask abstraction だけで十分開けることがある。
- ここでも責務分離が重要で、native engine 側の変更は最小限に留め、board editing や replay は補助ツール側に逃がしている。
- 既存の 6x6 Rust/WASM 参照と競合するのではなく補完関係にある。あちらは browser 配布に強く、こちらは custom rule / shape の早期検証に強い。
- 7 種類のサンプル盤面は、「解ける」と「遊んで気持ちよい」は別条件だと示している。特に `リング` のような手数の多い盤面は解析に最大 2 分ほどかかる。

## 具体アンカー

- engine baseline: `Edax`
- 強解決向け設定: 探索深さを `21` から `60` に引き上げ、book 使用を無効化
- 変形盤面表現: 穴用の第 3 bitboard `H`
- 合法手 patch: 最終 mask を `~(P|O|H)` に変更
- 盤面 authoring interface: `setboard "<64 cells><side-to-move>"`、穴は space で表現
- replay support: `display` コマンドで棋譜を出力
- サンプル形状: `4x4`、`十字`、`風車`、`卍`、`X`、`外側に石`、`リング`

## フォローアップ

- `reversi-adventure` で変則盤面や authored puzzle board を試すなら、browser 実装に直行する場合と、このような native solver prototyping loop を先に回す場合を比較するとよい。
- 今後 external engine を再利用する board-game 実験では、variant support patch を legal-move boundary 近辺に閉じ込め、solver correctness を追いやすく保つ。
