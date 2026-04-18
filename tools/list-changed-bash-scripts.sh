#!/bin/bash
set -euo pipefail

base_ref="${1:-origin/main}"

has_bash_shebang() {
    local path="$1"
    local first_line

    if ! IFS= read -r first_line < "$path"; then
        return 1
    fi

    case "$first_line" in
        '#!/bin/bash'|"#!/usr/bin/env bash")
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

is_bash_script() {
    local path="$1"

    [ -f "$path" ] || return 1

    case "$path" in
        *.sh|*.bash)
            return 0
            ;;
    esac

    has_bash_shebang "$path"
}

while IFS= read -r path; do
    [ -n "$path" ] || continue

    if is_bash_script "$path"; then
        printf '%s\n' "$path"
    fi
done < <(git diff --name-only --diff-filter=AMR "${base_ref}...HEAD")