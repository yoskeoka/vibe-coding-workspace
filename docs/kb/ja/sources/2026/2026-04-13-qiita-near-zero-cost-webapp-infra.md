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

# 趣味開発 Web アプリケーションのほぼ 0 円インフラ構成

## ここで重要な理由

これは generic な「安い hosting」記事より強い具体資料で、低トラフィックな趣味 web app 向けに、実際の multi-service 構成を示しています。

## 要約

- 小さな web app の形として、`SSG + CSR` の static frontend、container 化 backend、database-backed API を前提にしている。
- DB の中心には `TiDB Cloud` 上の `TiDB Serverless` を置き、小規模サービスには比較的大きな free tier があると報告している。
- backend は `Cloud Run` で minimum instance を `0` にし、アクセスが少ないときのコストを free tier 近辺に抑える。
- `Cloud Monitoring` の uptime check を keep-warm 策として使い、cold start を和らげる一方で、長期的に保証された挙動ではないと明記している。
- static asset 配信には、重い frontend hosting ではなく `Cloud Storage` を使っている。
- `Cloud Load Balancing` のコストを避けるため、`Cloudflare` DNS proxy mode と `Cloud Connector` を組み合わせて `Cloud Run` と `Cloud Storage` にルーティングしている。
- 普遍的に良いと装わず、重要な trade-off も列挙している。
- `Cloud Run domain mapping` は preview 段階
- 証明書発行と DNS 反映の間に downtime が必要になりうる
- `asia-northeast1` では domain mapping による latency 悪化がありうる
- `Cloud Run` と `TiDB Serverless` の cross-cloud traffic は network charge を生みうる
- `Cloud Connector` は beta
- さらに `Cloudflare WAF` を cost-protection layer として置き、不要な bot traffic で `Cloud Run` が起きるのを抑える。

## ワークスペースでの含意

- これは「service ごとの推奨一覧」ではなく、「multi-provider で低コストな趣味 web app」の具体例として保持する価値がある。
- `Cloudflare DNS Proxy + Cloud Connector` は、preview / beta 機能や多少の rough edge を許容できるなら、`Cloud Load Balancing` を払わずに済む具体的代替として覚えておける。
- keep-warm 用 uptime check は hack であり、platform guarantee とみなさないほうがよい。
- 実運用に使う前には、pricing、quota、feature maturity を再確認する必要がある。

## フォローアップ

- 将来 hobby-grade な web stack が必要になったら、この構成を `Cloudflare Pages + Render` や単一 provider のより単純な構成と並べて比較する。
