---
title: LLM ナレッジベース
last_reviewed: 2026-04-25
status: active
sources:
  - ../../sources/2026/2026-04-11-karpathy-llm-knowledge-bases-x-post.md
  - ../../sources/2026/2026-04-11-karpathy-llm-wiki-gist.md
  - ../../sources/2026/2026-04-11-classmethod-llm-knowledge-base.md
---

# LLM ナレッジベース

## 要約

中核パターンは、curated source から compiled Markdown wiki を LLM に継続保守させ、bookmark や ad hoc search から毎回同じ情報を掘り直さないようにすること。

## このワークスペースで合う理由

- このワークスペースはすでに git 管理の Markdown を AI context として使っている。
- 知識は複数の hobby project をまたいで蓄積されるべきで、chat history に閉じ込めないほうがよい。
- compiled wiki は、散在した source link 群より skim しやすい。

## 動作モデル

- 原点は Karpathy の X 投稿で、「code だけでなく knowledge の操作にもっと LLM を使う」という方向づけ。
- `sources/` には article、paper、repo、image、post などの source-oriented note と provenance を置く。
- `schema.md` は、他所の `CLAUDE.md` のように、agent の保守ルールを定義する。
- `wiki/` には durable な synthesized page を置き、`index.md` と `log.md` のような retrieval anchor を保つ。
- `ingest`、`query`、`lint` が運用ループを成す。

## 残すべき具体アンカー

- 閲覧やローカル知識 UX: `Obsidian`
- web capture の例: `Obsidian Web Clipper`
- 実務記事で触れられた search / memory layer: `Mem0`、`pgvector`
- 記事内の workspace command 例: `/kb-compile`
- 繰り返し重要視される core file: `schema.md`、`index.md`、`log.md`

## トレードオフの捉え方

- このパターンは anti-RAG ではなく、query のたびに全部を作り直すことに反対している。
- durable layer は compiled wiki であり、その隣に search layer があってもよい。
- small-to-medium scale では、とくに魅力が高い。corpus が巨大になって初めて重い retrieval 基盤が要る。

## 関連ページ

- [source-ingestion](../patterns/source-ingestion.md)
- [ww](../projects/ww.md)
