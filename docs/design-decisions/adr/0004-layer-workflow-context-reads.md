# Layer workflow context reads by task

## Status

Accepted — 2026-07-12

## Context

Long, overlapping workflow docs, skills, and PR templates made routine work
spend context on procedures unrelated to the current phase. The workflow still
needs explicit enforcement and durable routing.

## Decision

Use a three-layer contract: `AGENTS.md` supplies universal guards and routing,
`AI_WORKFLOW.md` supplies linkable lifecycle rules, and skills supply their own
procedures. Index ADRs and on-demand specs; keep PR template completion
conditional on the change type.

## Consequences

- Routine tasks can read a smaller, relevant set of documents.
- Lifecycle enforcement remains centralized and linkable.
- Future workflow additions must name their owner and avoid duplicating shared
  branch, PR, or completion prose.
