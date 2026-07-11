package main

import (
	"bytes"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"strings"
)

// primSpec describes one primitive function marked with //lg:native.
type primSpec struct {
	GoPkg      string   // Go package import path, e.g. "github.com/nooga/let-go/pkg/rt/builtins"
	GoIdent    string   // exported Go identifier, e.g. "UpperCase"
	Ns         string   // target let-go namespace, e.g. "clojure.string"
	LgName     string   // let-go name, e.g. "upper-case" (defaults to kebab-cased GoIdent)
	Arity      int      // fixed arity; -1 for variadic
	Variadic   bool     // true → final ParamSpec is the rest slice element type
	ParamSpecs []string // Go type strings as understood by lowering, e.g. []string{"string"}
	ResultSpec string   // Go type of single result, e.g. "string"
	NeedsError bool     // result tuple is (ResultSpec, error)
	NeedsEC    bool     // true → first param is *vm.ExecContext (and excluded from ParamSpecs)
	Private    bool     // true if marked with //lg:private
	Package    string   // Go package name (e.g. "builtins") from the source file
}

// hasBuildConstraint reports whether the file carries a //go:build (or legacy
// // +build) line. The generated registrar compiles unconditionally on every
// platform, so a primitive in a constrained file would emit references that
// break the targets the constraint excludes. Skipping ANY constrained file —
// rather than evaluating the constraint against the scanning host — keeps the
// committed output deterministic across GOOS/GOARCH.
func hasBuildConstraint(src []byte) bool {
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(src, '\n'); i >= 0 {
			line, src = src[:i], src[i+1:]
		} else {
			src = nil
		}
		trimmed := string(bytes.TrimSpace(line))
		if strings.HasPrefix(trimmed, "package ") {
			break
		}
		if constraint.IsGoBuild(trimmed) || constraint.IsPlusBuild(trimmed) {
			return true
		}
	}
	return false
}

// scanSource parses the Go source and extracts all //lg:native-annotated
// functions. A parse failure is an error: treating it as "no primitives"
// would silently replace valid generated output with an empty registrar.
func scanSource(path string, src []byte) ([]primSpec, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var specs []primSpec

	// Capture the package name from the file
	packageName := ""
	if file.Name != nil {
		packageName = file.Name.Name
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil { // skip methods, only process top-level funcs
			continue
		}

		if fn.Doc == nil {
			continue
		}

		ns, name, private, native := readDirectives(fn.Doc)
		if !native {
			continue
		}

		// Parse the signature
		sig := fn.Type
		goIdent := fn.Name.Name

		// Determine lg name (default: kebab-cased GoIdent)
		lgName := name
		if lgName == "" {
			lgName = kebabCase(goIdent)
		}

		// Parse parameters
		// Count individual parameter names, not fields (e.g., "start, end int" is 2 params, not 1 field)
		var paramSpecs []string
		needsEC := false
		arity := 0
		for _, field := range sig.Params.List {
			arity += len(field.Names)
		}

		if arity > 0 {
			// Check if first param is *vm.ExecContext
			firstField := sig.Params.List[0]
			firstType := typeString(firstField.Type)
			if firstType == "*vm.ExecContext" && len(firstField.Names) > 0 {
				needsEC = true
				arity--
				// Skip the EC param when building ParamSpecs
				// The EC param is the first name in the first field, but we need to handle remaining names
				// and all subsequent fields
				if len(firstField.Names) > 1 {
					// Multiple names in first field, only first is EC
					for _, name := range firstField.Names[1:] {
						paramSpecs = append(paramSpecs, typeString(firstField.Type))
						_ = name // use name to avoid unused
					}
				}
				// Add all params from remaining fields
				for _, field := range sig.Params.List[1:] {
					for range field.Names {
						paramSpecs = append(paramSpecs, typeString(field.Type))
					}
				}
			} else {
				// All params are regular params: collect from all fields
				for _, field := range sig.Params.List {
					for range field.Names {
						paramSpecs = append(paramSpecs, typeString(field.Type))
					}
				}
			}
		}

		// Check for variadic
		variadic := false
		if arity > 0 && len(sig.Params.List) > 0 {
			lastField := sig.Params.List[len(sig.Params.List)-1]
			if _, ok := lastField.Type.(*ast.Ellipsis); ok {
				variadic = true
				arity = -1
				// A variadic fn accepts 0+ of the rest element; only the
				// fixed params stay in ParamSpecs.
				if len(paramSpecs) > 0 {
					paramSpecs = paramSpecs[:len(paramSpecs)-1]
				}
			}
		}

		// Parse results
		needsError := false
		resultSpec := ""
		if sig.Results != nil && sig.Results.NumFields() > 0 {
			numResults := sig.Results.NumFields()
			if numResults >= 2 {
				// Check if last result is error
				lastResult := sig.Results.List[numResults-1]
				if len(lastResult.Names) == 0 { // unnamed param
					lastType := typeString(lastResult.Type)
					if lastType == "error" {
						needsError = true
						// Get the actual result type (first result)
						resultType := sig.Results.List[0]
						resultSpec = typeString(resultType.Type)
					}
				}
			} else if numResults == 1 {
				// Single result (no error)
				resultType := sig.Results.List[0]
				resultSpec = typeString(resultType.Type)
			}
		}

		specs = append(specs, primSpec{
			GoPkg:      "", // Will be filled in by the caller if needed, or derived from package context
			GoIdent:    goIdent,
			Ns:         ns,
			LgName:     lgName,
			Arity:      arity,
			Variadic:   variadic,
			ParamSpecs: paramSpecs,
			ResultSpec: resultSpec,
			NeedsError: needsError,
			NeedsEC:    needsEC,
			Private:    private,
			Package:    packageName, // Capture the package name
		})
	}

	return specs, nil
}

// readDirectives extracts //lg: directives from a comment group.
func readDirectives(doc *ast.CommentGroup) (ns, name string, private, native bool) {
	ns = "clojure.core"
	for _, c := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		switch {
		case line == "lg:native":
			native = true
		case line == "lg:private":
			private = true
		case strings.HasPrefix(line, "lg:ns "):
			ns = strings.TrimSpace(line[len("lg:ns "):])
		case strings.HasPrefix(line, "lg:name "):
			name = strings.TrimSpace(line[len("lg:name "):])
		}
	}
	return
}

// typeString converts an AST expression representing a type to a string.
func typeString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		// e.g. vm.String → "vm.String"
		if x, ok := e.X.(*ast.Ident); ok {
			return x.Name + "." + e.Sel.Name
		}
	case *ast.StarExpr:
		// e.g. *vm.ExecContext → "*vm.ExecContext"
		return "*" + typeString(e.X)
	case *ast.Ellipsis:
		// e.g. ...string → "...string"
		return "..." + typeString(e.Elt)
	case *ast.ArrayType:
		if e.Len == nil {
			// slice type []T
			return "[]" + typeString(e.Elt)
		}
		// array type [N]T — not used for params but included for completeness
		return "[...]" + typeString(e.Elt)
	}
	return "unknown"
}
