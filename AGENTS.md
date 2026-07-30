# AI Agent Entry Point

`AGENTS.md` is the always-read, universal workflow guard. Read the task-specific
document named below rather than loading every workflow document.

## Non-negotiable guards

- Work on a fresh `ww create <type>/<name>` worktree from current `main`; enter it
  with `ww cd`. Do not silently reuse a task branch. Record unexpected `ww`
  behavior under [the dogfooding contract](docs/specs/ww-dogfooding-workflow.md).
- Non-trivial work needs a reviewed execution plan. Update the black-box
  `docs/specs/` contract before code. Keep unrelated findings in `docs/issues/`.
- Every change, including docs-only work, goes through a PR. Run applicable
  non-AI quality gates, resolve fixable workflow-linter findings, push, and use
  `review-task` for the bounded latest-head follow-up before handoff.
- An execution branch deletes its resolved plan and linked local issues after
  verification and PR preparation; retrieve them through PR/Git history. Linked
  external GitHub issues need matching `Closes` metadata unless the PR explains
  why they remain open.
- Serialize Git writes in one worktree. Read design decisions before making a
  new architectural choice.

## Workspace routing

This is a meta-repo. Child repositories live below this root (`ww/`,
`ai-arena/`, and others); a GitHub `yoskeoka/<repo>` path resolves to
`<workspace-root>/<repo>/`. Use that child repository's guidance after routing.

## Minimum read by task

| Task | Read now | Then use |
| --- | --- | --- |
| New project or project scope | [AI workflow: project planning](AI_WORKFLOW.md#project-planning) | `plan-project` or `new-project-intake` |
| New feature, bug, or non-trivial docs/tooling work | [AI workflow: execution planning](AI_WORKFLOW.md#execution-planning) | `plan-execution` |
| Approved plan implementation | [AI workflow: execution](AI_WORKFLOW.md#execution) and the plan | `execute-task` |
| PR creation, update, or CI/review follow-up | [AI workflow: PR and follow-up](AI_WORKFLOW.md#pr-and-follow-up) | `review-task` |
| Workflow bootstrap | [AI workflow: workflow setup](AI_WORKFLOW.md#workflow-setup) | `manage-workflow` |
| Session priority review | [AI workflow: task triage](AI_WORKFLOW.md#task-triage) | `triage-tasks` |
| Significant completed work | [AI workflow: post-task review](AI_WORKFLOW.md#post-task-review) | `post-task-review` |

## On-demand references

- [Workflow context contract](docs/specs/workflow-context-contract.md) explains
  ownership and read boundaries.
- [Design decisions index](docs/design-decisions/README.md) and
  [core beliefs](docs/design-decisions/core-beliefs.md) inform design choices.
- [Specification index](docs/specs/README.md) routes product and operated-service
  contracts. `docs/project-plan.md` is the project north star.
