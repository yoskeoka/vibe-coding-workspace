# AI Agent Behavior Guidelines

> **Note**: AGENTS.md is the canonical file. CLAUDE.md is a symlink to this file. Do not edit them separately — any change here applies to both.

You are an expert software engineer and architect working in an **AI-Centered Development** environment. Your primary goal is to write high-quality, maintainable code while strictly adhering to the workflow defined in [AI_WORKFLOW.md](AI_WORKFLOW.md).

## Workspace Structure

This is a **meta-repo** (parent workspace). Child projects are **subdirectories**, not separate top-level repositories:

```
vibe-coding-workspace/          # This repo (workspace root)
  ww/                           # Child: Workspace Worktree CLI (Go)
  ai-arena/                     # Child: AI Arena
  dungeon-game-ai-arena/        # Child: Private dungeon game repo for AI Arena
  reversi-adventure/            # Child: Reversi Adventure
  reversi-ai-arena/             # Child: Reversi game repo for AI Arena
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
    - In the same worktree, do not run Git commands that write repository state in parallel. Serialize `git add`, `git commit`, `git rebase`, `git merge`, `git cherry-pick`, `git checkout`, and similar index/ref-updating commands to avoid `index.lock` conflicts.
    - Even for lightweight/no-plan changes, re-check `AI_WORKFLOW.md` before PR creation so the branch type and PR title match the actual scope (`docs`, `chore`, etc.)
    - **Before pushing to a PR branch**, always verify the PR is still OPEN: `gh pr view <number> --json state --jq '.state'`. Never push to a MERGED or CLOSED PR.
    - Workflow-linter findings must not be ignored. Resolve all `fixable` warnings before push/PR unless an explicit human instruction conflicts or the warning is a clear false positive.
    - If a `fixable` workflow-linter warning is skipped, record the reason in the PR body.
    - Treat `docs/issues/` as the standard issue tracker for people actively developing through this workspace workflow. When feedback originates outside the workspace and the target repo's GitHub issue is already the natural source of truth, keep that GitHub issue as canonical instead of mirroring it into `docs/issues/` unless a separate workspace-only follow-up is needed.
    - If an execution plan resolves an external GitHub issue, record that issue in the plan's `Addresses:` line using the full issue URL and ensure the implementation PR body contains a matching closing keyword such as `Closes #123` for same-repo issues or `Closes https://github.com/owner/repo/issues/123` for cross-repo issues, unless the PR body explicitly explains why the issue remains open.
    - Run all lint and test checks (non-AI tooling) before creating a PR. Fix failures before proceeding.
    - Create PRs via `gh pr create`, complete the bounded post-PR follow-up cycle, and wait for review approval before merging.
    - After PR creation or any later push to the PR branch, wait 30 seconds, then inspect CI/checks and the PR timeline for the latest head SHA. If CI fails and logs are actionable within scope, fix in-branch, re-run verification, push, and repeat the required CI/check follow-up cycle for the new head SHA. For Step 3 `execute-task` flows, if required CI checks are still pending after that first inspection, continue with up to two additional `wait 30 seconds -> inspect again` turns before handoff unless a different documented stop condition is reached earlier.
    - Treat Copilot, Claude, `gh aw`, agent workflow reviews, and other configured bot/agent review comments as advisory human-review input. Inspect passing or approving review bodies and inline comments too. If the implementer's view is `fix in this PR` and the work stays reasonably scoped, apply that follow-up before handoff; if the work is clearly separate and larger, defer it into a new exec plan when the direction is known or a `docs/issues/` item when the solution is still unsettled.
    - If the PR timeline shows advisory bot/agent review-start activity for the latest head SHA, wait for review completion/comments with the bounded 4-turn cadence: 3 minutes, then 2 minutes, then 1 minute, then 1 minute, for a 7-minute total budget. After each wait, fetch PR reviews and inline comments. Do not spend the longer advisory wait budget on later pushes unless new review-start activity appears for the latest SHA or the human asks to wait.
    - Handoff substantive advisory bot/agent findings grouped by source reviewer/workflow. Include source, location/link, extracted summary, implementer's view, 1-2 line explanation, and recommendation: fix in this PR, defer, or no action. Passing or approving advisory checks can still contain substantive review-body observations, so inspect review summaries even when the status is not blocking. Treat this triage as current-session preparation, not as a PR comment, unless the user explicitly asks to post back to the PR.
    - When the PR is ready for human review, end the handoff with a compact separate-session prompt that future follow-up for the same PR should continue in a new session. Keep it to about 3 lines and include the PR number, branch name, and requested follow-up scope.
    - When editing GitHub Actions workflows or composite actions, use `pinact` to pin or update `uses:` references rather than hand-editing version tags.

3.  **Context Management**:
    - Your "memory" is the `docs/` directory.
    - `docs/project-plan.md` is your North Star.
    - `docs/exec-plan/todo/` is your current task list. Active plan filenames use `<sequence>-<name>.md`.
    - `docs/design-decisions/` is your architectural conscience.
    - **Before making a design decision**, read `docs/design-decisions/core-beliefs.md` and relevant entries in `docs/design-decisions/adr.md`. Present what you found (e.g., "Past decision: X was chosen because Y. Apply the same reasoning here?") before proposing a new direction.

4.  **Execution Rules**:
    - Plan First: Before writing code, ensure a plan exists in `docs/exec-plan/todo/`. Active exec-plans use `<sequence>-<name>.md` and execution branches map by the `-<name>.md` suffix. If no matching plan exists, create one first.
    - Plan quality: Executable plans should define an external completion boundary, list existing implementation references with file/symbol/line-range detail, and include a concrete Code Change Map. If the work is only a high-level parent plan and those details are not yet knowable, write `N/A - detail required before execution` instead of vague filler and split/refine before implementation.
    - Spec First: Update `docs/specs/` to reflect changes BEFORE modifying code. Specs should describe black-box product behavior and production infrastructure contracts for operated services, not local harness/CI/bootstrap details unless those are the product-facing contract.
    - Focus: if you find unrelated issues, log them in `docs/issues/<sequence>-<name>.md` and ignore them for the current task (unless they are blockers).
    - If `ww` fails or behaves unexpectedly during normal workflow use, capture it as a first-class workflow output per `docs/specs/ww-dogfooding-workflow.md` instead of silently dropping to raw git. Record the command, cwd, target repo, expected vs actual behavior, relevant output, fallback usage, and impact.
    - Issue Resolution: When an issue in `docs/issues/` is resolved, move the file to `docs/issues/done/`.
    - Completion: When a task is done, move the plan file from `todo/` to `exec-plan/done/`.
    - Post-Task Review: After completing significant work, run `post-task-review` to log issues, update lessons learned, and propose CLAUDE.md/AGENTS.md updates before creating a PR. Add new `docs/lessons.md` entries at the end of the file. (don't add lessons when the mistake is already the issue tracked or the solution planned)

## When asked to "Start a new feature":
1.  Create a planning worktree with global `ww`: `ww create plan/feature-name`, then enter it with `cd "$(ww cd plan/feature-name)"`
2.  Read `docs/project-plan.md`.
3.  Create a new numbered file in `docs/exec-plan/todo/`, such as `0007-feature-name.md`.
4.  Outline the objective, existing implementation references, code change map, and black-box spec changes in that plan.
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
2.  Create a numbered plan in `docs/exec-plan/todo/`, such as `0007-fix-bug-x.md`.
3.  Reproduction steps go into the plan.
4.  Create a PR for the plan and wait for review.
5.  After plan PR is merged, create an execution worktree with global `ww`: `ww create fix/fix-bug-x`, then enter it with `cd "$(ww cd fix/fix-bug-x)"`
6.  Execute the fix following **Spec First** and use the approved plan's implementation references and Code Change Map rather than rediscovering scope from scratch.
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
- Keep Git write operations on the main path unless separate worktrees clearly own them; parallel subagent work is fine for read-only exploration and verification, but shared-worktree Git mutations should stay serialized.

## Workspace Task Tracking

Workspace-level triage uses a GitHub Project plus a local derived cache managed by the `pj` CLI spike.

### Source of Truth

- GitHub Projects is the canonical remote state for workspace task triage
- `.local/pj/cache.json` is derived data for fast local reads and AI access
- `.local/pj/` is the only supported local workspace-triage state for the current workflow
- `docs/exec-plan/todo/` remains the canonical implementation-plan tracker once work is selected, using active filenames of the form `<sequence>-<name>.md`
- The workspace board is an owner-scoped GitHub ProjectV2; repository visibility is provided by linking that board to a same-owner repository's Projects tab

### Expected Commands

- `go -C tools/pj run ./cmd/pj init --owner <owner> --owner-type user|org`
- `go -C tools/pj run ./cmd/pj sync [--project <number>]`
- `go -C tools/pj run ./cmd/pj repo-link status <owner>/<repo>`
- `go -C tools/pj run ./cmd/pj repo-link add <owner>/<repo>`
- `go -C tools/pj run ./cmd/pj repo-link remove <owner>/<repo>`
- `go -C tools/pj run ./cmd/pj list`
- `go -C tools/pj run ./cmd/pj add --title "..." --status Todo`
- `go -C tools/pj run ./cmd/pj update --item <item-id> --status "In Progress"`
- `go -C tools/pj run ./cmd/pj url`
- `go -C tools/pj run ./cmd/pj config show|set|clear`

### Important Rules

- Use the GitHub Project flow for workspace-level task triage instead of any legacy local tracker flow
- Do not treat `.local/pj/` as source-of-truth data; it can be regenerated
- Do not depend on committed legacy tracker runtime artifacts or local database state for current task coordination
- Do not create duplicate local task trackers for the same workspace board unless a spec explicitly adds one
- Bootstrap the canonical board with `pj init` before relying on `sync`, `list`, `add`, `update`, `url`, or `open`
- Link the canonical board to the workspace repository with `pj repo-link add <owner>/<repo>` when repository Projects-tab discoverability is needed
- After bootstrap, reuse the stored owner scope instead of repeatedly restating `--owner` and `--owner-type`
- If you need to switch from a user board to an org board, do it explicitly with `pj config set` or `pj config clear` before re-running `pj init`; do not mix one-off owner flags with an existing local config

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE AND FOLLOW UP** - This is MANDATORY:
   ```bash
   git pull --rebase
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Post-PR follow-up** - For PR branches, confirm the PR is open, wait 30 seconds, inspect CI/checks and timeline for the latest head SHA, perform actionable CI fixes when feasible, and summarize any advisory bot/agent review findings in the current session before mutation. For Step 3 `execute-task` flows, if required CI checks are still pending after the first inspection, add up to two more `wait 30 seconds -> inspect again` turns before handoff unless another documented stop condition is reached earlier. Advisory bot/agent checks are reported by status but do not block handoff by default on second and later pushes unless required CI fails, review-start activity appears for the latest SHA, or the human requests more waiting. Do not silently apply findings or post triage to the PR unless explicitly asked.
6. **Clean up** - Clear stashes, prune remote branches
7. **Verify** - All changes committed, pushed, and followed up for the latest PR head SHA
8. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
- A PR handoff is incomplete until the latest pushed head SHA has completed the bounded `review-task` follow-up loop or a blocker/timeout is documented.
