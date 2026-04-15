from __future__ import annotations

import json
import math
import os
import platform
import re
import shutil
import subprocess
import tempfile
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT_DIR = Path(__file__).resolve().parents[3]
SKILL_DIR = Path(__file__).resolve().parents[1]
PROMPTS_DIR = SKILL_DIR / "prompts"
DEFAULT_FRAME_INTERVAL_SEC = 30
DEFAULT_SEGMENT_LENGTH_SEC = 180
DEFAULT_OCR_BATCH_SIZE = 1
MAX_RECOMMENDED_OCR_BATCH_SIZE = 4
DEFAULT_MAX_FRAMES = 60
SUBTITLE_LANGS = "ja,en"
COMMON_ENTITY_STOPWORDS = {
    "a",
    "all",
    "an",
    "and",
    "anyway",
    "as",
    "back",
    "basically",
    "but",
    "by",
    "each",
    "every",
    "exactly",
    "for",
    "how",
    "in",
    "it",
    "its",
    "maybe",
    "my",
    "now",
    "of",
    "one",
    "or",
    "so",
    "than",
    "that",
    "the",
    "then",
    "there",
    "these",
    "this",
    "those",
    "to",
    "unless",
    "using",
    "well",
    "what",
}


@dataclass
class DependencyStatus:
    yt_dlp: str | None
    ffmpeg: str | None
    paddleocr_import: bool
    paddle_import: bool
    pillow_import: bool
    setuptools_import: bool

    @property
    def ok(self) -> bool:
        return bool(
            self.yt_dlp
            and self.ffmpeg
            and self.paddleocr_import
            and self.paddle_import
            and self.pillow_import
            and self.setuptools_import
        )

    def missing_items(self) -> list[str]:
        missing = []
        if not self.yt_dlp:
            missing.append("yt-dlp")
        if not self.ffmpeg:
            missing.append("ffmpeg")
        if not self.paddleocr_import:
            missing.append("paddleocr")
        if not self.paddle_import:
            missing.append("paddlepaddle")
        if not self.pillow_import:
            missing.append("Pillow")
        if not self.setuptools_import:
            missing.append("setuptools")
        return missing


@dataclass
class RuntimeProfile:
    platform_system: str
    platform_release: str
    machine: str
    is_wsl: bool
    total_memory_gb: float | None
    gpu_requested: bool
    gpu_available: bool
    recommended_ocr_batch_size: int
    recommended_frame_interval_sec: int


def utc_now() -> str:
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat()


def slugify(value: str) -> str:
    lowered = value.lower()
    normalized = re.sub(r"[^a-z0-9]+", "-", lowered).strip("-")
    return normalized or "video-source"


def format_seconds(seconds: float | int | None) -> str:
    if seconds is None:
        return "unknown"
    total = max(0, int(seconds))
    hours, rem = divmod(total, 3600)
    minutes, sec = divmod(rem, 60)
    if hours:
        return f"{hours:02d}:{minutes:02d}:{sec:02d}"
    return f"{minutes:02d}:{sec:02d}"


def run_command(
    args: list[str],
    *,
    cwd: Path | None = None,
    capture_output: bool = True,
    text: bool = True,
    check: bool = True,
) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        args,
        cwd=cwd,
        check=check,
        capture_output=capture_output,
        text=text,
    )


def write_json(path: Path, payload: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")


def read_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def load_prompt_contract(name: str) -> str:
    return (PROMPTS_DIR / name).read_text(encoding="utf-8").strip()


def ensure_dependency_status() -> DependencyStatus:
    import importlib.util

    return DependencyStatus(
        yt_dlp=shutil.which("yt-dlp"),
        ffmpeg=shutil.which("ffmpeg"),
        paddleocr_import=importlib.util.find_spec("paddleocr") is not None,
        paddle_import=importlib.util.find_spec("paddle") is not None,
        pillow_import=importlib.util.find_spec("PIL") is not None,
        setuptools_import=importlib.util.find_spec("setuptools") is not None,
    )


def dependency_help_text() -> str:
    requirements_path = SKILL_DIR / "requirements-video.txt"
    return "\n".join(
        [
            "Missing video-ingest dependencies.",
            f"- Python packages: python3 -m pip install -r {requirements_path}",
            "- macOS/Homebrew: brew install yt-dlp ffmpeg",
            "- Debian/Ubuntu: sudo apt-get install ffmpeg",
            "- If the distro yt-dlp package is too old, install a current standalone release or use Homebrew/Linuxbrew.",
            "- Install the Paddle runtime that matches the actual environment. Use a CUDA build only when Paddle and the host both support GPU execution.",
        ]
    )


def detect_runtime_profile() -> RuntimeProfile:
    system = platform.system()
    release = platform.release()
    machine = platform.machine()
    is_wsl = system == "Linux" and "microsoft" in release.lower()
    total_memory_gb = detect_total_memory_gb()
    gpu_available = detect_gpu_available()
    recommended_ocr_batch_size = recommend_ocr_batch_size(
        is_wsl=is_wsl,
        total_memory_gb=total_memory_gb,
        gpu_available=gpu_available,
    )
    recommended_frame_interval_sec = recommend_frame_interval_sec(
        is_wsl=is_wsl,
        total_memory_gb=total_memory_gb,
        gpu_available=gpu_available,
    )
    return RuntimeProfile(
        platform_system=system,
        platform_release=release,
        machine=machine,
        is_wsl=is_wsl,
        total_memory_gb=total_memory_gb,
        gpu_requested=False,
        gpu_available=gpu_available,
        recommended_ocr_batch_size=recommended_ocr_batch_size,
        recommended_frame_interval_sec=recommended_frame_interval_sec,
    )


def detect_total_memory_gb() -> float | None:
    meminfo = Path("/proc/meminfo")
    if meminfo.exists():
        for line in meminfo.read_text(encoding="utf-8", errors="ignore").splitlines():
            if line.startswith("MemTotal:"):
                parts = line.split()
                if len(parts) >= 2:
                    kib = int(parts[1])
                    return round(kib / (1024 * 1024), 1)
    return None


def detect_gpu_available() -> bool:
    try:
        import paddle

        if hasattr(paddle, "device") and paddle.device.is_compiled_with_cuda():
            return True
    except Exception:
        pass
    if shutil.which("nvidia-smi"):
        try:
            proc = run_command(
                ["nvidia-smi", "--query-gpu=name", "--format=csv,noheader"],
                check=False,
            )
            return proc.returncode == 0 and bool(proc.stdout.strip())
        except Exception:
            return False
    return False


def recommend_ocr_batch_size(
    *,
    is_wsl: bool,
    total_memory_gb: float | None,
    gpu_available: bool,
) -> int:
    if gpu_available:
        return 4
    if is_wsl:
        return 1
    if total_memory_gb is not None and total_memory_gb < 8:
        return 1
    if total_memory_gb is not None and total_memory_gb < 16:
        return 2
    return 2


def recommend_frame_interval_sec(
    *,
    is_wsl: bool,
    total_memory_gb: float | None,
    gpu_available: bool,
) -> int:
    if gpu_available:
        return 20
    if is_wsl:
        return 45
    if total_memory_gb is not None and total_memory_gb < 8:
        return 45
    return DEFAULT_FRAME_INTERVAL_SEC


def resolve_job_dir(
    *,
    job_dir: str | None,
    scratch_root: str | None,
    source_slug: str,
    mode: str,
) -> Path:
    if job_dir:
        resolved = Path(job_dir).expanduser().resolve()
        resolved.mkdir(parents=True, exist_ok=True)
        return resolved
    if scratch_root:
        root = Path(scratch_root).expanduser().resolve()
        root.mkdir(parents=True, exist_ok=True)
        timestamp = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
        resolved = root / f"{source_slug}-{mode}-{timestamp}"
        resolved.mkdir(parents=True, exist_ok=True)
        return resolved
    return Path(
        tempfile.mkdtemp(prefix=f"kb-video-{source_slug}-{mode}-")
    ).resolve()


def build_job_metadata(args: Any, *, job_dir: Path) -> dict[str, Any]:
    source_slug = slugify(args.source_slug or args.video_url or args.source_url)
    runtime_profile = detect_runtime_profile()
    return {
        "job_id": job_dir.name,
        "created_at": utc_now(),
        "mode": args.mode,
        "source_url": args.source_url,
        "video_url": args.video_url or args.source_url,
        "source_slug": source_slug,
        "scratch_root": str(job_dir.parent),
        "paths": {
            "downloads_dir": str(job_dir / "downloads"),
            "frames_dir": str(job_dir / "frames"),
            "prompts_dir": str(job_dir / "prompts"),
            "outputs_dir": str(job_dir / "outputs"),
            "artifacts_dir": str(job_dir / "artifacts"),
        },
        "settings": {
            "frame_interval_sec": args.frame_interval_sec,
            "segment_length_sec": args.segment_length_sec,
            "max_frames": args.max_frames,
            "ocr_batch_size": args.ocr_batch_size,
            "ocr_lang": args.ocr_lang,
            "transcribe_command": args.transcribe_command,
        },
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
        "workspace_relevance": args.workspace_relevance or "",
    }


def load_or_create_job_metadata(args: Any) -> tuple[Path, dict[str, Any]]:
    source_slug = slugify(args.source_slug or args.video_url or args.source_url)
    job_dir = resolve_job_dir(
        job_dir=args.job_dir,
        scratch_root=args.scratch_root,
        source_slug=source_slug,
        mode=args.mode,
    )
    metadata_path = job_dir / "job.json"
    if metadata_path.exists():
        metadata = read_json(metadata_path)
    else:
        metadata = build_job_metadata(args, job_dir=job_dir)
        write_json(metadata_path, metadata)
    for path_name in metadata["paths"].values():
        Path(path_name).mkdir(parents=True, exist_ok=True)
    return job_dir, metadata


def fetch_video_metadata(job_dir: Path, video_url: str) -> dict[str, Any]:
    proc = run_command(["yt-dlp", "--dump-single-json", "--no-warnings", video_url])
    payload = json.loads(proc.stdout)
    write_json(job_dir / "artifacts" / "video-metadata.json", payload)
    return payload


def fetch_subtitles(job_dir: Path, video_url: str) -> list[Path]:
    template = str(job_dir / "downloads" / "source.%(ext)s")
    proc = run_command(
        [
            "yt-dlp",
            "--skip-download",
            "--no-warnings",
            "--write-subs",
            "--write-auto-subs",
            "--sub-langs",
            SUBTITLE_LANGS,
            "--sub-format",
            "vtt/best",
            "--output",
            template,
            video_url,
        ],
        check=False,
    )
    subtitle_attempt = {
        "video_url": video_url,
        "returncode": proc.returncode,
        "stdout": proc.stdout,
        "stderr": proc.stderr,
    }
    write_json(job_dir / "artifacts" / "subtitle-fetch.json", subtitle_attempt)
    subtitle_paths = sorted((job_dir / "downloads").glob("*.vtt"))
    return subtitle_paths


def download_video(job_dir: Path, video_url: str) -> Path:
    downloads_dir = job_dir / "downloads"
    existing = [
        path
        for path in downloads_dir.iterdir()
        if path.is_file() and path.suffix.lower() not in {".vtt", ".json", ".part"}
    ]
    if existing:
        return sorted(existing)[0]
    template = str(downloads_dir / "source.%(ext)s")
    run_command(
        [
            "yt-dlp",
            "--no-warnings",
            "-f",
            "best[height<=480][ext=mp4]/best[height<=480]/best[ext=mp4]/best",
            "--output",
            template,
            video_url,
        ]
    )
    existing = [
        path
        for path in downloads_dir.iterdir()
        if path.is_file() and path.suffix.lower() not in {".vtt", ".json", ".part"}
    ]
    if not existing:
        raise RuntimeError("yt-dlp finished but no video file was downloaded.")
    return sorted(existing)[0]


TIMESTAMP_RE = re.compile(
    r"(?P<start>\d{2}:\d{2}:\d{2}\.\d{3}|\d{2}:\d{2}\.\d{3})\s+-->\s+"
    r"(?P<end>\d{2}:\d{2}:\d{2}\.\d{3}|\d{2}:\d{2}\.\d{3})"
)


def parse_vtt_timestamp(value: str) -> float:
    parts = value.split(":")
    if len(parts) == 2:
        minutes, rest = parts
        seconds, millis = rest.split(".")
        return int(minutes) * 60 + int(seconds) + int(millis) / 1000
    hours, minutes, rest = parts
    seconds, millis = rest.split(".")
    return int(hours) * 3600 + int(minutes) * 60 + int(seconds) + int(millis) / 1000


def parse_vtt_file(path: Path) -> list[dict[str, Any]]:
    lines = path.read_text(encoding="utf-8", errors="ignore").splitlines()
    cues: list[dict[str, Any]] = []
    current_time: tuple[float, float] | None = None
    buffer: list[str] = []
    for line in lines:
        stripped = line.strip()
        match = TIMESTAMP_RE.match(stripped)
        if match:
            if current_time and buffer:
                text = " ".join(buffer).strip()
                if text:
                    cues.append(
                        {
                            "start": current_time[0],
                            "end": current_time[1],
                            "text": cleanup_caption_text(text),
                        }
                    )
            current_time = (
                parse_vtt_timestamp(match.group("start")),
                parse_vtt_timestamp(match.group("end")),
            )
            buffer = []
            continue
        if not stripped or stripped == "WEBVTT" or stripped.startswith("Kind:") or stripped.startswith("Language:"):
            continue
        if "-->" in stripped:
            continue
        buffer.append(stripped)
    if current_time and buffer:
        text = cleanup_caption_text(" ".join(buffer))
        if text:
            cues.append({"start": current_time[0], "end": current_time[1], "text": text})
    return cues


def cleanup_caption_text(text: str) -> str:
    cleaned = re.sub(r"<[^>]+>", " ", text)
    cleaned = cleaned.replace("&nbsp;", " ").replace("&amp;", "&")
    cleaned = re.sub(r"\s+", " ", cleaned)
    return cleaned.strip()


def load_or_create_transcript(
    job_dir: Path,
    *,
    subtitle_paths: list[Path],
    video_path: Path,
    transcribe_command: str | None,
) -> list[dict[str, Any]]:
    transcript_path = job_dir / "artifacts" / "normalized-transcript.json"
    if transcript_path.exists():
        return read_json(transcript_path)
    cues: list[dict[str, Any]] = []
    if subtitle_paths:
        cues = parse_vtt_file(subtitle_paths[0])
    elif transcribe_command:
        cues = run_external_transcription(job_dir, video_path, transcribe_command)
    else:
        raise RuntimeError(
            "No subtitles were found and no transcription fallback command was provided. "
            "Pass --transcribe-command with a command that prints VTT to stdout."
        )
    write_json(transcript_path, cues)
    return cues


def run_external_transcription(
    job_dir: Path,
    video_path: Path,
    transcribe_command: str,
) -> list[dict[str, Any]]:
    command_text = (
        transcribe_command.replace("{video_path}", str(video_path)).replace("{job_dir}", str(job_dir))
    )
    proc = subprocess.run(
        command_text,
        shell=True,
        check=True,
        capture_output=True,
        text=True,
    )
    vtt_output = proc.stdout.strip()
    if not vtt_output:
        raise RuntimeError("Transcription command completed but produced no VTT output.")
    transcript_vtt_path = job_dir / "artifacts" / "transcription-fallback.vtt"
    transcript_vtt_path.write_text(vtt_output + "\n", encoding="utf-8")
    return parse_vtt_file(transcript_vtt_path)


def extract_frames(
    *,
    job_dir: Path,
    video_path: Path,
    frame_interval_sec: int,
    max_frames: int,
) -> list[dict[str, Any]]:
    frames_dir = job_dir / "frames" / "raw"
    frames_dir.mkdir(parents=True, exist_ok=True)
    existing = sorted(frames_dir.glob("frame-*.jpg"))
    if not existing:
        fps_filter = f"fps=1/{max(1, frame_interval_sec)},scale='min(960,iw)':-2"
        run_command(
            [
                "ffmpeg",
                "-hide_banner",
                "-loglevel",
                "error",
                "-i",
                str(video_path),
                "-vf",
                fps_filter,
                "-q:v",
                "3",
                str(frames_dir / "frame-%05d.jpg"),
            ]
        )
        existing = sorted(frames_dir.glob("frame-*.jpg"))
    selected = existing[: max(1, max_frames)]
    frames = []
    for index, path in enumerate(selected, start=1):
        timestamp = (index - 1) * frame_interval_sec
        frames.append(
            {
                "frame_index": index,
                "timestamp_sec": timestamp,
                "timestamp_label": format_seconds(timestamp),
                "path": str(path),
            }
        )
    write_json(job_dir / "artifacts" / "frames.json", frames)
    return frames


def dedupe_frames(job_dir: Path, frames: list[dict[str, Any]]) -> list[dict[str, Any]]:
    deduped_path = job_dir / "artifacts" / "frames-deduped.json"
    if deduped_path.exists():
        return read_json(deduped_path)
    from PIL import Image

    kept: list[dict[str, Any]] = []
    seen_hashes: set[str] = set()
    for frame in frames:
        frame_hash = compute_dhash(Path(frame["path"]), Image)
        if frame_hash in seen_hashes:
            continue
        seen_hashes.add(frame_hash)
        frame_copy = dict(frame)
        frame_copy["dhash"] = frame_hash
        kept.append(frame_copy)
    write_json(deduped_path, kept)
    return kept


def compute_dhash(path: Path, image_module: Any) -> str:
    image = image_module.open(path).convert("L").resize((9, 8))
    bits = []
    pixels = list(image.getdata())
    for row in range(8):
        row_start = row * 9
        for col in range(8):
            left = pixels[row_start + col]
            right = pixels[row_start + col + 1]
            bits.append("1" if left > right else "0")
    return f"{int(''.join(bits), 2):016x}"


def run_ocr(
    job_dir: Path,
    frames: list[dict[str, Any]],
    *,
    ocr_lang: str,
    ocr_batch_size: int,
) -> list[dict[str, Any]]:
    ocr_path = job_dir / "artifacts" / "frame-ocr.json"
    if ocr_path.exists():
        return read_json(ocr_path)
    paddle_home = job_dir / ".cache" / "paddleocr"
    paddle_home.mkdir(parents=True, exist_ok=True)
    os.environ.setdefault("PADDLE_OCR_BASE_DIR", str(paddle_home))
    from paddleocr import PaddleOCR

    batch_size = max(1, min(ocr_batch_size, MAX_RECOMMENDED_OCR_BATCH_SIZE))
    runtime_profile = detect_runtime_profile()
    ocr = PaddleOCR(
        use_angle_cls=True,
        lang=ocr_lang,
        use_gpu=runtime_profile.gpu_available,
        show_log=False,
    )
    results: list[dict[str, Any]] = []
    for offset in range(0, len(frames), batch_size):
        chunk = frames[offset : offset + batch_size]
        for frame in chunk:
            ocr_result = ocr.ocr(frame["path"], cls=True)
            lines: list[dict[str, Any]] = []
            for block in ocr_result or []:
                for line in block or []:
                    if len(line) < 2:
                        continue
                    text, confidence = line[1]
                    text = str(text).strip()
                    if not text:
                        continue
                    lines.append({"text": text, "confidence": float(confidence)})
            results.append(
                {
                    "path": frame["path"],
                    "timestamp_sec": frame["timestamp_sec"],
                    "timestamp_label": frame["timestamp_label"],
                    "line_count": len(lines),
                    "avg_confidence": round(
                        sum(line["confidence"] for line in lines) / len(lines), 4
                    )
                    if lines
                    else 0.0,
                    "lines": lines,
                }
            )
    write_json(ocr_path, results)
    return results


def build_segments(
    *,
    transcript: list[dict[str, Any]],
    deduped_frames: list[dict[str, Any]],
    ocr_results: list[dict[str, Any]],
    segment_length_sec: int,
) -> list[dict[str, Any]]:
    if not transcript:
        return []
    ocr_by_path = {entry["path"]: entry for entry in ocr_results}
    segments_map: dict[int, dict[str, Any]] = {}
    for cue in transcript:
        index = int(cue["start"] // segment_length_sec)
        segment = segments_map.setdefault(
            index,
            {
                "segment_index": index + 1,
                "start_sec": math.floor(cue["start"]),
                "end_sec": math.ceil(cue["end"]),
                "transcript_parts": [],
                "candidate_frames": [],
            },
        )
        segment["start_sec"] = min(segment["start_sec"], math.floor(cue["start"]))
        segment["end_sec"] = max(segment["end_sec"], math.ceil(cue["end"]))
        segment["transcript_parts"].append(cue["text"])
    for frame in deduped_frames:
        index = int(frame["timestamp_sec"] // segment_length_sec)
        if index not in segments_map:
            continue
        frame_copy = dict(frame)
        frame_copy["ocr"] = ocr_by_path.get(frame["path"], {})
        segments_map[index]["candidate_frames"].append(frame_copy)
    segments = []
    for index in sorted(segments_map):
        segment = segments_map[index]
        transcript_text = cleanup_caption_text(" ".join(segment.pop("transcript_parts")))
        named_entities = extract_named_entities(
            transcript_text,
            segment["candidate_frames"],
        )
        representative_frames = rank_representative_frames(segment["candidate_frames"])
        segments.append(
            {
                "segment_index": segment["segment_index"],
                "start_sec": segment["start_sec"],
                "end_sec": segment["end_sec"],
                "time_anchor_label": f"{format_seconds(segment['start_sec'])}-{format_seconds(segment['end_sec'])}",
                "transcript_text": transcript_text,
                "draft_summary": build_draft_summary(transcript_text, representative_frames),
                "named_entities": named_entities,
                "candidate_frames": segment["candidate_frames"],
                "representative_frames": representative_frames,
            }
        )
    return segments


def extract_named_entities(transcript_text: str, frames: list[dict[str, Any]]) -> list[str]:
    candidates: set[str] = set()
    text = transcript_text
    for frame in frames:
        ocr = frame.get("ocr", {})
        for line in ocr.get("lines", []):
            text += f" {line['text']}"
    for token in re.findall(r"\b[A-Z][A-Za-z0-9_.:/+-]{1,}\b", text):
        if token.lower() in {"the", "this", "that"}:
            continue
        candidates.add(token)
    return sorted(candidates)[:12]


def build_draft_summary(transcript_text: str, representative_frames: list[dict[str, Any]]) -> str:
    compact = transcript_text.strip()
    if len(compact) > 280:
        compact = compact[:277].rstrip() + "..."
    ocr_terms = []
    for frame in representative_frames:
        for line in frame.get("ocr", {}).get("lines", [])[:3]:
            ocr_terms.append(line["text"])
    if ocr_terms:
        return f"{compact} OCR anchors: {' | '.join(ocr_terms[:3])}"
    return compact


def rank_representative_frames(frames: list[dict[str, Any]]) -> list[dict[str, Any]]:
    ranked = []
    for frame in frames:
        ocr = frame.get("ocr", {})
        line_count = ocr.get("line_count", 0)
        confidence = ocr.get("avg_confidence", 0.0)
        ranked.append((line_count + confidence, frame))
    ranked.sort(key=lambda item: item[0], reverse=True)
    selected = []
    seen_paths: set[str] = set()
    for _, frame in ranked:
        if frame["path"] in seen_paths:
            continue
        seen_paths.add(frame["path"])
        selected.append(frame)
        if len(selected) >= 2:
            break
    return selected


def build_prompt_payloads(job_dir: Path, metadata: dict[str, Any], segments: list[dict[str, Any]]) -> None:
    prompts_dir = Path(metadata["paths"]["prompts_dir"])
    prompts_dir.mkdir(parents=True, exist_ok=True)
    segment_payload = {
        "contract": load_prompt_contract("segment-summary.md"),
        "source": {
            "source_url": metadata["source_url"],
            "video_url": metadata["video_url"],
            "workspace_relevance": metadata.get("workspace_relevance", ""),
        },
        "segments": [
            {
                "segment_index": segment["segment_index"],
                "time_anchor_label": segment["time_anchor_label"],
                "transcript_text": segment["transcript_text"],
                "ocr_lines": [
                    line["text"]
                    for frame in segment["candidate_frames"]
                    for line in frame.get("ocr", {}).get("lines", [])[:3]
                ],
            }
            for segment in segments
        ],
    }
    frame_payload = {
        "contract": load_prompt_contract("representative-frame-selection.md"),
        "segments": [
            {
                "segment_index": segment["segment_index"],
                "time_anchor_label": segment["time_anchor_label"],
                "draft_summary": segment["draft_summary"],
                "candidate_frames": [
                    {
                        "path": relativize_to_repo(Path(frame["path"])),
                        "timestamp_label": frame["timestamp_label"],
                        "ocr_lines": [line["text"] for line in frame.get("ocr", {}).get("lines", [])[:5]],
                    }
                    for frame in segment["candidate_frames"]
                ],
            }
            for segment in segments
        ],
    }
    compile_payload = {
        "contract": load_prompt_contract("kb-compile-review.md"),
        "source": {
            "source_url": metadata["source_url"],
            "video_url": metadata["video_url"],
            "source_slug": metadata["source_slug"],
            "workspace_relevance": metadata.get("workspace_relevance", ""),
        },
        "segments": [
            {
                "segment_index": segment["segment_index"],
                "time_anchor_label": segment["time_anchor_label"],
                "draft_summary": segment["draft_summary"],
                "named_entities": segment["named_entities"],
                "representative_frames": [
                    {
                        "path": relativize_to_repo(Path(frame["path"])),
                        "timestamp_label": frame["timestamp_label"],
                    }
                    for frame in segment["representative_frames"]
                ],
            }
            for segment in segments
        ],
    }
    write_json(prompts_dir / "segment-summary-payload.json", segment_payload)
    write_json(prompts_dir / "representative-frame-selection-payload.json", frame_payload)
    write_json(prompts_dir / "kb-compile-review-payload.json", compile_payload)


def relativize_to_repo(path: Path) -> str:
    try:
        return str(path.resolve().relative_to(ROOT_DIR))
    except ValueError:
        return str(path.resolve())


def render_skim_markdown(
    *,
    metadata: dict[str, Any],
    video_metadata: dict[str, Any],
    segments: list[dict[str, Any]],
) -> str:
    lines = [
        "# Video KB Skim",
        "",
        f"- Source URL: {metadata['source_url']}",
        f"- Canonical video URL: {metadata['video_url']}",
        f"- Title: {video_metadata.get('title', 'unknown')}",
        f"- Channel: {video_metadata.get('uploader') or video_metadata.get('channel') or 'unknown'}",
        f"- Duration: {format_seconds(video_metadata.get('duration'))}",
        f"- Workspace relevance: {metadata.get('workspace_relevance', '') or '(not provided)'}",
        "",
        "## Segment summaries",
        "",
    ]
    for segment in segments:
        lines.append(f"### Segment {segment['segment_index']} ({segment['time_anchor_label']})")
        lines.append("")
        lines.append(segment["draft_summary"] or "(no summary)")
        lines.append("")
        if segment["named_entities"]:
            lines.append(f"- Named entities: {', '.join(segment['named_entities'])}")
        if segment["representative_frames"]:
            lines.append("- Representative frames:")
            for frame in segment["representative_frames"]:
                lines.append(
                    f"  - `{relativize_to_repo(Path(frame['path']))}` at {frame['timestamp_label']}"
                )
        lines.append("")
    return "\n".join(lines).rstrip() + "\n"


def render_source_note_draft(
    *,
    metadata: dict[str, Any],
    video_metadata: dict[str, Any],
    segments: list[dict[str, Any]],
) -> str:
    ingested_on = datetime.now().date().isoformat()
    ingested_year = ingested_on[:4]
    source_type = "video" if metadata["source_url"] == metadata["video_url"] else "video_backed_article"
    summary_bullets = [segment["draft_summary"] for segment in segments[:4] if segment["draft_summary"]]
    selected_screenshots = build_durable_screenshot_paths(
        metadata=metadata,
        segments=segments,
        ingested_year=ingested_year,
    )
    time_anchors = [segment["time_anchor_label"] for segment in segments[:6]]
    named_entities = sorted(
        {
            entity
            for segment in segments
            for entity in segment.get("named_entities", [])
            if entity.lower() not in COMMON_ENTITY_STOPWORDS
        }
    )[:20]
    lines = [
        "---",
        f"title: {json.dumps(video_metadata.get('title', metadata['source_slug']))}",
        f"source_url: {metadata['source_url']}",
        f"source_type: {source_type}",
        f"video_url: {metadata['video_url']}",
        f"ingested_on: {ingested_on}",
        "status: active",
        "tags: []",
        "related_pages: []",
        f"channel: {json.dumps(video_metadata.get('uploader') or video_metadata.get('channel') or '')}",
        f"duration: {json.dumps(format_seconds(video_metadata.get('duration')))}",
        f"time_anchors: {json.dumps(time_anchors, ensure_ascii=False)}",
        f"selected_screenshots: {json.dumps(selected_screenshots[:8], ensure_ascii=False)}",
        f"named_entities: {json.dumps(named_entities, ensure_ascii=False)}",
        "---",
        "",
        "## Workspace relevance",
        "",
        metadata.get("workspace_relevance", "") or "- TODO",
        "",
        "## Summary",
        "",
    ]
    for bullet in summary_bullets:
        lines.append(f"- {bullet}")
    lines.extend(["", "## Follow-up", "", "- TODO", ""])
    return "\n".join(lines)


def build_durable_screenshot_paths(
    *,
    metadata: dict[str, Any],
    segments: list[dict[str, Any]],
    ingested_year: str,
) -> list[str]:
    durable_root = Path("docs/kb/assets/source-images") / ingested_year / metadata["source_slug"]
    durable_paths = []
    for segment in segments:
        for frame in segment["representative_frames"]:
            frame_name = Path(frame["path"]).name
            durable_paths.append(str(durable_root / frame_name))
    return durable_paths


def render_wiki_update_draft(
    *,
    metadata: dict[str, Any],
    video_metadata: dict[str, Any],
    segments: list[dict[str, Any]],
) -> str:
    lines = [
        f"# Wiki update draft for {video_metadata.get('title', metadata['source_slug'])}",
        "",
        "## Candidate durable takeaways",
        "",
    ]
    for segment in segments[:6]:
        lines.append(f"- {segment['time_anchor_label']}: {segment['draft_summary']}")
    lines.extend(["", "## Candidate wiki targets", "", "- topics/TODO", "- tools/TODO", "- patterns/TODO", ""])
    return "\n".join(lines)


def render_log_entry_draft(*, metadata: dict[str, Any], video_metadata: dict[str, Any]) -> str:
    date_text = datetime.now().date().isoformat()
    return (
        f"- {date_text}: skimmed video-backed source "
        f"`{video_metadata.get('title', metadata['source_slug'])}` ({metadata['video_url']})\n"
    )


def write_output(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")


def run_skim_pipeline(args: Any) -> dict[str, Any]:
    status = ensure_dependency_status()
    if not status.ok:
        raise RuntimeError(dependency_help_text())
    job_dir, metadata = load_or_create_job_metadata(args)
    video_metadata = fetch_video_metadata(job_dir, metadata["video_url"])
    subtitle_paths = fetch_subtitles(job_dir, metadata["video_url"])
    video_path = download_video(job_dir, metadata["video_url"])
    transcript = load_or_create_transcript(
        job_dir,
        subtitle_paths=subtitle_paths,
        video_path=video_path,
        transcribe_command=args.transcribe_command,
    )
    frames = extract_frames(
        job_dir=job_dir,
        video_path=video_path,
        frame_interval_sec=args.frame_interval_sec,
        max_frames=args.max_frames,
    )
    deduped_frames = dedupe_frames(job_dir, frames)
    ocr_results = run_ocr(
        job_dir,
        deduped_frames,
        ocr_lang=args.ocr_lang,
        ocr_batch_size=args.ocr_batch_size,
    )
    segments = build_segments(
        transcript=transcript,
        deduped_frames=deduped_frames,
        ocr_results=ocr_results,
        segment_length_sec=args.segment_length_sec,
    )
    write_json(job_dir / "artifacts" / "segments.json", segments)
    build_prompt_payloads(job_dir, metadata, segments)
    skim_markdown = render_skim_markdown(
        metadata=metadata,
        video_metadata=video_metadata,
        segments=segments,
    )
    skim_path = Path(metadata["paths"]["outputs_dir"]) / "skim.md"
    write_output(skim_path, skim_markdown)
    return {
        "job_dir": str(job_dir),
        "skim_path": str(skim_path),
        "segment_count": len(segments),
        "video_title": video_metadata.get("title", ""),
    }


def run_ingest_pipeline(args: Any) -> dict[str, Any]:
    job_dir, metadata = load_or_create_job_metadata(args)
    segments_path = job_dir / "artifacts" / "segments.json"
    video_metadata_path = job_dir / "artifacts" / "video-metadata.json"
    if not segments_path.exists() or not video_metadata_path.exists():
        run_skim_pipeline(args)
    metadata = read_json(job_dir / "job.json")
    segments = read_json(job_dir / "artifacts" / "segments.json")
    video_metadata = read_json(job_dir / "artifacts" / "video-metadata.json")
    outputs_dir = Path(metadata["paths"]["outputs_dir"])
    source_note_path = outputs_dir / "source-note-draft.md"
    wiki_update_path = outputs_dir / "wiki-update-draft.md"
    log_entry_path = outputs_dir / "log-entry-draft.md"
    write_output(
        source_note_path,
        render_source_note_draft(
            metadata=metadata,
            video_metadata=video_metadata,
            segments=segments,
        ),
    )
    write_output(
        wiki_update_path,
        render_wiki_update_draft(
            metadata=metadata,
            video_metadata=video_metadata,
            segments=segments,
        ),
    )
    write_output(
        log_entry_path,
        render_log_entry_draft(metadata=metadata, video_metadata=video_metadata),
    )
    return {
        "job_dir": str(job_dir),
        "source_note_path": str(source_note_path),
        "wiki_update_path": str(wiki_update_path),
        "log_entry_path": str(log_entry_path),
    }
