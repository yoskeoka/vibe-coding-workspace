# Workspace Knowledge Base

This directory holds a repo-native knowledge base for `vibe-coding-workspace`.

The operating model is simple:
- AI reads curated source notes in `sources/`
- AI compiles durable knowledge into `wiki/`
- humans browse the same Markdown on GitHub or via the rendered GitHub Pages site

## Primary workflows

- Ingest: ask the agent to ingest one or more URLs into the knowledge base
- Query: ask a question against the compiled wiki and file durable answers back when useful
- Lint: ask the agent to review the wiki for staleness, duplication, and weak links
- Build: run `tools/kb build`

See [schema.md](schema.md) for structure rules and [ingest.md](ingest.md) for the ingest flow.
