#!/usr/bin/env bash
# Regenerate the lowered Go tree, compile it, then lint it.
set -euo pipefail
cd "$(dirname "$0")/.."
go run -tags bootstrap ./cmd/lgbgen --target=go --source-paths pkg/rt/core:pkg/rt/gogen
go build ./pkg/rt/core_go_lowered/...                      # hard-fails on declared-and-not-used / undefined
golangci-lint run --timeout=5m ./pkg/rt/core_go_lowered/...
echo "lowered Go tree: builds + lints clean"
