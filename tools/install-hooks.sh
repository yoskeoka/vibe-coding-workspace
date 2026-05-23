#!/bin/bash
set -euo pipefail

# Install workflow hooks in a repository.
#
# Usage:
#   tools/install-hooks.sh [repo-path]
#
# If no repo-path is given, installs hooks in the current directory.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Determine target repo
if [ -n "${1:-}" ]; then
    TARGET_REPO="$(cd "$1" && pwd)"
else
    TARGET_REPO="$(pwd)"
fi

# Verify target is a git repo (regular checkout or git worktree)
if ! git -C "$TARGET_REPO" rev-parse --git-dir >/dev/null 2>&1; then
    echo "Error: $TARGET_REPO is not a git repository." >&2
    exit 1
fi

echo "Installing workflow hooks in: $TARGET_REPO"

# Create directories
mkdir -p "$TARGET_REPO/.githooks"
mkdir -p "$TARGET_REPO/tools"
mkdir -p "$TARGET_REPO/.github/workflows"

# Copy hook and linter
cp "$WORKSPACE_ROOT/.githooks/pre-push" "$TARGET_REPO/.githooks/pre-push"
chmod +x "$TARGET_REPO/.githooks/pre-push"

cp "$WORKSPACE_ROOT/tools/workflow-lint.sh" "$TARGET_REPO/tools/workflow-lint.sh"
chmod +x "$TARGET_REPO/tools/workflow-lint.sh"

cp "$WORKSPACE_ROOT/.github/workflows/workflow-lint.yml" "$TARGET_REPO/.github/workflows/workflow-lint.yml"

# Set core.hooksPath
git -C "$TARGET_REPO" config core.hooksPath .githooks

echo "Workflow runtime assets installed. core.hooksPath set to .githooks"
