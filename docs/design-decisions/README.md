# Architecture Decision Records

Records are immutable after acceptance. Create the next numbered file in `adr/`
using the Michael Nygard form: title, Status, Context, Decision, Consequences.

| ID | Status | Tags | Outcome |
| --- | --- | --- | --- |
| [0001](adr/0001-skip-child-workflow-prs-without-skill-changes.md) | Accepted | workflow, sync | Child sync PRs are created only for `skills/` changes. |
| [0002](adr/0002-store-knowledge-base-in-repo.md) | Accepted | knowledge-base, docs | The knowledge base lives under `docs/kb/`. |
| [0003](adr/0003-dogfood-released-global-ww.md) | Accepted | workflow, ww | Normal workflow startup uses released global `ww`. |
| [0004](adr/0004-layer-workflow-context-reads.md) | Accepted | workflow, context | Universal guards, lifecycle, and procedures are layered by read need. |
