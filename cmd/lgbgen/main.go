//go:build bootstrap

// lgbgen compiles core.lg and all embedded namespaces into a pre-compiled .lgb bundle,
// or with --target=go, generates Go source files from the Go-lowering pipeline.
//
// Usage:
//
//	go run -tags bootstrap ./cmd/lgbgen [output-path]              # .lgb bundle (default)
//	go run -tags bootstrap ./cmd/lgbgen --target=go <out-dir>       # Go source bootstrap
//
// Default .lgb output: pkg/rt/core_compiled.lgb (when run from repo root)
// Default Go output dir: pkg/rt/core_go_lowered
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// embeddedNS lists all embedded namespaces in compilation order.
// Dependencies must come before dependents (test depends on walk).
var embeddedNS = []struct {
	name string
	src  *string
}{
	{"core", &rt.CoreSrc},
	{"walk", &rt.WalkSrc},
	{"string", &rt.StringSrc},
	{"set", &rt.SetSrc},
	{"pprint", &rt.PprintSrc},
	{"edn", &rt.EdnSrc},
	{"io", &rt.IoSrc},
	{"async", &rt.AsyncSrc},
	{"test", &rt.TestSrc}, // depends on walk — must come after
	// ir.data is loaded from source on demand (like `zip` and `data`)
	// because precompiled ns chunks only replay nil stubs for defn,
	// not function bodies — the intern block at the bottom of data.lg
	// must run with live function values, which only happens via the
	// source-load path in the resolver's loadEmbedded.
	{"ir.zipper", &rt.IRZipperSrc},
	{"ir.passes", &rt.IRPassesSrc},
	{"ir.dominance", &rt.IRDominanceSrc},
	{"ir.lower", &rt.IRLowerSrc},
	{"ir.lower-go", &rt.IRLowerGoSrc},
	{"ir.validate", &rt.IRValidateSrc},
	{"ir.passes.dce", &rt.IRPassDCESrc},
	{"ir.passes.constfold", &rt.IRPassConstFoldSrc},
	{"ir.passes.cse", &rt.IRPassCSESrc},
	{"ir.passes.typeinfer", &rt.IRPassTypeInferSrc},
	{"ir.passes.licm", &rt.IRPassLICMSrc},
	{"ir.build", &rt.IRBuildSrc},
	{"ir.passes.pipeline", &rt.IRPassPipelineSrc},
	{"ir.dump", &rt.IRDumpSrc},
	// zip and data are loaded from source on demand (precompiled ns chunks
	// only replay nil stubs for defn, not the actual function bodies)
}

func main() {
	// Parse mode: --target=go <out-dir> or default bytecode mode
	targetGo := false
	outPath := "pkg/rt/core_compiled.lgb"
	goOutDir := "pkg/rt/core_go_lowered"

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--target=go":
			targetGo = true
			if i+1 < len(args) {
				goOutDir = args[i+1]
				i++
			}
		default:
			outPath = args[i]
		}
	}

	// rt.init() has already run — native builtins are registered in CoreNS.
	consts := vm.NewConsts()

	// Phase 0: bootstrap ir.data (same for both modes).
	{
		coreNS := rt.NS(rt.NameCoreNS)
		c := compiler.NewCompiler(consts, coreNS)
		c.SetSource("<embedded:ir.data:lgbgen-bootstrap>")
		if _, _, err := c.CompileMultiple(strings.NewReader(rt.IRDataSrc)); err != nil {
			fmt.Fprintf(os.Stderr, "ir.data bootstrap compilation failed: %v\n", err)
			os.Exit(1)
		}
	}

	// Phase 1: compile all namespaces from source (bytecode, sets up VM state).
	nsChunks := make(map[string]*vm.CodeChunk)
	nsOrder := make([]string, 0, len(embeddedNS))

	for _, ns := range embeddedNS {
		src := *ns.src
		if src == "" {
			continue
		}
		coreNS := rt.NS(rt.NameCoreNS)
		c := compiler.NewCompiler(consts, coreNS)
		c.SetSource("<embedded:" + ns.name + ">")

		chunk, _, err := c.CompileMultiple(strings.NewReader(src))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s compilation failed: %v\n", ns.name, err)
			os.Exit(1)
		}
		nsChunks[ns.name] = chunk
		nsOrder = append(nsOrder, ns.name)
		if targetGo {
			fmt.Printf("  compiled %-20s (%d bytecode + lowering to Go)\n", ns.name, len(chunk.Code())*4)
		} else {
			fmt.Printf("  compiled %-10s (%d bytes bytecode)\n", ns.name, len(chunk.Code())*4)
		}
	}

	if targetGo {
		runGoTarget(goOutDir)
		return
	}

	// Bytecode mode: write .lgb bundle.
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create %s: %v\n", outPath, err)
		os.Exit(1)
	}
	defer f.Close()

	if err := bytecode.EncodeBundleOrdered(f, consts, nsChunks, nsOrder); err != nil {
		fmt.Fprintf(os.Stderr, "encode failed: %v\n", err)
		os.Exit(1)
	}

	fi, _ := f.Stat()
	fmt.Printf("wrote %s (%d bytes, %d consts, %d namespaces)\n",
		outPath, fi.Size(), len(consts.Values()), len(nsChunks))
}

// nsToGoPkgName converts a let-go namespace name to a valid Go package name.
// Dots become underscores (ir.zipper → ir_zipper), hyphens become underscores.
func nsToGoPkgName(nsName string) string {
	s := strings.ReplaceAll(nsName, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

// runGoTarget re-pipelines each namespace's defn forms through the Go
// lowering pipeline and writes .go source files to outDir.
func runGoTarget(outDir string) {
	var failed []string
	pipelineVar := rt.LookupVar("ir.passes.pipeline", "lower-ns-to-go")
	if pipelineVar == nil {
		fmt.Fprintf(os.Stderr, "fatal: ir.passes.pipeline/lower-ns-to-go not found\n")
		os.Exit(1)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir %s: %v\n", outDir, err)
		os.Exit(1)
	}

	for _, ns := range embeddedNS {
		src := *ns.src
		if src == "" {
			continue
		}

		// Read forms from source and pick out defn forms only.
		// defmacro forms are skipped for now — their bodies are macro
		// template construction code that doesn't lower cleanly.
		r := compiler.NewLispReader(strings.NewReader(src), "<embedded:"+ns.name+">")
		var defnForms []vm.Value
		for {
			form, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Fprintf(os.Stderr, "%s: read error: %v\n", ns.name, err)
				os.Exit(1)
			}
			if isDefnOnly(form) {
				defnForms = append(defnForms, form)
			}
		}

		if len(defnForms) == 0 {
			fmt.Printf("  skip %-20s (no defn forms)\n", ns.name)
			continue
		}

		// Each namespace maps to its own Go package in a subdirectory.
		pkgName := nsToGoPkgName(ns.name)
		pkgDir := filepath.Join(outDir, pkgName)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "create pkg dir %s: %v\n", pkgDir, err)
			os.Exit(1)
		}

		// Batch all forms into one namespace file.
		formsVec := vm.NewPersistentVector(defnForms)
		args := []vm.Value{
			vm.String(pkgName),
			vm.Symbol(ns.name),
			formsVec,
		}
		result, err := pipelineVar.Invoke(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Go lowering failed: %v\n", ns.name, err)
			failed = append(failed, ns.name)
			continue
		}
		goSrc, ok := result.(vm.String)
		if !ok {
			fmt.Fprintf(os.Stderr, "%s: expected string result, got %T\n", ns.name, result)
			failed = append(failed, ns.name)
			continue
		}
		filename := filepath.Join(pkgDir, pkgName+".go")
		if err := os.WriteFile(filename, []byte(goSrc), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "%s: write %s: %v\n", ns.name, filename, err)
			os.Exit(1)
		}
		fmt.Printf("  wrote %s (%d bytes, %d defns)\n", filename, len(goSrc), len(defnForms))
	}
	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "lgbgen: %d namespace(s) failed: %v\n", len(failed), failed)
		os.Exit(1)
	}
}

// isDefnLike returns true if form is a list beginning with defn or defmacro.
func isDefnLike(form vm.Value) bool {
	return isDefnOnly(form) || isDefmacroForm(form)
}

func isDefmacroForm(form vm.Value) bool {
	list, ok := form.(vm.Sequable)
	if !ok {
		return false
	}
	seq := list.Seq()
	if seq == nil {
		return false
	}
	head := seq.First()
	sym, ok := head.(vm.Symbol)
	if !ok {
		return false
	}
	return string(sym) == "defmacro"
}

// isDefnOnly returns true if form is a list beginning with defn (not defmacro).
func isDefnOnly(form vm.Value) bool {
	list, ok := form.(vm.Sequable)
	if !ok {
		return false
	}
	seq := list.Seq()
	if seq == nil {
		return false
	}
	head := seq.First()
	sym, ok := head.(vm.Symbol)
	if !ok {
		return false
	}
	return string(sym) == "defn"
}
