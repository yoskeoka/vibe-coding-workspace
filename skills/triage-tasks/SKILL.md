---
name: triage-tasks
description: When starting a new session, reviewing priorities, running a daily briefing, triaging pending tasks across workspace and child projects, or deciding what to work on next.
metadata:
  author: yoskeoka
  version: '2.1.0'
---

# Triage Tasks — Daily Briefing (Workspace-Only)

**Position in workflow**: This is a **session-start ritual**, not a numbered workflow step. Run it at the beginning of a work session to decide what to tackle.

## When to Use

- Start of a new session (proposed automatically via AGENTS.md hook)
- When the user asks "what should I work on?" or "triage"
- When returning after a break and context is stale

## `pj` Quick Reference

| Command | Action |
|---------|--------|
| `go -C tools/pj run ./cmd/pj init --owner <owner> --owner-type user|org` | Resolve or create the canonical `Workspace Task Triage` board and write cache metadata |
| `go -C tools/pj run ./cmd/pj sync` | Refresh `.local/pj/cache.json` from the configured GitHub Project |
| `go -C tools/pj run ./cmd/pj list` | Show cached tasks with `Status`, `Repo`, `Kind`, and `Priority` |
| `go -C tools/pj run ./cmd/pj add --title "..." --status Todo --repo <repo> --kind <kind> --priority <priority>` | Create a new triage item |
| `go -C tools/pj run ./cmd/pj move --item <item-id> --status "In Progress"` | Claim and start a task |
| `go -C tools/pj run ./cmd/pj move --item <item-id> --status Done` | Close a completed task |

## What to Do

### Step 0: Bootstrap or refresh the workspace Project cache

Before checking tasks, make sure the canonical GitHub Project is configured locally:

```bash
git pull
```

- If `.local/pj/cache.json` does not exist yet, run `go -C tools/pj run ./cmd/pj init --owner <owner> --owner-type user|org`.
- If the cache already exists, run `go -C tools/pj run ./cmd/pj sync`.
- If `sync` fails because the cache is incomplete or stale, rerun `init` with explicit owner metadata, then `sync`.

### Step 1: Check current tasks

Run `go -C tools/pj run ./cmd/pj list` to inspect the current queue.

- Treat `Status=Todo` items as the default candidate set.
- Use the `Priority` column plus current session context to decide what is most actionable.
- If `Priority` is missing, rank a short list anyway using explicit heuristics:
  - active exec plans over vague future ideas
  - broken workflow or failing review items over speculative enhancements
  - the repo the user is currently focused on over unrelated backlog
- This spike does **not** have `bd ready`-style dependency resolution; do not infer blockers unless the task title/body or repo context makes them explicit.

- If tasks exist, present them to the user.
- Include the GitHub Project URL when owner scope and project number are known.
- Present a short prioritized top list before the full board dump when the board is large.
- Present the next step as explicit numbered choices and wait for the user's answer:
  1. **Pick a task** — proceed to Step 4 (execution handoff).
  2. **Update the list** — go to Step 2 to add/modify/close tasks.
  3. **Full re-triage** — go to Step 3 to collect fresh data from all repos.

If no tasks exist (empty board or all tasks are already `Done`), go directly to Step 3.

### Step 2: Update tasks

Interactively update the task list based on user input:

- Create new tasks with `go -C tools/pj run ./cmd/pj add`
- Move a selected task to `In Progress` with `go -C tools/pj run ./cmd/pj move --item <id> --status "In Progress"`
- Close completed tasks with `go -C tools/pj run ./cmd/pj move --item <id> --status Done`
- If `Repo`, `Kind`, or `Priority` on an existing item are wrong, treat that as a manual GitHub Project edit followed by `pj sync`, because the current CLI only updates `Status`

After updates, run `go -C tools/pj run ./cmd/pj list` again and return to Step 1.

### Step 3: Full re-triage

Only run this when:
- The Project cache is empty or the board has no relevant items yet
- The user explicitly requests a full re-triage
- Context is very stale (e.g., returning after a long break)

**Collect data** (1 subagent per repo):

Read `setup.sh` in the workspace root to get the `REPOS` array. Launch one **read-only** subagent per repo (plus the workspace itself) when that helps keep the main context smaller. Prefer an available low-cost small model for these exploration subagents; do not hard-code a specific model version in the skill. Each subagent:

1. Reads `docs/project-plan.md` — find unchecked milestones and unmet requirements.
2. Lists `docs/exec-plan/todo/` — summarize each pending plan file.
3. Lists `docs/issues/` — summarize each logged issue.
4. Runs `gh pr list --state open` — find open PRs, note review status.
5. Runs `gh issue list --state open` — find open issues.
6. Returns a structured summary.

**Subagent rules**: Do NOT modify any files. Do NOT inspect other repos. The main agent remains responsible for final prioritization, `pj` mutations, and the fresh-session handoff prompt.

**Populate the Project**: For each discovered item:
- Create a Project item with `go -C tools/pj run ./cmd/pj add --title "..." --status Todo --repo <repo> --kind <kind> --priority <priority>`
- Use the title and body to encode enough context for later triage; the current spike does not model dependency edges
- Keep repo/category information normalized through the `Repo` and `Kind` fields whenever possible

After populating, run `go -C tools/pj run ./cmd/pj list` and return to Step 1.

### Step 4: Execution handoff

For each confirmed task, propose one of:

1. **Fresh session prompt (default)**: Generate a copy-paste prompt for a new session so implementation starts with a clean context window.
2. **Needs exec-plan**: If the task is non-trivial and has no execution plan yet, make the prompt target plan creation first.
3. **Do now (exception)**: Only stay in the same session when the user explicitly wants immediate execution despite the broader triage context.

Claim the chosen task: `go -C tools/pj run ./cmd/pj move --item <id> --status "In Progress"`

#### Separate-session prompt template

Generate the prompt in the same language the user is currently using in the chat:

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
- Create PR via `gh pr create` using the **PR template** when done
- Do not modify other repos
```

## Rules

1. This skill is **workspace-only**. Do not distribute to child repos.
2. Subagents are **read-only** during data collection. No file modifications.
3. Do not auto-execute tasks without user confirmation.
4. Keep the briefing fast: aim for < 2 minutes to reach a task selection.
5. Prefer `pj list` plus the cached `Priority` field over full re-triage. Full re-triage is expensive.
6. Do not mention or rely on beads/Dolt backup, dependency commands, or other `bd`-only concepts in this workflow.
7. After a task is selected, default to emitting a fresh-session handoff prompt instead of starting the work in the same session.
