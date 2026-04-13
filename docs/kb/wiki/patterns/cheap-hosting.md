---
title: Cheap Hosting
last_reviewed: 2026-04-13
status: seed
sources:
  - ../../sources/2026/2026-04-11-qiita-free-personal-deploy-options.md
  - ../../sources/2026/2026-04-13-qiita-near-zero-cost-webapp-infra.md
---

# Cheap Hosting

## Pattern

For hobby projects, default to hosting combinations that preserve an easy path to launch while keeping financial downside near zero.

## Heuristics

- prefer managed services with durable free tiers for first public launches
- preserve concrete candidate names such as Cloudflare Pages, Render, Cloud Run, TiDB Cloud Starter, Supabase, Neon, and Upstash so future re-checks start from known options
- record price assumptions with dates
- separate "best DX" from "best default for a hobby budget"
- keep a concrete multi-provider fallback pattern in mind when one cheap service is not enough: `Cloudflare DNS Proxy + Cloud Connector` in front of `Cloud Run` and `Cloud Storage`, with `TiDB Serverless` as the DB
- treat cold-start mitigation tricks such as uptime-check keep-warm as opportunistic hacks and re-validate them before depending on them
- preserve beta or preview risk in the notes; low cost often comes from accepting maturity trade-offs, not just lower pricing

## Related pages

- [deployment-options](../topics/deployment-options.md)
