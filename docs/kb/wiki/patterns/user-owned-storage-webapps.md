---
title: User-Owned Storage Web Apps
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-05-17-zenn-remotestorage-user-owned-webapps.md
---

# User-Owned Storage Web Apps

## Pattern

For user-centric web apps that mostly store personal data, treat "the app does not own the backend" as a real architecture option rather than assuming every persistent app needs its own database and API service.

## Heuristics

- Prefer this pattern when data is primarily user-scoped notes, settings, progress, or other portable personal state.
- Preserve the concrete `remoteStorage` anchors:
  - `WebFinger` for discovery
  - `OAuth2` for authorization
  - REST-like file operations
  - `remotestorage.js` for browser implementation
- Record provider and ecosystem names such as `5apps`, `Nextcloud`, and `Armadietto` so future evaluation starts from concrete options.
- Compare this against local-only persistence first, then against the cheapest hosted backend. The point is not that user-owned storage always wins, but that it sometimes removes an entire class of infra obligations.
- Reject this pattern early when the product needs shared authoritative state, server-side jobs, or strong multi-user coordination.

## Why it matters here

This workspace has several hobby-project lines where infra cost and operational burden are real constraints. A "no app-owned backend" option deserves its own retrieval anchor.

## Related pages

- [deployment-options](../topics/deployment-options.md)
- [cheap-hosting](./cheap-hosting.md)
