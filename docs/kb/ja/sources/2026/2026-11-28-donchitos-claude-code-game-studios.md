---
title: Claude Code Game Studios
source_url: https://github.com/Donchitos/Claude-Code-Game-Studios
source_type: repo
original_language: en
ingested_on: 2026-11-28
status: active
tags:
  - claude-code
  - ai-game-dev
  - game-workflow
related_pages:
  - ../../wiki/topics/ai-game-dev.md
---

# Claude Code Game Studios

## 要約

`Donchitos/Claude-Code-Game-Studios` は、Claude Code でゲームのアイデアを試作し、段階的に育てるための project template です。`/prototype` は短い設計質問から始め、確認可能な `design/game-brief.md` を作成し、`prototypes/<name>/` に browser game を生成して限定的な browser check を行います。引き渡しでは自動チェックと人間によるプレイを区別します。

試作後は Vite + TypeScript の web project に育てられます。`/start`、`/next-step`、`/brainstorm`、`/setup-engine`、`/project-stage-detect` などの workflow もあります。Web、Godot、Unity、Unreal Engine の構成を記載し、`project.yaml` から文書量、QA、review などの rigor を調整します。repository 自身が、engine 経路と深いシステム設計は browser prototype と同等の評価がまだないと明記しています。

## ワークスペースでの意義

短い brief、すぐ遊べる試作、browser check と人間のプレイを別の証拠として扱う AI ゲーム開発 workflow の比較対象です。engine 対応の主張を評価する際は repository 自身の評価範囲も参照します。template を採用する場合は現在の setup と permission rules を確認します。

## 出典

- Repository: https://github.com/Donchitos/Claude-Code-Game-Studios
- License: MIT
