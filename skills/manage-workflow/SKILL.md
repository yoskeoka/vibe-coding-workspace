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

The AI-Centered Development workflow consists of 4 repeating phases:

1. **Project Plan** — Define goals, significance, and requirements in `docs/project-plan.md`
2. **Execution Plan** — Create task plans in `docs/exec-plan/todo/`
3. **Execution** — Spec First → Implement → Log issues → Move plan to done
4. **Review** — Create PR with code, spec updates, moved plan, and verification artifacts

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

1. Run `setup-skills.sh` to set up submodule, symlinks, and docs/ templates (mechanical setup).
2. Then follow the steps below to configure CLAUDE.md for this project.

## Configuring CLAUDE.md and AGENTS.md

When initializing the workflow for a project, ensure both `CLAUDE.md` and `AGENTS.md` contain the necessary workflow instructions. **Do not overwrite existing content** — merge these into the existing files, preserving project-specific instructions.

- **CLAUDE.md**: Read by Claude Code. Add all sections below.
- **AGENTS.md**: Read by other AI tools (Cursor, GitHub Copilot, Gemini CLI, etc.). Add the Workflow Adherence section. The Subagent Strategy section is optional (tool-specific).

### Required content for CLAUDE.md and AGENTS.md

The following blocks must be present. If the file doesn't exist, create it. If it already exists, add only the missing sections.

#### 1. Workflow Adherence (required — both CLAUDE.md and AGENTS.md)

```markdown
## AI-Centered Development Workflow

This project follows the AI-Centered Development workflow.

### Core Responsibilities

1. **Workflow Adherence**:
   - NEVER skip the "Execution Plan" phase for non-trivial changes.
   - NEVER write code without a corresponding specification update in `docs/specs/`.

2. **Context Management**:
   - Your "memory" is the `docs/` directory.
   - `docs/project-plan.md` is your North Star.
   - `docs/exec-plan/todo/` is your current task list.
   - `docs/design-decisions/` is your architectural conscience.

3. **Execution Rules**:
   - **Plan First**: Before writing code, ensure a plan exists in `docs/exec-plan/todo/`. If not, create one.
   - **Spec First**: Update `docs/specs/` to reflect changes BEFORE modifying code.
   - **Focus**: If you find unrelated issues, log them in `docs/issues/<name>.md` and ignore them for the current task (unless they are blockers).
   - **Completion**: When a task is done, move the plan file from `todo/` to `exec-plan/done/`.
```

#### 2. Subagent Strategy (recommended — CLAUDE.md only)

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

1. Confirm `CLAUDE.md` contains the "Workflow Adherence" section with all 3 core responsibilities.
2. Confirm `AGENTS.md` contains the "Workflow Adherence" section with all 3 core responsibilities.
3. Confirm no existing project-specific instructions were removed from either file.
4. Verify `CLAUDE.md` and `AGENTS.md` don't have conflicting instructions.

## Templates

This skill includes templates in the `templates/` directory that provide the initial file contents for each document in the `docs/` structure.
