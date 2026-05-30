###
# Auto install `go` into ./.cache/local/ if not available.
# Also if `make ... GO-VERSION=1.x.y` is used.

ifneq (,$(or $(if $(shell which go),,1),$(GO-VERSION)))
R := https://github.com/makeplus/makes
M := .cache/makes
$(shell [ -d '$M' ] || git clone -q $R '$M')
include $M/init.mk
# override default Go version with: `make ... GO-VERSION=1.x.y`
GO-VERSION ?= 1.26.3
include $M/go.mk
include $M/shell.mk
endif

# Prefer - to _ for make var names (won't conflict with env vars):
LG := lg
GOLANGCI-LINT := github.com/golangci/golangci-lint/cmd/golangci-lint
REPORT-SCRIPT := scripts/clojure_compat_report.sh

# Resource caps for test invocations. GOMEMLIMIT bounds the Go heap
# (soft cap — runtime aggressively GCs to stay under). GO-TEST-TIMEOUT
# matches CI's per-package timeout so runaway tests die locally too.
# Override on the command line: make test GOMEMLIMIT=4GiB.
GOMEMLIMIT ?= 2GiB
GO-TEST-TIMEOUT ?= 60s

# Standard flags + env for `go test`. Use as: $(GO-TEST-ENV) go test $(GO-TEST-FLAGS) ./...
GO-TEST-ENV := GOMEMLIMIT=$(GOMEMLIMIT)
GO-TEST-FLAGS := -timeout $(GO-TEST-TIMEOUT)


# Start repl by default:
default:: run

run: $(LG)
	./$<

build: $(LG)

# Bundle target. The runtime loads compiled bytecode for the core
# namespaces from this file, NOT from the .lg sources. Anyone editing
# a .lg under pkg/rt/core/ must regenerate the bundle or runtime
# behavior silently diverges from source. This prereq rule makes the
# regeneration automatic — `make test`, `make build`, etc. now keep
# the bundle in lockstep with the .lg sources.
CORE-LG-FILES := $(shell find pkg/rt/core -name '*.lg' -type f 2>/dev/null)
LGBGEN-SOURCES := $(shell find cmd/lgbgen -name '*.go' -type f 2>/dev/null)
pkg/rt/core_compiled.lgb: $(CORE-LG-FILES) $(LGBGEN-SOURCES)
	go run -tags bootstrap ./cmd/lgbgen

# Lowered-Go target. The -tags gogen_ir build path links these generated
# Go files in place of the bytecode-VM IR pipeline. Same staleness story
# as the .lgb bundle: regenerate after editing .lg under pkg/rt/core/ or
# the two engines silently disagree (parity-full diverges on bucket
# hashes even when pass/fail counts match). lower_go.go is the timestamp
# anchor for the whole tree — every regen rewrites it.
pkg/rt/core_go_lowered/ir/lower_go/lower_go.go: $(CORE-LG-FILES) $(LGBGEN-SOURCES)
	go run -tags bootstrap ./cmd/lgbgen --target=go

# Loud check used by CI / pre-commit to detect a bundle that was
# committed in a stale state (e.g. someone ran a manual `go test`
# without going through make and committed the half-stale state).
# Exits non-zero with a clear remediation when any .lg under
# pkg/rt/core is newer than the committed bundle.
check-bundle-fresh:
	@stale=$$(find pkg/rt/core -name '*.lg' -newer pkg/rt/core_compiled.lgb 2>/dev/null); \
	if [ -n "$$stale" ]; then \
		echo "ERROR: pkg/rt/core_compiled.lgb is stale relative to:"; \
		echo "$$stale" | sed 's/^/  /'; \
		echo "Run 'make generate' to regenerate the bundle."; \
		exit 1; \
	fi

# Sibling of check-bundle-fresh for the -tags gogen_ir lowered Go tree.
# parity-full silently fails on bucket hashes if the lgb is fresh but
# core_go_lowered/ is stale (untagged vs gogen_ir run two different
# versions of the IR pipeline). Caught the hard way 2026-05-28.
check-lowered-fresh:
	@stale=$$(find pkg/rt/core -name '*.lg' -newer pkg/rt/core_go_lowered/ir/lower_go/lower_go.go 2>/dev/null); \
	if [ -n "$$stale" ]; then \
		echo "ERROR: pkg/rt/core_go_lowered/ is stale relative to:"; \
		echo "$$stale" | sed 's/^/  /'; \
		echo "Run 'go run -tags bootstrap ./cmd/lgbgen --target=go' to regenerate."; \
		exit 1; \
	fi

generate: $(GO) generate-ir-ops generate-ir-bridge generate-ir-data pkg/rt/core_compiled.lgb pkg/rt/core_go_lowered/ir/lower_go/lower_go.go

# Regenerate pkg/ir/op_generated.go from examples/go-gen/ir_ops.lg.
# Requires ./lg to exist (built by `make build`).
generate-ir-ops: build
	./scripts/generate-ir-ops.sh

# Regenerate pkg/rt/ir_bridge_generated.go from examples/go-gen/ir_bridge.lg.
# Requires ./lg to exist (built by `make build`).
generate-ir-bridge: build
	./scripts/generate-ir-bridge.sh

# Regenerate pkg/rt/core/ir/data/generated.lg from examples/go-gen/ir_data.lg.
# Lisp output (mechanical accessor surface for the IR data types).
# Requires ./lg to exist (built by `make build`).
generate-ir-data: build
	./scripts/generate-ir-data.sh

$(LG): $(GO) lg.go pkg/**/* pkg/rt/core_compiled.lgb
	which go
	go build -ldflags="-s -w" -o $@ .

test: pkg/**/* pkg/rt/core_compiled.lgb $(GO)
	$(GO-TEST-ENV) go test $(GO-TEST-FLAGS) -count=1 -v ./test

clojure-compat-report: $(GO)
	@$(REPORT-SCRIPT)

# Performance ratchet.
#
#   bench-ratchet         compare current benchmarks against the
#                         committed docs/perf/baseline.json; exits
#                         non-zero on any regression > 5% (anchor-
#                         normalized). Suitable for CI.
#   bench-ratchet-update  re-run the sweep and ratchet-merge the
#                         baseline (per-(benchmark, metric) MIN).
#                         The ratchet only tightens; -force bypasses.
#   bench-ratchet-show    run the sweep, print the would-be baseline
#                         JSON to stdout, write nothing. For
#                         spot-checking before deciding to update.
#
# All three are anchor-normalized — see cmd/bench-ratchet/main.go
# and docs/perf/ratchet.md.
# Default gate (~1 min): the jank suite under BOTH VM variants (bytecode +
# gogen_ir-lowered) + the calibration anchor. This is what CI runs.
bench-ratchet:
	go run ./cmd/bench-ratchet check

bench-ratchet-update:
	go run ./cmd/bench-ratchet update

bench-ratchet-show:
	go run ./cmd/bench-ratchet show
  
# Parity checks: untagged vs -tags gogen_ir across jank + ir-stress.
# `parity-check` is the default cadence (~3 min); `parity-quick` for
# pre-commit smoke (~2 sec); `parity-full` for the long check (~5 min).
parity-quick:
	@scripts/gogen-parity.sh --quick

parity-check:
	@scripts/gogen-parity.sh

parity-full:
	@scripts/gogen-parity.sh --full

# Manual deep-dive (~25 min): the entire pkg/vm micro-benchmark fleet plus the
# suite under -tags. Not gated in CI — run by hand when investigating a specific
# regression. Pair with `update` to refresh the full baseline.
bench-ratchet-full:
	go run ./cmd/bench-ratchet -full check

bench-ratchet-full-update:
	go run ./cmd/bench-ratchet -full update

clean:
	$(RM) $(LG)

distclean: clean
ifneq (,$(wildcard .cache))
	chmod -R +w .cache
	$(RM) -r .cache
endif

lint: install-golangci-lint
	golangci-lint run

install-golangci-lint: $(GO)
	which golangci-lint || \
	  GO111MODULE=off go get -u $(GO111MODULE-LINT)

# PHONY targets are for ones that have conflicting files/dirs present:
.PHONY: test bench-ratchet bench-ratchet-update bench-ratchet-show check-bundle-fresh check-lowered-fresh
