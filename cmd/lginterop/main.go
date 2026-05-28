package main

import (
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
	"unicode"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/vm"
)

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
	skeletonFlag := flag.Bool("skeleton", false, "generate let-go skeleton files with defn- stubs for hand customization")
	flag.Parse()

	var entries []interopEntry
	if *packagesFlag != "" {
		for pkg := range strings.SplitSeq(*packagesFlag, ",") {
			entries = append(entries, interopEntry{pkg: strings.TrimSpace(pkg), smart: *smartFlag})
		}
	} else {
		var err error
		entries, err = gointeropFromDepsEdn(*dir, *smartFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lginterop: %v\n", err)
			os.Exit(1)
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

	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lginterop: %v\n", err)
		os.Exit(1)
	}

	lgBin, err := ensureLgBinary(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lginterop: %v\n", err)
		os.Exit(1)
	}

	okCount := 0
	for _, ent := range entries {
		if err := generatePackage(repoRoot, lgBin, ent, *out, *skeletonFlag); err != nil {
			fmt.Fprintf(os.Stderr, "lginterop: %s: %v\n", ent.pkg, err)
			continue
		}
		okCount++
	}

	fmt.Printf("lginterop: generated %d/%d package(s) in %s\n", okCount, len(entries), *out)
}

// --- repo root & lg binary discovery --------------------------------------

func findRepoRoot() (string, error) {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("cannot find repo root (no go.mod in ancestor directories)")
}

func ensureLgBinary(repoRoot string) (string, error) {
	lgPath := filepath.Join(repoRoot, "lg")
	if _, err := os.Stat(lgPath); err == nil {
		return lgPath, nil
	}
	cmd := exec.Command("go", "build", "-o", "lg", ".")
	cmd.Dir = repoRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build lg binary: %w\n%s", err, output)
	}
	return lgPath, nil
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
	pkg   string
	alias string
	smart bool
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

func generatePackage(repoRoot, lgBin string, ent interopEntry, outDir string, skeleton bool) error {
	pkgName := ent.pkg
	alias := ent.alias
	if alias == "" {
		alias = defaultAlias(pkgName)
	}

	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, "source", nil)
	pkg, err := imp.Import(pkgName)
	if err != nil {
		return fmt.Errorf("import: %w", err)
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

	fileName := goPackageToFileName(pkgName) + ".go"
	outPath := filepath.Join(outDir, fileName)

	// Write the Lisp script that drives gogen codegen.
	scriptPath, err := writeGenScript(repoRoot, pkgName, alias, exports, outPath, ent.smart)
	if err != nil {
		return fmt.Errorf("write script: %w", err)
	}
	// DEBUG: keep temp script for inspection
	// defer os.Remove(scriptPath)

	cmd := exec.Command(lgBin, "-source-paths", filepath.Join(repoRoot, "scripts"), scriptPath)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lg failed: %w\noutput:\n%s", err, output)
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

func writeGenScript(repoRoot, pkgName, alias string, exports []export, outPath string, smart bool) (string, error) {
	macroPath := filepath.Join(repoRoot, "scripts", "lginterop.lg")
	macroLib, err := os.ReadFile(macroPath)
	if err != nil {
		return "", fmt.Errorf("read macro library: %w", err)
	}

	var b strings.Builder
	b.Write(macroLib)
	b.WriteString("\n")
	b.WriteString("(def exports ")
	b.WriteString(serializeExports(exports))
	b.WriteString(")\n")
	fmt.Fprintf(&b, "(lginterop/generate %s %s exports %s %s)\n",
		strconv.Quote(pkgName), strconv.Quote(alias), strconv.Quote(outPath), strconv.FormatBool(smart))

	tmpFile, err := os.CreateTemp("", "lginterop-*.lg")
	if err != nil {
		return "", err
	}
	if _, err := tmpFile.WriteString(b.String()); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}
	return tmpFile.Name(), nil
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
				fmt.Fprintf(b, "(defn- %s\n", kebabCase(ex.name))
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
				fmt.Fprintf(b, "(defn- %s\n", kebabCase(ex.name))
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
			fmt.Fprintf(b, ";; (def %s %s)\n\n", kebabCase(ex.name), qname)
		case *types.Var:
			fmt.Fprintf(b, ";; Variable: %s\n", qname)
			fmt.Fprintf(b, ";; (def %s %s)\n\n", kebabCase(ex.name), qname)
		}
	}

	return b.String()
}

func kebabCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 {
			prev := rune(s[i-1])
			if unicode.IsUpper(r) {
				if unicode.IsLower(prev) {
					b.WriteByte('-')
				} else if i+1 < len(s) && unicode.IsLower(rune(s[i+1])) {
					b.WriteByte('-')
				}
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// Avoid "declared and not used" for runtime import.
var _ = runtime.GOOS
