# Spec: GitHub Projects Task CLI Spike

## Goal
Provide a small Go CLI that lets this workspace sync, inspect, create, update, and navigate GitHub Project tasks while keeping a local JSON cache optimized for AI reads.

## Location
- The CLI lives under `tools/pj/`.
- The command entrypoint is `tools/pj/cmd/pj/`.
- Supporting packages live under `tools/pj/internal/`.

This layout is preferred over `scripts/cmd/pj` because the spike is a compiled Go utility, not a shell-script collection, and the workspace root does not currently have a shared `go.mod` suitable for `cmd/pj/`.

## Implementation Testability
- `Run(args, stdout, stderr)` is the public command entrypoint.
- Remote-backed commands MUST obtain GitHub Project clients through dependencies owned by the internal command app, not through mutable package-level client factory state.
- Command tests MAY instantiate the internal app with stub client dependencies so remote-backed command tests can run independently when their config/cache fixtures are isolated.
- Remote-backed commands MUST continue to create clients lazily only after argument, config, and cache validation reaches the point where remote access is needed.

## Configuration
- Authentication MUST use `gh auth token`.
- The GitHub token needs ProjectV2 scopes:
  - write project scope for `init`, `add`, and `update`
  - `read:project` for `sync` and remote-backed reads
- The canonical workspace board is a dedicated ProjectV2 named `Workspace Task Triage`.
- The canonical workspace board is owner-scoped (`user` or `org`), not repository-owned.
- The owner-scoped board MAY be linked to repositories owned by the same user or organization so it appears in those repositories' Projects tabs.
- The CLI MUST provide an explicit bootstrap command, `pj init`, that resolves this board by name for the configured owner before creating a new board.
- The CLI MUST treat the owner target as explicit local configuration, not as an incidental cache side effect.
- The owner target configuration MUST live in `.local/pj/config.json`.
- The target project is configured by explicit owner configuration plus cached project metadata:
  - `owner`
  - `owner_type` (`user` or `org`)
  - `project_number`
- `pj init` MUST persist the resolved owner target into `.local/pj/config.json`.
- `pj init` or a successful `sync` MUST persist project identity in `.local/pj/cache.json` so later commands can reuse it without manual project-number lookup.
- When `.local/pj/config.json` is missing but the cache already contains owner metadata from an older version, the CLI MAY seed config from that cache data as a compatibility migration.
- Owner-target flags MUST NOT silently override stored owner configuration. Switching owner scope requires an explicit configuration action.

## Cache
- Cache root: `.local/pj/`
- Owner configuration file: `.local/pj/config.json`
- Primary file: `.local/pj/cache.json`
- The default cache path MUST resolve to the workspace-root `.local/pj/cache.json` even when the CLI is invoked from `tools/pj/`.
- The default config path MUST resolve to the workspace-root `.local/pj/config.json` even when the CLI is invoked from `tools/pj/`.
- The cache stores:
  - project identity and project ID
  - custom field IDs and single-select option IDs
  - ordered enriched repo-option metadata for runtime `repo` resolution
  - a normalized item list suitable for AI inspection
  - the last sync timestamp

The config stores:
- the active owner target:
  - `owner`
  - `owner_type`

## Commands

### `pj init`
- Requires `--owner` and `--owner-type` unless that metadata already exists in config
- Resolves the canonical `Workspace Task Triage` board by owner and title before creating a new board
- Creates the canonical board when absent
- Provisions missing custom `Workspace Repo`, `Kind`, and `Priority` fields as single-select workflow fields during bootstrap
- Provisions `Workspace Repo` display values from workspace repo metadata derived from `setup.sh` plus the workspace repository itself
- Derives enriched repo metadata from the same source:
  - `source_type`
  - `source_url`
  - `canonical_slug`
  - alias values such as basename and `owner/repo`
- Uses `github-repo` as the initial `source_type`
- Defines a `github-repo` canonical slug as `github.com/<owner>/<repo>`, derived from the repo option's `source_url`
- Keeps GitHub Project display values as repo basenames (`ww`, `ai-arena`, etc.) even though runtime resolution may use richer identity
- Uses these canonical option sets for provisioned non-repo fields:
  - `Kind`: `Feature`, `Bug`, `Chore`, `Research`
  - `Priority`: `High`, `Medium`, `Low`
- Keeps `--repo` as the CLI flag, and stores the normalized item property as `repo`, even though the remote ProjectV2 field name and cached field-metadata map key remain `Workspace Repo`
- Reuses existing compatible workflow fields when they already exist instead of creating duplicates
- Writes `.local/pj/config.json` with the resolved owner target
- Writes `.local/pj/cache.json` with the resolved project identity and current remote snapshot
- Fails clearly if more than one board with the canonical title exists for the same owner
- Fails clearly when required workflow fields are still missing or incompatible after bootstrap, naming the blocking fields and compatibility problem
- Fails clearly when `--owner` and/or `--owner-type` conflict with stored config; the operator must use `pj config set` or `pj config clear` before switching owner scope

### `pj sync`
- Resolves the target project through GitHub GraphQL
- Loads complete project field metadata with cursor-based pagination
- Loads complete project items with cursor-based pagination
- Loads complete normalized item field values with cursor-based pagination when GitHub reports more field values for an item
- Writes `.local/pj/cache.json`
- Reuses the stored owner target from `.local/pj/config.json` when owner flags are omitted
- Reuses cached project identity when `--project` is omitted
- MUST work after `pj init` without requiring manual project-number lookup
- Fails clearly when owner flags conflict with stored config; the operator must switch config explicitly instead of mixing old and new owner metadata

### `pj config`
- `pj config show` prints the active owner target from `.local/pj/config.json`
- `pj config set --owner <owner> --owner-type user|org` explicitly replaces the active owner target
- `pj config clear` removes the active owner target and cached project snapshot so the next `pj init` or `pj sync` must establish them again
- `pj config set` MUST also clear `.local/pj/cache.json` if the cached project belongs to a different owner target, so the workspace cannot accidentally reuse stale project identity after a scope switch

### `pj repo-link`
- `pj repo-link status <owner>/<repo>` reports whether the cached canonical ProjectV2 is linked to the target repository.
- `pj repo-link add <owner>/<repo>` links the cached canonical ProjectV2 to the target repository.
- `pj repo-link remove <owner>/<repo>` unlinks the cached canonical ProjectV2 from the target repository.
- The target repository is explicit command input and is separate from the configured Project owner target.
- The target repository owner MUST match the configured Project owner because GitHub only supports linking a ProjectV2 to repositories owned by the same user or organization.
- The commands MUST use the cached project identity from `.local/pj/cache.json`; operators must run `pj init` or `pj sync` first if the cache is missing project metadata.
- `status` MUST fail clearly if the linked repository list exceeds the current single-page limit instead of silently reporting a partial result.
- `add` and `remove` SHOULD be idempotent from an operator perspective: adding an already-linked repository or removing an unlinked repository should report the current state without treating it as a failure.
- Setting the linked repository as the Project's default repository is not required for the workspace workflow and is not part of the `pj repo-link` command set. The current workflow creates Project draft items directly, so repository-tab discoverability is the required behavior.

### `pj list`
- Reads `.local/pj/cache.json`
- Prints a readable task table or line-oriented summary
- Supports optional filtering by normalized fields when practical

### `pj add`
- Creates a draft item with a title and optional body
- Optionally sets `Status`, `Repo`, `Kind`, and `Priority`
- Supports `--body <text>` and `--body-file <path>` for draft item body input
- Rejects `--body` and `--body-file` when both are provided
- Fails clearly when the `--body-file` path cannot be read
- Refreshes the local cache after successful mutation

### `pj update`
- Updates an existing project item by Project item ID
- Requires `--item <id>`
- Supports partial updates: only explicitly provided fields are changed
- Supports these fields:
  - `title`
  - `body`
  - `status`
  - `repo`
  - `kind`
  - `priority`
- Supports `--body <text>` and `--body-file <path>` for body input
- Rejects `--body` and `--body-file` when both are provided
- Fails clearly when the `--body-file` path cannot be read
- Updates draft item title/body when GitHub supports it for the selected item
- Refreshes the local cache after successful mutation
- Replaces the former `pj move` command; `pj move` is not a supported command

### `pj url`
- Reads `.local/pj/cache.json`
- Prints the canonical GitHub Project URL to stdout without opening a browser

### `pj open`
- Reads `.local/pj/cache.json`
- Opens the canonical GitHub Project URL in the operator's browser
- Fails clearly when browser opening is unavailable in the current environment

## Data model
Each cached item MUST expose enough normalized data for AI and human use:
- project item ID
- draft issue ID when available and different from the Project item ID
- content type (`DraftIssue`, `Issue`, `PullRequest`, or unknown)
- title
- repository hint when available
- current values for `Status`, `Repo`, `Kind`, and `Priority`
- optional URL for non-draft content

The cache MUST also expose ordered repo-option metadata:
- The order is ascending by `canonical_slug`.
- The remote display value remains the basename used by the `Workspace Repo` field.
- Runtime resolution uses this enriched cache data instead of a hardcoded repo list.
- Cache refresh derives metadata from `setup.sh` plus the workspace repository itself; field mutation still requires the resolved display value to exist in the cached remote `Workspace Repo` field options.

## Value resolution
- `status`, `kind`, and `priority` values resolve against the current cached remote field options.
- `repo` values resolve against ordered enriched repo-option metadata from cache.
- Enum-like inputs are case-insensitive and normalize `-`, `_`, and whitespace together.
- Resolution order is:
  - exact canonical slug match for `repo`
  - exact alias match
  - integer index for `repo`, using the ordered metadata stored in cache
  - unique normalized prefix match against canonical values and aliases
- Ambiguous or unknown inputs must fail clearly.
- Ambiguous prefix errors must name the conflicting candidates.

## Error handling
- If `gh auth token` fails, commands that contact GitHub must fail with a clear authentication message.
- If the token lacks GitHub Projects scopes, the CLI must surface the GraphQL scope error clearly.
- If GitHub returns a non-success HTTP response, the CLI must include the HTTP status in the returned error even when the response body is not valid JSON.
- If the cache is missing, local-cache commands must tell the operator to run `pj init` or `pj sync`.
- If the owner config is missing, commands that require remote project resolution must tell the operator to run `pj init` or `pj config set`.
- If provided owner flags conflict with stored owner config, the CLI must fail with a clear mismatch error instead of partially mixing values.
- If a field is missing from the project, mutations must fail with a clear field-name error instead of silently skipping.
- If `pj init` resolves or creates the canonical board but cannot provision or reconcile the required workflow fields, it must name the blocking fields in the returned error after writing the latest cache snapshot.
- If a required field exists with an unsupported type or missing required single-select options, `pj init` must fail with a clear compatibility error instead of silently mutating an unknown schema.
- If GitHub returns a paginated project field, project item, or item field-value response, the CLI must follow cursors until the full connection is loaded instead of silently truncating the cache.
- If a paginated response is malformed and cannot provide the next cursor while reporting more pages, the CLI must fail clearly instead of writing an incomplete cache.

## CI verification
- The repository MUST provide a dedicated GitHub Actions workflow named `check-pj` for the workspace-local `tools/pj` Go module.
- The workflow MUST run on pull requests targeting `main` when any of these paths change:
  - `tools/pj/**`
  - `docs/specs/github-projects-task-cli.md`
  - `.github/workflows/check-pj.yml`
- The workflow MUST check out the repository with `actions/checkout@v6`.
- The workflow MUST set up Go from `tools/pj/go.mod`.
- All Go checks MUST run with `tools/pj` as the working directory.
- The workflow MUST install pinned tool versions before checking the module:
  - `goimports`: `golang.org/x/tools/cmd/goimports@v0.44.0`
  - `staticcheck`: `honnef.co/go/tools/cmd/staticcheck@v0.7.0`
- The workflow MUST be non-mutating. Formatting checks must report files that need changes instead of rewriting them in CI.
- The workflow MUST run these checks:
  - `goimports` formatting check over all Go files in `tools/pj`
  - `staticcheck ./...`
  - `go vet ./...`
  - `go test ./...`
- Tool upgrades SHOULD be done by editing the pinned versions in `.github/workflows/check-pj.yml`, updating this spec in the same PR, and rerunning the local equivalents of the workflow checks.

## Non-Goals
- Full parity with `gh project`
- Complex search syntax
- Automatic project discovery across multiple boards
- General-purpose schema management for arbitrary ProjectV2 fields beyond the required workspace workflow fields
- Cross-project aggregation beyond the single configured workspace board
