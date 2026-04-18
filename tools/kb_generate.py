#!/usr/bin/env python3
from __future__ import annotations

import argparse
import os
import posixpath
import re
import shutil
from pathlib import Path


ROOT_DIR = Path(__file__).resolve().parent.parent
SOURCE_DOCS_DIR = ROOT_DIR / "docs" / "kb"
GENERATED_PARENT = ROOT_DIR / ".local" / "kb-generated"
DEFAULT_GENERATED_ROOT = GENERATED_PARENT / "direct"
GENERATED_ROOT = DEFAULT_GENERATED_ROOT
GENERATED_DOCS_DIR = GENERATED_ROOT / "docs"
GENERATED_CONFIG_FILE = GENERATED_ROOT / "mkdocs.kb.generated.yml"
TEMPLATE_CONFIG_FILE = ROOT_DIR / "mkdocs.kb.template.yml"
SITE_DIR = ROOT_DIR / ".site" / "kb"


def configured_generated_root() -> str | None:
    value = os.environ.get("KB_GENERATED_ROOT")
    if value and value.strip():
        return value
    return None


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Generate MkDocs inputs for the workspace KB.")
    parser.add_argument(
        "--generated-root",
        default=configured_generated_root(),
        help="Invocation-owned directory for generated docs and config files.",
    )
    return parser.parse_args()


def configure_generated_root(path: str | None) -> None:
    global GENERATED_ROOT, GENERATED_DOCS_DIR, GENERATED_CONFIG_FILE

    if path:
        candidate = Path(path).expanduser().resolve()
    else:
        candidate = DEFAULT_GENERATED_ROOT.resolve()

    generated_parent = GENERATED_PARENT.resolve()
    try:
        candidate.relative_to(generated_parent)
    except ValueError as exc:
        raise SystemExit(f"generated root must be inside {generated_parent}: {candidate}") from exc
    if candidate == generated_parent:
        raise SystemExit(f"generated root must not be the shared parent directory: {generated_parent}")

    GENERATED_ROOT = candidate
    GENERATED_DOCS_DIR = GENERATED_ROOT / "docs"
    GENERATED_CONFIG_FILE = GENERATED_ROOT / "mkdocs.kb.generated.yml"


def parse_frontmatter(text: str) -> tuple[dict[str, object], str]:
    lines = text.splitlines()
    if not lines or lines[0].strip() != "---":
        return {}, text

    end = None
    for idx in range(1, len(lines)):
        if lines[idx].strip() == "---":
            end = idx
            break
    if end is None:
        return {}, text

    meta_lines = lines[1:end]
    body = "\n".join(lines[end + 1 :])
    meta: dict[str, object] = {}
    current_key: str | None = None

    for raw in meta_lines:
        if not raw.strip():
            continue
        if raw.startswith("  - ") and current_key:
            values = meta.setdefault(current_key, [])
            if isinstance(values, list):
                values.append(raw[4:].strip())
            continue

        current_key = None
        if ":" not in raw:
            continue

        key, value = raw.split(":", 1)
        key = key.strip()
        value = value.strip()
        if not value:
            meta[key] = []
            current_key = key
            continue
        meta[key] = value

    return meta, body


def first_heading(body: str) -> str | None:
    for line in body.splitlines():
        if line.startswith("# "):
            return line[2:].strip()
    return None


def doc_rel_path(path: Path) -> str:
    return path.relative_to(GENERATED_DOCS_DIR).as_posix()


def title_for_doc(path: Path) -> str:
    text = path.read_text(encoding="utf-8")
    meta, body = parse_frontmatter(text)
    title = meta.get("title")
    if isinstance(title, str) and title:
        return title
    heading = first_heading(body)
    return heading or path.stem


def build_title_map() -> dict[str, str]:
    mapping: dict[str, str] = {}
    for path in sorted(GENERATED_DOCS_DIR.rglob("*.md")):
        mapping[doc_rel_path(path)] = title_for_doc(path)
    return mapping


def resolve_target(current_rel: str, target: str) -> str:
    return posixpath.normpath(posixpath.join(posixpath.dirname(current_rel), target))


def render_link_section(current_rel: str, targets: list[str], title_map: dict[str, str], heading: str) -> str:
    lines = [f"## {heading}", ""]
    for target in targets:
        resolved = resolve_target(current_rel, target)
        label = title_map.get(resolved, Path(resolved).stem)
        lines.append(f"- [{label}]({target})")
    return "\n".join(lines)


def yaml_quote(text: str) -> str:
    escaped = text.replace("\\", "\\\\").replace('"', '\\"')
    return f'"{escaped}"'


def insert_before_heading(body: str, heading: str, section: str) -> str:
    pattern = re.compile(rf"^## {re.escape(heading)}\s*$", re.MULTILINE)
    match = pattern.search(body)
    if not match:
        return body.rstrip() + "\n\n" + section + "\n"
    return body[: match.start()].rstrip() + "\n\n" + section + "\n\n" + body[match.start() :].lstrip()


def enhance_markdown(current_rel: str, meta: dict[str, object], body: str, title_map: dict[str, str]) -> str:
    updated = body.rstrip() + "\n"

    sources = meta.get("sources")
    if current_rel.startswith("wiki/") and isinstance(sources, list) and sources and "## Sources" not in updated:
        section = render_link_section(current_rel, sources, title_map, "Sources")
        updated = insert_before_heading(updated, "Related pages", section)

    related_pages = meta.get("related_pages")
    if current_rel.startswith("sources/") and isinstance(related_pages, list) and related_pages and "## Related pages" not in updated:
        section = render_link_section(current_rel, related_pages, title_map, "Related pages")
        updated = updated.rstrip() + "\n\n" + section + "\n"

    return updated


def collect_source_notes() -> dict[str, list[Path]]:
    by_year: dict[str, list[Path]] = {}
    for path in sorted((GENERATED_DOCS_DIR / "sources").glob("*/*.md")):
        if path.name == "index.md":
            continue
        year = path.parent.name
        by_year.setdefault(year, []).append(path)
    for paths in by_year.values():
        paths.sort(key=lambda item: item.name, reverse=True)
    return dict(sorted(by_year.items(), reverse=True))


def generate_year_index(year: str, paths: list[Path], title_map: dict[str, str]) -> None:
    target = GENERATED_DOCS_DIR / "sources" / year / "index.md"
    lines = [f"# Sources from {year}", ""]
    for path in paths:
        source_rel = doc_rel_path(path)
        date_label = path.stem[:10]
        lines.append(f"- {date_label}: [{title_map[source_rel]}]({path.name})")
    lines.append("")
    target.write_text("\n".join(lines), encoding="utf-8")


def generate_sources_index(by_year: dict[str, list[Path]]) -> None:
    target = GENERATED_DOCS_DIR / "sources" / "index.md"
    lines = ["# Sources", ""]
    if not by_year:
        lines.append("No source notes have been ingested yet.")
        lines.append("")
        target.write_text("\n".join(lines), encoding="utf-8")
        return

    lines.append("Browse source notes by year.")
    lines.append("")
    for year, paths in by_year.items():
        lines.append(f"- [{year} ({len(paths)})]({year}/index.md)")
    lines.append("")
    target.write_text("\n".join(lines), encoding="utf-8")


def generate_root_index() -> None:
    source = GENERATED_DOCS_DIR / "wiki" / "index.md"
    target = GENERATED_DOCS_DIR / "index.md"
    text = source.read_text(encoding="utf-8")
    replacements = {
        "](projects/": "](wiki/projects/",
        "](topics/": "](wiki/topics/",
        "](tools/": "](wiki/tools/",
        "](patterns/": "](wiki/patterns/",
        "](log.md)": "](wiki/log.md)",
    }
    for old, new in replacements.items():
        text = text.replace(old, new)
    target.write_text(text, encoding="utf-8")
    source.unlink()


def relocate_readme() -> None:
    source = GENERATED_DOCS_DIR / "README.md"
    target = GENERATED_DOCS_DIR / "kb-docs" / "README.md"
    text = source.read_text(encoding="utf-8")
    text = text.replace("(schema.md)", "(../schema.md)")
    text = text.replace("(ingest.md)", "(../ingest.md)")
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(text, encoding="utf-8")
    source.unlink()


def build_sources_nav(by_year: dict[str, list[Path]], title_map: dict[str, str]) -> str:
    lines = ["      - Index: sources/index.md"]
    for year, paths in by_year.items():
        lines.append(f"      - {yaml_quote(f'{year} ({len(paths)})')}:")
        lines.append(f"          - Index: sources/{year}/index.md")
        for path in paths:
            source_rel = doc_rel_path(path)
            date_label = path.stem[:10]
            label = yaml_quote(f"{date_label} {title_map[source_rel]}")
            lines.append(f"          - {label}: {source_rel}")
    return "\n".join(lines)


def write_generated_config(nav_block: str) -> None:
    template = TEMPLATE_CONFIG_FILE.read_text(encoding="utf-8")
    rendered = template.replace("__GENERATED_DOCS_DIR__", str(GENERATED_DOCS_DIR))
    rendered = rendered.replace("__GENERATED_SITE_DIR__", str(SITE_DIR))
    rendered = rendered.replace("__GENERATED_SOURCES_NAV__", nav_block)
    GENERATED_CONFIG_FILE.write_text(rendered, encoding="utf-8")


def main() -> None:
    args = parse_args()
    configure_generated_root(args.generated_root)

    if GENERATED_ROOT.exists():
        shutil.rmtree(GENERATED_ROOT)
    GENERATED_DOCS_DIR.parent.mkdir(parents=True, exist_ok=True)
    shutil.copytree(SOURCE_DOCS_DIR, GENERATED_DOCS_DIR)
    generate_root_index()
    relocate_readme()

    initial_title_map = build_title_map()
    by_year = collect_source_notes()
    generate_sources_index(by_year)
    for year, paths in by_year.items():
        generate_year_index(year, paths, initial_title_map)

    title_map = build_title_map()
    for path in sorted(GENERATED_DOCS_DIR.rglob("*.md")):
        current_rel = doc_rel_path(path)
        text = path.read_text(encoding="utf-8")
        meta, body = parse_frontmatter(text)
        if not meta:
            continue
        updated_body = enhance_markdown(current_rel, meta, body, title_map)
        if updated_body == body.rstrip() + "\n":
            continue
        frontmatter = text[: text.find("---", 3) + 4]
        path.write_text(frontmatter + "\n" + updated_body, encoding="utf-8")

    write_generated_config(build_sources_nav(by_year, title_map))


if __name__ == "__main__":
    main()
