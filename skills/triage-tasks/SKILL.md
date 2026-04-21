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
| `go -C tools/pj run ./cmd/pj add --title "..." --body-file <path> --status Todo --repo <repo> --kind <kind> --priority <priority>` | Create a new triage item |
| `go -C tools/pj run ./cmd/pj update --item <item-id> --status "In Progress"` | Claim and start a task |
| `go -C tools/pj run ./cmd/pj update --item <item-id> --status Done` | Close a completed task |
| `go -C tools/pj run ./cmd/pj url` | Print the canonical GitHub Project URL |

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
- Exclude routine workflow-skill update PR items from triage. If an item points to a PR titled `chore: update workflow skills` (for example, `https://github.com/yoskeoka/vim-learning-game/pull/93`) or has that title, do not rank or present it as actionable work; close the Project item as `Done` when updating the list.
- Use the `Priority` column plus current session context to decide what is most actionable.
- The canonical Project must have a `Priority` field; if an item's `Priority` value is empty, unset, or shows as `-`, rank a short list anyway using explicit heuristics:
  - active exec plans over vague future ideas
  - broken workflow or failing review items over speculative enhancements
  - the repo the user is currently focused on over unrelated backlog
- This spike does **not** have dependency-aware "ready queue" resolution; do not infer blockers unless the task title/body or repo context makes them explicit.

- If tasks exist, present them to the user.
- Include the GitHub Project URL from `go -C tools/pj run ./cmd/pj url` when the cache is available.
- Present a short prioritized top list before the full board dump when the board is large.
- Present the next step as explicit numbered choices and wait for the user's answer:
  1. **Pick a task** — proceed to Step 4 (execution handoff).
  2. **Update the list** — go to Step 2 to add/modify/close tasks.
  3. **Full re-triage** — go to Step 3 to collect fresh data from all repos.

If no tasks exist (empty board or all tasks are already `Done`), go directly to Step 3.

### Step 2: Update tasks

Interactively update the task list based on user input:

- Create new tasks with `go -C tools/pj run ./cmd/pj add`
- Every new task created with `pj add` must include `--body` or `--body-file` using the compact handoff format from Step 3; prefer `--body-file` for generated multi-line bodies.
- Move a selected task to `In Progress` with `go -C tools/pj run ./cmd/pj update --item <id> --status "In Progress"`
- Close completed tasks with `go -C tools/pj run ./cmd/pj update --item <id> --status Done`
- Correct `Repo`, `Kind`, or `Priority` on an existing item with `pj update` followed by `pj sync` or `pj list`

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

When collecting open PRs, ignore PRs titled `chore: update workflow skills`. These are routine workflow propagation PRs that can be checked opportunistically from the repo PR list; they churn as workspace updates close older PRs and open replacements, so they should not create durable triage tasks.

**Subagent rules**: Do NOT modify any files. Do NOT inspect other repos. The main agent remains responsible for final prioritization, `pj` mutations, and the fresh-session handoff prompt.

**Populate the Project**: For each discovered item:
- Create a Project item with `go -C tools/pj run ./cmd/pj add --title "..." --body-file <path> --status Todo --repo <repo> --kind <kind> --priority <priority>`
- Use the title and body to encode enough context for later triage; the current spike does not model dependency edges
- Keep repo/category information normalized through the `Repo` and `Kind` fields whenever possible

#### Project item body format

Stored Project item bodies are durable workspace-board data. Write them in English for consistency across sessions, even when the current chat is in another language. The separate user-facing handoff prompt in Step 4 still follows the user's current chat language.

Keep bodies compact: the body should help a human or agent start from the GitHub Project item without becoming a full execution plan. Use this structure:

```markdown
Source: <local path, PR URL, issue URL, or discovery source>
Repo: <workspace repo>
Next: <plan-execution|execute-task|Manual triage>
Start: <ww create ... then cd "$(ww cd ...)" | Not yet specified>
Read: <short comma-separated list of first files/URLs to inspect>
Goal: <one sentence outcome>
```

Use one `Read:` line unless the task truly needs multiple source references. Prefer exact local paths for workspace-local sources and URLs for GitHub-native sources.

Tailor the body by item type:

- Exec-plan item: use `Next: execute-task`; set `Source` and `Read` to `docs/exec-plan/todo/<name>.md`; set `Start` to the matching execution branch command based on the plan type: `ww create feat/<name>` then `cd "$(ww cd feat/<name>)"` or `ww create fix/<name>` then `cd "$(ww cd fix/<name>)"` from the target repo root, or the corresponding `--repo <repo>` form from the workspace root for child repos.
- Local issue follow-up: use `Next: plan-execution` when the issue needs non-trivial work; set `Source` and `Read` to `docs/issues/<name>.md`; set `Start` to `ww create plan/<name>` then `cd "$(ww cd plan/<name>)"` from the target repo root, or the `--repo <repo>` form from the workspace root for child repos.
- Open PR review/follow-up: use `Next: Manual triage` unless the needed skill is clear; set `Source` and `Read` to the PR URL plus any relevant local plan/spec path; set `Start: Not yet specified` for pure review, approval, or post-review response items.
- Open GitHub Issue: use `Next: Manual triage` unless the issue clearly maps to planning or execution; set `Source` and `Read` to the issue URL; set `Start: Not yet specified` unless there is a concrete local workflow branch to create.

Do not invent local plan paths, issue paths, blockers, or dependency relationships that are not present in the collected data. If the next workflow step is unclear, prefer `Manual triage` over a misleading skill recommendation.

After populating, run `go -C tools/pj run ./cmd/pj list` and return to Step 1.

### Step 4: Execution handoff

For each confirmed task, propose one of:

1. **Fresh session prompt (default)**: Generate a copy-paste prompt for a new session so implementation starts with a clean context window.
2. **Needs exec-plan**: If the task is non-trivial and has no execution plan yet, make the prompt target plan creation first.
3. **Do now (exception)**: Only stay in the same session when the user explicitly wants immediate execution despite the broader triage context.

Claim the chosen task: `go -C tools/pj run ./cmd/pj update --item <id> --status "In Progress"`

#### Separate-session prompt template

Generate the prompt in the same language the user is currently using in the chat:

```
Use skill: <plan-execution|execute-task|other explicit next-step skill>

## Task: <task name>
**Target repo**: <owner/repo> at <local path>
**Worktree setup**:
- If starting in the target repo root: `ww create <type>/<name>` then `cd "$(ww cd <type>/<name>)"`
- If starting in the workspace root for a child repo: `ww create --repo <repo> <type>/<name>` then `cd "$(ww cd --repo <repo> <type>/<name>)"`
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

- If the task is non-trivial and does not already have an execution plan, use `plan-execution`.
- If the task already has an execution plan and the next session should implement it, use `execute-task`.

## Rules

1. This skill is **workspace-only**. Do not distribute to child repos.
2. Subagents are **read-only** during data collection. No file modifications.
3. Do not auto-execute tasks without user confirmation.
4. Keep the briefing fast: aim for < 2 minutes to reach a task selection.
5. Prefer `pj list` plus the cached `Priority` field over full re-triage. Full re-triage is expensive.
6. Do not mention or rely on legacy tracker backup, dependency commands, or other retired workflow concepts in this workflow.
7. After a task is selected, default to emitting a fresh-session handoff prompt instead of starting the work in the same session.
