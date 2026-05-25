//go:build bootstrap

// lgbgen compiles core.lg and all embedded namespaces into a pre-compiled .lgb bundle.
// Usage: go run -tags bootstrap ./cmd/lgbgen [output-path]
// Default output: pkg/rt/core_compiled.lgb (when run from repo root)
package main

import (
	"fmt"
	"os"
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
	{"ir.passes.dce", &rt.IRPassDCESrc},
	{"ir.passes.constfold", &rt.IRPassConstFoldSrc},
	{"ir.passes.cse", &rt.IRPassCSESrc},
	{"ir.passes.typeinfer", &rt.IRPassTypeInferSrc},
	{"ir.passes.licm", &rt.IRPassLICMSrc},
	{"ir.build", &rt.IRBuildSrc},
	{"ir.validate", &rt.IRValidateSrc},
	{"ir.passes.pipeline", &rt.IRPassPipelineSrc},
	{"ir.dump", &rt.IRDumpSrc},
	// zip and data are loaded from source on demand (precompiled ns chunks
	// only replay nil stubs for defn, not the actual function bodies)
}

func main() {
	outPath := "pkg/rt/core_compiled.lgb"
	if len(os.Args) > 1 {
		outPath = os.Args[1]
	}

	// rt.init() has already run — native builtins are registered in CoreNS.
	consts := vm.NewConsts()
	nsChunks := make(map[string]*vm.CodeChunk)
	nsOrder := make([]string, 0, len(embeddedNS))

	// Phase F: ir.data is not in the embeddedNS list (its precompiled
	// stubs would intern nil into the `ir` namespace, breaking
	// subsequent compilation). Instead, evaluate ir.data's source
	// up-front so its intern block populates `ir/op`, `ir/refs`, etc.
	// before any IR-using namespace below compiles.
	{
		coreNS := rt.NS(rt.NameCoreNS)
		c := compiler.NewCompiler(consts, coreNS)
		c.SetSource("<embedded:ir.data:lgbgen-bootstrap>")
		if _, _, err := c.CompileMultiple(strings.NewReader(rt.IRDataSrc)); err != nil {
			fmt.Fprintf(os.Stderr, "ir.data bootstrap compilation failed: %v\n", err)
			os.Exit(1)
		}
	}

	for _, ns := range embeddedNS {
		src := *ns.src
		if src == "" {
			continue
		}
		// Use CoreNS as starting namespace — the (ns ...) form will switch to the target
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
		fmt.Printf("  compiled %-10s (%d bytes bytecode)\n", ns.name, len(chunk.Code())*4)
	}

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
