# Follow Up: Make Knowledge-Base Ingest Follow Branch Startup Rules

Status: resolved in the PR that records this issue by adding a preparation step to `skills/knowledge-base/SKILL.md`.

## Summary

Knowledge-base ingest almost always changes files under `docs/kb/`, and may also add follow-up notes under `docs/issues/`. The current `knowledge-base` skill focuses on ingest mechanics, but it does not explicitly require the same startup discipline used by `execute-task` and `review-task`:

- check whether the current working tree is clean
- avoid starting file-changing ingest directly on `main`
- create a fresh branch/worktree from latest `origin/main`
- use the globally installed `ww` CLI for the branch/worktree operation

This gap caused the GPT-5.4 Phaser ingest work to begin on `main` with file changes, then require a later stash/branch recovery.

## Why It Matters

- KB ingest is file-changing work, even when it feels like "just research".
- Starting on `main` risks mixing ingest edits with unrelated local changes.
- The workspace workflow expects every change, including doc-only and KB changes, to go through a clean branch and PR.
- `knowledge-base` should not be a hidden bypass around the workflow branch rules.

## Proposed Solution

- Update `skills/knowledge-base/SKILL.md` so file-changing ingest starts with the same branch hygiene as execution work:
  - run `git status --short --branch`
  - if dirty, stop or preserve the work before branching
  - create a fresh branch from latest `main` with `ww create docs/<name>` or `ww create chore/<name>` depending on scope
  - enter it with `ww cd`
- Keep read-only KB queries exempt when they do not write files.
- Clarify that ingest PRs are usually `docs/*` unless they include tooling/dependency changes, in which case `chore/*` is appropriate.
