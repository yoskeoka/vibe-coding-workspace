# Manage AI Workflow Skill

This skill is responsible for synchronizing the AI-Centered Development workflow rules, directory structures, and templates to managed projects.

## Usage

Run the `run.sh` script to propagate the workflow configuration to all sub-projects in the workspace.

```bash
./skills/manage-workflow/run.sh [optional-target-directory]
```

## Maintenance

- `templates/`: Contains the directory structure and placeholder markdown files copied to `docs/`.
- `run.sh`: The logic for copying files and ensuring consistency.
