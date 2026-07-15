package main

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestScanDirectives(t *testing.T) {
	src := `package builtins
import "github.com/nooga/let-go/pkg/vm"
//lg:native
//lg:ns clojure.string
//lg:name upper-case
func UpperCase(s string) (string, error) { return "", nil }`
	specs, scanErr := scanSource("x.go", []byte(src))
	if scanErr != nil {
		t.Fatalf("scanSource: %v", scanErr)
	}
	if len(specs) != 1 {
		t.Fatalf("want 1 spec, got %d", len(specs))
	}
	s := specs[0]
	if s.Ns != "clojure.string" || s.LgName != "upper-case" || s.GoIdent != "UpperCase" ||
		s.Arity != 1 || s.NeedsEC || !s.NeedsError ||
		s.ParamSpecs[0] != "string" || s.ResultSpec != "string" {
		t.Fatalf("bad spec: %+v", s)
	}
}

func TestEmitAdapterAndRegistrar(t *testing.T) {
	spec := primSpec{
		GoPkg:      "github.com/nooga/let-go/pkg/rt/builtins",
		GoIdent:    "UpperCase",
		Ns:         "clojure.string",
		LgName:     "upper-case",
		Arity:      1,
		ParamSpecs: []string{"string"},
		ResultSpec: "string",
		NeedsError: true,
		Package:    "builtins",
	}
	output := emitFile([]primSpec{spec})

	// Check for adapter function
	if !strings.Contains(output, "func _adapt_UpperCase") {
		t.Errorf("missing adapter func _adapt_UpperCase\noutput:\n%s", output)
	}

	// Check for arity check
	if !strings.Contains(output, `len(vs) != 1`) {
		t.Errorf("missing arity check in adapter\noutput:\n%s", output)
	}

	// Check for type coercion
	if !strings.Contains(output, `vm.String`) {
		t.Errorf("missing type coercion to vm.String\noutput:\n%s", output)
	}

	// Check for RegisterGeneratedPrimitives function
	if !strings.Contains(output, "func RegisterGeneratedPrimitives()") {
		t.Errorf("missing RegisterGeneratedPrimitives function\noutput:\n%s", output)
	}

	// Check for RegisterNativeModule call
	if !strings.Contains(output, "RegisterNativeModule") {
		t.Errorf("missing RegisterNativeModule call\noutput:\n%s", output)
	}

	// Check for the NativeModule struct with correct fields
	if !strings.Contains(output, `GoIdent: "UpperCase"`) {
		t.Errorf("missing GoIdent in NativeDirectFn\noutput:\n%s", output)
	}

	if !strings.Contains(output, `ParamSpecs: []string{"string"}`) {
		t.Errorf("missing ParamSpecs in NativeDirectFn\noutput:\n%s", output)
	}

	if !strings.Contains(output, `ResultSpec: "string"`) {
		t.Errorf("missing ResultSpec in NativeDirectFn\noutput:\n%s", output)
	}

	// Check for the import of the builtins package
	if !strings.Contains(output, `"github.com/nooga/let-go/pkg/rt/builtins"`) {
		t.Errorf("missing import for builtins package\noutput:\n%s", output)
	}

	// Check that the package name is used to qualify the function call
	if !strings.Contains(output, "builtins.UpperCase(") {
		t.Errorf("missing qualified builtins.UpperCase call\noutput:\n%s", output)
	}
}

// TestEmitFileDeterministic guards against nondeterministic generator output.
// emitFile registers modules by ranging over moduleMap, a Go map whose
// iteration order is randomized. gofmt (applied before the file is written)
// canonicalizes the import block but does NOT reorder function-body statements,
// so without an explicit key sort the RegisterGeneratedPrimitives() body varies
// run-to-run, intermittently dirtying the checked-in artifact and failing
// freshness checks. This asserts the actual committed shape — the gofmt'd
// output — is byte-stable across many regenerations. Several distinct packages
// AND namespaces exercise both the import and registration paths.
func TestEmitFileDeterministic(t *testing.T) {
	specs := []primSpec{
		{GoPkg: "github.com/nooga/let-go/pkg/rt/alpha", GoIdent: "AlphaOne", Ns: "clojure.core", LgName: "alpha-one", Arity: 1, ParamSpecs: []string{"string"}, ResultSpec: "string", Package: "alpha"},
		{GoPkg: "github.com/nooga/let-go/pkg/rt/beta", GoIdent: "BetaTwo", Ns: "clojure.string", LgName: "beta-two", Arity: 1, ParamSpecs: []string{"int"}, ResultSpec: "int", Package: "beta"},
		{GoPkg: "github.com/nooga/let-go/pkg/rt/gamma", GoIdent: "GammaThree", Ns: "clojure.set", LgName: "gamma-three", Arity: 1, ParamSpecs: []string{"string"}, ResultSpec: "string", Package: "gamma"},
		{GoPkg: "github.com/nooga/let-go/pkg/rt/alpha", GoIdent: "AlphaFour", Ns: "clojure.core", LgName: "alpha-four", Arity: 1, ParamSpecs: []string{"int"}, ResultSpec: "int", Package: "alpha"},
	}

	canonical := func() string {
		out, err := gofmtCode(emitFile(specs))
		if err != nil {
			t.Fatalf("gofmt generated primitives: %v", err)
		}
		return out
	}

	want := canonical()
	for i := 0; i < 30; i++ {
		if got := canonical(); got != want {
			t.Fatalf("generated primitives are nondeterministic (run %d differs from run 0):\n--- run 0 ---\n%s\n--- run %d ---\n%s", i, want, i, got)
		}
	}
}

func TestScanDirectivesWithEC(t *testing.T) {
	src := `package builtins
import "github.com/nooga/let-go/pkg/vm"
//lg:native
//lg:ns clojure.core
//lg:name seq
func Seq(ec *vm.ExecContext, v vm.Value) (vm.Value, error) { return vm.NIL, nil }`
	specs, scanErr := scanSource("x.go", []byte(src))
	if scanErr != nil {
		t.Fatalf("scanSource: %v", scanErr)
	}
	if len(specs) != 1 {
		t.Fatalf("want 1 spec, got %d", len(specs))
	}
	s := specs[0]
	if s.GoIdent != "Seq" || s.Arity != 1 || !s.NeedsEC {
		t.Fatalf("bad spec: %+v", s)
	}
	if len(s.ParamSpecs) != 1 || s.ParamSpecs[0] != "vm.Value" {
		t.Fatalf("bad param specs: %v", s.ParamSpecs)
	}
}

func TestECEmission(t *testing.T) {
	// Test that EC-needing functions emit ec-aware adapters with the proper signature
	spec := primSpec{
		GoPkg:      "github.com/nooga/let-go/pkg/rt/builtins",
		GoIdent:    "Seq",
		Ns:         "clojure.core",
		LgName:     "seq",
		Arity:      1,
		ParamSpecs: []string{"vm.Value"},
		ResultSpec: "vm.Value",
		NeedsError: false,
		NeedsEC:    true,
		Package:    "builtins",
	}
	output := emitFile([]primSpec{spec})

	// Check that the ec-aware adapter has the correct signature
	// func _adapt_Seq(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error)
	if !strings.Contains(output, "func _adapt_Seq(ec *vm.ExecContext, vs []vm.Value)") {
		t.Errorf("missing ec-aware adapter signature\noutput:\n%s", output)
	}

	// Check that the adapter passes ec to the Go function call
	if !strings.Contains(output, "builtins.Seq(ec") {
		t.Errorf("adapter does not pass ec to builtins.Seq\noutput:\n%s", output)
	}

	// Check that vm.NewCtxNativeFn is used in registration
	if !strings.Contains(output, "vm.NewCtxNativeFn") {
		t.Errorf("missing vm.NewCtxNativeFn in registration\noutput:\n%s", output)
	}

	// Check for the import of the builtins package
	if !strings.Contains(output, `"github.com/nooga/let-go/pkg/rt/builtins"`) {
		t.Errorf("missing import for builtins package\noutput:\n%s", output)
	}

	// Parse-check the output
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "gen.go", output, 0); err != nil {
		t.Errorf("generated code does not parse: %v\noutput:\n%s", err, output)
	}
}

func TestMultiArity(t *testing.T) {
	// Test that multi-arity functions produce one dispatch adapter and multiple registry entries
	spec1 := primSpec{
		GoPkg:      "github.com/nooga/let-go/pkg/rt/builtins",
		GoIdent:    "Foo",
		Ns:         "clojure.core",
		LgName:     "foo",
		Arity:      2,
		ParamSpecs: []string{"string", "string"},
		ResultSpec: "string",
		NeedsError: false,
		NeedsEC:    false,
		Package:    "builtins",
	}
	spec2 := primSpec{
		GoPkg:      "github.com/nooga/let-go/pkg/rt/builtins",
		GoIdent:    "Foo",
		Ns:         "clojure.core",
		LgName:     "foo",
		Arity:      3,
		ParamSpecs: []string{"string", "string", "string"},
		ResultSpec: "string",
		NeedsError: false,
		NeedsEC:    false,
		Package:    "builtins",
	}

	output := emitFile([]primSpec{spec1, spec2})

	// Check for dispatch adapter with len(vs) switch
	if !strings.Contains(output, "switch len(vs)") {
		t.Errorf("missing dispatch adapter with len(vs) switch\noutput:\n%s", output)
	}

	// Check for arity-specific adapters
	if !strings.Contains(output, "_adapt_Foo_arity2") {
		t.Errorf("missing arity-specific adapter _adapt_Foo_arity2\noutput:\n%s", output)
	}
	if !strings.Contains(output, "_adapt_Foo_arity3") {
		t.Errorf("missing arity-specific adapter _adapt_Foo_arity3\noutput:\n%s", output)
	}

	// Check for registry entries with arity-qualified keys
	// In the Fns map, we should have entries for both arities with keys like "foo@2" and "foo@3"
	if !strings.Contains(output, `"foo@2": {GoIdent: "Foo", LgName: "foo", Arity: 2`) {
		t.Errorf("missing registry entry for foo@2\noutput:\n%s", output)
	}
	if !strings.Contains(output, `"foo@3": {GoIdent: "Foo", LgName: "foo", Arity: 3`) {
		t.Errorf("missing registry entry for foo@3\noutput:\n%s", output)
	}

	// Parse-check the output
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "gen.go", output, 0); err != nil {
		t.Errorf("generated code does not parse: %v\noutput:\n%s", err, output)
	}
}

func TestGeneratedCodeParses(t *testing.T) {
	// Test that all generated code passes go/parser.ParseFile check
	spec := primSpec{
		GoPkg:      "github.com/nooga/let-go/pkg/rt/builtins",
		GoIdent:    "UpperCase",
		Ns:         "clojure.string",
		LgName:     "upper-case",
		Arity:      1,
		ParamSpecs: []string{"string"},
		ResultSpec: "string",
		NeedsError: true,
		Package:    "builtins",
	}
	output := emitFile([]primSpec{spec})

	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "gen.go", output, 0); err != nil {
		t.Errorf("generated code does not parse: %v\noutput:\n%s", err, output)
	}
}

func TestMultiArityEC(t *testing.T) {
	// Test that multi-arity with NeedsEC generates ec-aware dispatch and arity-specific adapters
	spec1 := primSpec{
		GoPkg:      "github.com/nooga/let-go/pkg/rt/builtins",
		GoIdent:    "Foo",
		Ns:         "clojure.core",
		LgName:     "foo",
		Arity:      2,
		ParamSpecs: []string{"vm.Value", "vm.Value"},
		ResultSpec: "vm.Value",
		NeedsError: false,
		NeedsEC:    true, // This spec needs EC
		Package:    "builtins",
	}
	spec2 := primSpec{
		GoPkg:      "github.com/nooga/let-go/pkg/rt/builtins",
		GoIdent:    "Foo",
		Ns:         "clojure.core",
		LgName:     "foo",
		Arity:      3,
		ParamSpecs: []string{"vm.Value", "vm.Value", "vm.Value"},
		ResultSpec: "vm.Value",
		NeedsError: false,
		NeedsEC:    false, // This spec doesn't need EC
		Package:    "builtins",
	}

	output := emitFile([]primSpec{spec1, spec2})

	// Check that dispatch adapter is ec-aware (signature includes ec *vm.ExecContext)
	if !strings.Contains(output, "func _adapt_Foo(ec *vm.ExecContext, vs []vm.Value)") {
		t.Errorf("dispatch adapter should be ec-aware\noutput:\n%s", output)
	}

	// Check that dispatch adapter passes ec to arity-specific adapters
	if !strings.Contains(output, "_adapt_Foo_arity2(ec, vs)") {
		t.Errorf("dispatch should pass ec to _adapt_Foo_arity2\noutput:\n%s", output)
	}
	if !strings.Contains(output, "_adapt_Foo_arity3(ec, vs)") {
		t.Errorf("dispatch should pass ec to _adapt_Foo_arity3\noutput:\n%s", output)
	}

	// Check that both arity-specific adapters have ec-aware signature
	if !strings.Contains(output, "func _adapt_Foo_arity2(ec *vm.ExecContext, vs []vm.Value)") {
		t.Errorf("_adapt_Foo_arity2 should have ec-aware signature\noutput:\n%s", output)
	}
	if !strings.Contains(output, "func _adapt_Foo_arity3(ec *vm.ExecContext, vs []vm.Value)") {
		t.Errorf("_adapt_Foo_arity3 should have ec-aware signature\noutput:\n%s", output)
	}

	// Check that the NeedsEC arity-specific adapter passes ec to Go function
	if !strings.Contains(output, "builtins.Foo(ec") {
		t.Errorf("ec-needing arity should pass ec to builtins.Foo\noutput:\n%s", output)
	}

	// Check that the non-NeedsEC arity-specific adapter has "_ = ec" to suppress unused parameter
	if !strings.Contains(output, "_ = ec") {
		t.Errorf("non-ec arity should suppress unused ec parameter\noutput:\n%s", output)
	}

	// Check that vm.NewCtxNativeFn is used in registration
	if !strings.Contains(output, "vm.NewCtxNativeFn") {
		t.Errorf("missing vm.NewCtxNativeFn in registration\noutput:\n%s", output)
	}

	// Parse-check the output
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "gen.go", output, 0); err != nil {
		t.Errorf("generated code does not parse: %v\noutput:\n%s", err, output)
	}
}

func TestScanSourceParseErrorPropagates(t *testing.T) {
	if _, err := scanSource("x.go", []byte("package m\nfunc {")); err == nil {
		t.Fatal("expected parse error, got nil — a parse failure must not read as 'no primitives'")
	}
}

func TestHasBuildConstraint(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want bool
	}{
		{"go:build line", "//go:build plan9\n\npackage m\n", true},
		{"legacy plus-build", "// +build plan9\n\npackage m\n", true},
		{"unconstrained", "// just a comment\npackage m\n", false},
		{"constraint-looking line after package clause", "package m\n\n// comment mentioning //go:build syntax\nvar x = 1\n", false},
	}
	for _, tc := range cases {
		if got := hasBuildConstraint([]byte(tc.src)); got != tc.want {
			t.Errorf("%s: hasBuildConstraint = %v, want %v", tc.name, got, tc.want)
		}
	}
}
