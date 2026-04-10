# Ingest Flow

The intended UX is conversational:

> Ingest these URLs into the knowledge base.

## Expected agent behavior

1. Read the given URLs or the provided source material.
2. Capture the user's framing:
   - why the links matter
   - which projects or themes they seem relevant to
3. Create one source note per URL in `sources/<year>/`.
4. Update the most relevant compiled wiki pages in `wiki/`.
5. Add a short entry to `wiki/log.md`.
6. If the source suggests a concrete experiment, record it as a follow-up bullet.

## Classification hints

Ask these questions during ingest:
- Which workspace project does this help, if any?
- Is this primarily a tool note, topic note, or reusable pattern?
- Is it a durable insight or just a temporary watch item?
- Should it modify an existing page or create a new one?

## Human review expectations

- The human can skim source notes for accuracy.
- The human can browse the compiled wiki or the rendered Pages site.
- The human can redirect emphasis in the next ingest request instead of editing the wiki manually.
