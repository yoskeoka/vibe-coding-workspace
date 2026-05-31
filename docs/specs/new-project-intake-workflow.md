# Spec: New Project Intake Workflow Package

## Goal
Provide a reusable pre-`plan-project` workflow package for adding new projects to this meta-repo.

## Scope
This package covers only the steps before detailed project planning:
- Idea sparring and rough articulation of pain points and desired experience
- Lightweight market/reference research
- Go/No-Go decision support
- Research-checkpoint compression into durable local handoff notes
- New project bootstrap tasks in the workspace
- Handoff proposal to `plan-project`

## Requirements

### 1. New skill for pre-planning phase
A dedicated skill MUST exist to run this flow before `plan-project`.

Skill responsibilities:
1. Act as a sparring partner to clarify:
   - what problem should be solved
   - what experience should be delivered
2. Research existing solutions/experiences (e.g., games, products, tools).
3. Evaluate whether there is meaningful value to build:
   - novelty or differentiation
   - personal motivation / hobby value
   - reason users would care
4. Branch on decision:
   - GO: proceed to project creation/bootstrap
   - NO-GO: summarize findings and append to rejected-ideas log
5. For GO case, execute project bootstrap:
   - create GitHub repository if missing
   - ensure local child repo exists
   - initialize docs/workflow scaffold in child repo
   - ensure workspace meta config includes the new project (`setup.sh`, `.gitignore`, `README.md` Managed Projects, `docs/project-plan.md` Managed Projects, `AGENTS.md` workspace structure, and `.github/workflows/sync-workflow-to-child-repos.yml` when the repo should receive workflow sync PRs)
6. After bootstrap, explicitly propose continuing with `plan-project` in the new child repo.

### 2. Research-checkpoint handoff notes
When `new-project-intake` performs non-trivial idea sparring or external research, it MUST compress the current state into a durable `docs/issues/<descriptive-name>.md` note before asking the user what to do next.

This note exists to discard no-longer-needed exploration context while preserving the parts that matter for later sessions.

The note MUST capture:
- the current problem framing or pain-point hypothesis
- the current conclusion for this checkpoint
- the leading solution options still worth considering
- rejected or deprioritized options only when their rejection matters to future decisions
- the evidence needed to resume later, with source links for every external reference
- the next concrete research questions or decision points

The note SHOULD stay compressed and decision-oriented rather than turning into a raw research log.

The note MUST be created or updated before the agent asks the user a checkpoint question such as whether to continue research, compare another option, pivot, bootstrap, or stop.

When later intake work supersedes the earlier checkpoint, the same issue note SHOULD be updated instead of creating fragmented duplicate notes unless the topic has clearly split into separate tracks.

### 3. Rejected idea log
Workspace MUST contain an English-named file to store no-go ideas and research outcomes:
- `docs/design-decisions/rejected-ideas.md`

The file should be append-only and include:
- date
- idea name
- summary
- why no-go now
- conditions to revisit
- references

Before appending the final no-go decision, the intake flow SHOULD already have a `docs/issues/` checkpoint note that preserves the fuller research context and references.

### 4. Skill distribution boundary
The new skill is **meta-repo only** and MUST NOT be distributed to child repos via `setup-workspace.sh`.

Rationale:
- This skill is for deciding whether to create a project and bootstrapping workspace-level assets.
- Child repos should start from `plan-project` and not carry pre-creation intake logic.

### 5. Documentation updates
Workflow docs SHOULD mention this pre-`plan-project` package and when to use it.

## Non-Goals
- Full product requirements definition (handled by `plan-project`)
- Feature-level implementation planning (handled by `plan-execution`)
- Implementation work (handled by `execute-task`)
