# pj init still hardcodes workspace repo metadata

## Summary

`pj init` currently provisions the `Workspace Repo` Project field from a
hardcoded workspace repo list with only basename-style display values.

That is acceptable for the current spike, but it becomes a design problem if we
later want any of these:

- repo options injected during `pj init`
- repo options managed through future `pj repo` add/remove operations
- stronger internal repo identity than the remote display value

At that point, `pj init` needs a reliable source for richer repo metadata, not
just display names. That metadata now includes at least:

- remote display basename
- `source_type`
- `source_url`
- canonical slug (for example `github.com/yoskeoka/ww`)
- alias material derived from the source metadata

The current hardcoded list is not a good long-term source of truth for that.

## Why This Matters

The current runtime validation approach can correctly use the cached remote
field options from `.local/pj/cache.json` after `pj init` / `pj sync`.
However, schema provisioning during `pj init` still needs an authoritative list
of repo options and enriched repo metadata before that cache exists or before
the remote schema is updated.

If repo metadata becomes user-configurable or is expected to support stronger
internal identity than the remote basename label, a hardcoded list in `pj init`
becomes the wrong source of truth.

## Open Question

What should become the source of truth for the workspace repo metadata used by
`pj init` to provision `Workspace Repo` display values and derive canonical repo
identity?

Possible directions:

- derive it from `setup.sh`
- introduce explicit local config for `pj`
- add a future `pj repo` management surface and persist that state locally
- derive it from some other workflow-owned metadata file

## Priority

Medium. This does not block the current spike, but it should be resolved before
`Workspace Repo` options become operator-configurable.

## Resolution

Resolved by the `pj-update-and-open-commands` execution task. `pj init` now
derives repo option metadata from `setup.sh` plus the workspace repository
itself, stores ordered enriched repo metadata in cache, and uses that metadata
for runtime `repo` value resolution.
