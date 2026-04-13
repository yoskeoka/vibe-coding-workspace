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
  - `read:project` for `sync` and remote-backed reads
  - write project scope for `add` and `move`
- The target project is configured by command flags and/or cached metadata:
  - `owner`
  - `owner_type` (`user` or `org`)
  - `project_number`
- The first successful `sync` persists project identity in the local cache so later commands can reuse it.

## Cache
- Cache root: `.local/pj/`
- Primary file: `.local/pj/cache.json`
- The cache stores:
  - project identity and project ID
  - custom field IDs and single-select option IDs
  - a normalized item list suitable for AI inspection
  - the last sync timestamp

## Commands

### `pj sync`
- Resolves the target project through GitHub GraphQL
- Loads project field metadata
- Loads project items and normalized field values
- Writes `.local/pj/cache.json`

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
- If the cache is missing, `pj list` must tell the operator to run `pj sync`.
- If a field is missing from the project, mutations must fail with a clear field-name error instead of silently skipping.

## Non-Goals
- Full parity with `gh project`
- Complex search syntax
- Automatic project discovery across multiple boards
- Cross-project aggregation beyond the single configured workspace board
