---
title: slopless
last_reviewed: 2026-05-19
status: seed
sources:
  - ../../sources/2026/2026-05-19-agent-quality-controls-slopless.md
---

# slopless

## Current signal

`slopless` is a deterministic English-Markdown review tool aimed at catching vague, padded, or repetitive prose patterns without calling another model.

## Notes

- It ships as a CLI plus agent-installable skills for `Codex` and `Claude Code`.
- The JSON-only output shape is the main operational anchor because it makes the tool scriptable and inspectable.
- The repo's intended loop is specifically agent-oriented: install the skill, start a fresh session, run the checker, rewrite, and repeat until the findings are gone.
- The repository also exposes a broader evaluation surface through its wiki: rule inventory, philosophy, comparison, behavior, and ignore syntax.
- Its English-only scope matters. It should not be treated as a workspace-wide docs linter without a language-boundary decision first.

## Related pages

- [agentic-coding-workflows](../topics/agentic-coding-workflows.md)
