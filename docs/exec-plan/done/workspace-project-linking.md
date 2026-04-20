# Workspace Project Linking Clarification

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Clarify and operationalize how the workspace triage ProjectV2 is meant to work now that `tools/pj` has created the canonical `Workspace Task Triage` board under the configured owner account.

The goal is to remove the remaining ambiguity between:

- the Project owner scope (`user` / `org`)
- the repository where the workspace lives (`vibe-coding-workspace`)
- the repository UI surface where the board should be discoverable (`vibe-coding-workspace` repository Projects tab)

This plan supports the project-plan requirements that workspace task coordination uses a lightweight GitHub Projects-backed flow and that the CLI/workflow remain practical and low-friction for real use.

## Background

Current specs and guidance say the workspace uses a dedicated ProjectV2 named `Workspace Task Triage`, but the wording still makes it easy to assume the board itself is repository-owned.

GitHub Projects (ProjectV2) are owned by a user or organization, not by a repository. Repositories can instead link an owner-scoped project into the repository's `Projects` tab and can set a default repository for that project. GitHub's official documentation for "Adding your project to a repository" states that only projects owned by the same user or organization as the repository can be linked.

That means the current configured-owner `Workspace Task Triage` board is structurally valid for this repository, but the docs and workflow should make the distinction explicit and the repository-level linking step should be part of the operating model.

## Spec Changes

### `docs/specs/triage-tasks.md`

- Clarify that the canonical workspace board is owner-scoped (`user` or `org`), not repository-owned
- Document that the board should also be linked to a same-owner repository for discoverability, with `yoskeoka/vibe-coding-workspace` as this workspace's verification target
- Document the same-owner constraint for repository linking
- Clarify whether setting the repository as the project's default repository is required or only recommended

### `docs/specs/github-projects-task-cli.md`

- Clarify that `pj init` creates or resolves an owner-scoped ProjectV2
- Specify that `tools/pj` supports both:
  - checking whether the canonical Project is linked to a target repository
  - linking the canonical Project to a target repository
- Clarify that repository linking is exposed through `pj repo-link status|add|remove <owner>/<repo>` and is separate from `pj init`
- Capture any remaining safety gap around `pj init` creating a board under the wrong owner scope when the operator intended a different owner (`user` vs `org`)

### `AGENTS.md`

- Update workspace task tracking guidance so "workspace board" clearly means an owner-scoped Project linked to a same-owner repository for repository-tab visibility
- Keep the expected commands accurate if linking remains a manual GitHub UI step

### `README.md`

- Add or refine operator-facing guidance describing:
  - why the board appears under the user/org Projects page
  - how to link it to the repository Projects tab
  - any default-repository recommendation for this workspace

## Code Changes

### `tools/pj/`

- Add `pj repo-link status <owner>/<repo>` to inspect whether the canonical Project is already linked to a repository
- Add `pj repo-link add <owner>/<repo>` to link the canonical Project to a repository using GitHub's ProjectV2 GraphQL API
- Add `pj repo-link remove <owner>/<repo>` so the command group has a natural inverse operation even though this plan's verification does not remove the workspace link
- Keep repository linking separate from `pj init`
- Evaluate whether `pj init` should remain unchanged for owner targeting in this plan or whether the owner-selection safety gap should become a follow-up issue or a separate execution plan

## Design Decisions

- The canonical `Workspace Task Triage` board remains owner-scoped because ProjectV2 is not repository-owned
- The `vibe-coding-workspace` repository should expose the canonical board through its Projects tab by linking the owner-scoped board, but the CLI should accept any same-owner `<owner>/<repo>` target
- GitHub's public ProjectV2 GraphQL API should be used for repository link inspection/link creation rather than relying on a manual-only UI workflow
- If setting a default repository changes issue-creation behavior in ways that matter to the workspace, capture that explicitly in the specs during execution
- If `pj init` owner-targeting safety is not fixed in this plan, log it explicitly as a follow-up issue or separate plan rather than leaving it implicit

## Sub-tasks

- [x] Update `docs/specs/triage-tasks.md` to distinguish project ownership from repository visibility/linking
- [x] Update `docs/specs/github-projects-task-cli.md`, `AGENTS.md`, and `README.md` so they consistently describe the owner-scoped board model and the new repo-link capabilities
- [x] [parallel] Verify the exact GraphQL path for ProjectV2 repository-link inspection and link creation
- [x] [parallel] Verify how setting `vibe-coding-workspace` as the default repository affects project behavior and whether the workflow should require it
- [x] [depends on: repository-link GraphQL verification] Implement `tools/pj` support for checking whether the canonical Project is linked to a repository with `pj repo-link status <owner>/<repo>`
- [x] [depends on: repository-link GraphQL verification] Implement `tools/pj` support for linking the canonical Project to a repository with `pj repo-link add <owner>/<repo>`
- [x] [depends on: repository-link GraphQL verification] Implement `tools/pj` support for unlinking the canonical Project from a repository with `pj repo-link remove <owner>/<repo>`
- [x] [depends on: repo-link status command, repo-link mutation] Link `Workspace Task Triage` into the `vibe-coding-workspace` repository Projects tab using `tools/pj`
- [x] [depends on: docs/spec updates] Decide whether `pj init` owner-targeting safety belongs in this execution PR or should be logged as a follow-up issue / separate plan
- [x] [depends on: repository link applied] Verify the linked board is discoverable from both the owner's Projects page and the repository Projects tab

Owner-targeting safety decision: keep `pj init` owner targeting unchanged in this execution PR. The existing `pj config set` / `pj config clear` guard already rejects silent owner drift, and no new owner-selection gap was introduced by repository linking.

## Verification

- Confirm docs/specs no longer imply that the canonical board is repository-owned
- Confirm `tools/pj` can report whether the canonical board is linked to `vibe-coding-workspace`
- Confirm `tools/pj` can link the canonical board to `vibe-coding-workspace`
- Confirm the canonical board is linked from the `vibe-coding-workspace` repository Projects tab
- Confirm any default-repository decision is reflected consistently in docs and observed GitHub behavior
- Confirm `tools/pj` guidance still matches the actual bootstrap and day-to-day workflow after the clarification
- Confirm the owner-targeting safety decision for `pj init` is either implemented or explicitly recorded as follow-up work

## Expected Outcome

- The workspace workflow clearly treats `Workspace Task Triage` as an owner-scoped ProjectV2 linked to this repository
- Operators no longer misread the owner's Projects page as an implementation mistake
- The repository gains a discoverable Projects-tab entry point to the canonical workspace board without changing the board's source of truth
- `tools/pj` can both inspect and establish the repo link instead of relying on a manual GitHub UI step
