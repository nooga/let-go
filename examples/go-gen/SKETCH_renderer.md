# Sketch: Go-side renderer for gogen

The Clojure prototype in `gogen.lg` produces AST nodes that look like:

```clojure
{:tag :func
 :name "DotProduct"
 :params [{:tag :param :name "a" :type {:tag :type :spec "[]float64"}}
          {:tag :param :name "b" :type {:tag :type :spec "[]float64"}}]
 :ret {:tag :type :spec "float64"}
 :body [{:tag :var-decl :name "sum" :type {:tag :type :spec "float64"} :init {:tag :float-lit :value 0.0}}
        {:tag :assign-op :op "+=" :target {:tag :ident :name "sum"} :value ...}
        ...
        {:tag :return :values [{:tag :ident :name "sum"}]}]}
```

The real renderer replaces `gogen/render` with a Go-side function
that takes this map (as a `vm.Value`), produces `go/ast` nodes, and
calls `go/format.Node` to emit gofmt-canonical source.

## Where it lives

New file: `pkg/rt/gogen.go`. Installed in `init()` (in `pkg/rt/lang.go`)
alongside `installHttpNS`, `installJSONNS`, etc:

```go
func init() {
    // ... existing installs ...
    installGogenNS()
}
```

`installGogenNS` creates a namespace `gogen-go` (or installs `render`
directly into the `gogen` ns after it loads — the question of how to
add Go-side functions to a Clojure-defined namespace is worth checking).

## API surface

Two functions:

```clojure
(gogen/render-go ast-map)
;; -> String of gofmt-canonical Go source.

(gogen/render-file ns-name imports decls)
;; -> Full Go file: package decl, imports, then decls.
```

## Implementation skeleton

```go
package rt

import (
    "bytes"
    "fmt"
    "go/ast"
    "go/format"
    "go/token"

    "github.com/nooga/let-go/pkg/vm"
)

var goFset = token.NewFileSet()

// astFromValue dispatches on the :tag keyword inside a vm.PersistentMap
// and returns the corresponding go/ast node. It returns ast.Node since
// Expr/Stmt/Decl/Spec are all subtypes.
func astFromValue(v vm.Value) (ast.Node, error) {
    m, ok := v.(*vm.PersistentMap)
    if !ok {
        return nil, fmt.Errorf("expected AST map, got %T", v)
    }
    tag := tagOf(m)
    switch tag {
    case "func":      return buildFuncDecl(m)
    case "var-decl":  return buildVarDeclStmt(m)
    case "assign":    return buildAssignStmt(m, token.ASSIGN)
    case "assign-op": return buildAssignOp(m)
    case "return":    return buildReturnStmt(m)
    case "if":        return buildIfStmt(m)
    case "for":       return buildForStmt(m)
    case "expr-stmt": return buildExprStmt(m)
    case "binop":     return buildBinaryExpr(m)
    case "unop":      return buildUnaryExpr(m)
    case "index":     return buildIndexExpr(m)
    case "field":     return buildSelectorExpr(m)
    case "call":      return buildCallExpr(m)
    case "cast":      return buildTypeAssertExpr(m)  // or paren type conversion
    case "ident":     return ast.NewIdent(stringField(m, "name")), nil
    case "int-lit":   return buildBasicLit(m, token.INT)
    case "float-lit": return buildBasicLit(m, token.FLOAT)
    case "string-lit":return buildBasicLit(m, token.STRING)
    case "type":      return buildTypeExpr(m)  // parse :spec as type via parser.ParseExpr
    }
    return nil, fmt.Errorf("unknown AST tag: %s", tag)
}

func buildFuncDecl(m *vm.PersistentMap) (*ast.FuncDecl, error) {
    name := stringField(m, "name")
    paramsV := vectorField(m, "params")
    retV := mapField(m, "ret")
    bodyV := vectorField(m, "body")

    var params []*ast.Field
    for _, pv := range paramsV {
        pm := pv.(*vm.PersistentMap)
        pname := stringField(pm, "name")
        ptype, err := buildTypeExpr(mapField(pm, "type"))
        if err != nil { return nil, err }
        params = append(params, &ast.Field{
            Names: []*ast.Ident{ast.NewIdent(pname)},
            Type:  ptype,
        })
    }

    retExpr, err := buildTypeExpr(retV)
    if err != nil { return nil, err }

    var body []ast.Stmt
    for _, sv := range bodyV {
        n, err := astFromValue(sv)
        if err != nil { return nil, err }
        body = append(body, n.(ast.Stmt))
    }

    return &ast.FuncDecl{
        Name: ast.NewIdent(name),
        Type: &ast.FuncType{
            Params:  &ast.FieldList{List: params},
            Results: &ast.FieldList{List: []*ast.Field{{Type: retExpr}}},
        },
        Body: &ast.BlockStmt{List: body},
    }, nil
}

// buildTypeExpr lets us write Go types as strings ("[]float64",
// "map[string]int") and parse them via go/parser.ParseExpr — much
// simpler than building each type AST by hand.
func buildTypeExpr(m *vm.PersistentMap) (ast.Expr, error) {
    spec := stringField(m, "spec")
    expr, err := parser.ParseExpr(spec)
    if err != nil {
        return nil, fmt.Errorf("type %q: %w", spec, err)
    }
    return expr, nil
}

// ... similar builders for the rest ...

// Render: top-level entry point.
func goRender(v vm.Value) (string, error) {
    node, err := astFromValue(v)
    if err != nil { return "", err }
    var buf bytes.Buffer
    if err := format.Node(&buf, goFset, node); err != nil {
        return "", err
    }
    return buf.String(), nil
}
```

## What this buys

1. **Gofmt-canonical output**: parens disappear from binops, spacing is right, semicolons resolved correctly, line breaks rational.
2. **Real validation**: invalid AST shapes fail at `format.Node` rather than producing nonsense source.
3. **Free upgrade path**: when go/ast grows new node types (generics in 1.18 added type params), we can expose them by adding builders and new :tag values without touching the macro layer.
4. **Optional**: `go/types` validation as a separate `gogen/check` function, for codegen pipelines that want to catch type errors before `go build`.

## Trickiness to handle

- **Positions**: synthesized nodes have `token.NoPos`. `format.Node` mostly handles this but can pack lines tightly. If we hit it, allocate a `*token.File` and assign sequential positions to nodes during build.
- **Statement vs expression coercion**: a bare `call` could be either. We see `:expr-stmt` wrap calls used as statements; the renderer trusts that.
- **Imports**: emitting a full file needs an `*ast.File` with `Decls` including any `import` `*ast.GenDecl`s before user code. Add a `gogen/file` AST tag.
- **Comments**: `go/ast` carries comments separately from nodes; they don't survive a plain `format.Node` round-trip unless you maintain a `CommentMap`. For codegen this is fine — we don't need comments — but if we ever want generated-from-source-with-doc-comments we'll need to handle it.

## What's *not* in this sketch

- Type-parameterized funcs (generics) — straightforward additive change.
- Method receivers — extend `:func` with optional `:recv` field.
- Struct / interface / type declarations — new `:tag` values + builders.
- Switch statements — adds `:case` tags.
- Channel ops, select — adds `:send`, `:recv`, `:select`.

Add as needed. The core machinery (Clojure map → go/ast → format.Node)
doesn't change.

## Effort estimate

The full builder set for the common Go subset (no generics) is about
**500 lines of Go**. Half is mechanical mapping; half is helper utilities
for extracting fields from `vm.PersistentMap` and the type-parsing
shortcut. Probably a weekend.

Adding generics, methods, struct/interface decls, switch, select, and a
File wrapper brings it to ~1000 lines. Another weekend.

The whole gogen library — macro + Go renderer + a few worked examples —
is realistically a one-person two-week project to v1.
