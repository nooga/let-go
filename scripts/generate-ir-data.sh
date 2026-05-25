#!/usr/bin/env bash
# Generate pkg/rt/core/ir/data_generated.lg from examples/go-gen/ir_data.lg.
#
# Lisp output (no compile / format step beyond what the generator emits).
# Output is the mechanical accessor surface for the IR data types; the
# rest (constructors, structural mutators, uses cache) stays hand-written
# in pkg/rt/core/ir/data.lg and (def) sources the generated names.

set -euo pipefail

cd "$(dirname "$0")/.."

if [ ! -x ./lg ]; then
  echo "error: ./lg binary not found; run 'make build' first" >&2
  exit 1
fi

OUT=pkg/rt/core/ir/data_generated.lg

./lg -source-paths examples/go-gen examples/go-gen/ir_data.lg > "$OUT"

echo "wrote $OUT"
