## Plan / Issues

<!-- Link the exec-plan, issue, or project-plan that triggered this PR. If none, write N/A.
If the matching execution plan has an `Addresses:` line with local issue paths, list the same
paths under Issues unless you are intentionally leaving one open. In that case, keep the path
here and say `remains open: <reason>` somewhere in the PR body.
If the matching execution plan has an `Addresses:` line with external GitHub issues, list the
same issues under Issues. For implementation PRs that resolve those issues, also fill in the
Closes section below with explicit closing keywords. -->

- **Plan**: <!-- e.g., docs/exec-plan/todo/0042-feature-name.md -->
- **Issues**: <!-- e.g., docs/issues/0019-bug-name.md, or GitHub issue link -->

## Closes

<!--
For Step 3 implementation PRs that resolve external GitHub issues declared in the plan's
`Addresses:` line, add one explicit closing keyword per issue.

Examples:
- Closes #227
- Closes https://github.com/yoskeoka/ww/issues/227

Leave `N/A` when this PR does not close an external GitHub issue.
-->

N/A

## Type of Change

- [ ] Project Plan update
- [ ] Execution Plan (new/updated plan)
- [ ] Feature implementation
- [ ] Bug fix
- [ ] Refactor
- [ ] Documentation only
- [ ] Chore (CI, tooling, deps)

## Change Summary

<!--
Overall change summary. Recommended to use bullet list. Each explanations should be concise.
-->

## Human Instructions / Intent

<!--
Record the human instruction and intent that led to creating this PR.
This section documents how the PR was initiated, not a command that must
match text inside the diff.

If this PR creates or updates an execution plan, record the planning command
that produced the plan PR.
If this PR executes an approved plan, record the execution command that
produced the implementation PR.

Examples:
- Plan creation: `/plan-execution docs/exec-plan/todo/<sequence>-<plan-name>.md`
- Plan execution: `/execute-task docs/exec-plan/todo/<sequence>-<plan-name>.md`
- Project planning: `/plan-project docs/project-plan.md`

For chores/docs-only or other non-execution PRs, use:
- Human instruction: `N/A`
-->
Human instruction: `________________`
### Additional Context from Instructing Human

<!--
Record instructions, decisions, and intent from the human that are NOT already
captured in the exec-plan, specs, or code diff. This section preserves context
that would otherwise be lost when the conversation ends.

What belongs here:
- Library/framework preferences ("Use slog instead of log for structured logging")
- Implementation-style directives ("Keep it simple — no premature abstraction")
- Decisions made during AI-human dialogue — include the AI's question AND the
  human's answer so the context is self-contained
  (e.g., "AI proposed approach A (faster) vs B (simpler). Human chose B,
   reasoning: maintenance cost matters more than runtime speed for this component.")
- Scope decisions about discovered issues
  ("Found stale import in utils.go — out of scope, logged as docs/issues/stale-import.md")
- Quality/priority directives ("Ship this as-is; polish in next iteration")

If no additional instructions were given beyond the standard plan execution, write N/A.
-->

N/A

## Verification

<!-- How was this change verified? Fill in relevant items. -->

- [ ] Tests pass (command: `________________`)
- [ ] Lint passes (command: `________________`)
- [ ] Manual verification (describe below)

## Checklist

- [ ] Branch created from latest `origin/main`
- [ ] `docs/specs/` updated (Spec-Code Parity) — _if code changed_
- [ ] Plan moved from `todo/` to `done/` — _if executing a plan_
- [ ] Resolved linked local issues from the plan's `Addresses:` line were moved to `docs/issues/done/`, or this PR explains why they remain open
- [ ] External GitHub issues declared in the plan's `Addresses:` line are linked under `Issues`, and implementation PRs include matching `Closes` entries or explain why they remain open
- [ ] Workflow-linter warnings reviewed; all `fixable` warnings were resolved or explicitly justified in this PR
- [ ] New issues logged in `docs/issues/` — _if discovered during work_
- [ ] No unresolved blockers remain

## Dependencies

<!-- PRs or issues that must be merged before this one, or that are blocked by this one. -->

N/A

## Reviewer Notes

<!-- Specific areas to focus on during review, known trade-offs, or things that look wrong but are intentional. -->

N/A

## Links

<!-- External references: library docs, design references, related discussions, etc. -->

N/A

## Breaking Changes

<!-- Describe breaking changes. Delete this section if none. -->

N/A

## Screenshots / Logs

<!-- Attach verification artifacts if applicable (test output, screenshots, before/after metrics). Delete this section if not applicable. -->
