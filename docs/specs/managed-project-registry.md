# Spec: Managed Project Registry

## Goal

Define how the workspace tracks the set of child repositories that it manages.

## Scope

This spec covers the workspace-level registry only:
- which repositories are managed by the workspace
- how the registry is reflected in scripts and documentation
- how new child repositories are added to the workspace metadata

## Requirements

### 1. Registry source of truth
`setup.sh` MUST contain the canonical list of managed repository URLs.

### 2. Workspace metadata sync
Whenever a repository is added to the managed set, the workspace MUST update the following in the same change set:
- `.gitignore`
- `README.md` managed-project list
- `docs/project-plan.md` managed-project table
- `AGENTS.md` workspace structure list when it describes child projects

### 3. Local checkout behavior
Running `setup.sh` MUST clone or update each managed repository under the workspace root using the repository name as the directory name.

### 4. Ignored child checkouts
Child repository working trees MUST be ignored by Git from the workspace root so nested repositories are not tracked accidentally.

## Current Managed Projects

The workspace currently manages:
- `reversi-adventure`
- `ai-arena`
- `vim-learning-game`
- `ww`
- `homebrew-ww`

## Non-Goals

- Project-specific requirements for any child repository
- The contents of child repositories themselves
- New project intake decision-making
