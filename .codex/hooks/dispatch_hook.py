#!/usr/bin/env python3

import json
import os
import subprocess
import sys
from pathlib import Path


def repo_root() -> Path:
    return Path(__file__).resolve().parents[2]


def workspace_roots(root: Path) -> list[Path]:
    roots = [root]
    if root.parent.name == ".worktrees":
        roots.append(root.parents[1])
    return roots


def detect_ai_arena_root(payload: dict) -> Path | None:
    cwd = payload.get("cwd")
    if cwd:
        try:
            top = subprocess.run(
                ["git", "-C", cwd, "rev-parse", "--show-toplevel"],
                check=True,
                capture_output=True,
                text=True,
            ).stdout.strip()
            top_path = Path(top)
            if top_path.name == "ai-arena" or top_path.name.startswith("ai-arena@"):
                return top_path
        except Exception:
            pass

    for root in workspace_roots(repo_root()):
        candidate = root / "ai-arena"
        if (candidate / "tools/codex-hook-post-tool-use.sh").exists():
            return candidate
    return None


def should_dispatch(mode: str, payload: dict, child_root: Path) -> bool:
    cwd = payload.get("cwd", "")
    try:
        cwd_path = Path(cwd).resolve()
        if cwd_path == child_root or child_root in cwd_path.parents:
            return True
    except Exception:
        pass

    if mode == "post-tool-use":
        if payload.get("tool_name") != "apply_patch":
            return False
        command = str(payload.get("tool_input", {}).get("command", ""))
        return "ai-arena/" in command or str(child_root) in command

    return True


def child_script(mode: str, child_root: Path) -> Path:
    name = "codex-hook-post-tool-use.sh" if mode == "post-tool-use" else "codex-hook-stop.sh"
    return child_root / "tools" / name


def main() -> int:
    mode = sys.argv[1]
    payload = json.load(sys.stdin)
    child_root = detect_ai_arena_root(payload)
    if child_root is None:
        return 0
    if not should_dispatch(mode, payload, child_root):
        return 0

    script = child_script(mode, child_root)
    if not script.exists():
        return 0

    proc = subprocess.run([str(script)], input=json.dumps(payload), text=True)
    return proc.returncode


if __name__ == "__main__":
    raise SystemExit(main())
