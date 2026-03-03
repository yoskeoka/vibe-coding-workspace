# Spec: New Project Intake Workflow Package

## Goal
Provide a reusable pre-`plan-project` workflow package for adding new projects to this meta-repo.

## Scope
This package covers only the steps before detailed project planning:
- Idea sparring and rough articulation of pain points and desired experience
- Lightweight market/reference research
- Go/No-Go decision support
- New project bootstrap tasks in the workspace
- Handoff proposal to `plan-project`

## Requirements

### 1. New skill for pre-planning phase
A dedicated skill MUST exist to run this flow before `plan-project`.

Skill responsibilities:
1. Act as a sparring partner to clarify:
   - what problem should be solved
   - what experience should be delivered
2. Research existing solutions/experiences (e.g., games, products, tools).
3. Evaluate whether there is meaningful value to build:
   - novelty or differentiation
   - personal motivation / hobby value
   - reason users would care
4. Branch on decision:
   - **YES**: proceed to project creation/bootstrap
   - **NO**: summarize findings and append to rejected-ideas log
5. For YES case, execute project bootstrap:
   - create GitHub repository if missing
   - ensure local child repo exists
   - initialize docs/workflow scaffold in child repo
   - ensure workspace meta config includes the new project (`setup.sh`, `.gitignore`, README Managed Projects)
6. After bootstrap, explicitly propose continuing with `plan-project` in the new child repo.

### 2. Rejected idea log
Workspace MUST contain an English-named file to store no-go ideas and research outcomes:
- `docs/design-decisions/rejected-ideas.md`

The file should be append-only and include:
- date
- idea name
- summary
- why no-go now
- conditions to revisit
- references

### 3. Skill distribution boundary
The new skill is **meta-repo only** and MUST NOT be distributed to child repos via `setup-skills.sh`.

Rationale:
- This skill is for deciding whether to create a project and bootstrapping workspace-level assets.
- Child repos should start from `plan-project` and not carry pre-creation intake logic.

### 4. Documentation updates
Workflow docs SHOULD mention this pre-`plan-project` package and when to use it.

## Non-Goals
- Full product requirements definition (handled by `plan-project`)
- Feature-level implementation planning (handled by `plan-execution`)
- Implementation work (handled by `execute-task`)
