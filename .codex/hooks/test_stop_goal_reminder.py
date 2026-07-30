#!/usr/bin/env python3
"""Focused tests for the workspace-only Codex Stop reminder."""

import contextlib
import importlib.util
import io
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch


SCRIPT = Path(__file__).with_name("stop_goal_reminder.py")
SPEC = importlib.util.spec_from_file_location("stop_goal_reminder", SCRIPT)
assert SPEC and SPEC.loader
reminder = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(reminder)

DISPATCH_SCRIPT = Path(__file__).with_name("dispatch_hook.py")
DISPATCH_SPEC = importlib.util.spec_from_file_location("dispatch_hook", DISPATCH_SCRIPT)
assert DISPATCH_SPEC and DISPATCH_SPEC.loader
dispatch = importlib.util.module_from_spec(DISPATCH_SPEC)
DISPATCH_SPEC.loader.exec_module(dispatch)


def git(cwd: Path, *args: str) -> None:
    subprocess.run(["git", "-C", str(cwd), *args], check=True, capture_output=True, text=True)


class StopGoalReminderTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory()
        self.workspace = Path(self.temp.name) / "workspace"
        self.workspace.mkdir()
        git(self.workspace, "init", "-q")
        git(self.workspace, "config", "user.email", "test@example.invalid")
        git(self.workspace, "config", "user.name", "Test")
        (self.workspace / ".gitkeep").write_text("", encoding="utf-8")
        git(self.workspace, "add", ".gitkeep")
        git(self.workspace, "commit", "-qm", "initial")

    def tearDown(self) -> None:
        self.temp.cleanup()

    def branch(self, name: str) -> None:
        git(self.workspace, "checkout", "-qb", name)

    def plan(self, slug: str, text: str = "Prevent premature handoff.") -> None:
        directory = self.workspace / "docs/exec-plan/todo"
        directory.mkdir(parents=True, exist_ok=True)
        (directory / f"0047-{slug}.md").write_text(
            f"# Stop Reminder\n\n## Objective\n\n{text}\n", encoding="utf-8"
        )

    def evaluate(self, payload: dict) -> dict | None:
        with patch.object(reminder, "HOOK_ROOT", self.workspace):
            return reminder.evaluate(payload)

    def test_first_stop_blocks_with_matching_plan_summary(self) -> None:
        self.branch("feat/stop-goal-reminder")
        self.plan("stop-goal-reminder")
        result = self.evaluate({"cwd": str(self.workspace)})
        self.assertEqual("block", result["decision"])
        self.assertIn("Stop Reminder", result["reason"])
        self.assertIn("Prevent premature handoff", result["reason"])
        self.assertIn("review-task", result["reason"])

    def test_plan_branch_uses_plan_pr_boundary(self) -> None:
        self.branch("plan/stop-goal-reminder")
        result = self.evaluate({"cwd": str(self.workspace)})
        self.assertIn("reviewable plan PR", result["reason"])

    def test_missing_plan_uses_generic_reminder(self) -> None:
        self.branch("fix/no-plan")
        result = self.evaluate({"cwd": str(self.workspace)})
        self.assertEqual("block", result["decision"])
        self.assertNotIn("Matching plan:", result["reason"])

    def test_ambiguous_or_unsafe_plan_suffix_has_no_summary(self) -> None:
        self.branch("feat/ambiguous")
        self.plan("ambiguous")
        second = self.workspace / "docs/exec-plan/todo/0048-ambiguous.md"
        second.write_text("# Other Plan\n", encoding="utf-8")
        self.assertIsNone(reminder.matching_plan_summary(self.workspace, "feat/ambiguous"))
        self.assertIsNone(reminder.matching_plan_summary(self.workspace, "feat/"))
        self.assertIsNone(reminder.matching_plan_summary(self.workspace, "feat/foo/bar"))

    def test_second_stop_passes_through(self) -> None:
        self.branch("feat/stop-goal-reminder")
        self.assertIsNone(self.evaluate({"cwd": str(self.workspace), "stop_hook_active": True}))
        self.assertIsNone(self.evaluate({"cwd": str(self.workspace), "stop_hook_active": "true"}))

    def test_child_repository_is_out_of_scope(self) -> None:
        self.branch("feat/stop-goal-reminder")
        child = self.workspace / "child"
        child.mkdir()
        git(child, "init", "-q")
        self.assertIsNone(self.evaluate({"cwd": str(child)}))

    def test_existing_dispatcher_still_receives_child_stop_payload(self) -> None:
        child = self.workspace / "ai-arena"
        script = child / "tools/codex-hook-stop.sh"
        script.parent.mkdir(parents=True)
        script.write_text(
            "#!/usr/bin/env python3\n"
            "import json\n"
            "import sys\n"
            "raise SystemExit(0 if json.load(sys.stdin)['cwd'].endswith('ai-arena') else 1)\n",
            encoding="utf-8",
        )
        script.chmod(0o755)
        with patch.object(dispatch, "repo_root", return_value=self.workspace), patch.object(
            sys, "stdin", io.StringIO(json.dumps({"cwd": str(child)}))
        ), patch.object(sys, "argv", ["dispatch_hook.py", "stop"]):
            self.assertEqual(0, dispatch.main())

    def test_workspace_descendant_is_in_scope(self) -> None:
        self.branch("feat/stop-goal-reminder")
        nested = self.workspace / "nested"
        nested.mkdir()
        self.assertEqual("block", self.evaluate({"cwd": str(nested)})["decision"])

    def test_workspace_worktree_is_in_scope(self) -> None:
        worktree = Path(self.temp.name) / "feature-worktree"
        git(self.workspace, "worktree", "add", "-qb", "feat/worktree-scope", str(worktree))
        directory = worktree / "docs/exec-plan/todo"
        directory.mkdir(parents=True)
        (directory / "0047-worktree-scope.md").write_text(
            "# Worktree Scope\n\n## Completion boundary\n\nReady for review.\n",
            encoding="utf-8",
        )
        with patch.object(reminder, "HOOK_ROOT", worktree):
            result = reminder.evaluate({"cwd": str(worktree)})
        self.assertEqual("block", result["decision"])
        self.assertIn("Worktree Scope", result["reason"])

    def test_ineligible_branch_and_missing_cwd_pass_through(self) -> None:
        self.branch("chore/no-reminder")
        self.assertIsNone(self.evaluate({"cwd": str(self.workspace)}))
        self.assertIsNone(self.evaluate({}))

    def test_main_emits_valid_json_and_malformed_input_fails_open(self) -> None:
        self.branch("feat/stop-goal-reminder")
        with patch.object(reminder, "HOOK_ROOT", self.workspace), patch.object(sys, "stdin", io.StringIO(json.dumps({"cwd": str(self.workspace)}))), contextlib.redirect_stdout(io.StringIO()) as output:
            self.assertEqual(0, reminder.main())
        self.assertEqual("block", json.loads(output.getvalue())["decision"])

        with patch.object(sys, "stdin", io.StringIO("{")), contextlib.redirect_stderr(io.StringIO()) as error:
            self.assertEqual(0, reminder.main())
        self.assertIn("malformed input", error.getvalue())


if __name__ == "__main__":
    unittest.main()
