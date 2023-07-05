GOROOT := $(shell go env GOROOT)
GOOS=js GOARCH=wasm go build -o letgo.wasm ./pkg/compiler/**/* .
GOOS=js GOARCH=wasm go build -o letgo.wasm .
cp "$(GOROOT)/misc/wasm/wasm_exec.js" .
cp ../meta/logo.svg .

run: build
	npm install
	node index.js

pages: build 
	mkdir static
	cp index.html index.js letgo.wasm wasm_exec.js logo.svg static

