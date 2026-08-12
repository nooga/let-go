/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/internal/gofragment"
)

// TestNativeEntryASTGate is the iteration-1 native-entry gate.
//
// For every fixture under test/native-entry it makes ONE claim, end to end:
// the fixture's own source was lowered to Go by the production path
// (scripts/lg-compile --entry-frame), that generated Go has the expected
// structure, it was compiled into a real binary, and that binary produced the
// fixture's committed stdout.
//
// Four kinds of evidence are required per fixture:
//
//  1. native function exists — the fixture's defn appears in the generated
//     package as a Go function declaration matching a committed AST shape,
//     exactly once (gofragment, Kind=GoFunction);
//  2. the call is classified correctly — for a `direct` fixture the entry
//     function contains the direct Go call exactly once AND no
//     ec.Invoke(...) trampoline naming that defn; for a `vm` fixture the
//     inverse holds (trampoline exactly once, no direct call). The VM-backed
//     fixture is what scopes the direct-call claim: without it, "no
//     trampoline" could be an accident of the corpus;
//  3. it runs — the generated tree is compiled with `go build` and executed;
//     stdout must equal the committed <fixture>.expect byte for byte, exit 0;
//  4. it is falsifiable — the oracle mutants and the AST call-site mutant must
//     each be killed, for every fixture. Every mutation site is located
//     through go/ast, never by hardcoded temporaries or string offsets.
//
// The semantic mutant of the lowered Go body carries an INVERTED contract, and
// this is deliberate: the same mutation is applied to both kinds of fixture and
// opposite outcomes are required.
//
//   - lowering=direct — the mutant MUST change observed output
//     (live_native_body_mutant_dies). That is the evidence that the generated
//     Go body is what the entry path actually runs.
//   - lowering=vm — the SAME mutation MUST NOT change anything
//     (dead_native_body_mutant_survives): the binary must still build, exit 0
//     and print byte-exact .expect. That surviving mutant is required positive
//     evidence that the generated body is dead code and the ec.Invoke
//     trampoline is genuinely what runs. It is not a weakness in the gate; a
//     vm fixture whose body mutant DIES is a failure, because the body would
//     then be live and the fixture would no longer scope the direct-call claim.
//
// A fixture missing its .expect or .goexpect.json sidecar is a hard failure,
// never a skip; a toolchain that cannot build is a hard failure too.
func TestNativeEntryASTGate(t *testing.T) {
	if testing.Short() {
		t.Skip("native-entry gate lowers, builds and runs real binaries; run `make native-entry-gate`")
	}

	root := repoRoot(t)
	bin := buildLG(t)
	corpus := filepath.Join(root, "test", "native-entry")

	fixtures, err := filepath.Glob(filepath.Join(corpus, "*.lg"))
	if err != nil {
		t.Fatalf("glob native-entry corpus: %v", err)
	}
	if len(fixtures) == 0 {
		t.Fatalf("no fixtures found under %s", corpus)
	}
	fixtures = sortedStrings(fixtures)

	for _, fixture := range fixtures {
		fixture := fixture
		name := strings.TrimSuffix(filepath.Base(fixture), ".lg")
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			runNativeEntryFixture(ctx, t, bin, root, fixture)
		})
	}
}

// nativeEntryExpect is the committed structural contract for one fixture, held
// in <fixture>.goexpect.json beside the source. It is data, not code, so the
// corpus can grow without touching the gate; regenerate a fragment by running
// scripts/lg-compile on the fixture and copying the emitted declaration.
type nativeEntryExpect struct {
	// Note explains why this fixture exists. Not asserted.
	Note string `json:"note"`
	// Namespace is the fixture's (ns …) name.
	Namespace string `json:"namespace"`
	// GeneratedFile is the lowered package file, relative to the out dir.
	GeneratedFile string `json:"generatedFile"`
	// EntryFn is the generated Go function for the fixture's -main.
	EntryFn string `json:"entryFn"`
	// LispFn / GoFn name the fixture's lowered defn on both sides.
	LispFn string `json:"lispFn"`
	GoFn   string `json:"goFn"`
	// Lowering is "direct" (native Go call in the entry path) or "vm"
	// (intentionally still on the ec.Invoke/CachedVarFn trampoline).
	Lowering string `json:"lowering"`
	// FunctionFragment is evidence 1: the whole expected func declaration.
	FunctionFragment string `json:"functionFragment"`
	// CallFragment is evidence 2: the expected call expression inside EntryFn.
	CallFragment string `json:"callFragment"`
	// SemanticWrap builds the semantic mutant: __EXPR__ is replaced by the
	// existing returned expression of GoFn, so the mutant stays compilable
	// (temporaries keep being used) while changing the produced value.
	SemanticWrap string `json:"semanticWrap"`
}

const (
	nativeEntryDirect = "direct"
	nativeEntryVM     = "vm"
	nativeEntryModule = "tmpmod"
	semanticWrapHole  = "__EXPR__"
)

func runNativeEntryFixture(ctx context.Context, t *testing.T, bin, root, fixture string) {
	t.Helper()

	stem := strings.TrimSuffix(fixture, ".lg")
	want := mustReadFixtureFile(t, stem+".expect")

	var spec nativeEntryExpect
	if err := json.Unmarshal(mustReadFixtureFile(t, stem+".goexpect.json"), &spec); err != nil {
		t.Fatalf("parse %s.goexpect.json: %v", filepath.Base(stem), err)
	}
	if spec.Lowering != nativeEntryDirect && spec.Lowering != nativeEntryVM {
		t.Fatalf("fixture %s: lowering must be %q or %q, got %q",
			filepath.Base(fixture), nativeEntryDirect, nativeEntryVM, spec.Lowering)
	}

	// --- production lowering path: scripts/lg-compile --entry-frame ---
	outDir := t.TempDir()
	if out, err := runCmd(ctx, t, bin, root,
		[]string{"scripts/lg-compile", "--entry-frame", outDir, nativeEntryModule, fixture},
		"LG_SOURCE_PATHS=pkg/rt/core"); err != nil {
		t.Fatalf("lg-compile --entry-frame: %v\n%s", err, out)
	}
	generatedPath := filepath.Join(outDir, filepath.FromSlash(spec.GeneratedFile))
	generated, err := os.ReadFile(generatedPath)
	if err != nil {
		t.Fatalf("generated package %s missing: %v", spec.GeneratedFile, err)
	}

	// --- evidence 1: the native function exists, exactly once ---
	fnOracle := gofragment.GoMatchRequest{
		Kind:     gofragment.GoFunction,
		Target:   spec.GoFn,
		Expected: spec.FunctionFragment,
	}
	if err := gofragment.MatchGeneratedGoFragment(fnOracle, string(generated)); err != nil {
		t.Fatalf("evidence 1 (native function exists): %v", err)
	}

	// --- evidence 2: the call site is classified correctly ---
	callOracle := gofragment.GoMatchRequest{
		Kind:     gofragment.GoExpression,
		Target:   spec.EntryFn,
		Expected: spec.CallFragment,
	}
	if err := gofragment.MatchGeneratedGoFragment(callOracle, string(generated)); err != nil {
		t.Fatalf("evidence 2 (%s call site in %s): %v", spec.Lowering, spec.EntryFn, err)
	}
	assertNativeEntryCallClassification(t, spec, generated)

	// --- evidence 3: it builds and runs ---
	prepareNativeEntryModule(ctx, t, bin, root, outDir, fixture)
	// Build outside outDir so the mutant copies stay source-only.
	exe := filepath.Join(t.TempDir(), "app.native")
	if out, err := runCmd(ctx, t, "go", outDir, []string{"build", "-o", exe, "."}); err != nil {
		t.Fatalf("go build generated tree: %v\n%s", err, out)
	}
	stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
	if err != nil {
		t.Fatalf("run %s: %v", exe, err)
	}
	if exit != 0 {
		t.Fatalf("binary exit %d, want 0\nstdout=%q\nstderr=%q", exit, stdout, stderr)
	}
	if !bytes.Equal(stdout, want) {
		t.Fatalf("stdout %q, want %q (from %s.expect)\nstderr=%q",
			stdout, want, filepath.Base(stem), stderr)
	}

	// --- evidence 4: falsifiability ---
	t.Run("oracle mutant", func(t *testing.T) {
		assertOracleMutantRejected(t, fnOracle, generated)
		assertOracleMutantRejected(t, callOracle, generated)
	})

	t.Run("ast call-site mutant", func(t *testing.T) {
		mutated := mutateNativeEntryCallSite(t, spec, generated)
		assertNativeEntryMutantKilled(ctx, t, outDir, generatedPath, mutated, want, "ast")
	})

	// The semantic mutant asks the same question of both kinds of fixture —
	// "is the generated Go body what actually runs?" — and demands opposite
	// answers. See the inverted-contract note in the file comment.
	switch spec.Lowering {
	case nativeEntryDirect:
		t.Run("live_native_body_mutant_dies", func(t *testing.T) {
			mutated := mutateNativeEntryReturn(t, spec, generated)
			assertNativeEntryMutantKilled(ctx, t, outDir, generatedPath, mutated, want, "semantic")
		})
	case nativeEntryVM:
		t.Run("dead_native_body_mutant_survives", func(t *testing.T) {
			mutated := mutateNativeEntryReturn(t, spec, generated)
			assertNativeEntryMutantSurvived(ctx, t, outDir, generatedPath, mutated, want, spec)
		})
	}
}

func mustReadFixtureFile(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		// Never a skip: a fixture without its sidecar is a broken corpus.
		t.Fatalf("required fixture sidecar %s: %v", filepath.Base(path), err)
	}
	return body
}

// prepareNativeEntryModule completes the generated tree into a buildable Go
// module: program.lgb beside main.go (the frame loads it), a go.mod pointing at
// this checkout, and `go mod tidy`.
func prepareNativeEntryModule(ctx context.Context, t *testing.T, bin, root, outDir, fixture string) {
	t.Helper()
	lgb := filepath.Join(outDir, "program.lgb")
	if out, err := runCmd(ctx, t, bin, root, []string{"-c", lgb, fixture}); err != nil {
		t.Fatalf("lg -c program.lgb: %v\n%s", err, out)
	}
	mod := "module " + nativeEntryModule + "\n\ngo 1.23\n\n" +
		"require github.com/nooga/let-go v0.0.0\n" +
		"replace github.com/nooga/let-go => " + root + "\n"
	if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(mod), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCmd(ctx, t, "go", outDir, []string{"mod", "tidy"}); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
}

func runNativeEntryBinary(ctx context.Context, exe string) (stdout, stderr []byte, exit int, err error) {
	cmd := exec.CommandContext(ctx, exe)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	runErr := cmd.Run()
	if runErr != nil {
		exitErr, ok := runErr.(*exec.ExitError)
		if !ok {
			return outBuf.Bytes(), errBuf.Bytes(), -1, runErr
		}
		exit = exitErr.ExitCode()
	}
	return outBuf.Bytes(), errBuf.Bytes(), exit, nil
}

// assertNativeEntryMutantKilled rebuilds a COPY of the generated tree with one
// mutated file and requires the mutant to die: the build fails, the binary
// exits non-zero, or its stdout differs from the fixture's .expect.
func assertNativeEntryMutantKilled(
	ctx context.Context,
	t *testing.T,
	outDir, generatedPath string,
	mutatedFile, want []byte,
	kind string,
) {
	t.Helper()
	mutantDir := t.TempDir()
	copyNativeEntryTree(t, outDir, mutantDir)

	rel, err := filepath.Rel(outDir, generatedPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mutantDir, rel), mutatedFile, 0o644); err != nil {
		t.Fatal(err)
	}

	exe := filepath.Join(t.TempDir(), "mutant.native")
	buildOut, buildErr := runCmd(ctx, t, "go", mutantDir, []string{"build", "-o", exe, "."})
	if buildErr != nil {
		if kind == "semantic" {
			// A semantic mutant that cannot compile proves nothing about
			// behaviour — the wrap in the sidecar must keep the tree valid.
			t.Fatalf("semantic mutant must stay compilable; go build failed:\n%s", buildOut)
		}
		return // AST mutant killed at build time, as intended.
	}

	stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
	if err != nil {
		t.Fatalf("run %s mutant: %v", kind, err)
	}
	if exit != 0 {
		return // killed at run time.
	}
	if bytes.Equal(stdout, want) {
		if kind == "semantic" {
			t.Fatalf("semantic mutant survived on a %q fixture: the lowered Go body was "+
				"mutated yet stdout %q still equals .expect. For lowering=%q the generated "+
				"Go function MUST be what runs, so mutating it MUST change observed output; "+
				"a survivor means the entry path is really going somewhere else (VM "+
				"trampoline, cached var, constant-folded value).\nstderr=%q",
				nativeEntryDirect, stdout, nativeEntryDirect, stderr)
		}
		t.Fatalf("%s mutant survived: exit 0 and stdout %q equals %s\nstderr=%q",
			kind, stdout, ".expect", stderr)
	}
}

// assertNativeEntryMutantSurvived is the INVERSE of assertNativeEntryMutantKilled
// and is required, not tolerated, for lowering=vm fixtures.
//
// A VM-backed defn is still emitted as a Go function, but nothing calls it: the
// entry path goes through the ec.Invoke trampoline instead. So mutating that Go
// body must be invisible — the binary must still build, exit 0, and print
// byte-exact .expect. That surviving mutant is the positive evidence that the
// generated body is dead code and the VM trampoline is genuinely what runs.
// If it dies, the fixture is not VM-backed in the way the sidecar claims.
func assertNativeEntryMutantSurvived(
	ctx context.Context,
	t *testing.T,
	outDir, generatedPath string,
	mutatedFile, want []byte,
	spec nativeEntryExpect,
) {
	t.Helper()
	mutantDir := t.TempDir()
	copyNativeEntryTree(t, outDir, mutantDir)

	rel, err := filepath.Rel(outDir, generatedPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mutantDir, rel), mutatedFile, 0o644); err != nil {
		t.Fatal(err)
	}

	exe := filepath.Join(t.TempDir(), "mutant.native")
	if buildOut, buildErr := runCmd(ctx, t, "go", mutantDir, []string{"build", "-o", exe, "."}); buildErr != nil {
		t.Fatalf("dead-body mutant of %s must stay compilable; go build failed:\n%s",
			spec.GoFn, buildOut)
	}

	stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
	if err != nil {
		t.Fatalf("run dead-body mutant: %v", err)
	}
	if exit != 0 {
		t.Fatalf("dead-body mutant of %s exited %d, want 0. For lowering=%q the generated "+
			"Go body is expected to be unreachable, so mutating it must not affect the "+
			"program at all.\nstdout=%q\nstderr=%q",
			spec.GoFn, exit, nativeEntryVM, stdout, stderr)
	}
	if !bytes.Equal(stdout, want) {
		t.Fatalf("dead-body mutant of %s CHANGED output: got %q, want byte-exact .expect %q. "+
			"For lowering=%q the entry path must run the ec.Invoke trampoline for %q, not "+
			"the generated Go body; if mutating that body moves the output, the body is live "+
			"and this fixture no longer scopes the direct-call claim.\nstderr=%q",
			spec.GoFn, stdout, want, nativeEntryVM, spec.LispFn, stderr)
	}
}

func copyNativeEntryTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, body, info.Mode().Perm())
	})
	if err != nil {
		t.Fatalf("copy generated tree: %v", err)
	}
}

// --- structural assertions and mutations (all located through go/ast) ---

func parseNativeEntryFile(t *testing.T, generated []byte) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	// Comments are deliberately dropped: the mutated AST is re-printed, and
	// re-attaching comments to synthesized nodes can corrupt the output.
	file, err := parser.ParseFile(fset, "generated.go", generated, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse generated Go: %v", err)
	}
	return fset, file
}

func findNativeEntryFunc(t *testing.T, file *ast.File, name string) *ast.FuncDecl {
	t.Helper()
	var found []*ast.FuncDecl
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == name {
			found = append(found, fn)
		}
	}
	if len(found) != 1 {
		t.Fatalf("generated func %q: found %d, want exactly 1", name, len(found))
	}
	return found[0]
}

// directCallsTo returns every `<name>(…)` call inside scope.
func directCallsTo(scope ast.Node, name string) []*ast.CallExpr {
	var calls []*ast.CallExpr
	ast.Inspect(scope, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == name {
			calls = append(calls, call)
		}
		return true
	})
	return calls
}

// vmTrampolinesFor returns every `ec.Invoke(…)` call inside scope whose
// argument tree mentions the fixture's lisp-level function name — i.e. the
// VM trampoline for that defn.
func vmTrampolinesFor(scope ast.Node, lispName string) []*ast.CallExpr {
	quoted := strconv.Quote(lispName)
	var calls []*ast.CallExpr
	ast.Inspect(scope, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Invoke" {
			return true
		}
		if ident, ok := sel.X.(*ast.Ident); !ok || ident.Name != "ec" {
			return true
		}
		if nativeEntryMentionsString(call, quoted) {
			calls = append(calls, call)
		}
		return true
	})
	return calls
}

func nativeEntryMentionsString(scope ast.Node, quoted string) bool {
	found := false
	ast.Inspect(scope, func(node ast.Node) bool {
		lit, ok := node.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING && lit.Value == quoted {
			found = true
		}
		return !found
	})
	return found
}

// assertNativeEntryCallClassification is the half of evidence 2 that the
// fragment oracle cannot express: for a direct fixture no VM trampoline for
// that defn may exist anywhere in the package, and for a VM-backed fixture the
// entry path must NOT contain a direct call to it.
func assertNativeEntryCallClassification(t *testing.T, spec nativeEntryExpect, generated []byte) {
	t.Helper()
	_, file := parseNativeEntryFile(t, generated)
	entry := findNativeEntryFunc(t, file, spec.EntryFn)

	switch spec.Lowering {
	case nativeEntryDirect:
		if calls := directCallsTo(entry, spec.GoFn); len(calls) != 1 {
			t.Fatalf("direct calls to %s in %s: found %d, want exactly 1",
				spec.GoFn, spec.EntryFn, len(calls))
		}
		// Scope the claim over the whole package, not just the entry fn: a
		// surviving trampoline anywhere means the defn is still VM-dispatched.
		if calls := vmTrampolinesFor(file, spec.LispFn); len(calls) != 0 {
			t.Fatalf("%s is claimed direct but %d ec.Invoke trampoline(s) for %q remain",
				spec.GoFn, len(calls), spec.LispFn)
		}
	case nativeEntryVM:
		if calls := vmTrampolinesFor(entry, spec.LispFn); len(calls) != 1 {
			t.Fatalf("ec.Invoke trampolines for %q in %s: found %d, want exactly 1",
				spec.LispFn, spec.EntryFn, len(calls))
		}
		if calls := directCallsTo(entry, spec.GoFn); len(calls) != 0 {
			t.Fatalf("%s is a VM-backed fixture but %s contains %d direct call(s) to it",
				spec.GoFn, spec.EntryFn, len(calls))
		}
	}
}

// assertOracleMutantRejected perturbs a committed expectation structurally and
// requires the matcher to reject it. Without this, a fragment that silently
// stopped being asserted would still "pass".
func assertOracleMutantRejected(t *testing.T, req gofragment.GoMatchRequest, generated []byte) {
	t.Helper()
	mutant := req
	mutant.Expected = mutateExpectedFragment(t, req.Kind, req.Expected)
	if mutant.Expected == req.Expected {
		t.Fatalf("oracle mutation produced an identical fragment for %s target %q",
			req.Kind, req.Target)
	}
	if err := gofragment.MatchGeneratedGoFragment(mutant, string(generated)); err == nil {
		t.Fatalf("oracle mutant matched: %s target %q accepted %q",
			req.Kind, req.Target, mutant.Expected)
	}
}

// mutateExpectedFragment rewrites a committed fragment through go/ast: the
// first string literal is renamed when there is one (VM trampolines are
// identified by var name), otherwise the callee / declared identifier is.
func mutateExpectedFragment(t *testing.T, kind gofragment.GoFragmentKind, fragment string) string {
	t.Helper()
	const suffix = "__oracle_mutant"
	fset := token.NewFileSet()

	var node ast.Node
	switch kind {
	case gofragment.GoFunction:
		file, err := parser.ParseFile(fset, "expected.go", "package probe\n"+fragment, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse expected function: %v", err)
		}
		if len(file.Decls) != 1 {
			t.Fatalf("expected exactly one declaration, got %d", len(file.Decls))
		}
		fn, ok := file.Decls[0].(*ast.FuncDecl)
		if !ok {
			t.Fatalf("expected function fragment, got %T", file.Decls[0])
		}
		fn.Name = ast.NewIdent(fn.Name.Name + suffix)
		node = fn
	case gofragment.GoExpression:
		expr, err := parser.ParseExprFrom(fset, "expected.go", fragment, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse expected expression: %v", err)
		}
		if !mutateFirstStringLiteral(expr, suffix) && !mutateCalleeIdent(expr, suffix) {
			t.Fatalf("no structural mutation site in expected expression %q", fragment)
		}
		node = expr
	default:
		t.Fatalf("unsupported oracle mutation kind %q", kind)
	}

	var out bytes.Buffer
	if err := format.Node(&out, token.NewFileSet(), node); err != nil {
		t.Fatalf("render mutated fragment: %v", err)
	}
	return out.String()
}

func mutateFirstStringLiteral(scope ast.Node, suffix string) bool {
	done := false
	ast.Inspect(scope, func(node ast.Node) bool {
		lit, ok := node.(*ast.BasicLit)
		if done || !ok || lit.Kind != token.STRING {
			return !done
		}
		value, err := strconv.Unquote(lit.Value)
		if err != nil {
			return true
		}
		lit.Value = strconv.Quote(value + suffix)
		done = true
		return false
	})
	return done
}

func mutateCalleeIdent(scope ast.Node, suffix string) bool {
	done := false
	ast.Inspect(scope, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if done || !ok {
			return !done
		}
		ident, ok := call.Fun.(*ast.Ident)
		if !ok {
			return true
		}
		ident.Name += suffix
		done = true
		return false
	})
	return done
}

// mutateNativeEntryCallSite rewrites the entry path so it no longer reaches the
// fixture's function: a direct call is retargeted at an undeclared function
// (build must fail), a VM trampoline is retargeted at an undefined var name
// (the run must fail). The site is found through go/ast — never by identifier
// text such as v78 or by byte offsets.
func mutateNativeEntryCallSite(t *testing.T, spec nativeEntryExpect, generated []byte) []byte {
	t.Helper()
	fset, file := parseNativeEntryFile(t, generated)
	entry := findNativeEntryFunc(t, file, spec.EntryFn)

	switch spec.Lowering {
	case nativeEntryDirect:
		calls := directCallsTo(entry, spec.GoFn)
		if len(calls) != 1 {
			t.Fatalf("call-site mutation: %d direct calls to %s in %s, want 1",
				len(calls), spec.GoFn, spec.EntryFn)
		}
		calls[0].Fun = ast.NewIdent(spec.GoFn + "__absent_mutant")
	case nativeEntryVM:
		calls := vmTrampolinesFor(entry, spec.LispFn)
		if len(calls) != 1 {
			t.Fatalf("call-site mutation: %d trampolines for %q in %s, want 1",
				len(calls), spec.LispFn, spec.EntryFn)
		}
		if !mutateTrampolineName(calls[0], spec.LispFn, "__absent_mutant") {
			t.Fatalf("call-site mutation: no %q literal under the trampoline", spec.LispFn)
		}
	}
	return renderNativeEntryFile(t, fset, file)
}

func mutateTrampolineName(call *ast.CallExpr, lispName, suffix string) bool {
	quoted := strconv.Quote(lispName)
	done := false
	ast.Inspect(call, func(node ast.Node) bool {
		lit, ok := node.(*ast.BasicLit)
		if done || !ok || lit.Kind != token.STRING || lit.Value != quoted {
			return !done
		}
		lit.Value = strconv.Quote(lispName + suffix)
		done = true
		return false
	})
	return done
}

// mutateNativeEntryReturn wraps the LAST returned expression of the fixture's
// lowered function (AST walk order, so a closure body wins over its enclosing
// return) with the sidecar's SemanticWrap. Wrapping the existing expression
// instead of replacing it keeps generated temporaries used, so the mutant
// stays compilable and can only be killed by observed behaviour.
func mutateNativeEntryReturn(t *testing.T, spec nativeEntryExpect, generated []byte) []byte {
	t.Helper()
	return wrapLastReturnResult(t, generated, spec.GoFn, spec.SemanticWrap)
}

// wrapLastReturnResult is the one structural semantic-mutation primitive in
// this package: it rewrites the LAST returned expression of fnName (AST walk
// order, so a closure body wins over its enclosing return) by substituting that
// expression into wrap at semanticWrapHole. Wrapping the existing expression
// instead of replacing it keeps generated temporaries used, so the mutant stays
// compilable and can only be killed by observed behaviour. The site is located
// through go/ast, never by a gensym-numbered temporary such as v78 or by byte
// offsets; that is why the jank direct-ABI test in jank_direct_abi_test.go shares
// this helper instead of string-replacing generated identifiers.
func wrapLastReturnResult(t *testing.T, generated []byte, fnName, wrap string) []byte {
	t.Helper()
	if !strings.Contains(wrap, semanticWrapHole) {
		t.Fatalf("semantic wrap %q must contain %s", wrap, semanticWrapHole)
	}
	fset, file := parseNativeEntryFile(t, generated)
	fn := findNativeEntryFunc(t, file, fnName)

	var last *ast.ReturnStmt
	ast.Inspect(fn, func(node ast.Node) bool {
		if ret, ok := node.(*ast.ReturnStmt); ok && len(ret.Results) > 0 {
			last = ret
		}
		return true
	})
	if last == nil {
		t.Fatalf("semantic mutation: %s has no return with results", fnName)
	}

	var current bytes.Buffer
	if err := format.Node(&current, fset, last.Results[0]); err != nil {
		t.Fatalf("render returned expression: %v", err)
	}
	wrapped := strings.ReplaceAll(wrap, semanticWrapHole, current.String())
	expr, err := parser.ParseExprFrom(token.NewFileSet(), "mutant.go", wrapped, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse semantic mutant expression %q: %v", wrapped, err)
	}
	last.Results[0] = expr
	return renderNativeEntryFile(t, fset, file)
}

func renderNativeEntryFile(t *testing.T, fset *token.FileSet, file *ast.File) []byte {
	t.Helper()
	var out bytes.Buffer
	if err := format.Node(&out, fset, file); err != nil {
		t.Fatalf("render mutated generated file: %v", err)
	}
	return out.Bytes()
}
