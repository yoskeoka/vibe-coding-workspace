# workflow-lint CI fixable warnings are still non-blocking

## Summary

`tools/workflow-lint.sh` currently treats all findings as warnings-only even in
`--mode=ci`, so GitHub Actions stays green unless someone opens the job logs and
reads the emitted warning text.

Relevant implementation points:

- [tools/workflow-lint.sh](/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/tools/workflow-lint.sh:6)
  documents `All checks are warnings only (exit 0)`.
- [tools/workflow-lint.sh](/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/tools/workflow-lint.sh:761)
  only prints a reminder when `FIXABLE_WARN_COUNT > 0`.
- [tools/workflow-lint.sh](/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/tools/workflow-lint.sh:768)
  exits with `0` unconditionally.
- [docs/specs/workflow-linter.md](/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/docs/specs/workflow-linter.md:5)
  and [docs/specs/workflow-linter.md](/home/yoske/src/github.com/yoskeoka/vibe-coding-workspace/docs/specs/workflow-linter.md:72)
  explicitly describe the CI contract as non-blocking.

This makes `fixable` warnings easy to miss in PR review because the checks UI
shows success even when the linter found repo-state problems that the workflow
expects to be resolved before push/PR.

## Proposed Solution

- Change the workflow-linter contract so `--mode=ci` exits non-zero when
  `FIXABLE_WARN_COUNT > 0`.
- Keep `--mode=pre-push` non-blocking if we still want local visibility without
  blocking raw Git operations.
- Update `docs/specs/workflow-linter.md`, `AGENTS.md`, and any PR/review skill
  docs that currently describe `fixable` warnings as visible-but-non-blocking.
- Audit child-repo copies or installers that distribute `tools/workflow-lint.sh`
  so the behavior stays consistent across the workspace workflow.

## Priority

High. This is a workflow correctness gap, not just a naming issue. The current
contract says `fixable` warnings should normally be resolved before push/PR, but
CI still reports green unless a reviewer manually opens the logs.
