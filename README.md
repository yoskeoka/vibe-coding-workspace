# vibe-coding-workspace
My vibe coding workspace to maximize AI-Centered Development

## Meta-Repo Management

This repository acts as a meta-repository to manage other AI projects.

### Configuration
1. Edit [setup.sh](setup.sh) to add repository URLs to the `REPOS` array.
2. Run `./setup.sh` to clone or update all project repositories.
3. The script automatically adds cloned repositories to `.gitignore` to prevent nested tracking.

### Managed Projects
- [reversi-adventure](https://github.com/yoskeoka/reversi-adventure) (AI Agent config & skills workspace)
- [ai-arena](https://github.com/yoskeoka/ai-arena) (AI-vs-AI online game platform)
- [vim-learning-game](https://github.com/yoskeoka/vim-learning-game) (Learn Vim through gameplay)
- [ww](https://github.com/yoskeoka/ww) (Workspace worktree manager for multi-repo parallel development)

## AI-Centered Development Workflow

This workspace is managed using a strict **AI-first workflow**. The core idea is that all context needed for development (plans, specs, decisions) is stored in a structured way that AI agents can easily parse and reason about.

**Key Documents:**
- [AI_WORKFLOW.md](AI_WORKFLOW.md): The "Standard Operating Procedure" for this project. Read this first.
- [docs/project-plan.md](docs/project-plan.md): The high-level goals and roadmap.
- [docs/specs/](docs/specs/): Detailed specifications. Code must always match these specs.
- [AGENTS.md](AGENTS.md): Instructions for the AI agent itself.

**Workflow Summary:**
1.  **Goal**: Define or update the high-level roadmap in `docs/project-plan.md`.
2.  **Plan**: Create a specific execution plan in `docs/exec-plan/todo/` to achieve a goal.
3.  **Spec**: Update `docs/specs/` to reflect the change *before* coding.
4.  **Code**: Implement the change.
5.  **Verify**: Ensure specs and code match, then move the plan to `done/`.

For normal day-to-day task startup, use the globally installed `ww` CLI instead of raw `git switch -c` so the workflow continuously dogfoods the released tool:

```bash
ww create plan/example-task
cd "$(ww cd plan/example-task)"
```

From the workspace root when targeting a child repo:

```bash
ww create --repo ww feat/example-task
cd "$(ww cd --repo ww feat/example-task)"
```

Reserve repo-local `ww` development builds for work inside `ww/` itself or for reproducing/verifying a `ww` bug.

If global `ww` fails during normal workflow startup, treat that as a first-class workflow finding. Capture the command, cwd, target repo, expected behavior, actual behavior, relevant output, any raw-git fallback, and the impact on the blocked task so it can be filed back to `ww`.

For a new project, run `./skills/manage-workflow/run.sh <project-dir>` to apply this structure.

## Workspace Task Triage CLI

The workspace task triage spike lives under `tools/pj/` and uses `gh auth token` for GitHub Projects API access.

Before running `pj init`, make sure `gh` is authenticated with the required scopes:

```bash
gh auth login -h github.com -s project,read:project,repo,read:org,gist
```

If you are already logged in and only need to add the missing GitHub Projects write scope:

```bash
gh auth refresh -h github.com -s project
```

Verify the current auth state with:

```bash
gh auth status
```

Then bootstrap the canonical workspace board:

```bash
go -C tools/pj run ./cmd/pj init --owner <owner> --owner-type user|org
```

Notes:
- `read:project` is enough for read-only sync, but `pj init`, `pj add`, and `pj move` require the `project` scope.
- `pj init` can create the `Workspace Task Triage` board, but the current spike still expects the custom `Repo`, `Kind`, and `Priority` fields to be provisioned separately.

## Using Skills in Child Projects

To set up the AI workflow skills in a child repository:

```bash
# From inside the child repo
/path/to/vibe-coding-workspace/setup-workspace.sh

# Or specify the child repo path
/path/to/vibe-coding-workspace/setup-workspace.sh /path/to/child-repo
```

This will:
1. Add this repo as a shallow Git submodule at `.claude/vendor/workflow/`
2. Create symlinks in `.claude/skills/` for each workflow skill
3. Copy `docs/` templates, `CLAUDE.md`, and `AGENTS.md` if they don't exist
4. Install workflow hooks (pre-push linter)

Child repos can add project-specific skills directly to `.claude/skills/`.

To update skills later:
```bash
git submodule update --remote --depth 1 .claude/vendor/workflow
```
