# Follow Up: Make KB Video Ingest More Reliable on Low-Memory WSL Hosts

## Summary

The KB video-ingest flow could complete dependency checks and transcript fallback on the current WSL host, but the OCR-heavy `skim` path exited with status `137` during processing.

Observed runtime profile from `--check-deps`:

- `is_wsl: true`
- `gpu_available: false`
- `recommended_ocr_batch_size: 1`
- `recommended_frame_interval_sec: 45`

Even after reducing workload, the OCR phase remained too memory-constrained to finish reliably on this host.

## Why It Matters

- The current implementation works in principle, but the practical "first run" experience on weaker WSL machines is fragile.
- A video-backed ingest flow that fails after transcript/download setup creates wasted time and makes the skill less trustworthy.
- This weakens the intended skim-first workflow for thin article + embedded video sources, which is exactly the kind of source this pipeline should help with.

## Proposed Solution

- Reproduce the failure with explicit memory telemetry so the actual hot path is clear.
- Add a more conservative low-memory profile for WSL/CPU-only hosts, potentially including:
  - much lower default frame counts
  - earlier fallback to transcript-only skim
  - an OCR skip/degrade mode when memory pressure is too high
- Consider a CLI switch that allows transcript-first ingest even when OCR is unavailable or intentionally disabled.
- Document expected resource envelopes for video ingest so operators can predict whether full OCR skim is realistic on the current machine.
