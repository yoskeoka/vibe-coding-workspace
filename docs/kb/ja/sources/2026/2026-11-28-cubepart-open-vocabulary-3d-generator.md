---
title: "CubePart: An Open-Vocabulary Part-Controllable 3D Generator"
source_url: https://cubepart.github.io/
source_type: docs
original_language: en
ingested_on: 2026-11-28
status: active
tags:
  - 3d-generation
  - game-development
  - research
related_pages:
  - ../../wiki/topics/ai-game-dev.md
---

# CubePart: An Open-Vocabulary Part-Controllable 3D Generator

## 要約

CubePart は、文章プロンプトとユーザー定義の部品スキーマから、複数パーツで構成される3Dメッシュを生成・分解する研究システムです。指定した各部品に対応するメッシュを作り、組み合わせた状態でも一貫した形状を保つことを目指しています。部品構造はアニメーション、物理挙動、スクリプト制御に利用できます。

モデルは2段階構成です。まず部品スキーマを考慮した全体形状を生成し、次にパーツ間 attention を使う decoder で個別パーツを生成します。プロジェクトページは、462K assets / 2.02M parts の部品ラベル付き dataset、インタラクティブな mesh gallery、比較評価を紹介しています。

## ワークスペースでの意義

ロボットの手足や宝箱のふたのように、意味のある部品ごとの操作が必要なゲーム素材について、text-to-3D と下流の animation / game logic をつなぐ研究例です。production-ready な素材生成 pipeline とみなす前に、リンク先の論文・コード・demoを確認します。

## 出典

- Project page: https://cubepart.github.io/
- 論文とコードへのリンクは project page に掲載されています。
- project page では SIGGRAPH 2026 と CC BY 4.0 を掲示しています。
