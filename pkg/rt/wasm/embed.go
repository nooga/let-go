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
	"fmt"
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
	hostJS := mustReplaceOnce(lgHostJS, "__WASM_EXEC_JS__", string(execJSON))
	hostJS = mustReplaceOnce(hostJS, "__WASM_GZ_B64__", string(b64JSON))
	return mustReplaceOnce(htmlTemplate, "__LG_HOST_JS_BODY_PLACEHOLDER__", hostJS)
}

// mustReplaceOnce panics unless marker appears exactly once in s. The
// templates are embedded at build time, so a missing or duplicated
// marker is a developer error in lg-host.js or host.html — fail loud
// rather than silently shipping a half-substituted bundle.
func mustReplaceOnce(s, marker, replacement string) string {
	if n := strings.Count(s, marker); n != 1 {
		panic(fmt.Sprintf("wasm.AssembleHTML: marker %q expected exactly once, got %d", marker, n))
	}
	return strings.Replace(s, marker, replacement, 1)
}
