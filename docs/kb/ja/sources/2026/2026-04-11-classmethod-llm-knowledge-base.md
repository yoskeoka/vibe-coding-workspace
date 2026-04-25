---
title: Karpathy 氏が言語化した「LLM Knowledge Base」というパターン
source_url: https://dev.classmethod.jp/articles/karpathy-llm-knowledge-base/
source_type: article
original_language: ja
ingested_on: 2026-04-11
status: active
tags:
  - knowledge-base
  - llm
  - workflow
related_pages:
  - ../../wiki/topics/llm-knowledge-bases.md
  - ../../wiki/patterns/source-ingestion.md
---

# Karpathy 氏が言語化した「LLM Knowledge Base」というパターン

## ここで重要な理由

この記事は、抽象的な gist の考え方を、具体的なツールやコマンドが見えるワークスペース運用に落とし込んでいます。Karpathy の概念と、このワークスペースで実際にどう回すかをつなぐ、今のところ最良の橋渡しです。

## 要約

- 中心動作を、検索結果の断片に答えさせるのではなく、curated document を LLM に構造化 Markdown へ compile させることだと整理している。
- 層構造を raw source、schema、wiki の 3 層として示し、記事、論文、リポジトリ、画像などを具体例にしている。
- `Obsidian`、`Obsidian Web Clipper`、`CLAUDE.md`、`Mem0`、`pgvector`、`/kb-compile` といった具体名を残している点が重要。
- 永続的な compiled knowledge を重視しつつ、別途 search layer を併存させてもよいとして、純粋な RAG との差を説明している。
- project 単位 compile、全 project compile、矛盾や broken link、stale page を見る lint 風チェックなど、保守操作も具体的に挙げている。

## ワークスペースでの含意

- `tools/kb build` のようなコマンドは、タスク実行ワークフローとは独立したまま保つべき。
- 大きな製品や基盤を先に作らなくても、小さな command と skill の組み合わせでこのパターンは十分実体化できる。
- project page、topic page、`index.md`、`log.md` は、平坦な bookmark 一覧より強い検索アンカーになる。

## フォローアップ

- ingest は完全自動ではなく、対話的なまま保つ。
- 記事全文コピーより source note を優先する。
- 将来必要になったら、compiled wiki の横に別 search layer を持つか再検討する。
