# Build the browser-inspector viewer: lg compiles main.lg to a wasm web app and
# splices the host runtime into shell.html natively (no python) via
# `-w-shell shell.html`. Honors LG (prebuilt lg binary) and LG_WASM_BUILD_TAGS.
HERE := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
ROOT := $(abspath $(HERE)/../..)
LG ?=
LG_WASM_BUILD_TAGS ?=
# Per-build nonce stamped into the page so each build is uniquely identifiable
# (the whole app is one inlined index.html; a stale cache serves old code
# silently). Overridable; defaults to a UTC build timestamp.
BUILD_NONCE ?= $(shell date -u +%Y%m%dT%H%M%SZ)

.PHONY: build serve clean
build:
	mkdir -p "$(HERE)dist"
ifeq ($(LG),)
	cd "$(ROOT)" && env LG_WASM_BUILD_TAGS="$(LG_WASM_BUILD_TAGS)" \
	  go run ./ -w "$(HERE)dist" -w-shell "$(HERE)shell.html" -w-host-eval "$(HERE)main.lg"
else
	env LG_WASM_BUILD_TAGS="$(LG_WASM_BUILD_TAGS)" \
	  "$(LG)" -w "$(HERE)dist" -w-shell "$(HERE)shell.html" -w-host-eval "$(HERE)main.lg"
endif
	@python3 "$(HERE)stamp_nonce.py" "$(HERE)dist/index.html" "$(BUILD_NONCE)"
	@echo "wrote $(HERE)dist/index.html (build $(BUILD_NONCE))"

# Serve with caching disabled (serve.py sends Cache-Control: no-store) so the
# browser always refetches the freshly built page — no stale-wasm confusion.
serve: build
	cd "$(HERE)dist" && python3 "$(HERE)serve.py"

clean:
	rm -rf "$(HERE)dist"
