// lgbgen compiles core.lg and all embedded namespaces into a pre-compiled .lgb bundle.
// Usage: go run ./cmd/lgbgen [output-path]
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
		fmt.Printf("  compiled %-10s (%d bytes bytecode)\n", ns.name, len(chunk.Code())*4)
	}

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create %s: %v\n", outPath, err)
		os.Exit(1)
	}
	defer f.Close()

	if err := bytecode.EncodeBundle(f, consts, nsChunks); err != nil {
		fmt.Fprintf(os.Stderr, "encode failed: %v\n", err)
		os.Exit(1)
	}

	fi, _ := f.Stat()
	fmt.Printf("wrote %s (%d bytes, %d consts, %d namespaces)\n",
		outPath, fi.Size(), len(consts.Values()), len(nsChunks))
}
