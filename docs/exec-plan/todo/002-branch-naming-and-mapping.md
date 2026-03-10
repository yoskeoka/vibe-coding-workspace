# 002: Branch Naming Convention & Exec-Plan Mapping

## Objective

Formalize the branch naming convention and the branch-to-exec-plan mapping rule in `AI_WORKFLOW.md`. This unblocks the workflow linter checks that are currently out-of-scope in 001-workflow-linter (branch naming, exec-plan existence, exec-plan completion).

Addresses:
- Project plan requirement: "Branch naming convention is formally declared and enforceable"
- Project plan requirement: "Exec-plan-to-branch mapping convention is formally declared"
- Issues: `docs/issues/branch-naming-rule.md`, `docs/issues/exec-plan-branch-mapping.md`

## Branch Naming Rule

Branch names MUST match the pattern:

```
<type>/<description>
```

| Type   | Purpose                                     | Example                        |
|--------|---------------------------------------------|--------------------------------|
| `plan` | Execution plan creation/update              | `plan/002-feature-name`        |
| `feat` | Feature implementation (from an exec-plan)  | `feat/002-feature-name`        |
| `fix`  | Bug fix implementation (from an exec-plan)  | `fix/003-bug-name`             |
| `chore`| Non-functional changes (CI, tooling, deps)  | `chore/update-ci`              |
| `docs` | Documentation-only changes                  | `docs/update-readme`           |

### Description format

- For branches associated with an exec-plan: `<NNN>-<short-name>` where `<NNN>` is the exec-plan number (zero-padded to 3 digits).
- For branches without an exec-plan (chore, docs, trivial): free-form kebab-case.

## Exec-Plan-to-Branch Mapping Rule

When a branch name contains a plan number (`<type>/<NNN>-*`), the corresponding exec-plan file MUST exist:

- Planning phase: `docs/exec-plan/todo/<NNN>-*.md` must exist.
- After execution is complete: the file must be moved to `docs/exec-plan/done/<NNN>-*.md`.

Branches of type `chore` and `docs` are exempt from this rule (no exec-plan required).

## Code Changes

### `AI_WORKFLOW.md`

Update the "Branch Setup" section to replace the current informal examples with the formal rules declared above:

1. Add a "Branch Naming Convention" subsection with the type table and description format rules.
2. Add an "Exec-Plan Mapping" subsection declaring the branch-to-plan-file mapping rule.

### `CLAUDE.md` (AGENTS.md equivalent)

Update the "Branch & PR Rules" section to reference the new formal rules instead of using ad-hoc examples. No new rules needed — just point to AI_WORKFLOW.md as the source of truth.

## Spec Changes

None — no `docs/specs/` file needed. The rules live in `AI_WORKFLOW.md` which is the workflow specification.

## Sub-tasks

- [ ] [parallel] Update `AI_WORKFLOW.md` with branch naming convention and exec-plan mapping rules
- [ ] [parallel] Update `CLAUDE.md` to reference the formal rules
- [ ] [depends on: both above] Move `docs/issues/branch-naming-rule.md` and `docs/issues/exec-plan-branch-mapping.md` to `docs/issues/done/`
