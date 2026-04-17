# KB Tools Concurrent Build Race

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Make `tools/kb check`, `tools/kb build`, and `tools/kb serve` robust when more than one invocation runs at the same time.

This resolves `docs/issues/kb-tools-concurrent-build-race-follow-up.md`, where parallel `check` and `build` runs both used `.local/kb-generated` and one invocation removed generated inputs while the other was still reading them. The fix supports the project-plan requirement for a repo-native knowledge base that can be built and verified reliably by automated and agentic workflows.

## Reproduction

Current race-prone behavior can be reproduced by running two KB commands at once from the workspace root, for example:

```sh
tools/kb check & tools/kb build
```

Known observed failures include:

- `OSError: [Errno 39] Directory not empty: '.local/kb-generated'`
- MkDocs config errors when `.local/kb-generated/docs` disappears while another process is building from it

Sequential runs of the same commands succeed, so the fix should target shared generated-path coordination rather than content or MkDocs configuration.

## Code Changes

### `tools/kb`

- Replace the globally shared generated workspace passed to each command with an invocation-scoped generated root, or add a lock if a stable shared path is still required.
- Prefer invocation-scoped generated roots for `check` and `build` so independent verification tasks can run in parallel without blocking each other.
- Ensure `serve` keeps its generated root alive for the lifetime of the MkDocs process.
- Keep `UV_CACHE_DIR` shared and stable unless verification shows it is part of the race; the known failure is the generated docs/config workspace, not the dependency cache.
- Clean up per-run generated roots after successful `check` and `build` runs without deleting another process's workspace.
- If cleanup fails, report the exact generated path and leave a local-only artifact under `.local/` for inspection rather than masking the build result.

### `tools/kb_generate.py`

- Accept the generated root from the caller, such as through an environment variable or CLI argument, instead of hardcoding `.local/kb-generated`.
- Continue deriving `docs`, generated config, and site output paths from that generated root unless the site output intentionally remains `.site/kb`.
- Avoid removing a shared parent directory during generation; only replace the invocation-owned generated root.
- Preserve deterministic generated MkDocs config and generated source indexes for a single invocation.

### Tests or Verification Helpers

- Add a small automated regression check for concurrent KB command execution if the existing tooling can support it without adding heavy dependencies.
- At minimum, document and run a shell-level verification that starts `tools/kb check` and `tools/kb build` concurrently and confirms both exit successfully.
- Continue running the normal sequential `tools/kb check` and `tools/kb build` verification because they cover the expected user path and strict MkDocs validation.

## Spec Changes

### `docs/specs/knowledge-base.md`

- Define that generated KB build inputs are invocation-scoped or otherwise protected from concurrent mutation.
- Clarify that `check`, `build`, and `serve` must not delete or invalidate another live KB invocation's generated docs/config files.
- Document the cleanup behavior for local generated artifacts under `.local/`.
- Document any explicit concurrency limitation if `serve` keeps a long-lived generated workspace and should not share it with one-shot commands.

## Issue Lifecycle

- Move `docs/issues/kb-tools-concurrent-build-race-follow-up.md` to `docs/issues/done/` during execution after the fix and verification land.

## Design Decisions

Past belief: `Correctness over Speed` means specs must define the concurrency contract before changing code. Apply that here by first documenting whether generated KB inputs are per-run or locked.

Past decision: the KB lives in-repo and publishes from the same Markdown source of truth. Apply the same reasoning here by keeping durable source files under `docs/kb/` and treating generated MkDocs inputs as disposable local artifacts.

No new ADR is expected unless execution chooses a broader shared cache or locking abstraction that affects more than KB generation.

## Sub-tasks

- [ ] [parallel] Update `docs/specs/knowledge-base.md` with the KB command concurrency and cleanup contract
- [ ] [parallel] Add or design a focused regression verification for concurrent `tools/kb` invocations
- [ ] [depends on: spec] Update `tools/kb_generate.py` so the generated root is provided by the caller and only invocation-owned paths are replaced
- [ ] [depends on: kb_generate] Update `tools/kb` to create, pass, retain, and clean invocation-scoped generated paths for `check`, `build`, and `serve`
- [ ] [depends on: tools/kb, regression verification] Run sequential and concurrent KB verification
- [ ] [depends on: verification] Move the resolved issue file to `docs/issues/done/`
