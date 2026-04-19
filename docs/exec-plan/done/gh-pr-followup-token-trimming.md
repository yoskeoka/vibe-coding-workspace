# GH PR Follow-up Token Trimming

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Reduce context churn during `review-task` post-PR follow-up by replacing raw, large GitHub JSON inspection with a compact wrapper that returns only CI, review, timeline, and inline-comment fields that agents need for handoff decisions.

This supports the project goal of keeping hobby workflow costs low while preserving the existing post-PR correctness contract.

Addresses: `docs/issues/gh-pr-followup-token-trimming.md`

## Current State

- `docs/specs/pr-follow-up-workflow.md` requires `review-task` to inspect PR checks, issue timeline events, review summaries, and inline PR comments.
- `skills/review-task/SKILL.md` currently documents raw `gh pr view` and `gh api` calls for these reads.
- Raw issue timeline and pull-comment responses can include large nested issue, repository, diff hunk, user, reaction, and link objects.
- Repeated polling can return the same already-seen events and comments, consuming context without improving the handoff.

## Spec Changes

### `docs/specs/pr-follow-up-workflow.md`

- Add a compact polling helper contract for post-PR follow-up.
- State that `review-task` should prefer the helper when present and fall back to raw `gh` commands only when the helper is missing or fails.
- Define local non-canonical polling state under `.local/gh-pr-followup/`.
- Require state to be keyed by owner, repo, and PR number.
- Require state to record at least:
  - `head_sha`
  - `last_timeline_event_id`
  - `last_review_comment_id`
  - `last_checked_at`
- Require marker reset when the PR head SHA changes.
- Require compact output suitable for direct agent context, including:
  - current head SHA and review decision
  - compact check rollup fields
  - only relevant timeline events
  - only new inline review comments since the stored marker

## Code Changes

### `skills/review-task/scripts/gh-pr-followup`

Create a POSIX shell or Bash wrapper invoked as:

```sh
skills/review-task/scripts/gh-pr-followup poll <owner/repo> <pr-number>
```

The wrapper should:

- call `gh pr view` for `headRefOid`, `statusCheckRollup`, and `reviewDecision`
- call `gh api repos/<owner>/<repo>/issues/<pr-number>/timeline --paginate` with `--jq` to emit only AI-relevant timeline fields
- call `gh api repos/<owner>/<repo>/pulls/<pr-number>/comments --paginate` with `--jq` to emit only AI-relevant inline comment fields
- filter timeline events and inline comments whose IDs are older than or equal to the saved markers
- reset saved markers when `head_sha` changes
- update `.local/gh-pr-followup/<safe-owner-repo-pr>.json` after a successful poll
- emit compact JSON that can be pasted directly into agent context

The compact timeline event shape should include only fields needed to detect review-start/review-complete signals, such as:

- `id`
- `event`
- `created_at`
- `actor`
- `reviewer`
- `team`
- `app`
- `commit_id`
- `review_state`

The compact inline comment shape should include only:

- `id`
- `path`
- `line`
- `user`
- `body`
- `created_at`
- `updated_at`
- `commit_id`
- `html_url`

The compact check shape should include only:

- `name`
- `workflow`
- `status`
- `conclusion`
- `details_url`

### `skills/review-task/SKILL.md`

- Update the Post-PR Follow-up Loop to use the wrapper as the preferred polling path.
- Keep the existing raw `gh` commands as fallback guidance for missing/failing helper cases.
- Clarify that repeated polling should inspect new helper output first, not repeatedly paste raw timeline/comment JSON into the main context.

### `docs/issues/`

- Move `docs/issues/gh-pr-followup-token-trimming.md` to `docs/issues/done/gh-pr-followup-token-trimming.md` during execution after the helper and docs updates are complete.

## Design Decisions

Past decisions reviewed before planning:

- `core-beliefs.md` favors AI-first context retrieval and correctness over speed. This supports compact, repeatable context extraction without weakening PR follow-up checks.
- The 2026-04-14 ADR says normal workflow operations should dogfood the released global `ww` binary. This plan keeps workflow startup on `ww` and does not touch the child `ww/` repo.
- `review-task-post-pr-follow-up.md` established `review-task` as the owner of PR readiness and initial post-PR monitoring.
- `review-task-pr-workflow-guard.md` strengthened `review-task` as the shared PR gate instead of adding overlapping PR workflow skills.

Apply the same reasoning here:

- Keep `review-task` as the post-PR follow-up owner.
- Add a small helper under the `review-task` skill boundary instead of introducing a separate skill.
- Preserve raw `gh` fallback so the workflow remains usable if the helper is unavailable.

No ADR update is expected unless execution reveals that skill-local helper scripts should become a broader workflow tooling convention.

## Sub-tasks

- [x] [parallel] Update `docs/specs/pr-follow-up-workflow.md` with the compact polling helper contract, state behavior, fallback rule, and output requirements.
- [x] [parallel] Design the `gh-pr-followup poll` shell interface and local state filename/keying scheme.
- [x] [depends on: helper design] Implement `skills/review-task/scripts/gh-pr-followup` with compact `gh --jq` reads and marker updates.
- [x] [depends on: helper implementation] Update `skills/review-task/SKILL.md` to prefer the helper in the post-PR follow-up loop.
- [x] [depends on: helper implementation] Add focused verification for script syntax and, if practical, a fake-`gh` smoke test for marker reset/new-comment filtering.
- [x] [depends on: docs and helper] Move `docs/issues/gh-pr-followup-token-trimming.md` to `docs/issues/done/`.
- [x] [depends on: all above] Run workflow lint and relevant shell/script checks, then prepare the execution PR through `review-task`.

## Parallelism

- Spec update and helper interface design can proceed independently.
- Implementation depends on the helper interface.
- Skill wording depends on the final helper behavior.
- Issue resolution and PR verification depend on all implementation and docs updates.

## Verification

- `tools/workflow-lint.sh --mode=pre-push`
- shell syntax check for `skills/review-task/scripts/gh-pr-followup`
- if the script is written to support a `GH_BIN` override, run a fake-`gh` smoke test that verifies:
  - head SHA changes reset markers
  - repeated polls suppress already-seen timeline events and review comments
  - output contains compact fields and omits bulky fields such as `diff_hunk`
- manual spec-code parity check:
  - `docs/specs/pr-follow-up-workflow.md` matches the helper behavior
  - `skills/review-task/SKILL.md` names the helper and fallback behavior accurately

## Expected Outcome

- Agents still complete the bounded post-PR follow-up loop.
- The main context receives compact, decision-oriented PR follow-up data.
- Repeated CI/Copilot polling stops re-consuming old timeline and comment entries.
- Raw `gh` commands remain available as a fallback, but they are no longer the default documented path for `review-task` follow-up.
