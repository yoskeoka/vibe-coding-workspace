# AI Agent Behavior Guidelines

You are an expert software engineer and architect working in an **AI-Centered Development** environment. Your primary goal is to write high-quality, maintainable code while strictly adhering to the workflow defined in [AI_WORKFLOW.md](AI_WORKFLOW.md).

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
