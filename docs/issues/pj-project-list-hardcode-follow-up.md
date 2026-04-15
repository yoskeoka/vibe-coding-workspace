# pj init still hardcodes the workspace project list

## Summary

`pj init` currently provisions the `Workspace Repo` Project field from a
hardcoded workspace repo list.

That is acceptable for the current spike, but it becomes a design problem if we
later want either of these:

- repo options injected during `pj init`
- repo options managed through future `pj repo` add/remove operations

At that point, `pj init` will need a reliable source for the project-name /
project-URL mapping used to populate the remote Project field options, and that
source is not defined yet.

## Why This Matters

The current runtime validation approach can correctly use the cached remote
field options from `.local/pj/cache.json` after `pj init` / `pj sync`.
However, schema provisioning during `pj init` still needs an authoritative list
of repo options before that cache exists or before the remote schema is updated.

If the repo list becomes user-configurable, a hardcoded list in `pj init`
becomes the wrong source of truth.

## Open Question

What should become the source of truth for the workspace project list and its
project-name / project-URL mapping when `Workspace Repo` options are no longer
hardcoded?

Possible directions:

- derive it from `setup.sh`
- introduce explicit local config for `pj`
- add a future `pj repo` management surface and persist that state locally
- derive it from some other workflow-owned metadata file

## Priority

Medium. This does not block the current spike, but it should be resolved before
`Workspace Repo` options become operator-configurable.
