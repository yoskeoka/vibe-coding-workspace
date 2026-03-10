# Spec: Triage Tasks — Daily Briefing Workflow

## Goal
Provide an automated daily briefing that scans all managed repositories, classifies pending work, and proposes a prioritized Top 5 task list — with minimal user effort.

## Scope
- Automated collection of pending work across workspace + all child projects
- Classification and prioritization
- Volatile storage of the result (non-Git)
- Execution handoff (in-session subagent or separate-session prompt generation)

## Requirements

### 1. Session start hook (AGENTS.md)
A short entry in `AGENTS.md` MUST propose triage at the start of every new session:
- "At the start of a new session, if the user has not given you a specific task yet, suggest running `triage-tasks` to review priorities."
- Keep to 1–2 lines. Details live in the skill.

### 2. Skill: `triage-tasks`
A workspace-only skill (not distributed to child repos via `setup-workspace.sh`).

#### 2a. Data collection (subagent per repo)
For each managed repo (workspace itself + child projects listed in `setup.sh` REPOS):
- Launch one subagent per repo (read-only, no file edits).
- Each subagent collects:
  - `project-plan.md` — unchecked milestones / unmet requirements
  - `docs/exec-plan/todo/` — pending execution plans
  - `docs/issues/` — logged issues
  - GitHub PRs — open PRs, especially those awaiting review
  - GitHub Issues — including dependency/automated alerts (e.g., Dependabot)
- Each subagent returns a fixed-format table:
  `Item | Source | Type | Effort (S/M/L) | Risk (L/M/H) | Suggested action`

#### 2b. Classification
Aggregate subagent results and classify each item:
- **Must**: blocking, stale PRs, security alerts, broken CI
- **Should**: high-impact project-plan gaps, significant exec-plan items
- **Quick Win**: <30 min effort, low risk, independent

#### 2c. Prioritization
Score each item: P = (Impact × Urgency × Confidence) / Effort (1–5 scale).
Propose Top 5 ordered by P descending.

#### 2d. Decision support
Present Top 5 to user for quick 2–3 round wall-hitting to finalize.

#### 2e. Execution handoff
For each confirmed task, propose:
- **Do now (in-session)**: subagent or direct execution for Quick Wins
- **Needs exec-plan**: create `docs/exec-plan/todo/` entry
- **Separate session**: generate copy-paste prompts (Japanese + English) scoped to one repo

### 3. Volatile storage: `.local/priority.md`
- Stored in `.local/` directory (non-Git, added to `.gitignore`).
- Contains only the current Top 5 with TTL (default 72h).
- Format: table with columns `# | Task | Repo | Category | P-score | TTL | Status`
- Stale entries (past TTL) are dropped on next triage run.
- This file is a convenience memo, not a source of truth.

### 4. Separate-session prompt templates
The skill MUST be able to generate per-repo execution prompts in both Japanese and English.
Template structure:
- Target repo path/name
- Task description
- Relevant context files to read
- Expected deliverables
- Constraints (branch rules, spec-first, etc.)

## Non-Goals
- Replacing `exec-plan/todo` as the canonical task tracker
- Automated execution without user confirmation
- Distribution to child repos (workspace-only skill)
- Git-tracking the priority list
