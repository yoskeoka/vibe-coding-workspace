---
name: triage-tasks
description: When starting a new session, reviewing priorities, running a daily briefing, triaging pending tasks across workspace and child projects, or deciding what to work on next.
metadata:
  author: yoskeoka
  version: '2.0.0'
---

# Triage Tasks — Daily Briefing (Workspace-Only)

**Position in workflow**: This is a **session-start ritual**, not a numbered workflow step. Run it at the beginning of a work session to decide what to tackle.

## When to Use

- Start of a new session (proposed automatically via AGENTS.md hook)
- When the user asks "what should I work on?" or "triage"
- When returning after a break and context is stale

## bd Quick Reference

| Command | Action |
|---------|--------|
| `bd ready` | Show tasks with no open blockers (the "what next?" command) |
| `bd list` | Show all tasks |
| `bd create "Title" -p <0-4>` | Create task (0=critical, 4=low) |
| `bd update <id> --claim` | Claim and start a task |
| `bd close <id> --reason "Done"` | Close a completed task |
| `bd dep add <child> <parent>` | Add dependency (child is blocked by parent) |
| `bd show <id>` | View task details and audit trail |
| `bd search <query>` | Search tasks |

## What to Do

### Step 1: Check current tasks

Run `bd ready` to show actionable tasks (no open blockers).

- If tasks exist, present them to the user.
- Ask the user to choose:
  1. **Pick a task** — proceed to Step 4 (execution handoff).
  2. **Update the list** — go to Step 2 to add/modify/close tasks.
  3. **Full re-triage** — go to Step 3 to collect fresh data from all repos.

If no tasks exist (empty DB or all closed), go directly to Step 3.

### Step 2: Update tasks

Interactively update the task list based on user input:

- Create new tasks with `bd create`
- Close completed tasks with `bd close`
- Adjust priorities with `bd update <id> -p <0-4>`
- Add/remove dependencies with `bd dep add` / `bd dep remove`

After updates, run `bd ready` again and return to Step 1.

### Step 3: Full re-triage

Only run this when:
- The beads DB is empty (first run or fresh start)
- The user explicitly requests a full re-triage
- Context is very stale (e.g., returning after a long break)

**Collect data** (1 subagent per repo):

Read `setup.sh` in the workspace root to get the `REPOS` array. Launch one **read-only** subagent per repo (plus the workspace itself). Each subagent:

1. Reads `docs/project-plan.md` — find unchecked milestones and unmet requirements.
2. Lists `docs/exec-plan/todo/` — summarize each pending plan file.
3. Lists `docs/issues/` — summarize each logged issue.
4. Runs `gh pr list --state open` — find open PRs, note review status.
5. Runs `gh issue list --state open` — find open issues.
6. Returns a structured summary.

**Subagent rules**: Do NOT modify any files. Do NOT inspect other repos.

**Populate beads**: For each discovered item:
- `bd create "Title" -p <priority>` with appropriate priority (0=critical, 4=low)
- Use `bd dep add` to set dependencies where applicable
- Tag items by repo/category in the title or description

After populating, run `bd ready` and return to Step 1.

### Step 4: Execution handoff

For each confirmed task, propose one of:

1. **Do now (Quick Win)**: Execute immediately via subagent or directly.
   - For trivial fixes (e.g., small dependency bumps), just do it.
2. **Needs exec-plan**: Create `docs/exec-plan/todo/<name>.md` following `plan-execution` skill.
3. **Separate session**: Generate copy-paste prompts for a fresh session.

Claim the chosen task: `bd update <id> --claim`

#### Separate-session prompt template

Generate **two versions** (Japanese + English) per task:

```
## Task: <task name>
**Target repo**: <owner/repo> at <local path>
**Branch**: `git fetch origin && git switch -c <type>/<name> origin/main`
**Context to read first**:
- docs/project-plan.md
- docs/exec-plan/todo/<relevant plan if any>
- <other relevant files>

**Goal**: <what to accomplish>

**Deliverables**:
- <list of expected outputs>

**Constraints**:
- Follow Spec First rule
- Create PR via `gh pr create` when done
- Do not modify other repos
```

## Rules

1. This skill is **workspace-only**. Do not distribute to child repos.
2. Subagents are **read-only** during data collection. No file modifications.
3. Do not auto-execute tasks without user confirmation.
4. Keep the briefing fast: aim for < 2 minutes to reach a task selection.
5. Prefer `bd ready` over full re-triage. Full re-triage is expensive.
