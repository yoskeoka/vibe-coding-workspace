---
title: LLM Knowledge Bases
source_url: https://x.com/karpathy/status/2039805659525644595
source_type: post
original_language: en
ingested_on: 2026-04-11
status: active
tags:
  - knowledge-base
  - llm
  - x
related_pages:
  - ../../wiki/topics/llm-knowledge-bases.md
  - ../../wiki/patterns/source-ingestion.md
  - ../../wiki/projects/ww.md
---

# LLM Knowledge Bases

## Why it matters here

This is the original framing post that named the pattern. It matters less as a full spec than as the motivation: spend LLM token budget on maintaining knowledge, not only on generating one-off answers or code.

## Summary

- Karpathy frames the shift as using LLMs to build personal knowledge bases for active research topics.
- The memorable claim is that more of his recent token throughput is going into manipulating knowledge instead of manipulating code.
- The post functions as an announcement and pointer to the follow-up gist, where the file layout and operating loop are spelled out in more detail.
- The surrounding discussion makes the pattern legible as a small-scale alternative to reaching for full RAG infrastructure too early.

## Workspace takeaways

- Keep this post as the "why now" source for the workspace KB.
- Use the gist and article for implementation detail, but keep the X post linked as the origin point for the pattern name.

## Follow-up

- When explaining the KB to future collaborators, cite this post together with the gist instead of only the gist.
