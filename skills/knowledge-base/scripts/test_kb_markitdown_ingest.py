from __future__ import annotations

import tempfile
import unittest
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))

import kb_markitdown_ingest as mod


class MarkItDownIngestTests(unittest.TestCase):
    def test_sanitize_slug(self) -> None:
        self.assertEqual(mod.sanitize_slug("Hello, KB World!"), "hello-kb-world")
        self.assertEqual(mod.sanitize_slug("___"), "kb-source")

    def test_default_slug_for_url(self) -> None:
        self.assertEqual(
            mod.default_slug_for_source("https://example.com/files/My Deck.pptx"),
            "my-deck",
        )

    def test_resolve_source_for_local_path(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            source = root / "sample.pdf"
            source.write_text("placeholder", encoding="utf-8")
            job_dir = root / "job"
            result = mod.resolve_source(str(source), job_dir)
            self.assertEqual(result["source_kind"], "local_path")
            self.assertEqual(result["source_url"], "")
            self.assertEqual(result["source_file"], str(source.resolve()))
            self.assertEqual(result["input_path"], str(source.resolve()))

    def test_resolve_source_for_file_uri(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            source = root / "sample.docx"
            source.write_text("placeholder", encoding="utf-8")
            job_dir = root / "job"
            uri = source.resolve().as_uri()
            result = mod.resolve_source(uri, job_dir)
            self.assertEqual(result["source_kind"], "file_uri")
            self.assertEqual(result["source_url"], uri)
            self.assertEqual(result["source_file"], str(source.resolve()))

    def test_build_source_context_mentions_durable_rules(self) -> None:
        metadata = {
            "source_argument": "sample.pdf",
            "source_kind": "local_path",
            "source_url": "",
            "source_file": "sample.pdf",
            "converter_version": "markitdown 0.1.5",
            "converted_markdown": "/tmp/job/outputs/converted.md",
            "workspace_relevance": "Useful for KB coverage.",
        }
        context = mod.build_source_context(metadata)
        self.assertIn("do not commit `converted.md`", context)
        self.assertIn("conversion_method", context)
        self.assertIn("Workspace relevance", context)


if __name__ == "__main__":
    unittest.main()
