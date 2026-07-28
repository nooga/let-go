// Command hoist-natives is a codemod that lifts the inline anonymous closures
// registered as clojure.core primitives in pkg/rt/lang.go into named,
// //lg:native-annotated top-level functions.
//
// Motivation (three wins at once):
//  1. Stack traces read `rt.corePlus` instead of `installCore.func42`.
//  2. Registration becomes declarative — the `ns.Def("+", plus)` call is
//     dropped and lginterop emits the registration from the //lg:native +
//     //lg:name annotations.
//  3. It is the #411 "Fn-adapter" done as a codemod: hand-built Fn values
//     become annotated decls lginterop already consumes, with no runtime
//     adapter.
//
// It only hoists closures that are SAFE to lift: no genuine enclosing-scope
// capture (package-level refs and params are fine; an installCore local is
// not), and a discoverable lg-name (an `ns.Def("<name>", <var>)` site). Anything
// else is reported and left untouched — the codemod never guesses.
//
// Default is a dry run: it prints a classification report and the proposed
// hoisted function for each safe site. `-apply` performs the rewrite.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nooga/let-go/internal/primgen"
)

func main() {
	file := flag.String("file", "pkg/rt/lang.go", "Go file to codemod")
	pkgDir := flag.String("pkg", "pkg/rt", "package dir (for resolving package-level names)")
	only := flag.String("only", "", "regex; limit to lg-names matching it")
	helpers := flag.String("helpers", "", "comma-separated installCore-local helper closures to lift to package level (instead of hoisting primitives)")
	apply := flag.Bool("apply", false, "apply the rewrite (default: dry run + report)")
	flag.Parse()

	var onlyRe *regexp.Regexp
	if *only != "" {
		onlyRe = regexp.MustCompile(*only)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *file, nil, parser.ParseComments)
	if err != nil {
		die("parse %s: %v", *file, err)
	}

	pkgNames, err := packageLevelNames(*pkgDir)
	if err != nil {
		die("scan package: %v", err)
	}
	for name := range importNames(f) {
		pkgNames[name] = true
	}

	if *helpers != "" {
		liftHelpers(*file, fset, f, splitSet(*helpers), pkgNames, *apply)
		return
	}

	defNames := collectNsDefNames(f)             // fn-var name -> lg name (last wins)
	allNames := collectNsDefAllNames(f)          // fn-var name -> every lg name
	lookupNames := lookupNamesInPackage(*pkgDir) // lg-names fetched via .Lookup("…")
	// lg-names a pre-existing //lg:native decl already owns (native_prims.go etc.),
	// excluding the file we're hoisting into. A candidate colliding with one is a
	// redundant second registration; hoisting it would emit a duplicate.
	existingNative, err := primgen.ScanNativeNames(*pkgDir, filepath.Base(*file))
	if err != nil {
		die("scan existing //lg:native names: %v", err)
	}
	// lg-names a stdlib .lg source redefines — shadowed at load and not covered
	// by the eager core reapply, so hoisting+guarding them would deviate.
	lgRedefined := lgRedefinedNames(filepath.Join(*pkgDir, "core"))

	sites := findNativeFnSites(f)

	var safe []*site
	skipCap := map[string][]string{} // var -> captured names
	var skipNoName []string
	var skipMulti []string
	for _, s := range sites {
		if onlyRe != nil {
			lg := defNames[s.varName]
			if lg == "" || !onlyRe.MatchString(lg) {
				continue
			}
		}
		lg, ok := defNames[s.varName]
		if !ok {
			skipNoName = append(skipNoName, s.varName)
			continue
		}
		if len(allNames[s.varName]) > 1 {
			// Registered under multiple names — hoisting to one //lg:name would
			// silently drop the others from clojure.core. Leave it in place.
			skipMulti = append(skipMulti, fmt.Sprintf("%s (%s)", s.varName, strings.Join(allNames[s.varName], "/")))
			continue
		}
		if lookupNames[lg] {
			// Fetched by lg-name via `.Lookup("<lg>")` somewhere in the package
			// (possibly another file) at init time, before the generated
			// registrar runs — hoisting would leave that lookup nil.
			skipMulti = append(skipMulti, fmt.Sprintf("%s (Lookup %q)", s.varName, lg))
			continue
		}
		if existingNative[lg] {
			// Name already owned by a pre-existing //lg:native decl — hoisting
			// this redundant registration would emit a duplicate in the registrar.
			skipMulti = append(skipMulti, fmt.Sprintf("%s (dup of //lg:native %q)", s.varName, lg))
			continue
		}
		if lgRedefined[lg] {
			// Shadowed by a stdlib .lg (defn %q) at load; the eager core reapply
			// skips it, so a guarded native here deviates permanently.
			skipMulti = append(skipMulti, fmt.Sprintf("%s (redefined in .lg %q)", s.varName, lg))
			continue
		}
		s.lgName = lg
		captured := freeVars(s.lit, pkgNames)
		if len(captured) > 0 {
			skipCap[s.varName] = captured
			continue
		}
		safe = append(safe, s)
	}
	if len(skipMulti) > 0 {
		sort.Strings(skipMulti)
		fmt.Printf("-- skipped: registered under multiple lg-names (would drop names) --\n  %s\n\n", strings.Join(skipMulti, ", "))
	}

	// Demote any safe var still referenced by code that survives the rewrite.
	// The rewrite removes only the var's own `X, err := …Wrap(…)` assignment and
	// its single `ns.Def("name", X)` registration; a var referenced anywhere else
	// — a second registration (`lgCore.Def("gt", gt)`), a helper, or a skipped
	// site's closure body — would become an undefined identifier. freeVars flags
	// captures WITHIN a hoisted closure; this flags uses OF the var by others.
	safe, skipRef := demoteSurvivingRefs(f, safe)

	report(fset, safe, skipCap, skipNoName, defNames)
	if len(skipRef) > 0 {
		sort.Strings(skipRef)
		fmt.Printf("-- skipped: still referenced by surviving code (2nd registration / helper / captured by a skip) --\n  %s\n\n", strings.Join(skipRef, ", "))
	}

	if *apply {
		if err := rewrite(*file, fset, f, safe, defNames); err != nil {
			die("apply: %v", err)
		}
		fmt.Printf("\napplied: hoisted %d functions into %s\n", len(safe), *file)
	} else {
		fmt.Printf("\ndry run — pass -apply to write. proposed hoists for the %d safe sites:\n\n", len(safe))
		for _, s := range safe {
			fmt.Println(genFunc(fset, s))
		}
	}
}

// --- site discovery -------------------------------------------------------

type site struct {
	varName string       // the local the closure was assigned to (plus, mul, …)
	lgName  string       // the clojure.core name (from ns.Def)
	ctx     bool         // true for NewCtxNativeFn (first param is ec)
	lit     *ast.FuncLit // the closure
	assign  *ast.AssignStmt
}

// findNativeFnSites finds `X, _ := vm.NativeFnType.Wrap(func…)` /
// `.WrapNoErr(func…)` / `vm.NewCtxNativeFn(func…)` assignments.
func findNativeFnSites(f *ast.File) []*site {
	var out []*site
	ast.Inspect(f, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 2 || len(as.Rhs) != 1 {
			return true
		}
		lhs, ok := as.Lhs[0].(*ast.Ident)
		if !ok || lhs.Name == "_" {
			return true
		}
		call, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}
		lit, ok := call.Args[0].(*ast.FuncLit)
		if !ok {
			return true
		}
		boxer, ctx := classifyBoxer(call.Fun)
		if boxer == "" {
			return true
		}
		out = append(out, &site{varName: lhs.Name, ctx: ctx, lit: lit, assign: as})
		return true
	})
	return out
}

// demoteSurvivingRefs returns the subset of safe whose local var is referenced
// ONLY by its own assignment target and the single `ns.Def("name", X)` call the
// rewrite cuts, plus the names it demoted. A var referenced by any surviving
// identifier (a second registration, a helper, or a skipped closure body) is
// not hoistable — removing its local would leave that reference undefined.
func demoteSurvivingRefs(f *ast.File, safe []*site) (kept []*site, demoted []string) {
	byVar := map[string]*site{}
	for _, s := range safe {
		byVar[s.varName] = s
	}
	// Identifiers the rewrite removes: each safe assign's target and the arg of
	// the single ns.Def(name, X) it cuts. Every OTHER occurrence survives.
	accounted := map[*ast.Ident]bool{}
	for _, s := range safe {
		if id, ok := s.assign.Lhs[0].(*ast.Ident); ok {
			accounted[id] = true
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 2 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Def" {
			return true
		}
		if id, ok := sel.X.(*ast.Ident); !ok || id.Name != "ns" {
			return true
		}
		if arg, ok := call.Args[1].(*ast.Ident); ok && byVar[arg.Name] != nil {
			accounted[arg] = true
		}
		return true
	})
	survives := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok || byVar[id.Name] == nil || accounted[id] {
			return true
		}
		survives[id.Name] = true
		return true
	})

	// A primitive is also un-hoistable if surviving code fetches it BY LG-NAME
	// via `ns.Lookup("name")` — for post-registration mutation (`.SetMacro()`)
	// or to alias it under another name. installLangNS runs before the generated
	// registrar, so such a lookup would hit an unregistered var. These references
	// are string literals, invisible to the identifier scan above.
	byLg := map[string]*site{}
	for _, s := range safe {
		byLg[s.lgName] = s
	}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Lookup" {
			return true
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		if name, err := strconv.Unquote(lit.Value); err == nil && byLg[name] != nil {
			survives[byLg[name].varName] = true
		}
		return true
	})
	for _, s := range safe {
		if survives[s.varName] {
			demoted = append(demoted, s.varName)
			continue
		}
		kept = append(kept, s)
	}
	return kept, demoted
}

// classifyBoxer returns ("Wrap"|"WrapNoErr"|"NewCtxNativeFn", isCtx) or ("",_).
func classifyBoxer(fun ast.Expr) (string, bool) {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	switch sel.Sel.Name {
	case "Wrap", "WrapNoErr":
		// vm.NativeFnType.Wrap
		if inner, ok := sel.X.(*ast.SelectorExpr); ok && inner.Sel.Name == "NativeFnType" {
			return sel.Sel.Name, false
		}
	case "NewCtxNativeFn":
		if id, ok := sel.X.(*ast.Ident); ok && id.Name == "vm" {
			return sel.Sel.Name, true
		}
	}
	return "", false
}

// collectNsDefAllNames maps a registered fn-var to EVERY clojure name it is
// Def'd under. A var registered under several names (`ns.Def("int", intf)`,
// `ns.Def("byte", intf)`, `ns.Def("short", intf)`) can't be hoisted to a single
// //lg:native decl without dropping names, so the caller skips multi-name vars.
func collectNsDefAllNames(f *ast.File) map[string][]string {
	out := map[string][]string{}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 2 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Def" {
			return true
		}
		if nsID, ok := sel.X.(*ast.Ident); !ok || nsID.Name != "ns" {
			return true
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		valID, ok := call.Args[1].(*ast.Ident)
		if !ok {
			return true
		}
		if name, err := strconv.Unquote(lit.Value); err == nil {
			out[valID.Name] = append(out[valID.Name], name)
		}
		return true
	})
	return out
}

// collectNsDefNames maps a registered fn-var to its clojure name from
// `ns.Def("name", fnvar)` call sites.
func collectNsDefNames(f *ast.File) map[string]string {
	out := map[string]string{}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 2 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Def" {
			return true
		}
		nsID, ok := sel.X.(*ast.Ident)
		if !ok || nsID.Name != "ns" {
			return true
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		valID, ok := call.Args[1].(*ast.Ident)
		if !ok {
			return true
		}
		name, err := strconv.Unquote(lit.Value)
		if err == nil {
			out[valID.Name] = name
		}
		return true
	})
	return out
}

// --- free-variable analysis ----------------------------------------------

// freeVars returns identifiers used in the closure body that are neither
// declared within it nor package-level nor universe builtins — i.e. genuine
// enclosing-scope captures. Over-reports on shadowing rather than under, so a
// flagged site is skipped for human review, never mis-hoisted.
func freeVars(lit *ast.FuncLit, pkgNames map[string]bool) []string {
	declared := map[string]bool{}
	for _, fld := range lit.Type.Params.List {
		for _, nm := range fld.Names {
			declared[nm.Name] = true
		}
	}
	// collect names bound anywhere inside the body
	ast.Inspect(lit.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			if x.Tok == token.DEFINE {
				for _, l := range x.Lhs {
					if id, ok := l.(*ast.Ident); ok {
						declared[id.Name] = true
					}
				}
			}
		case *ast.RangeStmt:
			for _, e := range []ast.Expr{x.Key, x.Value} {
				if id, ok := e.(*ast.Ident); ok {
					declared[id.Name] = true
				}
			}
		case *ast.DeclStmt:
			if gd, ok := x.Decl.(*ast.GenDecl); ok {
				for _, spec := range gd.Specs {
					switch sp := spec.(type) {
					case *ast.ValueSpec:
						for _, nm := range sp.Names {
							declared[nm.Name] = true
						}
					case *ast.TypeSpec:
						declared[sp.Name.Name] = true
					}
				}
			}
		case *ast.FuncLit:
			for _, fld := range x.Type.Params.List {
				for _, nm := range fld.Names {
					declared[nm.Name] = true
				}
			}
		}
		return true
	})

	// Idents that are NOT value references: a selector's field/method name
	// (`x.Type` — `Type` is not a capture) and a struct composite's field key.
	skip := map[*ast.Ident]bool{}
	ast.Inspect(lit.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			skip[x.Sel] = true
		case *ast.CompositeLit:
			// keys of a struct composite are field names. (Map/slice keys ARE
			// values, so only skip when the type is a struct/named type — being
			// conservative and NOT skipping when unsure would just over-flag.)
			if isNamedType(x.Type) {
				for _, elt := range x.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						if id, ok := kv.Key.(*ast.Ident); ok {
							skip[id] = true
						}
					}
				}
			}
		}
		return true
	})

	freeSet := map[string]bool{}
	ast.Inspect(lit.Body, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && !skip[id] {
			considerIdent(id, declared, pkgNames, freeSet)
		}
		return true
	})
	out := make([]string, 0, len(freeSet))
	for nm := range freeSet {
		out = append(out, nm)
	}
	sort.Strings(out)
	return out
}

// isNamedType reports whether a composite-literal type is a named/struct type
// (T{…} or pkg.T{…}), whose keys are field names — as opposed to []T{…} /
// map[K]V{…} whose keys are values.
func isNamedType(e ast.Expr) bool {
	switch e.(type) {
	case *ast.Ident, *ast.SelectorExpr:
		return true
	default:
		return false
	}
}

func considerIdent(id *ast.Ident, declared, pkg, free map[string]bool) {
	nm := id.Name
	if nm == "_" || declared[nm] || pkg[nm] || universe[nm] {
		return
	}
	// treat obvious non-values out: capitalized single tokens that are types are
	// usually package-level; if unknown, we flag (conservative).
	free[nm] = true
}

var universe = strSet(
	"append", "cap", "clear", "close", "complex", "copy", "delete", "imag",
	"len", "make", "max", "min", "new", "panic", "print", "println", "real",
	"recover", "nil", "true", "false", "iota",
	"bool", "byte", "rune", "string", "error", "int", "int8", "int16", "int32",
	"int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
	"float32", "float64", "complex64", "complex128", "any", "comparable",
)

// --- package-level name collection ---------------------------------------

// lookupNamesInPackage collects every string literal passed to a `.Lookup("…")`
// call across all non-test .go files in dir. A primitive named by any of these
// lgRedefinedNames collects every name defined by a (def…)/(defn…)/(defmacro…)
// form in the stdlib .lg sources under dir/core. A primitive whose lg-name is
// redefined in core.lg is shadowed by that bootstrap definition at load; the
// eager core-namespace load does NOT reapply the native root over it (see
// coreload.go, which skips NameCoreNS), so a hoisted+guarded native there would
// deviate permanently (native-prims-intact? => false). On main these primitives
// were hand-registered but unguarded, so the .lg definition silently won — the
// behavior-preserving choice is to leave them un-hoisted.
func lgRedefinedNames(coreDir string) map[string]bool {
	names := map[string]bool{}
	re := regexp.MustCompile(`\(def(?:n|macro)?-?\s+([^\s()]+)`)
	entries, err := os.ReadDir(coreDir)
	if err != nil {
		return names
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".lg") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(coreDir, e.Name()))
		if err != nil {
			continue
		}
		for _, m := range re.FindAllStringSubmatch(string(src), -1) {
			names[m[1]] = true
		}
	}
	return names
}

// is fetched by lg-name at init time — possibly from ANOTHER file (host_core_fns
// aliases `class` to `ns.Lookup("type")`) before the generated registrar runs —
// so it must not be hoisted. Best-effort union across files.
func lookupNamesInPackage(dir string) map[string]bool {
	names := map[string]bool{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return names
	}
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, 0)
		if err != nil {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 1 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Lookup" {
				return true
			}
			if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				if name, err := strconv.Unquote(lit.Value); err == nil {
					names[name] = true
				}
			}
			return true
		})
	}
	return names
}

func packageLevelNames(dir string) (map[string]bool, error) {
	names := map[string]bool{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, 0)
		if err != nil {
			continue // tolerate build-tag / partial parse issues; names union is best-effort
		}
		for _, d := range f.Decls {
			switch decl := d.(type) {
			case *ast.FuncDecl:
				if decl.Recv == nil {
					names[decl.Name.Name] = true
				}
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch sp := spec.(type) {
					case *ast.ValueSpec:
						for _, nm := range sp.Names {
							names[nm.Name] = true
						}
					case *ast.TypeSpec:
						names[sp.Name.Name] = true
					}
				}
			}
		}
	}
	return names, nil
}

func importNames(f *ast.File) map[string]bool {
	out := map[string]bool{}
	for _, imp := range f.Imports {
		if imp.Name != nil {
			out[imp.Name.Name] = true
			continue
		}
		p, _ := strconv.Unquote(imp.Path.Value)
		base := p
		if i := strings.LastIndex(p, "/"); i >= 0 {
			base = p[i+1:]
		}
		out[base] = true
	}
	return out
}

// --- generation -----------------------------------------------------------

// genFunc renders the hoisted, annotated function for a safe site.
func genFunc(fset *token.FileSet, s *site) string {
	name := hoistName(s.varName)
	params := convertSignature(fset, s.lit.Type.Params)
	results := renderResults(fset, s.lit.Type.Results)
	body := nodeString(fset, s.lit.Body)

	var b strings.Builder
	b.WriteString("//lg:native\n")
	b.WriteString("//lg:name " + s.lgName + "\n")
	fmt.Fprintf(&b, "func %s(%s) %s %s", name, params, results, body)
	// gofmt the fragment
	src := "package rt\n" + b.String() + "\n"
	if out, err := format.Source([]byte(src)); err == nil {
		return strings.TrimPrefix(string(out), "package rt\n\n")
	}
	return b.String()
}

// convertSignature renders params, turning a trailing `vs []vm.Value` into the
// variadic `vs ...vm.Value` shape lginterop's native model consumes. Typed
// params (ec *vm.ExecContext, coll vm.Value, i int) pass through unchanged.
func convertSignature(fset *token.FileSet, params *ast.FieldList) string {
	var parts []string
	for i, fld := range params.List {
		typ := fld.Type
		if i == len(params.List)-1 {
			if arr, ok := typ.(*ast.ArrayType); ok && arr.Len == nil && nodeString(fset, arr.Elt) == "vm.Value" {
				name := "args"
				if len(fld.Names) > 0 {
					name = fld.Names[0].Name
				}
				parts = append(parts, name+" ...vm.Value")
				continue
			}
		}
		parts = append(parts, renderField(fset, fld))
	}
	return strings.Join(parts, ", ")
}

// renderField renders `[name[, name] ]type` for a param/result field.
// format.Node handles the type (an ast.Expr) but not *ast.Field itself.
func renderField(fset *token.FileSet, fld *ast.Field) string {
	typ := nodeString(fset, fld.Type)
	if len(fld.Names) == 0 {
		return typ
	}
	names := make([]string, len(fld.Names))
	for i, nm := range fld.Names {
		names[i] = nm.Name
	}
	return strings.Join(names, ", ") + " " + typ
}

// renderResults renders the result list, parenthesised when needed.
func renderResults(fset *token.FileSet, results *ast.FieldList) string {
	if results == nil || len(results.List) == 0 {
		return ""
	}
	parts := make([]string, len(results.List))
	named := false
	for i, fld := range results.List {
		if len(fld.Names) > 0 {
			named = true
		}
		parts[i] = renderField(fset, fld)
	}
	joined := strings.Join(parts, ", ")
	if len(results.List) > 1 || named {
		return "(" + joined + ") "
	}
	return joined + " "
}

func hoistName(varName string) string {
	// plus -> CorePlus, excludeInCurrentNs -> CoreExcludeInCurrentNs.
	// EXPORTED (capital C): the gogen_ir lowered tree lives in sibling packages
	// that direct-call these primitives as `rt.Core…`, so an unexported name
	// would be invisible there (undefined: rt.corePlus).
	if varName == "" {
		return "CoreFn"
	}
	return "Core" + strings.ToUpper(varName[:1]) + varName[1:]
}

// --- rewrite (apply) ------------------------------------------------------

// rewrite deletes each safe site's `X, _ := …` assignment and its
// `ns.Def("name", X)` statement, then appends the hoisted funcs, at the text
// level (position-based) so formatting elsewhere is untouched. gofmt at the end.
func rewrite(path string, fset *token.FileSet, f *ast.File, safe []*site, defNames map[string]string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	type span struct{ start, end int }
	var cuts []span
	byVar := map[string]*site{}
	for _, s := range safe {
		byVar[s.varName] = s
		cuts = append(cuts, span{fset.Position(s.assign.Pos()).Offset, lineEnd(src, fset.Position(s.assign.End()).Offset)})
	}
	// find the ns.Def statements to cut
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 2 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Def" {
			return true
		}
		if id, ok := sel.X.(*ast.Ident); !ok || id.Name != "ns" {
			return true
		}
		valID, ok := call.Args[1].(*ast.Ident)
		if !ok || byVar[valID.Name] == nil {
			return true
		}
		cuts = append(cuts, span{lineStart(src, fset.Position(call.Pos()).Offset), lineEnd(src, fset.Position(call.End()).Offset)})
		return true
	})

	// Cut `if err != nil { … }` guards left dead by the removal. An `err` local
	// is typically threaded through every Wrap assignment in a block — each with
	// its own `if err != nil { panic(err) }`, plus a trailing consolidated
	// `if err != nil { panic("… failed") }`. Once every assignment that DECLARES
	// that err (via `:=`) is cut and no surviving `:=` redeclares it, the name is
	// undefined and all its guards are dead. Work per block: a name declared only
	// by cut assigns (and by no surviving `:=`) is orphaned, so cut every
	// `if <name> != nil {…}` guard in that block. Scoped this way, a guard whose
	// err still has a live declaration is never touched.
	safeAssigns := map[*ast.AssignStmt]bool{}
	for _, s := range safe {
		safeAssigns[s.assign] = true
	}
	declaredNames := func(as *ast.AssignStmt) []string {
		if as.Tok != token.DEFINE {
			return nil
		}
		var ns []string
		for _, l := range as.Lhs {
			if id, ok := l.(*ast.Ident); ok && id.Name != "_" {
				ns = append(ns, id.Name)
			}
		}
		return ns
	}
	ast.Inspect(f, func(n ast.Node) bool {
		blk, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		cutDecl, surviveDecl := map[string]bool{}, map[string]bool{}
		for _, st := range blk.List {
			as, ok := st.(*ast.AssignStmt)
			if !ok {
				continue
			}
			names := declaredNames(as)
			target := surviveDecl
			if safeAssigns[as] {
				target = cutDecl
			}
			for _, nm := range names {
				target[nm] = true
			}
		}
		for _, st := range blk.List {
			name := errGuardName(st)
			if name != "" && cutDecl[name] && !surviveDecl[name] {
				cuts = append(cuts, span{lineStart(src, fset.Position(st.Pos()).Offset), lineEnd(src, fset.Position(st.End()).Offset)})
			}
		}
		return true
	})

	sort.Slice(cuts, func(i, j int) bool { return cuts[i].start > cuts[j].start })
	out := append([]byte(nil), src...)
	for _, c := range cuts {
		out = append(out[:c.start], out[c.end:]...)
	}

	var appended bytes.Buffer
	appended.WriteString("\n// --- hoisted native primitives (see cmd/hoist-natives) ---\n\n")
	for _, s := range safe {
		appended.WriteString(genFunc(fset, s))
		appended.WriteString("\n\n")
	}
	out = append(out, appended.Bytes()...)

	formatted, err := format.Source(out)
	if err != nil {
		// leave unformatted so the author can inspect; still write
		formatted = out
	}
	return os.WriteFile(path, formatted, 0644)
}

// errGuardName returns the identifier name X for a statement of the exact shape
// `if X != nil { … }` (no init/else), or "" otherwise. The caller decides
// whether X is orphaned before cutting, so this only classifies the shape.
func errGuardName(s ast.Stmt) string {
	ifs, ok := s.(*ast.IfStmt)
	if !ok || ifs.Init != nil || ifs.Else != nil {
		return ""
	}
	bin, ok := ifs.Cond.(*ast.BinaryExpr)
	if !ok || bin.Op != token.NEQ {
		return ""
	}
	x, ok := bin.X.(*ast.Ident)
	if !ok {
		return ""
	}
	if y, ok := bin.Y.(*ast.Ident); !ok || y.Name != "nil" {
		return ""
	}
	return x.Name
}

func lineStart(src []byte, off int) int {
	for off > 0 && src[off-1] != '\n' {
		off--
	}
	return off
}
func lineEnd(src []byte, off int) int {
	for off < len(src) && src[off] != '\n' {
		off++
	}
	if off < len(src) {
		off++ // include newline
	}
	return off
}

// --- reporting ------------------------------------------------------------

func report(fset *token.FileSet, safe []*site, skipCap map[string][]string, skipNoName []string, defNames map[string]string) {
	fmt.Printf("hoist-natives report\n====================\n")
	fmt.Printf("safe to hoist:        %d\n", len(safe))
	fmt.Printf("skipped (captures):   %d\n", len(skipCap))
	fmt.Printf("skipped (no lg name): %d\n\n", len(skipNoName))

	if len(skipCap) > 0 {
		fmt.Println("-- skipped: genuine enclosing-scope captures (need ec-threading / manual) --")
		keys := sortedKeys(skipCap)
		for _, k := range keys {
			fmt.Printf("  %-24s captures: %s   (lg: %q)\n", k, strings.Join(skipCap[k], ", "), defNames[k])
		}
		fmt.Println()
	}
	if len(skipNoName) > 0 {
		sort.Strings(skipNoName)
		fmt.Println("-- skipped: no ns.Def(\"name\", var) found (registered differently / internal) --")
		fmt.Printf("  %s\n\n", strings.Join(skipNoName, ", "))
	}
}

// --- helpers --------------------------------------------------------------

func nodeString(fset *token.FileSet, n ast.Node) string {
	if n == nil {
		return ""
	}
	var b bytes.Buffer
	if err := format.Node(&b, fset, n); err != nil {
		return ""
	}
	return b.String()
}

func strSet(xs ...string) map[string]bool {
	m := make(map[string]bool, len(xs))
	for _, x := range xs {
		m[x] = true
	}
	return m
}

func sortedKeys[V any](m map[string]V) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// --- helper lift ----------------------------------------------------------

type helperSite struct {
	name   string
	lit    *ast.FuncLit
	assign *ast.AssignStmt
}

// liftHelpers lifts named installCore-local helper closures (X := func…{…}) to
// package-level funcs. Helpers may call each other, so the requested set is
// treated as already-package-level for the capture check.
func liftHelpers(path string, fset *token.FileSet, f *ast.File, set, pkgNames map[string]bool, apply bool) {
	for name := range set {
		pkgNames[name] = true
	}
	sites := findHelperDecls(f, set)
	found := map[string]bool{}
	var safe []*helperSite
	skipCap := map[string][]string{}
	for _, h := range sites {
		found[h.name] = true
		if captured := freeVars(h.lit, pkgNames); len(captured) > 0 {
			skipCap[h.name] = captured
			continue
		}
		safe = append(safe, h)
	}

	fmt.Printf("lift-helpers report\n===================\n")
	fmt.Printf("requested: %d, found: %d, liftable: %d, capture-blocked: %d\n\n",
		len(set), len(found), len(safe), len(skipCap))
	for name := range set {
		if !found[name] {
			fmt.Printf("  NOT FOUND as `%s := func…`\n", name)
		}
	}
	for _, k := range sortedKeys(skipCap) {
		fmt.Printf("  capture-blocked %-20s captures: %s\n", k, strings.Join(skipCap[k], ", "))
	}
	fmt.Println()

	if apply {
		if err := rewriteHelpers(path, fset, safe); err != nil {
			die("lift-helpers apply: %v", err)
		}
		fmt.Printf("applied: lifted %d helpers into %s\n", len(safe), path)
		return
	}
	fmt.Printf("dry run — proposed package-level funcs for the %d liftable helpers:\n\n", len(safe))
	for _, h := range safe {
		fmt.Println(genHelper(fset, h))
	}
}

func findHelperDecls(f *ast.File, set map[string]bool) []*helperSite {
	var out []*helperSite
	ast.Inspect(f, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			return true
		}
		id, ok := as.Lhs[0].(*ast.Ident)
		if !ok || !set[id.Name] {
			return true
		}
		lit, ok := as.Rhs[0].(*ast.FuncLit)
		if !ok {
			return true
		}
		out = append(out, &helperSite{name: id.Name, lit: lit, assign: as})
		return true
	})
	return out
}

// genHelper renders a lifted helper as a package-level func — no annotations, no
// signature conversion (helpers keep their exact interface).
func genHelper(fset *token.FileSet, h *helperSite) string {
	var params []string
	for _, fld := range h.lit.Type.Params.List {
		params = append(params, renderField(fset, fld))
	}
	results := renderResults(fset, h.lit.Type.Results)
	body := nodeString(fset, h.lit.Body)
	src := fmt.Sprintf("package rt\nfunc %s(%s) %s%s\n", h.name, strings.Join(params, ", "), results, body)
	if out, err := format.Source([]byte(src)); err == nil {
		return strings.TrimPrefix(string(out), "package rt\n\n")
	}
	return src
}

func rewriteHelpers(path string, fset *token.FileSet, safe []*helperSite) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	type span struct{ start, end int }
	var cuts []span
	for _, h := range safe {
		cuts = append(cuts, span{
			lineStart(src, fset.Position(h.assign.Pos()).Offset),
			lineEnd(src, fset.Position(h.assign.End()).Offset),
		})
	}
	sort.Slice(cuts, func(i, j int) bool { return cuts[i].start > cuts[j].start })
	out := append([]byte(nil), src...)
	for _, c := range cuts {
		out = append(out[:c.start], out[c.end:]...)
	}
	var appended bytes.Buffer
	appended.WriteString("\n// --- lifted native helpers (see cmd/hoist-natives -helpers) ---\n\n")
	for _, h := range safe {
		appended.WriteString(genHelper(fset, h))
		appended.WriteString("\n\n")
	}
	out = append(out, appended.Bytes()...)
	formatted, err := format.Source(out)
	if err != nil {
		formatted = out
	}
	return os.WriteFile(path, formatted, 0644)
}

func splitSet(csv string) map[string]bool {
	m := map[string]bool{}
	for _, p := range strings.Split(csv, ",") {
		if p = strings.TrimSpace(p); p != "" {
			m[p] = true
		}
	}
	return m
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "hoist-natives: "+format+"\n", args...)
	os.Exit(1)
}
