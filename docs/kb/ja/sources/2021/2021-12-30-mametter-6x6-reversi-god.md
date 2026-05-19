---
title: 6x6リバーシの神
source_url: https://mametter.hatenablog.com/entry/2021/12/30/114111
source_type: post
original_language: ja
ingested_on: 2026-05-19
status: active
tags:
  - reversi
  - rust
  - webassembly
  - browser-games
  - exact-search
published_on: 2021-12-30
named_entities:
  - Rust
  - WebAssembly
  - TypeScript
  - three.js
  - wasm-bindgen
  - wasm-pack
  - Web worker
  - MTD(f)
  - LOUDS
  - Joel Feinstein
related_pages:
  - ../../wiki/projects/reversi-adventure.md
  - ../../wiki/patterns/browser-perfect-play-games.md
---

# ブラウザで動く 6x6 リバーシ完全解析の実装メモ

## ここで重要な理由

これは、小さな完全情報ゲームを server なしで browser に載せるときの強い参照である。特に `reversi-adventure` に対して、解かれた小盤面ゲームを「絶対に勝てないが触って楽しい体験」に変える実装の形を具体化している。

## 要約

- この投稿は、白番を完全プレイする 6x6 リバーシ AI を説明しており、黒番の人間は勝てない。
- 盤面表現の中心は Rust の `u64` 2 本による bitboard で、黒石位置と白石位置を分けて持つ。
- 探索は局面段階ごとに分けられており、終盤は exact negamax、中盤は置換表付き `MTD(f)` とヒューリスティックな move ordering、序盤は 16 手ぶんの事前計算 tree を使う。
- 序盤 tree は `LOUDS` で圧縮し、`include_bytes!` で binary に埋め込んで、browser 配布サイズを現実的に保っている。
- browser 実行は `wasm-bindgen`、`wasm-pack`、`Web worker` を使い、UI は別レイヤとして TypeScript と `three.js` で実装している。

## ワークスペースでの含意

- browser で遊べる完全情報ゲームは、盤面表現を十分に圧縮し、重い探索を段階分割できれば server を必須にしない。
- この投稿は、最終的な exactness を捨てずに heuristic acceleration を混ぜる具体例でもある。評価関数は最終判定そのものではなく、exact search に入る前の move ordering 改善に使われる。
- 序盤の事前計算は単なる速度改善ではなく、interactive な実行時間を browser 側に残すための配布戦略でもある。
- platform 差分の罠も明確で、native では通った実装が WASM では 32-bit の `usize` 前提で壊れた。
- UI と solver を分ける設計も重要で、探索 engine は Rust/WASM に残し、presentation は TypeScript + `three.js` に委ねている。

## 具体アンカー

- 盤面表現: `u64` 2 本の bitboard と、36 マス向けに特化生成した flip code
- 終盤 solver: 空きマス 11 以下で exact negamax
- 中盤探索: `edge2x` と着手可能数を使う move ordering 付き `MTD(f)`
- 学習済み重み: 約 3000 個の浮動小数点重みを Rust 配列として埋め込み
- 序盤戦略: 約 16 手の best-response data を `LOUDS` で圧縮
- browser packaging: `wasm-bindgen`、`wasm-pack`、`Web worker`、TypeScript、`three.js`
- 移植時の落とし穴: native 64-bit では見えなかった WASM 32-bit の `usize` overflow

## フォローアップ

- `reversi-adventure` で deterministic な強敵モードを考えるなら、現在の board representation と search split をこの `bitboard -> phased solver -> browser worker` 形と比較する。
- 今後 client-only hosting を目指す board game を作るなら、WASM に載せる前提の整数幅と payload size を最初から test に含める。
