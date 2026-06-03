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
pkg/rt/core_go_lowered/ir_lower_go/ir_lower_go.go: $(CORE-LG-FILES) $(LGBGEN-SOURCES)
	go run -tags bootstrap ./cmd/lgbgen --target=go

# Regenerate every committed code-gen artifact via the let-go orchestrator
# scripts/generate.lg: the three Go-gen files (op_generated.go,
# ir_bridge_generated.go, ir/data/generated.lg), the core_compiled.lgb bundle,
# and the lowered-Go tree. Requires ./lg, so it builds first. The orchestrator
# uses os/sh + os/exec* to run ./lg on the examples/go-gen sources and to run
# `go run -tags bootstrap ./cmd/lgbgen [--target=go]` for the bundle/lowered
# tree. (Replaces the former generate-ir-{ops,bridge,data}.sh shell scripts.)
generate: build
	./lg scripts/generate.lg --go "$$(command -v go)" --lg ./lg --source-paths examples/go-gen

# Short commit for `-X main.commit` so SHA-pin require-letgo checks can fire on
# `make` builds. Release builds get this from goreleaser; a bare `make` build
# previously reported commit="none", so SHA pins always warn-and-skipped.
# `version` deliberately stays "dev" (no honest release version on an untagged
# build). Falls back to "none" outside a git checkout.
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)

$(LG): $(GO) lg.go pkg/**/* pkg/rt/core_compiled.lgb
	which go
	go build -ldflags="-s -w -X main.commit=$(COMMIT)" -o $@ .

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
# Regenerate the gitignored gogen_ir lowered tree (a build artifact, not
# committed — see check-generated). Any target that builds -tags gogen_ir
# depends on this. Cheap relative to the runs that follow.
.PHONY: lowered
lowered:
	@go run -tags bootstrap ./cmd/lgbgen --target=go >/dev/null

# Default gate (~1 min): the jank suite under BOTH VM variants (bytecode +
# gogen_ir-lowered) + the calibration anchor. This is what CI runs.
bench-ratchet: lowered
	go run ./cmd/bench-ratchet check

bench-ratchet-update: lowered
	go run ./cmd/bench-ratchet update

bench-ratchet-show: lowered
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
bench-ratchet-full: lowered
	go run ./cmd/bench-ratchet -full check

bench-ratchet-full-update: lowered
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

# Register the local git merge driver that resolves pkg/rt/core_compiled.lgb
# conflicts by regenerating the bundle from the merged .lg sources (see
# .gitattributes `merge=lgb` and scripts/git-merge-lgb.sh). A merge driver
# lives in .git/config, which is not shared, so each clone must run this once.
install-hooks:
	git config merge.lgb.name "regenerate core_compiled.lgb from merged .lg sources"
	git config merge.lgb.driver "scripts/git-merge-lgb.sh %O %A %B %L %P"
	@echo "Registered the 'lgb' merge driver. core_compiled.lgb conflicts now auto-regenerate."

# Single gate for every generated artifact. One target to remember, and it
# treats the two artifacts by their actual nature. VCS-agnostic by design:
# it shells out to no `git` (this repo is used with jj, whose secondary
# workspaces have no .git, so a `git diff` gate breaks there).
#
#   * core_compiled.lgb is byte-deterministic, so it gets a CONTENT gate:
#     stash the committed bytes, regenerate, and `cmp`. A difference means the
#     committed bundle was stale. Survives a fresh checkout (an mtime
#     `find -newer` check silently passes after any VCS checkout — in fact the
#     old check-lowered-fresh pointed at a path that no longer exists and had
#     been a silent no-op).
#
#   * core_go_lowered/ (+ the gogen_ir wireup files) is NOT committed — it is
#     a build artifact, regenerated on demand and gitignored. Its self-lower
#     trips the wall-clock *typeinfer-budget-ms*, so the bytes are not
#     reproducible and committing it would churn. Here it is regenerated fresh,
#     then gated behaviorally: it must compile under -tags gogen_ir AND
#     dispatch natively (dce -> NativeFn). gogen_ir consumers (this gate, the
#     parity job, any -tags gogen_ir build) regenerate it first; the untagged
#     build and the shipped bytecode binary never need it.
#
# This is the gate CI runs. After a merge/rebase touching pkg/rt/core/**, run
# `make check-generated` (or `make generate` to refresh, then commit).
check-generated: $(GO)
	@echo ">> bundle: regenerate + verify lockstep (content-based, VCS-agnostic)"
	@cp pkg/rt/core_compiled.lgb pkg/rt/.core_compiled.lgb.committed
	@go run -tags bootstrap ./cmd/lgbgen >/dev/null
	@if cmp -s pkg/rt/core_compiled.lgb pkg/rt/.core_compiled.lgb.committed; then \
		rm -f pkg/rt/.core_compiled.lgb.committed; \
		echo "OK: core_compiled.lgb in lockstep with the .lg sources."; \
	else \
		rm -f pkg/rt/.core_compiled.lgb.committed; \
		echo "ERROR: pkg/rt/core_compiled.lgb is stale — the regenerated bytes differ."; \
		echo "       Run 'make generate' and commit the regenerated bundle."; \
		exit 1; \
	fi
	@echo ">> lowered tree: regenerate (gitignored build artifact), then compile + dispatch natively under -tags gogen_ir"
	@go run -tags bootstrap ./cmd/lgbgen --target=go >/dev/null
	@go build -tags gogen_ir ./...
	@out=$$(printf '(require (quote ir.passes.dce)) (println "DCE-TYPE:" (type ir.passes.dce/dce))' \
	        | go run -tags gogen_ir . /dev/stdin 2>&1); \
	if echo "$$out" | grep -q "DCE-TYPE: let-go.lang.NativeFn"; then \
		echo "OK: core_go_lowered/ compiles + dispatches natively (dce -> NativeFn)."; \
	else \
		echo "FAIL: ir.passes.dce/dce did not dispatch to a NativeFn override"; \
		echo "$$out" | tail -5; \
		exit 1; \
	fi

# Fanout ratchet: gate on the size of the generated -tags gogen_ir lowered tree.
# Gates the byte-sum of ALREADY-TRACKED modules against a percent band; new
# modules are exempt. Baseline is docs/perf/fanout-baseline.edn (committed).
#   fanout-ratchet         regenerate tree, fail on tracked-module bloat > band
#   fanout-ratchet-update  recompute + MIN-merge the baseline (tighten-only)
#   fanout-ratchet-show    print current metrics, write nothing
fanout-ratchet: build
	./lg scripts/fanout-ratchet.lg check --go "$$(command -v go)"

fanout-ratchet-update: build
	./lg scripts/fanout-ratchet.lg update --go "$$(command -v go)"

fanout-ratchet-show: build
	./lg scripts/fanout-ratchet.lg show --go "$$(command -v go)"

# PHONY targets are for ones that have conflicting files/dirs present:
.PHONY: test bench-ratchet bench-ratchet-update bench-ratchet-show install-hooks check-generated fanout-ratchet fanout-ratchet-update fanout-ratchet-show
