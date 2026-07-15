package resolver

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestParseSearchPaths(t *testing.T) {
	sep := string(os.PathListSeparator)
	got := ParseSearchPaths("a" + sep + "b" + sep + "" + sep + "c")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseSearchPaths() = %#v, want %#v", got, want)
	}
}

func TestPathsFromInputs_UsesFallbackWhenNotExplicit(t *testing.T) {
	sep := string(os.PathListSeparator)
	got := PathsFromInputs("ignored", "x"+sep+"y", false)
	want := []string{"x", "y"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PathsFromInputs() = %#v, want %#v", got, want)
	}
}

func TestPathsFromInputs_ExplicitOverridesFallback(t *testing.T) {
	sep := string(os.PathListSeparator)
	got := PathsFromInputs("a"+sep+"b", "x"+sep+"y", true)
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PathsFromInputs() = %#v, want %#v", got, want)
	}
}

func TestPathsFromInputs_ExplicitEmptyMeansNoPaths(t *testing.T) {
	got := PathsFromInputs("", "x", true)
	if len(got) != 0 {
		t.Fatalf("PathsFromInputs() = %#v, want empty (no paths)", got)
	}
}

func TestForceSourceNS(t *testing.T) {
	cases := []struct {
		env  string
		name string
		want bool
	}{
		{"", "ir.passes.typeinfer", false},                   // unset: never force
		{"ir.passes.typeinfer", "ir.passes.typeinfer", true}, // single match
		{"ir.passes.typeinfer", "ir.build", false},           // single non-match
		{"core, ir.build ,string", "ir.build", true},         // trimmed middle entry
		{"core,ir.build", "string", false},                   // not listed
		{"core,ir.build", "", false},                         // empty name never matches a listed ns
	}
	for _, c := range cases {
		t.Setenv("LG_FORCE_SOURCE_NS", c.env)
		if got := forceSourceNS(c.name); got != c.want {
			t.Errorf("forceSourceNS(%q) with LG_FORCE_SOURCE_NS=%q = %v, want %v", c.name, c.env, got, c.want)
		}
	}
}

// TestRequireUnregisteredTermReportsUnavailable guards the wasip1 regression
// from nooga/let-go#466: gating pkg/rt/term.go off wasip1 leaves no
// installTermNS, so `term` is never registered. loadEmbedded's term special
// case must report it unavailable via a non-registering lookup; the previous
// rt.NS("term") re-registered a placeholder and re-entered the loader,
// recursing until the wasm stack was exhausted. Simulated here by removing the
// natively-installed term ns and requiring it — expect a clean error, not a
// stack overflow.
func TestRequireUnregisteredTermReportsUnavailable(t *testing.T) {
	consts := vm.NewConsts()
	ctx := compiler.NewCompiler(consts, rt.NS("user"))
	rt.SetNSLoader(NewNSResolver(ctx, []string{"."}))
	ctx.SetSource("<test>")

	// Drop the natively-installed term ns to mimic a platform without an
	// installTermNS (e.g. wasip1), then restore it so other tests are unaffected.
	saved := rt.LookupNS("term")
	rt.RemoveNS("term")
	defer func() {
		if saved != nil {
			rt.RegisterNS(saved)
		}
	}()

	_, _, err := ctx.CompileMultiple(strings.NewReader("(require 'term)"))
	if err == nil {
		t.Fatal("requiring an unregistered term ns should error, got nil (regressed to the recursive loader path?)")
	}
}
