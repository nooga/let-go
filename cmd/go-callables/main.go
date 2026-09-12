// Command go-callables measures Go callables for the code-quality report.
//
//	go-callables <path>...    # files, or directories walked recursively
//
// It writes one EDN map to stdout, keyed by concern so that adding a concern
// later is additive and breaks no consumer:
//
//	{:version 1
//	 :files     [{:path "pkg/vm/map.go" :lines 120 :sloc 98 :blank 14}]
//	 :functions [{:path "pkg/vm/map.go" :name "(*Map).Get" :kind :method
//	              :line 12 :end 40 :sloc 21 :cc 7}]}
//
// This is THE GO ANALYSIS TOOL, not "the erosion input provider". The
// boundary between it and the .lg side is the language being analysed, not
// the metric: anything that needs to understand Go should come from here,
// where go/ast and go/scanner are exact, rather than from a generic
// tokenizer that is close enough. Two keys are deliberately absent and
// expected: `:tokens`, a normalized per-file token stream from go/scanner so
// duplication can winnow real Go tokens, and `:comments`, comment blocks
// with their line ranges for scripts/lint.lg.
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

	type measurement struct {
		path   string
		counts fileCounts
		calls  []callable
	}
	var ms []measurement
	failed := 0
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go-callables: %s: %v\n", f, err)
			failed++
			continue
		}
		cs, counts, err := goAnalyze(string(src))
		if err != nil {
			fmt.Fprintf(os.Stderr, "go-callables: %s: %v\n", f, err)
			failed++
			continue
		}
		ms = append(ms, measurement{path: f, counts: counts, calls: cs})
	}

	fmt.Fprintln(w, "{:version 1")
	fmt.Fprint(w, " :files [")
	for _, m := range ms {
		fmt.Fprintf(w, "\n  {:path %s :lines %d :sloc %d :blank %d}",
			ednString(m.path), m.counts.lines, m.counts.code, m.counts.blank)
	}
	fmt.Fprintln(w, "]")
	fmt.Fprint(w, " :functions [")
	for _, m := range ms {
		for _, c := range m.calls {
			fmt.Fprintf(w, "\n  {:path %s :name %s :kind :%s :line %d :end %d :sloc %d :cc %d}",
				ednString(m.path), ednString(c.name), c.kind, c.line, c.end, c.sloc, c.cc)
		}
	}
	fmt.Fprintln(w, "]}")

	measured := len(ms)

	if measured == 0 && failed > 0 {
		fmt.Fprintln(os.Stderr, "go-callables: nothing could be measured")
		w.Flush()
		os.Exit(1)
	}
}
