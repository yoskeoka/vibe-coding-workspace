# KB Compile / Review Contract

You are converting a skimmed video job into durable knowledge-base updates.

## Goal

Produce source-note and wiki-update drafts that keep retrieval anchors while staying concise.

## Input contract

- normalized source metadata
- segment summaries with time anchors
- selected representative screenshots, if any
- workspace relevance framing from the user

## Output contract

Return JSON:

```json
{
  "source_note_summary": ["bullet", "bullet"],
  "follow_up_actions": ["bullet"],
  "wiki_targets": ["topics/foo", "tools/bar"],
  "wiki_update_notes": ["bullet", "bullet"],
  "log_entry": "one line"
}
```

## Rules

- Preserve concrete tools, libraries, commands, services, and document names.
- Keep only durable takeaways that belong in the KB.
- Reference time anchors when they make later retrieval easier.
- Mention screenshots only when they materially improve human skimming.
