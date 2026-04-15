#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import sys

from kb_video_pipeline import (
    DEFAULT_FRAME_INTERVAL_SEC,
    DEFAULT_MAX_FRAMES,
    DEFAULT_SEGMENT_LENGTH_SEC,
    MAX_RECOMMENDED_OCR_BATCH_SIZE,
    dependency_help_text,
    detect_runtime_profile,
    ensure_dependency_status,
    run_ingest_pipeline,
    run_skim_pipeline,
)


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="Skim or ingest video-backed sources for the workspace knowledge base."
    )
    parser.add_argument(
        "mode",
        nargs="?",
        choices=["skim", "ingest"],
        help="Pipeline mode. Use skim first, then ingest when the source is worth keeping.",
    )
    parser.add_argument("source_url", nargs="?", help="Original source URL or direct video URL.")
    parser.add_argument(
        "--video-url",
        help="Canonical video URL when the source URL is a thin article wrapper.",
    )
    parser.add_argument(
        "--source-slug",
        help="Optional override for the durable source slug and job directory prefix.",
    )
    parser.add_argument(
        "--workspace-relevance",
        help="Short note about why this source matters to the workspace.",
    )
    parser.add_argument(
        "--job-dir",
        help="Existing or fixed job directory for resume/debug workflows.",
    )
    parser.add_argument(
        "--scratch-root",
        help="Optional persistent scratch root such as .local/kb-ingest.",
    )
    parser.add_argument(
        "--frame-interval-sec",
        type=int,
        default=None,
        help="Frame extraction interval in seconds. When omitted, the CLI chooses a default from the detected runtime profile.",
    )
    parser.add_argument(
        "--segment-length-sec",
        type=int,
        default=DEFAULT_SEGMENT_LENGTH_SEC,
        help="Transcript chunk size for segment building.",
    )
    parser.add_argument(
        "--max-frames",
        type=int,
        default=DEFAULT_MAX_FRAMES,
        help="Hard cap on extracted frames kept for OCR and review.",
    )
    parser.add_argument(
        "--ocr-batch-size",
        type=int,
        default=None,
        help=(
            "OCR scheduler batch size. When omitted, the CLI chooses a default from the detected runtime profile; "
            f"{MAX_RECOMMENDED_OCR_BATCH_SIZE} remains the hard ceiling."
        ),
    )
    parser.add_argument(
        "--ocr-lang",
        default="en",
        help="PaddleOCR language code. Keep 'en' unless the source materially benefits from another model.",
    )
    parser.add_argument(
        "--transcribe-command",
        help=(
            "Fallback shell command when subtitles are missing. "
            "Use {video_path} and {job_dir} placeholders. The command must print VTT to stdout."
        ),
    )
    parser.add_argument(
        "--check-deps",
        action="store_true",
        help="Validate external and Python dependencies, then exit.",
    )
    return parser


def validate_args(args: argparse.Namespace) -> None:
    if args.check_deps:
        return
    if not args.mode or not args.source_url:
        raise SystemExit("mode and source_url are required unless --check-deps is used.")
    runtime_profile = detect_runtime_profile()
    if args.frame_interval_sec is None:
        args.frame_interval_sec = runtime_profile.recommended_frame_interval_sec
    if args.ocr_batch_size is None:
        args.ocr_batch_size = runtime_profile.recommended_ocr_batch_size
    if args.ocr_batch_size > MAX_RECOMMENDED_OCR_BATCH_SIZE:
        raise SystemExit(
            f"--ocr-batch-size {args.ocr_batch_size} is too high for this pipeline. "
            f"Use {MAX_RECOMMENDED_OCR_BATCH_SIZE} or less."
        )
    if args.frame_interval_sec < 5:
        raise SystemExit("--frame-interval-sec below 5 is too aggressive for the intended lightweight pipeline.")


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    validate_args(args)

    if args.check_deps:
        status = ensure_dependency_status()
        runtime_profile = detect_runtime_profile()
        payload = {
            "ok": status.ok,
            "yt_dlp": status.yt_dlp,
            "ffmpeg": status.ffmpeg,
            "paddleocr_import": status.paddleocr_import,
            "paddle_import": status.paddle_import,
            "pillow_import": status.pillow_import,
            "setuptools_import": status.setuptools_import,
            "missing": status.missing_items(),
            "runtime_profile": {
                "platform_system": runtime_profile.platform_system,
                "platform_release": runtime_profile.platform_release,
                "machine": runtime_profile.machine,
                "is_wsl": runtime_profile.is_wsl,
                "total_memory_gb": runtime_profile.total_memory_gb,
                "gpu_available": runtime_profile.gpu_available,
                "recommended_ocr_batch_size": runtime_profile.recommended_ocr_batch_size,
                "recommended_frame_interval_sec": runtime_profile.recommended_frame_interval_sec,
            },
        }
        print(json.dumps(payload, indent=2))
        if not status.ok:
            print("", file=sys.stderr)
            print(dependency_help_text(), file=sys.stderr)
            return 1
        return 0

    try:
        if args.mode == "skim":
            result = run_skim_pipeline(args)
        else:
            result = run_ingest_pipeline(args)
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        return 1

    print(json.dumps(result, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
