---
title: 低コストホスティング
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-04-11-qiita-free-personal-deploy-options.md
  - ../../sources/2026/2026-04-13-qiita-near-zero-cost-webapp-infra.md
  - ../../sources/2026/2026-05-17-zenn-remotestorage-user-owned-webapps.md
  - ../../sources/2026/2026-05-13-zenn-low-cost-minecraft-server-aws.md
---

# 低コストホスティング

## パターン

趣味 project では、まず公開までの道筋を単純に保ちつつ、金銭的な downside をほぼゼロに抑える hosting の組み合わせを既定にする。

## ヒューリスティクス

- 最初の公開では、耐久的な free tier を持つ managed service を優先する。
- Cloudflare Pages、Render、Cloud Run、TiDB Cloud Starter、Supabase、Neon、Upstash のような具体名を残し、次回の再確認を既知候補から始められるようにする。
- 価格前提には日付を添える。
- 「最高の DX」と「趣味予算に対する最良の既定値」を分けて考える。
- app が user-owned personal data だけで成立するなら、hosted stack を比べる前に `remoteStorage` で backend 自体を消せないか試す。
- 単一サービスでは足りないときの具体 fallback として、`Cloudflare DNS Proxy + Cloud Connector` を front に置き、`Cloud Run` と `Cloud Storage`、DB に `TiDB Serverless` を組み合わせる構成を覚えておく。
- bursty な multiplayer backend には、event trigger で on-demand `EC2` spot を起動し、idle uptime ではなく session 時間に主に課金される構成も別系統として覚えておく。
- uptime check による cold-start 緩和のような工夫は opportunistic hack と見なし、依存前に再検証する。
- beta / preview のリスクもメモに残す。低コストは単に価格が安いだけでなく、成熟度の trade-off を受け入れて成り立つことが多い。

## 関連ページ

- [deployment-options](../topics/deployment-options.md)
