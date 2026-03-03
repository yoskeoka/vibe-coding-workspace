# AI-Centered Development Workflow

This document outlines the workflow for developing projects with AI as the central driver.

## Core Principles

1.  **AI-Centric Context**: All necessary information must be immediately accessible to AI in the `docs/` directory. Files should be structured for easy parsing and context retrieval.
2.  **Spec-Code Parity**: `docs/specs/` must strictly match the actual code. No PR is reviewable without verification that specs and code are in sync.
3.  **Verification First**: Human review happens *after* mechanical tests and "visual" verification data are ready.

## Directory Structure

- `docs/project-plan.md`: The single source of truth for the project's goals, significance, and requirements.
- `docs/design-decisions/`:
    - `adr.md`: Append-only log of architectural decisions.
    - `core-beliefs.md`: Guiding principles for trade-offs.
    - `rejected-ideas.md`: Append-only log for no-go project ideas (with rationale and revisit conditions).
- `docs/specs/`: Detailed specifications. Must match implementation.
- `docs/exec-plan/`:
    - `todo/`: Active execution plans.
    - `done/`: Completed execution plans.
- `docs/references/`: Context for external tools/protocols (e.g., `fastapi-llms.txt`).
- `docs/issues/`: Local issue tracking. Avoids confusion with GitHub Issues during active "exec-plan" cycles.

## Workflow Cycle

### 0. New Project Intake (Pre-Step, optional but recommended for vague ideas)
- Use this when an idea is still fuzzy and not ready for a full `project-plan`.
- Activities:
    1. Idea sparring (pain points, desired experience, target users)
    2. Existing-solution research
    3. Go/No-Go decision
- If **No-Go**: append findings to `docs/design-decisions/rejected-ideas.md` and stop.
- If **Go**: create/bootstrap the child project repo, update workspace meta entries, then continue to Step 1 (`plan-project`).

> **Rule**: Every step that produces changes MUST go through a GitHub PR review — including doc-only changes like Project Plan and Execution Plan updates. AI Agents must always create a new clean branch from the latest `main` before starting any work.

### Branch Setup (applies to every step below)
- **Always** start from the latest `main`:
    ```sh
    git fetch origin
    git switch -c <branch-name> origin/main
    ```
- Use a descriptive branch name (e.g., `plan/project-plan-update`, `plan/002-feature-x`, `feat/002-feature-x`, `fix/003-bug-x`).
- Never reuse an existing feature branch; always create a fresh one.

### PR Workflow (applies to every step below)
1. **Verify** — Run **all** project lint and test commands using non-AI tooling (e.g., `make lint`, `npm run lint`, `go vet`, `pytest`, `npm test`, or whatever the project defines). If any check fails, fix the issue in the same branch and re-run until **all pass**. Skip this for doc-only PRs when no lint/test tooling covers documentation.
2. **Create PR** — Push the branch and create a PR via `gh pr create`.
3. **Review** — Wait for GitHub PR review approval before merging into `main`.

---

### 1. Project Plan (`docs/project-plan.md`) — **requires PR**
- Create a new branch (e.g., `plan/project-plan-v1`).
- Define or update the Goal, Significance, and Requirements.
- Follow the **PR Workflow** above to merge the plan into `main`.
- Update this as the project evolves (each update = new branch + PR).

### 2. Execution Plan (`docs/exec-plan/todo/`) — **requires PR**
- Create a new branch (e.g., `plan/001-initial-setup`).
- Create a new plan file (e.g., `001-initial-setup.md`) in `todo/`.
- Detail:
    - Code changes.
    - Spec changes (How `docs/specs/` will change).
    - Break down large tasks into smaller sub-plans if needed.
- Review/Update `design-decisions/` if architectural choices are made.
- Follow the **PR Workflow** above to merge the plan into `main`.

### 3. Execution — **requires PR**
- Create a new branch (e.g., `feat/001-initial-setup`).
- **Spec First**: Update `docs/specs/` *before* modifying code.
- **Implement**: Write the code to match the spec.
- **Issues**: If unrelated problems are found, log them in `docs/issues/<name>.md`. Do not fix them within the current plan unless blocking.
- **Completion**: Move the plan file from `docs/exec-plan/todo/` to `docs/exec-plan/done/`.
- Follow the **PR Workflow** above (Verify → Create PR → Review).
- The PR must include:
    - Code changes.
    - Spec updates.
    - The plan file moved to `done/`.
    - Verification artifacts (test results, screenshots, logs) for human review.

Repeat steps 1–3 until the Project Plan is complete.
