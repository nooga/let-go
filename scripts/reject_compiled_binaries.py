#!/usr/bin/env python3
"""Reject staged compiled executables (Mach-O / ELF / PE).

Catches an accidentally-committed `go build` output (e.g. an `lg` binary)
before it bloats history. Detection is magic-byte based, so it does NOT flag
the repo's legitimate binary assets: the core bundle starts with `LGB\x01`
and images (PNG/JPEG) are not executable formats.

pre-commit / prek passes the staged filenames as argv; we inspect each file's
leading bytes and fail if any is an executable. Bypass an intentional commit
with `git commit --no-verify`.
"""
import sys

# 4-byte magic prefixes of executable formats we refuse to commit.
EXECUTABLE_MAGICS = {
    b"\x7fELF": "ELF executable",
    b"\xfe\xed\xfa\xce": "Mach-O executable (32-bit)",
    b"\xfe\xed\xfa\xcf": "Mach-O executable (64-bit)",
    b"\xce\xfa\xed\xfe": "Mach-O executable (32-bit, LE)",
    b"\xcf\xfa\xed\xfe": "Mach-O executable (64-bit, LE)",
    b"\xca\xfe\xba\xbe": "Mach-O universal binary",
    b"\xbe\xba\xfe\xca": "Mach-O universal binary (LE)",
}


def classify(path):
    """Return a human label if `path` is a compiled executable, else None."""
    try:
        with open(path, "rb") as fh:
            head = fh.read(4)
    except (OSError, IsADirectoryError):
        return None
    for magic, label in EXECUTABLE_MAGICS.items():
        if head.startswith(magic):
            return label
    # PE/DOS ("MZ") — only the first two bytes are load-bearing.
    if head[:2] == b"MZ":
        return "PE/DOS executable"
    return None


def main(argv):
    offenders = [(p, label) for p in argv if (label := classify(p))]
    if not offenders:
        return 0
    sys.stderr.write("refusing to commit compiled executable binaries:\n")
    for path, label in offenders:
        sys.stderr.write(f"  {path}  ({label})\n")
    sys.stderr.write(
        "these look like build artifacts — unstage them "
        "(git rm --cached <file>, then gitignore), or commit with --no-verify.\n"
    )
    return 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
