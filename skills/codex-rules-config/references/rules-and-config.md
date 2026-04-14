# Codex Rules And Config Reference

This file keeps the minimum official guidance needed for the `codex-rules-config` skill.

## Config locations

- User config lives in `~/.codex/config.toml`.
- Project-scoped overrides can live in `.codex/config.toml`.
- Project config only applies when the project is trusted.

## Rules locations

- Create `.rules` files under `.codex/rules/`.
- User-level example: `~/.codex/rules/default.rules`.
- Project/team-level example: `<repo>/.codex/rules/default.rules`.
- Codex scans `rules/` under team config locations at startup.
- After adding or changing rule files, restart Codex to ensure the new rules are loaded.

## `prefix_rule()` fields

- `pattern` required: command prefix tokens to match
- `decision` optional, defaults to `"allow"`
- decisions:
  - `"allow"`: run outside the sandbox without prompting
  - `"prompt"`: ask before each matching invocation
  - `"forbidden"`: block without prompting
- `justification` optional but recommended
- `match` and `not_match` optional inline tests

Codex applies the most restrictive matching decision:
- `forbidden` > `prompt` > `allow`

## Validation

Use:

```bash
codex execpolicy check --rules /path/to/default.rules --pretty -- <command> ...
```

Notes:
- Use more than one `--rules` flag to test combined files.
- The command emits JSON with the effective decision and matching rules.

## Practical guidance

- Put `prefix_rule(...)` in `.rules`, not in `config.toml`.
- Use project scope for repo-specific workflow allowances.
- Use user scope for personal defaults and global safety rules.
- For forbidden rules, make the justification actionable when possible.
