---
title: サーバー代0円・バックエンド不要！？remoteStorageで実現する「ユーザー主権」なWebアプリの実装パターン
source_url: https://zenn.dev/konnyaku256/articles/web-app-impl-patttern-using-remote-storage
source_type: article
original_language: ja
ingested_on: 2026-05-19
status: active
tags:
  - deployment
  - web-app
  - remotestorage
  - privacy
related_pages:
  - ../../wiki/topics/deployment-options.md
  - ../../wiki/patterns/user-owned-storage-webapps.md
---

# サーバー代0円・バックエンド不要！？remoteStorageで実現する「ユーザー主権」なWebアプリの実装パターン

## ここで重要な理由

これは単なる「安い backend」論の別案ではなく、特定の種類の web app では backend 自体を持たないという比較軸を与えてくれる。

## 要約

- 記事は `remoteStorage` を、app 開発者ではなく user 自身が data ownership を持つための open protocol として位置づけている。
- 実装モデルは app と storage を明確に分離する。frontend app が user に storage endpoint への接続を求め、その user 管理の backend に対して直接 read/write する。
- 残すべき protocol anchor は具体的で、endpoint discovery の `WebFinger`、authorization の `OAuth2`、file / directory 操作用の REST-like HTTP API が核になる。
- JavaScript 実装の中心は `remotestorage.js` で、`notes` のような access scope、接続 widget、`storeFile()` や `getFile()` のような client API が使われる。
- provider 例として `5apps`、plugin を入れた `Nextcloud`、`Armadietto` が挙げられ、adapter 的に `Dropbox` や `Google Drive` を使う案にも触れている。
- 開発者側の利点は、backend hosting cost をほぼゼロにできること、運用する database が不要になること、そして「data を app 作者の server に預けない」という privacy 説明がしやすいこと。

## ワークスペースでの含意

- user ごとの note、progress、setting のような personal data を主に扱う project では、「user-owned storage」を選択肢に残しておく価値がある。
- これは server-side system の一般解ではない。shared authoritative state、重い backend computation、multi-user coordination が重要なら適合しにくい。
- hobby web app が `static frontend + user-chosen storage` で成立するなら、最安 backend stack を足す前に比較すべき候補になる。

## フォローアップ

- 将来 sync が欲しい project では、local-only persistence や薄い hosted backend と並べて `remoteStorage` を比較する。
