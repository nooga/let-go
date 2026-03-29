// lgbgen compiles core.lg into a pre-compiled .lgb bytecode file.
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

func main() {
	outPath := "pkg/rt/core_compiled.lgb"
	if len(os.Args) > 1 {
		outPath = os.Args[1]
	}

	// rt.init() has already run — native builtins are registered in CoreNS.
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("<core>")

	chunk, _, err := c.CompileMultiple(strings.NewReader(rt.CoreSrc))
	if err != nil {
		fmt.Fprintf(os.Stderr, "compilation failed: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create %s: %v\n", outPath, err)
		os.Exit(1)
	}
	defer f.Close()

	if err := bytecode.EncodeCompilation(f, consts, chunk); err != nil {
		fmt.Fprintf(os.Stderr, "encode failed: %v\n", err)
		os.Exit(1)
	}

	fi, _ := f.Stat()
	fmt.Printf("wrote %s (%d bytes, %d consts)\n", outPath, fi.Size(), len(consts.Values()))
}
