# ShellCheck changed Bash scripts in CI
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Add a CI check that runs ShellCheck only against Bash script files changed by a pull request. This supports the workspace tooling goal in `docs/project-plan.md`: workflow/tooling is a first-class deliverable, and this repo expects to accumulate shell scripts for automation.

## Context

- Current CI includes `workflow-lint.yml`, which enforces AI workflow rules.
- `docs/specs/workflow-linter.md` explicitly treats project-specific lint as out of scope for the workflow linter.
- Past decision in `docs/design-decisions/adr.md`: use focused workflows that only fan out when their inputs matter, such as the child-repo sync workflow checking `skills/` changes before opening PRs.
- Core belief: correctness over speed. Applying ShellCheck to changed Bash scripts gives useful mechanical feedback without making unrelated script history a blocker.

## Code changes

### `.github/workflows/shellcheck.yml`

Create a GitHub Actions workflow that:

- runs on `pull_request` targeting `main`
- checks out the PR head and explicitly fetches the PR base branch before diffing, rather than assuming `origin/main` already exists after checkout
- installs ShellCheck explicitly on the runner instead of relying on runner image contents
- computes changed candidate files with `git diff --name-only --diff-filter=AMR <fetched-base-ref>...HEAD`
- filters to existing Bash script files:
  - files ending in `.sh`, treated as Bash for this workspace's automation scripts
  - files ending in `.bash`
  - files with a Bash shebang, including `#!/bin/bash` and `#!/usr/bin/env bash`
- exits successfully with a clear "no changed Bash scripts" message when no files match
- runs ShellCheck against the matching files and fails CI on ShellCheck findings
- invokes ShellCheck in Bash mode for extension-based matches that do not have their own shell directive or shebang

### `tools/list-changed-bash-scripts.sh`

Add a small helper script so the detection logic is testable locally and not embedded only in YAML.

Expected behavior:

- accepts an optional base ref, defaulting to `origin/main`
- compares `${base_ref}...HEAD`
- emits one matching path per line
- ignores deleted files
- handles filenames safely enough for the current repository convention, which uses ordinary path names without embedded newlines
- uses structured shell logic rather than fragile one-off YAML-only pipelines

### Optional local documentation

If useful during execution, add a short mention to the relevant tooling spec rather than creating operator-facing README churn.

## Spec changes

Add `docs/specs/shellcheck-ci.md` documenting:

- purpose and scope of the ShellCheck CI
- exact definition of a "changed Bash script"
- the rule that `.sh` files are considered Bash automation scripts in this workspace unless the implementation deliberately narrows the policy
- skipped/no-op behavior when a PR changes no Bash scripts
- relationship to `workflow-lint.yml` and why this is a separate project lint workflow
- local verification command for the helper script and ShellCheck invocation

Do not fold this behavior into `docs/specs/workflow-linter.md` except for an optional cross-reference. The workflow linter remains focused on AI workflow rules.

## Design decisions

No ADR update is required unless execution chooses a materially different architecture, such as introducing a generalized changed-file lint framework. The planned decision is local and follows existing ADR guidance: keep automation narrowly scoped to relevant changed inputs.

## Sub-tasks

- [ ] [parallel] Add `docs/specs/shellcheck-ci.md`.
- [ ] [parallel] Add `tools/list-changed-bash-scripts.sh`.
- [ ] [depends on: helper script] Add `.github/workflows/shellcheck.yml`.
- [ ] [depends on: spec, helper script, workflow] Run local verification for the helper script against representative changed files.
- [ ] [depends on: workflow] Run ShellCheck locally on existing Bash scripts when ShellCheck is available, or document the missing local dependency and rely on CI installation for PR verification.

## Verification plan

- Run `tools/list-changed-bash-scripts.sh origin/main` on the execution branch and confirm it lists changed Bash scripts only.
- Create or modify one temporary Bash-script candidate during execution verification if needed, then remove it before commit, to confirm `.sh` and Bash shebang detection.
- Run `shellcheck` against the helper and any changed Bash scripts when available locally.
- Run the existing workflow linter in CI mode or pre-push mode to confirm the new plan and later implementation do not introduce fixable workflow warnings.

## Non-goals

- Do not lint every historical shell script on every PR.
- Do not add a pre-commit hook.
- Do not make ShellCheck part of `tools/workflow-lint.sh`.
- Do not expand this into a generalized multi-language changed-file lint runner in this task.
