---
name: review-task
description: When creating a pull request, preparing a PR for review, generating verification artifacts, collecting test results or screenshots for review, submitting changes for human review, or checking that all PR requirements (code, specs, plan, verification) are met.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Review Task (Workflow Step 4)

**Position in workflow**: This is **Step 4** (final step) of the AI-Centered Development cycle. Execution (Step 3) must be complete and the plan must be in `docs/exec-plan/done/` before creating a PR. After review, repeat from Step 2 for the next task.

## What to Do

Create a Pull Request that includes all required artifacts for human review.

### PR Must Include

1. **Code changes**: The implementation from Step 3.
2. **Spec updates**: The updated `docs/specs/` files that match the code.
3. **Plan file moved to `done/`**: The execution plan in `docs/exec-plan/done/` proving the task was completed through the proper workflow.
4. **Verification artifacts**: Test results, screenshots, logs, or other evidence for human reviewers. Human review happens _after_ mechanical tests and verification data are ready.

### Verification Standards by Task Type

| Task Type   | Minimum Verification             |
| ----------- | -------------------------------- |
| Bug fix     | Reproduce → Fix → Verify fixed   |
| Feature     | Tests pass + manual demo         |
| Refactor    | Behavior unchanged + tests pass  |
| Performance | Before/after metrics             |
| Security    | Specific vulnerability addressed |

### Pre-PR Checklist

Before creating the PR, verify:

- [ ] `docs/specs/` matches the implementation (Spec-Code Parity).
- [ ] Plan file has been moved from `docs/exec-plan/todo/` to `docs/exec-plan/done/`.
- [ ] Tests pass and test results are available.
- [ ] Any visual or behavioral changes have screenshots/logs attached.
- [ ] No unresolved blockers remain (non-blockers should be in `docs/issues/`).

### Verification First Principle

Human review happens **after** mechanical tests and "visual" verification data are ready. The PR should contain enough evidence that a reviewer can confirm correctness without running the code themselves.

## Next Step

After the PR is reviewed and merged, return to **Execution Plan** (Step 2) to pick up the next task. Repeat Steps 2–4 until the Project Plan is complete.
