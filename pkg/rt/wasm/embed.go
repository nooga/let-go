// Package wasm holds build-time JS and HTML assets for the `lg -w` WASM
// bundler. AssembleHTML returns a single self-contained HTML page given
// the Go runtime support source and the gzipped-base64 program WASM.
//
// lg-host.js carries two markers (__WASM_EXEC_JS__ and __WASM_GZ_B64__)
// the assembler substitutes with JSON-encoded JS strings. host.html
// carries a single marker (__LG_HOST_JS_BODY_PLACEHOLDER__) where the
// populated host JS is inlined.
package wasm

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed lg-host.js
var lgHostJS string

//go:embed host.html
var htmlTemplate string

// AssembleHTML returns the full self-contained HTML page produced by
// `lg -w`. Pure function: same inputs produce same output. Tested via
// a golden file in testdata/.
func AssembleHTML(wasmExecJS, wasmGzB64 string) string {
	execJSON, _ := json.Marshal(wasmExecJS)
	b64JSON, _ := json.Marshal(wasmGzB64)
	hostJS := strings.Replace(lgHostJS, "__WASM_EXEC_JS__", string(execJSON), 1)
	hostJS = strings.Replace(hostJS, "__WASM_GZ_B64__", string(b64JSON), 1)
	return strings.Replace(htmlTemplate, "__LG_HOST_JS_BODY_PLACEHOLDER__", hostJS, 1)
}
