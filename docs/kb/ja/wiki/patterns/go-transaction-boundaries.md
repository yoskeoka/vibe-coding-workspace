---
title: Go トランザクション境界
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-18-tonbi-go-clean-architecture-transactions.md
---

# Go トランザクション境界

## パターン

layered な Go service では、トランザクション所有権は通常、transport handler ではなく business operation 側に置くべきである。handler や repository interface に transaction primitive が漏れるほど、責務境界は曖昧になる。

## 望ましい進め方

- 1 つの business operation がきれいに 1 transaction に対応するなら、まず usecase 所有の transaction boundary から始める。
- 多くの usecase が同じ transaction pattern を必要としたり、usecase から ORM 固有知識を減らしたくなったら、`TransactionManager` 抽象を導入する。
- 大きな repository 群が繰り返し一緒に移動する場合にだけ `UnitOfWork` を検討する。

## 既定で避けるもの

- endpoint ごとに散らばる handler 主導の `Begin` / `Commit` / `Rollback` 制御。
- 明確な正当化なしに transaction control method を抱え込んだ repository interface。
