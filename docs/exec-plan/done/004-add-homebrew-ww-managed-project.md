**Execution**: Use `/execute-task` to implement this plan.

# Add `homebrew-ww` to the managed projects registry

## Objective

Register `https://github.com/yoskeoka/homebrew-ww` as a managed project in the workspace so it is cloned, ignored, and documented alongside the other child repositories.

## Code Changes

- Update `setup.sh` so `homebrew-ww` is included in the managed repository list.
- Update `.gitignore` so the `homebrew-ww/` working tree is ignored when present locally.
- Update `README.md`, `AGENTS.md`, and `docs/project-plan.md` so the managed-project list stays consistent.

## Spec Changes

- Add a workspace-level spec describing the managed-project registry and the required documentation/configuration sync when a project is added.

## Sub-tasks

- [parallel] Draft the workspace managed-project registry spec.
- [depends on: managed-project registry spec] Update the workspace metadata files and repository list.
- [depends on: workspace metadata update] Verify the new project appears consistently in docs and setup configuration.

## Notes

- No architectural decision is expected for this change.
- The local `homebrew-ww/` checkout is not required for the registry update itself; `setup.sh` remains the source of truth for cloning.
