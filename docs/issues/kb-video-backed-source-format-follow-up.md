# Follow Up: Normalize Existing KB Video-Backed Source Notes

## Summary

The Phaser Claude Code space-shooter ingest established a richer durable format for `video_backed_article` source notes:

- visible `## Source` section with source URL, video URL, and subtitle/transcript retrieval status
- selected skim screenshots, normally at least one per important time anchor
- detailed notes grouped under the same time anchors used in frontmatter
- frontmatter `selected_screenshots` entries that point to the durable image assets

Existing video-backed source notes should be reviewed and brought into the same shape where the source still has enough value to justify the work.

## Known Affected Source

- `docs/kb/sources/2026/2026-04-11-openai-gpt-5-4-phaser-tactical-rpg.md`

## Why It Matters

- `source_url` and `video_url` currently live mostly in YAML frontmatter, so they are easy to miss when browsing the rendered wiki.
- Source notes are meant to be a drill-down layer for KB search hits. When a video source is relevant, the reader should not need to rerun video processing just to recover basic source links, screenshots, or anchor-level detail.
- Keeping video-backed notes in one format makes future KB maintenance and renderer improvements simpler.

## Proposed Solution

- Revisit each existing `source_type: video_backed_article` or `source_type: video` note.
- Add a visible `## Source` section that exposes article/video links and transcript or subtitle retrieval status.
- Add selected screenshots under `docs/kb/assets/source-images/<year>/<source-slug>/` when the frames improve human skim.
- Add anchor-aligned detailed notes using existing `time_anchors`.
- Avoid committing full third-party verbatim transcripts unless the source license or user-provided material permits it; preserve detailed segment notes and retrieval anchors instead.

## Verification

- Run `tools/kb check`.
- Run `tools/kb build`.
- Confirm rendered source pages show the source/video links and screenshots.
