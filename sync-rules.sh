#!/bin/bash

# sync-rules.sh
# Synchronizes the AI-Centered Development workflow rules and directory structure
# to managed sub-projects.

WORKSPACE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATES_DIR="$WORKSPACE_ROOT/templates"
WORKFLOW_DOC="$WORKSPACE_ROOT/AI_WORKFLOW.md"

# Function to sync a single project
sync_project() {
    local project_dir="$1"
    
    # Check if directory exists
    if [ ! -d "$project_dir" ]; then
        echo "Error: Directory $project_dir does not exist."
        return 1
    fi

    echo "Syncing rules to: $project_dir"

    # 1. Copy the Workflow Documentation and Agent Guidelines
    cp "$WORKFLOW_DOC" "$project_dir/AI_WORKFLOW.md"
    echo "  - Updated AI_WORKFLOW.md"

    # Copy AGENTS.md (if it doesn't exist, or maybe we want to enforce it? 
    # For now, let's copy it as AGENTS_RULES.md and include it in the project's main AGENTS.md if possible,
    # or just overwrite if it's a managed project. The user said: "rules ... manage in this repository".
    # So we'll enforce it.
    cp "$WORKSPACE_ROOT/AGENTS.md" "$project_dir/AGENTS.md"
    echo "  - Updated AGENTS.md"

    # Ensure CLAUDE.md symlink exists
    if [ ! -L "$project_dir/CLAUDE.md" ] && [ ! -f "$project_dir/CLAUDE.md" ]; then
        ln -s AGENTS.md "$project_dir/CLAUDE.md"
        echo "  - Created CLAUDE.md -> AGENTS.md symlink"
    fi

    # 2. Ensure docs/ directory structure exists
    # We use rsync to copy the structure. 
    # -a: archive mode
    # --ignore-existing: Do not overwrite existing files
    # --exclude: exclude specific files if needed
    
    if command -v rsync &> /dev/null; then
        rsync -a --ignore-existing "$TEMPLATES_DIR/docs" "$project_dir/"
        echo "  - Ensured docs/ structure and templates exist"
    else
        # Fallback if rsync is not available (mostly for very minimal environments)
        echo "  - rsync not found, using basic copy (no overwrite)"
        cp -r -n "$TEMPLATES_DIR/docs" "$project_dir/"
    fi
     
    echo "  - Sync complete."
}

# Main execution

# If arguments are provided, sync those specific directories
if [ "$#" -gt 0 ]; then
    for dir in "$@"; do
        sync_project "$dir"
    done
else
    # Otherwise, scan for git repositories in immediate subdirectories
    echo "Scanning for sub-projects..."
    for dir in */ ; do
        # remove trailing slash
        dir=${dir%/}
        
        # Skip if it's the current directory or not a git repo
        if [ "$dir" == "." ] || [ "$dir" == ".." ] || [ "$dir" == "templates" ] || [ "$dir" == "docs" ]; then
            continue
        fi

        # Check if it has a .git folder or is clearly a project
        if [ -d "$dir/.git" ]; then
            sync_project "$dir"
        fi
    done
fi
