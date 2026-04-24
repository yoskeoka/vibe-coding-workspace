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
SOURCE_JA_DOCS_DIR = SOURCE_DOCS_DIR / "ja"
GENERATED_PARENT = ROOT_DIR / ".local" / "kb-generated"
DEFAULT_GENERATED_ROOT = GENERATED_PARENT / "direct"
GENERATED_ROOT = DEFAULT_GENERATED_ROOT
GENERATED_DOCS_DIR = GENERATED_ROOT / "docs"
GENERATED_CONFIG_FILE = GENERATED_ROOT / "mkdocs.kb.generated.yml"
TEMPLATE_CONFIG_FILE = ROOT_DIR / "mkdocs.kb.template.yml"
SITE_DIR = ROOT_DIR / ".site" / "kb"

DEFAULT_LOCALE = "en"
JA_LOCALE = "ja"
LOCALE_SUFFIX = {JA_LOCALE: ".ja.md"}
VISIBLE_SECTION_HEADINGS = {
    DEFAULT_LOCALE: {"sources": "Sources", "related_pages": "Related pages"},
    JA_LOCALE: {"sources": "ソース", "related_pages": "関連ページ"},
}
SECTION_LABELS = {
    DEFAULT_LOCALE: {
        "topics": "Topics",
        "tools": "Tools",
        "patterns": "Patterns",
        "projects": "Projects",
        "log": "Log",
        "sources": "Sources",
        "index": "Index",
        "kb_docs": "KB Docs",
        "home": "Home",
        "schema": "Schema",
        "ingest": "Ingest",
    },
    JA_LOCALE: {
        "topics": "トピック",
        "tools": "ツール",
        "patterns": "パターン",
        "projects": "プロジェクト",
        "log": "ログ",
        "sources": "ソース",
        "index": "一覧",
        "kb_docs": "KB Docs",
        "home": "ホーム",
        "schema": "Schema",
        "ingest": "Ingest",
    },
}
NAV_PATHS = {
    "topics": [
        "wiki/topics/llm-knowledge-bases.md",
        "wiki/topics/ai-game-dev.md",
        "wiki/topics/deployment-options.md",
        "wiki/topics/pixel-art-ui.md",
    ],
    "tools": [
        "wiki/tools/phaser.md",
        "wiki/tools/godot.md",
        "wiki/tools/next2d.md",
    ],
    "patterns": [
        "wiki/patterns/source-ingestion.md",
        "wiki/patterns/cheap-hosting.md",
        "wiki/patterns/branching-narrative-authoring.md",
        "wiki/patterns/component-oriented-gameplay.md",
        "wiki/patterns/go-transaction-boundaries.md",
        "wiki/patterns/go-container-workflows.md",
    ],
    "projects": [
        "wiki/projects/ww.md",
        "wiki/projects/reversi-adventure.md",
        "wiki/projects/vim-learning-game.md",
    ],
}


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


def locale_for_rel(rel: str) -> str:
    for locale, suffix in LOCALE_SUFFIX.items():
        if rel.endswith(suffix):
            return locale
    return DEFAULT_LOCALE


def localized_rel(rel: str, locale: str) -> str:
    if locale == DEFAULT_LOCALE or not rel.endswith(".md"):
        return rel
    return rel[:-3] + LOCALE_SUFFIX[locale]


def best_rel_for_locale(rel: str, locale: str, title_map: dict[str, str]) -> str:
    localized = localized_rel(rel, locale)
    if localized in title_map:
        return localized
    return rel


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
    current_locale = locale_for_rel(current_rel)
    current_dir = posixpath.dirname(current_rel)

    for target in targets:
        resolved = resolve_target(current_rel, target)
        preferred = best_rel_for_locale(resolved, current_locale, title_map)
        label = title_map.get(preferred, title_map.get(resolved, Path(resolved).stem))
        rendered_target = posixpath.relpath(preferred, current_dir)
        lines.append(f"- [{label}]({rendered_target})")
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
    locale = locale_for_rel(current_rel)
    headings = VISIBLE_SECTION_HEADINGS[locale]
    updated = body.rstrip() + "\n"

    sources = meta.get("sources")
    if current_rel.startswith("wiki/") and isinstance(sources, list) and sources and f"## {headings['sources']}" not in updated:
        section = render_link_section(current_rel, sources, title_map, headings["sources"])
        updated = insert_before_heading(updated, headings["related_pages"], section)

    related_pages = meta.get("related_pages")
    if current_rel.startswith("sources/") and isinstance(related_pages, list) and related_pages and f"## {headings['related_pages']}" not in updated:
        section = render_link_section(current_rel, related_pages, title_map, headings["related_pages"])
        updated = updated.rstrip() + "\n\n" + section + "\n"

    return updated


def collect_source_notes(locale: str, title_map: dict[str, str]) -> dict[str, list[Path]]:
    by_year: dict[str, list[Path]] = {}
    for path in sorted((GENERATED_DOCS_DIR / "sources").glob("*/*.md")):
        name = path.name
        if name == "index.md" or name.endswith(".ja.md"):
            continue
        if locale == JA_LOCALE:
            localized = localized_rel(doc_rel_path(path), JA_LOCALE)
            if localized in title_map:
                path = GENERATED_DOCS_DIR / localized
        year = path.parent.name
        by_year.setdefault(year, []).append(path)

    for paths in by_year.values():
        paths.sort(key=lambda item: item.name, reverse=True)
    return dict(sorted(by_year.items(), reverse=True))


def generate_year_index(year: str, paths: list[Path], title_map: dict[str, str], locale: str) -> None:
    filename = "index.md" if locale == DEFAULT_LOCALE else f"index.{locale}.md"
    target = GENERATED_DOCS_DIR / "sources" / year / filename
    heading = f"# Sources from {year}" if locale == DEFAULT_LOCALE else f"# {year}年のソース"
    lines = [heading, ""]
    for path in paths:
        source_rel = doc_rel_path(path)
        date_label = path.stem[:10]
        lines.append(f"- {date_label}: [{title_map[source_rel]}]({path.name})")
    lines.append("")
    target.write_text("\n".join(lines), encoding="utf-8")


def generate_sources_index(by_year: dict[str, list[Path]], locale: str) -> None:
    filename = "index.md" if locale == DEFAULT_LOCALE else f"index.{locale}.md"
    target = GENERATED_DOCS_DIR / "sources" / filename
    heading = "# Sources" if locale == DEFAULT_LOCALE else "# ソース"
    intro = "Browse source notes by year." if locale == DEFAULT_LOCALE else "年ごとにソースノートを参照します。"
    empty = "No source notes have been ingested yet." if locale == DEFAULT_LOCALE else "まだソースノートはありません。"
    lines = [heading, ""]
    if not by_year:
        lines.append(empty)
        lines.append("")
        target.write_text("\n".join(lines), encoding="utf-8")
        return

    lines.append(intro)
    lines.append("")
    for year, paths in by_year.items():
        year_index = "index.md" if locale == DEFAULT_LOCALE else f"index.{locale}.md"
        lines.append(f"- [{year} ({len(paths)})]({year}/{year_index})")
    lines.append("")
    target.write_text("\n".join(lines), encoding="utf-8")


def require_generated_file(path: Path, description: str) -> None:
    if not path.is_file():
        rel = path.relative_to(GENERATED_DOCS_DIR)
        raise SystemExit(f"missing generated KB {description}: {rel}")


def rewrite_root_links(text: str) -> str:
    replacements = {
        "](projects/": "](wiki/projects/",
        "](topics/": "](wiki/topics/",
        "](tools/": "](wiki/tools/",
        "](patterns/": "](wiki/patterns/",
        "](log.md)": "](wiki/log.md)",
    }
    for old, new in replacements.items():
        text = text.replace(old, new)
    return text


def generate_root_indexes() -> None:
    source = GENERATED_DOCS_DIR / "wiki" / "index.md"
    target = GENERATED_DOCS_DIR / "index.md"
    require_generated_file(source, "wiki index")
    target.write_text(rewrite_root_links(source.read_text(encoding="utf-8")), encoding="utf-8")
    source.unlink()

    ja_source = GENERATED_DOCS_DIR / "wiki" / "index.ja.md"
    if ja_source.exists():
        ja_target = GENERATED_DOCS_DIR / "index.ja.md"
        ja_target.write_text(rewrite_root_links(ja_source.read_text(encoding="utf-8")), encoding="utf-8")
        ja_source.unlink()


def relocate_readme() -> None:
    source = GENERATED_DOCS_DIR / "README.md"
    target = GENERATED_DOCS_DIR / "kb-docs" / "README.md"
    require_generated_file(source, "README")
    text = source.read_text(encoding="utf-8")
    text = text.replace("(schema.md)", "(../schema.md)")
    text = text.replace("(ingest.md)", "(../ingest.md)")
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(text, encoding="utf-8")
    source.unlink()


def copy_source_docs() -> None:
    GENERATED_DOCS_DIR.parent.mkdir(parents=True, exist_ok=True)
    shutil.copytree(SOURCE_DOCS_DIR, GENERATED_DOCS_DIR, ignore=shutil.ignore_patterns("ja"))


def copy_japanese_mirror() -> None:
    if not SOURCE_JA_DOCS_DIR.exists():
        return

    for source in sorted(SOURCE_JA_DOCS_DIR.rglob("*")):
        rel = source.relative_to(SOURCE_JA_DOCS_DIR)
        target = GENERATED_DOCS_DIR / rel
        if source.is_dir():
            target.mkdir(parents=True, exist_ok=True)
            continue

        if source.suffix == ".md":
            target = target.with_name(f"{target.stem}.ja.md")
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(source, target)


def build_sources_nav(by_year: dict[str, list[Path]], title_map: dict[str, str], locale: str) -> list[str]:
    labels = SECTION_LABELS[locale]
    lines = [f"      - {labels['index']}: {best_rel_for_locale('sources/index.md', locale, title_map)}"]
    for year, paths in by_year.items():
        lines.append(f"      - {yaml_quote(f'{year} ({len(paths)})')}:")
        year_index = best_rel_for_locale(f"sources/{year}/index.md", locale, title_map)
        lines.append(f"          - {labels['index']}: {year_index}")
        for path in paths:
            source_rel = doc_rel_path(path)
            date_label = path.stem[:10]
            lines.append(f"          - {yaml_quote(f'{date_label} {title_map[source_rel]}')}: {source_rel}")
    return lines


def build_nav(locale: str, title_map: dict[str, str], sources_nav: list[str]) -> str:
    labels = SECTION_LABELS[locale]
    lines: list[str] = [f"  - {labels['home']}: {best_rel_for_locale('index.md', locale, title_map)}"]

    for section in ["topics", "tools", "patterns", "projects"]:
        lines.append(f"  - {labels[section]}:")
        for rel in NAV_PATHS[section]:
            preferred = best_rel_for_locale(rel, locale, title_map)
            lines.append(f"      - {yaml_quote(title_map[preferred])}: {preferred}")

    lines.append(f"  - {labels['log']}: {best_rel_for_locale('wiki/log.md', locale, title_map)}")
    lines.append(f"  - {labels['sources']}:")
    lines.extend(sources_nav)
    lines.append(f"  - {labels['kb_docs']}:")
    lines.append(f"      - {labels['home']}: kb-docs/README.md")
    lines.append(f"      - {labels['schema']}: schema.md")
    lines.append(f"      - {labels['ingest']}: ingest.md")
    return "\n".join(lines)


def indent_block(text: str, spaces: int) -> str:
    prefix = " " * spaces
    return "\n".join(f"{prefix}{line}" if line else "" for line in text.splitlines())


def write_generated_config(en_nav: str, ja_nav: str) -> None:
    template = TEMPLATE_CONFIG_FILE.read_text(encoding="utf-8")
    rendered = template.replace("__GENERATED_DOCS_DIR__", str(GENERATED_DOCS_DIR))
    rendered = rendered.replace("__GENERATED_SITE_DIR__", str(SITE_DIR))
    rendered = rendered.replace("__GENERATED_EN_NAV__", en_nav)
    rendered = rendered.replace("__GENERATED_JA_NAV__", indent_block(ja_nav, 12))
    GENERATED_CONFIG_FILE.write_text(rendered, encoding="utf-8")


def main() -> None:
    args = parse_args()
    configure_generated_root(args.generated_root)

    if GENERATED_ROOT.exists():
        shutil.rmtree(GENERATED_ROOT)
    copy_source_docs()
    copy_japanese_mirror()
    generate_root_indexes()
    relocate_readme()

    title_map = build_title_map()
    by_year_en = collect_source_notes(DEFAULT_LOCALE, title_map)
    by_year_ja = collect_source_notes(JA_LOCALE, title_map)
    generate_sources_index(by_year_en, DEFAULT_LOCALE)
    generate_sources_index(by_year_ja, JA_LOCALE)
    for year, paths in by_year_en.items():
        generate_year_index(year, paths, title_map, DEFAULT_LOCALE)
    for year, paths in by_year_ja.items():
        generate_year_index(year, paths, title_map, JA_LOCALE)

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

    title_map = build_title_map()
    en_nav = build_nav(DEFAULT_LOCALE, title_map, build_sources_nav(by_year_en, title_map, DEFAULT_LOCALE))
    ja_nav = build_nav(JA_LOCALE, title_map, build_sources_nav(by_year_ja, title_map, JA_LOCALE))
    write_generated_config(en_nav, ja_nav)


if __name__ == "__main__":
    main()
