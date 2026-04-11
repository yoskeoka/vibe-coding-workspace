---
title: Deployment Options for Hobby Projects
last_reviewed: 2026-04-11
status: seed
sources:
  - ../../sources/2026/2026-04-11-qiita-free-personal-deploy-options.md
---

# Deployment Options for Hobby Projects

## Current default posture

Prefer deployment choices that minimize cost risk first, then minimize setup friction, and only then optimize for scale.

## Useful defaults

- Static frontends: evaluate Cloudflare Pages first
- Simple backends: start from Render if setup simplicity and a free-tier default matter most; compare against Cloud Run when lower latency or finer cloud control matters more
- Databases: keep TiDB Cloud Starter, Supabase, and Neon as named candidates, then re-check current free-tier terms before choosing
- Redis: keep Upstash as a named free-tier candidate when a hosted Redis-like service is needed

## Why it matters here

The workspace goal includes keeping hobby-project cost low. A reusable deployment page avoids redoing the same evaluation for every new web project.

## Why the concrete names stay

Even when pricing and quotas change, named candidates such as Render, Cloud Run, TiDB Cloud Starter, Supabase, Neon, and Upstash are useful retrieval anchors. They make future re-evaluation faster than starting again from generic searches like "free cloud DB".

## Related pages

- [cheap-hosting](../patterns/cheap-hosting.md)
