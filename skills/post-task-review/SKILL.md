---
name: post-task-review
description: After completing significant work (bug fix, feature, investigation), review findings, log issues in docs/issues/, update lessons learned, and propose CLAUDE.md updates. Should be performed as part of task completion, not only when explicitly invoked.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Post-Task Review (Workflow Step 3.5)

**Position in workflow**: This step runs between **Execution** (Step 3) and **Review/PR** (Step 4). After moving a plan from `todo/` to `done/` and before creating a PR, perform this self-retrospective.

## When to Use

- After completing a bug fix, feature, or investigation
- When the user says work is done, asks to wrap up, or asks for a review
- After touching multiple files and gaining codebase insight
- **Automatically** after moving an exec-plan from `todo/` to `done/`

Do NOT use for trivial changes (typo fixes, single-line edits).

## Workflow

```
Task completed
    │
    ├─ 1. Capture Unrecorded User Intent
    │     └─ Capture user intent/rationale not yet persisted to project memory
    │
    ├─ 2. Review Findings
    │     └─ Present prioritized summary to user
    │
    ├─ 3. Log issues in docs/issues/
    │     └─ Optionally mirror to GitHub Issues (with user approval)
    │
    ├─ 4. Update docs/issues/lessons.md
    │     └─ Document patterns from corrections encountered during task
    │
    └─ 5. Propose CLAUDE.md / AGENTS.md updates
          └─ Apply with user approval
```

### 1. Capture Unrecorded User Intent

Review the session for user knowledge that was expressed or implied but NOT yet persisted to project memory (`CLAUDE.md`, `AGENTS.md`, `docs/design-decisions/`, `docs/specs/`).

**What to look for:**

- **Context injections**: The user provided background, motivation, or "why" that guided the work — is it captured in project docs?
- **Unexplained choices**: The user selected an option (e.g., "skip this", "use approach A") without stating the reasoning. Ask: "You chose X over Y — what was your reasoning?" and persist the answer.
- **Corrected assumptions**: The user corrected the agent's understanding of goals, scope, or priorities — is the corrected understanding now reflected in project memory?
- **Implicit goals**: Objectives or constraints that the user "just knows" but that don't appear anywhere in `docs/` or `CLAUDE.md`. These surfaced naturally during conversation but would be lost at session end.

**Process:**

1. Scan the session for moments where the user injected context, made a choice, or corrected direction.
2. For each, check: is this knowledge already in project memory?
3. If not, ask the user to verbalize the intent/rationale (if they haven't already).
4. Collect these items — they feed into Step 5 (propose where to persist them).

> The goal is **knowledge capture**, not task-completion tracking. "Did we finish everything?" belongs in the exec-plan status, not here.

### 2. Review Findings

Identify issues discovered during work. Categories to check:

- **Spec-code parity gaps**: Public APIs listed, Input/Output mismatches and behavior inconsistencies between `docs/specs/` and code
- **Duplicated logic**: Same business logic in multiple files
- **Inconsistent patterns**: Different approaches to the same problem across files
- **Missing tests**: Untested critical paths found during investigation
- **Tight coupling**: Components that should be separated
- **Dependency concerns**: Version mismatches, deprecated APIs

Present a prioritized summary to the user. Ask which items to log.

### 3. Log Issues in docs/issues/

Create `docs/issues/<descriptive-name>.md` for each approved finding. Each file must include:

- **Summary**: What the problem is, with specific file paths and line numbers
- **Proposed Solution**: Concrete direction, not vague suggestions
- **Priority**: Why it matters (data integrity, performance, maintainability)

**Important**: `docs/issues/` is the AI's primary memory for issue tracking. Always create files here first. Optionally ask the user if they also want GitHub Issues created via `gh issue create`.

### 4. Update Lessons Learned

Check if corrections occurred during the task. If so, create or update `docs/issues/lessons.md` using this format:

- **Mistake**: What went wrong (be specific)
- **Pattern**: The underlying cause or anti-pattern
- **Rule**: Concrete, actionable rule to prevent recurrence
- **Applied**: Where this rule applies (specific files, patterns, situations)

> "Be more careful" is not a rule. Rules must be specific and testable.

### 5. Propose CLAUDE.md / AGENTS.md, Skills, Design Decisions Updates

Check if the work revealed knowledge that would reduce future investigation time:

- **Project-specific insights**: Non-obvious information about the codebase, architecture, or design decisions
  - Use "Step 1. Capture Unrecorded User Intent" result and refer to `docs/design-decisions/adr.md` and `docs/design-decisions/core-beliefs.md` for existing design knowledge
- **Build/test commands**: New crates, test targets, or lint configurations
- **Architecture notes**: How subsystems connect, data flow, key design decisions
- **Duplication risks**: List of files that must be updated together
- **Tech stack changes**: New dependencies, version requirements

Think where is the best place to document these knowledge items and propose specific additions to the user before editing.

- Propose creating skills for recurring patterns or named workflows that require multiple steps to execute
  - To find patterns, look session logs and `docs/issues/lessons.md` for repeated sequences of actions that could be abstracted into a skill
- Update `CLAUDE.md` for project specific knowledge and best practices that AI agents should follow on every task
  - Update both `CLAUDE.md` and `AGENTS.md` to keep them in sync.
- Propose to update User's Global `CLAUDE.md` and `AGENTS.md` if the insight is user's core belief or fundamental principle that's not related to the specific project but can guide design decisions in general.
- Otherwise add to AI agents'Memory or just dismiss based on user's choice.

## What NOT to Do

- Do not create issues without user approval
- Do not add speculative or hypothetical issues
- Do not update CLAUDE.md with information already documented
- Do not add generic best practices — only project-specific knowledge discovered during the task
- Do not use `gh issue create` without first creating the corresponding `docs/issues/` file
