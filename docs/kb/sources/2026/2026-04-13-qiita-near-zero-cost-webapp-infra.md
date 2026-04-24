---
title: 趣味開発Webアプリケーションのほぼ0円インフラ構成
source_url: https://qiita.com/mazrean/items/f4a48d43b2d680a92216
source_type: article
original_language: ja
ingested_on: 2026-04-13
status: watch
tags:
  - deployment
  - personal-dev
  - cloudflare
  - cloud-run
  - tidb
related_pages:
  - ../../wiki/topics/deployment-options.md
  - ../../wiki/patterns/cheap-hosting.md
---

# 趣味開発Webアプリケーションのほぼ0円インフラ構成

## Why it matters here

This is a stronger concrete reference than a generic "cheap hosting" article because it describes a specific multi-service architecture for low-traffic hobby web apps.

## Summary

- Targets a small web app shape: static frontend by `SSG + CSR`, containerized backend, and a database-backed API.
- Uses `TiDB Serverless` on `TiDB Cloud` as the database anchor because the article reports a relatively large free tier for small services.
- Uses `Cloud Run` with minimum instances set to `0` so backend cost stays near the free tier when traffic is sparse.
- Uses `Cloud Monitoring` uptime checks as a practical keep-warm tactic to reduce `Cloud Run` cold starts, while explicitly noting that this behavior is not guaranteed long-term.
- Uses `Cloud Storage` for static asset delivery instead of a heavier frontend hosting setup.
- Avoids `Cloud Load Balancing` cost by combining `Cloudflare` DNS proxy mode with `Cloud Connector` to route requests between `Cloud Run` and `Cloud Storage`.
- Calls out important trade-offs instead of pretending the stack is universally good:
  - `Cloud Run domain mapping` is preview-stage
  - domain cutover can require downtime during certificate issuance and DNS propagation
  - `asia-northeast1` can see noticeable latency increase with domain mapping
  - cross-cloud traffic between `Cloud Run` and `TiDB Serverless` can create network charges
  - `Cloud Connector` is beta
- Adds `Cloudflare WAF` as a cost-protection layer to reduce useless bot traffic that would otherwise wake `Cloud Run`.

## Workspace takeaways

- Keep this article as the concrete example for "multi-provider, low-cost hobby web app" rather than only storing service-by-service recommendations.
- Treat `Cloudflare DNS Proxy + Cloud Connector` as a concrete alternative to paying for `Cloud Load Balancing` when the app can tolerate preview or beta features and some operational rough edges.
- Treat keep-warm uptime checks as a hack, not a platform guarantee.
- Re-check pricing, quota, and feature maturity before using this stack for any real deployment decision.

## Follow-up

- If a future workspace project needs a hobby-grade web stack, compare this exact pattern against `Cloudflare Pages + Render` and a simpler single-provider setup before choosing.
