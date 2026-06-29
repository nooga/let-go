#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OUT_DIR="$ROOT/examples/browser-inspector/dist"
OUT_HTML="$OUT_DIR/index.html"
LG="${LG:-$ROOT/lg}"

mkdir -p "$OUT_DIR"
if [ -d "$OUT_HTML" ]; then
  rm -rf "$OUT_HTML"
fi

"$LG" -w "$OUT_DIR" -w-shell none -w-host-eval "$ROOT/examples/browser-inspector/main.lg"

python3 - "$OUT_HTML" "$ROOT/examples/browser-inspector/shell.html" <<'PY'
from pathlib import Path
import re
import sys

bundle = Path(sys.argv[1])
template = Path(sys.argv[2]).read_text()
html = bundle.read_text()
match = re.search(r"<script>\s*(.*?)\s*</script>\s*</body>\s*</html>\s*$", html, re.S)
if not match:
    raise SystemExit("could not extract host script from generated bundle")
host_script = match.group(1).strip()
if "__LG_HOST_SCRIPT__" not in template:
    raise SystemExit("shell template missing __LG_HOST_SCRIPT__ placeholder")
bundle.write_text(template.replace("__LG_HOST_SCRIPT__", host_script, 1))
PY

echo "wrote $OUT_HTML"
