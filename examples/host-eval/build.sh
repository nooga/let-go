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

# Inline a λ favicon in <head> (the bundle ships none) so browsers don't fall
# back to a /favicon.ico request. base64 SVG: amber λ on the brand dark.
# An icon link only counts in <head>; a body-placed one is ignored.
ico='<link rel="icon" href="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAzMiAzMiI+PHJlY3Qgd2lkdGg9IjMyIiBoZWlnaHQ9IjMyIiByeD0iNiIgZmlsbD0iIzFhMTYxMCIvPjx0ZXh0IHg9IjE2IiB5PSIyNCIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZm9udC1mYW1pbHk9InVpLW1vbm9zcGFjZSxtb25vc3BhY2UiIGZvbnQtc2l6ZT0iMjIiIGZpbGw9IiNmNWIwNDEiPiYjOTU1OzwvdGV4dD48L3N2Zz4=">'
awk -v ico="$ico" '/<\/head>/ && !d { print ico; d=1 } { print }' \
  "$page" > "$page.tmp" && mv "$page.tmp" "$page"

# Graft shell.html in before </body>. Rebuilds regenerate the page from scratch,
# so neither inject ever stacks.
awk -v shell="$here/shell.html" '
  /<\/body>/ && !done { while ((getline line < shell) > 0) print line; done = 1 }
  { print }
' "$page" > "$page.tmp" && mv "$page.tmp" "$page"

echo "built $page"
echo "serve with cross-origin isolation (COOP: same-origin, COEP: require-corp)"
echo "— the wasm input ring is SharedArrayBuffer-backed and needs crossOriginIsolated."
