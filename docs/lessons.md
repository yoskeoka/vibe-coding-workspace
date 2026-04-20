# Lessons Learned

## Issue-to-Implementation Requires Planning First

- **Mistake**: I treated a non-trivial `docs/issues/` item as ready for `execute-task` even though no execution plan existed yet.
- **Pattern**: An issue file can describe a useful fix, but it is not the same artifact as `docs/exec-plan/todo/<name>.md`; jumping straight to execution skips the workflow's planning review gate.
- **Rule**: When the user points at a non-trivial `docs/issues/<name>.md` item and no matching `docs/exec-plan/todo/<name>.md` exists, start with `plan-execution`, create `docs/exec-plan/todo/<name>.md`, and only use `execute-task` after that plan PR is merged.
- **Applied**: Workflow issue follow-ups, especially issues that change specs, skills, scripts, or multiple files.

## Skipped Detecting target project

- **Mistake**: I started investigating from the workspace root instead of the target child repository.
- **Pattern**: In this meta-repo, child projects live under the workspace root and task scope must be resolved before running repo-specific workflow or tests.
- **Rule**: When a request is about a child project, switch into that child repo first and perform branching, docs, specs, and verification there.
- **Applied**: All tasks in the `vibe-coding-workspace` meta-repo, especially requests involving `ww/`, `ai-arena/`, `reversi-adventure/`, or `vim-learning-game/`.

## Lost Shell State Between Workflow Steps

- **Mistake**: I referenced `CURRENT_SHA` in a later GitHub Actions step without passing it through `env` or `outputs`, so the variable expanded to an empty string and `git diff` failed with `bad revision ''`.
- **Pattern**: Shell variables only live for the duration of a single `run` block; step boundaries drop local state unless it is explicitly exported.
- **Rule**: When a later workflow step needs a value from an earlier step, store it in `GITHUB_OUTPUT` and pass it back via `steps.<id>.outputs.*` or re-compute it in the later step.
- **Applied**: GitHub Actions workflows, especially multi-step jobs that build PR bodies, diff summaries, or other derived metadata.

## Summaries Dropped Concrete Retrieval Anchors

- **Mistake**: I summarized a source into generic categories and dropped the concrete service names that made the source useful for later search and comparison.
- **Pattern**: Over-compressing external references can preserve the high-level recommendation while destroying the practical lookup value of the note.
- **Rule**: When ingesting references that compare concrete services, tools, libraries, APIs, or documents, keep those proper names in both the source note and the relevant compiled wiki page unless there is a strong reason not to.
- **Applied**: `docs/kb/sources/`, `docs/kb/wiki/`, and the `knowledge-base` skill's ingest behavior.

## Dogfooding Tasks Must Stay on the Released Global `ww`

- **Mistake**: I could have treated a workspace workflow task as if it should use raw git or an in-repo `ww` build instead of the released global `ww` binary.
- **Pattern**: When the repo contains the tool being dogfooded, it is easy to conflate "work on the workflow" with "work inside the tool repo" and silently switch operator paths.
- **Rule**: For workspace workflow tasks, use the globally installed `ww` CLI as the default operator path and avoid touching the `ww/` repo unless the task explicitly targets unreleased `ww` behavior.
- **Applied**: Workspace-level planning and execution tasks, especially dogfooding changes that mention `ww` but do not modify files under `ww/`.

## Python Cache Files Must Be Ignored Before Staging

- **Mistake**: I let `skills/knowledge-base/scripts/__pycache__/` get staged after running Python verification commands.
- **Pattern**: When verification creates interpreter cache files in a new script directory, they can slip into the index if the repo does not already ignore them.
- **Rule**: When adding Python files or running Python verification in a repo, ensure `.gitignore` covers `__pycache__/` and `*.py[cod]` before staging, and remove any generated cache files from the worktree before `git add`.
- **Applied**: Any workspace or child-repo task that adds Python code or runs Python commands that emit bytecode caches.

## Environment Constraints Must Not Become Product Constraints

- **Mistake**: I documented the video pipeline as if `CPU/WSL2` were a product-level assumption, even though that was only true for the current execution environment.
- **Pattern**: Temporary local constraints can leak into skills and specs if I optimize for the machine in front of me instead of defining environment-aware behavior.
- **Rule**: When the user gives machine-specific constraints, reflect them as runtime detection or tuning logic unless the product intentionally targets only that environment.
- **Applied**: Skill docs, CLI defaults, dependency guidance, and any implementation that chooses performance-sensitive options such as GPU usage, frame cadence, or batch sizing.

## PR Creation Must Check Bot Review Timeline

- **Mistake**: I created a PR, checked CI and `reviewDecision`, but did not check the GitHub timeline for a pending Copilot review before calling the PR ready.
- **Pattern**: `gh pr view` can show `reviewRequests: []` before or after Copilot review activity, while the UI still shows Copilot review-requested / review-started timeline events. Bot reviews can later leave actionable inline comments without changing `reviewDecision` to a blocking state.
- **Rule**: After creating or updating a PR, wait 30 seconds, then inspect CI/checks and the issue timeline. Only use the longer Copilot wait when the timeline shows bot review-start activity such as `copilot_work_started`; then wait 5 minutes, 1 minute, and 1 minute, checking PR reviews and inline comments after each wait. Present Copilot comment response options in the current session; do not post triage back to the PR unless explicitly asked.
- **Applied**: All PR handoffs, especially after `review-task` or `execute-task` creates a new PR.

## Knowledge-Base Ingest Is File-Changing Workflow Work

- **Mistake**: I started a knowledge-base ingest directly on `main` even though ingest predictably modifies `docs/kb/` and often adds follow-up notes.
- **Pattern**: Research-oriented tasks can look read-only at the start, but KB ingest is normally a documentation change that must follow the same branch and PR discipline as other workflow work.
- **Rule**: Before a file-changing KB ingest, check the working tree, ensure local changes are preserved or absent, create a fresh branch/worktree from latest `origin/main` with global `ww`, and only then write KB files.
- **Applied**: `knowledge-base` skill usage, especially URL/video ingest requests that will create or update source notes, wiki pages, logs, issues, or renderer/tooling files.

## Bot/Agent Review Comments Need Human Triage Before Edits

- **Mistake**: I applied an advisory bot review suggestion during PR follow-up before first extracting the comments, adding the implementer's view, and asking whether the item should be fixed in the current PR.
- **Pattern**: Treating advisory bot review like an already-approved implementation instruction turns PR follow-up into unreviewed scope changes; treating an approving/pass review state as "no comments to inspect" can also hide minor observations and notes in the review body.
- **Rule**: When Copilot, `gh aw`, agent workflow, or another advisory bot leaves review comments or review-body observations, first present each item grouped by source with the implementer's view, a 1-2 line explanation, and a recommendation of whether it should be fixed in the current PR; do not edit or push changes for those comments unless the human explicitly approves.
- **Applied**: `review-task` follow-up, `execute-task` PR monitoring, and any workflow that reads bot/agent PR review comments.

## Token-Saving Helpers Should Fail Closed

- **Mistake**: I documented raw `gh` API reads as an automatic fallback when the compact PR follow-up helper is missing or fails.
- **Pattern**: A helper created to avoid large token usage can lose its value if its failure path silently returns to raw, high-volume JSON polling.
- **Rule**: When a workflow helper exists specifically to limit context/token cost, its missing/failing path should report the failure reason and stop automatic monitoring unless the user explicitly asks for raw inspection or targeted helper diagnosis.
- **Applied**: `skills/review-task/scripts/gh-pr-followup`, `skills/review-task/SKILL.md`, and `docs/specs/pr-follow-up-workflow.md`.

## Plan Details Must Match Declared CI Contract

- **Mistake**: I left the execution-plan workflow trigger and tool-install guidance more permissive than the CI contract documented earlier in the same plan.
- **Pattern**: Drafting implementation details later in a plan can drift from earlier scope constraints unless they are rechecked for exact parity.
- **Rule**: When a plan defines CI trigger scope or deterministic tool behavior, restate the same constraints in the workflow section (explicit `paths` filters and pinned tool versions) instead of leaving open-ended wording like unscoped triggers or `@latest`.
- **Applied**: Execution-plan docs that introduce or modify GitHub Actions workflows, especially `docs/exec-plan/todo/pj-ci-checks.md`.
