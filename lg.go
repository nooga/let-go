/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/nrepl"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"

	_ "github.com/nooga/let-go/pkg/rt/corefns"
)

// Trailers appended to standalone binaries.
//
// Legacy (no resources): [lgb data][8-byte lgbSize][4-byte 'LGBX']  (12-byte trailer)
// v2    (resources):     [lgb data][resource archive][8-byte lgbSize][8-byte resSize][4-byte 'LGB2']  (20-byte trailer)
//
// A v2 trailer is written only when the bundle carries resources; resource-less
// bundles keep the byte-identical legacy trailer. Readers recognize both.
var bundleMagic = [4]byte{'L', 'G', 'B', 'X'}
var bundleMagicV2 = [4]byte{'L', 'G', 'B', '2'}

// bundleKind classifies a standalone binary's appended trailer.
type bundleKind int

const (
	bundleNone   bundleKind = iota // no recognized trailer (a plain, non-bundled binary)
	bundleLegacy                   // 12-byte 'LGBX' trailer (lgb only)
	bundleV2                       // 20-byte 'LGB2' trailer (lgb + resource archive)
)

// trailerLen returns the on-disk size of the trailer for this kind.
func (k bundleKind) trailerLen() int64 {
	switch k {
	case bundleV2:
		return 20
	case bundleLegacy:
		return 12
	default:
		return 0
	}
}

// parseBundleTrailer reads and validates the trailer appended to f, the single
// place that discriminates the LGB2/LGBX formats. It returns bundleNone with a
// nil error when f carries no recognized trailer (a normal, non-bundled
// binary). For a recognized trailer it validates that the claimed payload fits
// within the file and returns a "corrupt bundle" error otherwise — so callers
// never seek to a bogus offset or allocate a garbage-sized slice.
func parseBundleTrailer(f *os.File) (kind bundleKind, lgbSize, resSize uint64, err error) {
	fi, err := f.Stat()
	if err != nil {
		return bundleNone, 0, 0, err
	}
	total := fi.Size()
	if total < bundleLegacy.trailerLen() {
		return bundleNone, 0, 0, nil
	}

	// Discriminate by the trailing 4-byte magic.
	if _, err := f.Seek(-4, io.SeekEnd); err != nil {
		return bundleNone, 0, 0, err
	}
	var magic [4]byte
	if _, err := io.ReadFull(f, magic[:]); err != nil {
		return bundleNone, 0, 0, err
	}

	switch magic {
	case bundleMagicV2:
		if total < bundleV2.trailerLen() {
			return bundleNone, 0, 0, nil
		}
		if _, err := f.Seek(-bundleV2.trailerLen(), io.SeekEnd); err != nil {
			return bundleNone, 0, 0, err
		}
		var tr [20]byte
		if _, err := io.ReadFull(f, tr[:]); err != nil {
			return bundleNone, 0, 0, err
		}
		kind = bundleV2
		lgbSize = binary.LittleEndian.Uint64(tr[0:8])
		resSize = binary.LittleEndian.Uint64(tr[8:16])
	case bundleMagic:
		if _, err := f.Seek(-bundleLegacy.trailerLen(), io.SeekEnd); err != nil {
			return bundleNone, 0, 0, err
		}
		var tr [12]byte
		if _, err := io.ReadFull(f, tr[:]); err != nil {
			return bundleNone, 0, 0, err
		}
		kind = bundleLegacy
		lgbSize = binary.LittleEndian.Uint64(tr[0:8])
	default:
		return bundleNone, 0, 0, nil
	}

	// Size guard: the claimed payload plus trailer must fit within the file.
	// A crafted size that fails this can no longer reach a make([]byte, lgbSize)
	// or a negative seek offset.
	if !payloadFitsFile(lgbSize, resSize, kind.trailerLen(), total) {
		return bundleNone, 0, 0, fmt.Errorf("corrupt bundle: payload size exceeds file size")
	}
	return kind, lgbSize, resSize, nil
}

// payloadFitsFile reports whether a payload of lgbSize + resSize bytes plus a
// trailerLen-byte trailer fits within a total-byte file. It subtracts step by
// step instead of summing, so the test can't overflow uint64 even on a huge
// (e.g. sparse) file where the individual sizes are valid but their sum wraps.
func payloadFitsFile(lgbSize, resSize uint64, trailerLen, total int64) bool {
	if total < 0 || trailerLen < 0 {
		return false
	}
	avail := uint64(total)
	if lgbSize > avail {
		return false
	}
	avail -= lgbSize
	if resSize > avail {
		return false
	}
	avail -= resSize
	return uint64(trailerLen) <= avail
}

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
		if data, resData := readBundledLGB(path); data != nil {
			return data, resData
		}
	}
	return nil, nil
}

// readBundledLGB extracts the appended payload from the file at path. It
// recognizes both the v2 (resource-carrying) trailer and the legacy trailer;
// res is nil for legacy bundles.
func readBundledLGB(path string) (lgb []byte, res []byte) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	kind, lgbSize, resSize, err := parseBundleTrailer(f)
	if err != nil || kind == bundleNone {
		return nil, nil // not a bundle, or a corrupt one — behave as no payload
	}

	// Payload layout: [lgb][resArc][trailer]. Sizes are validated to fit the
	// file, so the seek offset is a valid negative and make() can't overrun.
	if _, err := f.Seek(-kind.trailerLen()-int64(resSize)-int64(lgbSize), io.SeekEnd); err != nil {
		return nil, nil
	}
	lgb = make([]byte, lgbSize)
	if _, err := io.ReadFull(f, lgb); err != nil {
		return nil, nil
	}
	if resSize > 0 {
		res = make([]byte, resSize)
		if _, err := io.ReadFull(f, res); err != nil {
			return nil, nil
		}
	}
	return lgb, res
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
	resFiles, err := collectResources(resRoots, dstAbs)
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
	binSize, err := getBaseBinarySize(srcBin)
	if err != nil {
		return err
	}
	srcBin.Seek(0, io.SeekStart)
	if _, err := io.CopyN(out, srcBin, binSize); err != nil {
		return err
	}

	// Append LGB data
	if _, err := out.Write(lgbData); err != nil {
		return err
	}

	// Embed the resource archive (if any) between the lgb and the trailer.
	if len(resArc) > 0 {
		if _, err := out.Write(resArc); err != nil {
			return err
		}
		// v2 trailer: [8-byte lgbSize][8-byte resSize][4-byte 'LGB2']
		var tr [20]byte
		binary.LittleEndian.PutUint64(tr[0:8], uint64(len(lgbData)))
		binary.LittleEndian.PutUint64(tr[8:16], uint64(len(resArc)))
		copy(tr[16:], bundleMagicV2[:])
		if _, err := out.Write(tr[:]); err != nil {
			return err
		}
		return nil
	}

	// No resources: legacy trailer [8-byte lgbSize][4-byte 'LGBX'].
	var footer [12]byte
	binary.LittleEndian.PutUint64(footer[:8], uint64(len(lgbData)))
	copy(footer[8:], bundleMagic[:])
	if _, err := out.Write(footer[:]); err != nil {
		return err
	}

	return nil
}

// collectResources returns a map of slash-relative path → file bytes for every
// regular file reachable under the resource roots. It follows symlinks
// everywhere the dev FS provider does (which resolves names with os.Stat) — to
// symlinked roots, symlinked sub-directories, and symlinked files — so a -b
// bundle embeds exactly what dev lookup would find. Symlink cycles are guarded
// by real-path. When the same relative path exists under multiple roots the
// first root wins (matching the provider's precedence). The bundle's own output
// file (excludeAbs) is never embedded, so a dst inside a resource root can't
// embed itself.
func collectResources(roots []string, excludeAbs string) (map[string][]byte, error) {
	files := map[string][]byte{}
	var exclude os.FileInfo
	if excludeAbs != "" {
		// os.Stat (not Lstat) so the comparison is symlink/hardlink robust.
		if fi, err := os.Stat(excludeAbs); err == nil {
			exclude = fi
		}
	}
	for _, root := range roots {
		realRoot, err := filepath.EvalSymlinks(root)
		if err != nil {
			continue // missing or dangling root
		}
		absRoot, err := filepath.Abs(realRoot)
		if err != nil {
			continue
		}
		if info, err := os.Stat(absRoot); err != nil || !info.IsDir() {
			continue
		}
		ancestors := map[string]bool{}
		if err := collectResourceDir(absRoot, absRoot, "", files, exclude, ancestors); err != nil {
			return nil, err
		}
	}
	return files, nil
}

// collectResourceDir recursively collects regular files under dir (an absolute,
// symlink-resolved directory inside rootReal) into files, keyed by relPrefix +
// entry name (slash-separated). Symlinks to files and sub-directories are
// followed, but only when they resolve to a path inside rootReal — a symlink
// escaping the root (e.g. `up -> ..`) is skipped, so a bundle never embeds
// files outside the declared resource tree. ancestors holds the resolved
// directory paths on the current descent path; revisiting one is a cycle and is
// skipped (while distinct names aliasing the same in-root dir are both kept).
func collectResourceDir(dir, rootReal, relPrefix string, files map[string][]byte, exclude os.FileInfo, ancestors map[string]bool) error {
	if ancestors[dir] {
		return nil // symlink cycle
	}
	ancestors[dir] = true
	defer delete(ancestors, dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		p := filepath.Join(dir, e.Name())
		key := e.Name()
		if relPrefix != "" {
			key = relPrefix + "/" + e.Name()
		}
		real, err := filepath.EvalSymlinks(p)
		if err != nil {
			continue // dangling symlink, race — skip
		}
		// Containment: never follow a symlink that escapes the resource root.
		if !rt.WithinRoot(rootReal, real) {
			continue
		}
		info, err := os.Stat(real)
		if err != nil {
			continue
		}
		if info.IsDir() {
			if err := collectResourceDir(real, rootReal, key, files, exclude, ancestors); err != nil {
				return err
			}
			continue
		}
		if !info.Mode().IsRegular() {
			continue // FIFO, device, socket
		}
		if exclude != nil && os.SameFile(info, exclude) {
			continue // never embed the bundle's own output file
		}
		if _, exists := files[key]; exists {
			continue // first root wins
		}
		data, err := os.ReadFile(real)
		if err != nil {
			return err
		}
		files[key] = data
	}
	return nil
}

// getBaseBinarySize returns the size of the lg binary without any appended
// bundle, so re-bundling can strip an existing payload. A corrupt trailer
// surfaces as an error rather than a silently wrong size.
func getBaseBinarySize(f *os.File) (int64, error) {
	kind, lgbSize, resSize, err := parseBundleTrailer(f)
	if err != nil {
		return 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	total := fi.Size()
	if kind == bundleNone {
		return total, nil
	}
	// Sizes are validated to fit the file, so this can't go negative.
	return total - int64(lgbSize) - int64(resSize) - kind.trailerLen(), nil
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
		paths := resolver.PathsFromInputs(sourcePaths, envVal, explicitSet)
		// Transition notice for the dropped implicit ".". Tooling that owns the
		// search path deliberately omits "." and can set
		// LG_SUPPRESS_SOURCE_PATHS_WARNING to silence this; the notice is
		// removed in a future release.
		if !slices.Contains(paths, ".") && os.Getenv("LG_SUPPRESS_SOURCE_PATHS_WARNING") == "" {
			fmt.Fprintln(os.Stderr, `WARNING: the current directory (".") is no `+
				`longer added to the namespace search path automatically when `+
				`-source-paths or LG_SOURCE_PATHS is set; add "." to the list to `+
				`keep searching it. This notice will be removed in a future release `+
				`(set LG_SUPPRESS_SOURCE_PATHS_WARNING=1 to silence).`)
		}
		return paths
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
		if err := buildWasm(context, nsResolver, files[0], wasmOutput); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

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
			fmt.Printf("nREPL server running at tcp://127.0.0.1:%d\n", nreplPort)
		}
		repl(context)
	}

}
