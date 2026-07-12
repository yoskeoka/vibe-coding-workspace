---
name: triage-tasks
description: Run a workspace daily briefing and reconcile the GitHub Project task queue.
---

# Triage Tasks

Use [task triage](../../AI_WORKFLOW.md#task-triage). This is a session-start
briefing, not automatic execution.

## Procedure

1. Bootstrap the canonical board once with `go -C tools/pj run ./cmd/pj init
   --owner <owner> --owner-type user|org`; otherwise run `pj sync`, then `pj
   list` and `pj url`. GitHub Projects is canonical; `.local/pj/` is cache.
2. Present a concise prioritized Todo shortlist, excluding routine
   `workflow-sync` maintenance from the default recommendation. Prefer active
   plans, active issues, broken workflow/review work, and the current repo.
3. On explicit full re-triage, collect each repository read-only: project-plan
   gaps, active plans/issues, open PRs, and open GitHub issues. Reconcile by
   source before mutation: mark stale items done, update equivalents, and add
   only genuinely missing work. Keep bodies compact: `Source`, `Repo`, `Next`,
   `Start`, `Read`, `Goal`.
4. Ask the human to pick, update, or re-triage. On selection, mark the Project
   item In Progress and provide a fresh-session prompt that names the target
   repo, `ww` command, first reads, goal, and next skill. Do not auto-execute.
