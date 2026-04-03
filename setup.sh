#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ---- Prerequisites ----

# Ensure beads (bd) is installed
if ! command -v bd &>/dev/null; then
    if command -v brew &>/dev/null; then
        echo "Installing beads..."
        brew install beads
    else
        echo "Error: beads (bd) is not installed and Homebrew is not available."
        echo "Please install Homebrew first (https://brew.sh) or install beads manually:"
        echo "  https://github.com/steveyegge/beads"
        exit 1
    fi
fi

# ---- Child repositories ----

# List of repositories to manage (use HTTPS URLs)
REPOS=(
    "https://github.com/yoskeoka/reversi-adventure"
    "https://github.com/yoskeoka/ai-arena"
    "https://github.com/yoskeoka/vim-learning-game"
    "https://github.com/yoskeoka/ww"
    "https://github.com/yoskeoka/homebrew-ww"
)

# Convert an HTTPS GitHub URL to SSH format
# e.g. https://github.com/owner/repo -> git@github.com:owner/repo.git
convert_to_ssh() {
    local url="$1"
    echo "$url" | sed -E 's|https://github.com/([^/]+)/(.+?)(.git)?$|git@github.com:\1/\2.git|'
}

# Try to clone a repository using multiple methods, in order:
#   1. gh CLI
#   2. git + HTTPS URL
#   3. git + SSH URL (converted from HTTPS)
# Returns 0 on success, 1 on failure.
try_clone() {
    local url="$1"
    local dir="$2"

    if command -v gh &>/dev/null; then
        echo "  [1/3] Trying: gh repo clone ..."
        gh repo clone "$url" "$dir" && return 0
    else
        echo "  [1/3] Skipped: gh not available."
    fi

    echo "  [2/3] Trying: git clone (HTTPS) ..."
    git clone "$url" "$dir" && return 0

    local ssh_url
    ssh_url=$(convert_to_ssh "$url")
    echo "  [3/3] Trying: git clone (SSH: $ssh_url) ..."
    git clone "$ssh_url" "$dir" && return 0

    echo "  Error: All clone methods failed for $url"
    return 1
}

# Loop through the list
for repo_url in "${REPOS[@]}"; do
    # Skip empty lines if any
    [[ -z "$repo_url" ]] && continue

    # Extract the directory name from the URL
    dir_name=$(basename "$repo_url" .git)

    if [ -d "$dir_name" ]; then
        echo "Updating $dir_name..."
        cd "$dir_name" || continue

        # Check if it's a git repository (check for .git dir)
        if [ -d ".git" ]; then
            git pull
        else
            echo "Warning: Directory $dir_name exists but is not a git repository."
        fi

        cd ..
    else
        echo "Cloning $dir_name..."
        try_clone "$repo_url" "$dir_name"
    fi

    # Add to .gitignore if directory exists and not already present
    if [ -d "$dir_name" ] && ! grep -q "^$dir_name/$" .gitignore; then
        echo "Adding $dir_name/ to .gitignore..."
        echo "$dir_name/" >> .gitignore
    fi

    # Run setup-workspace.sh to sync AI workflow skills and submodule
    if [ -d "$dir_name/.git" ]; then
        echo "Running setup-workspace.sh for $dir_name..."
        if ! "$SCRIPT_DIR/setup-workspace.sh" "$dir_name"; then
            echo "Error: setup-workspace.sh failed for $dir_name" >&2
            exit 1
        fi
    fi

    echo "-----------------------------------"
done
