# Beads Task Management Integration

## Objective

Replace the manual `.local/priority.md` task list with [beads](https://github.com/steveyegge/beads) (`bd` command) — a git-backed graph issue tracker designed for AI coding agents. This gives us persistent, dependency-aware task tracking without burning context window on full re-triage every session.

Addresses:
- Slow triage (subagent data collection for all repos every session)
- Loss of information (only top 5 kept, rest discarded)
- Manual priority file maintenance by LLM

## Changes

### `setup.sh`

Add beads installation via Homebrew:

```bash
if command -v brew &>/dev/null; then
    if ! command -v bd &>/dev/null; then
        echo "Installing beads..."
        brew install beads
    fi
fi
```

### Beads initialization

- Run `bd init` in the workspace root
- `git add .beads/` (git-tracking is the intended default)

### Claude Code plugin

Install as project-level plugin:

```
claude plugin install beads --scope project
```

This writes to `.claude/settings.json` (git-tracked) and provides SessionStart hooks that run `bd prime`.

### `triage-tasks` skill rewrite

Replace the current 8-step flow with a beads-centered approach:

**Step 0: Quick status** — Run `bd ready` to show actionable tasks. If tasks exist, present them and ask user to pick one, update the list, or do a full re-triage.

**Full re-triage (only when requested or when beads DB is empty):**
1. Collect data from repos (same subagent pattern as before)
2. For each discovered item, `bd create` with appropriate priority and tags
3. Set dependencies with `bd dep add` where applicable
4. Run `bd ready` to present the prioritized, unblocked list

**Minimal bd reference in the skill** (not a full manual, just what triage needs):
- `bd ready` — show tasks with no open blockers
- `bd list` — show all tasks
- `bd create "Title" -p <0-4>` — create task (0=critical, 4=low)
- `bd update <id> --claim` — claim and start a task
- `bd close <id> --reason "Done"` — close a completed task
- `bd dep add <child> <parent>` — add dependency

### Remove `.local/priority.md` workflow

- Delete `.local/priority.md`
- Remove references from triage-tasks skill
- Remove the `.local/` directory if empty

## Spec Changes

None — no `docs/specs/` file needed. The workflow changes live in the skill file and setup.sh.

## Sub-tasks

- [ ] [parallel] Update `setup.sh` with brew install for beads
- [ ] [parallel] Rewrite `triage-tasks` skill to use `bd` commands
- [ ] [depends on: setup.sh] Run `bd init` and commit `.beads/`
- [ ] [depends on: bd init] Install beads Claude Code plugin at project scope
- [ ] [depends on: all above] Remove `.local/priority.md` and clean up references
- [ ] [depends on: all above] Migrate existing priority items to beads
