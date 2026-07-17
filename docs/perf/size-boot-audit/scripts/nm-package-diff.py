#!/usr/bin/env python3
"""Diff two `go tool nm -size` dumps, bucketed by package.

Answers "which packages grew between two builds". let-go internal packages
are collapsed to `letgo:pkg/<sub>`; stdlib/vendor keep their import path.

Generate the dumps from ELF builds (nm's size field is unreliable on
darwin/macho, so cross-compile to linux):

    GOOS=linux GOARCH=arm64 go build -o /tmp/lg-old . && go tool nm -size /tmp/lg-old > old.txt
    GOOS=linux GOARCH=arm64 go build -o /tmp/lg-new . && go tool nm -size /tmp/lg-new > new.txt
    ./nm-package-diff.py old.txt new.txt

Note: nm-visible symbols are only part of a Go binary; pclntab and some
rodata are not attributed here, so the total delta undercounts the real
file-size delta. Use it for ATTRIBUTION (what grew), not absolute totals.
"""
import sys, re

TYPES = "TtRrDdBbCcUuGgSs?"

def load(path):
    sizes = {}
    for line in open(path):
        t = line.split()
        if len(t) < 3:
            continue
        ti = next((i for i, tok in enumerate(t) if len(tok) == 1 and tok in TYPES), None)
        if not ti:
            continue
        try:
            sz = int(t[ti - 1])
        except ValueError:
            continue
        name = " ".join(t[ti + 1:])
        if name.startswith(("type:", "go:", "gofile:", "runtime.")):
            pkg = re.split(r"[.:]", name)[0]
        else:
            ls = name.rfind("/")
            if ls >= 0:
                dot = name.find(".", ls)
                pkg = name[:dot] if dot > 0 else name
            else:
                dot = name.find(".")
                pkg = name[:dot] if dot > 0 else name
        m = re.search(r"let-go/(pkg/[^/]+)", pkg)
        if m:
            pkg = "letgo:" + m.group(1)
        elif pkg.startswith("github.com/nooga/let-go"):
            pkg = "letgo:root"
        sizes[pkg] = sizes.get(pkg, 0) + sz
    return sizes

def main():
    if len(sys.argv) != 3:
        sys.exit(__doc__)
    a, b = load(sys.argv[1]), load(sys.argv[2])
    rows = sorted(((b.get(k, 0) - a.get(k, 0), a.get(k, 0), b.get(k, 0), k)
                   for k in set(a) | set(b)), reverse=True)
    print(f"{'DELTA':>10} {'old':>10} {'new':>10}  package")
    print("-" * 60)
    for d, av, bv, k in rows:
        if d == 0:
            continue
        print(f"{d:>+10} {av:>10} {bv:>10}  {k}")
    tot = sum(b.values()) - sum(a.values())
    print("-" * 60)
    print(f"nm-visible delta: {tot:+} bytes ({tot/1048576:+.2f} MB)")

if __name__ == "__main__":
    main()
