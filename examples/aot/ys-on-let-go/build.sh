#!/usr/bin/env bash
# Regenerate the .lg files in this directory from their .ys / .yaml sources.
# Requires `ys` (https://yamlscript.org) in PATH.

set -euo pipefail
cd "$(dirname "$0")"

if ! command -v ys >/dev/null 2>&1; then
  echo "ys not found. Install from https://yamlscript.org/install or skip — checked-in .lg files are already runnable on let-go." >&2
  exit 1
fi

for src in 01_hello.ys 02_fizzbuzz.ys 03_classify.ys; do
  out="${src%.ys}.lg"
  {
    echo ";; Auto-generated from $src via: ys --compile"
    echo ";; Do not edit by hand. Regenerate with: ./build.sh"
    echo "(require '[ys-runtime :refer [say ARGS]])"
    ys --compile "$src"
  } > "$out"
done

# Data-mode YAML: emit a driver around the compiled (% ...) form.
{
  echo ";; Auto-generated from 04_data.yaml via: ys --compile"
  echo ";; Do not edit by hand. Regenerate with: ./build.sh"
  echo "(require '[ys-runtime :refer [%]])"
  echo "(def loaded"
  ys --compile 04_data.yaml
  echo ")"
  echo "(println loaded)"
  echo "(println \"name:\" (get loaded \"name\"))"
  echo "(println \"likes:\" (get loaded \"likes\"))"
} > 04_data.lg

echo "Regenerated 01_hello.lg 02_fizzbuzz.lg 03_classify.lg 04_data.lg"
