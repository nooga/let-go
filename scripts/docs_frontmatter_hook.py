#!/usr/bin/env python3
"""
Frontmatter maintenance for docs/**/*.md.

Two modes:

  default: For each file in argv, stub the frontmatter if missing,
           or bump `last-verified:` to today if older. Mutates files
           in place. Used as a pre-commit hook.

  --check: Validate frontmatter is present and well-formed on each
           file in argv. Never mutates. Exits non-zero on any
           failure. Used in CI.

The mode split is deliberate: CI enforces the *floor* (frontmatter
present, parseable), pre-commit handles *maintenance* (last-verified
bump). The two must not conflict — a doc authored a week ago is
CI-clean even though the bump hasn't run.

Both modes never touch: `status:` (when present),
`authoritative-for:`, `supersedes:`, `superseded-by:`, `shipped:`,
`remaining-open:`, or `human-verified:`. Those are human-authored.
The bump matches only top-level `last-verified:` (no leading
whitespace), so indented occurrences inside block scalars or nested
mappings are left alone.

Stdlib only — no external dependencies.
"""

from __future__ import annotations

import argparse
import datetime
import re
import sys
from pathlib import Path

TODAY = datetime.date.today().isoformat()
DELIM = "---"
LAST_VERIFIED_RE = re.compile(r"^last-verified\s*:\s*(.*?)\s*$")
STATUS_RE = re.compile(r"^status\s*:\s*\S")


class MalformedFrontmatter(Exception):
    """Opening `---` present but no closing delimiter (or other broken shape)."""


def find_close(lines: list[str]) -> int | None:
    """
    Locate the closing `---` line of the frontmatter block.

    Returns:
      None       — file has no frontmatter at all.
      int        — index of the closing delim line.

    Raises:
      MalformedFrontmatter — opening `---` but no closing delim.
    """
    if not lines or lines[0] != DELIM:
        return None
    for i in range(1, len(lines)):
        if lines[i] == DELIM:
            return i
    raise MalformedFrontmatter("opening `---` present but no closing delimiter")


def stub_block() -> str:
    return (
        f"{DELIM}\n"
        f"status: active\n"
        f"last-verified: {TODAY}\n"
        f"human-verified:\n"
        f"{DELIM}\n\n"
    )


def bump(lines: list[str], close_idx: int) -> bool:
    """Bump top-level `last-verified:` inside the frontmatter. True if changed."""
    for i in range(1, close_idx):
        match = LAST_VERIFIED_RE.match(lines[i])
        if not match:
            continue
        existing = match.group(1).strip()
        if existing and existing >= TODAY:
            return False
        lines[i] = f"last-verified: {TODAY}"
        return True
    lines.insert(close_idx, f"last-verified: {TODAY}")
    return True


def read_text(path: Path) -> str:
    # utf-8-sig strips a BOM if present; the file is written back without one.
    return path.read_text(encoding="utf-8-sig")


def maintain(path: Path) -> str | None:
    """Stub or bump as needed. Returns an action description or None."""
    text = read_text(path)
    lines = text.split("\n")
    close_idx = find_close(lines)

    if close_idx is None:
        path.write_text(stub_block() + text, encoding="utf-8")
        return f"stubbed: {path} (status: active, last-verified: {TODAY})"

    if bump(lines, close_idx):
        path.write_text("\n".join(lines), encoding="utf-8")
        return f"bumped last-verified: {path}"

    return None


def check(path: Path) -> list[str]:
    """Validate frontmatter is present and well-formed. Returns a list of issues."""
    text = read_text(path)
    lines = text.split("\n")
    try:
        close_idx = find_close(lines)
    except MalformedFrontmatter as e:
        return [str(e)]

    if close_idx is None:
        return ["no frontmatter (expected opening `---` block)"]

    issues: list[str] = []
    if not any(STATUS_RE.match(lines[i]) for i in range(1, close_idx)):
        issues.append("missing top-level `status:`")
    if not any(LAST_VERIFIED_RE.match(lines[i]) for i in range(1, close_idx)):
        issues.append("missing top-level `last-verified:`")
    return issues


def eligible(path: Path) -> bool:
    # Symlinks are skipped: a staged `docs/*.md` symlink would otherwise
    # cause the hook to write through to whatever the symlink points at.
    return (
        not path.is_symlink()
        and path.is_file()
        and path.suffix.lower() == ".md"
    )


def run_maintain(paths: list[str]) -> int:
    actions: list[str] = []
    errors: list[str] = []
    for arg in paths:
        path = Path(arg)
        if not eligible(path):
            continue
        try:
            result = maintain(path)
        except MalformedFrontmatter as e:
            errors.append(f"{path}: {e}")
            continue
        except (OSError, UnicodeDecodeError) as e:
            errors.append(f"{path}: {e}")
            continue
        if result:
            actions.append(result)

    for a in actions:
        print(f"[docs-frontmatter] {a}")
    if actions:
        print(
            "[docs-frontmatter] Note: human-verified is set only by explicit "
            "human action; this hook and any automation must leave it blank. "
            "See docs/frontmatter-hook.md."
        )
    for err in errors:
        print(f"[docs-frontmatter] error: {err}", file=sys.stderr)
    return 1 if errors else 0


def run_check(paths: list[str]) -> int:
    failures: list[tuple[Path, list[str]]] = []
    for arg in paths:
        path = Path(arg)
        if not eligible(path):
            continue
        try:
            issues = check(path)
        except (OSError, UnicodeDecodeError) as e:
            issues = [str(e)]
        if issues:
            failures.append((path, issues))

    if not failures:
        return 0

    print("Frontmatter check failed:", file=sys.stderr)
    for path, issues in failures:
        for issue in issues:
            print(f"  {path}: {issue}", file=sys.stderr)
    print(
        "\nSee docs/frontmatter-hook.md. To fix locally, run the pre-commit "
        "hook (or `python3 scripts/docs_frontmatter_hook.py <file>`) and "
        "re-commit.",
        file=sys.stderr,
    )
    return 1


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[1])
    parser.add_argument(
        "--check",
        action="store_true",
        help="Validate frontmatter without mutating. Exit non-zero on failure.",
    )
    parser.add_argument("paths", nargs="*")
    args = parser.parse_args(argv)
    if args.check:
        return run_check(args.paths)
    return run_maintain(args.paths)


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
