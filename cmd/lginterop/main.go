package main

import (
	_ "embed"
	"flag"
	"fmt"
	"go/build"
	"go/constant"
	"go/importer"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/nooga/let-go/internal/primgen"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// macroLib is the gogen-based codegen library that emits the interop files.
// It is embedded rather than read from the repo so this command works from
// any working directory — `go run github.com/nooga/let-go/cmd/lginterop@v...`
// has no let-go checkout to read it out of. go:embed cannot reach outside
// the package directory, which is why the script lives here and not in
// scripts/.
//
//go:embed lginterop.lg
var macroLib string

// init ensures build.Default.GOROOT matches the user's on-PATH `go` binary,
// so the source importer resolves against the actual Go install.
func init() {
	if out, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		if g := strings.TrimSpace(string(out)); g != "" {
			build.Default.GOROOT = g
		}
	}
}

func main() {
	dir := flag.String("dir", ".", "directory containing deps.edn")
	out := flag.String("out", ".lg-interop", "output directory for generated Go files")
	packagesFlag := flag.String("packages", "", "comma-separated list of packages (overrides deps.edn :gointerop)")
	smartFlag := flag.Bool("smart", false, "generate explicit wrappers with type-specific unboxing/boxing")
	opaqueFlag := flag.Bool("opaque-structs", false, "skip vm.RegisterStruct: struct types stay boxed and dispatch methods reflectively")
	buildTagsFlag := flag.String("build-tags", "", "emit //go:build <constraint> as the first line of each generated file (e.g. '!tinygo')")
	outPkgFlag := flag.String("out-pkg", "", "emit a self-contained package with this name for use OUTSIDE the let-go tree (imports pkg/rt, installs from init); empty emits the in-tree `package rt` form")
	skeletonFlag := flag.Bool("skeleton", false, "generate let-go skeleton files with defn- stubs for hand customization")
	primitivesDir := flag.String("primitives", "", "directory containing //lg:-annotated Go sources (generates zz_primitives_generated.go)")
	primitivesOut := flag.String("primitives-out", "pkg/rt/zz_primitives_generated.go", "output file for generated primitives")
	goPkg := flag.String("go-pkg", "", "Go import path of the scanned sources (used with -primitives)")
	flag.Parse()

	// Handle -primitives mode (separate from the external interop path).
	// Delegates to the runtime-free primgen package — prefer the standalone
	// cmd/lgprimgen binary, which does not import pkg/compiler and so can run
	// mid-migration when the runtime this binary boots cannot.
	if *primitivesDir != "" {
		if err := primgen.Generate(*primitivesDir, *primitivesOut, *goPkg); err != nil {
			fmt.Fprintf(os.Stderr, "lginterop: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := validateOutPkg(*outPkgFlag); err != nil {
		fmt.Fprintf(os.Stderr, "lginterop: %v\n", err)
		os.Exit(1)
	}

	var entries []interopEntry
	if *packagesFlag != "" {
		for pkg := range strings.SplitSeq(*packagesFlag, ",") {
			entries = append(entries, interopEntry{
				pkg: strings.TrimSpace(pkg), smart: *smartFlag, opaque: *opaqueFlag,
				buildTags: *buildTagsFlag, outPkg: *outPkgFlag,
			})
		}
	} else {
		var err error
		entries, err = gointeropFromDepsEdn(*dir, *smartFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lginterop: %v\n", err)
			os.Exit(1)
		}
		for i := range entries {
			entries[i].opaque = *opaqueFlag
			entries[i].buildTags = *buildTagsFlag
			entries[i].outPkg = *outPkgFlag
		}
	}

	if len(entries) == 0 {
		fmt.Println("lginterop: no packages to generate")
		return
	}

	if err := os.MkdirAll(*out, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "lginterop: mkdir %s: %v\n", *out, err)
		os.Exit(1)
	}

	// Output filenames derive from the alias; two entries sharing one would
	// silently clobber each other's interop_<alias>.go. Refuse up front.
	seenAlias := map[string]string{}
	for _, ent := range entries {
		alias := ent.alias
		if alias == "" {
			alias = defaultAlias(ent.pkg)
		}
		if prev, dup := seenAlias[alias]; dup {
			fmt.Fprintf(os.Stderr, "lginterop: alias %q used by both %s and %s — outputs would collide; set a distinct alias\n", alias, prev, ent.pkg)
			os.Exit(1)
		}
		seenAlias[alias] = ent.pkg
	}

	okCount := 0
	for _, ent := range entries {
		if err := generatePackage(ent, *out, *skeletonFlag); err != nil {
			fmt.Fprintf(os.Stderr, "lginterop: %s: %v\n", ent.pkg, err)
			continue
		}
		okCount++
	}

	fmt.Printf("lginterop: generated %d/%d package(s) in %s\n", okCount, len(entries), *out)
	// Per-package failures only log and continue, so without this a run that
	// generated nothing still exited 0 — a caller driving this from a build
	// pipeline would see success and a missing file.
	if okCount < len(entries) {
		os.Exit(1)
	}
}

// goKeywords are rejected as package names: the emitter would happily write
// `package package`, and the tool would report success for output that does
// not parse.
var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

// validateOutPkg rejects anything that would not compile as a package clause.
// An empty value is the in-tree default and always valid.
func validateOutPkg(name string) error {
	if name == "" {
		return nil
	}
	// rt is the in-tree package, which the empty default already selects.
	// Accepting it here would give two spellings for two DIFFERENT outputs
	// (`-out-pkg rt` would import rt into itself), so refuse it outright.
	if name == "rt" {
		return fmt.Errorf("-out-pkg rt is ambiguous: omit the flag for the in-tree `package rt` output")
	}
	if name == "_" {
		return fmt.Errorf("-out-pkg %q is not a valid package name", name)
	}
	if goKeywords[name] {
		return fmt.Errorf("-out-pkg %q is a Go keyword, not a valid package name", name)
	}
	for i, r := range name {
		ok := r == '_' ||
			(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(i > 0 && r >= '0' && r <= '9')
		if !ok {
			return fmt.Errorf("-out-pkg %q is not a valid Go identifier "+
				"(letters, digits and _ only; may not start with a digit)", name)
		}
	}
	return nil
}

// --- deps.edn parsing -----------------------------------------------------

func gointeropFromDepsEdn(dir string, globalSmart bool) ([]interopEntry, error) {
	depsPath := path.Join(dir, "deps.edn")
	data, err := os.ReadFile(depsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("deps.edn not found in %s", dir)
		}
		return nil, err
	}
	val, err := compiler.ReadString(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse deps.edn: %w", err)
	}
	m, ok := val.(*vm.PersistentMap)
	if !ok {
		return nil, fmt.Errorf("deps.edn root is not a map")
	}

	var out []interopEntry
	if m.Contains(vm.Keyword("gointerop")) {
		v := m.ValueAt(vm.Keyword("gointerop"))
		vec, ok := v.(vm.ArrayVector)
		if !ok {
			return nil, fmt.Errorf(":gointerop is not a vector")
		}
		for _, item := range vec {
			ent := parseInteropItem(item)
			ent.smart = globalSmart
			if ent.pkg != "" {
				out = append(out, ent)
			}
		}
	}

	if m.Contains(vm.Keyword("gointerop-wrappers")) {
		v := m.ValueAt(vm.Keyword("gointerop-wrappers"))
		vec, ok := v.(vm.ArrayVector)
		if !ok {
			return nil, fmt.Errorf(":gointerop-wrappers is not a vector")
		}
		for _, item := range vec {
			ent := parseInteropItem(item)
			ent.smart = true
			if ent.pkg != "" {
				found := false
				for i := range out {
					if out[i].pkg == ent.pkg {
						out[i].smart = true
						if ent.alias != "" {
							out[i].alias = ent.alias
						}
						found = true
						break
					}
				}
				if !found {
					out = append(out, ent)
				}
			}
		}
	}

	return out, nil
}

func parseInteropItem(item vm.Value) interopEntry {
	switch it := item.(type) {
	case vm.String:
		if it != "" {
			return interopEntry{pkg: string(it)}
		}
	case *vm.PersistentMap:
		for s := it.Seq(); s != nil; s = s.Next() {
			entry, ok := s.First().(vm.MapEntry)
			if !ok {
				continue
			}
			if pkgStr, ok := entry.Key.(vm.String); ok && pkgStr != "" {
				alias := ""
				if aliasVal, ok := entry.Value.(vm.String); ok {
					alias = string(aliasVal)
				}
				return interopEntry{pkg: string(pkgStr), alias: alias}
			}
		}
	case vm.ArrayVector:
		if len(it) >= 2 {
			if pkgStr, ok := it[0].(vm.String); ok && pkgStr != "" {
				alias := ""
				if aliasVal, ok := it[1].(vm.String); ok {
					alias = string(aliasVal)
				}
				return interopEntry{pkg: string(pkgStr), alias: alias}
			}
		}
	}
	return interopEntry{}
}

type interopEntry struct {
	pkg       string
	alias     string
	smart     bool
	opaque    bool   // skip vm.RegisterStruct: structs stay Boxed, methods dispatch reflectively
	buildTags string // emitted as //go:build <constraint> when non-empty
	outPkg    string // non-empty: emit a self-contained package for out-of-tree use
}

type export struct {
	name string
	obj  types.Object
}

func defaultAlias(pkg string) string {
	if i := strings.LastIndex(pkg, "/"); i >= 0 {
		return pkg[i+1:]
	}
	return pkg
}

// --- package generation ---------------------------------------------------

func generatePackage(ent interopEntry, outDir string, skeleton bool) error {
	pkgName := ent.pkg
	alias := ent.alias
	if alias == "" {
		alias = defaultAlias(pkgName)
	}

	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, "source", nil)
	pkg, err := imp.Import(pkgName)
	if err != nil {
		// The go/types source importer resolves imports against the module
		// context of the CURRENT working directory. Scanning a third-party
		// package therefore only works from inside a module that requires it.
		return fmt.Errorf("import %s: %w\n"+
			"hint: the scanned package must be resolvable from the current "+
			"directory's module — run lginterop inside a module that requires "+
			"it (`go get %s` first)", pkgName, err, pkgName)
	}

	var exports []export
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)
		if !obj.Exported() {
			continue
		}
		if _, isBuiltin := obj.(*types.Builtin); isBuiltin {
			continue
		}
		if c, ok := obj.(*types.Const); ok && !constBoxable(c) {
			continue
		}
		exports = append(exports, export{name: name, obj: obj})
	}

	if len(exports) == 0 {
		fmt.Printf("lginterop: %s — no eligible exports, skipping\n", pkgName)
		return nil
	}

	sort.Slice(exports, func(i, j int) bool {
		return exports[i].name < exports[j].name
	})

	fileName := "interop_" + goPackageToFileName(alias) + ".go"
	outPath := filepath.Join(outDir, fileName)

	// Build the let-go script that drives gogen codegen and evaluate it in
	// THIS process. Running it in-process (rather than building an lg binary
	// and shelling out) is what makes the command self-contained: gogen is
	// registered by pkg/rt, which any binary importing the runtime already
	// links, so there is no repo root to find and no emitter/binary drift to
	// guard against.
	script := buildGenScript(pkgName, alias, exports, outPath, ent.smart, ent.opaque, ent.buildTags, ent.outPkg)
	if err := keepScript(script); err != nil {
		return err
	}
	if err := evalGenScript(script); err != nil {
		return err
	}

	mode := "direct"
	if ent.smart {
		mode = "smart"
	}
	fmt.Printf("lginterop: %s (as %s, %s) → %s (%d exports)\n", pkgName, alias, mode, outPath, len(exports))

	if skeleton {
		skelPath := filepath.Join(outDir, alias+"_skeleton.lg")
		skel := buildSkeleton(alias, exports, ent.smart)
		if err := os.WriteFile(skelPath, []byte(skel), 0644); err != nil {
			return fmt.Errorf("write skeleton %s: %w", skelPath, err)
		}
		fmt.Printf("lginterop: skeleton → %s\n", skelPath)
	}

	return nil
}

// --- Lisp script generation -----------------------------------------------

func buildGenScript(pkgName, alias string, exports []export, outPath string, smart, opaque bool, buildTags, outPkg string) string {
	var b strings.Builder
	b.WriteString(macroLib)
	b.WriteString("\n")
	b.WriteString("(def exports ")
	b.WriteString(serializeExports(exports))
	b.WriteString(")\n")
	fmt.Fprintf(&b, "(lginterop/generate %s %s exports %s %s %s %s %s)\n",
		strconv.Quote(pkgName), strconv.Quote(alias), strconv.Quote(outPath),
		strconv.FormatBool(smart), strconv.FormatBool(opaque), strconv.Quote(buildTags),
		strconv.Quote(outPkg))
	return b.String()
}

// evalGenScript runs the codegen script through the runtime this binary
// already links. The script writes its output with `spit`, so nothing comes
// back through the return value; only the error matters.
//
// CompileMultiple, not pkg/api's Run: the script is a sequence of top-level
// forms (an (ns ...), the defns, the driver call), and Run compiles only the
// first one. This mirrors lg.go's own runForm.
func evalGenScript(script string) error {
	ns := rt.NS("user")
	if ns == nil {
		return fmt.Errorf("codegen: namespace \"user\" not found")
	}
	ctx := compiler.NewCompiler(vm.NewConsts(), ns)
	if _, _, err := ctx.CompileMultiple(strings.NewReader(script)); err != nil {
		return fmt.Errorf("codegen: %w", err)
	}
	return nil
}

// keepScript honors LGINTEROP_KEEP_SCRIPT by dumping the generated script to a
// temp file for debugging. The script is no longer written to disk as part of
// normal operation, so this is the only way to inspect it.
func keepScript(script string) error {
	if os.Getenv("LGINTEROP_KEEP_SCRIPT") == "" {
		return nil
	}
	f, err := os.CreateTemp("", "lginterop-*.lg")
	if err != nil {
		return fmt.Errorf("keep script: %w", err)
	}
	defer f.Close()
	if _, err := f.WriteString(script); err != nil {
		return fmt.Errorf("keep script: %w", err)
	}
	fmt.Printf("lginterop: keeping temp script %s\n", f.Name())
	return nil
}

// serializeExports emits a compact positional vector for each export:
//
//	[:func  "Name" [type-params] [params] [results] variadic?]
//	[:type  "Name" [type-params] struct? [methods] [fields]]
//	[:const "Name"]
//	[:var   "Name"]
//
// Methods:  ["Name" [params] [results] variadic?]
// Fields:   ["Name" "type" embedded?]
// Type-params: [{:name "T" :constraint "any"} ...]
func serializeExports(exports []export) string {
	var b strings.Builder
	b.WriteString("[")
	for i, ex := range exports {
		if i > 0 {
			b.WriteString("\n ")
		}
		switch obj := ex.obj.(type) {
		case *types.Func:
			sig := obj.Type().(*types.Signature)
			b.WriteString("[:func ")
			b.WriteString(strconv.Quote(ex.name))
			b.WriteString(" ")
			b.WriteString(serializeTypeParams(sig.TypeParams()))
			b.WriteString(" ")
			b.WriteString(serializeTypeSlice(sig.Params()))
			b.WriteString(" ")
			b.WriteString(serializeTypeSlice(sig.Results()))
			if sig.Variadic() {
				b.WriteString(" :variadic")
			}
			b.WriteString("]")
		case *types.TypeName:
			b.WriteString("[:type ")
			b.WriteString(strconv.Quote(ex.name))
			b.WriteString(" ")
			if named, ok := obj.Type().(*types.Named); ok {
				b.WriteString(serializeTypeParams(named.TypeParams()))
				b.WriteString(" ")
				if isStructType(obj.Type()) {
					b.WriteString(":struct ")
					if strct, ok := named.Underlying().(*types.Struct); ok {
						b.WriteString(serializeFields(strct))
					} else {
						b.WriteString("[]")
					}
				} else {
					b.WriteString("nil []")
				}
				b.WriteString(" ")
				b.WriteString(serializeMethods(named))
			} else {
				b.WriteString("nil nil [] []")
			}
			b.WriteString("]")
		case *types.Const:
			b.WriteString("[:const ")
			b.WriteString(strconv.Quote(ex.name))
			b.WriteString("]")
		case *types.Var:
			b.WriteString("[:var ")
			b.WriteString(strconv.Quote(ex.name))
			b.WriteString("]")
		default:
			b.WriteString("[:unknown ")
			b.WriteString(strconv.Quote(ex.name))
			b.WriteString("]")
		}
	}
	b.WriteString("]")
	return b.String()
}

func serializeTypeParams(tplist *types.TypeParamList) string {
	if tplist == nil || tplist.Len() == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteString("[")
	for i := 0; i < tplist.Len(); i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		tp := tplist.At(i)
		b.WriteString("{:name ")
		b.WriteString(strconv.Quote(tp.String()))
		b.WriteString(" :constraint ")
		b.WriteString(strconv.Quote(types.TypeString(tp.Constraint(), nil)))
		b.WriteString("}")
	}
	b.WriteString("]")
	return b.String()
}

func serializeTypeSlice(list *types.Tuple) string {
	if list == nil {
		return "[]"
	}
	var b strings.Builder
	b.WriteString("[")
	for i := 0; i < list.Len(); i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(strconv.Quote(types.TypeString(list.At(i).Type(), nil)))
	}
	b.WriteString("]")
	return b.String()
}

func serializeMethods(named *types.Named) string {
	var b strings.Builder
	b.WriteString("[")
	first := true
	for m := range named.Methods() {
		if !m.Exported() {
			continue
		}
		if !first {
			b.WriteString(" ")
		}
		first = false
		sig := m.Type().(*types.Signature)
		b.WriteString("[")
		b.WriteString(strconv.Quote(m.Name()))
		b.WriteString(" ")
		b.WriteString(serializeTypeParams(sig.TypeParams()))
		b.WriteString(" ")
		b.WriteString(serializeTypeSlice(sig.Params()))
		b.WriteString(" ")
		b.WriteString(serializeTypeSlice(sig.Results()))
		if sig.Variadic() {
			b.WriteString(" :variadic")
		}
		b.WriteString("]")
	}
	b.WriteString("]")
	return b.String()
}

func serializeFields(strct *types.Struct) string {
	var b strings.Builder
	b.WriteString("[")
	for i := 0; i < strct.NumFields(); i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		f := strct.Field(i)
		b.WriteString("[")
		b.WriteString(strconv.Quote(f.Name()))
		b.WriteString(" ")
		b.WriteString(strconv.Quote(types.TypeString(f.Type(), nil)))
		if f.Embedded() {
			b.WriteString(" :embedded")
		}
		b.WriteString("]")
	}
	b.WriteString("]")
	return b.String()
}

func constBoxable(c *types.Const) bool {
	// Untyped integer constants that overflow int64 can't be passed as `any`
	// (they default to int and fail at compile time). Skip them.
	if c.Val().Kind() == constant.Int {
		if _, ok := constant.Int64Val(c.Val()); !ok {
			return false
		}
	}
	return true
}

func goPackageToFileName(pkg string) string {
	return strings.NewReplacer("/", "_", ".", "_", "-", "_").Replace(pkg)
}

func isStructType(t types.Type) bool {
	if named, ok := t.(*types.Named); ok {
		_, isStruct := named.Underlying().(*types.Struct)
		return isStruct
	}
	return false
}

// --- skeleton generation --------------------------------------------------

func buildSkeleton(alias string, exports []export, smart bool) string {
	b := &strings.Builder{}
	fmt.Fprintf(b, ";; Generated by lginterop for package %q.\n", alias)
	b.WriteString(";; Hand-customize the defn- stubs below as needed.\n")
	fmt.Fprintf(b, "(ns %s)\n\n", alias)

	for _, ex := range exports {
		qname := alias + "/" + ex.name
		switch obj := ex.obj.(type) {
		case *types.Func:
			sig := obj.Type().(*types.Signature)
			params := sig.Params()
			arity := params.Len()
			variadic := sig.Variadic()

			argNames := make([]string, arity)
			for i := range arity {
				argNames[i] = fmt.Sprintf("a%d", i)
			}

			if smart {
				fmt.Fprintf(b, "(defn- %s\n", primgen.KebabCase(ex.name))
				fmt.Fprintf(b, "  \"Wrapper for %s. Customize as needed.\"\n", qname)
				if variadic {
					fmt.Fprintf(b, "  [& args]\n")
					fmt.Fprintf(b, "  (apply %s args))\n\n", qname)
				} else if arity == 0 {
					fmt.Fprintf(b, "  []\n")
					fmt.Fprintf(b, "  (%s))\n\n", qname)
				} else {
					fmt.Fprintf(b, "  [%s]\n", strings.Join(argNames, " "))
					fmt.Fprintf(b, "  (%s %s))\n\n", qname, strings.Join(argNames, " "))
				}
			} else {
				fmt.Fprintf(b, "(defn- %s\n", primgen.KebabCase(ex.name))
				fmt.Fprintf(b, "  \"Wrapper for %s. Customize as needed.\"\n", qname)
				if variadic {
					fmt.Fprintf(b, "  [& args]\n")
					fmt.Fprintf(b, "  (apply %s args))\n\n", qname)
				} else if arity == 0 {
					fmt.Fprintf(b, "  []\n")
					fmt.Fprintf(b, "  (%s))\n\n", qname)
				} else {
					fmt.Fprintf(b, "  [%s]\n", strings.Join(argNames, " "))
					fmt.Fprintf(b, "  (%s %s))\n\n", qname, strings.Join(argNames, " "))
				}
			}
		case *types.TypeName:
			if isStructType(obj.Type()) {
				fmt.Fprintf(b, ";; Struct type registered: %s\n", qname)
				fmt.Fprintf(b, ";; Use (make-record %s {...}) after registration.\n\n", qname)
			}
		case *types.Const:
			fmt.Fprintf(b, ";; Constant: %s\n", qname)
			fmt.Fprintf(b, ";; (def %s %s)\n\n", primgen.KebabCase(ex.name), qname)
		case *types.Var:
			fmt.Fprintf(b, ";; Variable: %s\n", qname)
			fmt.Fprintf(b, ";; (def %s %s)\n\n", primgen.KebabCase(ex.name), qname)
		}
	}

	return b.String()
}

// Avoid "declared and not used" for runtime import.
var _ = runtime.GOOS
