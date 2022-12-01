build:
	GOOS=js GOARCH=wasm go build -o letgo.wasm .

run: build
	npm install
	node index.js

pages: build 
	mkdir static
	cp index.html index.js letgo.wasm wasm_exec.js static