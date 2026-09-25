package resolver

import (
	"io"
	"os"
	"path/filepath"
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

func TestLoadReturnsCompileErrorWithoutPrinting(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "broken.lg")
	source := "(ns broken)\n(def broken-value\n  (fn []\n    (let [:tag 1] 1)))\n"
	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}

	ctx := compiler.NewCompiler(vm.NewConsts(), rt.CoreNS)
	resolver := NewNSResolver(ctx, []string{dir})
	readEnd, writeEnd, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatal(pipeErr)
	}
	originalStderr := os.Stderr
	defer func() {
		os.Stderr = originalStderr
		_ = writeEnd.Close()
		_ = readEnd.Close()
	}()
	os.Stderr = writeEnd

	_, err := resolver.LoadWithError("broken")
	if closeErr := writeEnd.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	os.Stderr = originalStderr
	printed, readErr := io.ReadAll(readEnd)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if err == nil || !strings.Contains(err.Error(), "let binding name must be a symbol") {
		t.Fatalf("Load error = %v, want original compile chain", err)
	}
	if len(printed) != 0 {
		t.Fatalf("resolver printed an error it should propagate: %s", printed)
	}
}

func TestRequirePropagatesTaggedReaderRegistry(t *testing.T) {
	const nsName = "pr770-registry-dep"
	dir := t.TempDir()
	file := filepath.Join(dir, "pr770_registry_dep.lg")
	if err := os.WriteFile(file, []byte("(ns "+nsName+")\n(def value #review/probe 1)\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	registry := compiler.NewTaggedReaderRegistry()
	if err := registry.RegisterData("review/probe", func(vm.Value) (vm.Value, error) {
		return vm.Int(42), nil
	}); err != nil {
		t.Fatal(err)
	}

	parentNS := rt.NS("pr770-registry-parent")
	ctx := compiler.NewCompiler(vm.NewConsts(), parentNS).SetTaggedReaders(registry)
	ctx.SetSource("<require-registry-test>")
	loader := NewNSResolver(ctx, []string{dir})
	previousLoader := rt.GetNSLoader()
	rt.SetNSLoader(loader)
	t.Cleanup(func() {
		rt.SetNSLoader(previousLoader)
		rt.RemoveNS(nsName)
		rt.RemoveNS("pr770-registry-parent")
	})

	_, got, err := ctx.CompileMultiple(strings.NewReader("(do (require '" + nsName + ") " + nsName + "/value)"))
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.Int(42) {
		t.Fatalf("required namespace value = %v, want 42 from the parent tagged-reader registry", got)
	}
}

func TestExecPrecompiledPropagatesTaggedReaderRegistry(t *testing.T) {
	const (
		nsName     = "pr770-precompiled-registry-dep"
		parentName = "pr770-precompiled-registry-parent"
	)
	registry := compiler.NewTaggedReaderRegistry()
	if err := registry.RegisterData("review/probe", func(vm.Value) (vm.Value, error) {
		return vm.Int(42), nil
	}); err != nil {
		t.Fatal(err)
	}

	parentNS := rt.NS(parentName)
	ctx := compiler.NewCompiler(vm.NewConsts(), parentNS).SetTaggedReaders(registry)
	ctx.SetSource("<precompiled-registry-test>")
	chunk, _, err := ctx.CompileMultiple(strings.NewReader(
		"(ns " + nsName + ")\n(def value (read-string \"#review/probe 1\"))\n"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		rt.RemoveNS(nsName)
		rt.RemoveNS(parentName)
	})

	ctx.SetCurrentNS(parentNS)
	loaded, err := NewNSResolver(ctx, nil).execPrecompiled(nsName, chunk)
	if err != nil {
		t.Fatal(err)
	}
	value := loaded.LookupLocal(vm.Symbol("value"))
	if value == nil {
		t.Fatal("precompiled namespace has no value var")
	}
	if got := value.Deref(); got != vm.Int(42) {
		t.Fatalf("precompiled namespace value = %v, want 42 from the parent tagged-reader registry", got)
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

// loadScopeFixture writes files (relative path -> source) under a temp dir and
// returns a compiler whose require resolves against it. It restores the
// namespace loader, removes the fixture namespaces, and restores the root of
// each named core var afterwards, so a leaked set! cannot bleed into other tests.
func loadScopeFixture(t *testing.T, files map[string]string, nsNames []string, vars ...string) *compiler.Context {
	t.Helper()
	dir := t.TempDir()
	for rel, source := range files {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range vars {
		v := rt.CoreNS.LookupLocal(vm.Symbol(name))
		if v == nil {
			t.Fatalf("%s var not defined", name)
		}
		old := v.Root()
		t.Cleanup(func() { v.SetRoot(old) })
	}
	const parentName = "load-scope-parent"
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(parentName))
	ctx.SetSource("<load-scope-test>")
	previousLoader := rt.GetNSLoader()
	rt.SetNSLoader(NewNSResolver(ctx, []string{dir}))
	t.Cleanup(func() {
		rt.SetNSLoader(previousLoader)
		for _, name := range nsNames {
			rt.RemoveNS(name)
		}
		rt.RemoveNS(parentName)
	})
	return ctx
}

func captureReflectionWarnings(t *testing.T) *strings.Builder {
	t.Helper()
	var output strings.Builder
	restore := rt.SetReflectionWarningWriter(&output)
	t.Cleanup(restore)
	rt.ResetReflectionWarnings()
	return &output
}

func coreVarTruthy(t *testing.T, name string) bool {
	t.Helper()
	return vm.IsTruthy(rt.CoreNS.LookupLocal(vm.Symbol(name)).Deref())
}

func TestRequireScopesWarnOnReflectionToTheLoadedFile(t *testing.T) {
	const flag = "*warn-on-reflection*"

	t.Run("set! does not leak into a later file", func(t *testing.T) {
		ctx := loadScopeFixture(t, map[string]string{
			"wscope/setter.lg": "(ns wscope.setter)\n(set! *warn-on-reflection* true)\n",
			"wscope/later.lg":  "(ns wscope.later)\n(defn f [s] (.length s))\n",
		}, []string{"wscope.setter", "wscope.later"}, flag)
		output := captureReflectionWarnings(t)
		if _, _, err := ctx.CompileMultiple(strings.NewReader("(require 'wscope.setter)\n(require 'wscope.later)\n")); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(output.String(), "reflection warning") {
			t.Fatalf("a file that never set the flag warned:\n%s", output.String())
		}
		if coreVarTruthy(t, flag) {
			t.Fatal("*warn-on-reflection* is still true after the setting file finished loading")
		}
	})

	t.Run("set! still warns inside the file that sets it", func(t *testing.T) {
		ctx := loadScopeFixture(t, map[string]string{
			"wscope/self.lg": "(ns wscope.self)\n(set! *warn-on-reflection* true)\n(defn g [s] (.length s))\n",
		}, []string{"wscope.self"}, flag)
		output := captureReflectionWarnings(t)
		if _, _, err := ctx.CompileMultiple(strings.NewReader("(require 'wscope.self)\n")); err != nil {
			t.Fatal(err)
		}
		got := output.String()
		if count := strings.Count(got, "reflection warning"); count != 1 {
			t.Fatalf("warning count = %d, want 1:\n%s", count, got)
		}
		if !strings.Contains(got, "self.lg") {
			t.Fatalf("warning does not name the setting file:\n%s", got)
		}
	})

	t.Run("the per-load binding inherits an enclosing true", func(t *testing.T) {
		ctx := loadScopeFixture(t, map[string]string{
			"wscope/plain.lg": "(ns wscope.plain)\n(defn h [s] (.length s))\n",
		}, []string{"wscope.plain"}, flag)
		warnVar := rt.CoreNS.LookupLocal(vm.Symbol(flag))
		warnVar.SetRoot(vm.FALSE)
		vm.RootExecContext.PushBinding(warnVar, vm.TRUE)
		defer vm.RootExecContext.PopBinding(warnVar)
		output := captureReflectionWarnings(t)
		if _, _, err := ctx.CompileMultiple(strings.NewReader("(require 'wscope.plain)\n")); err != nil {
			t.Fatal(err)
		}
		if count := strings.Count(output.String(), "reflection warning"); count != 1 {
			t.Fatalf("warning count = %d, want 1 from the inherited binding:\n%s", count, output.String())
		}
	})

	t.Run("the binding is popped when the load fails", func(t *testing.T) {
		ctx := loadScopeFixture(t, map[string]string{
			"wscope/broken.lg": "(ns wscope.broken)\n(set! *warn-on-reflection* true)\n(defn oops [] (let [:tag 1] 1))\n",
		}, []string{"wscope.broken"}, flag)
		captureReflectionWarnings(t)
		if _, _, err := ctx.CompileMultiple(strings.NewReader("(require 'wscope.broken)\n")); err == nil {
			t.Fatal("requiring a file with a compile error succeeded")
		}
		if coreVarTruthy(t, flag) {
			t.Fatal("*warn-on-reflection* is still true after the failed load")
		}
	})
}

func TestRequireScopesUncheckedMathToTheLoadedFile(t *testing.T) {
	const flag = "*unchecked-math*"
	ctx := loadScopeFixture(t, map[string]string{
		"mscope/setter.lg": "(ns mscope.setter)\n(set! *unchecked-math* true)\n",
		"mscope/later.lg":  "(ns mscope.later)\n(defn f [a b] (+ a b))\n",
	}, []string{"mscope.setter", "mscope.later"}, flag)
	if _, _, err := ctx.CompileMultiple(strings.NewReader("(require 'mscope.setter)\n(require 'mscope.later)\n")); err != nil {
		t.Fatal(err)
	}
	if coreVarTruthy(t, flag) {
		t.Fatal("*unchecked-math* is still true after the setting file finished loading")
	}
}
