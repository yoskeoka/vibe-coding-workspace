---
title: Knowledge Base Log
last_reviewed: 2026-04-24
status: seed
sources: []
---

# Knowledge Base Log

- 2026-04-11: seeded `docs/kb/` with schema, ingest guide, initial source notes, topic pages, project pages, tool pages, MkDocs config, and Pages publishing workflow
- 2026-04-11: re-ingested Karpathy's `LLM Knowledge Bases` X post, the `llm-wiki` gist, and the Classmethod article; added the original X post as a source note and updated compiled pages to preserve concrete anchors like `Obsidian`, `Obsidian Web Clipper`, `Mem0`, `pgvector`, `/kb-compile`, `index.md`, and `log.md`
- 2026-04-13: ingested Qiita article `趣味開発Webアプリケーションのほぼ0円インフラ構成`; added a source note and updated deployment pages with the concrete `Cloudflare DNS Proxy + Cloud Connector` -> `Cloud Run` / `Cloud Storage` routing pattern, `TiDB Serverless`, `Cloudflare WAF`, and explicit notes about cold-start, preview/beta maturity, and cross-cloud cost trade-offs
- 2026-04-16: refreshed `OpenAI built a Tactical RPG with Phaser` as a `video_backed_article`; preserved the canonical YouTube URL, `10:58-11:16` demo anchor, and the concrete workflow tuple `GPT-5.4 + Codex + Playwright Interactive + image generation + Phaser`
- 2026-04-18: ingested `Build a 2D Space Shooter with Phaser and Claude Code`; added a `video_backed_article` source note and updated Phaser / AI-assisted game development pages with the concrete workflow `Cursor + Claude Code + AskUserQuestion + Legacy Pixel Collection + GitHub + Vercel`
- 2026-04-22: normalized the `OpenAI built a Tactical RPG with Phaser` source note to the current video-backed format with visible source/video links, subtitle availability, selected screenshots, and anchor-aligned segment notes for the full embedded video.
- 2026-04-22: ingested Unity's `Esoteric Ebb: tackling interactive writing challenges`; added a source note and new `branching-narrative-authoring` pattern page with anchors for `Ink`, story-variable naming, choice feedback, text search, and pre-launch validation.
- 2026-04-24: added bilingual KB rendering with `docs/kb/ja/` mirror content, MkDocs static i18n publishing for `/` and `/ja/`, locale-aware source indexes, and English-only retrieval guidance.
