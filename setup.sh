#!/bin/bash

# List of repositories to manage
REPOS=(
    "https://github.com/yoskeoka/reversi-adventure"
)

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
        git clone "$repo_url"
    fi

    # Add to .gitignore if directory exists and not already present
    if [ -d "$dir_name" ] && ! grep -q "^$dir_name/$" .gitignore; then
        echo "Adding $dir_name/ to .gitignore..."
        echo "$dir_name/" >> .gitignore
    fi

    echo "-----------------------------------"
done
