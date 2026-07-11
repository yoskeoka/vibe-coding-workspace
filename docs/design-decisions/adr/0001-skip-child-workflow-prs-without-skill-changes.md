# Skip child-repo workflow PRs when `skills/` is unchanged

## Status

Accepted — 2026-04-07

## Context

`sync-workflow-to-child-repos.yml` opened PRs in every child repository for
every workspace `main` push, but `setup-workspace.sh` symlinks only `skills/`
from the workflow submodule. Changes outside `skills/` do not affect children.

## Decision

Compare each child's recorded workflow commit with the pushed commit and create
a sync PR only when the diff includes `skills/` changes.

## Consequences

- Children no longer receive no-op PRs for unrelated workspace changes.
- A change intended for children must live under `skills/` or explicitly update
  the sync rule.
