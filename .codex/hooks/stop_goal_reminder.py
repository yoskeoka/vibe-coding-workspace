#!/usr/bin/env python3
"""Emit one workflow-completion reminder for an eligible Codex Stop event."""

import json
import subprocess
import sys
from pathlib import Path


HOOK_ROOT = Path(__file__).resolve().parents[2]
ELIGIBLE_PREFIXES = ("plan/", "feat/", "fix/")
MAX_SUMMARY_CHARS = 320


def git_output(cwd: Path, *args: str) -> str | None:
    try:
        result = subprocess.run(
            ["git", "-C", str(cwd), *args],
            check=True,
            capture_output=True,
            text=True,
        )
    except (OSError, subprocess.SubprocessError):
        return None
    return result.stdout.strip() or None


def git_top_level(cwd: Path) -> Path | None:
    output = git_output(cwd, "rev-parse", "--show-toplevel")
    return Path(output).resolve() if output else None


def compact(text: str) -> str:
    return " ".join(text.split())[:MAX_SUMMARY_CHARS].rstrip()


def matching_plan_summary(workspace_root: Path, branch: str) -> str | None:
    description = branch.split("/", 1)[1]
    plans = sorted((workspace_root / "docs/exec-plan/todo").glob(f"*-{description}.md"))
    if not plans:
        return None

    try:
        lines = plans[0].read_text(encoding="utf-8").splitlines()
    except OSError:
        return None

    heading = next((line[2:].strip() for line in lines if line.startswith("# ")), "")
    detail = ""
    for index, line in enumerate(lines):
        if line in ("## Objective", "## Completion boundary"):
            following = []
            for candidate in lines[index + 1 :]:
                if candidate.startswith("#"):
                    break
                if candidate.strip():
                    following.append(candidate.strip())
            detail = compact(" ".join(following))
            if detail:
                break
        if line.startswith("Completion boundary:"):
            detail = compact(line)
            break

    summary = compact(". ".join(part for part in (heading, detail) if part))
    return summary or None


def reminder(branch: str, summary: str | None) -> str:
    if branch.startswith("plan/"):
        boundary = "the reviewable plan PR and its initial latest-head follow-up"
    else:
        boundary = "verification, PR creation, and review-task's latest-head stop condition"
    plan_context = f" Matching plan: {summary}." if summary else ""
    return (
        "Continue unless you can affirm the user-requested goal is complete, "
        f"identify a genuine blocker, or require a user decision. Confirm {boundary}."
        f"{plan_context}"
    )


def evaluate(payload: dict[object, object]) -> dict[str, str] | None:
    if payload.get("stop_hook_active") is True:
        return None
    cwd = payload.get("cwd")
    if not isinstance(cwd, str) or not cwd:
        return None

    workspace_root = git_top_level(HOOK_ROOT)
    session_root = git_top_level(Path(cwd))
    if workspace_root is None or session_root is None or session_root != workspace_root:
        return None

    branch = git_output(workspace_root, "branch", "--show-current")
    if branch is None or not branch.startswith(ELIGIBLE_PREFIXES):
        return None
    return {"decision": "block", "reason": reminder(branch, matching_plan_summary(workspace_root, branch))}


def main() -> int:
    try:
        payload = json.load(sys.stdin)
    except (json.JSONDecodeError, OSError):
        sys.stderr.write("stop goal reminder: ignoring malformed input\n")
        return 0
    if not isinstance(payload, dict):
        sys.stderr.write("stop goal reminder: ignoring malformed input\n")
        return 0

    result = evaluate(payload)
    if result is not None:
        print(json.dumps(result))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
