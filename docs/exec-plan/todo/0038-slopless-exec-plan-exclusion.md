# Exclude Execution Plans from Slopless CI

> **Execution**: Use `/execute-task` to implement this plan. After implementation is complete, use `/review-task` to prepare and create the PR.

## Objective

Stop Slopless warnings for Markdown under `docs/exec-plan/` across the workspace and managed child projects. CI must not start for a plan-only change or pass plan files when other Markdown starts the job. The local post-edit hook must also skip plan files.

The completion boundary is the merged workspace change, merged CI workflow updates in all six child projects that run Slopless, and merged workflow-sync updates to every repository in the seven-repository child sync matrix. Existing Slopless coverage for other Markdown remains intact.

## Context

The workspace and `ai-arena` workflows already limit their triggers and candidate selectors to selected documentation directories. The other five child workflows use a broad Markdown trigger and selection loop, so execution plans are currently included there. The local post-edit hook accepts edited Markdown paths without a plan-directory exclusion. The shared workflow submodule can deliver that hook change to child projects through its existing sync PR flow once the hook path is registered as child-consumed.

The report motivating this change is the [Slopless comment on reversi-adventure PR #228](https://github.com/yoskeoka/reversi-adventure/pull/228#issuecomment-5853375835), which reported 73 findings across four changed files, including active execution plans.

## References

- `docs/specs/slopless-ci.md` — current scope, changed-Markdown definition, root helper, and workflow reporting contract.
- `docs/specs/workflow-sync-to-child-repos.md` — child-consumed source paths and automatic submodule update PR contract.
- `tools/list-changed-markdown.sh` — root candidate selector; its path `case` currently permits only the documented long-lived Markdown directories.
- `.codex/hooks/slopless_post_tool_use.py:28-79,202-225` — `normalize_path`, `collect_paths`, and `main` select local post-edit Markdown candidates.
- `.github/workflows/sync-workflow-to-child-repos.yml` — the `git diff --quiet` child-consumed path list and `COMMIT_LIST` path list decide whether each child receives a workflow-sync PR and which source commits it reports.
- `.github/workflows/slopless.yml` — root PR path filters and changed-file helper invocation.
- Broad child `pull_request.paths` patterns (`:6-8`) and the `mapfile` / `case "$path"` candidate loop:
  - `dungeon-game-ai-arena/.github/workflows/slopless.yml`
  - `envdiff/.github/workflows/slopless.yml`
  - `reversi-adventure/.github/workflows/slopless.yml`
  - `reversi-ai-arena/.github/workflows/slopless.yml`
  - `vim-learning-game/.github/workflows/slopless.yml`
- `ai-arena/.github/workflows/slopless.yml:6-11` — scoped trigger and runtime allowlist, which must remain intact apart from the explicit plan exclusion.

## Change Map

- (MODIFY) `docs/specs/slopless-ci.md` — define `docs/exec-plan/` as a common exclusion for the workspace and managed child workflows.
- (MODIFY) `docs/specs/workflow-sync-to-child-repos.md` — register the shared Slopless hook as a child-consumed source path.
- (MODIFY) `tools/list-changed-markdown.sh` — explicitly skip `docs/exec-plan/` before applying the workspace path allowlist.
- (MODIFY) `.codex/hooks/slopless_post_tool_use.py` — reject execution-plan paths before local Slopless runs.
- (MODIFY) `.github/workflows/sync-workflow-to-child-repos.yml` — recognize the hook source path so child submodule sync PRs are created.
- (MODIFY) the six child `.github/workflows/slopless.yml` files listed above — explicitly exclude plan paths from runtime candidates; add a negative trigger pattern where the existing trigger is broad, while preserving each project's current non-plan scope.
- (MODIFY) workspace `.github/workflows/slopless.yml` only if review finds its existing positive path allowlist does not already prevent execution-plan-only runs.

## Black-Box Specification Changes

- Markdown rooted under `docs/exec-plan/` is never linted by Slopless in any managed repository.
- A pull request that changes only files under that path does not trigger Slopless.
- When a pull request changes both execution plans and other eligible Markdown, Slopless receives only the other eligible Markdown.
- Local post-edit checks skip execution-plan files, including when one edit event names both plan and non-plan Markdown.
- A shared local hook change is eligible for automatic child submodule sync.
- The workspace and each child project's current non-plan Markdown scope remains unchanged.

## Subtasks and Dependencies

1. Keep the common exclusion contract in `docs/specs/slopless-ci.md` ahead of implementation.
2. Make the workspace helper reject the execution-plan path explicitly while preserving its current allowlist. Add the same path rejection to the local post-edit hook.
3. Register the local hook source in the workflow-sync contract and path check so child repositories receive it through the vendored submodule.
4. Update all six current child CI workflows. Add the negative `pull_request.paths` pattern to workflows with broad Markdown triggers and a runtime path guard to every workflow; keep `ai-arena`'s narrower allowlist and the other five workflows' broad Markdown scope intact.
5. Inspect each target repository's current guidance and latest `main` before opening its fresh worktree. Prepare separate CI workflow PRs because those workflows are owned by separate repositories; review and merge the automated local-hook sync PRs for all seven sync-matrix repositories as well.
6. Use `review-task` for the workspace and child PRs; the aggregate task is complete only after each applicable PR is merged and its latest-head status is checked.

The workspace selector and child workflow updates are independent after the shared contract is approved. Git writes remain serialized within each repository. No new repository-wide Slopless workflow is introduced for projects that do not currently run one.

## Verification

- Run the workspace's applicable workflow, shell, Python syntax, and repository quality gates, including `git diff --check`.
- Validate each changed workflow with the repository's configured workflow linter and review the negative path pattern together with the runtime candidate guard.
- Confirm the local hook and root helper skip `docs/exec-plan/todo/example.md` while retaining an eligible non-plan Markdown path.
- Confirm the sync path check includes `.codex/hooks/slopless_post_tool_use.py` and inspect the generated child sync PRs.
- Confirm a plan-only CI diff yields no Slopless candidate while a mixed diff passes only eligible non-plan Markdown.
- Review the resulting Slopless check behavior on the latest PR heads; findings for eligible non-plan Markdown remain reportable.
