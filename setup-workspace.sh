#!/bin/bash
set -euo pipefail

# Setup AI workflow skills and hooks in a child repository.
#
# Usage:
#   Initial setup (from inside a child repository):
#     /path/to/vibe-coding-workspace/setup-workspace.sh
#
#   Initial setup (specify child repo path):
#     /path/to/vibe-coding-workspace/setup-workspace.sh /path/to/child-repo
#
# Note: Submodule updates are handled automatically via GitHub Actions.
#       See .github/workflows/sync-workflow-to-child-repos.yml in vibe-coding-workspace.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKFLOW_REPO_HTTPS="https://github.com/yoskeoka/vibe-coding-workspace.git"
WORKFLOW_REPO_SSH="git@github.com:yoskeoka/vibe-coding-workspace.git"
VENDOR_DIR=".claude/vendor/workflow"

# Skills to symlink from the workflow repo
SKILLS=(
    manage-workflow
    plan-project
    plan-execution
    execute-task
    post-task-review
    review-task
)

# ----------------------------------------
# Parse arguments
# ----------------------------------------
CHILD_REPO=""

while [ $# -gt 0 ]; do
    case "$1" in
        --help|-h)
            echo "Usage: setup-workspace.sh [child-repo-path]"
            echo ""
            echo "Options:"
            echo "  --help      Show this help message"
            echo ""
            echo "If no child-repo-path is given, the current directory is used."
            exit 0
            ;;
        *)
            CHILD_REPO="$1"
            shift
            ;;
    esac
done

# Determine child repo directory
if [ -n "$CHILD_REPO" ]; then
    CHILD_REPO="$(cd "$CHILD_REPO" && pwd)"
else
    CHILD_REPO="$(pwd)"
fi

# Verify we're in a git repo
if [ ! -d "$CHILD_REPO/.git" ]; then
    echo "Error: $CHILD_REPO is not a git repository."
    exit 1
fi

# Don't run on the workflow repo itself
if [ "$CHILD_REPO" = "$SCRIPT_DIR" ]; then
    echo "Error: Cannot run setup-workspace on the workflow repo itself."
    exit 1
fi

cd "$CHILD_REPO"

echo "Setting up AI workflow skills in: $CHILD_REPO"
echo "---"

# ----------------------------------------
# Step 1: Add submodule
# ----------------------------------------
if [ -f ".gitmodules" ] && grep -q "$VENDOR_DIR" .gitmodules 2>/dev/null; then
    echo "Submodule already exists at $VENDOR_DIR. Skipping."
    git submodule update --init "$VENDOR_DIR"
else
    echo "Adding workflow submodule (shallow clone)..."
    if git submodule add --depth 1 "$WORKFLOW_REPO_HTTPS" "$VENDOR_DIR" 2>/dev/null; then
        echo "  Added via HTTPS."
    else
        echo "  HTTPS failed. Trying SSH..."
        git submodule add --depth 1 "$WORKFLOW_REPO_SSH" "$VENDOR_DIR"
        echo "  Added via SSH."
    fi
fi

echo "Submodule ready at $VENDOR_DIR"
echo "---"

# ----------------------------------------
# Step 2: Create skill directories
# ----------------------------------------
# Skills are symlinked into three directories so that Claude Code, Gemini,
# and other agents that read .agents/skills/ all discover them.
SKILL_DIRS=(
    ".claude/skills"
    ".agents/skills"
    ".gemini/skills"
)

for dir in "${SKILL_DIRS[@]}"; do
    mkdir -p "$dir"
done

# ----------------------------------------
# Step 3: Create symlinks for each skill
# ----------------------------------------
for skill in "${SKILLS[@]}"; do
    for dir in "${SKILL_DIRS[@]}"; do
        link="$dir/$skill"

        # Compute relative path from the skill dir back to the vendor skills.
        # .claude/skills/ → ../vendor/workflow/skills/<skill>
        # .agents/skills/ → ../../.claude/vendor/workflow/skills/<skill>
        # .gemini/skills/ → ../../.claude/vendor/workflow/skills/<skill>
        if [ "$dir" = ".claude/skills" ]; then
            target="../vendor/workflow/skills/$skill"
        else
            target="../../.claude/vendor/workflow/skills/$skill"
        fi

        if [ -L "$link" ]; then
            echo "Symlink already exists: $dir/$skill"
        elif [ -e "$link" ]; then
            echo "Warning: $link exists but is not a symlink. Skipping."
        else
            ln -s "$target" "$link"
            echo "Linked: $dir/$skill → $target"
        fi
    done
done

echo "---"

# ----------------------------------------
# Step 4: Copy docs/ templates if not exist
# ----------------------------------------
TEMPLATE_DIR="$VENDOR_DIR/skills/manage-workflow/templates/docs"

if [ -d "$TEMPLATE_DIR" ] && [ ! -f "docs/project-plan.md" ]; then
    echo "Initializing docs/ from templates..."
    cp -rn "$TEMPLATE_DIR" docs/ 2>/dev/null || cp -r "$TEMPLATE_DIR" docs/
    echo "docs/ structure created."
else
    echo "docs/ already exists or templates not found. Skipping."
fi

echo "---"

# ----------------------------------------
# Step 5: Ensure CLAUDE.md symlinks to AGENTS.md
# ----------------------------------------
if [ -L "CLAUDE.md" ]; then
    echo "CLAUDE.md is already a symlink. OK."
elif [ -f "CLAUDE.md" ] && [ -f "AGENTS.md" ]; then
    echo "Warning: Both CLAUDE.md and AGENTS.md exist as regular files."
    echo "  Please merge CLAUDE.md content into AGENTS.md, then replace CLAUDE.md with a symlink."
    echo "  Ask the AI agent: 'Sync CLAUDE.md and AGENTS.md with a symlink'"
elif [ -f "CLAUDE.md" ] && [ ! -f "AGENTS.md" ]; then
    echo "Migrating: CLAUDE.md → AGENTS.md (with symlink back)"
    mv CLAUDE.md AGENTS.md
    ln -s AGENTS.md CLAUDE.md
    echo "  AGENTS.md is now the canonical file. CLAUDE.md → AGENTS.md"
elif [ ! -f "CLAUDE.md" ] && [ -f "AGENTS.md" ]; then
    echo "Creating symlink: CLAUDE.md → AGENTS.md"
    ln -s AGENTS.md CLAUDE.md
elif [ ! -f "CLAUDE.md" ] && [ ! -f "AGENTS.md" ]; then
    echo "Note: Neither CLAUDE.md nor AGENTS.md exist yet."
    echo "  The 'manage-workflow' skill will create AGENTS.md and symlink CLAUDE.md."
    echo "  Ask the AI agent: 'Initialize the AI workflow for this project'"
fi

# ----------------------------------------
# Step 6: Install workflow hooks
# ----------------------------------------
echo "Installing workflow hooks..."
"$SCRIPT_DIR/tools/install-hooks.sh" "$CHILD_REPO"

echo "---"
echo "Done! AI workflow skills and hooks are ready."
echo ""
echo "Next steps:"
echo "  1. Ask the AI agent to initialize the workflow (triggers manage-workflow skill)."
echo "     This will set up AGENTS.md with proper workflow instructions."
echo "     CLAUDE.md will be a symlink to AGENTS.md (single source of truth)."
echo "  2. Edit docs/project-plan.md with your project goals."
echo "  3. Add project-specific skills to .claude/skills/ as needed."
echo "  4. Commit the changes: git add -A && git commit -m 'Add AI workflow skills'"
