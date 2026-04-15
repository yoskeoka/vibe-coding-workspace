# Segment Summary Contract

You are summarizing one segment from a technical video for the workspace knowledge base.

## Goal

Produce a compact, durable summary that helps later retrieval and human skimming.

## Input contract

- segment metadata with start/end times
- transcript excerpt for that segment
- OCR lines from representative candidate frames
- optional notes about why the source matters to the workspace

## Output contract

Return JSON:

```json
{
  "summary": "2-4 sentences",
  "highlights": ["bullet", "bullet"],
  "named_entities": ["tool", "command", "service"],
  "time_anchor_label": "MM:SS-MM:SS"
}
```

## Rules

- Preserve concrete tool names, command names, product names, APIs, and document names.
- Prefer durable technical takeaways over filler narration.
- Keep wording specific enough that a later human can jump back into the source.
- Do not invent facts that are not supported by transcript or OCR context.
