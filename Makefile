###
# Auto install `go` into ./.cache/local/ if not available.
# Also if `make ... GO-VERSION=1.x.y` is used.

ifneq (,$(or $(if $(shell which go),,1),$(GO-VERSION)))
R := https://github.com/makeplus/makes
M := .cache/makes
# Pin makeplus/makes so a HEAD change there can't break `make` with no change
# in this repo. Same rationale as the Clojure toolchain pin in go.yml
# (unpinned makeplus install → setup-clojure). Bump MAKES-SHA after verifying
# the new tip; delete .cache/makes first so the next `make` reclones.
# The trailing `rm -rf` matters: $(shell ...) discards exit status, so without
# it a failed checkout would leave a valid-looking clone at floating HEAD — the
# exact thing this pin defends against — and the [ -d ] guard would make that
# state stick forever. Removing it fails the include loudly and retries next run.
MAKES-SHA := 7f28494955c20e5a87ca82ae351464842a7236de
$(shell [ -d '$M' ] || (git clone -q $R '$M' && git -C '$M' -c advice.detachedHead=false checkout -q $(MAKES-SHA)) || rm -rf '$M')
include $M/init.mk
# go.mod is the single repository-owned Go version pin: prefer the toolchain
# directive; `go mod tidy` erases it when it equals the go directive, so fall
# back to a full-patch go directive. Anything weaker fails loudly rather than
# guessing. Override only for an explicit one-off bootstrap:
# `make ... GO-VERSION=1.x.y`.
GO-VERSION ?= $(shell sed -n 's/^toolchain go//p' go.mod)
ifeq ($(strip $(GO-VERSION)),)
GO-VERSION := $(shell sed -n 's/^go \([0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\)$$/\1/p' go.mod)
endif
ifeq ($(strip $(GO-VERSION)),)
$(error cannot derive GO-VERSION: go.mod needs 'toolchain goX.Y.Z' or a full 'go X.Y.Z' directive)
endif
include $M/go.mk
include $M/shell.mk
endif

# Prefer - to _ for make var names (won't conflict with env vars):
# Build outputs live under build/; bin/ holds the promoted current copy that a
# PATH entry or a script can point at. Both are gitignored as directories, so a
# new binary needs no new ignore entry, and `make clean` removes both.
BUILD-DIR := build
BIN-DIR := bin
# $(LG) is the freshly built CANDIDATE. Gates and ratchets below run against it
# deliberately: they validate the compiler that was just built, not the last
# copy that happened to pass promotion. $(LG-PROMOTED) is the promoted copy.
LG := $(BUILD-DIR)/lg
LG-PROFILE ?= $(BUILD-DIR)/lg-profile
LG-PROMOTED := $(BIN-DIR)/lg
BOOTPROBE := $(BUILD-DIR)/bootprobe
BOOTPROBE-SOURCES := $(wildcard cmd/bootprobe/*.go)
SMOKE-SOURCES := scripts/smoke.lg scripts/smoke-boot.sh
# Boot budget sits ~2x above the measured median-of-5 ceiling (4.10ms on an
# idle M3) so it does not flake, while still catching the #663 class.
SMOKE-BOOT-BUDGET-MS ?= 8
SMOKE-BOOT-SAMPLES ?= 5
GOLANGCI-LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint
GOLANGCI-LINT-VERSION ?= v2.12.2
GOLANGCI-LINT-VERSION-NO-V := $(patsubst v%,%,$(GOLANGCI-LINT-VERSION))
GOLANGCI-LINT-BIN := $(CURDIR)/.cache/local/bin/golangci-lint
GOLANGCI-LINT-CACHE := $(CURDIR)/.cache/local/golangci-lint
GOLANGCI-LINT-GOENV := GOPATH=$(CURDIR)/.cache/local/go GOBIN=$(CURDIR)/.cache/local/bin GOMODCACHE=$(CURDIR)/.cache/local/go/pkg/mod GOCACHE=$(CURDIR)/.cache/local/cache/go-build
REPORT-SCRIPT := scripts/clojure_compat_report.sh

# Resource caps for test invocations. GOMEMLIMIT bounds the Go heap
# (soft cap — runtime aggressively GCs to stay under).
# Override on the command line: make test GOMEMLIMIT=4GiB.
GOMEMLIMIT ?= 2GiB
# `make test` runs the full suite, including the expensive lowering harnesses.
# The determinism test compiles the whole stdlib twice concurrently (~60-120s
# on CI, ~200s on a contended local machine). CI's fast lanes use -short at 60s
# and run these in a dedicated long-timeout step; the full local run keeps a
# generous per-package ceiling matching that step.
GO-TEST-TIMEOUT ?= 600s

# Export the heap cap to EVERY recipe's environment, not just `make test`.
# The bootstrap/lowering targets (generate, lowered, parity) shell out to
# `go run -tags bootstrap` / `go test`, which compile the whole .lg stdlib
# from source and balloon the heap; uncapped, parallel invocations OOM a
# 16GB machine. `export` makes GOMEMLIMIT visible to those go subprocesses
# too. Sub-make/scripts inherit it unless they override.
export GOMEMLIMIT

# Standard flags + env for `go test`. Use as: $(GO-TEST-ENV) go test $(GO-TEST-FLAGS) ./...
# GO-TEST-ENV is retained for explicitness at test call sites; the value is
# already exported above, so it is now belt-and-suspenders.
GO-TEST-ENV := GOMEMLIMIT=$(GOMEMLIMIT)
GO-TEST-FLAGS := -timeout $(GO-TEST-TIMEOUT)


.PHONY: all default run build

all: build

# Start repl by default:
default:: run

run: $(LG)
	./$<

# Promotion is part of the build rule rather than a separate step, so bin/lg
# cannot silently lag build/lg.
build: $(LG-PROMOTED)

# Promotion gate. build/lg is the CANDIDATE; bin/lg is the promoted copy.
# A candidate is promoted only after `smoke` passes, so bin/lg means "booted,
# booted fast, and not obviously broken" — NOT "verified correct". Deep
# correctness is make test / ir-stress-gate / gogen-diff, each ~2 min, which is
# too slow to run on every promotion.
$(LG-PROMOTED): $(LG) $(BOOTPROBE) $(SMOKE-SOURCES)
	@$(MAKE) --no-print-directory smoke
	@mkdir -p $(BIN-DIR)
	install -m 0755 $(LG) $@
	@echo "promoted $(LG) -> $@"

$(BOOTPROBE): $(GO) $(BOOTPROBE-SOURCES) $(ROOT-GO-FILES) pkg/**/* pkg/rt/core_compiled.lgb
	@mkdir -p $(BUILD-DIR)
	go build -o $@ ./cmd/bootprobe

# Fast promotion gate: minimal correctness + a boot-time budget. Seconds, not
# minutes. Override the budget with: make smoke SMOKE-BOOT-BUDGET-MS=12
.PHONY: smoke
smoke: $(LG) $(BOOTPROBE)
	@echo ">> smoke: correctness"
	@$(LG) scripts/smoke.lg
	@echo ">> smoke: boot budget"
	@./scripts/smoke-boot.sh $(BOOTPROBE) $(SMOKE-BOOT-BUDGET-MS) $(SMOKE-BOOT-SAMPLES)

build-profile: $(LG-PROFILE)

# Bundle target. The runtime loads compiled bytecode for the core
# namespaces from this file, NOT from the .lg sources. Anyone editing
# a .lg under pkg/rt/core/ must regenerate the bundle or runtime
# behavior silently diverges from source. This prereq rule makes the
# regeneration automatic — `make test`, `make build`, etc. now keep
# the bundle in lockstep with the .lg sources.
CORE-LG-FILES := $(shell find pkg/rt/core -name '*.lg' -type f 2>/dev/null)
LGBGEN-SOURCES := $(shell find cmd/lgbgen -name '*.go' -type f 2>/dev/null)
ROOT-GO-FILES := $(shell find . -maxdepth 1 -name '*.go' -type f 2>/dev/null)

# Primitive registrar. lgprimgen scans pkg/rt's //lg:native-annotated Go sources
# and emits this file. It is a PURE go/ast codegen tool (no pkg/compiler / pkg/rt
# import), so it MUST run before the lgbgen bootstrap that consumes it: when a
# primitive is hoisted out of installLangNS into a //lg:native decl, the bootstrap
# runtime cannot boot until this file carries its registration. A runtime-coupled
# generator here would depend on the very artifact it produces. Making the bundle
# depend on it breaks that cycle.
PRIMGEN-SOURCES := $(shell find pkg/rt -maxdepth 1 -name '*.go' -type f -not -name 'zz_*' -not -name '*_test.go' 2>/dev/null) $(shell find cmd/lgprimgen internal/primgen -name '*.go' -type f 2>/dev/null)
pkg/rt/zz_primitives_generated.go: $(PRIMGEN-SOURCES) $(GO)
	go run ./cmd/lgprimgen -primitives pkg/rt -go-pkg github.com/nooga/let-go/pkg/rt -primitives-out pkg/rt/zz_primitives_generated.go

COREFNS-SOURCES := $(shell find pkg/rt/corefns -maxdepth 1 -name '*.go' -type f -not -name 'zz_*' -not -name '*_test.go' 2>/dev/null) $(shell find cmd/lgprimgen internal/primgen -name '*.go' -type f 2>/dev/null)
pkg/rt/corefns/zz_primitives_generated.go: $(COREFNS-SOURCES) $(GO)
	go run ./cmd/lgprimgen -primitives pkg/rt/corefns -go-pkg github.com/nooga/let-go/pkg/rt/corefns -primitives-out pkg/rt/corefns/zz_primitives_generated.go

pkg/rt/core_compiled.lgb: pkg/rt/zz_primitives_generated.go pkg/rt/corefns/zz_primitives_generated.go $(CORE-LG-FILES) $(LGBGEN-SOURCES) $(GO)
	go run -tags bootstrap ./cmd/lgbgen

# Lowered-Go target. The -tags gogen_ir build path links these generated
# Go files in place of the bytecode-VM IR pipeline. Same staleness story
# as the .lgb bundle: regenerate after editing .lg under pkg/rt/core/ or
# the two engines silently disagree (parity-full diverges on bucket
# hashes even when pass/fail counts match). lower_go.go is the timestamp
# anchor for the whole tree — every regen rewrites it.
pkg/rt/core_go_lowered/ir/lower_go/lower_go.go: pkg/rt/zz_primitives_generated.go pkg/rt/corefns/zz_primitives_generated.go $(CORE-LG-FILES) $(LGBGEN-SOURCES) $(GO)
	go run -tags bootstrap ./cmd/lgbgen --target=go

# Regenerate every committed code-gen artifact via the let-go orchestrator
# scripts/generate.lg: the three Go-gen files (op_generated.go,
# ir_bridge_generated.go, ir/data/generated.lg), the core_compiled.lgb bundle,
# and the lowered-Go tree. Requires ./lg, so it builds first. The orchestrator
# uses os/sh + os/exec* to run ./lg on the IR-gen specs (which require the
# embedded gogen macros, #425) and to run
# `go run -tags bootstrap ./cmd/lgbgen [--target=go]` for the bundle/lowered
# tree. (Replaces the former generate-ir-{ops,bridge,data}.sh shell scripts.)
generate: build
	$(LG) scripts/generate.lg --go "$$(command -v go)" --lg $(LG)

# Short commit for `-X main.commit` so SHA-pin require-letgo checks can fire on
# `make` builds. Release builds get this from goreleaser; a bare `make` build
# previously reported commit="none", so SHA pins always warn-and-skipped.
# `version` deliberately stays "dev" (no honest release version on an untagged
# build). Falls back to "none" outside a git checkout.
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
PERF-TIMELINE-DIR ?= docs/perf/timeline
PERF-SNAPSHOT ?= $(PERF-TIMELINE-DIR)/$(shell date -u +%Y%m%dT%H%M%SZ)-$(COMMIT).json

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

$(LG): $(GO) $(ROOT-GO-FILES) pkg/**/* pkg/rt/core_compiled.lgb
	which go
	@mkdir -p $(BUILD-DIR)
	go build -ldflags="-s -w -X main.commit=$(COMMIT)" -o $@ .

$(LG-PROFILE): $(GO) $(ROOT-GO-FILES) pkg/**/* pkg/rt/core_compiled.lgb
	which go
	@mkdir -p $(BUILD-DIR)
	go build -tags lg_profile -ldflags="-s -w -X main.commit=$(COMMIT)" -o $@ .

# -short skips the heavyweight e2e tests that each shell out to a full core
# lowering (TestLoweringDeterminism ~150s, TestGogenAOTDiff, TestDeftypeSkeletonNativeLowering).
# They are testing.Short()-gated and run full in dedicated CI jobs
# (.github/workflows/go.yml "Expensive lowering e2e" + the gogen-diff job), so
# the local `make test` loop stays fast without losing CI coverage.

.PHONY: test-gogen-diff-gate
test-gogen-diff-gate:
	@scripts/test-gogen-diff-short.sh

test: pkg/**/* pkg/rt/core_compiled.lgb $(GO)
	@scripts/test-gogen-diff-short.sh
	$(GO-TEST-ENV) go test $(GO-TEST-FLAGS) -short -count=1 -v ./test/...

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
perf-page: $(GO)
	go run ./cmd/perf-page -out docs/perf/index.html

perf-snapshot: lowered $(GO)
	mkdir -p $(PERF-TIMELINE-DIR)
	go run ./cmd/bench-ratchet -full -baseline $(PERF-SNAPSHOT) snapshot

# Regenerate the gitignored gogen_ir lowered tree (a build artifact, not
# committed — see check-generated). Any target that builds -tags gogen_ir
# depends on this. Cheap relative to the runs that follow.
.PHONY: lowered
lowered: $(GO)
	@go run -tags bootstrap ./cmd/lgbgen --target=go >/dev/null

# Self-hosting (3b + 4): re-lower the whole core with the NATIVE compiler and
# verify the fixpoint. `lowered` (3a) bootstraps the tree from source with the
# interpreted passes; here we build lgbgen -tags "bootstrap gogen_ir" so the
# lowering pipeline runs on THAT tree's native passes, re-lower everything
# (including the passes themselves), and byte-compare against the committed
# tree. Byte-identical = the native compiler reproduces its own output, so the
# passes are correct native code AND fast enough to self-host (bench-ratchet
# IRCompile: native ~16% faster than bytecode). A diff means a pass changed and
# the tree is stale relative to source — re-run `make lowered` (3a) first.
.PHONY: check-selfhost
check-selfhost: lowered $(GO)
	@echo ">> 3b: build native-passes lgbgen (-tags 'bootstrap gogen_ir')"
	@mkdir -p $(BUILD-DIR)
	@go build -tags "bootstrap gogen_ir" -o $(BUILD-DIR)/selfhost-lgbgen ./cmd/lgbgen
	@echo ">> 3b: re-lower the whole core with the NATIVE compiler"
	@tmp=$$(mktemp -d); $(BUILD-DIR)/selfhost-lgbgen --target=go "$$tmp" >/dev/null; \
	echo ">> 4: fixpoint — native re-lowering vs the committed tree"; \
	if diff -rq "$$tmp" pkg/rt/core_go_lowered >/dev/null 2>&1; then \
		echo "OK: native self-hosting fixpoint holds (byte-identical)."; rc=0; \
	else \
		echo "FAIL: native re-lowering differs — a pass changed; run 'make lowered' first:"; \
		diff -rq "$$tmp" pkg/rt/core_go_lowered 2>&1 | head; rc=1; \
	fi; \
	rm -rf "$$tmp" $(BUILD-DIR)/selfhost-lgbgen; exit $$rc

# Differential engine-output gate: first run the complete-corpus strict
# capability/parity census, then retain the existing non-strict differential
# check. Both mandatory legs force full mode even when GOFLAGS contains -short.
.PHONY: gogen-diff
gogen-diff: engine-parity-gate
	go test -short=false -run TestGogenAOTDiff -count=1 -v ./test/e2e/

# Complete-corpus engine-output parity and strict-capability gate. The strict
# wrapper measures whether each fixture can execute with fallback forbidden; it
# does not prove fixture-source generated-Go or native-entry execution.
# The gate always logs a deterministic sorted coverage report: per-fixture
# category plus the measured lowering reason behind every non-full fixture, and
# a bucket census of distinct reasons. Report content is informational: it never
# changes classification, ratchets, or rederive. The one exit-status coupling is
# deliberate — an explicitly requested report write that fails is reported as a
# test error. Save a copy with:
#   LETGO_PARITY_REPORT=/tmp/parity-report.txt make engine-parity-gate
# The gate also enforces the committed case ledger test/parity-ledger.txt: a
# diffable record of every measured case plus tombstones for removed fixtures.
# A normal run fails with per-case ADDED / REMOVED / STATUS-CHANGED lines, on a
# reappearing tombstoned fixture, and on any UNEXPLAINED tombstone.
# Re-derive the ledger and both neutral shrink-only baselines (written as one
# rollback-safe set) after a reviewed improvement with:
#   LETGO_PARITY_REDERIVE=1 make engine-parity-gate
# Explain fixtures that disappear in the same rederive; unexplained tombstones
# fail later runs:
#   LETGO_PARITY_TOMBSTONE_REASON="a.lg=merged into b.lg;c.lg=obsolete" \
#     LETGO_PARITY_REDERIVE=1 make engine-parity-gate
# Drop tombstones inherited from earlier revisions (keeping this revision's) with:
#   LETGO_PARITY_LEDGER_CLEANUP=1 LETGO_PARITY_REDERIVE=1 make engine-parity-gate
# A revision is the jj change id when available, else the git commit; when the
# revision cannot be resolved, cleanup is skipped and nothing is dropped.
.PHONY: engine-parity-gate
engine-parity-gate: lowered $(GO)
	go test -short=false -run TestEngineParityGate -count=1 -v ./test/e2e/

# Historical compatibility aliases for existing local automation. The names
# predate the complete result-driven census; neither denotes a phase-limited or
# separately curated audit tier.
.PHONY: parity-gate-phase1 strict-audit
parity-gate-phase1 strict-audit: engine-parity-gate

# Native-entry AST gate: for every fixture in test/native-entry/, lower the
# fixture through the production path (scripts/lg-compile --entry-frame),
# assert the generated Go's structure with internal/gofragment + go/ast, build
# a real binary, and require its stdout to equal the committed <fixture>.expect.
# Three AST-located mutants (oracle / call site / returned value) must each
# fail, so a green run cannot be vacuous.
#
# The same target also runs TestJankSuiteDirectABIGeneratedGo: the pinned jank
# `identical?` deftest is lowered to a real Go test package (generated into a
# temp dir, never into the repo), matched against an inline Go AST oracle (plus
# its falsifier), executed with `go test`, and killed by an AST-located
# semantic mutant. It is gated here rather than left
# incidental. It needs test/clojure-test-suite; that submodule is only
# materialized in the primary git worktree, so in a jj workspace run
# `scripts/link-clojure-test-suite.sh <workspace>` first. The harness FAILS
# loudly (never skips) when the suite is missing.
#
# The gate is testing.Short()-gated (it builds and runs binaries), so it must
# run with -short=false. An ambient GOFLAGS=-short — exported by fast CI lanes
# and by some local shells — would otherwise turn the whole gate into a silent
# skip: strip it from GOFLAGS AND pass -short=false explicitly (the later flag
# wins over anything GOFLAGS injects).
.PHONY: native-entry-gate
native-entry-gate: $(GO)
	GOFLAGS="$(filter-out -short -test.short,$(GOFLAGS))" $(GO-TEST-ENV) go test $(GO-TEST-FLAGS) -run 'TestNativeEntryASTGate|TestJankSuiteDirectABIGeneratedGo' -short=false -count=1 -v ./test/e2e/

.PHONY: bench-ratchet

# Default gate (~1 min): the jank suite under BOTH VM variants (bytecode +
# gogen_ir-lowered) + the calibration anchor. This is what CI runs.
bench-ratchet: lowered $(GO)
	go run ./cmd/bench-ratchet check

bench-ratchet-update: lowered $(GO)
	go run ./cmd/bench-ratchet update

bench-ratchet-show: lowered $(GO)
	go run ./cmd/bench-ratchet show

# Parity checks: untagged vs -tags gogen_ir across jank + ir-stress.
# `parity-check` is the default cadence (~3 min); `parity-quick` for
# pre-commit smoke (~2 sec); `parity-full` for the long check (~5 min).
parity-quick: $(GO)
	@scripts/gogen-parity.sh --quick

parity-check: $(GO)
	@scripts/gogen-parity.sh

parity-full: $(GO)
	@scripts/gogen-parity.sh --full

# Manual deep-dive (~25 min): the pkg/vm fleet plus suite/IR variants. Not
# gated in PR CI — run by hand when investigating a specific regression. Pair
# with `update` to refresh the full baseline.
bench-ratchet-full: lowered $(GO)
	go run ./cmd/bench-ratchet -full check

bench-ratchet-full-update: lowered $(GO)
	go run ./cmd/bench-ratchet -full update

.PHONY: clean-lowered clean distclean

clean-lowered:
	$(RM) -r pkg/rt/core_go_lowered
	$(RM) lg_gogen_ir.go lg_gogen_accel.go cmd/lgbgen/main_gogen_ir.go cmd/lgbgen/main_gogen_accel.go pkg/ir/zz_gogen_ir_wire_test.go pkg/ir/zz_gogen_accel_wire_test.go pkg/wasmhost/zz_gogen_ir_wire_test.go pkg/rt/generated.provenance
	@echo "Cleaned lowered Go tree and wireup files"

clean: clean-lowered
	$(RM) -r $(BUILD-DIR) $(BIN-DIR)

distclean: clean
ifneq (,$(wildcard .cache))
	chmod -R +w .cache
	$(RM) -r .cache
endif

lint: install-golangci-lint
	GOLANGCI_LINT_CACHE=$(GOLANGCI-LINT-CACHE) $(GOLANGCI-LINT-BIN) run

install-golangci-lint: $(GO)
	@mkdir -p $(dir $(GOLANGCI-LINT-BIN)) $(GOLANGCI-LINT-CACHE)
	@if ! test -x $(GOLANGCI-LINT-BIN) || ! $(GOLANGCI-LINT-BIN) --version | grep -q 'version $(GOLANGCI-LINT-VERSION-NO-V)'; then \
	  $(GOLANGCI-LINT-GOENV) $(GO) install $(GOLANGCI-LINT)@$(GOLANGCI-LINT-VERSION); \
	fi

# Register the local git merge drivers for the generated artifacts (see
# .gitattributes). Merge drivers live in .git/config, which is not shared, so
# each clone must run this once.
#   * lgb   — regenerates pkg/rt/core_compiled.lgb from the merged .lg sources
#             (scripts/git-merge-lgb.sh).
#   * sums  — recomputes pkg/rt/generated.sums' digest from the merged sources
#             (scripts/git-merge-sums.sh) so the signature is never stale.
# NOTE: git merge drivers only fire on the plain-git merge path. jj does not run
# them, so under jj both artifacts are reconciled by `make generate` post-rebase.
install-hooks:
	git config merge.lgb.name "regenerate core_compiled.lgb from merged .lg sources"
	git config merge.lgb.driver "scripts/git-merge-lgb.sh %O %A %B %L %P"
	git config merge.sums.name "recompute generated.sums digest from merged sources"
	git config merge.sums.driver "scripts/git-merge-sums.sh %O %A %B %L %P"
	@echo "Registered the 'lgb' + 'sums' merge drivers (core_compiled.lgb / generated.sums)."

# Non-mutating front gate: fail before any target has a chance to refresh
# pkg/rt/generated.sums. This catches the exact drift that a prior go generate
# or lgbgen invocation would otherwise mask in CI.
check-generated-manifest: $(GO)
	@go run ./cmd/check-generated || { \
		echo "ERROR: dependency manifest stale or check errored — run 'make generate'."; \
		exit 1; \
	}

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
#     uses *typeinfer-max-drains* (a deterministic drain-count, not wall-clock),
#     ensuring the bytes are reproducible across processes. Here it is regenerated
#     fresh, then gated behaviorally: it must compile under -tags gogen_ir AND
#     dispatch natively (dce -> NativeFn). gogen_ir consumers (this gate, the
#     parity job, any -tags gogen_ir build) regenerate it first; the untagged
#     build and the shipped bytecode binary never need it.
#
# Default-build dependency gate: an untagged build must not link the optional
# heavy subsystems. See scripts/check-default-deps.sh for why the patterns are
# anchored — the stdlib's vendored x/text copies would otherwise trip it.
.PHONY: check-default-deps
check-default-deps: $(GO)
	@scripts/check-default-deps.sh

# This is the gate CI runs. After a merge/rebase touching pkg/rt/core/**, run
# `make check-generated` (or `make generate` to refresh, then commit).
check-generated: check-generated-manifest $(GO)
	@echo ">> regenerate bundle + lowered tree from a SINGLE core compile (--target=both)"
	@cp pkg/rt/core_compiled.lgb pkg/rt/.core_compiled.lgb.committed
	@cp pkg/rt/corefns/zz_primitives_generated.go pkg/rt/corefns/.zz_primitives_generated.go.committed
	@cp pkg/rt/zz_primitives_generated.go pkg/rt/.zz_primitives_generated.go.committed
	@go run ./cmd/lgprimgen -primitives pkg/rt -go-pkg github.com/nooga/let-go/pkg/rt -primitives-out pkg/rt/zz_primitives_generated.go
	@go run ./cmd/lgprimgen -primitives pkg/rt/corefns -go-pkg github.com/nooga/let-go/pkg/rt/corefns -primitives-out pkg/rt/corefns/zz_primitives_generated.go
	@go run -tags bootstrap ./cmd/lgbgen --target=both >/dev/null
	@echo ">> rt primitives: verify lockstep (content-based, VCS-agnostic)"
	@if cmp -s pkg/rt/zz_primitives_generated.go pkg/rt/.zz_primitives_generated.go.committed; then \
		rm -f pkg/rt/.zz_primitives_generated.go.committed; \
		echo "OK: pkg/rt/zz_primitives_generated.go in lockstep with source."; \
	else \
		rm -f pkg/rt/.zz_primitives_generated.go.committed; \
		echo "ERROR: pkg/rt/zz_primitives_generated.go is stale — the regenerated bytes differ."; \
		echo "       Run 'make generate' and commit the regenerated file."; \
		exit 1; \
	fi
	@echo ">> corefns registrar: verify lockstep (content-based, VCS-agnostic)"
	@if cmp -s pkg/rt/corefns/zz_primitives_generated.go pkg/rt/corefns/.zz_primitives_generated.go.committed; then \
		rm -f pkg/rt/corefns/.zz_primitives_generated.go.committed; \
		echo "OK: pkg/rt/corefns/zz_primitives_generated.go in lockstep with source."; \
	else \
		rm -f pkg/rt/corefns/.zz_primitives_generated.go.committed; \
		echo "ERROR: pkg/rt/corefns/zz_primitives_generated.go is stale — the regenerated bytes differ."; \
		echo "       Run 'make generate' and commit the regenerated file."; \
		exit 1; \
	fi
	@echo ">> bundle: verify lockstep (content-based, VCS-agnostic)"
	@if cmp -s pkg/rt/core_compiled.lgb pkg/rt/.core_compiled.lgb.committed; then \
		rm -f pkg/rt/.core_compiled.lgb.committed; \
		echo "OK: core_compiled.lgb in lockstep with the .lg sources."; \
	else \
		rm -f pkg/rt/.core_compiled.lgb.committed; \
		echo "ERROR: pkg/rt/core_compiled.lgb is stale — the regenerated bytes differ."; \
		echo "       Run 'make generate' and commit the regenerated bundle."; \
		exit 1; \
	fi
	@echo ">> lowered tree: compile + dispatch natively under -tags gogen_ir"
	@go build -tags gogen_ir ./pkg/rt/...
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
	$(LG) scripts/fanout-ratchet.lg check --go "$$(command -v go)"

fanout-ratchet-update: build
	$(LG) scripts/fanout-ratchet.lg update --go "$$(command -v go)"

fanout-ratchet-show: build
	$(LG) scripts/fanout-ratchet.lg show --go "$$(command -v go)"

# IR-stress: lower-go AOT pass-rate over the committed corpus allow-list
# (scripts/ir-stress-corpus.edn = every shipped + test/example/script .lg minus
# :exclude). Failures are real lowering gaps. Env overridable: LG_STRESS_PASSES
# (default 1), LG_STRESS_TIMEOUT_MS (15000), LG_STRESS_LOG (/tmp/ir-stress.log).
ir-stress: build
	LG_STRESS_PASSES=$${LG_STRESS_PASSES:-1} \
	  LG_STRESS_TIMEOUT_MS=$${LG_STRESS_TIMEOUT_MS:-15000} \
	  LG_STRESS_LOG=$${LG_STRESS_LOG:-/tmp/ir-stress.log} \
	  $(LG) scripts/ir-stress.lg corpus scripts/ir-stress-corpus.edn

# Jank lowering-coverage gate: lower-go AOT pass-rate over the vendored jank
# Clojure-compat suite (test/clojure-test-suite, a git submodule). Unlike the
# internal corpus, these .cljc files exercise the broader Clojure surface, so
# the buckets surface lowering gaps the internal corpus can't (e.g. BigDecimal
# literals, multimethods). LG_SOURCE_PATHS lists the repo's compat shim FIRST so
# its portability.lg shadows the suite's own portability.cljc. Fixtures are the
# macro-generated test fns (deftest bodies), enumerated by the pipeline's
# canonical lowerable-fn-forms. Env overridable like ir-stress.
JANK_SUITE_DIR := test/clojure-test-suite/test/clojure
jank-stress: build
	LG_SOURCE_PATHS="test/compat:test/clojure-test-suite/test" \
	  LG_STRESS_PASSES=$${LG_STRESS_PASSES:-1} \
	  LG_STRESS_TIMEOUT_MS=$${LG_STRESS_TIMEOUT_MS:-15000} \
	  LG_STRESS_LOG=$${LG_STRESS_LOG:-/tmp/jank-lowering.log} \
	  $(LG) scripts/ir-stress.lg lower-go $(JANK_SUITE_DIR) \
	    $$(cd $(JANK_SUITE_DIR) && ls core_test/*.cljc string_test/*.cljc)

# ITER-0021 lowering-coverage ratchet gate. Runs the ir-stress corpus and fails
# (exit 1) if native-lowering failures grew vs docs/perf/ir-stress-baseline.edn —
# a form newly failing to lower (whole-function fallback to bytecode). Reuses the
# ir-stress fallback census; the ratchet only tightens. Automates BE-2.
ir-stress-gate: build
	LG_STRESS_PASSES=1 \
	  LG_STRESS_TIMEOUT_MS=$${LG_STRESS_TIMEOUT_MS:-15000} \
	  LG_STRESS_BASELINE=docs/perf/ir-stress-baseline.edn \
	  $(LG) scripts/ir-stress.lg corpus scripts/ir-stress-corpus.edn

# Bytecode-path census. Same corpus, but lowers through ir.lower (the backend
# `(set! *ir-compile* true)` actually drives at runtime) instead of lower_go.
# The lower-go ratchet is NOT a proxy for this: the two backends fail on
# different things — lower.lg materializes by stack position, lower_go.lg works
# off structurize's control tree, which lower.lg doesn't use — so a form can
# lower natively to Go and still fail to bytecode. Separate baseline file so
# the two ratchets can't overwrite each other.
ir-stress-bytecode: build
	LG_STRESS_PASSES=$${LG_STRESS_PASSES:-1} \
	  LG_STRESS_TIMEOUT_MS=$${LG_STRESS_TIMEOUT_MS:-15000} \
	  LG_STRESS_BACKEND=lower \
	  LG_STRESS_LOG=$${LG_STRESS_LOG:-/tmp/ir-stress-bytecode.log} \
	  $(LG) scripts/ir-stress.lg corpus scripts/ir-stress-corpus.edn

ir-stress-bytecode-gate: build
	LG_STRESS_PASSES=1 \
	  LG_STRESS_TIMEOUT_MS=$${LG_STRESS_TIMEOUT_MS:-15000} \
	  LG_STRESS_BACKEND=lower \
	  LG_STRESS_BASELINE=docs/perf/ir-stress-bytecode-baseline.edn \
	  $(LG) scripts/ir-stress.lg corpus scripts/ir-stress-corpus.edn

ir-stress-bytecode-rebaseline: build
	LG_STRESS_PASSES=1 \
	  LG_STRESS_TIMEOUT_MS=$${LG_STRESS_TIMEOUT_MS:-15000} \
	  LG_STRESS_BACKEND=lower \
	  LG_STRESS_BASELINE=docs/perf/ir-stress-bytecode-baseline.edn \
	  LG_STRESS_REBASELINE=1 \
	  LG_STRESS_DATE=$$(date +%F) \
	  $(LG) scripts/ir-stress.lg corpus scripts/ir-stress-corpus.edn

# Rewrite the committed coverage baseline from a fresh census (tool-maintained;
# never hand-edit the EDN). Run after an intentional corpus or coverage change,
# review the diff, and commit it with the change that caused it.
ir-stress-rebaseline: build
	LG_STRESS_PASSES=1 \
	  LG_STRESS_TIMEOUT_MS=$${LG_STRESS_TIMEOUT_MS:-15000} \
	  LG_STRESS_BASELINE=docs/perf/ir-stress-baseline.edn \
	  LG_STRESS_REBASELINE=1 \
	  LG_STRESS_DATE=$$(date +%F) \
	  $(LG) scripts/ir-stress.lg corpus scripts/ir-stress-corpus.edn

# Combined speed + size gates. Both ratchets need the gogen_ir lowered tree, and
# each would otherwise regenerate it (the dominant cost). `ratchets` regenerates
# it ONCE via `lowered`, runs the speed gate against it, then runs the size gate
# with --no-regen so it reuses the same tree — ~halving wall time vs running
# `make bench-ratchet fanout-ratchet`. Use this in CI.
ratchets: build lowered $(GO)
	go run ./cmd/bench-ratchet check
	$(LG) scripts/fanout-ratchet.lg check --go "$$(command -v go)" --no-regen

ratchets-update: build lowered $(GO)
	go run ./cmd/bench-ratchet update
	$(LG) scripts/fanout-ratchet.lg update --go "$$(command -v go)" --no-regen

# PHONY targets are for ones that have conflicting files/dirs present:
.PHONY: test clean clean-lowered ir-stress-gate

# Build the browser-inspector example (delegates to its own Makefile).
.PHONY: browser-inspector
browser-inspector:
	$(MAKE) -C examples/browser-inspector build
