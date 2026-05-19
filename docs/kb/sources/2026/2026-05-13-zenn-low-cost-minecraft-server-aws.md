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

## Why it matters here

This is a concrete "only pay while playing" game-server deployment pattern, which is different from the always-on hobby web app references already in the KB.

## Summary

- The article optimizes for two goals at once: very low cost and full infrastructure-as-code control.
- The remembered cost anchors are strong:
  - fixed cost about `22 yen/month` for `AMI + Snapshot`
  - variable cost about `2.3 yen/hour` for `t4g.medium` spot plus small `EBS`
- The architecture is event-driven and worth preserving: player action through `Discord` triggers `Lambda`, which starts the cheapest availability-zone `EC2` spot instance, runs the Minecraft workload in containers, and saves state back to `EBS Snapshot` on stop.
- The article explicitly treats the server as mostly-off infrastructure: when the instance is stopped, ongoing cost is close to zero.
- The tag set `AWS`, `Terraform`, `Ansible`, and `Packer` is a useful retrieval anchor because the value is not only the price point but also the repeatable provisioning posture.

## Workspace takeaways

- Keep this as the concrete comparison point for "ephemeral multiplayer server" rather than mixing it into generic cheap hosting advice.
- The pattern is most useful when sessions are bursty and user-triggered, not when a service must stay online continuously.
- This should be compared against managed dedicated-server offerings and against simpler always-on VPS setups before reuse, because spot availability and state-snapshot handling are real trade-offs.

## Follow-up

- If a future game project needs a low-duty-cycle backend, compare this event-driven start/stop pattern against serverless turn-based designs and static-shareable game artifacts.
