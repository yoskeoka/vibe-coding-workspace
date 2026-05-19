---
title: ユーザー所有ストレージ型Webアプリ
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-05-17-zenn-remotestorage-user-owned-webapps.md
---

# ユーザー所有ストレージ型 Web アプリ

## パターン

個人データ中心の web app では、「永続化する app には必ず自前 backend がいる」と決めつけず、user 自身が持つ storage を app が使う構成を現実的な選択肢として扱う。

## ヒューリスティクス

- note、setting、progress のような user-scoped data が中心なら、このパターンを先に検討する。
- `remoteStorage` の具体 anchor を残す。
- discovery は `WebFinger`
- authorization は `OAuth2`
- browser 実装は `remotestorage.js`
- provider / ecosystem 名として `5apps`、`Nextcloud`、`Armadietto` を覚えておく。
- 比較順は、まず local-only persistence、次に最安 hosted backend。狙いは「常に勝つ architecture」を選ぶことではなく、infra obligation 自体を消せる場面を見落とさないこと。
- shared authoritative state、server-side job、強い multi-user coordination が必要なら早めに不適と判断する。

## ここで重要な理由

このワークスペースには infra cost と運用負荷を抑えたい hobby project が複数あるため、「backend を持たない」という案を検索しやすい形で残しておく価値がある。

## 関連ページ

- [deployment-options](../topics/deployment-options.md)
- [cheap-hosting](./cheap-hosting.md)
