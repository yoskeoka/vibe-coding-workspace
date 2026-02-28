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
- `docs/specs/`: Detailed specifications. Must match implementation.
- `docs/exec-plan/`:
    - `todo/`: Active execution plans.
    - `done/`: Completed execution plans.
- `docs/references/`: Context for external tools/protocols (e.g., `fastapi-llms.txt`).
- `docs/issues/`: Local issue tracking. Avoids confusion with GitHub Issues during active "exec-plan" cycles.

## Workflow Cycle

### 1. Project Plan (`docs/project-plan.md`)
- Define the Goal, Significance, and Requirements.
- Update this as the project evolves.

### 2. Execution Plan (`docs/exec-plan/todo/`)
- Create a new plan file (e.g., `001-initial-setup.md`) in `todo/`.
- Detail:
    - Code changes.
    - Spec changes (How `docs/specs/` will change).
    - Break down large tasks into smaller sub-plans if needed.
- Review/Update `design-decisions/` if architectural choices are made.

### 3. Execution
- **Spec First**: Update `docs/specs/` *before* modifying code.
- **Implement**: Write the code to match the spec.
- **Issues**: If unrelated problems are found, log them in `docs/issues/<name>.md`. Do not fix them within the current plan unless blocking.
- **Completion**: Move the plan file from `docs/exec-plan/todo/` to `docs/exec-plan/done/` at the end of the work.

### 4. Review
- Create a PR.
- The PR must include:
    - Code changes.
    - Spec updates.
    - The plan file moved to `done/`.
    - Verification artifacts (test results, screenshots, logs) for human review.

repeat steps 2-4 until Project Plan is complete.
