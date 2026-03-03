---
name: plan-project
description: When starting a new project, defining project goals, writing or updating project requirements, setting the project vision, updating the project plan, defining what to build and why, or reviewing project scope and significance.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Plan Project (Workflow Step 1)

**Position in workflow**: This is **Step 1** of the AI-Centered Development cycle. This step requires its own branch and PR. After the project plan is merged, proceed to Step 2 (Execution Plan).

## What to Do

Define the project's **Goal**, **Significance**, and **Requirements** in `docs/project-plan.md`.

### Branch Setup

Before making any changes, create a fresh branch from the latest `main`:

```sh
git fetch origin
git switch -c plan/project-plan-<description> origin/main
```

Use a descriptive branch name (e.g., `plan/project-plan-v1`, `plan/project-plan-add-requirements`).

### Rules

1. **Read first**: If `docs/project-plan.md` already exists, read it before making changes.
2. **Single source of truth**: `docs/project-plan.md` is the authoritative document for what the project aims to achieve.
3. **Keep it updated**: As the project evolves, update this document to reflect new understanding of goals and requirements.

### Content Structure for `docs/project-plan.md`

The project plan should contain:

- **Goal**: What are we building? A concise statement of the project's purpose.
- **Significance**: Why does this matter? Business value, user impact, or technical motivation.
- **Requirements**: What must be true for the project to be considered complete? Functional and non-functional requirements.

### When to Update

- At project inception
- When requirements change
- When scope is adjusted
- After significant learnings that affect project direction

## PR Workflow

After making changes to `docs/project-plan.md`:

1. Commit the changes on the branch.
2. Push the branch and create a PR via `gh pr create`.
3. Wait for GitHub PR review approval before merging into `main`.

## Next Step

After the project plan PR is merged, proceed to **Execution Plan** (Step 2): create task plans in `docs/exec-plan/todo/` to break the work into implementable units.
