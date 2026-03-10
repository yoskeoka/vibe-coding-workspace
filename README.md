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

For a new project, run `./skills/manage-workflow/run.sh <project-dir>` to apply this structure.

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
