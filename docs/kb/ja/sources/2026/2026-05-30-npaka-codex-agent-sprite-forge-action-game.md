---
title: Codex + Agent Sprite Forge でアクションゲームの生成を試す
source_url: https://note.com/npaka/n/ne0c0daafbfa3
source_type: article
original_language: ja
ingested_on: 2026-11-28
status: active
tags:
  - agent-sprite-forge
  - codex
  - phaser
  - game-development
related_pages:
  - ../../wiki/tools/agent-sprite-forge.md
  - ../../wiki/topics/ai-game-dev.md
---

# Codex + Agent Sprite Forge でアクションゲームの生成を試す

## 要約

npaka 氏による、Codex、Agent Sprite Forge、Vite、TypeScript、Phaser を使った横スクロールアクションゲーム制作の反復例。ゲーム画面の参考画像と、キャラクター・ステージ全景の画像を役割付きで渡し、16:9 のゲームと idle/run アニメーションを生成します。その後、キャラクター画像を参照してスプライトを作り直します。

改善は、タイル間隔とジャンプ可能性、ダメージ後のライフ・再開処理、敵の向きや踏みつけ、クリア・ゲームオーバー演出、UI、マップエディタ、効果音と BGM に分けて進めています。背景は Photoshop の生成塗りつぶしで横ループの継ぎ目を調整しています。

## 再利用できる点

- 参照画像ごとの役割と、`Vite + TypeScript + Phaser`、16:9 などの制約を最初の指示に含める。
- 移動・地形、敵、プレイヤーアニメ、UI、エディタ、音声のように、確認する関心ごとごとに反復する。
- 生成物は試作品としてプレイ確認する。記事では、最初の穴がジャンプ不能、エディタとゲーム画面でタイル画像が不一致、保存後に操作不能といった不具合が記録されています。
- 歩行アニメは同じ手足が続けて前に出ることがあるため、手作業で整えたモーション画像を参照にするか、Agent Sprite Forge に画像生成なしの後処理だけを依頼する方法を試す。
- スプライト、背景ループ、エディタ用タイル、音声を実際のゲーム内利用につなげて調整する。

## ワークスペースでの意義

既存の Agent Sprite Forge と AI ゲーム開発のページを補う、Phaser での具体的な一連の制作例です。画像を使った初期生成からゲームプレイ上の修正、素材調整までを一つのループとして扱っています。マップエディタについては、編集結果と実行時表示の一致、および保存後の操作性を確認する材料になります。

## 出典

- 記事: https://note.com/npaka/n/ne0c0daafbfa3
- 著者: npaka
- 公開日: 2026-05-30
