---
name: manage-workflow
description: When setting up a new project workspace, initializing the AI-Centered Development workflow, creating the docs/ directory structure, scaffolding a project from scratch, or bootstrapping the AI workflow for a new repository.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Manage AI-Centered Development Workflow

This skill initializes and manages the AI-Centered Development workflow structure for a project.

## Overview

The AI-Centered Development workflow consists of 3 repeating steps, each requiring its own branch and GitHub PR review:

1. **Project Plan** — Define goals, significance, and requirements in `docs/project-plan.md` (branch + PR)
2. **Execution Plan** — Create task plans in `docs/exec-plan/todo/` (branch + PR)
3. **Execution** — Spec First → Implement → Verify (lint/test) → PR with code, spec updates, moved plan, and verification artifacts

Every change goes through: **Branch from latest `main`** → **Work** → **Verify** → **PR** → **Review** → **Merge**.

## Core Principles

### Process Principles

1. **AI-Centric Context**: All necessary information must be immediately accessible to AI in the `docs/` directory. Files should be structured for easy parsing and context retrieval.
2. **Spec-Code Parity**: `docs/specs/` must strictly match the actual code. No PR is reviewable without verification that specs and code are in sync.
3. **Verification First**: Human review happens _after_ mechanical tests and "visual" verification data are ready.

### Code Quality Principles

4. **Simplicity First**: Make every change as simple as possible. Prefer the straightforward approach.
5. **No Laziness**: Find root causes. No temporary fixes. Maintain senior developer standards.
6. **Minimal Impact**: Only touch what's necessary. Avoid introducing bugs through unnecessary changes.

## Directory Structure to Initialize

When setting up a new project, create the following structure under `docs/`:

```
docs/
  project-plan.md          # Goals, significance, requirements
  design-decisions/
    adr.md                 # Append-only architectural decision log
    core-beliefs.md        # Guiding principles for trade-offs
  specs/                   # Detailed specifications (must match implementation)
  exec-plan/
    todo/                  # Active execution plans
    done/                  # Completed execution plans
  references/              # External context (e.g., fastapi-llms.txt)
  issues/                  # Local issue tracking
```

## Setup Instructions

1. Run `setup-workspace.sh` to set up submodule, symlinks, and docs/ templates (mechanical setup).
2. Then follow the steps below to configure CLAUDE.md for this project.

## Updating the Workflow Submodule

Child repos that use `.claude/vendor/workflow` as a submodule will frequently see diffs like:

```
diff --git a/.claude/vendor/workflow b/.claude/vendor/workflow
--- a/.claude/vendor/workflow
+++ b/.claude/vendor/workflow
@@ -1 +1 @@
-Subproject commit abc1234...
+Subproject commit def5678...
```

### Quick update (recommended)

From inside the child repository, run:

```bash
/path/to/vibe-coding-workspace/setup-workspace.sh --update
```

This will:
1. Fetch the latest workflow commit from remote
2. Update the submodule pointer
3. Auto-commit the change with message `chore: update workflow submodule to <sha>`

Then push with `git push`.

### Update all child repos at once

If you manage multiple child repos, create a simple loop:

```bash
WORKSPACE="/path/to/vibe-coding-workspace"
for repo in /path/to/child-repo-1 /path/to/child-repo-2; do
  echo "=== Updating $repo ==="
  "$WORKSPACE/setup-workspace.sh" --update "$repo"
  (cd "$repo" && git push)
done
```

### When to update

- **Before starting a new task**: ensures skills and hooks are current.
- **When you see a dirty submodule diff**: run `--update` to apply and commit it cleanly.
- **After pushing changes to the workflow repo**: update child repos to pick up the new version.

## Configuring AGENTS.md and CLAUDE.md

`AGENTS.md` is the single source of truth for AI agent instructions. `CLAUDE.md` is a symlink pointing to `AGENTS.md`, so both Claude Code and other AI tools (Cursor, GitHub Copilot, Gemini CLI, etc.) read the same content.

### Symlink strategy

- **`AGENTS.md`**: The canonical file. All workflow instructions go here.
- **`CLAUDE.md`**: A symlink to `AGENTS.md`. Do NOT create as a separate file.

When initializing the workflow:

1. If `AGENTS.md` does not exist, create it with the required content below.
2. If `AGENTS.md` already exists, merge the missing sections — **do not overwrite** existing project-specific instructions.
3. If `CLAUDE.md` does not exist, create it as a symlink: `ln -s AGENTS.md CLAUDE.md`
4. If `CLAUDE.md` exists as a regular file (not a symlink), replace it with a symlink:
   - Merge any unique content from `CLAUDE.md` into `AGENTS.md`
   - Remove the regular `CLAUDE.md` file
   - Create the symlink: `ln -s AGENTS.md CLAUDE.md`

### Required content for AGENTS.md

The following blocks must be present. Add only the missing sections.

#### 1. Workflow Adherence (required)

```markdown
## AI-Centered Development Workflow

This project follows the AI-Centered Development workflow.

### Core Responsibilities

1. **Workflow Adherence**:
   - NEVER skip the "Execution Plan" phase for non-trivial changes.
   - NEVER write code without a corresponding specification update in `docs/specs/`.
   - ALWAYS create a new branch from the latest `main` before starting any work.
   - ALWAYS go through GitHub PR review for every change — including doc-only changes.

2. **Branch & PR Rules**:
   - Create a fresh branch from `origin/main` for every task: `git fetch origin && git switch -c <branch-name> origin/main`
   - Never reuse an existing feature branch.
   - Run all lint and test checks (non-AI tooling) before creating a PR.
   - Create PRs via `gh pr create` using the **PR template** (project-level `.github/PULL_REQUEST_TEMPLATE.md` if present, otherwise workspace-level) and wait for review approval before merging.

3. **Context Management**:
   - Your "memory" is the `docs/` directory.
   - `docs/project-plan.md` is your North Star.
   - `docs/exec-plan/todo/` is your current task list.
   - `docs/design-decisions/` is your architectural conscience.

4. **Execution Rules**:
   - **Plan First**: Before writing code, ensure a plan exists in `docs/exec-plan/todo/`. If not, create one.
   - **Spec First**: Update `docs/specs/` to reflect changes BEFORE modifying code.
   - **Focus**: If you find unrelated issues, log them in `docs/issues/<name>.md` and ignore them for the current task (unless they are blockers).
   - **Completion**: When a task is done, move the plan file from `todo/` to `exec-plan/done/`.
```

#### 2. Subagent Strategy (recommended)

```markdown
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
```

### Verification steps after configuration

1. Confirm `AGENTS.md` exists as a regular file and contains all required sections.
2. Confirm `CLAUDE.md` is a symlink pointing to `AGENTS.md` (verify with `ls -la CLAUDE.md`).
3. Confirm no existing project-specific instructions were removed from `AGENTS.md`.

## Templates

This skill includes templates in the `templates/` directory that provide the initial file contents for each document in the `docs/` structure.
