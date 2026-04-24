---
title: 世界初!?バイブコーディング・縦スクロールシューティングを作った
source_url: https://note.com/shi3zblog/n/nc93e34da423f
source_type: article
original_language: ja
ingested_on: 2026-04-25
status: active
tags:
  - ai-game-dev
  - vertical-shooter
  - voice-ui
  - llm
published_on: 2026-04-11
named_entities:
  - GIKEN STAR
  - Bonsai
  - FLUX Klein 4B
  - Qwen3.5-9B
  - Suno
related_pages:
  - ../../wiki/topics/ai-game-dev.md
---

# shi3z on an AI-driven vertical shooter prototype

## Why it matters here

This note captures a concrete "AI game" interaction loop that is more ambitious than ordinary code generation. Voice input, realtime image generation, LLM-based rule interpretation, and generated boss behavior are all inside the live play loop.

## Summary

- The prototype `GIKEN STAR` is a vertical shooter controlled with lever movement, a shot button, and a bomb button tied to spoken keywords.
- When the player shouts a keyword during bomb activation, speech recognition feeds an LLM that decides the power-up, while the player's ship art and enemy imagery are regenerated in realtime to match the keyword.
- The article names `Bonsai` as the system that interprets spoken keywords and derives parameterized equipment behavior, so phrases like stronger or more specific variants can change the resulting loadout.
- `FLUX Klein 4B` is used for the regenerated player-ship visuals, and `Qwen3.5-9B` is used to generate attack programs, including boss attack behavior, with automatic retries on coding failures.
- The article also describes pulling fresh topics from internet news so early-stage enemies can reflect current themes.

## Workspace takeaways

- This is a strong retrieval anchor for "AI inside the gameplay loop" rather than AI only assisting development.
- Voice keyword design becomes part of game balance. The article frames successful play around inventing prompts that produce better equipment.
- For future prototypes, preserve the stack split explicitly: speech recognition, LLM interpretation, code generation, image generation, and gameplay fallback behavior when generation fails.

