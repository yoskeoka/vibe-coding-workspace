# Branch Naming Convention & Exec-Plan Mapping

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
| `plan` | Execution plan creation/update              | `plan/feature-name`            |
| `feat` | Feature implementation (from an exec-plan)  | `feat/feature-name`            |
| `fix`  | Bug fix implementation (from an exec-plan)  | `fix/bug-name`                 |
| `chore`| Non-functional changes (CI, tooling, deps)  | `chore/update-ci`              |
| `docs` | Documentation-only changes                  | `docs/update-readme`           |

### Description format

- Free-form kebab-case describing the work (e.g., `workflow-linter`, `branch-naming-and-mapping`).
- No numeric prefixes. Priority and ordering are determined by plan file content, not by naming.

## Exec-Plan-to-Branch Mapping Rule

The branch description and the exec-plan filename MUST share the same name:

- `plan/<name>` branch → creates `docs/exec-plan/todo/<name>.md`
- `feat/<name>` or `fix/<name>` branch → expects `docs/exec-plan/todo/<name>.md` (or `done/<name>.md` if already completed)

After execution is complete, the plan file is moved from `todo/` to `done/`.

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

## Migration

Existing plan files with numeric prefixes (`001-workflow-linter.md`) will be renamed during execution to drop the prefix (`workflow-linter.md`).

## Sub-tasks

- [ ] [parallel] Update `AI_WORKFLOW.md` with branch naming convention and exec-plan mapping rules
- [ ] [parallel] Update `CLAUDE.md` to reference the formal rules
- [ ] [parallel] Rename existing exec-plan files to drop numeric prefixes
- [ ] [depends on: all above] Move `docs/issues/branch-naming-rule.md` and `docs/issues/exec-plan-branch-mapping.md` to `docs/issues/done/`
