---
title: Build a 2D Space Shooter with Phaser and Claude Code
source_url: https://phaser.io/news/2026/02/phaser-claude-code-tutorial
source_type: video_backed_article
ingested_on: 2026-04-18
status: active
tags:
  - phaser
  - ai-game-dev
  - claude-code
  - browser-games
video_url: https://youtu.be/247Z3jdw_hs
video_platform: youtube
channel: Peter Yang
published_on: 2026-02-04
duration: 20m
time_anchors:
  - "00:00: why Peter Yang likes building games with his child using Claude Code"
  - "01:19: project setup"
  - "02:56: free pixel-art asset selection"
  - "05:33: drafting the spec with Claude Code's AskUserQuestion flow"
  - "08:56: MVP build and iteration"
  - "13:06: boss-battle debugging"
  - "15:24: shipping with GitHub and Vercel"
selected_screenshots: []
named_entities:
  - Phaser
  - Claude Code
  - Cursor
  - AskUserQuestion
  - Animus
  - Legacy Pixel Collection
  - GitHub
  - Vercel
  - Peter Yang
related_pages:
  - ../../wiki/tools/phaser.md
  - ../../wiki/topics/ai-game-dev.md
---

# Build a 2D Space Shooter with Phaser and Claude Code

## Why it matters here

This is a compact example of AI-assisted browser-game prototyping that starts from a human-readable game idea, turns it into a spec through questions, builds playable milestones, and ships a hosted demo.

## Summary

- Phaser's article highlights Peter Yang's 20-minute tutorial for building a retro vertical space shooter with Claude Code.
- The workflow uses Cursor, Claude Code, free pixel art assets, Claude Code's `AskUserQuestion` pattern for spec drafting, iterative playable milestones, and deployment through GitHub and Vercel.
- The resulting game includes scrolling shooter basics, enemy waves, power-ups, boss battles with health bars, screen shake, and debugging of a more complex boss encounter.
- The original Peter Yang post adds useful concrete anchors: a playable demo at `space-pixel-shooter.vercel.app`, the `Legacy Pixel Collection` from Animus, and a five-step tutorial sequence from setup to shipping.

## Workspace takeaways

- The reusable pattern is not "ask AI to make a game" in one prompt. It is: gather assets, have the agent ask design questions, write a spec, build a thin MVP, then add features in controlled milestones.
- Phaser remains a good target for fast browser-first game experiments because a shipped demo can be as simple as a Vercel deployment.
- For child-friendly or voice-driven co-creation, the article is a useful example of pairing an AI coding agent with dictated design intent rather than expecting the human collaborator to type or code.
- For workspace game projects, preserve this as a reference for milestone discipline: MVP first, then waves/power-ups, then boss mechanics and polish.

## Source access note

The Phaser page redirected direct browser fetches to a preview gate during ingest, but the article metadata and body summary were available through search indexing. The related Peter Yang post and Class Central entry corroborated the tutorial structure, video duration, author, and syllabus-level sequence.

## Follow-up

- If a workspace Phaser prototype starts, test whether this flow works better as a written `docs/specs/` questionnaire before implementation or as an interactive agent loop inside the coding tool.
- Capture the actual prompts and milestone boundaries if the team applies this to a local game project.
