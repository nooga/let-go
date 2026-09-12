// Command go-callables measures Go callables for the code-quality report.
//
//	go-callables <path>...    # files, or directories walked recursively
//
// It writes EDN to stdout: one record per top-level function or method
// declaration, plus one per dead nested binding.
//
//	[{:path "pkg/vm/map.go" :name "(*Map).Get" :kind :method
//	  :line 12 :end 40 :sloc 21 :cc 7} ...]
//
// It MEASURES and judges nothing. Mass, the complexity-above-ten split, the
// erosion ratio, weights, thresholds and reporting all live in .lg; see
// scripts/quality/. The division exists because cyclomatic complexity is a
// pure function of AST node kinds, which Go's own parser already gives us,
// while every decision about what those numbers MEAN belongs with the rest
// of the metric code.
//
// It is a command rather than a runtime native because a Go parser bridged
// into pkg/rt would be runtime scope creep for a tooling need, and because
// standalone AST tools are the established pattern here (cmd/hoist-natives,
// cmd/lgprimgen). quality.lg already shells out for git and the coverage
// profile, so one more subprocess is not a new dependency direction.
//
// Files that fail to parse are reported on stderr and skipped; the exit
// status is non-zero only if NOTHING could be measured, so a single
// unparseable file cannot silently empty the corpus.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ednString quotes a string for EDN: the same escapes EDN and Go agree on.
func ednString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString("\\\"")
		case '\\':
			b.WriteString("\\\\")
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		case '\r':
			b.WriteString("\\r")
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// goFiles expands each argument: a file is itself, a directory is every .go
// file under it. Generated files and the lowered-core tree are skipped, the
// same exclusions the .lg side applies.
func goFiles(args []string) ([]string, error) {
	var out []string
	for _, a := range args {
		info, err := os.Stat(a)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			out = append(out, a)
			continue
		}
		err = filepath.Walk(a, func(p string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				if name := fi.Name(); name == ".git" || name == ".jj" ||
					name == ".workspaces" || name == "core_go_lowered" {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_generated.go") {
				out = append(out, p)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(out)
	return out, nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: go-callables <path>...")
		os.Exit(2)
	}
	files, err := goFiles(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "go-callables:", err)
		os.Exit(2)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	fmt.Fprint(w, "[")

	measured, failed := 0, 0
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go-callables: %s: %v\n", f, err)
			failed++
			continue
		}
		cs, err := goCallables(string(src))
		if err != nil {
			fmt.Fprintf(os.Stderr, "go-callables: %s: %v\n", f, err)
			failed++
			continue
		}
		measured++
		for _, c := range cs {
			fmt.Fprintf(w, "\n {:path %s :name %s :kind :%s :line %d :end %d :sloc %d :cc %d}",
				ednString(f), ednString(c.name), c.kind, c.line, c.end, c.sloc, c.cc)
		}
	}
	fmt.Fprintln(w, "]")

	if measured == 0 && failed > 0 {
		fmt.Fprintln(os.Stderr, "go-callables: nothing could be measured")
		w.Flush()
		os.Exit(1)
	}
}
