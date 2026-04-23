---
title: ソース取り込み
last_reviewed: 2026-04-24
status: active
sources:
  - ../../sources/2026/2026-04-11-karpathy-llm-knowledge-bases-x-post.md
  - ../../sources/2026/2026-04-11-karpathy-llm-wiki-gist.md
  - ../../sources/2026/2026-04-11-classmethod-llm-knowledge-base.md
---

# ソース取り込み

## パターン

ユーザーが URL と短い文脈を渡し、エージェントが source note を作り、既存 wiki page を更新します。孤立した bookmark を増やさず、durable な整理面に編み込みます。

## 良い取り込み動作

- 具体的な製品名や文書名などの retrieval anchor を落とさない
- 既存ページに自然な置き場所があるなら新規ページを乱立させない
- provenance を source note で保ち、wiki page では再利用可能な要点に圧縮する
