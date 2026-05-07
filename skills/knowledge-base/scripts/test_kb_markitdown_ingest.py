from __future__ import annotations

import argparse
import io
import tempfile
import unittest
import unittest.mock
import urllib.error
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

    def test_ensure_job_dir_resolves_job_dir_path(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            target = root / "my-job"
            args = argparse.Namespace(
                job_dir=str(target),
                scratch_root=None,
                source="sample.pdf",
                source_slug=None,
            )
            result = mod.ensure_job_dir(args)
            self.assertEqual(result, target.resolve())
            self.assertTrue(result.is_dir())

    def test_ensure_job_dir_resolves_scratch_root_path(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            scratch = root / "scratch"
            args = argparse.Namespace(
                job_dir=None,
                scratch_root=str(scratch),
                source="sample.pdf",
                source_slug=None,
            )
            result = mod.ensure_job_dir(args)
            self.assertTrue(result.is_dir())
            self.assertEqual(result.parent, scratch.resolve())

    def test_download_source_streams_content(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            dest = Path(tmp) / "inputs" / "file.pdf"
            fake_content = b"fake pdf content"
            mock_response = unittest.mock.MagicMock()
            mock_response.__enter__ = lambda s: s
            mock_response.__exit__ = unittest.mock.MagicMock(return_value=False)
            mock_response.read.side_effect = [fake_content, b""]
            with unittest.mock.patch("urllib.request.urlopen", return_value=mock_response):
                result = mod.download_source("https://example.com/file.pdf", dest)
            self.assertEqual(result, dest)
            self.assertEqual(dest.read_bytes(), fake_content)

    def test_download_source_raises_on_oversized_download(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            dest = Path(tmp) / "inputs" / "big.pdf"
            oversized_chunk = b"x" * (mod.DOWNLOAD_MAX_BYTES + 1)
            mock_response = unittest.mock.MagicMock()
            mock_response.__enter__ = lambda s: s
            mock_response.__exit__ = unittest.mock.MagicMock(return_value=False)
            mock_response.read.return_value = oversized_chunk
            with unittest.mock.patch("urllib.request.urlopen", return_value=mock_response):
                with self.assertRaises(RuntimeError):
                    mod.download_source("https://example.com/big.pdf", dest)

    def test_default_markitdown_with_includes_epub(self) -> None:
        self.assertIn("epub", mod.DEFAULT_MARKITDOWN_WITH)


if __name__ == "__main__":
    unittest.main()
