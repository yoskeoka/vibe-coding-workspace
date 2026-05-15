# AI-Centered Development Workflow

This document outlines the workflow for developing projects with AI as the central driver.

## Core Principles

1.  **AI-Centric Context**: All necessary information must be immediately accessible to AI in the `docs/` directory. Files should be structured for easy parsing and context retrieval.
2.  **Spec-Code Parity**: `docs/specs/` must strictly match the actual code. No PR is reviewable without verification that specs and code are in sync.
3.  **Verification First**: Human review happens *after* mechanical tests and "visual" verification data are ready.

## Directory Structure

- `docs/project-plan.md`: The single source of truth for the project's goals, significance, and requirements.
- `docs/design-decisions/`:
    - `adr.md`: Append-only log of architectural decisions.
    - `core-beliefs.md`: Guiding principles for trade-offs.
    - `rejected-ideas.md`: Append-only log for no-go project ideas (with rationale and revisit conditions).
- `docs/specs/`: Detailed specifications. Must match implementation.
- `docs/exec-plan/`:
    - `todo/`: Active execution plans.
    - `done/`: Completed execution plans.
- `docs/references/`: Context for external tools/protocols (e.g., `fastapi-llms.txt`).
- `docs/lessons.md`: Accumulated lessons learned. Reviewed at session start, with new lessons appended at the end of the file.
- `docs/issues/`: Local issue tracking for people developing through this workspace workflow. Avoids confusion with GitHub Issues during active "exec-plan" cycles.
    - When feedback comes from outside this workspace and naturally lives in the target repository's GitHub Issues, keep that GitHub issue as the canonical tracker instead of mirroring it into `docs/issues/` unless separate workspace-only follow-up is needed.
    - `done/`: Resolved issues (moved here after fix is merged).
- `.local/pj/`: Derived local cache for workspace-level GitHub Projects triage. This cache is non-canonical and must stay untracked.

## Workflow Cycle

### 0. New Project Intake (Pre-Step, optional but recommended for vague ideas)
- Use this when an idea is still fuzzy and not ready for a full `project-plan`.
- Activities:
    1. Idea sparring (pain points, desired experience, target users)
    2. Existing-solution research
    3. Go/No-Go decision
- Before asking what to do next after a non-trivial intake checkpoint, compress the current findings into `docs/issues/<descriptive-name>.md` with the preserved problem framing, conclusion, next questions, and source links so later sessions can resume without replaying the full research context.
- If **No-Go**: append findings to `docs/design-decisions/rejected-ideas.md` and stop.
- If **Go**: create/bootstrap the child project repo, update workspace meta entries, then continue to Step 1 (`plan-project`).

> **Rule**: Every step that produces changes MUST go through a GitHub PR review — including doc-only changes like Project Plan and Execution Plan updates. AI Agents must always create a new clean branch from the latest `main` before starting any work.

### Branch Setup (applies to every step below)
- **Default operator path**: use the globally installed `ww` CLI for normal plan/execution startup instead of raw git branch creation.
- From the target repo root:
    ```sh
    ww create <type>/<description>
    cd "$(ww cd <type>/<description>)"
    ```
- From the workspace root when targeting a child repo:
    ```sh
    ww create --repo <repo> <type>/<description>
    cd "$(ww cd --repo <repo> <type>/<description>)"
    ```
- Never reuse an existing feature branch or task worktree silently; each active task should get its own `ww` worktree.
- Raw git branch creation is reserved for `ww` bootstrap/recovery cases and for developing or verifying unreleased `ww` behavior inside `ww/`. If `ww` fails unexpectedly, follow `docs/specs/ww-dogfooding-workflow.md` instead of bypassing it silently.
- When capturing a `ww` failure, record the command, cwd, target repo, expected behavior, actual behavior, relevant output, fallback usage, and impact so the finding can feed back into `ww`.
- This startup contract is fully migrated for the workflow docs and skills covered by `docs/specs/ww-dogfooding-workflow.md`, including adjacent project planning, PR review, and workflow bootstrap skills.

#### Branch Naming Convention

Branch names MUST match the pattern `<type>/<description>`:

| Type   | Purpose                                     | Example                        |
|--------|---------------------------------------------|--------------------------------|
| `plan` | Execution plan creation/update              | `plan/feature-name`            |
| `feat` | Feature implementation (from an exec-plan)  | `feat/feature-name`            |
| `fix`  | Bug fix implementation (from an exec-plan)  | `fix/bug-name`                 |
| `chore`| Non-functional changes (CI, tooling, deps)  | `chore/update-ci`              |
| `docs` | Documentation-only changes                  | `docs/update-readme`           |

The `<description>` is free-form kebab-case. Do not add workflow sequence numbers to branch names; ordering lives in active plan and issue filenames instead.

#### Active Plan / Issue Naming

- Active execution plans under `docs/exec-plan/todo/` MUST use `<sequence>-<name>.md`.
- Active local issues under `docs/issues/` MUST use `<sequence>-<name>.md`.
- Sequence numbers use four-digit zero padding from `0001` through `9999`.
- Sequence numbers at `10000` or above MUST be written without zero padding.
- `README.md` is exempt in both directories.
- New active files should take the next available sequence number in their directory family so creation order remains visible without opening Git history.

#### Exec-Plan Mapping

The branch description and the exec-plan basename suffix MUST share the same name:

- `plan/<name>` branch creates `docs/exec-plan/todo/<sequence>-<name>.md`
- `feat/<name>` or `fix/<name>` branch expects one matching active or completed plan whose basename suffix is `-<name>.md`
- Historical completed plans in `docs/exec-plan/done/` may still use older non-numbered filenames and remain valid
- After execution is complete, the plan file is moved from `todo/` to `done/`
- Branches of type `chore` and `docs` are exempt (no exec-plan required)

### PR Workflow (applies to every step below)
1. **Verify** — Run **all** project lint and test commands using non-AI tooling (e.g., `make lint`, `npm run lint`, `go vet`, `pytest`, `npm test`, or whatever the project defines). If any check fails, fix the issue in the same branch and re-run until **all pass**. Skip this for doc-only PRs when no lint/test tooling covers documentation.
2. **Create PR and follow up** — Push the branch and create or update the PR through the shared `review-task` gate when using project workflow skills. Use the **PR template** and fill in all sections. Template priority: **current repo template > workspace root repo template > vendored workflow template > workflow-repo template** — if the current repo has `.github/PULL_REQUEST_TEMPLATE.md`, use that; otherwise, when working in a child repo inside this workspace, use the workspace root repo template at `<workspace-root>/.github/PULL_REQUEST_TEMPLATE.md`; otherwise, in child repos that vendor this workflow, use `.claude/vendor/workflow/.github/PULL_REQUEST_TEMPLATE.md`; otherwise use the workflow repo's `.github/PULL_REQUEST_TEMPLATE.md`. After PR creation or update, `review-task` owns the bounded initial monitoring loop for the current PR head SHA: wait for CI/checks to settle, inspect for configured advisory bot/agent review activity, and collect review signals before handoff. A workflow step that routes into `review-task` is not complete until that gate reaches a documented stop condition for the latest pushed PR head SHA.
3. **Review** — Wait for GitHub PR review approval before merging into `main`.

#### Post-PR Follow-up

PR creation is not the terminal workflow action. For every new PR, updated PR, or later push to the PR branch, required CI/check inspection restarts for the new head SHA:

- Planning PRs and execution PRs use the same completion boundary: the caller must not report the task complete before `review-task` reaches a documented stop condition for the latest pushed head SHA.
- CI failures are mechanical verification failures. Investigate failing checks and fix them in-branch when the logs are actionable and the fix stays within scope; then re-run verification, push, and restart monitoring for the new head SHA.
- Advisory bot/agent review comments from Copilot, Claude, `gh aw`, agent workflow reviews, or other configured automation are human-review input. Do not edit files, apply suggestions, commit, or push in response to those findings unless the human explicitly approves that specific follow-up or a prior human instruction already authorized implementing that exact review feedback.
- Handoffs must group substantive advisory findings by source reviewer/workflow and include the source, location/link, extracted summary, implementer's view, 1-2 line explanation, and recommendation: fix in this PR, defer, or no action.
- Passing or approving advisory bot/agent checks can still contain substantive observations in review bodies, so inspect review summaries and inline comments even when the overall state is not blocking.
- After PR creation or a later push, wait 30 seconds before the first CI/check and timeline inspection. Only spend the longer advisory-review wait budget when the timeline shows bot/agent review-start activity for the latest head SHA, such as `copilot_work_started`, or when the human explicitly asks to wait.
- For a non-blocked PR creation or update flow, the shared minimum landing path starts as: `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll`, with checks, timeline events, review summaries, and inline comments inspected from that poll.
- Planning PRs may stop after that initial poll when no other stop-condition work remains. Execution PRs use a bounded CI-settling window in the landing check: if required checks are still pending after the initial poll, add up to two more `30-second wait -> poll` turns before handoff, unless checks finish earlier, advisory-review waiting starts, the head SHA changes, the helper fails, or the human asks to stop.
- When advisory bot/agent review has started, use 3 polling turns: wait 5 minutes once, then 1 minute twice, for a 7-minute total budget. After each wait, fetch PR reviews and inline comments. If no review/comments arrive by the end, stop waiting and document the timeout state.
- Advisory review triage is a session handoff, not a PR comment. After implementation, summarize comment response options in the current session unless the user explicitly asks to post them back to the PR.
- Polling-style wait loops should use a low-cost subagent only when the platform supports delegation and the current session explicitly authorizes subagent use; final decisions, fixes, and handoff stay with the main agent.

---

### 1. Project Plan (`docs/project-plan.md`) — **requires PR**
- Create a new `ww` worktree/branch (e.g., `ww create plan/project-plan-v1`).
- Define or update the Goal, Significance, and Requirements.
- Follow the **PR Workflow** above to merge the plan into `main`.
- Update this as the project evolves (each update = new branch + PR).

### 2. Execution Plan (`docs/exec-plan/todo/`) — **requires PR**
- Create a new `ww` worktree/branch (e.g., `ww create plan/initial-setup`).
- Create a new numbered plan file (e.g., `0007-initial-setup.md`) in `todo/`.
- Detail:
    - Code changes.
    - Spec changes (How `docs/specs/` will change).
    - `Addresses:` entries for any tracked issues that this execution work is expected to resolve:
      - local workspace issues under `docs/issues/`
      - external GitHub issues when the canonical feedback lives in the target repo issue tracker
    - Use full GitHub issue URLs in `Addresses:` for external issues, for example `Addresses: https://github.com/yoskeoka/ww/issues/227`.
    - Break down large tasks into smaller sub-plans if needed.
- Review/Update `design-decisions/` if architectural choices are made.
- Follow the **PR Workflow** above to merge the plan into `main`.
- If the plan is driven by an external GitHub issue, the plan PR should link that issue in its PR body under `Issues` so reviewers can trace the execution target before implementation starts.

### 3. Execution — **requires PR**
- Create a new `ww` worktree/branch (e.g., `ww create feat/initial-setup`).
- **Spec First**: Update `docs/specs/` *before* modifying code.
- **Implement**: Write the code to match the spec.
- **Issues**: If unrelated problems are found, log them in `docs/issues/<sequence>-<name>.md`. Do not fix them within the current plan unless blocking.
- **Issue Resolution**: When an issue is resolved, move its file from `docs/issues/` to `docs/issues/done/`. If the matching execution plan declares that issue in `Addresses:`, the implementation branch should include the move unless the PR body explicitly explains why the issue remains open.
- **External GitHub Issue Resolution**: When the matching execution plan declares external GitHub issues in `Addresses:`, the implementation PR must include corresponding closing keywords unless the PR body explicitly explains why the issue remains open.
  - For same-repo issues, use `Closes #<number>` such as `Closes #227`.
  - For cross-repo issues, use the full URL such as `Closes https://github.com/yoskeoka/ww/issues/227`.
- **Completion**: Move the plan file from `docs/exec-plan/todo/` to `docs/exec-plan/done/`.
- Follow the **PR Workflow** above (Verify → Create PR → Review).
- The PR must include:
    - Code changes.
    - Spec updates.
    - The plan file moved to `done/`.
    - Any resolved linked local issue files moved to `docs/issues/done/`, or an explicit PR-body justification for leaving them open.
    - `Closes` entries for external GitHub issues declared in the plan's `Addresses:` line, or an explicit PR-body justification for leaving them open.
    - Verification artifacts (test results, screenshots, logs) for human review.
    - Post-PR follow-up status for the latest pushed head SHA.

Repeat steps 1–3 until the Project Plan is complete.

## Lessons Maintenance

- `docs/lessons.md` is append-only in practice for new lessons: add new entries at the end of the file rather than inserting them near the top.
- When a new rule is learned during planning, execution, or review follow-up, update `docs/lessons.md` in the same branch that captured the lesson.

## Workspace Task Tracking

- GitHub Projects is the canonical remote state for workspace task triage.
- `.local/pj/` is the only supported local workspace-triage state in the current workflow.
- Committed legacy tracker runtime artifacts and local database state are not part of the supported workflow contract.
