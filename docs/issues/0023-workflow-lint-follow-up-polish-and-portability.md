# workflow-lint follow-up polish and portability cleanup

## Summary

The workflow-lint rollout fixed the functional CI visibility gap, but a few
low-priority polish and portability items remain in the shared workflow source.

## Details

- `.github/workflows/workflow-lint.yml`
  - The annotation helper currently escapes `%`, `\r`, and `\n`.
  - GitHub Actions workflow-command property values such as `title=...` should
    also escape `:` and `,` to avoid malformed annotations when those
    characters appear in titles.
- `tools/install-hooks.sh`
  - The script now installs runtime assets beyond git hooks, but its leading log
    line still says `Installing workflow hooks`.
- `tools/workflow-lint.sh`
  - The same-repo issue warning path still says `external GitHub issue` even
    when the URL belongs to the current repository.
  - The same-repo URL detection uses a `grep` pattern with `\+` without
    enabling extended regex, which is less portable than necessary.

## Proposed Solution

Apply a shared follow-up pass in the workspace workflow source that:

1. hardens annotation property escaping in `.github/workflows/workflow-lint.yml`
2. updates `tools/install-hooks.sh` log wording to match its current behavior
3. makes GitHub-issue wording neutral in `tools/workflow-lint.sh`
4. replaces the same-repo URL detection with a portable pattern

Then re-sync the copied workflow assets into child repos through the existing
workflow rollout path.
