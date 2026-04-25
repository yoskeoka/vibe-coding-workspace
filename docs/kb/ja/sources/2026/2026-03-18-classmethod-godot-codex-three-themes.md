---
title: Godot と Codex で 3 つの題材を試してみたら、初回出力の強さと人手確認の重さが見えてきた
source_url: https://dev.classmethod.jp/articles/godot-codex-three-theme-verification/
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - ai-game-dev
  - godot
  - codex
  - agentic-dev
published_on: 2026-03-18
named_entities:
  - Godot
  - Codex
  - GPT-5.4
  - EditorPlugin
  - GraphEdit
  - Multiplayer Bomber Demo
related_pages:
  - ../../wiki/topics/ai-game-dev.md
  - ../../wiki/tools/godot.md
---

# 3 つのゲーム開発課題における Godot と Codex

## ここで重要な理由

これは、ゲーム開発において AI コーディングエージェントが初手でどこまでできて、どこから人間の確認が必要になるかを比較した KB 内でもかなり具体的な資料です。

## 要約

- `Godot 4.6.1` と `Codex` を使い、2D アクション、小規模マルチプレイ、`EditorPlugin` ベースの制作ツールという 3 つの題材を検証している。
- 結論は 3 題材で一貫しており、初回出力は有用な土台を素早く作れるが、細かなチューニング、操作感、エディタの使い勝手、最終確認は人手が必要だとしている。
- 比較が混ざらないよう、題材ごとにセッションとディレクトリを分けて検証している。
- 見るべき観点も具体的で、アクションでは移動感、壁アクション、コヨーテタイム、カメラ挙動、マルチプレイでは所有権と同期整合性、プラグインでは実用的な editor UX が焦点になっている。
- ランタイムコードと editor 拡張の両方を同じ環境で試せる点が、この種の実験で Godot を有用にしていると位置づけている。

## ワークスペースでの含意

- 大きく性質の違う問題群でエージェント性能を比べるときは、題材ごとにセッションを分けるべき。
- 成否は、操作感、衝突の正しさ、同期整合性、エディタ運用の粗さのようなごまかしにくい面で判定するべき。
- 今後 Godot で AI 補助開発を試すなら、最初のプロンプトだけでなく、人手で見つけた逸脱も一緒に残すとよい。その差分こそ実際のベンチマークになる。
