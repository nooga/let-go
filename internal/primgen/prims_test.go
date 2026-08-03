package primgen

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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
	output := emitFile([]primSpec{spec}, "rt", "github.com/nooga/let-go/pkg/rt", true)

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
		out, err := formatGo(emitFile(specs, "rt", "github.com/nooga/let-go/pkg/rt", true))
		if err != nil {
			t.Fatalf("format generated primitives: %v", err)
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
	output := emitFile([]primSpec{spec}, "rt", "github.com/nooga/let-go/pkg/rt", true)

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

	output := emitFile([]primSpec{spec1, spec2}, "rt", "github.com/nooga/let-go/pkg/rt", true)

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
	output := emitFile([]primSpec{spec}, "rt", "github.com/nooga/let-go/pkg/rt", true)

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

	output := emitFile([]primSpec{spec1, spec2}, "rt", "github.com/nooga/let-go/pkg/rt", true)

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

func TestEmitFileUsesTargetPackageName(t *testing.T) {
	specs := []primSpec{{
		GoPkg: "github.com/nooga/let-go/pkg/rt/corefns", Package: "corefns",
		GoIdent: "Seq", Ns: "clojure.core", LgName: "seq",
		Arity: 1, ParamSpecs: []string{"vm.Value"}, ResultSpec: "vm.Value", NeedsError: true,
	}}
	out := emitFile(specs, "corefns", "github.com/nooga/let-go/pkg/rt/corefns", false)
	if !strings.Contains(out, "package corefns") {
		t.Fatalf("expected `package corefns` header, got:\n%s", out)
	}
	if strings.Contains(out, "package rt\n") {
		t.Fatalf("must not emit hardcoded package rt")
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

func TestContributeModeEmitsMetadataOnly(t *testing.T) {
	specs := []primSpec{{
		GoPkg: "github.com/nooga/let-go/pkg/rt/corefns", Package: "corefns",
		GoIdent: "Seq", Ns: "clojure.core", LgName: "seq",
		Arity: 1, ParamSpecs: []string{"vm.Value"}, ResultSpec: "vm.Value", NeedsError: true,
	}}
	out := emitFile(specs, "corefns", "github.com/nooga/let-go/pkg/rt/corefns", false) // contribute
	if !strings.Contains(out, "RegisterNativeModule(") {
		t.Fatal("contribute mode must still register module metadata")
	}
	if strings.Contains(out, "defGeneratedPrimitive(") {
		t.Fatal("contribute mode must NOT emit var binding")
	}
	if strings.Contains(out, "_adapt_Seq") {
		t.Fatal("contribute mode must NOT emit adapters (binding-only)")
	}
}

func TestOwnModeEmitsBinding(t *testing.T) {
	specs := []primSpec{{
		GoPkg: "github.com/nooga/let-go/pkg/rt", Package: "rt",
		GoIdent: "CorePlus", Ns: "clojure.core", LgName: "+",
		Arity: -1, Variadic: true, ResultSpec: "vm.Value", NeedsError: true,
	}}
	out := emitFile(specs, "rt", "github.com/nooga/let-go/pkg/rt", true) // own
	if !strings.Contains(out, "defGeneratedPrimitive(") {
		t.Fatal("own mode must emit var binding")
	}
}

func TestInitFormRtUsesInstaller(t *testing.T) {
	specs := []primSpec{{GoPkg: "github.com/nooga/let-go/pkg/rt", Package: "rt",
		GoIdent: "CorePlus", Ns: "clojure.core", LgName: "+", Arity: -1, Variadic: true,
		ResultSpec: "vm.Value", NeedsError: true}}
	out := emitFile(specs, "rt", "github.com/nooga/let-go/pkg/rt", true)
	if !strings.Contains(out, "func init() {\n\tRegisterInstaller(RegisterGeneratedPrimitives)\n}") {
		t.Fatalf("rt target must self-register via the installer queue:\n%s", out)
	}
	if !strings.Contains(out, "vm.SetSuppressShadowWarn(true)") {
		t.Fatal("own-mode registrar must suppress the core-shadow warning")
	}
}

func TestInitFormExternalCallsDirectly(t *testing.T) {
	specs := []primSpec{{GoPkg: "github.com/nooga/let-go/pkg/rt/corefns", Package: "corefns",
		GoIdent: "Seq", Ns: "clojure.core", LgName: "seq", Arity: 1,
		ParamSpecs: []string{"vm.Value"}, ResultSpec: "vm.Value", NeedsError: true}}
	out := emitFile(specs, "corefns", "github.com/nooga/let-go/pkg/rt/corefns", false)
	if !strings.Contains(out, "func init() {\n\tRegisterGeneratedPrimitives()\n}") {
		t.Fatalf("non-rt target must self-register with a direct call:\n%s", out)
	}
}

func TestExternalContributeQualifiesRtAndImportsRtOnly(t *testing.T) {
	specs := []primSpec{{
		GoPkg: "github.com/nooga/let-go/pkg/rt/corefns", Package: "corefns",
		GoIdent: "Seq", Ns: "clojure.core", LgName: "seq",
		Arity: 1, ParamSpecs: []string{"vm.Value"}, ResultSpec: "vm.Value", NeedsError: true,
	}}
	out := emitFile(specs, "corefns", "github.com/nooga/let-go/pkg/rt/corefns", false)
	for _, want := range []string{
		`"github.com/nooga/let-go/pkg/rt"`, // rt import present
		"rt.RegisterNativeModule(&rt.NativeModule{",
		"map[string]rt.NativeDirectFn{",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("external contribute output missing %q:\n%s", want, out)
		}
	}
	// contribute mode emits no adapters → fmt/vm are unused → must NOT be imported
	for _, bad := range []string{`"fmt"`, `"github.com/nooga/let-go/pkg/vm"`} {
		if strings.Contains(out, bad) {
			t.Fatalf("external contribute output must not import unused %q:\n%s", bad, out)
		}
	}
}

func TestRtTargetStaysUnqualified(t *testing.T) {
	specs := []primSpec{{
		GoPkg: "github.com/nooga/let-go/pkg/rt", Package: "rt",
		GoIdent: "CorePlus", Ns: "clojure.core", LgName: "+",
		Arity: -1, Variadic: true, ResultSpec: "vm.Value", NeedsError: true,
	}}
	out := emitFile(specs, "rt", "github.com/nooga/let-go/pkg/rt", true)
	if !strings.Contains(out, "RegisterNativeModule(&NativeModule{") {
		t.Fatal("rt target must emit unqualified RegisterNativeModule")
	}
	// rt target must not import its own package (check import block, not GoPkg field)
	importEnd := strings.Index(out, "\nfunc")
	if importEnd > 0 {
		importSection := out[:importEnd]
		if strings.Contains(importSection, `"github.com/nooga/let-go/pkg/rt"`) {
			t.Fatal("rt target must not import its own package")
		}
	}
	if !strings.Contains(out, `"github.com/nooga/let-go/pkg/vm"`) {
		t.Fatal("rt own-mode uses vm (adapters + SetSuppressShadowWarn) and must import it")
	}
}

// TestGeneratePreservesSourcePackageName guards against deriving the emitted
// package clause from the -go-pkg import-path basename. Go allows a package
// clause to differ from the final path component (e.g. ".../primitives/v2"
// declaring `package primitives`); the basename would emit `package v2` and
// the directory would fail to compile ("found packages primitives and v2").
func TestGeneratePreservesSourcePackageName(t *testing.T) {
	dir := t.TempDir()
	// A real scanned spec (bare //lg:native + //lg:ns / //lg:name — the scanner
	// matches `line == "lg:native"` exactly), so this exercises the spec-based
	// package-name path (primSpec.Package), not the empty-stub readPackageName
	// fallback.
	src := `package primitives

import "github.com/nooga/let-go/pkg/vm"

//lg:native
//lg:ns clojure.core
//lg:name noop
func Noop(v vm.Value) (vm.Value, error) { return v, nil }
`
	if err := os.WriteFile(filepath.Join(dir, "prims.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "zz_primitives_generated.go")
	// Import path ends in /v2, but the source declares `package primitives`.
	if err := Generate(dir, out, "example.com/primitives/v2"); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "package primitives") {
		t.Fatalf("expected `package primitives` from the source clause, got:\n%s", got)
	}
	if strings.Contains(got, "package v2") {
		t.Fatalf("must not derive `package v2` from the import-path basename:\n%s", got)
	}
}
