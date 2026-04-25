---
title: 世界初!?バイブコーディング・縦スクロールシューティングを作った
source_url: https://note.com/shi3zblog/n/nc93e34da423f
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - ai-game-dev
  - vertical-shooter
  - voice-ui
  - llm
published_on: 2026-04-11
named_entities:
  - GIKEN STAR
  - Bonsai
  - FLUX Klein 4B
  - Qwen3.5-9B
  - Suno
related_pages:
  - ../../wiki/topics/ai-game-dev.md
---

# shi3z による AI 駆動縦シューティング試作

## ここで重要な理由

このノートは、通常のコード生成より一歩踏み込んだ「AI ゲーム」ループを具体的に捉えています。音声入力、リアルタイム画像生成、LLM によるルール解釈、生成されたボス行動がすべて live play の中に入っています。

## 要約

- 試作ゲーム `GIKEN STAR` は、レバー移動、ショットボタン、爆弾ボタンに加えて、爆弾発動中の spoken keyword を使う縦スクロールシューティングとして作られている。
- keyword を叫ぶと speech recognition の結果を LLM が解釈して power-up を決め、さらに player ship と enemy の見た目も keyword に合わせて realtime 再生成される。
- 話された keyword を解釈し、パラメータ化された装備挙動に落とす仕組みとして `Bonsai` が明示されており、より強い・より具体的な表現が loadout を変える。
- `FLUX Klein 4B` が player ship の見た目再生成に使われ、`Qwen3.5-9B` は攻撃プログラムや boss の行動生成に使われ、コード生成失敗時には自動 retry も入っている。
- 記事では internet news から新しい話題を引き、序盤 enemy に時事テーマを反映させる流れにも触れている。

## ワークスペースでの含意

- これは、AI が開発を補助するだけでなく、gameplay loop 自体の中に入る例として強い retrieval anchor になる。
- voice keyword の設計は game balance の一部になる。記事では、よりよい装備を引き出す prompt 発明そのものが上達要素として描かれている。
- 将来の試作では、speech recognition、LLM interpretation、code generation、image generation、generation failure 時の fallback behavior という stack 分割を明示して残すべき。
