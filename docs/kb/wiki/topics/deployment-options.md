---
title: Deployment Options for Hobby Projects
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-04-11-qiita-free-personal-deploy-options.md
  - ../../sources/2026/2026-04-13-qiita-near-zero-cost-webapp-infra.md
  - ../../sources/2026/2026-05-17-zenn-remotestorage-user-owned-webapps.md
  - ../../sources/2026/2026-05-13-zenn-low-cost-minecraft-server-aws.md
---

# Deployment Options for Hobby Projects

## Current default posture

Prefer deployment choices that minimize cost risk first, then minimize setup friction, and only then optimize for scale.

## Useful defaults

- Static frontends: evaluate Cloudflare Pages first
- Simple backends: start from Render if setup simplicity and a free-tier default matter most; compare against Cloud Run when lower latency or finer cloud control matters more
- Databases: keep TiDB Cloud Starter, Supabase, and Neon as named candidates, then re-check current free-tier terms before choosing
- Redis: keep Upstash as a named free-tier candidate when a hosted Redis-like service is needed
- User-owned sync: if the app only needs per-user persistence, compare `remoteStorage` before adding any app-owned backend at all

## Concrete low-cost stack to remember

When a hobby web app needs one domain, a static frontend, a container backend, and a hosted relational DB, keep this specific stack as a comparison point:

- `Cloud Storage` for static file hosting
- `Cloud Run` for the backend with minimum instances set to `0`
- `TiDB Serverless` for the database
- `Cloudflare DNS Proxy + Cloud Connector` for request routing
- `Cloudflare WAF` to reduce bot-triggered backend wakeups

This matters because it is a more concrete architecture than a generic "pick a cheap host" note, and it explicitly aims to avoid `Cloud Load Balancing` cost.

## Concrete zero-backend stack to remember

When the product mostly stores user-owned personal data and does not need authoritative shared state, keep `remoteStorage` in the option set:

- browser app plus `remotestorage.js`
- endpoint discovery through `WebFinger`
- authorization through `OAuth2`
- user-chosen providers such as `5apps`, `Nextcloud`, or `Armadietto`

This matters because it changes the decision from "which backend should we host?" to "can this product avoid owning a backend entirely?"

## Concrete bursty game-server stack to remember

When a multiplayer game server only needs to exist during actual play sessions, keep this AWS pattern in mind:

- `Discord -> Lambda -> EC2 spot instance`
- containerized game workload on the instance
- state persistence through `AMI + EBS Snapshot`
- infrastructure managed through `Terraform`, `Ansible`, and `Packer`

This matters because the pricing model is usage-shaped rather than always-on: fixed storage cost stays small, while compute cost only accrues during active sessions.

## Trade-offs worth preserving

- `Cloud Run` keep-warm via uptime checks is a workaround, not a guaranteed platform contract
- `Cloud Run domain mapping` was described as preview-stage and can involve downtime during setup
- `asia-northeast1` latency can worsen with domain mapping
- cross-cloud traffic between `Cloud Run` and `TiDB Serverless` can turn into real egress cost
- `Cloud Connector` maturity should be re-checked before adopting it
- `remoteStorage` is a poor fit for shared server authority, coordinated multiplayer, or heavy backend jobs
- bursty `EC2` spot designs trade away simplicity and uptime guarantees in exchange for lower cost

## Why it matters here

The workspace goal includes keeping hobby-project cost low. A reusable deployment page avoids redoing the same evaluation for every new web project.

## Why the concrete names stay

Even when pricing and quotas change, named candidates such as Render, Cloud Run, TiDB Cloud Starter, Supabase, Neon, and Upstash are useful retrieval anchors. They make future re-evaluation faster than starting again from generic searches like "free cloud DB".

## Related pages

- [cheap-hosting](../patterns/cheap-hosting.md)
- [user-owned-storage-webapps](../patterns/user-owned-storage-webapps.md)
