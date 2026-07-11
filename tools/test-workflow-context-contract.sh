#!/bin/bash
set -euo pipefail

root=$(git rev-parse --show-toplevel)
cd "$root"

required=(
    AGENTS.md
    AI_WORKFLOW.md
    docs/specs/workflow-context-contract.md
    docs/design-decisions/README.md
    docs/lessons.md
)
for file in "${required[@]}"; do
    test -f "$file"
done

test ! -e docs/design-decisions/adr.md
test -d docs/design-decisions/adr
test "$(find docs/design-decisions/adr -maxdepth 1 -name '*.md' -type f | wc -l | tr -d ' ')" -gt 0
grep -q 'AI_WORKFLOW.md#' AGENTS.md
grep -q 'workflow-context-contract.md' AI_WORKFLOW.md

for file in skills/{plan-execution,execute-task,review-task,post-task-review,manage-workflow,plan-project,triage-tasks}/SKILL.md; do
    grep -q 'AI_WORKFLOW.md#' "$file"
done

lint_output=$(./tools/workflow-lint.sh --mode=pre-push 2>&1)
if echo "$lint_output" | grep -qE 'Workflow context contract|ADR .*is not indexed|Monolithic ADR|Active exceptions'; then
    echo "$lint_output" >&2
    echo "workflow context contract warning" >&2
    exit 1
fi

echo "workflow context contract: passed"
