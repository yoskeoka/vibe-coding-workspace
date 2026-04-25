---
title: Go コンテナワークフロー
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-18-yashikota-go-docker-compose.md
---

# Go コンテナワークフロー

## パターン

新しい Go service では、modern で低摩擦な既定値として、multistage build、digest pinning、`BuildKit` 時代の Dockerfile syntax、明示的な build check、現在の Compose 命名と command 慣習を採る。

## 具体アンカー

- `# syntax=docker/dockerfile:1`
- 公式の `Build Check`
- `.dockerignore`
- digest pinning
- `compose.yaml`
- `docker compose`
- `Air` / `Delve` を含む development image

## メモ

- local development 用 stage と production 用 stage は分ける。
- floating tag や legacy Compose 命名のような convenience alias より、再現可能な既定値を優先する。
