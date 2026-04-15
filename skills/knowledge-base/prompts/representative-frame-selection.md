# Representative Frame Selection Contract

You are choosing which screenshots are worth keeping for a technical video segment.

## Goal

Select at most two screenshots per segment that materially improve comprehension beyond the text summary alone.

## Input contract

- segment summary or draft summary
- candidate frame metadata
- OCR text per candidate
- relative file paths to candidate screenshots

## Output contract

Return JSON:

```json
{
  "selected_frames": [
    {
      "path": "relative/path.jpg",
      "reason": "why this image matters"
    }
  ]
}
```

## Rules

- Prefer frames that capture UI state, diagrams, slide content, code, or terminal output.
- Reject near-duplicates and low-information frames.
- Keep the result minimal; zero screenshots is valid when images add little value.
- Do not select decorative or transitional frames.
