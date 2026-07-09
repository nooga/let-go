#!/usr/bin/env python3
"""Stamp a per-build nonce into the generated index.html.

The browser-inspector app is a single self-contained index.html with the wasm
inlined, so a cached page silently runs old code. Stamping a unique nonce per
build gives every version a distinct, verifiable identity:
  - <meta name="lg-build" content="NONCE"> so tooling/tests can read it
  - a console.log so a human (or an automated check) can confirm which build
    is actually live in the tab

Idempotent: re-stamping replaces any prior stamp rather than accumulating.

Usage: stamp_nonce.py <index.html> <nonce>
"""
import re
import sys

path, nonce = sys.argv[1], sys.argv[2]
html = open(path, encoding="utf-8").read()

# Remove any previous stamp so repeated builds don't accumulate markers.
html = re.sub(r'\n?<meta name="lg-build"[^>]*>', "", html)
html = re.sub(r'\n?<script>console\.log\("lg-build",[^<]*</script>', "", html)

stamp = (
    '\n<meta name="lg-build" content="%s">'
    '\n<script>console.log("lg-build", "%s")</script>' % (nonce, nonce)
)
# Inject right after <head> (the shell template always has one).
if "<head>" not in html:
    raise SystemExit("stamp_nonce: no <head> in %s" % path)
html = html.replace("<head>", "<head>" + stamp, 1)

open(path, "w", encoding="utf-8").write(html)
print("stamped lg-build nonce: %s" % nonce)
