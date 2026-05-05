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

For AI-facing query or filing-back work, treat English `docs/kb/**` as canonical and exclude `docs/kb/ja/**` unless the user explicitly asks to inspect or update the Japanese mirror.

## Preparation for file-changing work

Knowledge-base ingest normally changes files, so use the same branch hygiene as other workspace work before editing `docs/kb/`, `docs/issues/`, renderer dependencies, or skill files.

1. Check the current state:

```sh
git status --short --branch
```

2. If the working tree is dirty, do not start ingest on top of it. Preserve or hand off the existing work first.
3. Create a fresh worktree from the latest `origin/main` with the globally installed `ww` CLI:

```sh
ww create <type>/kb-<short-name>
cd "$(ww cd <type>/kb-<short-name>)"
```

Keep the existing workflow branch types and use `kb-` as the description prefix:

- `docs/kb-<source-or-topic>` for source notes, wiki pages, and KB logs only.
- `chore/kb-<short-name>` when changing `tools/kb`, `requirements-kb.txt`, `skills/knowledge-base`, video ingest tooling, or other KB infrastructure.
- `feat/kb-<plan-name>` or `fix/kb-<plan-name>` only when executing an approved plan.

PR titles should keep the workflow type visible and use KB as the scope, for example `docs: ingest KB source for <topic>` or `chore: update KB video ingest workflow`.

Read-only KB queries that do not write files may skip branch creation.

## Ingest workflow

Follow `docs/kb/ingest.md` for the ingest flow and `docs/kb/schema.md` for durable source-note and wiki-page formats. Do not duplicate or override those rules here.

Choose the acquisition path before drafting:
- Use the normal conversational ingest flow for ordinary article URLs the agent can read directly.
- Use the video-backed workflow below for direct videos and thin article wrappers whose value lives in the video.
- Use the `markitdown` fallback below for unsupported file-like sources such as local PDF, DOCX, PPTX, XLSX, EPUB, or direct document URLs.

## Video-backed workflow

Use the skill-local pipeline when `docs/kb/ingest.md` calls for video-backed processing. Treat the generated artifacts as drafts; final durable output must still follow `docs/kb/schema.md` and `docs/kb/ingest.md`.

### Setup

- External binaries:
  - macOS/Homebrew: `brew install yt-dlp ffmpeg`
  - Debian/Ubuntu: `sudo apt-get install ffmpeg`
  - If `yt-dlp` from the distro is too old, install a current standalone release or use Homebrew/Linuxbrew.
- Python runtime:
  - `python3 -m pip install -r skills/knowledge-base/requirements-video.txt`
  - Install the Paddle runtime that matches the actual environment. Use a CUDA build only when the host and Paddle installation both support GPU execution.

### Dependency check

Run:

```sh
python3 skills/knowledge-base/scripts/kb_video_ingest.py --check-deps
```

The command fails fast with setup guidance if `yt-dlp`, `ffmpeg`, `paddleocr`, `paddlepaddle`, or `Pillow` is missing.
It also reports the runtime profile that the pipeline will use to choose safe defaults.

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
- If `--frame-interval-sec` or `--ocr-batch-size` is omitted, the CLI picks defaults from the detected runtime profile.
- Review `outputs/skim.md` before deciding whether to ingest.

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

## MarkItDown fallback workflow

Use this only for file-like sources that the normal conversational flow handles poorly. It is a preprocessing path, not a durable artifact format.

### Dependency check

Run:

```sh
python3 skills/knowledge-base/scripts/kb_markitdown_ingest.py --check-deps
```

The command passes when either:
- `markitdown` is already installed as a local executable, or
- `uv` is available so the helper can run `markitdown[pdf,docx,pptx,xlsx]` in an isolated environment

The default fallback intentionally excludes OCR-heavy plugins and Document Intelligence paths. If the source is scanned or structurally complex, stop and choose a higher-fidelity path instead of forcing this helper.

### Convert first

For a local file:

```sh
python3 skills/knowledge-base/scripts/kb_markitdown_ingest.py \
  convert \
  path/to/source.pdf \
  --workspace-relevance "why this matters here" \
  --scratch-root .local/kb-ingest
```

For a direct document URL:

```sh
python3 skills/knowledge-base/scripts/kb_markitdown_ingest.py \
  convert \
  "https://example.com/files/source.docx" \
  --workspace-relevance "why this matters here" \
  --scratch-root .local/kb-ingest
```

`convert` writes:
- `outputs/converted.md`
- `outputs/source-context.md`
- `metadata.json`

Use `source-context.md` as the drafting boundary for the durable KB note. Do not commit `converted.md` or other raw conversion artifacts into `docs/kb/`.

### Stop conditions

Do not use this fallback as the final answer when:
- the PDF is scanned and needs OCR
- multi-column ordering is badly broken
- slides, diagrams, or tables lose too much meaning
- the source is actually better handled by the video-backed pipeline

## Query filing-back workflow

Follow the query filing-back rules in `docs/kb/schema.md`.

## Lint workflow

Follow the lint rules in `docs/kb/schema.md`.

## Rendering

- Use `tools/kb check` to validate the structure.
- Use `tools/kb build` to build the Pages site.
- Use `tools/kb serve` for local preview if needed.
- `tools/kb` will prefer `uv` when available and otherwise use the local Python MkDocs installation.
