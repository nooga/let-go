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
