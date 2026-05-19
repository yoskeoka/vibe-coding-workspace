---
title: Cheap Hosting
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-04-11-qiita-free-personal-deploy-options.md
  - ../../sources/2026/2026-04-13-qiita-near-zero-cost-webapp-infra.md
  - ../../sources/2026/2026-05-17-zenn-remotestorage-user-owned-webapps.md
  - ../../sources/2026/2026-05-13-zenn-low-cost-minecraft-server-aws.md
---

# Cheap Hosting

## Pattern

For hobby projects, default to hosting combinations that preserve an easy path to launch while keeping financial downside near zero.

## Heuristics

- prefer managed services with durable free tiers for first public launches
- preserve concrete candidate names such as Cloudflare Pages, Render, Cloud Run, TiDB Cloud Starter, Supabase, Neon, and Upstash so future re-checks start from known options
- record price assumptions with dates
- separate "best DX" from "best default for a hobby budget"
- if the app only needs user-owned personal data, test whether `remoteStorage` can remove the backend entirely before comparing hosted stacks
- keep a concrete multi-provider fallback pattern in mind when one cheap service is not enough: `Cloudflare DNS Proxy + Cloud Connector` in front of `Cloud Run` and `Cloud Storage`, with `TiDB Serverless` as the DB
- keep a separate pattern for bursty multiplayer backends: on-demand `EC2` spot start via an event trigger, then pay mainly for session time rather than idle uptime
- treat cold-start mitigation tricks such as uptime-check keep-warm as opportunistic hacks and re-validate them before depending on them
- preserve beta or preview risk in the notes; low cost often comes from accepting maturity trade-offs, not just lower pricing

## Related pages

- [deployment-options](../topics/deployment-options.md)
- [user-owned-storage-webapps](./user-owned-storage-webapps.md)
