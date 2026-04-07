# Architectural Decision Records (ADR)

## 2026-04-07 Skip child-repo workflow PRs when `skills/` is unchanged

### Context
`sync-workflow-to-child-repos.yml` currently opens PRs in every child repository whenever `vibe-coding-workspace` receives a push to `main`. However, `setup-workspace.sh` only symlinks skills from the workflow submodule, so changes outside `skills/` do not affect child repos.

### Decision
The sync workflow will compare the child repo's recorded workflow submodule commit against the pushed workflow commit and only create a PR when the diff contains changes under `skills/`.

### Consequences
- Child repos stop receiving no-op workflow PRs for unrelated repo changes.
- The workflow becomes slightly more selective, so future changes that should sync to child repos must live under `skills/` or update the sync rule explicitly.

---
