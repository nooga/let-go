build:
	GOOS=js GOARCH=wasm go build -o letgo.wasm .

run: build
	npm install
	node index.js