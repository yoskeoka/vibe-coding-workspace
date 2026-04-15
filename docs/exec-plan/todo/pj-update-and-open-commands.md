# pj Update and Open Commands

> **Execution**: Use `/execute-task` to implement this plan.

## Objective

Replace the narrow `pj move` command with a general `pj update` mutation path,
and add lightweight remote-navigation commands so operators can jump from the
local spike CLI to the canonical GitHub Project quickly.

This plan covers three operator-facing improvements that belong together:

- replace `pj move` with `pj update`
- add `pj open`
- add `pj url`

It also updates option-value resolution so humans and agents can use shorter,
more forgiving input forms without depending on GitHub's exact display labels.

## Background

The current spike can:

- create items with `pj add`
- move items between statuses with `pj move`
- sync and list cached task state

That leaves two practical gaps:

1. Existing items cannot be updated when title/body or custom field values need
   correction, even though remote state changes and full re-triage can easily
   invalidate an earlier `pj add` body.
2. Humans who want to switch from CLI triage to the GitHub Projects UI do not
   have a dedicated command to open or print the canonical Project URL.

Recent discussion also settled these product directions for the spike:

- breaking CLI changes are acceptable at this stage
- `pj move` should be removed rather than kept as a parallel command
- `pj update` should become the single mutation path for existing items
- `pj open` / `pj url` are preferred over a `pj remote ...` namespace
- option values should be more forgiving:
  - case-insensitive
  - `-`, `_`, and spaces normalized together
  - unique prefix matching allowed
  - `repo` additionally accepts integer selection using a stable alphabetical
    order persisted in cache

## Spec Changes

### `docs/specs/github-projects-task-cli.md`

- Replace `pj move` with `pj update` in the command set and examples.
- Define `pj update --item <id>` as the mutation command for existing Project
  items.
- Specify supported update fields:
  - `title`
  - `body`
  - `status`
  - `repo`
  - `kind`
  - `priority`
- Define `pj open` as the command that opens the canonical GitHub Project URL
  in the operator's browser.
- Define `pj url` as the command that prints the canonical GitHub Project URL to
  stdout without opening it.
- Define `pj update --body-file <path>` as a supported way to load longer body
  text from disk instead of inline flag text.
- Define the `pj update` body-input conflict rule:
  - reject `--body` and `--body-file` when both are provided
  - fail clearly when the `--body-file` path cannot be read
- Define value-resolution rules for enum-like inputs:
  - for `repo`, canonical slug exact match using the full source identity
  - recognized alias exact match
  - for `repo`, integer index using ordered repo-option metadata stored in
    cache, sorted ascending by canonical slug
  - unique prefix match after normalization against canonical slugs and aliases
  - ambiguous or unknown input must fail clearly
- Define repo-option metadata for runtime resolution:
  - remote display value remains the basename (`ww`, `ai-arena`, etc.)
  - each repo option also carries:
    - `source_type`
    - `source_url`
    - `canonical_slug`
    - alias values derived from the source metadata
- Define `github-repo` as the initial `source_type`.
- Define the `github-repo` canonical slug as `github.com/<owner>/<repo>`,
  derived from the repo option's `source_url`.
- Clarify that `repo` runtime resolution uses this enriched ordered metadata
  from cache rather than only the current remote field display names.
- Clarify that `pj init` provisions `Workspace Repo` display values from
  `setup.sh`, while also deriving the enriched repo metadata from the same
  source.
- Clarify that runtime resolution uses cached remote field options and enriched
  repo metadata, not a hardcoded repo list, and therefore requires a valid
  cache/field snapshot.

### `docs/specs/triage-tasks.md`

- Replace references to `pj move` with `pj update --status ...`.
- Clarify that the skill may correct existing items through `pj update` instead
  of forcing manual GitHub edits or item recreation for `Repo`, `Kind`, and
  `Priority`.
- Update any `Workspace Repo` contract that currently hardcodes a fixed repo
  option list so it instead points to `setup.sh`-derived provisioning plus
  cached enriched repo metadata.
- Optionally reference `pj url` / `pj open` as the preferred Project URL path
  once the commands exist, instead of hand-building the URL from cache metadata.

## Code Changes

### `tools/pj/internal/pj/app.go`

- Remove the `move` subcommand registration.
- Add `update`, `open`, and `url` subcommands.
- Add dedicated usage/help text for the new commands instead of relying on the
  current minimal top-level usage line.
- Implement `runUpdate` with partial-update semantics so only explicitly passed
  fields are changed.
- Keep cache refresh behavior after successful mutation.

### `tools/pj/internal/pj/`

- Extend the GitHub client interface and implementation to update:
  - draft item title/body when supported
  - `Status`
  - `Workspace Repo`
  - `Kind`
  - `Priority`
- Reuse a shared value resolver for `repo`, `status`, `kind`, and `priority`.
- Replace map-only repo option metadata with an ordered cache representation
  that preserves stable alphabetical ordering by canonical slug.
- Derive repo-option metadata from `setup.sh`, including:
  - remote display basename
  - `source_type=github-repo`
  - `source_url`
  - `canonical_slug`
  - aliases such as basename and `owner/repo`
- Resolve repo integers against the ordered cached repo metadata.
- Normalize alias input before repo option lookup.
- Apply prefix matching to both canonical slugs and aliases, with clear
  ambiguity errors.
- Add URL helpers that derive the canonical Project URL from cached project
  metadata for `pj url` and `pj open`.

### Command ergonomics

- Support both `--flag value` and `--flag=value` forms through the Go flag
  parser.
- Add `--body-file` to `pj update` so longer handoff text can avoid shell-
  expansion hazards and awkward inline quoting.
- Reject `--body` and `--body-file` when both are provided, with a clear error.
- Ensure errors for ambiguous prefix matches name the conflicting candidates.

## Design Decisions

Past decisions:
- The GitHub Project is the canonical remote source of truth.
- The runtime field/value surface should follow the current remote Project
  schema after `pj init` / `pj sync`.
- Breaking CLI changes are acceptable during this spike.

Apply the same reasoning here:

- `pj update` should become the single existing-item mutation path.
- Runtime enum resolution should use the cached remote field options rather than
  code-local constants, especially for `Workspace Repo`.
- For `repo`, internal identity should be stronger than the remote display
  label: keep basename values in GitHub Projects, but resolve against canonical
  source-backed slugs such as `github.com/yoskeoka/ww`.
- `pj open` / `pj url` should prefer the canonical Project metadata already
  stored in cache instead of re-deriving the URL ad hoc in higher-level skills.

## Sub-tasks

- [ ] Update `docs/specs/github-projects-task-cli.md` for `pj update`, `pj open`,
      `pj url`, and the new value-resolution rules
- [ ] Update `docs/specs/triage-tasks.md` so the skill uses `pj update --status`
      instead of `pj move`, and can rely on `pj open` / `pj url`
- [ ] [parallel] Design the normalized resolver behavior for `repo`, `status`,
      `kind`, and `priority`, including ambiguous-input errors
- [ ] [parallel] Define the ordered cache shape for enriched repo-option
      metadata and how it is derived from `setup.sh`
- [ ] [parallel] Design the operator-facing help text and command examples for
      `pj update`, `pj open`, and `pj url`
- [ ] [depends on: spec updates, resolver design] Implement `pj update` and
      remove `pj move`
- [ ] [depends on: URL helper design] Implement `pj open` and `pj url`
- [ ] [depends on: update implementation] Add or update tests covering partial
      field updates, alias resolution, prefix resolution, and repo integer
      selection
- [ ] [depends on: all implementation tasks] Verify the CLI can:
      - update status through `pj update`
      - update title/body/custom fields on an existing item
      - print the canonical Project URL
      - open the canonical Project URL

## Verification

- Confirm `pj help` and command-specific help no longer mention `pj move`
- Confirm `pj update --item <id> --status todo` resolves to `Todo`
- Confirm `pj update --item <id> --status in-progress` resolves to
  `In Progress`
- Confirm `pj update --item <id> --repo 1` resolves using cache metadata saved
  in ascending alphabetical order by repo canonical slug
- Confirm `pj update --item <id> --repo github.com/yoskeoka/ww` resolves to the
  `ww` display value
- Confirm `pj update --item <id> --repo yoskeoka/ww` and `--repo ww` both
  resolve through alias metadata
- Confirm `pj update --item <id> --repo <unique-prefix>` succeeds and
  ambiguous prefixes fail clearly
- Confirm `pj url` prints the same canonical Project URL that `triage-tasks`
  currently computes from cache metadata
- Confirm `pj open` launches that URL successfully in the current environment
  or fails clearly when browser opening is unavailable

## Expected Outcome

- The spike has one coherent update command for existing Project items
- Operators no longer need GitHub UI edits or item recreation for ordinary
  Project corrections
- Humans can jump to the canonical Project from the CLI with a dedicated command
- Enum-like CLI values become easier for both humans and agents to provide
  correctly without memorizing GitHub's exact display strings
