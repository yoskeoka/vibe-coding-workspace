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

1. Copy the templates from `skills/manage-workflow/templates/docs/` into the project's `docs/` directory.
2. Fill in `docs/project-plan.md` with the project's goals, significance, and requirements.
3. Begin the workflow cycle starting with an Execution Plan.

## Templates

This skill includes templates in the `templates/` directory that provide the initial file contents for each document in the `docs/` structure.
