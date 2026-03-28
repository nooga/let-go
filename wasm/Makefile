GOROOT := $(shell go env GOROOT)

build:
	GOOS=js GOARCH=wasm go build -o letgo.wasm .
	cp "$(GOROOT)/lib/wasm/wasm_exec.js" . 2>/dev/null || cp "$(GOROOT)/misc/wasm/wasm_exec.js" .
	cp ../meta/logo.svg .

pages: build
	mkdir -p static
	cp index.html letgo.wasm wasm_exec.js logo.svg static

clean:
	rm -f letgo.wasm wasm_exec.js logo.svg
	rm -rf static

.PHONY: build pages clean
