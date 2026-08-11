#!/usr/bin/env python3
"""Fail if mise.toml pins a Go version that disagrees with go.mod.

go.mod is the single repository-owned Go version pin: the Makefile derives
GO-VERSION from it and every CI job uses `go-version-file: go.mod`. mise.toml
therefore carries no `go` entry — a second pin is a copy that drifts, which is
exactly what happened before this was single-sourced (Makefile said 1.26.3,
mise.toml said 1.26.5, go.mod said neither).

Deleting the pin is the intended state and passes silently. Re-adding one is
allowed but must agree with go.mod, so the two can never disagree unnoticed:

  * no `go` entry in mise.toml            -> pass (the normal state)
  * `go` matches go.mod                   -> pass
  * `go` disagrees with go.mod            -> fail, naming both values
  * go.mod states only major.minor        -> a mise patch release under that
                                             minor is accepted, since go.mod
                                             does not pin a patch to match

Stdlib-only, like the other hooks here. Run with no arguments; any filenames
passed by pre-commit are ignored, because the answer depends on both files
regardless of which one was staged.
"""

import os
import re
import sys

TOOLCHAIN = re.compile(r"^toolchain\s+go(\d+\.\d+(?:\.\d+)?)\s*$", re.M)
GO_DIRECTIVE = re.compile(r"^go\s+(\d+\.\d+(?:\.\d+)?)\s*$", re.M)
# [tools] entry, e.g.  go = "1.26.5"  /  go = '1.26'  /  go = "1.26.5"  # note
MISE_GO = re.compile(r"""^\s*go\s*=\s*["']([^"']+)["']""", re.M)


def repo_root() -> str:
    return os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def read(path: str) -> str:
    try:
        with open(path, encoding="utf-8") as fh:
            return fh.read()
    except FileNotFoundError:
        return ""


def gomod_version(text: str) -> tuple[str, str] | tuple[None, None]:
    """Authoritative version and the directive it came from."""
    m = TOOLCHAIN.search(text)
    if m:
        return m.group(1), "toolchain"
    m = GO_DIRECTIVE.search(text)
    if m:
        return m.group(1), "go"
    return None, None


def mise_go_version(text: str) -> str | None:
    m = MISE_GO.search(text)
    return m.group(1) if m else None


def compatible(mise: str, gomod: str) -> bool:
    if mise == gomod:
        return True
    # go.mod pins only major.minor: accept any patch release of that minor.
    if gomod.count(".") == 1 and mise.startswith(gomod + "."):
        return True
    return False


def check(gomod_text: str, mise_text: str) -> list[str]:
    mise = mise_go_version(mise_text)
    if mise is None:
        return []
    gomod, directive = gomod_version(gomod_text)
    if gomod is None:
        return [
            "mise.toml pins go = \"%s\" but go.mod declares no Go version to check it "
            "against (expected a `toolchain goX.Y.Z` or `go X.Y.Z` directive)." % mise
        ]
    if compatible(mise, gomod):
        return []
    return [
        'mise.toml pins go = "%s" but go.mod\'s %s directive says %s.' % (mise, directive, gomod),
        "go.mod is the single source for the Go version: the Makefile derives GO-VERSION",
        "from it and CI uses `go-version-file: go.mod`.",
        "Fix by deleting the `go` line from mise.toml (preferred), or by making it match.",
    ]


def main(argv: list[str]) -> int:
    root = repo_root()
    problems = check(read(os.path.join(root, "go.mod")), read(os.path.join(root, "mise.toml")))
    if not problems:
        return 0
    print("Go version pins disagree:", file=sys.stderr)
    for line in problems:
        print("  " + line, file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
