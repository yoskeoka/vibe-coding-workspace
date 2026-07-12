# Spec: Workflow Context Contract

## Observable contract

Workflow guidance is layered so an agent can identify mandatory behavior without
reading unrelated procedures.

- `AGENTS.md` is the short always-read entrypoint. It contains non-negotiable
  guards, workspace routing, and a task-to-document table.
- `AI_WORKFLOW.md` is the canonical, linkable lifecycle reference. Agents read
  the section for their current phase, not the entire workflow by default.
- Workflow skills own their phase procedure and link back to lifecycle sections
  instead of repeating branch, PR, and completion rules.
- `docs/specs/` and `docs/design-decisions/README.md` are indexed on-demand
  references. A task reads only the relevant spec and decision records.
- PR authors fill only template sections applicable to their change. `review-task`
  identifies conditional issue-closure, evidence, and follow-up details.

## Preserved enforcement

The compact entrypoint does not relax the workflow: fresh `ww` worktrees,
non-trivial planning, spec-first changes, issue lifecycle, PR review, quality
gates, and latest-head follow-up remain mandatory. `tools/workflow-lint.sh`
checks the structural contract for this workspace without imposing generic
line-count limits.

## ADR and lessons records

`docs/design-decisions/README.md` indexes immutable numbered records in
`docs/design-decisions/adr/`. Each record uses Michael Nygard's `Status`,
`Context`, `Decision`, and `Consequences` form. `docs/lessons.md` is an active
exceptions register: no more than ten unresolved recurring risks, each with a
canonical remediation link and review trigger. Git history is the archive.
