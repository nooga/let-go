/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nooga/let-go/pkg/gomod"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// compileCommand is the argv[1] that selects `lg compile`. Only an
// exact first argument dispatches, so every existing invocation — a script
// path, -e, flags, bare lg for the REPL — reaches the flag parser unchanged.
// A script literally named "compile" now needs a path prefix: lg ./compile.
const compileCommand = "compile"

// runCompile runs `lg compile`. Argument parsing, orchestration, and
// diagnostics are let-go code (lg.commands.compile over lg.compiler);
// this side boots the runtime and supplies the steps that need this binary or
// the Go toolchain, as the host map lg.compiler/build-program documents.
func runCompile(args []string) int {
	defer rt.ShutdownAllPods()
	ctx := initCompiler(false)
	rt.SetNSLoader(resolver.NewNSResolver(ctx, buildSearchPaths()))

	if _, err := runForm(ctx, "(require 'lg.commands.compile)"); err != nil {
		fmt.Fprint(os.Stderr, vm.FormatError(err))
		return 1
	}
	ns := rt.LookupNS("lg.commands.compile")
	if ns == nil {
		fmt.Fprintln(os.Stderr, "error: lg.commands.compile did not load")
		return 1
	}
	var fn vm.Fn
	if mainVar := ns.LookupLocal(vm.Symbol("main")); mainVar != nil {
		fn, _ = mainVar.Deref().(vm.Fn)
	}
	if fn == nil {
		fmt.Fprintln(os.Stderr, "error: lg.commands.compile/main is not a function")
		return 1
	}

	argv := make([]vm.Value, len(args))
	for i, a := range args {
		argv[i] = vm.String(a)
	}
	res, err := vm.RootExecContext.Invoke(fn, []vm.Value{vm.NewPersistentVector(argv), compileHost()})
	if err != nil {
		fmt.Fprint(os.Stderr, vm.FormatError(err))
		return 1
	}
	code, ok := res.(vm.Int)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: lg.commands.compile/main returned %s, not an exit code\n", res.Type())
		return 1
	}
	return int(code)
}

// compileHost is the host map lg.compiler/build-program takes. Each fn
// returns nil or an error, which surfaces in let-go as an exception.
func compileHost() vm.Value {
	return vm.NewArrayMap([]vm.Value{
		vm.Keyword("compile-lgb"), hostFn("compile-lgb", 3, func(s []string) error {
			return compileLGBSubprocess(s[0], s[1], s[2])
		}),
		vm.Keyword("scaffold-module"), hostFn("scaffold-module", 2, func(s []string) error {
			return writeGoModule(s[0], s[1])
		}),
		vm.Keyword("go-build"), hostFn("go-build", 2, func(s []string) error {
			return goBuild(s[0], s[1], false)
		}),
		vm.Keyword("open-module"), vm.NewCtxNativeFn("open-module", func(_ *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
			if len(vs) != 2 {
				return vm.NIL, fmt.Errorf("open-module expects 2 args, got %d", len(vs))
			}
			dir, ok := vs[0].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("open-module: arg 1 is %s, not a string", vs[0].Type())
			}
			pkgs, err := vm.SeqToSlice(vs[1])
			if err != nil {
				return vm.NIL, fmt.Errorf("open-module: arg 2: %w", err)
			}
			imports := make([]string, len(pkgs))
			for i, p := range pkgs {
				str, ok := p.(vm.String)
				if !ok {
					return vm.NIL, fmt.Errorf("open-module: import %d is %s, not a string", i+1, p.Type())
				}
				imports[i] = string(str)
			}
			m, err := openCallerModule(string(dir), imports)
			if err != nil {
				return vm.NIL, err
			}
			return vm.NewArrayMap([]vm.Value{
				vm.Keyword("dir"), vm.String(m.GenDir()),
				vm.Keyword("import-path"), vm.String(m.ImportPath()),
			}), nil
		}),
		vm.Keyword("go-build-readonly"), hostFn("go-build-readonly", 2, func(s []string) error {
			return goBuild(s[0], s[1], true)
		}),
		vm.Keyword("make-temp-dir"), vm.NewCtxNativeFn("make-temp-dir", func(_ *vm.ExecContext, _ []vm.Value) (vm.Value, error) {
			dir, err := os.MkdirTemp("", "lg-compile-*")
			if err != nil {
				return vm.NIL, err
			}
			return vm.String(dir), nil
		}),
		vm.Keyword("remove-dir"), hostFn("remove-dir", 1, func(s []string) error {
			return os.RemoveAll(s[0])
		}),
	})
}

// hostFn wraps a step taking n string arguments as a let-go fn returning nil.
func hostFn(name string, n int, step func([]string) error) *vm.NativeFn {
	return vm.NewCtxNativeFn(name, func(_ *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != n {
			return vm.NIL, fmt.Errorf("%s expects %d args, got %d", name, n, len(vs))
		}
		s := make([]string, n)
		for i, v := range vs {
			str, ok := v.(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("%s: arg %d is %s, not a string", name, i+1, v.Type())
			}
			s[i] = string(str)
		}
		return vm.NIL, step(s)
	})
}

// compileLGBSubprocess writes program.lgb by running this binary's own -c in a
// child process. It cannot run in-process: lowering has already loaded ir.*,
// gogen and lg.compiler into this process's resolver, and compileLG bundles
// every namespace the resolver loaded.
func compileLGBSubprocess(src, dst, entry string) error {
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding lg executable: %w", err)
	}
	cmd := exec.Command(self, "-c", dst, "-entry-frame-entry", entry, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compiling %s to bytecode: %w", src, err)
	}
	return nil
}

// goBuild builds the generated package in dir to out. It skips `go mod tidy`
// for the reason -w does: writeGoModule already resolved the one require the
// module has, and tidy would also walk let-go's test-only dependencies.
// readonly is for a caller-supplied module: an explicit -mod=readonly keeps
// GOFLAGS=-mod=mod from editing the caller's go.mod.
func goBuild(dir, out string, readonly bool) error {
	abs, err := filepath.Abs(out)
	if err != nil {
		return err
	}
	// go build -o into an existing directory writes <dir>/<module name>
	// instead, which would leave the binary somewhere other than reported.
	if fi, err := os.Stat(abs); err == nil && fi.IsDir() {
		return fmt.Errorf("output %s is a directory; name the binary with -o", out)
	}
	args := []string{"build"}
	if readonly {
		args = append(args, "-mod=readonly")
	}
	cmd := exec.Command(gomod.GoToolPath(), append(args, "-o", abs, ".")...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}
	return nil
}
