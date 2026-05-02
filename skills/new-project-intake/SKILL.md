---
name: new-project-intake
description: When evaluating a brand-new project idea before plan-project, doing idea sparring, researching existing solutions, making a go/no-go call, bootstrapping a new project repo/workflow, and handing off to plan-project.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# New Project Intake (Pre-Step to Workflow Step 1)

**Position in workflow**: This is a **pre-step** before `plan-project` (Step 1). Use this skill when the user has a raw idea but project goals/requirements are not yet ready.

## What to Do

Run a structured intake flow:

1. **Idea sparring**
   - Help the user verbalize:
     - unresolved pain points
     - desired user experience
     - target user and usage context
   - Keep this lightweight; do not force detailed requirements yet.

2. **Existing-solution research**
   - Investigate current alternatives (products, games, tools, communities).
   - Summarize what experience already exists and what is missing.
   - Before asking the user what to do next after a meaningful research or sparring checkpoint, compress the current findings into `docs/issues/<descriptive-name>.md`.
   - Use that note as the durable next-session handoff so web-search results, discarded comparisons, and transient technical-detail context can be dropped from chat context safely.

3. **Value test (Go/No-Go)**
   - Judge whether the idea has build value using:
     - novelty/differentiation
     - personal motivation (hobby/learning value)
     - user-facing payoff
   - Make the decision explicit:
     - **GO**: proceed to bootstrap.
     - **NO-GO**: log and stop.

4. **NO-GO path**
   - Append summary to `docs/design-decisions/rejected-ideas.md` with:
     - date
     - idea name
     - short summary
     - why no-go now
     - conditions to revisit
     - references
   - End with a concise recommendation (pause, pivot, or re-scope).

5. **GO path: project bootstrap**
   - Create GitHub repo if missing.
   - Ensure local project directory exists under this workspace.
   - Initialize docs/workflow scaffold (use `setup-workspace.sh` or template structure).
   - Update meta-repo management entries:
     - `setup.sh` (`REPOS`)
     - `.gitignore`
     - `README.md` Managed Projects

6. **Handoff to Step 1**
   - After bootstrap, explicitly ask:
     - "Move to the new project directory and continue with `plan-project` now?"

7. **Checkpoint handoff note**
   - For any non-trivial research checkpoint, create or update `docs/issues/<descriptive-name>.md` before asking the user whether to continue, pivot, compare another option, bootstrap, or stop.
   - Keep the note compressed and decision-oriented.
   - Include:
     - date
     - idea / topic name
     - current pain point or problem framing
     - checkpoint conclusion
     - live options still worth considering
     - rejected/deprioritized options only when the rejection matters later
     - key technical facts worth preserving
     - source links for every external reference
     - next questions / next candidate actions
   - Update the existing note when the same line of inquiry continues. Create a separate note only when the work clearly splits into independent tracks.

## Rules

1. Do not skip research and go/no-go framing when the idea is still vague.
2. Keep pre-step output concise and decision-oriented.
3. Do not produce a full project plan in this step; hand off to `plan-project`.
4. If user asks only for bootstrap work (already decided GO), skip straight to Step 5.
5. Do not ask the user a checkpoint "what next?" question after substantial research until the compressed `docs/issues/` handoff note is updated.
6. Keep durable notes in English unless the user explicitly asks for a different artifact language.

## Deliverables

- GO case:
  - `docs/issues/<name>.md` checkpoint handoff note with preserved research context and source links
  - New repo and local project skeleton ready
  - Meta-repo references updated
  - Proposal to continue with `plan-project`
- NO-GO case:
  - `docs/issues/<name>.md` checkpoint handoff note with preserved research context and source links
  - Entry appended to `docs/design-decisions/rejected-ideas.md`
  - Research-backed summary of why it is not worth starting now
