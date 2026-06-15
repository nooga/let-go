#!/usr/bin/env python3
"""Fail if any test file lives at the repository root.

Go and let-go test files belong next to the package they test (`pkg/<x>/…`,
`test/`, `test/e2e/`), never at the repo root — root `*_test.go` files force
`package main` test helpers into the binary's package and clutter the top level.
This guard keeps the root free of `*_test.go` and `*_test.lg`.

Run with no arguments to scan the repo root (used by the pre-commit hook and the
CI job). Any filenames passed as arguments are also checked, so it composes with
pre-commit's changed-file list. Exits non-zero and lists offenders on failure.
"""

import glob
import os
import sys

# Patterns considered "test files" for placement purposes.
TEST_GLOBS = ("*_test.go", "*_test.lg")


def repo_root() -> str:
    # scripts/ sits directly under the repo root.
    return os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def root_offenders(root: str) -> list[str]:
    found: list[str] = []
    for pat in TEST_GLOBS:
        found.extend(os.path.basename(p) for p in glob.glob(os.path.join(root, pat)))
    return sorted(found)


def main(argv: list[str]) -> int:
    root = repo_root()
    offenders = set(root_offenders(root))

    # Also honor any explicit paths (pre-commit passes changed files): flag a
    # test file whose directory is the repo root.
    for arg in argv:
        ap = os.path.abspath(arg)
        if os.path.dirname(ap) == root and (
            ap.endswith("_test.go") or ap.endswith("_test.lg")
        ):
            offenders.add(os.path.basename(ap))

    if offenders:
        print("ERROR: test files are not allowed at the repository root:", file=sys.stderr)
        for f in sorted(offenders):
            print(f"  {f}", file=sys.stderr)
        print(
            "\nMove them next to the package they test:\n"
            "  - black-box end-to-end tests (build + run the lg binary) -> test/e2e/\n"
            "  - package-internal/white-box tests                       -> that package's dir\n"
            "  - language behavior tests (.lg)                          -> test/ (TestRunner)\n",
            file=sys.stderr,
        )
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
