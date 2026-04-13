# Issue: `pj init` does not provision custom ProjectV2 fields yet

## Summary

`pj init` now resolves or creates the canonical `Workspace Task Triage` board and writes its identity into `.local/pj/cache.json`, but the spike still relies on the board already having the workflow field model.

GitHub's default ProjectV2 setup does not guarantee the custom `Repo`, `Kind`, and `Priority` fields required by the workspace triage flow, so a freshly created board can bootstrap successfully and then fail validation with a clear missing-field error.

## Why This Is Separate

- This execution plan focused on removing the manual board-creation prerequisite.
- Custom field provisioning is a second slice of GitHub Projects mutation work and can be implemented in a follow-up plan without changing the bootstrap contract.

## Follow-up

- Add field-provisioning support for `Repo`, `Kind`, and `Priority` during `pj init`
- Decide whether the fields should be created as single-selects with canonical option sets or whether some should remain free-text in the spike
