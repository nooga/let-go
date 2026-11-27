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

	"github.com/nooga/let-go/pkg/compiler"
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
	searchPaths := buildSearchPaths()
	nsResolver := resolver.NewNSResolver(ctx, searchPaths)
	rt.SetNSLoader(nsResolver)

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
	res, err := vm.RootExecContext.Invoke(fn, []vm.Value{vm.NewPersistentVector(argv), compileHost(nsResolver, searchPaths)})
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

// hostSourcePath is compileHost's :source-path step: the file a namespace
// loads from, or nil for one the loader would not read from disk.
func hostSourcePath(nsResolver *resolver.NSResolver) vm.Value {
	return vm.NewCtxNativeFn("source-path", func(_ *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("source-path expects 1 arg, got %d", len(vs))
		}
		ns, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("source-path: arg is %s, not a string", vs[0].Type())
		}
		p, err := programSourcePath(nsResolver, string(ns))
		if err != nil || p == "" {
			return vm.NIL, err
		}
		return vm.String(p), nil
	})
}

// hostAbsPath is compileHost's :abs-path step: a path made absolute with its
// symlinks resolved, so it compares with :source-path's answers.
func hostAbsPath() vm.Value {
	return vm.NewCtxNativeFn("abs-path", func(_ *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("abs-path expects 1 arg, got %d", len(vs))
		}
		p, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("abs-path: arg is %s, not a string", vs[0].Type())
		}
		abs, err := resolvedPath(string(p))
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(abs), nil
	})
}

// hostSourceRoots is compileHost's :source-roots step: the search paths,
// resolved as hostAbsPath resolves a path.
func hostSourceRoots(searchPaths []string) vm.Value {
	return vm.NewCtxNativeFn("source-roots", func(_ *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("source-roots expects no args, got %d", len(vs))
		}
		// An entry that does not exist roots nothing, so it is dropped.
		var roots []vm.Value
		for _, sp := range searchPaths {
			if abs, err := resolvedPath(sp); err == nil {
				roots = append(roots, vm.String(abs))
			}
		}
		return vm.NewPersistentVector(roots), nil
	})
}

// compileHost is the host map lg.compiler/build-program takes. Each fn
// returns nil or an error, which surfaces in let-go as an exception.
func compileHost(nsResolver *resolver.NSResolver, searchPaths []string) vm.Value {
	return vm.NewArrayMap([]vm.Value{
		vm.Keyword("source-path"), hostSourcePath(nsResolver),
		vm.Keyword("abs-path"), hostAbsPath(),
		vm.Keyword("source-roots"), hostSourceRoots(searchPaths),
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
		vm.Keyword("clear-lowered-core"), hostFn("clear-lowered-core", 1, func(s []string) error {
			return clearLoweredCore(s[0])
		}),
		vm.Keyword("lower-core"), hostFn("lower-core", 2, func(s []string) error {
			return lowerCoreInto(s[0], s[1], false)
		}),
		vm.Keyword("lower-core-readonly"), hostFn("lower-core-readonly", 2, func(s []string) error {
			return lowerCoreInto(s[0], s[1], true)
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
		vm.Keyword("wasm-build"), vm.NewCtxNativeFn("wasm-build", hostWasmBuild),
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

// hostWasmBuild is compileHost's :wasm-build step for lg compile -w:
// (src dir out entry module? wasm-args), as lg.compiler/build-web-program
// documents.
func hostWasmBuild(_ *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
	if len(vs) != 6 {
		return vm.NIL, fmt.Errorf("wasm-build expects 6 args, got %d", len(vs))
	}
	var s [4]string
	for i := range s {
		str, ok := vs[i].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("wasm-build: arg %d is %s, not a string", i+1, vs[i].Type())
		}
		s[i] = string(str)
	}
	module, ok := vs[4].(vm.Boolean)
	if !ok {
		return vm.NIL, fmt.Errorf("wasm-build: arg 5 is %s, not a boolean", vs[4].Type())
	}
	argVals, err := vm.SeqToSlice(vs[5])
	if err != nil {
		return vm.NIL, fmt.Errorf("wasm-build: arg 6: %w", err)
	}
	args := make([]string, len(argVals))
	for i, a := range argVals {
		str, ok := a.(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("wasm-build: flag %d is %s, not a string", i+1, a.Type())
		}
		args[i] = string(str)
	}
	return vm.NIL, wasmBuildSubprocess(s[0], s[1], s[2], s[3], bool(module), args)
}

// wasmBuildSubprocess runs this binary's own lg -w on src in a child process,
// for the reason compileLGBSubprocess gives. The child builds in dir, where
// the lowered packages already are, instead of a directory of its own, and
// calls entry from an entry frame. callerModule says dir is in a caller's
// module, which the child builds read-only.
func wasmBuildSubprocess(src, dir, out, entry string, callerModule bool, wasmArgs []string) error {
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding lg executable: %w", err)
	}
	args := append([]string{"-w", out, "-entry-frame-entry", entry}, wasmArgs...)
	cmd := exec.Command(self, append(args, src)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), preparedDirEnv+"="+dir)
	if callerModule {
		cmd.Env = append(cmd.Env, preparedModuleEnv+"=1")
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("building the web app: %w", err)
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

// programSourcePath is the absolute path of the file the resolver would load
// for namespace ns, or "" when it would load none: an embedded namespace
// (clojure.core, lg.*, a precompiled core namespace) resolves before any file,
// so it is never a program file even when a file of that name exists.
func programSourcePath(nsResolver *resolver.NSResolver, ns string) (string, error) {
	if _, ok := rt.EmbeddedSource(ns); ok {
		return "", nil
	}
	chunk, err := compiler.PrecompiledNSChunk(ns)
	if err != nil {
		return "", err
	}
	if chunk != nil {
		return "", nil
	}
	p := nsResolver.SourcePath(ns)
	if p == "" {
		return "", nil
	}
	return resolvedPath(p)
}

// resolvedPath is p made absolute with symlinks resolved, so that a program
// root and a namespace's file compare as the same tree however each was
// reached (macOS's /var is a link to /private/var, for one).
func resolvedPath(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(abs)
}
