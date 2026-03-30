# ww Status Fails Hard When Default Base Cannot Be Confirmed

## Summary

Running `./ww/ww list` from the workspace can fail completely when one detected repository does not have `refs/remotes/origin/HEAD` configured, even if the repository is otherwise healthy and usable.

In the observed case, `ww` detected `/home/yoske/src/github.com/yoskeoka` as the workspace root and included `/home/yoske/src/github.com/yoskeoka/reversi-adventure` plus `/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace` as child repositories. The sibling `reversi-adventure` repository has `origin/main` and local branches tracking `origin/main`, but does not have `origin/HEAD` configured. `ww list` aborted instead of returning a partial or degraded status result.

## Reproduction

From `/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace`:

```text
./ww/ww list
resolving base branch for reversi-adventure: cannot detect default branch: git symbolic-ref refs/remotes/origin/HEAD: exit status 128
fatal: ref refs/remotes/origin/HEAD is not a symbolic ref
```

## Cause Analysis

- `ww list` computes `STATUS` for every worktree, not just a raw worktree listing.
- Status computation currently requires a resolved base branch so it can run `git branch --merged <base>`.
- Base resolution currently hard-depends on `default_base` config or `git symbolic-ref refs/remotes/origin/HEAD`.
- If that lookup fails for any detected repository, `listRepo()` returns an error and the entire `ww list` / `ww clean` command stops.
- This is too strict for repositories that are missing `origin/HEAD` but still expose enough information to describe worktrees safely.

## Why This Is A Problem

- `origin/HEAD` missing is not inherently a broken Git state.
- Workspace-wide commands become fragile: one repository with incomplete remote metadata can block visibility into all other repositories.
- `ww clean` depends on the same status pipeline, so cleanup can also be blocked.

## Expected Behavior

- `ww list` should continue and report worktrees even when the default base cannot be confirmed.
- Status should degrade to `unknown(...)` instead of hard-failing.
- `ww clean` and `ww list --cleanable` should continue to act only on confirmed `merged` and `stale` worktrees, never on `unknown`.

## Planned Resolution

This issue is addressed by [001-ww-unknown-status-default-base.md](/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/docs/exec-plan/todo/001-ww-unknown-status-default-base.md), which adds:

- `status=unknown`
- `status_detail` for degraded classification reasons
- best-effort default-base inference
- non-fatal handling when base or remote metadata cannot be confirmed
