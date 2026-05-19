---
title: slopless
source_url: https://github.com/agent-quality-controls/slopless
source_type: repo
original_language: en
ingested_on: 2026-05-19
status: active
tags:
  - writing
  - lint
  - ai-coding
  - markdown
related_pages:
  - ../../wiki/topics/agentic-coding-workflows.md
  - ../../wiki/tools/slopless.md
---

# slopless

## Why it matters here

This is a concrete quality-control tool for AI-assisted writing loops: deterministic checks, structured findings, and an explicit agent integration path.

## Summary

- `slopless` is a Markdown prose checker built on deterministic `textlint` rules rather than LLM calls.
- The repo describes it as catching both AI and human prose slop in English Markdown, with more than 50 rules and JSON-only output.
- The intended loop is concrete: install the package, install the Codex or Claude skill, start a fresh agent session, and keep rewriting until the findings JSON is clean.
- The CLI contract is worth preserving:
  - `npx slopless "docs/**/*.md"` checks files or globs
  - exit `0` means clean
  - exit `1` means findings
  - exit `2` means failure or invalid invocation
- The repo recommends saving raw findings JSON under `.slopless/findings/` so review artifacts remain inspectable outside the agent loop.
- The repository wiki pointers are useful retrieval anchors for later evaluation: `Philosophy`, `Comparison`, `Rules`, `Behavior`, and `Ignore-Rules`.

## Workspace takeaways

- This is a candidate guardrail when English docs need repeatable prose cleanup without paying for another model call.
- The deterministic JSON contract makes it a better fit for CI or scripted review loops than vague style guidance.
- Its scope is explicitly English-only, so it should not be assumed to cover this workspace's Japanese internal docs.

## Follow-up

- If the workspace wants prose QA for English KB pages or PR-facing docs, compare `slopless` against existing `Vale`, `textlint`, or review-only workflows.
