# Lessons Learned

## New Lessons Append At The End

- **Mistake**: I inserted new workflow lessons near the top of `docs/lessons.md`, which made recency and review order harder to trust across sessions.
- **Pattern**: When a guidance file is read top-down for context, it is easy to optimize for visibility in the current session and accidentally turn the file into a manually re-sorted document.
- **Rule**: Add every new `docs/lessons.md` entry at the end of the file. Do not insert new lessons above older ones unless the task is explicitly reorganizing the file.
- **Applied**: `docs/lessons.md` maintenance across plan creation, execution, post-task review, and workflow-doc updates.

## Video-Backed Notes Need Whole-Video Context

- **Mistake**: I normalized a video-backed source note around only the article-highlighted segment and did not first capture the broader embedded video structure.
- **Pattern**: When a wrapper article points to one useful moment, it is easy to treat that moment as the entire durable source and miss surrounding model, workflow, pricing, benchmark, or prompting context that future searches need.
- **Rule**: For `source_type: video` or `video_backed_article`, first skim or caption-review the full video, then decide which anchors are durable; preserve the article-specific anchor plus any broader segments needed to make the source useful without reprocessing the video.
- **Applied**: `docs/kb/sources/`, especially video-backed article maintenance and follow-up normalization tasks.

## Command Names Should Avoid Repeating the Same Noun

- **Mistake**: I proposed `pj link-repo --repo <owner>/<repo>`, which repeated `repo` in both the command and flag and made the UX feel clumsy.
- **Pattern**: When extending a small CLI, copying a flag-oriented shape from another command can produce redundant names instead of matching the command tree's natural noun/action structure.
- **Rule**: For new `pj` command groups, sketch the command tree first and choose either `noun action <target>` or `verb <target>` so the target noun is not repeated in both the command and primary argument flag.
- **Applied**: `tools/pj/internal/pj/app.go` command additions and `docs/specs/github-projects-task-cli.md` CLI UX specs.

## Issue-to-Implementation Requires Planning First

- **Mistake**: I treated a non-trivial `docs/issues/` item as ready for `execute-task` even though no execution plan existed yet.
- **Pattern**: An issue file can describe a useful fix, but it is not the same artifact as `docs/exec-plan/todo/<sequence>-<name>.md`; jumping straight to execution skips the workflow's planning review gate.
- **Rule**: When the user points at a non-trivial `docs/issues/<sequence>-<name>.md` item and no matching `docs/exec-plan/todo/<sequence>-<name>.md` exists, start with `plan-execution`, create `docs/exec-plan/todo/<sequence>-<name>.md`, and only use `execute-task` after that plan PR is merged.
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

## Workflow Tasks Must Land the Initial PR Follow-up Loop

- **Mistake**: I treated a workflow implementation as complete after local edits, verification, or lesson updates, and stopped before the branch lifecycle reached `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll`.
- **Pattern**: Once the local change is correct, it is easy to narrate intended PR follow-up instead of treating `review-task` completion for the latest PR head SHA as part of the same task.
- **Rule**: When a task routes into `review-task`, do not report completion before the latest pushed PR head SHA reaches a documented `review-task` stop condition; for non-blocked flows, that minimum landing path is `commit -> push -> PR create/update -> 30-second wait -> initial follow-up poll` with checks, timeline, review summaries, and inline comments inspected.
- **Applied**: `AI_WORKFLOW.md`, `docs/specs/pr-follow-up-workflow.md`, `skills/plan-execution/SKILL.md`, `skills/execute-task/SKILL.md`, and any session handoff that claims a planning or execution task is complete.

## Portable Docs Should Not Hard-Code Local Absolute Paths

- **Mistake**: I updated workflow docs with `/home/yoske/...` absolute paths that matched this machine but were not portable to other workspace locations.
- **Pattern**: When fixing a repo-local workflow gap, it is easy to optimize for the current checkout path and miss that the documentation is a reusable contract for other environments.
- **Rule**: In durable docs and skills, describe workspace-relative paths with placeholders such as `<workspace-root>/...` unless the task explicitly requires a machine-specific path example.
- **Applied**: Workflow docs, PR template fallback guidance, and any skill text that explains file locations across repos or workspaces.

## Parallel Git Writes Need Worktree-Level Serialization

- **Mistake**: I ran `git add` and `git commit` against the same worktree in parallel and treated the resulting `index.lock` failure like a surprising worktree issue.
- **Pattern**: In a multi-tool or multi-agent workflow, it is easy to parallelize Git commands that look independent even though they both mutate the same worktree index and refs.
- **Rule**: Within a single worktree, serialize Git commands that write repository state such as `git add`, `git commit`, `git rebase`, `git merge`, `git cherry-pick`, `git checkout`, and similar operations. Restrict parallelism to read-only inspection, verification, or clearly separate worktrees.
- **Applied**: `AGENTS.md`, `execute-task` / `review-task` style execution flows, and any Codex session that uses parallel tool calls or subagents around Git operations.
