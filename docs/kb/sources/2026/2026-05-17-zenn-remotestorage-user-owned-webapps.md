---
title: サーバー代0円・バックエンド不要！？remoteStorageで実現する「ユーザー主権」なWebアプリの実装パターン
source_url: https://zenn.dev/konnyaku256/articles/web-app-impl-patttern-using-remote-storage
source_type: article
original_language: ja
ingested_on: 2026-05-19
status: active
tags:
  - deployment
  - web-app
  - remotestorage
  - privacy
related_pages:
  - ../../wiki/topics/deployment-options.md
  - ../../wiki/patterns/user-owned-storage-webapps.md
---

# サーバー代0円・バックエンド不要！？remoteStorageで実現する「ユーザー主権」なWebアプリの実装パターン

## Why it matters here

This is a useful counterexample to the usual "cheap backend" discussion because it removes the backend entirely for a specific class of user-data-centric web apps.

## Summary

- The article frames `remoteStorage` as an open protocol for web apps where users keep ownership of their own data instead of handing storage custody to the app developer.
- The implementation model splits app and storage cleanly: the frontend app asks the user to connect a storage endpoint, then reads and writes directly to that user-controlled backend.
- The protocol stack is concrete and worth preserving: `WebFinger` for endpoint discovery, `OAuth2` for authorization, and a REST-like HTTP API for file and directory operations.
- The practical JavaScript anchor is `remotestorage.js`, using access scopes such as `notes` plus a built-in connection widget and client methods like `storeFile()` and `getFile()`.
- The article calls out concrete provider anchors: `5apps`, `Nextcloud` with a plugin, and `Armadietto`, while also mentioning adapter-style use with `Dropbox` or `Google Drive`.
- The strongest benefits are zero backend hosting cost for the app developer, no app-owned database to operate, and a privacy story that is simpler to explain because user data is not stored on the app author's server.

## Workspace takeaways

- Keep this as the "user-owned storage" option when a project mostly stores personal notes, progress, settings, or other user-scoped data.
- This is not a generic replacement for server-side systems. It fits best when shared authoritative state, heavy server computation, and multi-user coordination are not core requirements.
- When a hobby web app can live inside `static frontend + user-chosen storage`, this may be a better comparison point than adding the cheapest possible backend stack.

## Follow-up

- If a future project wants sync without app-managed infra, compare `remoteStorage` against local-only persistence and against thin hosted backends before choosing.
