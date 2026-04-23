---
title: llm-wiki
source_url: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f
source_type: docs
original_language: en
ingested_on: 2026-04-11
status: active
tags:
  - knowledge-base
  - llm
  - wiki
related_pages:
  - ../../wiki/topics/llm-knowledge-bases.md
  - ../../wiki/patterns/source-ingestion.md
---

# llm-wiki

## ここで重要な理由

このソースはワークスペース KB の実装観に最も近いです。X 投稿のアイデアを、repo-native Markdown と agent skill に落とし込める運用モデルとして明文化しています。

## 要点

- `sources`、`schema`、compiled `wiki` の 3 層モデルを提案している
- `ingest`、`query`、`lint` を継続保守ループとして扱っている
- `index.md` と `log.md` のような durable navigation anchor を重視している
- 大きな RAG を先に入れず、小中規模では保守された wiki で十分役立つと示している

## ワークスペースでの含意

- agent が直接編集できる Markdown を source of truth に保つ
- `schema.md` で maintainer としての振る舞いを縛る
- source-oriented note と durable wiki page を明確に分ける
