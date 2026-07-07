// Package wasm holds build-time JS and HTML assets for the `lg -w` WASM
// bundler. AssembleHTML returns a single self-contained HTML page given
// the Go runtime support source and the gzipped-base64 program WASM.
//
// The host JS is split into two assets:
//   - lg-host-core.js: the host-agnostic glue (COI, wasm decode, the
//     window.LetGoHost surface, the worker/main-thread boot). Always emitted.
//   - lg-shell-xterm.js: the default xterm.js shell, which binds to the
//     runtime only through LetGoHost. Emitted unless shell == false, i.e.
//     `lg -w -w-shell none`, where the client supplies its own shell.
//
// lg-host-core.js carries two markers (__WASM_EXEC_JS__ and __WASM_GZ_B64__)
// the assembler substitutes with JSON-encoded JS strings. host.html carries
// __LG_XTERM_CSS__ / __LG_XTERM_JS__ (the xterm CDN tags, dropped in shell-less
// builds) and __LG_HOST_JS_BODY_PLACEHOLDER__ where the host JS is inlined.
package wasm

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed lg-host-core.js
var lgHostCoreJS string

//go:embed lg-shell-xterm.js
var lgShellXtermJS string

//go:embed host.html
var htmlTemplate string

const (
	xtermCSS = `  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/css/xterm.min.css">
  <style>.xterm{height:100%}</style>`
	xtermJS = `  <script src="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/lib/xterm.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/@xterm/addon-fit@0.10.0/lib/addon-fit.min.js"></script>`
)

// HostBodyMarker is where the assembled host JS is spliced into an HTML
// template (the embedded host.html, or a caller's custom -w-shell template).
const HostBodyMarker = "__LG_HOST_JS_BODY_PLACEHOLDER__"

// buildHostJS assembles the host-agnostic runtime JS (wasm_exec, the inline or
// external wasm payload, the LetGoHost surface, and the optional eval hook).
// Shared by AssembleHTML (built-in shell) and AssembleHTMLWithTemplate (custom).
func buildHostJS(wasmExecJS, wasmGzB64 string, externalWasm, hostEval bool) string {
	execJSON, _ := json.Marshal(wasmExecJS)

	mode := "inline"
	if externalWasm {
		mode = "external"
		wasmGzB64 = ""
	}
	modeJSON, _ := json.Marshal(mode)
	b64JSON, _ := json.Marshal(wasmGzB64)
	hostEvalJSON, _ := json.Marshal(hostEval)

	hostJS := mustReplaceOnce(lgHostCoreJS, "__WASM_EXEC_JS__", string(execJSON))
	hostJS = mustReplaceOnce(hostJS, "__WASM_MODE__", string(modeJSON))
	hostJS = mustReplaceOnce(hostJS, "__WASM_GZ_B64__", string(b64JSON))
	hostJS = mustReplaceOnce(hostJS, "__LG_HOST_EVAL__", string(hostEvalJSON))
	return hostJS
}

// AssembleHTML returns the full self-contained HTML page produced by
// `lg -w`. With shell == true the default xterm shell and its CDN tags are
// included; with shell == false only the host-agnostic core ships and the
// client binds its own shell to window.LetGoHost. With externalWasm == true
// the payload is delivered as a separate main.wasm the loader fetches and
// streams (wasmGzB64 is ignored); otherwise it is gzip-base64 inlined into the
// page. With hostEval == true the page exposes LetGoHost.eval (the -w-host-eval
// bundle); when false it is omitted, so feature detection stays honest. Pure
// function: same inputs produce same output. Tested via golden files in testdata/.
func AssembleHTML(wasmExecJS, wasmGzB64 string, shell, externalWasm, hostEval bool) string {
	hostJS := buildHostJS(wasmExecJS, wasmGzB64, externalWasm, hostEval)

	css, js := "", ""
	if shell {
		hostJS = hostJS + "\n" + lgShellXtermJS
		css, js = xtermCSS, xtermJS
	}

	out := mustReplaceOnce(htmlTemplate, "__LG_XTERM_CSS__", css)
	out = mustReplaceOnce(out, "__LG_XTERM_JS__", js)
	return mustReplaceOnce(out, HostBodyMarker, hostJS)
}

// AssembleHTMLWithTemplate injects the assembled host JS into a caller-supplied
// HTML template at HostBodyMarker. The custom template owns its own shell/markup,
// so the built-in xterm assets are never added and no __LG_XTERM_* markers are
// touched. externalWasm/hostEval behave as in AssembleHTML.
func AssembleHTMLWithTemplate(template, wasmExecJS, wasmGzB64 string, externalWasm, hostEval bool) string {
	hostJS := buildHostJS(wasmExecJS, wasmGzB64, externalWasm, hostEval)
	return mustReplaceOnce(template, HostBodyMarker, hostJS)
}

// mustReplaceOnce panics unless marker appears exactly once in s. The
// templates are embedded at build time, so a missing or duplicated
// marker is a developer error in the JS or HTML assets — fail loud
// rather than silently shipping a half-substituted bundle.
func mustReplaceOnce(s, marker, replacement string) string {
	if n := strings.Count(s, marker); n != 1 {
		panic(fmt.Sprintf("wasm.AssembleHTML: marker %q expected exactly once, got %d", marker, n))
	}
	return strings.Replace(s, marker, replacement, 1)
}
