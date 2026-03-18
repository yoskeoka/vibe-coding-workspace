# AI Agent Behavior Guidelines

> **Note**: AGENTS.md is the canonical file. CLAUDE.md is a symlink to this file. Do not edit them separately — any change here applies to both.

You are an expert software engineer and architect working in an **AI-Centered Development** environment. Your primary goal is to write high-quality, maintainable code while strictly adhering to the workflow defined in [AI_WORKFLOW.md](AI_WORKFLOW.md).

## Workspace Structure

This is a **meta-repo** (parent workspace). Child projects are **subdirectories**, not separate top-level repositories:

```
vibe-coding-workspace/          # This repo (workspace root)
  ww/                           # Child: Workspace Worktree CLI (Go)
  ai-arena/                     # Child: AI Arena
  reversi-adventure/            # Child: Reversi Adventure
  vim-learning-game/            # Child: Vim Learning Game
```

**CRITICAL**: When given a GitHub URL like `yoskeoka/ww`, the repo lives at `<workspace-root>/ww/`, NOT at `~/src/github.com/yoskeoka/ww/`. Always resolve child project paths relative to this workspace.

## Skill Priority

When multiple skills could apply, **project-level skills take precedence over global skills**. Project skills (e.g., `plan-execution`, `execute-task`) define the authoritative workflow for this workspace. Global skills (e.g., `superpowers:writing-plans`) are fallbacks for cases where no project skill exists.

## Session Start

At the start of a new session, if the user has not given a specific task, suggest running `triage-tasks` to review priorities across all managed projects.

## Core Responsibilities

1.  **Workflow Adherence**:
    - ALWAYS read and follow `AI_WORKFLOW.md`.
    - NEVER skip the "Execution Plan" phase for non-trivial changes. Trivial changes (single-line fixes, typos, doc-only updates) may skip the plan and go directly to execution.
    - NEVER write code without a corresponding specification update in `docs/specs/`.
    - ALWAYS create a new branch from the latest `main` before starting any work.
    - ALWAYS go through GitHub PR review for every change — including doc-only changes (Project Plan, Execution Plan).

2.  **Branch & PR Rules**:
    - Create a fresh branch from `origin/main` for every task: `git fetch origin && git switch -c <branch-name> origin/main`
    - Never reuse an existing feature branch; always create a fresh one.
    - **Before pushing to a PR branch**, always verify the PR is still OPEN: `gh pr view <number> --json state --jq '.state'`. Never push to a MERGED or CLOSED PR.
    - Run all lint and test checks (non-AI tooling) before creating a PR. Fix failures before proceeding.
    - Create PRs via `gh pr create` and wait for review approval before merging.

3.  **Context Management**:
    - Your "memory" is the `docs/` directory.
    - `docs/project-plan.md` is your North Star.
    - `docs/exec-plan/todo/` is your current task list.
    - `docs/design-decisions/` is your architectural conscience.

4.  **Execution Rules**:
    - **Plan First**: Before writing code, ensure a plan exists in `docs/exec-plan/todo/`. If not, creates one.
    - **Spec First**: Update `docs/specs/` to reflect changes BEFORE modifying code.
    - **Focus**: if you find unrelated issues, log them in `docs/issues/<name>.md` and ignore them for the current task (unless they are blockers).
    - **Issue Resolution**: When an issue in `docs/issues/` is resolved, move the file to `docs/issues/done/`.
    - **Completion**: When a task is done, move the plan file from `todo/` to `exec-plan/done/`.
    - **Post-Task Review**: After completing significant work, run `post-task-review` to log issues, update lessons learned, and propose CLAUDE.md/AGENTS.md updates before creating a PR.

## When asked to "Start a new feature":
1.  Create a branch: `git fetch origin && git switch -c plan/feature-name origin/main`
2.  Read `docs/project-plan.md`.
3.  Create a new file in `docs/exec-plan/todo/feature-name.md`.
4.  Outline the changes to specs and code in that plan.
5.  Create a PR for the plan and wait for review.
6.  After plan PR is merged, create a new branch for execution: `git fetch origin && git switch -c feat/feature-name origin/main`
7.  Execute the plan following **Spec First** rule.
8.  Run lint/tests, fix any failures, then create a PR.

## When asked to "Start a new project":
1.  Start with the `new-project-intake` skill (idea sparring → research → go/no-go).
2.  If no-go, log the result in `docs/design-decisions/rejected-ideas.md` and stop.
3.  If go, complete bootstrap via `new-project-intake`, then move to the child repo and continue with `plan-project`.

## When asked to "Fix a bug":
1.  Create a branch: `git fetch origin && git switch -c plan/fix-bug-x origin/main`
2.  Create a plan in `docs/exec-plan/todo/fix-bug-x.md`.
3.  Reproduction steps go into the plan.
4.  Create a PR for the plan and wait for review.
5.  After plan PR is merged, create a new branch: `git fetch origin && git switch -c fix/fix-bug-x origin/main`
6.  Execute the fix following the **Spec First** rule.
7.  Run lint/tests, fix any failures, then create a PR.
8.  Move plan to `done/`.

## Subagent Strategy

Keep the main context window clean by delegating to subagents.

### Delegate to subagents:
- Codebase exploration and search (grep, file structure investigation)
- Documentation research
- Parallel analysis of multiple files
- Independent verification tasks (test execution, lint checks)
- Any research that might add >1000 tokens to main context

### Keep in main context:
- Final implementation decisions
- User communication
- State that needs to persist across steps
- Sequential dependent operations (spec update → code implementation ordering)

### Rules:
- One task per subagent for focused execution
- Clear, specific instructions with expected output format
- Set scope boundaries — subagents must not modify files without explicit instruction

<!-- BEGIN BEADS INTEGRATION v:1 profile:full hash:d4f96305 -->
## Issue Tracking with bd (beads)

**IMPORTANT**: This project uses **bd (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why bd?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Dolt-powered version control with native sync
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**

```bash
bd ready --json
```

**Create new issues:**

```bash
bd create "Issue title" --description="Detailed context" -t bug|feature|task -p 0-4 --json
bd create "Issue title" --description="What this issue is about" -p 1 --deps discovered-from:bd-123 --json
```

**Claim and update:**

```bash
bd update <id> --claim --json
bd update bd-42 --priority 1 --json
```

**Complete work:**

```bash
bd close bd-42 --reason "Completed" --json
```

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `bd ready` shows unblocked issues
2. **Claim your task atomically**: `bd update <id> --claim`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" --description="Details about what was found" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`

### Auto-Sync

bd automatically syncs via Dolt:

- Each write auto-commits to Dolt history
- Use `bd dolt push`/`bd dolt pull` for remote sync
- No manual export/import needed!

### Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use `--json` flag for programmatic use
- ✅ Link discovered work with `discovered-from` dependencies
- ✅ Check `bd ready` before asking "what should I work on?"
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems

For more details, see README.md and docs/QUICKSTART.md.

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd dolt push
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

<!-- END BEADS INTEGRATION -->
