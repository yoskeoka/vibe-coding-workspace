---
title: 令和最新版 GoでのDockerfile / Docker Composeの書き方
source_url: https://zenn.dev/yashikota/articles/cb4a49553bd368
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - go
  - docker
  - compose
  - containers
published_on: 2026-03-18
named_entities:
  - Docker
  - Docker Compose
  - BuildKit
  - Build Check
  - Trivy
  - cosign
  - Air
  - Delve
  - Postgres
related_pages:
  - ../../wiki/patterns/go-container-workflows.md
---

# 現代的な Go 向け Dockerfile / Compose 指針

## ここで重要な理由

これは、今どきの Go コンテナ運用の既定値を短く押さえた実務メモです。大きなプラットフォーム事情ではなく、すぐ使える低摩擦な慣習に絞っている点が有用です。

## 要約

- 既定で multistage build を勧め、再現性のために floating tag ではなく image digest を pin するべきだとしている。
- `BuildKit` 時代の Dockerfile 機能を明示的に扱い、`# syntax=docker/dockerfile:1` ヘッダや公式の `Build Check` 対応を挙げている。
- build context を健全に保つための `.dockerignore` と、必要に応じた bind mount の使い分けを勧めている。
- Go 固有の注意として、linker flag で `-s` が `-w` を含意することや、最近の `go` モジュール系コマンドの実質的な readonly 性にも触れている。
- Compose 側では、`compose.yaml`、top-level `version:` の省略、`docker-compose` ではなく `docker compose` の利用を推している。
- 例では Go と Postgres を使い、開発用イメージに `Air` と `Delve` を含めつつ、本番用ステージは小さく保っている。

## ワークスペースでの含意

- 重い基盤を持ち込まずに、ローカル開発と本番のコンテナ既定値を整えたい新規 Go サービスのよい土台になる。
- 今後コンテナ文書を増やすときは、`compose.yaml`、`docker compose`、digest pinning、`Build Check`、multistage build といった検索アンカーを明示して残すべき。
