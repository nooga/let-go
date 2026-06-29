#!/usr/bin/env bash
# Build the host-eval demo: a client-owned-shell WASM bundle, with shell.html
# grafted in to make one page (wasm inlined; xterm + web font load from a CDN).
#
#   LG=lg ./build.sh
#
# Override LG to point at a specific let-go binary (defaults to `lg` on PATH).
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
lg="${LG:-lg}"
out="$here/dist"
page="$out/index.html"

# -w-shell none : emit the runtime + window.LetGoHost with NO built-in UI; the
#                 host page supplies the shell.
# -w-host-eval  : expose LetGoHost.eval(code) to run forms in the live image.
"$lg" -w "$out" -w-shell none -w-host-eval "$here/main.lg"

# Graft shell.html in before </body> so the result is one self-contained page.
# Rebuilds regenerate the page, so this never double-injects.
awk -v shell="$here/shell.html" '
  /<\/body>/ && !done { while ((getline line < shell) > 0) print line; done = 1 }
  { print }
' "$page" > "$page.tmp" && mv "$page.tmp" "$page"

echo "built $page"
echo "serve with cross-origin isolation (COOP: same-origin, COEP: require-corp)"
echo "— the wasm input ring is SharedArrayBuffer-backed and needs crossOriginIsolated."
