#!/bin/bash
set -euo pipefail

base_ref="${1:-origin/main}"

while IFS= read -r path; do
    [ -n "$path" ] || continue
    [ -f "$path" ] || continue

    case "$path" in
        docs/specs/*.md|docs/design-decisions/*.md|docs/kb/*.md|docs/references/*.md)
            printf '%s\n' "$path"
            ;;
    esac
done < <(git diff --name-only --diff-filter=AMR "${base_ref}...HEAD")
