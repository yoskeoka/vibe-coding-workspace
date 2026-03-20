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
    ├─ 1. Find Missing User Intent
    │     └─ Identify any user intent that was not fulfilled during the task
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

### 1. Find Missing User Intent

Identify any user intent that was not fulfilled during the task. This could be:

- System Requirements: Missing required constraints, performance targets, or compatibility needs
- Reasons of Decision: Unstated rationale behind a choice that would be important for future work

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
  - Use "Step 1. Find Missing User Intent" result and refer to `docs/design-decisions/adr.md` and `core-beliefs.md` for existing design knowledge
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
