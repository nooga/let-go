/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	runtimeDebug "runtime/debug"
	"strings"
	"sync"

	"github.com/nooga/let-go/pkg/buildmeta"
	"github.com/nooga/let-go/pkg/bundle"
	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/gomod"
	"github.com/nooga/let-go/pkg/nrepl"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	wasm "github.com/nooga/let-go/pkg/rt/wasm"
	"github.com/nooga/let-go/pkg/vm"

	_ "github.com/nooga/let-go/pkg/rt/corefns"
)

func versionString() string {
	if commit != "none" && len(commit) > 7 {
		return fmt.Sprintf("%s (%s)", version, commit[:7])
	}
	return version
}

func applyBuildInfoMetadata() {
	info, ok := runtimeDebug.ReadBuildInfo()
	if !ok {
		return
	}
	version, commit = buildmeta.Resolve(version, commit, info)
}

// runtimeMetadata resolves the version/commit rt exposes as let-go.version /
// let-go.commit. For the stock lg binary let-go IS the main module and the
// host stamp describes it; a custom host importing pkg/cli stamps its OWN
// identity, so the runtime's comes from let-go's dep entry in build info
// instead.
func runtimeMetadata(hostVersion, hostCommit string) (string, string) {
	info, ok := runtimeDebug.ReadBuildInfo()
	if !ok {
		return hostVersion, hostCommit
	}
	return runtimeMetadataFrom(info, hostVersion, hostCommit)
}

func runtimeMetadataFrom(info *runtimeDebug.BuildInfo, hostVersion, hostCommit string) (string, string) {
	if info == nil || info.Main.Path == gomod.ModulePath {
		return hostVersion, hostCommit
	}
	return letgoDepMeta(info)
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

func printResult(value vm.Value) error {
	rendered, err := vm.SafeString(value)
	if err != nil {
		return err
	}
	fmt.Println(rendered)
	return nil
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
	unit, err := rt.DecodeExecUnitWithDebugFile(data, filename)
	if err != nil {
		return fmt.Errorf("decoding %s: %w", filename, err)
	}
	return rt.RunExecUnit(unit)
}

// runScript dispatches a positional script to the .lgb loader or the source
// runner. Split out so the caller's error/exit-code handling stays flat rather
// than nesting the format-dispatch inside it.
func runScript(ctx *compiler.Context, script string) error {
	if filepath.Ext(script) == ".lgb" {
		return runLGB(script)
	}
	return runFile(ctx, script)
}

// bundleBinary creates a standalone executable by copying the lg binary
// and appending the compiled LGB + footer.
func bundleBinary(ctx *compiler.Context, nsRes *resolver.NSResolver, src string, dst string, basePath string) error {
	ctx.SetSource(src)

	// Snapshot the resource roots, output path, and storage store id *before*
	// compiling — and absolutize the paths against the current cwd — because
	// CompileMultiple runs the program's top-level forms, which may change the
	// working directory. Relative roots resolved afterward would point at the
	// wrong place, and the store id (which falls back to the cwd basename for
	// main.lg/init.lg) would key off the changed directory.
	resRoots := buildResourcePaths()
	for i, r := range resRoots {
		if abs, aerr := filepath.Abs(r); aerr == nil {
			resRoots[i] = abs
		}
	}
	dstAbs, _ := filepath.Abs(dst)
	bundleStoreID := storageIDForScript(src)

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
		if err := bytecode.EncodeBundleOrderedCompressed(&lgbBuf, ctx.Consts(), nsChunks, nsOrder, compressBundle); err != nil {
			return err
		}
	} else {
		if err := bytecode.EncodeCompilationCompressed(&lgbBuf, ctx.Consts(), chunk, compressBundle); err != nil {
			return err
		}
	}
	lgbData := lgbBuf.Bytes()
	var debugData []byte
	if stripDebug {
		stripped, companion, err := bytecode.SplitDebug(lgbData)
		if err != nil {
			return fmt.Errorf("splitting debug sections: %w", err)
		}
		lgbData = stripped
		debugData = companion
		if _, err := debugCompanionPath(dst); err != nil {
			return err
		}
	}

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
	outClosed := false
	defer func() {
		if !outClosed {
			_ = out.Close()
		}
	}()

	// Copy the base binary (strip any existing bundle first)
	binSize, err := bundle.BaseBinarySize(srcBin)
	if err != nil {
		return err
	}
	srcBin.Seek(0, io.SeekStart)
	if _, err := io.CopyN(out, srcBin, binSize); err != nil {
		return err
	}

	// Append the lgb payload + optional resource archive + trailer. The store
	// id (explicit -storage-id, else the source basename) was snapshotted
	// before compilation so the bundle keys its storage by the app rather than
	// by whatever the binary is renamed to at runtime, or by a directory the
	// compiled forms chdir'd into. Mirrors the WASM build, which bakes the same id.
	if err := bundle.AppendTrailer(out, lgbData, resArc, bundleStoreID); err != nil {
		return err
	}
	closeErr := out.Close()
	outClosed = true
	if closeErr != nil {
		return closeErr
	}
	if debugData != nil {
		return writeDebugCompanion(dst, debugData)
	}
	return nil
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
	var buf bytes.Buffer
	// If namespaces were loaded during compilation, use bundle format
	if len(nsRes.LoadedChunks) > 0 {
		// Include the main chunk under its namespace name, last in order
		mainNS := ctx.CurrentNS().Name()
		nsChunks := make(map[string]*vm.CodeChunk, len(nsRes.LoadedChunks)+1)
		maps.Copy(nsChunks, nsRes.LoadedChunks)
		nsChunks[mainNS] = chunk
		nsOrder := append(nsRes.LoadOrder, mainNS)
		if err := bytecode.EncodeBundleOrderedCompressed(&buf, ctx.Consts(), nsChunks, nsOrder, compressBundle); err != nil {
			return err
		}
	} else if err := bytecode.EncodeCompilationCompressed(&buf, ctx.Consts(), chunk, compressBundle); err != nil {
		return err
	}
	data := buf.Bytes()
	var debugData []byte
	if stripDebug {
		stripped, companion, err := bytecode.SplitDebug(data)
		if err != nil {
			return fmt.Errorf("splitting debug sections: %w", err)
		}
		data = stripped
		debugData = companion
		if _, err := debugCompanionPath(dst); err != nil {
			return err
		}
	}
	if err := os.WriteFile(dst, data, 0666); err != nil {
		return err
	}
	if debugData != nil {
		return writeDebugCompanion(dst, debugData)
	}
	return nil
}

func writeDebugCompanion(artifactPath string, data []byte) error {
	path, err := debugCompanionPath(artifactPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0666); err != nil {
		return fmt.Errorf("writing debug companion %s: %w", path, err)
	}
	return nil
}

func debugCompanionPath(artifactPath string) (string, error) {
	path := debugOutput
	if path == "" {
		path = artifactPath + bytecode.DebugCompanionSuffix
	}
	artifactAbs, artifactErr := filepath.Abs(artifactPath)
	debugAbs, debugErr := filepath.Abs(path)
	if artifactErr == nil && debugErr == nil && artifactAbs == debugAbs {
		return "", fmt.Errorf("debug output must differ from artifact output %s", artifactPath)
	}
	return path, nil
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

// Build metadata, supplied by the caller through Main. The ldflags contract
// (-X main.version / -X main.commit) stays on the root package: goreleaser and
// the Makefile both target `main`, and ldflags can only reach the package that
// declares the variable.
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
var compressBundle bool
var bundleBase string
var wasmOutput string
var wasmShell string
var wasmPayload string
var wasmHostEval bool
var storageID string
var stripDebug bool
var debugOutput string
var sourcePaths string
var resourcePaths string

// flagsOnce guards registerFlags. Registration used to happen in package
// init(), which is fine for a binary that IS this package but wrong for an
// importable one: merely importing pkg/cli would mutate the global flag set,
// so a custom main could not register its own flags before calling Main.
// Deferring to Main keeps the import side-effect-free.
//
// Flags go on flag.CommandLine rather than a private FlagSet so a custom main
// CAN add its own and have them parsed. The Once only prevents the duplicate
// -registration panic if Main is called twice; it does not make a second call
// clean (see Main).
var flagsOnce sync.Once

func registerFlags() {
	flag.BoolVar(&runREPL, "r", false, "attach REPL after running given files")
	flag.StringVar(&expr, "e", "", "eval given expression")
	flag.BoolVar(&debug, "d", false, "enable VM debug mode")
	flag.BoolVar(&runNREPL, "n", false, "enable nREPL server")
	flag.IntVar(&nreplPort, "p", 2137, "set nREPL port, default is 2137")
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.StringVar(&compileOutput, "c", "", "compile .lg file to .lgb bytecode (specify output path)")
	flag.StringVar(&bundleOutput, "b", "", "bundle .lg file into a standalone executable (specify output path)")
	flag.BoolVar(&compressBundle, "z", false, "with -c/-b: DEFLATE-compress the bundle body (smaller .lgb / standalone binary; transparently inflated at load)")
	flag.StringVar(&bundleBase, "bundle-base", "", "path to target-platform lg binary for cross-OS bundling (defaults to current executable)")
	flag.BoolVar(&stripDebug, "strip", false, "split source maps and local-variable tables into a digest-bound .debug companion (with -c/-b)")
	flag.StringVar(&debugOutput, "debug-output", "", "path for the -strip debug companion (defaults to <output>.debug)")
	flag.StringVar(&wasmOutput, "w", "", "build .lg file into a WASM web app (specify output directory)")
	flag.StringVar(&wasmShell, "w-shell", "xterm", "shell for -w: 'xterm' (default), 'none' (emit core only; client supplies its own shell via window.LetGoHost), or a path to a custom HTML template containing __LG_HOST_JS_BODY_PLACEHOLDER__")
	flag.StringVar(&wasmPayload, "w-wasm", "inline", "wasm delivery for -w: 'inline' (default; gzip-base64 baked into index.html) or 'external' (emit a separate main.wasm the loader fetches + streams)")
	flag.BoolVar(&wasmHostEval, "w-host-eval", false, "for -w: expose LetGoHost.eval(code) to call into the loaded image and keep it live (park after the program's main returns); works in both boot modes. Pair with -w-shell none")
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
	registerProfileFlags()
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

func storageIDForScript(script string) string {
	cwd, _ := os.Getwd()
	exe, _ := os.Executable()
	return rt.StorageIDFrom(storageID, script, cwd, exe)
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

func emitRuntimeStats() {
	if os.Getenv("LG_LOOKUP_STATS") != "" {
		fmt.Fprint(os.Stderr, vm.SnapshotLookupStats().Summary())
	}
}

func runMain() int {
	applyBuildInfoMetadata()

	// version/commit describe the HOST binary — the stock lg, or whatever a
	// custom main passed to Main — and drive -v output. What the runtime
	// exposes as let-go.version / let-go.commit must describe let-go itself,
	// so a custom host's identity cannot masquerade as a runtime version in
	// System/getProperty feature checks.
	rt.Version, rt.Commit = runtimeMetadata(version, commit)

	// Check for appended LGB payload before anything else.
	// If found, we're a standalone binary — run it directly.
	if lgbData, resData, bakedStoreID := bundle.ReadBundledSelf(); lgbData != nil {
		// Set up resolver so embedded namespaces (string, set, etc.) can load
		ctx := initCompiler(false)
		nsResolver := resolver.NewNSResolver(ctx, buildSearchPaths())
		rt.SetNSLoader(nsResolver)
		defer rt.ShutdownAllPods()

		// A bundle skips flag parsing, so every arg after the program name is a
		// user arg. Set this before any chunk runs — top-level forms read it.
		rt.SetCommandLineArgs(os.Args[1:])
		// Prefer the store id baked into the bundle at build time. Fall back to
		// the exe basename for bundles built before the v3 trailer carried one.
		bundleStoreID := bakedStoreID
		if bundleStoreID == "" {
			bundleStoreID = storageIDForScript("")
		}
		rt.InstallPersistentStorage(bundleStoreID)

		// Resources are self-contained in a bundle: serve io/resource from the
		// embedded archive only, ignoring the filesystem and -resource-paths.
		if resData != nil {
			files, err := rt.DecodeResourceArchive(resData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: decoding embedded resources: %v\n", err)
				return 1
			}
			rt.SetResourceProvider(rt.NewEmbeddedResourceProvider(files))
		}

		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: finding executable: %v\n", err)
			return 1
		}
		unit, err := rt.DecodeExecUnitWithDebugFile(lgbData, exe)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		if err := rt.RunExecUnit(unit); err != nil {
			fmt.Fprint(os.Stderr, vm.FormatError(err))
			return 1
		}
		return 0
	}

	flag.Parse()
	if compressBundle && compileOutput == "" && bundleOutput == "" {
		fmt.Fprintln(os.Stderr, "error: -z requires -c or -b")
		return 2
	}
	if stripDebug && compileOutput == "" && bundleOutput == "" {
		fmt.Fprintln(os.Stderr, "error: -strip requires -c or -b")
		return 2
	}
	if debugOutput != "" && !stripDebug {
		fmt.Fprintln(os.Stderr, "error: -debug-output requires -strip")
		return 2
	}

	if showVersion {
		fmt.Print(bytecode.FormatVersionReport("lg", versionString()))
		return 0
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
	rt.SetCommandLineArgs(userArgs)
	scriptForStorage := ""
	if len(files) >= 1 {
		scriptForStorage = files[0]
	}
	rt.InstallPersistentStorage(storageIDForScript(scriptForStorage))

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
			return 1
		}
		if err := compileLG(context, nsResolver, files[0], compileOutput); err != nil {
			fmt.Fprint(os.Stderr, vm.FormatError(err))
			return 1
		}
		return 0
	}

	// Bundle mode: compile .lg → standalone executable
	if bundleOutput != "" {
		if len(files) != 1 {
			fmt.Fprintln(os.Stderr, "error: -b requires exactly one input file")
			return 1
		}
		if err := bundleBinary(context, nsResolver, files[0], bundleOutput, bundleBase); err != nil {
			fmt.Fprint(os.Stderr, vm.FormatError(err))
			return 1
		}
		return 0
	}

	// WASM mode: compile .lg → web app directory
	if wasmOutput != "" {
		if len(files) != 1 {
			fmt.Fprintln(os.Stderr, "error: -w requires exactly one input file")
			return 1
		}
		customShellTemplate, xtermShell, shellErr := wasm.ResolveShell(wasmShell)
		if shellErr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", shellErr)
			return 1
		}
		if wasmPayload != "inline" && wasmPayload != "external" {
			fmt.Fprintf(os.Stderr, "error: -w-wasm must be 'inline' or 'external', got %q\n", wasmPayload)
			return 1
		}
		if err := buildWasm(context, nsResolver, files[0], wasmOutput, xtermShell, wasmPayload == "external", wasmHostEval, storageIDForScript(files[0]), customShellTemplate); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		return 0
	}

	// In profiling builds, profile only the script/REPL execution below.
	// Default builds compile this to a no-op so the release binary stays small.
	startProfiling()

	// Script mode: treat only the first positional as the script to run.
	// Any further positionals belong to the script (it reads os/args).
	ranSomething := false
	runFailed := false
	if len(files) >= 1 {
		if err := runScript(context, files[0]); err != nil {
			fmt.Print(vm.FormatError(err))
			runFailed = true
		}
		ranSomething = true
	}

	if expr != "" {
		context.SetSource("EXPR")
		val, err := runForm(context, expr)
		if err == nil {
			err = printResult(val)
		}
		if err != nil {
			fmt.Print(vm.FormatError(err))
			runFailed = true
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
	// A failed script or -e expression exits nonzero, but an interactive
	// session (-r) that recovered in the REPL still exits clean.
	if runFailed && !runREPL {
		return 1
	}
	return 0
}

// Main is the whole lg command line: flag registration, argument handling, and
// every mode the binary supports (eval, run, REPL, nREPL, -c compile, -b bundle,
// -w wasm). It returns an exit code rather than calling os.Exit itself, so an
// embedding main stays in control of shutdown.
//
// Two limits worth knowing. Flag parsing uses flag.CommandLine, which is
// ExitOnError: -h and a malformed flag still terminate the process from inside
// Main, as they do for any Go command. And Main is meant to be called once per
// process — option state lives in package variables and is not reset, so a
// second call inherits the first call's flags.
//
// ver and com are the build metadata the ldflags contract puts on the root
// package; pass "dev" and "none" when you have none.
//
// A third-party binary is then just:
//
//	func main() { os.Exit(cli.Main("dev", "none")) }
//
// with the generated interop packages blank-imported alongside it.
func Main(ver, com string) int {
	version, commit = ver, com
	flagsOnce.Do(registerFlags)
	code := runMain()
	emitRuntimeStats()
	return code
}
