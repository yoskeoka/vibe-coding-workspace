# Issue: Sync Child Repos To Latest Slopless Workflow

## Summary

The workspace `slopless` workflow gained follow-up fixes after the first child-repo rollout merged.

The merged child-repo PRs still carry the earlier workflow shape. They should be updated to match the workspace version.

## Why This Matters

- child repos still rely on the runner-default Node toolchain
- child repos still invoke `npx` once per file instead of once per changed file set
- child repos still read only the first page of PR comments before upserting the marker comment
- child repos still keep broader-than-needed workflow permissions
- child repos still lack explicit subprocess and GitHub API timeouts

## Target Repositories

- `ai-arena`
- `dungeon-game-ai-arena`
- `envdiff`
- `reversi-adventure`
- `reversi-ai-arena`
- `vim-learning-game`
- `ww`

## Expected Update

Port the current workspace `.github/workflows/slopless.yml` behavior to each target repo without expanding scope beyond the workflow file unless a repo-specific compatibility fix is required.
