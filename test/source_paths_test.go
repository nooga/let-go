package test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestParseSearchPaths(t *testing.T) {
	sep := string(os.PathListSeparator)
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"single", "a", []string{"a"}},
		{"two", "a" + sep + "b", []string{"a", "b"}},
		{"consecutive seps", "a" + sep + sep + "b", []string{"a", "b"}},
		{"leading sep", sep + "a", []string{"a"}},
		{"trailing sep", "a" + sep, []string{"a"}},
		{"only seps", sep + sep, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolver.ParseSearchPaths(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPathsFromInputs(t *testing.T) {
	cases := []struct {
		name        string
		explicit    string
		fallback    string
		explicitSet bool
		want        []string
	}{
		{"both empty", "", "", false, []string{"."}},
		{"explicit only", "a:b", "", true, []string{".", "a", "b"}},
		{"fallback only", "", "x:y", false, []string{".", "x", "y"}},
		{"explicit wins over fallback", "a:b", "x:y", true, []string{".", "a", "b"}},
		// Empty explicit must still override the fallback when the flag was
		// explicitly set; otherwise "flag wins over env" silently leaks env.
		{"explicit empty wins over fallback", "", "x:y", true, []string{"."}},
	}
	// Adjust separators for the host platform without rebuilding cases.
	sep := string(os.PathListSeparator)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			explicit := strings.ReplaceAll(tc.explicit, ":", sep)
			fallback := strings.ReplaceAll(tc.fallback, ":", sep)
			got := resolver.PathsFromInputs(explicit, fallback, tc.explicitSet)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResolverWithExtraPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lgsp-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := "(ns sptest.foo)\n(def x 42)\n"
	nsDir := filepath.Join(tmpDir, "sptest")
	if err := os.MkdirAll(nsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nsDir, "foo.lg"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	consts := vm.NewConsts()
	coreNS := rt.NS(rt.NameCoreNS)
	if coreNS == nil {
		t.Fatal("core ns not registered")
	}
	ctx := compiler.NewCompiler(consts, coreNS)
	r := resolver.NewNSResolver(ctx, []string{".", tmpDir})

	ns := r.Load("sptest.foo")
	if ns == nil {
		t.Fatalf("Load returned nil; expected to find sptest.foo under %s", tmpDir)
	}
	if ns.Name() != "sptest.foo" {
		t.Errorf("ns.Name() = %q, want %q", ns.Name(), "sptest.foo")
	}
}
