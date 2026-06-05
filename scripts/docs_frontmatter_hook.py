#!/usr/bin/env python3
"""
Frontmatter maintenance for docs/**/*.md.

Run from pre-commit. For each staged file passed on argv:

  - If the file has no YAML frontmatter, prepend a minimal stub:
        status: active
        last-verified: <today>
        human-verified:
  - If the file has frontmatter, bump `last-verified:` to today (only if
    the existing value is older or missing). Idempotent on docs already
    stamped today.

Never touches: `status:` (when present), `authoritative-for:`,
`supersedes:`, `superseded-by:`, `shipped:`, `remaining-open:`, or
`human-verified:`. Those fields are human-authored.

Stdlib only — no external dependencies.
"""

from __future__ import annotations

import datetime
import re
import sys
from pathlib import Path

TODAY = datetime.date.today().isoformat()
DELIM = "---"
LAST_VERIFIED_RE = re.compile(r"^(\s*)last-verified\s*:\s*(.*?)\s*$")


def find_frontmatter_close(lines: list[str]) -> int | None:
    """Return the index of the closing `---` line, or None if no frontmatter."""
    if not lines or lines[0] != DELIM:
        return None
    for i in range(1, len(lines)):
        if lines[i] == DELIM:
            return i
    return None


def stub_block() -> str:
    return (
        f"{DELIM}\n"
        f"status: active\n"
        f"last-verified: {TODAY}\n"
        f"human-verified:\n"
        f"{DELIM}\n\n"
    )


def bump(lines: list[str], close_idx: int) -> bool:
    """Bump `last-verified:` inside the frontmatter block. Returns True if changed."""
    for i in range(1, close_idx):
        match = LAST_VERIFIED_RE.match(lines[i])
        if not match:
            continue
        indent, existing = match.group(1), match.group(2).strip()
        if existing and existing >= TODAY:
            return False
        lines[i] = f"{indent}last-verified: {TODAY}"
        return True
    lines.insert(close_idx, f"last-verified: {TODAY}")
    return True


def process(path: Path) -> str | None:
    """Process one file; return a short action string, or None if unchanged."""
    try:
        text = path.read_text(encoding="utf-8")
    except (OSError, UnicodeDecodeError):
        return None

    lines = text.split("\n")
    close_idx = find_frontmatter_close(lines)

    if close_idx is None:
        path.write_text(stub_block() + text, encoding="utf-8")
        return f"stubbed: {path} (status: active, last-verified: {TODAY})"

    if bump(lines, close_idx):
        path.write_text("\n".join(lines), encoding="utf-8")
        return f"bumped last-verified: {path}"

    return None


def main(argv: list[str]) -> int:
    actions: list[str] = []
    for arg in argv:
        path = Path(arg)
        if not path.is_file() or path.suffix.lower() != ".md":
            continue
        result = process(path)
        if result:
            actions.append(result)

    if actions:
        for a in actions:
            print(f"[docs-frontmatter] {a}")
        print(
            "[docs-frontmatter] Note: human-verified is set only by explicit human "
            "action; this hook and any automation must leave it blank. "
            "See docs/frontmatter-hook.md."
        )
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
