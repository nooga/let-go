package resolver

import (
	"os"
	"reflect"
	"testing"
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
