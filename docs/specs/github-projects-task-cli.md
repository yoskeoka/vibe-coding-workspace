# Spec: GitHub Projects Task CLI Spike

## Goal
Provide a small Go CLI that lets this workspace sync, inspect, create, and move GitHub Project tasks while keeping a local JSON cache optimized for AI reads.

## Location
- The CLI lives under `tools/pj/`.
- The command entrypoint is `tools/pj/cmd/pj/`.
- Supporting packages live under `tools/pj/internal/`.

This layout is preferred over `scripts/cmd/pj` because the spike is a compiled Go utility, not a shell-script collection, and the workspace root does not currently have a shared `go.mod` suitable for `cmd/pj/`.

## Configuration
- Authentication MUST use `gh auth token`.
- The GitHub token needs ProjectV2 scopes:
  - write project scope for `init`, `add`, and `move`
  - `read:project` for `sync` and remote-backed reads
- The canonical workspace board is a dedicated ProjectV2 named `Workspace Task Triage`.
- The CLI MUST provide an explicit bootstrap command, `pj init`, that resolves this board by name for the configured owner before creating a new board.
- If `pj init` creates the canonical board, it MUST persist the resulting project identity into the local cache so later commands can reuse it without manual project-number lookup.
- The target project is configured by command flags and/or cached metadata:
  - `owner`
  - `owner_type` (`user` or `org`)
  - `project_number`
- `pj init` or the first successful `sync` persists project identity in the local cache so later commands can reuse it.

## Cache
- Cache root: `.local/pj/`
- Primary file: `.local/pj/cache.json`
- The default cache path MUST resolve to the workspace-root `.local/pj/cache.json` even when the CLI is invoked from `tools/pj/`.
- The cache stores:
  - project identity and project ID
  - custom field IDs and single-select option IDs
  - a normalized item list suitable for AI inspection
  - the last sync timestamp

## Commands

### `pj init`
- Requires `--owner` and `--owner-type` unless that metadata already exists in cache
- Resolves the canonical `Workspace Task Triage` board by owner and title before creating a new board
- Creates the canonical board when absent
- Provisions missing custom `Workspace Repo`, `Kind`, and `Priority` fields as single-select workflow fields during bootstrap
- Uses these canonical option sets for provisioned fields:
  - `Workspace Repo`: `vibe-coding-workspace`, `ww`, `ai-arena`, `reversi-adventure`, `vim-learning-game`, `envdiff`
  - `Kind`: `Feature`, `Bug`, `Chore`, `Research`
  - `Priority`: `High`, `Medium`, `Low`
- Keeps `--repo` as the CLI flag and normalized cache key even though the remote ProjectV2 field name is `Workspace Repo`
- Reuses existing compatible workflow fields when they already exist instead of creating duplicates
- Writes `.local/pj/cache.json` with the resolved project identity and current remote snapshot
- Fails clearly if more than one board with the canonical title exists for the same owner
- Fails clearly when required workflow fields are still missing or incompatible after bootstrap, naming the blocking fields and compatibility problem

### `pj sync`
- Resolves the target project through GitHub GraphQL
- Loads project field metadata
- Loads project items and normalized field values
- Writes `.local/pj/cache.json`
- MAY reuse cached project identity when `owner`, `owner_type`, or `project_number` flags are omitted
- MUST work after `pj init` without requiring manual project-number lookup

### `pj list`
- Reads `.local/pj/cache.json`
- Prints a readable task table or line-oriented summary
- Supports optional filtering by normalized fields when practical

### `pj add`
- Creates a draft item with a title and optional body
- Optionally sets `Status`, `Repo`, `Kind`, and `Priority`
- Refreshes the local cache after successful mutation

### `pj move`
- Updates the `Status` field for a given project item ID
- Refreshes the local cache after successful mutation

## Data model
Each cached item MUST expose enough normalized data for AI and human use:
- project item ID
- content type (`DraftIssue`, `Issue`, `PullRequest`, or unknown)
- title
- repository hint when available
- current values for `Status`, `Repo`, `Kind`, and `Priority`
- optional URL for non-draft content

## Error handling
- If `gh auth token` fails, commands that contact GitHub must fail with a clear authentication message.
- If the token lacks GitHub Projects scopes, the CLI must surface the GraphQL scope error clearly.
- If GitHub returns a non-success HTTP response, the CLI must include the HTTP status in the returned error even when the response body is not valid JSON.
- If the cache is missing, local-cache commands must tell the operator to run `pj init` or `pj sync`.
- If a field is missing from the project, mutations must fail with a clear field-name error instead of silently skipping.
- If `pj init` resolves or creates the canonical board but cannot provision or reconcile the required workflow fields, it must name the blocking fields in the returned error after writing the latest cache snapshot.
- If a required field exists with an unsupported type or missing required single-select options, `pj init` must fail with a clear compatibility error instead of silently mutating an unknown schema.
- If a query result exceeds the current single-page limits, the CLI must fail clearly instead of silently truncating the cache.

## Non-Goals
- Full parity with `gh project`
- Complex search syntax
- Automatic project discovery across multiple boards
- General-purpose schema management for arbitrary ProjectV2 fields beyond the required workspace workflow fields
- Cross-project aggregation beyond the single configured workspace board
