# Follow Up: Retire Legacy `.beads/` Workspace Artifacts

## Context

`github-projects-task-cli-spike` switches workspace task tracking from beads to GitHub Projects, but this execution task does not remove the historical `.beads/` data directory or every legacy ignore pattern.

## Why It Is Deferred

- The current task is focused on proving the GitHub Projects flow end to end
- Deleting historical tracker artifacts is operational cleanup, not required for the spike to function
- We should decide whether any old beads history is worth preserving before deleting local state

## Follow-up Work

- Decide whether `.beads/` should be deleted, archived, or left ignored
- Remove obsolete beads-specific ignore rules and warnings once the cleanup path is chosen
- Sweep any remaining docs and helper flows for stale beads references
