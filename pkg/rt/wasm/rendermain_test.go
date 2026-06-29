package wasm

import (
	"go/format"
	"strings"
	"testing"
)

// RenderMain must always substitute every placeholder — a leaked __LG_* marker
// would be a syntax error in the generated main.go.
func TestRenderMainSubstitutesAllMarkers(t *testing.T) {
	for _, hostEval := range []bool{false, true} {
		got := RenderMain("store-x", hostEval)
		if strings.Contains(got, "__LG_") {
			t.Fatalf("hostEval=%v: unsubstituted marker remains:\n%s", hostEval, got)
		}
		if !strings.Contains(got, `"store-x"`) {
			t.Fatalf("hostEval=%v: storage id not substituted", hostEval)
		}
	}
}

// Default (no -w-host-eval) must not pull in the eval bridge or its import —
// syscall/js unused would fail to compile, and the bundle should stay as it was.
func TestRenderMainHostEvalOff(t *testing.T) {
	got := RenderMain("s", false)
	for _, leak := range []string{"syscall/js", "encoding/json", "_lgEval", "_lgRequest", "select {}"} {
		if strings.Contains(got, leak) {
			t.Fatalf("host-eval code leaked into default bundle: %q", leak)
		}
	}
}

// With -w-host-eval the generated main imports syscall/js, installs the _lgEval
// and _lgRequest hooks, signals readiness, and parks so the runtime stays callable.
func TestRenderMainHostEvalOn(t *testing.T) {
	got := RenderMain("s", true)
	for _, want := range []string{`"syscall/js"`, `"encoding/json"`, `js.Global().Set("_lgEval", hostEval)`, `js.Global().Set("_lgRequest", hostRequestFn)`, "_lgRuntimeReady", "select {}"} {
		if !strings.Contains(got, want) {
			t.Fatalf("host-eval output missing %q", want)
		}
	}
	if strings.Contains(got, "hostRequest := js.FuncOf") {
		t.Fatalf("host-eval output reuses hostRequest as both type and value")
	}
}

// The generated main is never run through gofmt by the build, so a bad indent
// or dangling token from substitution would ship as-is. format.Source parses
// it — a syntax error from the splice fails here. (Exact gofmt equality isn't
// asserted: wasmMainTmpl predates this change and isn't gofmt-sorted, e.g. the
// `_ "embed"` import; that's upstream's to normalize, not this flag's.)
func TestRenderMainIsValidGo(t *testing.T) {
	for _, hostEval := range []bool{false, true} {
		if _, err := format.Source([]byte(RenderMain("store-x", hostEval))); err != nil {
			t.Fatalf("hostEval=%v: generated main is not valid Go: %v", hostEval, err)
		}
	}
}
