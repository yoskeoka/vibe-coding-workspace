# Manage AI Workflow Skill

This skill initializes and manages the AI-Centered Development workflow structure for a project. It provides the `docs/` directory templates for bootstrapping new projects.

## Related Skills

The workflow is split into step-specific skills for precise triggering:

| Skill | Workflow Step | Trigger |
|---|---|---|
| `skills/manage-workflow/` | Setup | Project workspace initialization, scaffolding |
| `skills/new-project-intake/` | Pre-Step 1 | New project idea sparring, go/no-go, bootstrap handoff |
| `skills/new-project-intake/` | Pre-Step 1 (workspace-only) | New project idea sparring, go/no-go, bootstrap handoff before child repo creation |
| `skills/plan-project/` | Step 1 | Defining goals, requirements, project vision |
| `skills/plan-execution/` | Step 2 | Planning features, bug fixes, creating task plans |
| `skills/execute-task/` | Step 3 | Implementing code, updating specs, logging issues |
| `skills/review-task/` | Step 4 | Creating PRs, preparing verification artifacts |

## Structure


## Usage in Sub-Projects

To use this workflow in a sub-project:

1.  Clone this repository (or subtree) into the sub-project's skills directory.
2.  Or use `setup.sh` in the meta-repo to manage the sub-project alongside this repo.

Note: `new-project-intake` is intentionally not installed into child repos because it is only used in the meta-repo before project creation.
