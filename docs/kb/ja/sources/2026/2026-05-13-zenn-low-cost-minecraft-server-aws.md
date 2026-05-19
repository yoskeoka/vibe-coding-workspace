---
title: 変動費2.3円/時、固定費22円/月のマイクラサーバーを構築した話
source_url: https://zenn.dev/daikitchen/articles/ac794d03b9baf3
source_type: article
original_language: ja
ingested_on: 2026-05-19
status: watch
tags:
  - deployment
  - aws
  - game-server
  - iac
related_pages:
  - ../../wiki/topics/deployment-options.md
  - ../../wiki/patterns/cheap-hosting.md
---

# 変動費2.3円/時、固定費22円/月のマイクラサーバーを構築した話

## ここで重要な理由

これは既存の「常時稼働する hobby web app」の低コスト構成とは別系統の、「遊ぶときだけ起動する game server」パターンとして保持する価値がある。

## 要約

- 記事は「とにかくコストを削る」と「すべて IaC で管理する」を両立させる構成を狙っている。
- 残すべき cost anchor はかなり具体的である。
- 固定費は `AMI + Snapshot` で月 `22円` 前後
- 変動費は `t4g.medium` の spot と小さな `EBS` で `2.3円/時` 前後
- architecture も覚えやすく、`Discord` からの操作が `Lambda` を起動し、最安 availability zone の `EC2` spot instance を立ち上げ、container 化された Minecraft workload を走らせ、停止時に `EBS Snapshot` へ state を戻す。
- server は基本的に停止前提で、instance を落としている間の継続コストを極小化している。
- `AWS`、`Terraform`、`Ansible`、`Packer` という tag 群も、価格だけでなく repeatable provisioning の資料として重要な retrieval anchor になる。

## ワークスペースでの含意

- これは generic な cheap hosting ではなく、「ephemeral multiplayer server」の具体比較対象として残すとよい。
- 有効なのは、session が断続的で user-triggered な場合であり、常時オンライン性が前提の service ではない。
- reuse する前には、managed dedicated server や常時稼働 VPS と比べて、spot availability と state snapshot 運用の trade-off を確認する必要がある。

## フォローアップ

- 将来 low-duty-cycle な backend が必要になったら、serverless な turn-based design や static-shareable artifact と並べて比較する。
