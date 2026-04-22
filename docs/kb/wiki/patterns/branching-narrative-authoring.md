---
title: Branching Narrative Authoring
last_reviewed: 2026-04-22
status: seed
sources:
  - ../../sources/2026/2026-04-22-unity-interactive-writing-esoteric-ebb.md
---

# Branching Narrative Authoring

## Current signal

Large nonlinear RPG writing succeeds when the authoring loop stays close to the speed of ordinary writing while still preserving player-state feedback. The `Esoteric Ebb` postmortem is useful because it connects narrative ambition to concrete authoring mechanics: `Ink`, Unity integration, text search, naming conventions, and late bug-fix pressure.

## Why it matters here

Workspace game projects may eventually include dialogue, choice, tutorials, or branching scenario content. Even if they are not CRPGs, they can inherit the same failure mode: choices are easy to add locally, but hard to validate globally.

## Reusable pattern

- Define the interaction contract before writing content: choice density, dynamic-content target, and whether the design is authored branching or systemic emergence.
- Pick tools that make branching cheap. Important properties are plain-text review, fast recompilation, project-wide search, Unity or engine integration, and easy custom commands.
- Treat narrative state as an explicit data model. Use stable prefixes for area, quest, quest-point, character, or tutorial state so variables can be found and audited.
- Require visible feedback for important choices. A variable that never changes later text, quest state, UI, or ending logic is probably narrative debt.
- Add validation before content scale grows: syntax checks, unused-variable reports, branch reachability checks, and playtest scripts for high-risk paths.

## Tool anchors

- `Ink`: plain-text interactive fiction scripting with strong Unity integration.
- `Yarn Spinner`: dialogue system often used with Unity projects.
- `Arcweave`: visual and text-oriented interactive narrative planning.
- `articy:draft`: narrative-design and game-content planning suite.

## Open questions

- Should workspace game specs include a standard "narrative state and validation" section when choices or tutorials branch?
- Which validation checks can be generated cheaply from `Ink` or `Yarn Spinner` project files?
- Is a visual planning tool useful here, or does plain text plus search fit the workspace better?
