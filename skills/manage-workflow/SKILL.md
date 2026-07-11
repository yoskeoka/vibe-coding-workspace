---
name: manage-workflow
description: Bootstrap or maintain the AI-Centered Development workflow in a repository.
---

# Manage Workflow

Use [workflow setup](../../AI_WORKFLOW.md#workflow-setup) for lifecycle rules.

## Procedure

1. Run `setup-workspace.sh` for mechanical submodule, skill-link, template, and
   hook setup. Do not overwrite existing project docs.
2. Ensure `AGENTS.md` is the short universal entrypoint and `CLAUDE.md` is its
   symlink. Add project-specific routing below the compact workflow guards.
3. Initialize the template tree: project plan, specs, plans, issues, references,
   `docs/design-decisions/README.md`, and `docs/design-decisions/adr/` records.
   ADRs use the Michael Nygard form: title, Status, Context, Decision,
   Consequences.
4. Confirm templates and guidance point to the compact-read contract; submit
   setup changes through `review-task`.
