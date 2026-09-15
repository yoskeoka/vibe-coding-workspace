# Preserve GitHub review visibility in PR follow-up polling

## Summary

`gh-pr-followup poll` can report no new review material even when Copilot has
submitted a review. The compact mode does not fetch review summaries at all,
and it advances the persisted timeline/comment cursor after every successful
API poll. If the caller loses the helper's output, a later poll on the same
head omits the already-cursored review comments.

## Evidence

On 2026-09-16, Copilot submitted inline reviews for `yoskeoka/ai-arena#368`
and `yoskeoka/reversi-ai-arena#51`. A polling invocation chained after a
30-second sleep advanced its state without returning the emitted payload to
the caller. The following compact polls returned empty `timeline_events` and
`inline_comments`; direct GitHub review endpoints still returned the Copilot
review body and comments.

The helper resets comment/timeline cursors to zero when the head changes, so a
push is not itself the direct loss mechanism. The durable gap is that review
bodies are absent from compact output, and cursor advancement is not coupled to
confirmed delivery of that output.

## Follow-up boundary

Update the helper or its caller contract so normal PR follow-up reliably shows
new review summaries and comments. Define a recoverable inspection path for
the current head, such as review-ID cursor state plus an explicit replay mode,
and ensure compact output includes actionable review bodies or a clear pointer
to them. Preserve bounded output and avoid duplicate noise across ordinary
polls.

## Status

Open.
