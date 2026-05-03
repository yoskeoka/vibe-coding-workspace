# Follow Up: `ww cd` Cannot Find a Just-Created Worktree

## Summary

During normal workflow startup for this task, the released global `ww` binary created the execution worktree successfully but failed to resolve that same branch immediately afterward.

- cwd: `/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace`
- command: `ww create feat/kb-bilingual-rendering`
- expected: create the worktree and allow `ww cd feat/kb-bilingual-rendering` to return that worktree path immediately
- actual:
  - `ww create feat/kb-bilingual-rendering` succeeded and reported `/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/.worktrees/vibe-coding-workspace@feat-kb-bilingual-rendering`
  - `ww cd feat/kb-bilingual-rendering` then failed with `no worktree found for branch "feat/kb-bilingual-rendering"`
- fallback used: continue from the created path reported by `ww create`
- impact: workflow dogfooding loses the documented `ww create` -> `ww cd` handoff and forces manual path recovery even though the worktree already exists

## Proposed Solution

- Reproduce the failure in `ww` with a workspace-root task on the meta-repo.
- Verify whether `ww cd` needs an explicit repo hint after `ww create`, or whether the created worktree was not indexed/discovered correctly.
- Fix `ww cd` so a newly created worktree is discoverable immediately from the same cwd without requiring manual path extraction.

## Priority

Medium. This does not block execution once the reported path is used manually, but it directly breaks the default dogfooding workflow that the workspace now documents and expects.

## Resolution

Resolved upstream in `ww` on `main`.

- fix commit: `720ab5d` `fix: cover git-backed workspace-root create/cd parity (#213)`
- regression test: `ww/integration_test.go` `TestCdFindsJustCreatedWorktreeFromGitBackedWorkspaceRoot`

Verified in this session with:

- `go test -C ww ./... -run TestCdFindsJustCreatedWorktreeFromGitBackedWorkspaceRoot -count=1`
