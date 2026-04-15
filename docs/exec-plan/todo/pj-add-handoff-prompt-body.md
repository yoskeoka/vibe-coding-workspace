# pj add Handoff Prompt Body

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Make workspace triage items created through `triage-tasks` more actionable by
including a concise task-start handoff prompt in the Project item body when the
skill uses `pj add`.

This supports the workspace goal that triage should help a human or agent move
from backlog review into the next workflow step with minimal re-discovery. When
someone opens the GitHub Project item directly, they should be able to see not
just the source reference, but also the recommended starting context and branch
setup for the task.

## Background

The current `triage-tasks` flow already rebuilds the Project board from
workspace-local sources such as:

- `docs/exec-plan/todo/`
- `docs/issues/`
- open GitHub PRs
- open GitHub Issues

When it adds items back into the Project, the body currently carries only short
source notes. That is enough for cache reconstruction, but it is weak as a
remote-facing handoff: a human looking at the Project in GitHub still has to
reconstruct where to start, which skill to use next, and which files to read.

## Spec Changes

### `docs/specs/triage-tasks.md`

- Require `triage-tasks` to include a compact handoff block in the item body
  when creating Project items with `pj add`.
- Define the minimum body structure:
  - source reference
  - target repo
  - recommended next-step skill when clear
  - suggested `ww create` / `ww cd` startup command
  - initial files or docs to read
  - short goal statement
- Clarify that the body should stay concise enough for GitHub Project scanning,
  not a full execution plan.
- Clarify fallback behavior for items where a concrete startup prompt is not yet
  meaningful, such as generic PR follow-up or externally tracked GitHub Issues.

### `skills/triage-tasks/SKILL.md`

- Update the full re-triage guidance so every `pj add` call emits the new body
  format rather than only a one-line source note.
- Define how the skill should tailor the body by item type:
  - exec-plan item
  - local issue follow-up
  - open PR review/follow-up
  - open GitHub Issue
- Keep the same-language rule for user-facing handoff, and decide whether the
  stored Project body should always use a stable workspace language or mirror
  the current session language.

## Code Changes

### `skills/triage-tasks/SKILL.md`

- Add explicit rules for composing Project item bodies during full re-triage.
- Add body templates that encode the minimum useful startup prompt without
  bloating the board.
- Clarify when the skill should prefer:
  - `plan-execution`
  - `execute-task`
  - no explicit skill recommendation

### Optional supporting docs

- If needed, update `AGENTS.md` or `docs/specs/github-projects-task-cli.md`
  only to keep terminology aligned with the new item-body expectation.
- Do not expand this plan into `tools/pj` CLI changes unless implementation
  reveals the current `pj add --body` path cannot carry the needed text.

## Design Decisions

Past decisions:
- `docs/specs/triage-tasks.md` already defines the default handoff as a
  fresh-session prompt after task selection.
- The workspace already decided to dogfood global `ww` for normal startup.

Apply the same reasoning here: Project items should preserve enough of that
handoff contract to be useful even when someone starts from the remote board
instead of the current chat session.

Open design choice to settle during execution:
- Should Project item bodies always use a stable canonical language for storage,
  or should they mirror the language of the session that created them? The
  current skill prefers same-language output for chat handoff, but stored board
  text may benefit from consistency.

## Sub-tasks

- [ ] Update `docs/specs/triage-tasks.md` to require the richer `pj add` body
      structure and define its minimum fields
- [ ] Update `skills/triage-tasks/SKILL.md` so full re-triage emits the richer
      body format for newly created Project items
- [ ] [parallel] Decide and document the storage-language rule for Project item
      bodies
- [ ] [parallel] Decide the minimal body variants for exec plans, local issues,
      PRs, and GitHub Issues
- [ ] [depends on: spec and skill updates] Re-run `triage-tasks` against the
      workspace Project and inspect a sample of created item bodies in the cache
      and/or GitHub UI
- [ ] [depends on: verification] Confirm the richer body still stays concise
      enough for Project scanning and does not make board maintenance noisy

## Verification

- Confirm the updated spec and skill agree on the required body structure
- Confirm a re-triaged exec-plan item includes enough information for a human
  to start work directly from the Project body
- Confirm at least one PR-backed item and one issue-backed item receive an
  appropriately scoped body rather than a misleading execution prompt
- Confirm `pj list` / `pj sync` behavior remains unchanged apart from the stored
  body text

## Expected Outcome

- Newly created triage items in GitHub Projects become self-contained enough to
  act as a lightweight remote handoff
- Humans checking the Project board in GitHub can decide and start work with
  less local rediscovery
- The workspace preserves a consistent `ww`-based startup contract across chat
  handoff and Project-backed triage
