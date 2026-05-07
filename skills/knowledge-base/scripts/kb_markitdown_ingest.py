#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
import urllib.parse
import urllib.request
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path

DEFAULT_MARKITDOWN_WITH = "markitdown[pdf,docx,pptx,xlsx,epub]"
REPO_ROOT = Path(__file__).resolve().parents[3]
DEFAULT_UV_CACHE_DIR = REPO_ROOT / ".uv-cache"
URL_SCHEMES = {"http", "https", "file"}
DOWNLOAD_TIMEOUT_SECONDS = 60
DOWNLOAD_MAX_BYTES = 100 * 1024 * 1024  # 100 MB
DOWNLOAD_CHUNK_SIZE = 64 * 1024  # 64 KB


@dataclass
class DependencyStatus:
    ok: bool
    runner_kind: str | None
    version: str | None
    missing: list[str]
    help_text: str


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="Convert file-like KB sources into temporary Markdown drafting input."
    )
    parser.add_argument(
        "mode",
        nargs="?",
        choices=["convert"],
        help="Pipeline mode. Use convert for file-like KB sources.",
    )
    parser.add_argument(
        "source",
        nargs="?",
        help="Local path, file:// URI, or direct document URL to convert.",
    )
    parser.add_argument(
        "--workspace-relevance",
        help="Short note describing why this source matters to the workspace.",
    )
    parser.add_argument(
        "--source-slug",
        help="Optional slug override for the job directory name.",
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
        "--keep-data-uris",
        action="store_true",
        help="Pass through markitdown data URIs instead of truncating them.",
    )
    parser.add_argument(
        "--check-deps",
        action="store_true",
        help="Validate markitdown runner availability, then exit.",
    )
    return parser


def validate_args(args: argparse.Namespace) -> None:
    if args.check_deps:
        return
    if not args.mode or not args.source:
        raise SystemExit("mode and source are required unless --check-deps is used.")


def sanitize_slug(value: str) -> str:
    slug = re.sub(r"[^a-z0-9]+", "-", value.lower()).strip("-")
    return slug or "kb-source"


def is_supported_url(value: str) -> bool:
    parsed = urllib.parse.urlparse(value)
    return parsed.scheme in URL_SCHEMES


def default_slug_for_source(source: str) -> str:
    if is_supported_url(source):
        parsed = urllib.parse.urlparse(source)
        name = Path(parsed.path).stem or parsed.netloc or "document"
        return sanitize_slug(name)
    return sanitize_slug(Path(source).stem or "document")


def ensure_job_dir(args: argparse.Namespace) -> Path:
    if args.job_dir:
        job_dir = Path(args.job_dir).expanduser().resolve()
        job_dir.mkdir(parents=True, exist_ok=True)
        return job_dir
    if args.scratch_root:
        scratch_root = Path(args.scratch_root).expanduser().resolve()
        scratch_root.mkdir(parents=True, exist_ok=True)
        slug = sanitize_slug(args.source_slug or default_slug_for_source(args.source))
        stamp = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
        return Path(tempfile.mkdtemp(prefix=f"{stamp}-{slug}-", dir=scratch_root))
    return Path(
        tempfile.mkdtemp(
            prefix="kb-markitdown-",
            dir=None,
        )
    )


def command_exists(name: str) -> bool:
    return shutil.which(name) is not None


def markitdown_runner() -> tuple[str | None, list[str] | None]:
    if command_exists("markitdown"):
        return "binary", ["markitdown"]
    if command_exists("uv"):
        return "uv", ["uv", "run", "--with", DEFAULT_MARKITDOWN_WITH, "python", "-m", "markitdown"]
    return None, None


def run_command(command: list[str], *, extra_env: dict[str, str] | None = None) -> subprocess.CompletedProcess[str]:
    env = os.environ.copy()
    if extra_env:
        env.update(extra_env)
    return subprocess.run(
        command,
        capture_output=True,
        text=True,
        check=False,
        env=env,
    )


def dependency_help_text() -> str:
    return "\n".join(
        [
            "Missing MarkItDown runtime.",
            "Preferred setup:",
            f"  env UV_CACHE_DIR=.uv-cache uv run --with '{DEFAULT_MARKITDOWN_WITH}' python -m markitdown --help",
            "Alternative setup:",
            f"  python3 -m pip install '{DEFAULT_MARKITDOWN_WITH}'",
            "",
            "This fallback intentionally excludes OCR-heavy plugins from the default path.",
        ]
    )


def ensure_dependency_status() -> DependencyStatus:
    runner_kind, runner = markitdown_runner()
    if not runner_kind or not runner:
        missing = []
        if not command_exists("markitdown"):
            missing.append("markitdown")
        if not command_exists("uv"):
            missing.append("uv")
        return DependencyStatus(
            ok=False,
            runner_kind=None,
            version=None,
            missing=missing or ["markitdown runner"],
            help_text=dependency_help_text(),
        )

    extra_env: dict[str, str] = {}
    if runner_kind == "uv" and "UV_CACHE_DIR" not in os.environ:
        extra_env["UV_CACHE_DIR"] = str(DEFAULT_UV_CACHE_DIR)
        DEFAULT_UV_CACHE_DIR.mkdir(parents=True, exist_ok=True)

    result = run_command([*runner, "--version"], extra_env=extra_env)
    if result.returncode != 0:
        return DependencyStatus(
            ok=False,
            runner_kind=runner_kind,
            version=None,
            missing=["markitdown runtime"],
            help_text=f"{dependency_help_text()}\n\nmarkitdown version check failed:\n{result.stderr.strip()}",
        )

    return DependencyStatus(
        ok=True,
        runner_kind=runner_kind,
        version=result.stdout.strip(),
        missing=[],
        help_text="",
    )


def download_source(url: str, destination: Path) -> Path:
    destination.parent.mkdir(parents=True, exist_ok=True)
    request = urllib.request.Request(url, headers={"User-Agent": "kb-markitdown-ingest/1.0"})
    with urllib.request.urlopen(request, timeout=DOWNLOAD_TIMEOUT_SECONDS) as response:
        with destination.open("wb") as out_file:
            downloaded = 0
            while True:
                chunk = response.read(DOWNLOAD_CHUNK_SIZE)
                if not chunk:
                    break
                downloaded += len(chunk)
                if downloaded > DOWNLOAD_MAX_BYTES:
                    raise RuntimeError(
                        f"Download of {url!r} exceeded maximum allowed size of"
                        f" {DOWNLOAD_MAX_BYTES // (1024 * 1024)} MB"
                    )
                out_file.write(chunk)
    return destination


def resolve_source(source: str, job_dir: Path) -> dict[str, str]:
    inputs_dir = job_dir / "inputs"
    inputs_dir.mkdir(parents=True, exist_ok=True)

    if is_supported_url(source):
        parsed = urllib.parse.urlparse(source)
        if parsed.scheme == "file":
            local_path = Path(urllib.request.url2pathname(parsed.path)).expanduser().resolve()
            if not local_path.is_file():
                raise FileNotFoundError(f"local file does not exist: {local_path}")
            return {
                "source_kind": "file_uri",
                "source_url": source,
                "source_file": str(local_path),
                "input_path": str(local_path),
            }

        suffix = Path(parsed.path).suffix or ".bin"
        filename = sanitize_slug(Path(parsed.path).stem or parsed.netloc or "download") + suffix
        downloaded_path = download_source(source, inputs_dir / filename)
        return {
            "source_kind": "download_url",
            "source_url": source,
            "source_file": downloaded_path.name,
            "input_path": str(downloaded_path),
        }

    local_path = Path(source).expanduser().resolve()
    if not local_path.is_file():
        raise FileNotFoundError(f"local file does not exist: {local_path}")
    return {
        "source_kind": "local_path",
        "source_url": "",
        "source_file": str(local_path),
        "input_path": str(local_path),
    }


def build_source_context(metadata: dict[str, object]) -> str:
    lines = [
        "# KB MarkItDown Source Context",
        "",
        "Use this file as the drafting boundary for the durable KB note.",
        "",
        f"- Original source argument: `{metadata['source_argument']}`",
        f"- Source kind: `{metadata['source_kind']}`",
    ]
    if metadata.get("source_url"):
        lines.append(f"- Source URL: `{metadata['source_url']}`")
    lines.append(f"- Source file identity: `{metadata['source_file']}`")
    lines.append(f"- Conversion method: `{metadata['converter_version']}`")
    lines.append(f"- Converted Markdown path: `{metadata['converted_markdown']}`")
    if metadata.get("workspace_relevance"):
        lines.append(f"- Workspace relevance: {metadata['workspace_relevance']}")
    lines.extend(
        [
            "",
            "## Durable note requirements",
            "",
            "- Keep the curated source note in `docs/kb/sources/<year>/`; do not commit `converted.md`.",
            "- Preserve the original source URL when one exists.",
            "- When the source is local-only, record a durable file identity in `source_file` instead of depending on a private absolute path if that path is not meaningful later.",
            "- Add `conversion_method` and any real cleanup caveats when structure loss matters.",
            "- Stop and escalate if tables, ordering, or slide structure are too damaged to trust.",
        ]
    )
    return "\n".join(lines) + "\n"


def run_convert(args: argparse.Namespace) -> dict[str, object]:
    status = ensure_dependency_status()
    if not status.ok:
        raise RuntimeError(status.help_text)

    job_dir = ensure_job_dir(args)
    outputs_dir = job_dir / "outputs"
    outputs_dir.mkdir(parents=True, exist_ok=True)

    source_info = resolve_source(args.source, job_dir)
    converted_path = outputs_dir / "converted.md"
    context_path = outputs_dir / "source-context.md"
    metadata_path = job_dir / "metadata.json"

    runner_kind, runner = markitdown_runner()
    assert runner_kind and runner

    extra_env: dict[str, str] = {}
    if runner_kind == "uv" and "UV_CACHE_DIR" not in os.environ:
        extra_env["UV_CACHE_DIR"] = str(DEFAULT_UV_CACHE_DIR)
        DEFAULT_UV_CACHE_DIR.mkdir(parents=True, exist_ok=True)

    command = [*runner, source_info["input_path"], "-o", str(converted_path)]
    if args.keep_data_uris:
        command.append("--keep-data-uris")
    result = run_command(command, extra_env=extra_env)
    if result.returncode != 0:
        raise RuntimeError(
            "markitdown conversion failed.\n"
            f"command: {' '.join(command)}\n"
            f"stderr:\n{result.stderr.strip()}"
        )

    metadata: dict[str, object] = {
        "mode": "convert",
        "created_at": datetime.now(timezone.utc).isoformat(),
        "job_dir": str(job_dir),
        "source_argument": args.source,
        "source_kind": source_info["source_kind"],
        "source_url": source_info["source_url"],
        "source_file": source_info["source_file"],
        "input_path": source_info["input_path"],
        "workspace_relevance": args.workspace_relevance or "",
        "runner_kind": runner_kind,
        "converter_version": status.version,
        "converted_markdown": str(converted_path),
        "source_context": str(context_path),
    }
    metadata_path.write_text(json.dumps(metadata, indent=2) + "\n", encoding="utf-8")
    context_path.write_text(build_source_context(metadata), encoding="utf-8")
    return metadata


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    validate_args(args)

    if args.check_deps:
        status = ensure_dependency_status()
        print(json.dumps(asdict(status), indent=2))
        if not status.ok:
            print("", file=sys.stderr)
            print(status.help_text, file=sys.stderr)
            return 1
        return 0

    try:
        result = run_convert(args)
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        return 1

    print(json.dumps(result, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
