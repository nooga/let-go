package bytecode

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestFlagsLayout locks the positional flag-bit convention so two branches
// cannot each assign the same bit without a test failure. The failure mode it
// prevents is a clean git merge of two Flag* declarations that both use the
// same shift (the #501 / #624 shape).
func TestFlagsLayout(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse package: %v", err)
	}
	pkg := pkgs["bytecode"]
	if pkg == nil {
		t.Fatal("bytecode package not found")
	}

	var (
		flagsEndBlock  *ast.GenDecl
		flagsEndIsLast bool
		outsideFlags   []string
		explicitValued []string
	)

	for _, f := range pkg.Files {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}

			hasFlagsEnd := false
			for _, spec := range gd.Specs {
				vs := spec.(*ast.ValueSpec)
				for _, name := range vs.Names {
					if name.Name == "flagsEnd" {
						hasFlagsEnd = true
					}
				}
			}

			if !hasFlagsEnd {
				for _, spec := range gd.Specs {
					vs := spec.(*ast.ValueSpec)
					for _, name := range vs.Names {
						if strings.HasPrefix(name.Name, "Flag") {
							outsideFlags = append(outsideFlags, name.Name)
						}
					}
				}
				continue
			}

			flagsEndBlock = gd
			for i, spec := range gd.Specs {
				vs := spec.(*ast.ValueSpec)
				for _, name := range vs.Names {
					if name.Name == "flagsEnd" {
						flagsEndIsLast = i == len(gd.Specs)-1
					}
				}
				// Bits are positional: only the first entry may carry an
				// explicit value (`= 1 << iota`). Everything after inherits.
				if i == 0 {
					continue
				}
				if len(vs.Values) > 0 {
					for _, name := range vs.Names {
						explicitValued = append(explicitValued, name.Name)
					}
				}
			}
		}
	}

	if flagsEndBlock == nil {
		t.Fatal("flagsEnd const block not found in package bytecode")
	}
	if !flagsEndIsLast {
		last := flagsEndBlock.Specs[len(flagsEndBlock.Specs)-1].(*ast.ValueSpec).Names[0].Name
		t.Errorf("flagsEnd must be the last entry in its const block so appended flags take the next bit; found %s last", last)
	}
	for _, name := range outsideFlags {
		t.Errorf("flag %s is declared outside the flagsEnd block; a flag declared elsewhere can silently take a bit that is already in use", name)
	}
	for _, name := range explicitValued {
		t.Errorf("flag %s carries an explicit value; bits are positional and must be left to iota", name)
	}
}

func TestKnownFlagsMatchesNewestVersionSet(t *testing.T) {
	// Appending a flag grows knownFlags; this forces a coordinated edit to
	// v3Flags (and a conscious choice about older version sets) before the
	// suite goes green again.
	if knownFlags != v3Flags {
		t.Fatalf("knownFlags=%#04x v3Flags=%#04x; the newest admitted set must cover every flag in the iota block", knownFlags, v3Flags)
	}
}
