---
title: Go Transaction Boundaries
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-18-tonbi-go-clean-architecture-transactions.md
---

# Go Transaction Boundaries

## Pattern

In layered Go services, transaction ownership should usually live with the business operation, not the transport handler. The more transaction primitives leak into handlers and repository interfaces, the more responsibility boundaries blur.

## Preferred progression

- Start with a usecase-owned transaction boundary when one business operation maps cleanly to one transaction.
- Introduce a `TransactionManager` abstraction when many usecases need the same transaction pattern or when reducing ORM-specific knowledge in usecases becomes valuable.
- Consider `UnitOfWork` only when a large repository set repeatedly travels together.

## Avoid by default

- Handler-led `Begin` / `Commit` / `Rollback` orchestration spread across endpoints.
- Repository interfaces polluted with transaction control methods unless that trade-off is clearly justified.

