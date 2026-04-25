---
title: クリーンアーキテクチャのGoで複数テーブルを更新するトランザクションの貼り方
source_url: https://zenn.dev/tonbi_attack/articles/fcca8d12fac711
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - go
  - transactions
  - clean-architecture
  - gorm
published_on: 2026-03-18
named_entities:
  - Go
  - GORM
  - Gin
  - MySQL
  - TransactionManager
  - UnitOfWork
related_pages:
  - ../../wiki/patterns/go-transaction-boundaries.md
---

# クリーンアーキテクチャでの Go トランザクション境界パターン

## ここで重要な理由

これは、レイヤードな Go サービスでトランザクション制御をどこに置くべきかを実務的に比較した資料です。抽象論ではなく、責務境界を具体的に比べている点がワークスペースのバックエンド実装にも役立ちます。

## 要約

- 3 テーブルを更新する記事作成例を使って、handler 主導、usecase 主導、`TransactionManager`、`UnitOfWork` の 4 パターンを比較している。
- handler 主導の問題点は責務漏れで、transport 層が永続化の詳細を知り始め、`Begin`、`Commit`、`Rollback` が各 endpoint に散らばることだと批判している。
- 推奨の既定値は `db.Transaction(...)` を使う usecase 主導で、rollback の扱いを中央化しつつ handler を薄く保てるとしている。
- さらに抽象化するなら、インフラ層実装の後ろに `db.Transaction(...)` を隠す `TransactionManager` を置き、usecase 側には「トランザクションが必要」という意図だけを残す。
- `UnitOfWork` は、多数の repository が毎回まとまって動くケースでのみ有効で、それ以外では抽象化コストのほうが大きいと整理している。

## ワークスペースでの含意

- ビジネス操作が自然に 1 トランザクションに対応するなら、まず usecase 所有の境界を既定にするのがよい。
- 複数 usecase に同じトランザクション処理が繰り返し出てきたり、usecase 層から ORM 依存知識を減らしたくなったら `TransactionManager` を検討する。
- 明確な理由がない限り、トランザクションプリミティブを handler や repository interface に押し込まないほうがよい。
