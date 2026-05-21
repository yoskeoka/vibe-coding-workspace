#!/usr/bin/env python3

import json
import os
import re
import subprocess
import sys
from pathlib import Path


POST_TOOL_NAMES = {"apply_patch", "Edit", "Write"}
SLOPLESS_VERSION = "0.2.10"
SLOPLESS_SHA = "c40c40f3127d0c61cbfc1c34cacf0a5f49ed7e26"
SLOPLESS_PACKAGE = f"slopless@{SLOPLESS_VERSION}"
JAPANESE_TEXT_RE = re.compile(r"[\u3040-\u30ff\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff]")
PATCH_PATH_RE = re.compile(r"^\*\*\* (?:Add|Delete|Update) File: (.+\.md)$", re.MULTILINE)
MOVE_PATH_RE = re.compile(r"^\*\*\* Move to: (.+\.md)$", re.MULTILINE)
PATH_KEYS = {
    "file_path",
    "path",
    "new_path",
    "old_path",
    "target_path",
    "destination_path",
    "source_path",
}


def repo_root() -> Path:
    return Path(__file__).resolve().parents[2]


def normalize_path(candidate: str, cwd: Path, root: Path) -> Path | None:
    raw = candidate.strip()
    if not raw.endswith(".md"):
        return None
    path = Path(raw)
    if not path.is_absolute():
        path = (cwd / path).resolve()
    else:
        path = path.resolve()
    try:
        path.relative_to(root)
    except ValueError:
        return None
    return path


def collect_paths(value: object, cwd: Path, root: Path, key: str | None = None) -> set[Path]:
    results: set[Path] = set()
    if isinstance(value, dict):
        for child_key, child_value in value.items():
            results.update(collect_paths(child_value, cwd, root, child_key))
        return results
    if isinstance(value, list):
        for item in value:
            results.update(collect_paths(item, cwd, root, key))
        return results
    if not isinstance(value, str):
        return results

    if key in PATH_KEYS:
        path = normalize_path(value, cwd, root)
        if path is not None:
            results.add(path)

    for match in PATCH_PATH_RE.findall(value):
        path = normalize_path(match, cwd, root)
        if path is not None:
            results.add(path)
    for match in MOVE_PATH_RE.findall(value):
        path = normalize_path(match, cwd, root)
        if path is not None:
            results.add(path)
    return results


def is_english_markdown(path: Path) -> bool:
    try:
        content = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        content = path.read_text(encoding="utf-8", errors="ignore")
    return JAPANESE_TEXT_RE.search(content) is None


def slopless_env() -> dict[str, str]:
    env = os.environ.copy()
    env.setdefault("npm_config_cache", "/tmp/vibe-coding-workspace-slopless-npm-cache")
    env.setdefault("npm_config_update_notifier", "false")
    env.setdefault("npm_config_fund", "false")
    env.setdefault("npm_config_audit", "false")
    return env


def sha_stamp_path() -> Path:
    return Path(f"/tmp/vibe-coding-workspace-slopless-{SLOPLESS_VERSION}-{SLOPLESS_SHA}.stamp")


def ensure_pinned_git_head(env: dict[str, str]) -> int:
    stamp = sha_stamp_path()
    if stamp.exists():
        return 0

    proc = subprocess.run(
        ["npm", "view", SLOPLESS_PACKAGE, "gitHead"],
        capture_output=True,
        text=True,
        env=env,
    )
    if proc.returncode != 0:
        sys.stderr.write(f"failed to resolve gitHead for {SLOPLESS_PACKAGE}\n")
        if proc.stdout:
            sys.stderr.write(proc.stdout)
        if proc.stderr:
            sys.stderr.write(proc.stderr)
        return proc.returncode

    actual = proc.stdout.strip()
    if actual != SLOPLESS_SHA:
        sys.stderr.write(
            f"slopless gitHead mismatch for {SLOPLESS_PACKAGE}: expected {SLOPLESS_SHA}, got {actual}\n"
        )
        return 2

    stamp.write_text(actual, encoding="utf-8")
    return 0


def run_slopless(path: Path) -> int:
    env = slopless_env()
    verify = ensure_pinned_git_head(env)
    if verify != 0:
        return verify

    proc = subprocess.run(
        [
            "npx",
            "--yes",
            f"--package={SLOPLESS_PACKAGE}",
            "slopless",
            str(path),
        ],
        capture_output=True,
        text=True,
        env=env,
    )
    if proc.returncode != 0:
        sys.stderr.write(f"slopless failed for {path.relative_to(repo_root())}\n")
        if proc.stdout:
            sys.stderr.write(proc.stdout)
        if proc.stderr:
            sys.stderr.write(proc.stderr)
    return proc.returncode


def main() -> int:
    payload = json.load(sys.stdin)
    if str(payload.get("tool_name", "")) not in POST_TOOL_NAMES:
        return 0

    root = repo_root()
    cwd = Path(payload.get("cwd") or root).resolve()
    candidates = sorted(collect_paths(payload.get("tool_input", {}), cwd, root))
    if not candidates:
        return 0

    for path in candidates:
        if not path.exists() or not path.is_file():
            continue
        if not is_english_markdown(path):
            continue
        result = run_slopless(path)
        if result != 0:
            return result
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
