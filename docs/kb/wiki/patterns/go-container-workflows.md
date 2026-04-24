---
title: Go Container Workflows
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-18-yashikota-go-docker-compose.md
---

# Go Container Workflows

## Pattern

For fresh Go services, the modern low-friction baseline is multistage builds, digest-pinned images, `BuildKit`-era Dockerfile syntax, explicit build checks, and current Compose naming and command conventions.

## Concrete anchors

- `# syntax=docker/dockerfile:1`
- official `Build Check`
- `.dockerignore`
- digest pinning
- `compose.yaml`
- `docker compose`
- development image with `Air` / `Delve`

## Notes

- Keep local-development and production stages separate.
- Prefer reproducible defaults over convenience aliases like floating tags or legacy Compose naming.

