---
name: knowledge-base
description: Use when the user wants to ingest URLs or notes into the workspace knowledge base, update compiled wiki pages, file durable answers back into the wiki, or lint the knowledge base for staleness and weak structure.
metadata:
  author: yoskeoka
  version: '1.0.0'
---

# Knowledge Base

This is a workspace-only skill for maintaining `docs/kb/`.

## When to use

- The user says to ingest URLs, articles, notes, or references into the knowledge base
- The user asks a question that should be answered from `docs/kb/` and preserved
- The user asks to lint, review, or reorganize the knowledge base

## Read first

1. `docs/kb/schema.md`
2. `docs/kb/ingest.md`
3. `docs/kb/wiki/index.md`

## Ingest workflow

1. Read the provided URLs or source material.
2. Create one source note per source under `docs/kb/sources/<year>/`.
3. Update the most relevant pages under `docs/kb/wiki/`.
4. Update `docs/kb/wiki/index.md` if navigation should change.
5. Append a dated line to `docs/kb/wiki/log.md`.
6. If the source suggests a concrete experiment, record it in the relevant page.

Prefer updating existing wiki pages over creating new ones. Preserve provenance. Avoid copying long source text.
Keep concrete retrieval anchors such as service names, library names, product names, and document names when they are part of the source's value. Do not collapse `Render vs Cloud Run` into `easy backend vs fast backend` if the original concrete options are useful later.

## Video-backed workflow

Use the skill-local pipeline when the source is a direct video or a thin article whose value lives in an embedded/linked video.

### Setup

- External binaries:
  - macOS/Homebrew: `brew install yt-dlp ffmpeg`
  - Debian/Ubuntu/WSL2: `sudo apt-get install ffmpeg`
  - If `yt-dlp` from the distro is too old, install a current standalone release or use Homebrew/Linuxbrew.
- Python runtime:
  - `python3 -m pip install -r skills/knowledge-base/requirements-video.txt`
  - This pipeline is CPU-first. On WSL2/non-GPU machines, keep the CPU `paddlepaddle` package and do not swap in a GPU-only build.

### Dependency check

Run:

```sh
python3 skills/knowledge-base/scripts/kb_video_ingest.py --check-deps
```

The command fails fast with setup guidance if `yt-dlp`, `ffmpeg`, `paddleocr`, `paddlepaddle`, or `Pillow` is missing.

### Skim first

Start with `skim` when KB value is uncertain:

```sh
python3 skills/knowledge-base/scripts/kb_video_ingest.py \
  skim \
  "https://example.com/video-or-article" \
  --video-url "https://youtube.com/watch?v=..." \
  --workspace-relevance "why this matters here" \
  --scratch-root .local/kb-ingest
```

Notes:
- Omit `--video-url` when the source URL itself is the canonical video URL.
- Use `--video-url` for video-backed articles so the durable note can keep both the wrapper article and the actual video.
- The default frame interval is intentionally conservative for CPU-only machines.
- Keep `--ocr-batch-size` low on WSL2. The script rejects values above 4.

`skim` produces:
- `job.json` with job metadata
- normalized transcript, OCR, frame, and segment artifacts under the job directory
- prompt payloads for segment summaries, representative frame selection, and KB compile/review
- `outputs/skim.md` for human review

### Ingest after skim

When the skimmed source is worth keeping:

```sh
python3 skills/knowledge-base/scripts/kb_video_ingest.py \
  ingest \
  "https://example.com/video-or-article" \
  --video-url "https://youtube.com/watch?v=..." \
  --job-dir .local/kb-ingest/<existing-job>
```

`ingest` reuses the prepared job and writes:
- `outputs/source-note-draft.md`
- `outputs/wiki-update-draft.md`
- `outputs/log-entry-draft.md`

### Subtitle-first, transcription fallback

- The pipeline prefers existing subtitles or auto-subs from `yt-dlp`.
- If subtitles are missing, pass `--transcribe-command` with a backend that prints VTT to stdout.
- Temporary artifacts stay in OS temp storage by default, or under `.local/kb-ingest/` when `--scratch-root` is set.
- Raw transcripts, OCR dumps, and bulk candidate frames stay outside `docs/kb/`.

## Query filing-back workflow

If a user asks a question and the answer is durable:
- create or update a compact page under `docs/kb/wiki/`
- link it from a relevant topic, pattern, tool, or project page
- note the change in `docs/kb/wiki/log.md`

## Lint workflow

Review the knowledge base for:
- orphan pages
- duplicated topics
- stale date-sensitive claims
- weak cross-linking
- source notes that are not reflected in the wiki

## Rendering

- Use `tools/kb check` to validate the structure.
- Use `tools/kb build` to build the Pages site.
- Use `tools/kb serve` for local preview if needed.
- `tools/kb` will prefer `uv` when available and otherwise use the local Python MkDocs installation.
