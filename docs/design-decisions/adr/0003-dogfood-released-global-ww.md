# Dogfood the released global `ww` binary in normal workflow startup

## Status

Accepted — 2026-04-14

## Context

The workflow taught raw Git branch creation, bypassing the tool built for this
workspace's parallel multi-repository worktree flow. `ww` already supports
workspace targeting, centralized worktrees, path-oriented navigation, and
git-native behavior.

## Decision

Use the released global `ww` binary by default for normal planning and
execution. Start with `ww create`, enter with `ww cd`, reserve in-repo builds
for unreleased `ww` work, and record normal-workflow `ww` failures rather than
silently bypassing them.

## Consequences

- Everyday work continuously dogfoods the released `ww` UX.
- Parallel tasks become the normal model.
- Workflow docs and skills must stay aligned on `ww` startup.
- Documented fallback remains necessary for `ww` debugging or blocked recovery.
