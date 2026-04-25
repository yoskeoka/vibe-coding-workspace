---
title: LLM Knowledge Bases
source_url: https://x.com/karpathy/status/2039805659525644595
source_type: post
original_language: en
ingested_on: 2026-04-11
status: active
tags:
  - knowledge-base
  - llm
  - x
related_pages:
  - ../../wiki/topics/llm-knowledge-bases.md
  - ../../wiki/patterns/source-ingestion.md
  - ../../wiki/projects/ww.md
---

# LLM Knowledge Bases

## ここで重要な理由

これは、このパターンに名前を与えた原点の投稿です。完全な仕様書というより、「単発の回答やコード生成だけでなく、知識の保守に LLM の token 予算を使う」という動機づけに価値があります。

## 要約

- Karpathy は、LLM を active research topic 向けの personal knowledge base 構築に使う方向転換を示している。
- 印象的な主張は、最近の token 消費の多くが code ではなく knowledge の操作に向いているという点である。
- 投稿自体は発表と導線の役割で、より詳しいファイル構成や運用ループは続く gist に譲っている。
- 周辺議論も含めると、重い RAG 基盤を早くから持ち込む代わりに、小規模向けの現実的な代替として読める。

## ワークスペースでの含意

- この投稿は、ワークスペース KB の「なぜ今やるのか」を示すソースとして残しておく価値がある。
- 実装詳細は gist や記事に譲りつつ、パターン名の起点として X 投稿も必ず紐づけておくべき。

## フォローアップ

- 将来コラボレータに KB を説明するときは、gist 単体ではなくこの投稿とセットで参照する。
