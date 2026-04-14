---
name: codex-rules-config
description: Update and validate Codex configuration and execpolicy rules for local, project, and team setup. Use when Codex needs to edit `~/.codex/config.toml`, `.codex/config.toml`, `~/.codex/rules/*.rules`, or `.codex/rules/*.rules`; add or review `prefix_rule(...)` entries; format `.rules` files consistently; or test command decisions with `codex execpolicy check`.
---

# Codex Rules Config

## Overview

Use this skill to change Codex configuration and command-approval rules safely.

Treat `config.toml` and `.rules` as separate surfaces:
- use `config.toml` for settings such as sandbox mode, approvals, model, and features
- use `.rules` files for `prefix_rule(...)` command policies

## Quick Start

1. Identify the target scope:
   - user scope: `~/.codex/`
   - project/team scope: `<repo>/.codex/`
2. Decide whether the change belongs in:
   - `config.toml`
   - `rules/*.rules`
   - both
3. Edit the file in the narrowest scope that matches the request.
4. Keep `.rules` entries in a consistent multi-line style.
5. Validate rule behavior with `codex execpolicy check`.
6. Tell the user when a Codex restart is needed.

## File Placement

- Put persistent user-specific settings in `~/.codex/config.toml`.
- Put project-scoped defaults in `<repo>/.codex/config.toml`.
- Put user-scoped command rules in `~/.codex/rules/default.rules`.
- Put project/team command rules in `<repo>/.codex/rules/default.rules` unless the repo already uses a more specific rules file layout.
- Do not put `prefix_rule(...)` inside `config.toml`.

Read `references/rules-and-config.md` when you need the official Codex docs details for file placement, field semantics, or restart behavior.

## Editing Workflow

### 1. Decide the right surface

Use `config.toml` when the request is about:
- sandbox defaults
- approval behavior
- model selection
- feature flags
- MCP or provider configuration

Use `.rules` when the request is about:
- allowing a command prefix without prompting
- prompting for a command prefix
- forbidding a command prefix
- testing which command policy applies

### 2. Prefer the narrowest scope

- Put personal safety or convenience rules in `~/.codex/`.
- Put repo-specific workflow rules in `<repo>/.codex/`.
- If a rule is only useful in one repository, do not add it globally.
- If a rule should protect every Codex session, prefer the user scope.

### 3. Format `.rules` consistently

Use this house style unless the file already follows a different stable style:

```python
prefix_rule(
    pattern = ["git"],
    decision = "allow",
    justification = "Allow git throughout this workspace project",
)
```

Formatting rules:
- one `prefix_rule(...)` per block
- multi-line function call
- one field per line
- spaces around `=`
- trailing comma after every field
- blank line between logical rule groups when it helps readability
- add short comments only when they explain policy intent, not obvious syntax

For forbidden rules, prefer a justification that explains the safer alternative when one exists.

### 4. Add inline tests when the rule is non-obvious

Use `match` and `not_match` when:
- the prefix is subtle
- multiple similar commands exist
- the rule could over-match

Keep tests short and representative.

### 5. Validate before finishing

Use `codex execpolicy check` for rule validation. Typical pattern:

```bash
codex execpolicy check --rules /path/to/default.rules --pretty -- git status
```

Validation checklist:
- the file parses successfully
- the reported decision matches intent
- the matched rule is the expected one
- restrictive rules still win over broader allow rules

Use multiple `--rules` flags when you need to reason about combined files.

### 6. Tell the user about reload behavior

- After changing `.rules`, tell the user Codex may need a restart because Codex scans `rules/` at startup.
- After changing `config.toml`, prefer telling the user whether the change is for future sessions unless you explicitly verify live reload behavior for that setting.

## Output Expectations

When making a Codex config or rules change, report:
- which file you changed
- why that scope was chosen
- whether the change belongs to `config.toml` or `.rules`
- what validation you ran
- whether the user should restart Codex

## References

- Official placement and field details: `references/rules-and-config.md`
