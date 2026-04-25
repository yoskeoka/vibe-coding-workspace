---
title: 趣味プロジェクト向けデプロイ選択肢
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-04-11-qiita-free-personal-deploy-options.md
  - ../../sources/2026/2026-04-13-qiita-near-zero-cost-webapp-infra.md
---

# 趣味プロジェクト向けデプロイ選択肢

## 現在の既定姿勢

まずコストの downside を最小化し、その次に setup friction を最小化し、最後に scale を最適化する順で deployment を選ぶ。

## 使いやすい既定値

- static frontend: まず Cloudflare Pages を評価する。
- simple backend: setup の簡単さと free-tier の既定値を重視するなら Render から始め、低 latency やより細かい cloud control を重視するなら Cloud Run と比較する。
- database: TiDB Cloud Starter、Supabase、Neon を具体候補として持ち、選ぶ直前に free-tier 条件を再確認する。
- Redis: hosted な Redis 互換が必要なら、free-tier 候補として Upstash を覚えておく。

## 覚えておくべき具体的な低コスト構成

1 つの domain、static frontend、container backend、hosted relational DB を必要とする hobby web app では、次の具体 stack を比較基準として残しておく。

- `Cloud Storage` で static file hosting
- minimum instance を `0` にした `Cloud Run` を backend に使う
- database に `TiDB Serverless`
- request routing に `Cloudflare DNS Proxy + Cloud Connector`
- bot 起因の backend 起床を減らすための `Cloudflare WAF`

この構成が重要なのは、単なる「安い host を選ぶ」メモより具体的であり、しかも `Cloud Load Balancing` コストを避けることを明示的に目標にしているから。

## 残しておくべきトレードオフ

- uptime check による `Cloud Run` keep-warm は workaround であり、保証された platform contract ではない。
- `Cloud Run domain mapping` は preview 段階として説明されており、setup 中に downtime が生じうる。
- `asia-northeast1` では domain mapping により latency が悪化する可能性がある。
- `Cloud Run` と `TiDB Serverless` の cross-cloud traffic は実コストになりうる。
- `Cloud Connector` の成熟度は導入前に再確認するべき。

## このワークスペースで重要な理由

ワークスペースの目標には hobby project のコストを低く保つことが含まれる。再利用可能な deployment page があれば、新しい web project のたびに同じ比較をやり直さずに済む。

## 具体名を残す理由

pricing や quota が変わっても、Render、Cloud Run、TiDB Cloud Starter、Supabase、Neon、Upstash のような具体候補は useful retrieval anchor である。generic な「無料 cloud DB」を毎回検索し直すより、既知候補から再評価したほうが速い。

## 関連ページ

- [cheap-hosting](../patterns/cheap-hosting.md)
