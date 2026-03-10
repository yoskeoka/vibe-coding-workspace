---
name: triage-tasks
description: When starting a new session, reviewing priorities, running a daily briefing, triaging pending tasks across workspace and child projects, or deciding what to work on next.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Triage Tasks — Daily Briefing (Workspace-Only)

**Position in workflow**: This is a **session-start ritual**, not a numbered workflow step. Run it at the beginning of a work session to decide what to tackle.

## When to Use

- Start of a new session (proposed automatically via AGENTS.md hook)
- When the user asks "what should I work on?" or "triage"
- When returning after a break and context is stale

## What to Do

### Step 0: Check existing priorities

Read `.local/priority.md`. If it exists and has non-expired entries (check TTL against current date):

1. Present the existing priority list to the user.
2. Ask if they want to **use it as-is**, **update it**, or **re-triage from scratch**.
3. If using as-is or updating, skip to **Step 5** (wall-hit and finalize).
4. If re-triaging, continue to Step 1.

If the file does not exist or all entries have expired, continue to Step 1.

### Step 1: Identify managed repos

Read `setup.sh` in the workspace root to get the `REPOS` array.
The workspace itself is also a target. Build the full list:

- `vibe-coding-workspace` (this workspace)
- Each child project from `REPOS`

### Step 2: Collect data (1 subagent per repo)

Launch one **read-only** subagent per repo. Each subagent:

1. Reads `docs/project-plan.md` — find unchecked milestones and unmet requirements.
2. Lists `docs/exec-plan/todo/` — summarize each pending plan file.
3. Lists `docs/issues/` — summarize each logged issue.
4. Runs `gh pr list --state open` — find open PRs, note review status.
5. Runs `gh issue list --state open` — find open issues (including Dependabot / automated alerts).
6. Returns a table in this exact format:

```
| Item | Source | Type | Effort | Risk | Suggested action |
|------|--------|------|--------|------|------------------|
| ...  | ...    | ...  | S/M/L  | L/M/H| ...              |
```

**Subagent rules**:

- Do NOT modify any files.
- Do NOT inspect other repos.
- If a source file/directory does not exist, skip it silently.

### Step 3: Classify

Merge all subagent results and tag each item:

| Category      | Criteria                                                       |
| ------------- | -------------------------------------------------------------- |
| **Must**      | Blocking work, stale PRs (>3 days), security alerts, broken CI |
| **Should**    | High-impact project-plan gaps, significant exec-plan items     |
| **Quick Win** | <30 min effort, low risk, independent of other tasks           |

### Step 4: Prioritize

Score each item:

$$P = \frac{Impact \times Urgency \times Confidence}{Effort}$$

Scale: 1–5 for each factor. Sort descending by P.

Present the **Top 5** to the user in a compact table:

```
| # | Task | Repo | Category | P | Effort | Action |
|---|------|------|----------|---|--------|--------|
```

### Step 5: Wall-hit and finalize

- Allow 2–3 rounds of quick discussion to adjust the Top 5.
- User confirms the final list.

### Step 6: Save to `.local/priority.md`

Write the confirmed Top 5 to `.local/priority.md` with TTL (default: 72h).
Drop any entries whose TTL has expired from a previous run.

Format:

```markdown
# Current Priorities

Updated: YYYY-MM-DD HH:MM

| #   | Task | Repo | Category | P   | TTL | Status      |
| --- | ---- | ---- | -------- | --- | --- | ----------- |
| 1   | ...  | ...  | Must     | 25  | 72h | not-started |
```

### Step 7: Execution handoff

For each confirmed task, propose one of:

1. **Do now (Quick Win)**: Execute immediately via subagent or directly.
   - For trivial GitHub issue fixes (e.g., small dependency bumps), just do it.
2. **Needs exec-plan**: Create `docs/exec-plan/todo/<name>.md` following `plan-execution` skill.
3. **Separate session**: Generate copy-paste prompts for a fresh session.

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
2. Subagents are **read-only**. No file modifications during collection.
3. The `.local/priority.md` file is **not Git-tracked**. It is a volatile convenience memo.
4. Do not auto-execute tasks without user confirmation.
5. Keep the briefing fast: aim for < 5 minutes to reach a confirmed Top 5.
