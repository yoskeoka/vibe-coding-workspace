---
title: AI-Assisted Game Development
last_reviewed: 2026-04-18
status: seed
sources:
  - ../../sources/2026/2026-04-11-openai-gpt-5-4-phaser-tactical-rpg.md
  - ../../sources/2026/2026-04-18-phaser-claude-code-space-shooter.md
---

# AI-Assisted Game Development

## Current signal

Recent examples point to a workflow where an AI agent helps with gameplay systems, browser playtesting, and asset generation in one loop.

## Why it matters here

- Game projects are a major part of the workspace portfolio.
- Agent-driven prototyping can keep hobby projects moving with limited human time.
- The Phaser tactical RPG reference is thin on internals, but it is still a useful proof point for the loop `code generation -> browser playtest -> UI/combat iteration -> asset generation`.
- Peter Yang's Claude Code space-shooter tutorial adds a more procedural pattern: use an agent to ask design questions, draft a concrete game spec, build an MVP, then iterate through feature milestones before deploying.

## Reusable pattern

- Start from a small genre with obvious mechanics, such as a vertical shooter.
- Add or select assets before implementation so the agent can map actual filenames to player, enemy, projectile, power-up, and background roles.
- Have the agent ask clarifying questions before writing the spec.
- Treat the first build as a playable MVP, then add features in narrow passes such as enemy waves, power-ups, boss health bars, screen shake, and deployment.
- Keep deployment in the loop early; a static Phaser build hosted on Vercel is enough for sharing and review.

## Open questions

- Which parts of the loop should be automated first: prototype generation, playtesting, or art support?
- Which project should become the first serious AI-assisted game-production testbed?
- What is the minimum artifact set worth preserving from these demos: code, prompts, playtest traces, or curated screenshots?
- Should workspace game plans include an explicit "agent asks design questions before spec" checkpoint?

## Related pages

- [phaser](../tools/phaser.md)
- [reversi-adventure](../projects/reversi-adventure.md)
