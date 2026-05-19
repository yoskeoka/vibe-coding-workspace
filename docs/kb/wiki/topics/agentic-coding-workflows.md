---
title: Agentic Coding Workflows
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-05-19-agent-quality-controls-slopless.md
  - ../../sources/2026/2026-05-11-zenn-claude-code-local-operations-guide.md
---

# Agentic Coding Workflows

## Current signal

The most useful workflow notes are no longer only about prompting. They are about runtime placement, privacy boundaries, deterministic guardrails, and how agents fit inside a repeatable loop.

## Why it matters here

- This workspace already uses agent-first development, so operational guidance has direct reuse value.
- The same repo may need opposite capabilities at different times: strong local privacy boundaries on one task, and strong prose-quality enforcement on another.
- These sources are useful because they preserve concrete commands, environment variables, and integration loops rather than generic "AI coding is useful" commentary.

## Reusable pattern

- Treat agent workflow design as a systems problem:
  - choose where the model runs
  - choose which outputs need deterministic review
  - keep the invocation contract simple enough to repeat
- For privacy-sensitive work, redirect the agent CLI to a local `Messages API`-compatible endpoint rather than assuming the hosted default is always acceptable.
- For prose or reviewable Markdown, prefer deterministic lint feedback such as `slopless` when the goal is repeatable cleanup rather than open-ended rewriting.
- Preserve raw review artifacts such as JSON findings when they help post-run inspection or CI integration.

## Concrete anchors worth preserving

- `Claude Code`
- `LM Studio`
- `Qwen3-Coder-30B-A3B-Instruct`
- `ANTHROPIC_BASE_URL`
- `slopless`
- `textlint`
- `.slopless/findings/`

## Open questions

- Which repos in this workspace actually justify local-only agent routing despite the capability drop?
- Where should deterministic writing checks run: pre-commit, CI, or only inside agent loops?
- Is there a narrow English-doc subset where `slopless` gives enough value to justify adding another lint surface?

## Related pages

- [llm-knowledge-bases](./llm-knowledge-bases.md)
- [slopless](../tools/slopless.md)
