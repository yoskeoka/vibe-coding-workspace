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
  envdiff/                      # Child: envdiff (Go)
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
    - For normal planning/execution work, use the globally installed `ww` CLI as the default branch/worktree operator path.
    - ALWAYS go through GitHub PR review for every change — including doc-only changes (Project Plan, Execution Plan).

2.  **Branch & PR Rules**:
    - Create a fresh task worktree from `main` for every task with global `ww`: `ww create <type>/<description>` from the target repo, or `ww create --repo <repo> <type>/<description>` from the workspace root
    - Enter task worktrees with `ww cd` rather than guessing paths manually
    - Never reuse an existing feature branch or primary checkout silently; each active task should have its own `ww` worktree
    - Even for lightweight/no-plan changes, re-check `AI_WORKFLOW.md` before PR creation so the branch type and PR title match the actual scope (`docs`, `chore`, etc.)
    - **Before pushing to a PR branch**, always verify the PR is still OPEN: `gh pr view <number> --json state --jq '.state'`. Never push to a MERGED or CLOSED PR.
    - Workflow-linter findings must not be ignored. Resolve all `fixable` warnings before push/PR unless an explicit human instruction conflicts or the warning is a clear false positive.
    - If a `fixable` workflow-linter warning is skipped, record the reason in the PR body.
    - Run all lint and test checks (non-AI tooling) before creating a PR. Fix failures before proceeding.
    - Create PRs via `gh pr create` and wait for review approval before merging.

3.  **Context Management**:
    - Your "memory" is the `docs/` directory.
    - `docs/project-plan.md` is your North Star.
    - `docs/exec-plan/todo/` is your current task list.
    - `docs/design-decisions/` is your architectural conscience.
    - **Before making a design decision**, read `docs/design-decisions/core-beliefs.md` and relevant entries in `docs/design-decisions/adr.md`. Present what you found (e.g., "Past decision: X was chosen because Y. Apply the same reasoning here?") before proposing a new direction.

4.  **Execution Rules**:
    - **Plan First**: Before writing code, ensure a plan exists in `docs/exec-plan/todo/`. If not, create one.
    - **Spec First**: Update `docs/specs/` to reflect changes BEFORE modifying code.
    - **Focus**: if you find unrelated issues, log them in `docs/issues/<name>.md` and ignore them for the current task (unless they are blockers).
    - If `ww` fails or behaves unexpectedly during normal workflow use, capture it as a first-class workflow output per `docs/specs/ww-dogfooding-workflow.md` instead of silently dropping to raw git. Record the command, cwd, target repo, expected vs actual behavior, relevant output, fallback usage, and impact.
    - **Issue Resolution**: When an issue in `docs/issues/` is resolved, move the file to `docs/issues/done/`.
    - **Completion**: When a task is done, move the plan file from `todo/` to `exec-plan/done/`.
    - **Post-Task Review**: After completing significant work, run `post-task-review` to log issues, update lessons learned, and propose CLAUDE.md/AGENTS.md updates before creating a PR.

## When asked to "Start a new feature":
1.  Create a planning worktree with global `ww`: `ww create plan/feature-name`, then enter it with `cd "$(ww cd plan/feature-name)"`
2.  Read `docs/project-plan.md`.
3.  Create a new file in `docs/exec-plan/todo/feature-name.md`.
4.  Outline the changes to specs and code in that plan.
5.  Create a PR for the plan and wait for review.
6.  After plan PR is merged, create an execution worktree with global `ww`: `ww create feat/feature-name`, then enter it with `cd "$(ww cd feat/feature-name)"`
7.  Execute the plan following **Spec First** rule.
8.  Run lint/tests, fix any failures, then create a PR.

## When asked to "Start a new project":
1.  Start with the `new-project-intake` skill (idea sparring → research → go/no-go).
2.  If no-go, log the result in `docs/design-decisions/rejected-ideas.md` and stop.
3.  If go, complete bootstrap via `new-project-intake`, then move to the child repo and continue with `plan-project`.

## When asked to "Fix a bug":
1.  Create a planning worktree with global `ww`: `ww create plan/fix-bug-x`, then enter it with `cd "$(ww cd plan/fix-bug-x)"`
2.  Create a plan in `docs/exec-plan/todo/fix-bug-x.md`.
3.  Reproduction steps go into the plan.
4.  Create a PR for the plan and wait for review.
5.  After plan PR is merged, create an execution worktree with global `ww`: `ww create fix/fix-bug-x`, then enter it with `cd "$(ww cd fix/fix-bug-x)"`
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

## Workspace Task Tracking

Workspace-level triage uses a GitHub Project plus a local derived cache managed by the `pj` CLI spike.

### Source of Truth

- GitHub Projects is the canonical remote state for workspace task triage
- `.local/pj/cache.json` is derived data for fast local reads and AI access
- `docs/exec-plan/todo/` remains the canonical implementation-plan tracker once work is selected

### Expected Commands

- `go -C tools/pj run ./cmd/pj init --owner <owner> --owner-type user|org`
- `go -C tools/pj run ./cmd/pj sync [--project <number>]`
- `go -C tools/pj run ./cmd/pj list`
- `go -C tools/pj run ./cmd/pj add --title "..." --status Todo`
- `go -C tools/pj run ./cmd/pj move --item <item-id> --status "In Progress"`
- `go -C tools/pj run ./cmd/pj config show|set|clear`

### Important Rules

- Use the GitHub Project flow for workspace-level task triage instead of `bd`
- Do not treat `.local/pj/` as source-of-truth data; it can be regenerated
- Do not create duplicate local task trackers for the same workspace board unless a spec explicitly adds one
- Bootstrap the canonical board with `pj init` before relying on `sync`, `list`, `add`, or `move`
- After bootstrap, reuse the stored owner scope instead of repeatedly restating `--owner` and `--owner-type`
- If you need to switch from a user board to an org board, do it explicitly with `pj config set` or `pj config clear` before re-running `pj init`; do not mix one-off owner flags with an existing local config

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
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
