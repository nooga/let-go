/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/nooga/let-go/pkg/bundle"
	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/nrepl"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"

	_ "github.com/nooga/let-go/pkg/rt/corefns"
)

func versionString() string {
	if commit != "none" && len(commit) > 7 {
		return fmt.Sprintf("%s (%s)", version, commit[:7])
	}
	return version
}

func motd() {
	banner := "" +
		" " + ansiBold + " λ" + ansiReset + "   " + ansiBold + "let-go" + ansiReset + " %s\n" +
		" " + ansiBoldCyan + "GO" + ansiReset + "   " + ansiDim + bannerQuitHint + ansiReset + "\n"
	fmt.Printf(banner, versionString())
}

func runForm(ctx *compiler.Context, in string) (vm.Value, error) {
	_, val, err := ctx.CompileMultiple(strings.NewReader(in))
	if err != nil {
		return nil, err
	}
	// if debug {
	// 	val, err = vm.NewDebugFrame(chunk, nil).Run()
	// } else {
	// 	val, err = vm.NewFrame(chunk, nil).Run()
	// }
	// if err != nil {
	// 	return nil, err
	// }
	return val, err
}

func runFile(ctx *compiler.Context, filename string) error {
	ctx.SetSource(filename)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, _, err = ctx.CompileMultiple(f)
	errc := f.Close()
	if err != nil {
		return err
	}
	if errc != nil {
		return errc
	}
	return nil
}

func runLGB(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	resolve := func(nsName, name string) *vm.Var {
		n := rt.DefNSBare(nsName)
		v := n.LookupLocal(vm.Symbol(name))
		if v == nil {
			return n.DefStub(name)
		}
		return v
	}
	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(data), resolve)
	if err != nil {
		return fmt.Errorf("decoding %s: %w", filename, err)
	}

	// For bundles with multiple namespaces, execute each NS chunk in
	// dependency order (NSOrder). Skip the main chunk — it runs last.
	if len(unit.NSOrder) > 0 {
		for _, name := range unit.NSOrder {
			chunk := unit.NSChunks[name]
			if chunk == nil || chunk == unit.MainChunk {
				continue
			}
			f := vm.NewFrame(chunk, nil)
			_, err := f.RunProtected()
			vm.ReleaseFrame(f)
			if err != nil {
				return fmt.Errorf("loading namespace %s: %w", name, err)
			}
		}
	}

	f := vm.NewFrame(unit.MainChunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	return err
}

// checkBundledLGB checks if the current executable has an appended payload.
// Returns the LGB bytecode and (for a v2 bundle) the gzipped resource archive;
// both are nil when no payload is found. The resource archive is nil for a
// legacy bundle or one built without resources.
func checkBundledLGB() (lgb []byte, res []byte) {
	candidates := make([]string, 0, 3)
	if exe, err := os.Executable(); err == nil && exe != "" {
		candidates = append(candidates, exe)
	}
	if len(os.Args) > 0 && os.Args[0] != "" {
		candidates = append(candidates, os.Args[0])
	}
	candidates = append(candidates, "/proc/self/exe")

	seen := map[string]bool{}
	for _, path := range candidates {
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true
		if data, resData := bundle.ReadBundled(path); data != nil {
			return data, resData
		}
	}
	return nil, nil
}

// bundleBinary creates a standalone executable by copying the lg binary
// and appending the compiled LGB + footer.
func bundleBinary(ctx *compiler.Context, nsRes *resolver.NSResolver, src string, dst string, basePath string) error {
	ctx.SetSource(src)

	// Snapshot the resource roots and output path *before* compiling — and
	// absolutize them against the current cwd — because CompileMultiple runs
	// the program's top-level forms, which may change the working directory.
	// Relative roots resolved afterward would point at the wrong place.
	resRoots := buildResourcePaths()
	for i, r := range resRoots {
		if abs, aerr := filepath.Abs(r); aerr == nil {
			resRoots[i] = abs
		}
	}
	dstAbs, _ := filepath.Abs(dst)

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	chunk, _, err := ctx.CompileMultiple(f)
	f.Close()
	if err != nil {
		return err
	}

	// Serialize LGB to memory
	var lgbBuf bytes.Buffer
	if len(nsRes.LoadedChunks) > 0 {
		mainNS := ctx.CurrentNS().Name()
		nsChunks := make(map[string]*vm.CodeChunk, len(nsRes.LoadedChunks)+1)
		maps.Copy(nsChunks, nsRes.LoadedChunks)
		nsChunks[mainNS] = chunk
		nsOrder := append(nsRes.LoadOrder, mainNS)
		if err := bytecode.EncodeBundleOrdered(&lgbBuf, ctx.Consts(), nsChunks, nsOrder); err != nil {
			return err
		}
	} else {
		if err := bytecode.EncodeCompilation(&lgbBuf, ctx.Consts(), chunk); err != nil {
			return err
		}
	}
	lgbData := lgbBuf.Bytes()

	// Collect resources under the -resource-paths roots *before* creating the
	// output file, and exclude the output path itself — otherwise a dst that
	// lives inside a resource root would embed its own (in-progress) binary.
	// resRoots/dstAbs were snapshot before user code ran (see top of func).
	resFiles, err := bundle.CollectResources(resRoots, dstAbs)
	if err != nil {
		return fmt.Errorf("collecting resources: %w", err)
	}
	var resArc []byte
	if len(resFiles) > 0 {
		resArc, err = rt.EncodeResourceArchive(resFiles)
		if err != nil {
			return fmt.Errorf("encoding resources: %w", err)
		}
	}

	// Base binary: user-supplied target (for cross-OS bundling) or our own exe.
	if basePath == "" {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("finding executable: %w", err)
		}
		basePath = exe
	}
	srcBin, err := os.Open(basePath)
	if err != nil {
		return err
	}
	defer srcBin.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the base binary (strip any existing bundle first)
	binSize, err := bundle.BaseBinarySize(srcBin)
	if err != nil {
		return err
	}
	srcBin.Seek(0, io.SeekStart)
	if _, err := io.CopyN(out, srcBin, binSize); err != nil {
		return err
	}

	// Append the lgb payload + optional resource archive + trailer.
	return bundle.AppendTrailer(out, lgbData, resArc)
}

func compileLG(ctx *compiler.Context, nsRes *resolver.NSResolver, src string, dst string) error {
	ctx.SetSource(src)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	chunk, _, err := ctx.CompileMultiple(f)
	f.Close()
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// If namespaces were loaded during compilation, use bundle format
	if len(nsRes.LoadedChunks) > 0 {
		// Include the main chunk under its namespace name, last in order
		mainNS := ctx.CurrentNS().Name()
		nsChunks := make(map[string]*vm.CodeChunk, len(nsRes.LoadedChunks)+1)
		maps.Copy(nsChunks, nsRes.LoadedChunks)
		nsChunks[mainNS] = chunk
		nsOrder := append(nsRes.LoadOrder, mainNS)
		return bytecode.EncodeBundleOrdered(out, ctx.Consts(), nsChunks, nsOrder)
	}
	return bytecode.EncodeCompilation(out, ctx.Consts(), chunk)
}

var nreplServer *nrepl.NreplServer

func nreplServe(ctx *compiler.Context, port int) error {
	nreplServer = nrepl.NewNreplServer(ctx)
	err := nreplServer.Start(port)
	if err != nil {
		return err
	}
	return nil
}

// Set by goreleaser via ldflags
var (
	version = "dev"
	commit  = "none"
)

var nreplPort int
var runNREPL bool
var runREPL bool
var expr string
var debug bool
var showVersion bool
var compileOutput string
var bundleOutput string
var bundleBase string
var wasmOutput string
var wasmShell string
var wasmPayload string
var storageID string
var sourcePaths string
var resourcePaths string

func init() {
	flag.BoolVar(&runREPL, "r", false, "attach REPL after running given files")
	flag.StringVar(&expr, "e", "", "eval given expression")
	flag.BoolVar(&debug, "d", false, "enable VM debug mode")
	flag.BoolVar(&runNREPL, "n", false, "enable nREPL server")
	flag.IntVar(&nreplPort, "p", 2137, "set nREPL port, default is 2137")
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.StringVar(&compileOutput, "c", "", "compile .lg file to .lgb bytecode (specify output path)")
	flag.StringVar(&bundleOutput, "b", "", "bundle .lg file into a standalone executable (specify output path)")
	flag.StringVar(&bundleBase, "bundle-base", "", "path to target-platform lg binary for cross-OS bundling (defaults to current executable)")
	flag.StringVar(&wasmOutput, "w", "", "build .lg file into a WASM web app (specify output directory)")
	flag.StringVar(&wasmShell, "w-shell", "xterm", "shell for -w: 'xterm' (default) or 'none' (emit core only; client supplies its own shell via window.LetGoHost)")
	flag.StringVar(&wasmPayload, "w-wasm", "inline", "wasm delivery for -w: 'inline' (default; gzip-base64 baked into index.html) or 'external' (emit a separate main.wasm the loader fetches + streams)")
	flag.StringVar(&storageID, "storage-id", "", "logical storage store id for the storage namespace (default: script name, or current directory for main.lg)")
	flag.StringVar(&sourcePaths, "source-paths", "",
		"namespace search paths separated by the OS path-list separator "+
			"(':' on Unix, ';' on Windows). When given, fully defines the search "+
			"path: the current directory is NOT searched implicitly — include '.' "+
			"to search it. Falls back to LG_SOURCE_PATHS if unset. "+
			"If flag or env var not given, it defaults to '.'")
	flag.StringVar(&resourcePaths, "resource-paths", "",
		"resource root directories for io/resource, separated by the OS path-list "+
			"separator (':' on Unix, ';' on Windows). Falls back to LG_RESOURCE_PATHS "+
			"if unset. With -b, resources under these roots are embedded in the binary.")
}

// buildSearchPaths resolves the resolver's path list from the -source-paths
// flag (preferred), the LG_SOURCE_PATHS env var, or deps.edn in the current
// directory (fallback). When the path is supplied explicitly — the
// -source-paths flag is present, or LG_SOURCE_PATHS is set (even to an empty
// value) — it fully defines the search path: "." is NOT included implicitly
// (list it to search the current directory), and an empty value yields no
// paths. Only a truly absent env var with no flag falls through to deps.edn
// and the "." default. Presence is detected the same way on both channels:
// flag.Visit for the flag, os.LookupEnv for the env var.
func buildSearchPaths() []string {
	explicitSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "source-paths" {
			explicitSet = true
		}
	})
	envVal, envSet := os.LookupEnv("LG_SOURCE_PATHS")
	if explicitSet || envSet {
		return resolver.PathsFromInputs(sourcePaths, envVal, explicitSet)
	}
	if depsPaths := resolver.PathsFromDepsEdn("."); depsPaths != nil {
		return append([]string{"."}, depsPaths...)
	}
	return []string{"."}
}

// buildResourcePaths resolves the io/resource search roots from the
// -resource-paths flag (preferred) or the LG_RESOURCE_PATHS env var. Unlike
// buildSearchPaths it is explicit-only: it does NOT prepend "." and does NOT
// consult deps.edn. Returns nil when neither is set. Project-level config
// (e.g. a conventional resources/ dir) is owned by external tools,
// which passes this flag.
func buildResourcePaths() []string {
	// An explicit -resource-paths wins even when empty, so `-resource-paths ""`
	// clears the LG_RESOURCE_PATHS fallback (the flag is documented as
	// preferred). Mirrors buildSearchPaths.
	explicitSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "resource-paths" {
			explicitSet = true
		}
	})
	raw := resourcePaths
	if !explicitSet {
		raw = os.Getenv("LG_RESOURCE_PATHS")
	}
	return resolver.ParseSearchPaths(raw)
}

// commandLineArgsValue converts the user's CLI args — the positionals after
// the script — into the value of core/*command-line-args*: nil when there are
// none, else a seq of strings.
func commandLineArgsValue(args []string) vm.Value {
	if len(args) == 0 {
		return vm.NIL
	}
	vs := make([]vm.Value, len(args))
	for i, a := range args {
		vs[i] = vm.String(a)
	}
	return vm.NewList(vs)
}

// setCommandLineArgs publishes the user's CLI args to core/*command-line-args*.
// lg is the only layer that knows authoritatively where the script ends and
// the user's args begin, so it computes them once and every consumer reads the
// var instead of slicing os/args by hand.
func setCommandLineArgs(args []string) {
	rt.CoreNS.Lookup("*command-line-args*").(*vm.Var).SetRoot(commandLineArgsValue(args))
}

func storageIDForScript(script string) string {
	cwd, _ := os.Getwd()
	exe, _ := os.Executable()
	return rt.StorageIDFrom(storageID, script, cwd, exe)
}

func installPersistentStorage(storeID string) {
	store, err := rt.NewDefaultFileStorage(storeID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: storage disabled: %v\n", err)
		return
	}
	if v := rt.LookupCoreVar("*storage*"); v != nil {
		v.SetRoot(vm.NewBoxed(store))
	}
}

func initCompiler(debug bool) *compiler.Context {
	consts := vm.NewConsts()
	ns := rt.NS("user")
	if ns == nil {
		fmt.Println("namespace not found")
		return nil
	}
	if debug {
		return compiler.NewDebugCompiler(consts, ns)
	} else {
		return compiler.NewCompiler(consts, ns)
	}
}

func main() {
	// Propagate version metadata to runtime so System/getProperty exposes it.
	rt.Version = version
	rt.Commit = commit

	// Check for appended LGB payload before anything else.
	// If found, we're a standalone binary — run it directly.
	if lgbData, resData := checkBundledLGB(); lgbData != nil {
		// Set up resolver so embedded namespaces (string, set, etc.) can load
		ctx := initCompiler(false)
		nsResolver := resolver.NewNSResolver(ctx, buildSearchPaths())
		rt.SetNSLoader(nsResolver)
		defer rt.ShutdownAllPods()

		// A bundle skips flag parsing, so every arg after the program name is a
		// user arg. Set this before any chunk runs — top-level forms read it.
		setCommandLineArgs(os.Args[1:])
		installPersistentStorage(storageIDForScript(""))

		// Resources are self-contained in a bundle: serve io/resource from the
		// embedded archive only, ignoring the filesystem and -resource-paths.
		if resData != nil {
			files, err := rt.DecodeResourceArchive(resData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: decoding embedded resources: %v\n", err)
				os.Exit(1)
			}
			rt.SetResourceProvider(rt.NewEmbeddedResourceProvider(files))
		}

		resolve := func(nsName, name string) *vm.Var {
			n := rt.DefNSBare(nsName)
			v := n.LookupLocal(vm.Symbol(name))
			if v == nil {
				return n.DefStub(name)
			}
			return v
		}
		unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(lgbData), resolve)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		// Execute namespace chunks in dependency order before main
		for _, name := range unit.NSOrder {
			chunk := unit.NSChunks[name]
			if chunk == nil || chunk == unit.MainChunk {
				continue
			}
			f := vm.NewFrame(chunk, nil)
			_, err := f.RunProtected()
			vm.ReleaseFrame(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: loading namespace %s: %v\n", name, err)
				os.Exit(1)
			}
		}
		f := vm.NewFrame(unit.MainChunk, nil)
		_, err = f.RunProtected()
		vm.ReleaseFrame(f)
		if err != nil {
			fmt.Fprint(os.Stderr, vm.FormatError(err))
			os.Exit(1)
		}
		return
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("lg %s\n", versionString())
		os.Exit(0)
	}

	files := flag.Args()

	// Ensure all pods are shut down on exit
	defer rt.ShutdownAllPods()

	context := initCompiler(debug)
	nsResolver := resolver.NewNSResolver(context, buildSearchPaths())
	rt.SetNSLoader(nsResolver)

	// files[0] is the script; the rest are the user's args. Set unconditionally
	// so script, -e, compile, bundle, and wasm modes all see it, and before any
	// user code runs.
	var userArgs []string
	if len(files) >= 1 {
		userArgs = files[1:]
	}
	setCommandLineArgs(userArgs)
	scriptForStorage := ""
	if len(files) >= 1 {
		scriptForStorage = files[0]
	}
	installPersistentStorage(storageIDForScript(scriptForStorage))

	// Dev/run resources: serve io/resource from the -resource-paths roots on
	// the filesystem. (In a bundled binary this branch is never reached — the
	// embedded provider is installed earlier, before flag.Parse.)
	if rp := buildResourcePaths(); len(rp) > 0 {
		rt.SetResourceProvider(rt.NewFSResourceProvider(rp))
	}

	// Compile mode: compile .lg → .lgb
	if compileOutput != "" || bundleOutput != "" || wasmOutput != "" {
		// Set *compiling-aot* so user code can detect AOT compilation
		rt.CoreNS.Lookup("*compiling-aot*").(*vm.Var).SetRoot(vm.TRUE)
	}
	if compileOutput != "" {
		if len(files) != 1 {
			fmt.Fprintln(os.Stderr, "error: -c requires exactly one input file")
			os.Exit(1)
		}
		if err := compileLG(context, nsResolver, files[0], compileOutput); err != nil {
			fmt.Fprint(os.Stderr, vm.FormatError(err))
			os.Exit(1)
		}
		return
	}

	// Bundle mode: compile .lg → standalone executable
	if bundleOutput != "" {
		if len(files) != 1 {
			fmt.Fprintln(os.Stderr, "error: -b requires exactly one input file")
			os.Exit(1)
		}
		if err := bundleBinary(context, nsResolver, files[0], bundleOutput, bundleBase); err != nil {
			fmt.Fprint(os.Stderr, vm.FormatError(err))
			os.Exit(1)
		}
		return
	}

	// WASM mode: compile .lg → web app directory
	if wasmOutput != "" {
		if len(files) != 1 {
			fmt.Fprintln(os.Stderr, "error: -w requires exactly one input file")
			os.Exit(1)
		}
		if wasmShell != "xterm" && wasmShell != "none" {
			fmt.Fprintf(os.Stderr, "error: -w-shell must be 'xterm' or 'none', got %q\n", wasmShell)
			os.Exit(1)
		}
		if wasmPayload != "inline" && wasmPayload != "external" {
			fmt.Fprintf(os.Stderr, "error: -w-wasm must be 'inline' or 'external', got %q\n", wasmPayload)
			os.Exit(1)
		}
		if err := buildWasm(context, nsResolver, files[0], wasmOutput, wasmShell == "xterm", wasmPayload == "external", storageIDForScript(files[0])); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// In profiling builds, profile only the script/REPL execution below.
	// Default builds compile this to a no-op so the release binary stays small.
	startProfiling()

	// Script mode: treat only the first positional as the script to run.
	// Any further positionals belong to the script (it reads os/args).
	ranSomething := false
	if len(files) >= 1 {
		script := files[0]
		if filepath.Ext(script) == ".lgb" {
			if err := runLGB(script); err != nil {
				fmt.Print(vm.FormatError(err))
			}
		} else {
			if err := runFile(context, script); err != nil {
				fmt.Print(vm.FormatError(err))
			}
		}
		ranSomething = true
	}

	if expr != "" {
		context.SetSource("EXPR")
		val, err := runForm(context, expr)
		if err != nil {
			fmt.Print(vm.FormatError(err))
		} else {
			fmt.Println(val)
		}
		ranSomething = true
	}

	if !ranSomething || runREPL {
		motd()
		if runNREPL {
			err := nreplServe(context, nreplPort)
			if err != nil {
				fmt.Println("failed to run nREPL server on port", nreplPort, err)
			}
			fmt.Printf("nREPL server running at tcp://127.0.0.1:%d\n", nreplServer.Port())
		}
		repl(context)
	}

	stopProfiling()
}
