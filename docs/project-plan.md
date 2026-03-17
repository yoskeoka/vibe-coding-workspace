# Project Plan

## Goal

A personal vibe-coding workspace where AI agents do the heavy lifting to build games, tools, and apps that I want — while I have fun directing and reviewing. The workspace itself evolves as a living lab for AI-centered development practices.

### What I want from this
1. **Have fun vibe-coding** — the process should be enjoyable, not bureaucratic
2. **Build things I want** — games, tools, apps; not limited to any category
3. **Accumulate AI-centered development knowledge** — learn what works and what doesn't through real hobby projects
4. **Keep costs low** — this is a hobby, not a full-time job; minimize token spend and human time per output

## Significance

- **Workflow as product**: The AI-Centered Development workflow (AI_WORKFLOW.md, skills, tooling) is a first-class deliverable. Child projects are both the point and the proving ground.
- **One human, many projects**: Structured plans, specs, and mechanical enforcement let a single reviewer manage multiple projects in parallel without burning out.
- **Reusable patterns**: What works here can be carried to other repos or shared publicly.

## Requirements

### Workflow framework
- [ ] AI_WORKFLOW.md defines the complete development lifecycle (project plan, exec-plan, execution, PR review)
- [x] CLAUDE.md / AGENTS.md configures AI agent behavior to follow the workflow
- [x] docs/ directory structure supports AI context retrieval (specs, plans, issues, decisions)
- [x] Branch naming convention is formally declared and enforceable
- [x] Exec-plan-to-branch mapping convention is formally declared

### Tooling
- [x] Workflow linter enforces declared rules mechanically (pre-push hook + CI)
- [x] setup.sh clones and updates child project repos
- [ ] Shared tooling (hooks, linter) is distributable to child repos
- [ ] Triage skill aggregates status across all managed projects

### Child project management
- [x] Child projects are listed in setup.sh and gitignored
- [x] Each child project has its own docs/ structure and project-plan.md
- [ ] New project intake workflow is defined and tested (new-project-intake skill)

## Managed Projects

| Project | Description | Status |
|---------|-------------|--------|
| [reversi-adventure](https://github.com/yoskeoka/reversi-adventure) | Reversi game with AI opponent and move explanation | Phase 1 complete, Phase 2 next |
| [ai-arena](https://github.com/yoskeoka/ai-arena) | AI competition platform with multiple game types | Phase 1 in planning |
| [vim-learning-game](https://github.com/yoskeoka/vim-learning-game) | Gamified Vim learning experience | Greenfield, plan defined |

## Milestones

- [x] Phase 1: Bootstrap workspace structure (docs/, AI_WORKFLOW.md, CLAUDE.md, setup.sh)
- [x] Phase 2: Formalize workflow rules (branch naming, exec-plan mapping conventions)
- [x] Phase 3: Implement workflow linter (pre-push hook + CI checks)
- [ ] Phase 4: Distribute shared tooling to child repos (hooks, linter, workspace setup)
- [ ] Phase 5: Iterate on workflow based on child project feedback
