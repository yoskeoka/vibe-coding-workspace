---
title: ソース取り込み
last_reviewed: 2026-04-25
status: active
sources:
  - ../../sources/2026/2026-04-11-karpathy-llm-knowledge-bases-x-post.md
  - ../../sources/2026/2026-04-11-karpathy-llm-wiki-gist.md
  - ../../sources/2026/2026-04-11-classmethod-llm-knowledge-base.md
---

# ソース取り込み

## パターン

ユーザーが URL と少量の文脈を渡し、エージェントがそれを source note に変換し、孤立した bookmark のまま放置せず durable な wiki page を更新する。

## 良い取り込み動作

- URL ごとに 1 つの source note を保つ
- 記事の完全再現ではなく、ワークスペースでの有用性を基準に要約する
- `Obsidian`、`Obsidian Web Clipper`、`Mem0`、`pgvector`、`/kb-compile`、`index.md`、`log.md` のような concrete retrieval anchor が重要なら落とさない
- 既存の concept page を新設より優先して更新する
- source が experiment を示唆するなら、短い follow-up idea も残す

## アンチパターン

- source の長文をそのまま repository に複製する
- link ごとに新しい wiki page を作る
- synthesized page 更新時に provenance を失う
- 後で検索や比較に効く tool 名、command 名、file 名を圧縮しすぎて消す
