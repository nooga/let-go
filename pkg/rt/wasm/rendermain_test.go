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
		for _, audio := range []bool{false, true} {
			got := RenderMain("store-x", hostEval, audio)
			if strings.Contains(got, "__LG_") {
				t.Fatalf("hostEval=%v audio=%v: unsubstituted marker remains:\n%s", hostEval, audio, got)
			}
			if !strings.Contains(got, `"store-x"`) {
				t.Fatalf("hostEval=%v audio=%v: storage id not substituted", hostEval, audio)
			}
		}
	}
}

// Default (no -w-host-eval) must not pull in the eval bridge or its import —
// syscall/js unused would fail to compile, and the bundle should stay as it was.
func TestRenderMainHostEvalOff(t *testing.T) {
	got := RenderMain("s", false, false)
	for _, leak := range []string{"syscall/js", "pkg/wasmhost", "_lgEval", "_lgRequest", "select {}"} {
		if strings.Contains(got, leak) {
			t.Fatalf("host-eval code leaked into default bundle: %q", leak)
		}
	}
}

// With -w-host-eval the generated main imports syscall/js, installs the _lgEval
// and _lgRequest hooks, signals readiness, and parks so the runtime stays callable.
func TestRenderMainHostEvalOn(t *testing.T) {
	got := RenderMain("s", true, false)
	for _, want := range []string{`"syscall/js"`, `"github.com/nooga/let-go/pkg/wasmhost"`, `js.Global().Set("_lgEval", hostEval)`, `js.Global().Set("_lgRequest", hostRequestFn)`, "_lgRuntimeReady", "select {}"} {
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
		for _, audio := range []bool{false, true} {
			if _, err := format.Source([]byte(RenderMain("store-x", hostEval, audio))); err != nil {
				t.Fatalf("hostEval=%v audio=%v: generated main is not valid Go: %v", hostEval, audio, err)
			}
		}
	}
}

// The template is filled by string replacement, never by a Printf, so a doubled
// %% reaches the generated main as written and its error messages print a
// literal %v instead of the error.
func TestRenderMainFormatVerbsAreSingle(t *testing.T) {
	for _, hostEval := range []bool{false, true} {
		for _, audio := range []bool{false, true} {
			if got := RenderMain("s", hostEval, audio); strings.Contains(got, "%%") {
				t.Fatalf("hostEval=%v audio=%v: generated main contains %%%%:\n%s", hostEval, audio, got)
			}
		}
	}
}

// The generated main must decode the program through rt.DecodeExecUnit. That
// is the decode that marks the embedded core namespaces a program references
// for loading; an inline resolver skips it, and a web build that requires zip
// gets every zip var nil (#954).
func TestRenderMainDecodesThroughRuntime(t *testing.T) {
	for _, hostEval := range []bool{false, true} {
		for _, audio := range []bool{false, true} {
			got := RenderMain("s", hostEval, audio)
			if !strings.Contains(got, "rt.DecodeExecUnit(lgbData)") {
				t.Fatalf("hostEval=%v audio=%v: generated main does not decode through rt.DecodeExecUnit", hostEval, audio)
			}
			if strings.Contains(got, "bytecode.DecodeToExecUnit") {
				t.Fatalf("hostEval=%v audio=%v: generated main decodes with its own resolver", hostEval, audio)
			}
		}
	}
}

// -w-audio blank-imports pkg/audio, whose init binds the browser backend; the
// default bundle must not, so it links none of the mixer.
func TestRenderMainAudio(t *testing.T) {
	const imp = `_ "github.com/nooga/let-go/pkg/audio"`
	if got := RenderMain("s", false, false); strings.Contains(got, "pkg/audio") {
		t.Fatalf("the default bundle imports pkg/audio")
	}
	for _, hostEval := range []bool{false, true} {
		if got := RenderMain("s", hostEval, true); !strings.Contains(got, imp) {
			t.Fatalf("hostEval=%v: -w-audio output missing %s", hostEval, imp)
		}
	}
}

// An entry-frame main loads the program through LoadProgramNamespaces, which
// drains the Go overrides of linked lowered packages, and invokes the entry
// itself. The default replay loop must be absent: it would run the program
// without draining them.
func TestRenderMainForEntryFrame(t *testing.T) {
	for _, hostEval := range []bool{false, true} {
		for _, audio := range []bool{false, true} {
			got := RenderMainForEntryFrame("s", hostEval, audio, "app.core", "-main")
			if _, err := format.Source([]byte(got)); err != nil {
				t.Fatalf("hostEval=%v audio=%v: generated main is not valid Go: %v", hostEval, audio, err)
			}
			if strings.Contains(got, "__LG_") {
				t.Fatalf("hostEval=%v audio=%v: unsubstituted marker remains:\n%s", hostEval, audio, got)
			}
			load := strings.Index(got, "rt.LoadProgramNamespaces(unit)")
			main := strings.Index(got, "rt.RunProgramMainChunk(unit)")
			entry := strings.Index(got, `rt.InvokeProgramEntry(vm.NewExecContext(), "app.core", "-main", nil)`)
			if load < 0 || main < load || entry < main {
				t.Fatalf("hostEval=%v audio=%v: want load, main chunk, then entry, in that order:\n%s", hostEval, audio, got)
			}
			if strings.Contains(got, "unit.NSOrder") {
				t.Fatalf("hostEval=%v audio=%v: entry-frame main still has the replay loop", hostEval, audio)
			}
			if hostEval && !strings.Contains(got, "_lgEval") {
				t.Fatalf("audio=%v: -w-host-eval body missing from the entry-frame main", audio)
			}
		}
	}
}

// The entry's namespace and name are spliced in as Go string literals.
func TestRenderMainForEntryFrameQuotesEntry(t *testing.T) {
	got := RenderMainForEntryFrame("s", false, false, `odd"ns`, "main")
	if !strings.Contains(got, `"odd\"ns", "main"`) {
		t.Fatalf("entry not quoted:\n%s", got)
	}
}
