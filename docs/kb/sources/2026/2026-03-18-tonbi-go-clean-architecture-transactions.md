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

# Go transaction-boundary patterns with clean architecture

## Why it matters here

This is a practical reference for where transaction control should live in a layered Go service. It is relevant to workspace backend code because it compares concrete responsibility boundaries instead of discussing transactions only abstractly.

## Summary

- The article uses a three-table article-creation example to compare four transaction patterns: handler-led, usecase-led, `TransactionManager`, and `UnitOfWork`.
- Its main critique of handler-led transactions is responsibility leakage: transport code starts knowing persistence internals, and `Begin` / `Commit` / `Rollback` logic spreads across endpoints.
- The recommended default is `usecase-led` transactions via `db.Transaction(...)`, because it centralizes rollback behavior and keeps handlers thin.
- A stronger abstraction is a `TransactionManager` interface that hides `db.Transaction(...)` behind an infrastructure-layer implementation while preserving the usecase's intent that a transaction is needed.
- `UnitOfWork` is framed as useful only when many repositories have to travel together in one transaction; otherwise it adds abstraction cost without much payoff.

## Workspace takeaways

- Default to usecase-owned transaction boundaries when the business operation naturally maps to one transaction.
- Reach for `TransactionManager` when repeated transaction logic starts appearing across multiple usecases or when the codebase wants to reduce ORM-specific knowledge in the usecase layer.
- Avoid pushing transaction primitives into handlers or repository interfaces unless there is a very specific reason and the cost is accepted explicitly.

