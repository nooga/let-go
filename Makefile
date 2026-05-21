###s
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


# Start repl by default:
default:: run

run: $(LG)
	./$<

build: $(LG)

generate: $(GO)
	go run ./cmd/lgbgen

$(LG): $(GO) lg.go pkg/**/*
	which go
	go build -ldflags="-s -w" -o $@ .

test: pkg/**/* $(GO)
	go test -count=1 -v ./test

clojure-compat-report: $(GO)
	@$(REPORT-SCRIPT)

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
.PHONY: test
