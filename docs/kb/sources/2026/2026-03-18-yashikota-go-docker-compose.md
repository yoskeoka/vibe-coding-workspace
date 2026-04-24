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

# Modern Go Dockerfile and Compose guidance

## Why it matters here

This is a compact operational reference for current Go container defaults. It is useful because it focuses on concrete low-friction conventions rather than a large platform-specific stack.

## Summary

- The article recommends multistage builds by default and argues for pinning image digests instead of floating tags for reproducibility.
- It calls out `BuildKit`-era Dockerfile features explicitly, including the `# syntax=docker/dockerfile:1` header and official `Build Check` support.
- It recommends `.dockerignore` for keeping build context clean and highlights bind mounts for dependency and source handling where appropriate.
- For Go specifically, it notes newer defaults such as `-s` implying `-w` in linker flags and the effective readonly behavior of modern `go` module commands.
- On the Compose side, the article prefers `compose.yaml`, dropping the legacy top-level `version:` field and using `docker compose` instead of `docker-compose`.
- The example stack uses Go plus Postgres, with a development image that includes `Air` and `Delve` alongside a smaller production stage.

## Workspace takeaways

- This is a good baseline for fresh Go service repos that need sane local-dev and production container defaults without heavy platform machinery.
- Keep the retrieval anchors explicit in future container docs: `compose.yaml`, `docker compose`, digest pinning, `Build Check`, and multistage builds.

