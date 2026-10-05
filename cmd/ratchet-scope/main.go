// Command ratchet-scope decides whether a change set can affect what the
// performance ratchet measures.
//
// The rule: run the ratchet when the change set can change what the ratchet
// measures — the IR and the VM plus everything they transitively depend on.
//
// Scoping matters because the gate compares a fresh capture against a
// COMMITTED bar. An unscoped gate fails every push from a machine whose bar has
// drifted, whatever the push contains, so unrelated work can only proceed by
// rebaselining first. Scope keeps that coupling to changes that could actually
// have moved the numbers.
//
// "Transitively" is computed, not guessed. The benchmarked packages are the
// roots; `go list -deps -test` expands them to the packages actually compiled
// into those benchmark binaries, under each build-tag configuration the
// ratchet captures, and reports each package's real file list. A file the
// toolchain says is compiled into or embedded by one of those packages is in
// scope; anything else is not. That is why cmd/ratchet-scope,
// cmd/bench-ratchet, docs/ and .github/ do not trigger a run: nothing in the
// measured binaries compiles them.
//
// Matching the file list rather than the directory tree matters in both
// directions. pkg/rt EMBEDS every core/**/*.lg file and core_compiled.lgb, so
// those are in scope automatically rather than through a hand-kept list that
// can drift from what is really embedded. Conversely test/*.lg sits in the
// `test` package's directory while belonging to no build, and the suite
// benchmark reads only the vendored corpus, so it is correctly out of scope.
//
// Only inputs that belong to no package are matched by path: the
// generated-artifact manifests, and the module graph and toolchain files that
// change generated code beneath every package.
//
// The bias is deliberately toward running. A false "affected" costs one
// benchmark run; a false "not affected" lets a regression through the gate
// that exists to catch it. So an unreadable dependency graph, an unknown path,
// or any error is treated as affected.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// benchmarkRoots are the packages the ratchet actually benchmarks. They are
// the roots of the dependency expansion, and they mirror the job list in
// cmd/bench-ratchet (anchorPackage, suitePackage, irCompilePackage,
// initPackage). If a benchmark moves to a new package, add it here — a root
// that is missing silently narrows the gate.
var benchmarkRoots = []string{
	"github.com/nooga/let-go/pkg/vm",
	"github.com/nooga/let-go/pkg/ir",
	"github.com/nooga/let-go/pkg/compiler",
	"github.com/nooga/let-go/test",
}

// tagConfigs are the build-tag configurations the ratchet captures under. The
// gogen_ir build compiles a different set of files, so its dependency closure
// differs from the untagged one and both have to be unioned.
var tagConfigs = []string{"", "gogen_ir"}

// pathTriggers are repo-relative paths that feed the measured binaries but
// belong to no package's file list.
//
// This set is deliberately small. The .lg sources under pkg/rt/core and the
// compiled bundle are NOT listed: `go list` reports them as pkg/rt's
// EmbedFiles, so the package file set already covers them. Repeating them here
// would be a second copy of a fact the toolchain already knows, and the copy
// would be the one that goes stale.
var pathTriggers = []string{
	// Freshness manifests for the generated artifacts. Not embedded by any
	// package, but a change to them means the artifacts were regenerated.
	"pkg/rt/generated.sums",
	"pkg/rt/generated.manifest",

	// The module graph and the toolchain both change generated code beneath
	// every package. A Go version bump can move allocation counts on its own.
	"go.mod",
	"go.sum",
	"mise.toml",

	// BenchmarkClojureTestSuite reads its workload from disk at run time, so
	// no package's file list names it. The namespace loader is rooted at
	// test/compat and the vendored suite, and runSuiteOnce compiles
	// test/compat/clojure/core-test/portability.lg directly. A submodule bump
	// reaches the change set as the bare gitlink path, which the prefix form
	// does not match, so both spellings are listed.
	"test/compat/",
	"test/clojure-test-suite",
	"test/clojure-test-suite/",
}

func main() {
	var (
		filesFrom = flag.String("files-from", "-", "read newline-separated changed paths from this file (\"-\" for stdin)")
		verbose   = flag.Bool("v", false, "explain the decision on stderr")
		summary   = flag.Bool("summary", false, "read jj diff --summary lines instead of plain paths")
	)
	flag.Parse()

	paths, err := readPaths(*filesFrom, *summary)
	if err != nil {
		// Cannot read the change set: assume the worst and run.
		fmt.Fprintf(os.Stderr, "ratchet-scope: cannot read changed paths (%v); assuming affected\n", err)
		fmt.Println("RUN")
		return
	}
	if len(paths) == 0 {
		// An empty change set is not evidence of safety — it usually means the
		// range was computed wrongly.
		fmt.Fprintln(os.Stderr, "ratchet-scope: empty change set; assuming affected")
		fmt.Println("RUN")
		return
	}

	files, err := dependencyFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ratchet-scope: cannot expand dependencies (%v); assuming affected\n", err)
		fmt.Println("RUN")
		return
	}

	root, err := moduleRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ratchet-scope: cannot locate the module root (%v); assuming affected\n", err)
		fmt.Println("RUN")
		return
	}
	exists := func(p string) bool {
		_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(p)))
		return err == nil
	}

	if reason := affected(paths, files, exists); reason != "" {
		if *verbose {
			fmt.Fprintf(os.Stderr, "ratchet-scope: affected — %s\n", reason)
		}
		fmt.Println("RUN")
		return
	}
	if *verbose {
		fmt.Fprintf(os.Stderr, "ratchet-scope: %d changed path(s), none compiled into or embedded by the measured binaries\n", len(paths))
	}
	fmt.Println("SKIP")
}

// affected returns a human-readable reason naming the first path that puts the
// change set in scope, or "" when none does. It reports a reason rather than a
// bool so the hook can say WHY it is spending several minutes on benchmarks.
//
// exists reports whether a changed path is present in the checked-out tree.
// The hook only asks for a decision when that tree is the one being pushed, so
// a changed path that is absent was deleted by the push.
func affected(paths []string, depFiles map[string]string, exists func(string) bool) string {
	var dirs map[string]bool
	for _, p := range paths {
		if t := matchedTrigger(p); t != "" {
			return fmt.Sprintf("%s matches %s", p, t)
		}
		if pkg := matchedDepFile(p, depFiles); pkg != "" {
			return fmt.Sprintf("%s is compiled into or embedded by %s", p, pkg)
		}
		if !exists(normalize(p)) {
			if dirs == nil {
				dirs = closureDirs(depFiles)
			}
			if d := deletedUnderClosure(p, dirs); d != "" {
				return fmt.Sprintf("%s was deleted from %s, which holds measured files", p, d)
			}
		}
	}
	return ""
}

// deletedUnderClosure reports the nearest ancestor directory of a deleted path
// that holds a file of the measured closure, or "" when there is none.
//
// The closure is computed from the tree after the push, so a deleted file is
// never in it, even when the push removed a benchmark or an embedded asset.
// Rebuilding the closure at the base would need a second checkout. A deletion
// anywhere beneath a directory that still holds measured files is treated as
// affected instead; this over-approximates, which is the safe direction.
// Ancestors rather than the immediate directory, because deleting a whole
// embedded subtree such as pkg/rt/core/<dir>/ leaves no file in that directory
// but does change what pkg/rt embeds.
func deletedUnderClosure(p string, dirs map[string]bool) string {
	for d := path.Dir(normalize(p)); d != "." && d != "/"; d = path.Dir(d) {
		if dirs[d] {
			return d
		}
	}
	return ""
}

// closureDirs returns every directory that holds at least one file of the
// measured closure.
func closureDirs(depFiles map[string]string) map[string]bool {
	dirs := make(map[string]bool, len(depFiles))
	for f := range depFiles {
		dirs[path.Dir(f)] = true
	}
	return dirs
}

// matchedTrigger reports the pathTriggers entry covering p, if any.
func matchedTrigger(p string) string {
	p = normalize(p)
	for _, t := range pathTriggers {
		if strings.HasSuffix(t, "/") {
			if strings.HasPrefix(p, t) {
				return t
			}
			continue
		}
		if p == t {
			return t
		}
	}
	return ""
}

// matchedDepFile reports the package that compiles or embeds p, if any.
//
// Matching is on the package's ACTUAL FILE LIST, as `go list` reports it,
// which covers //go:embed assets — an embedded file changes a measurement with
// no .go file changing — while excluding files that merely share a directory
// with a package. Both halves are load-bearing:
//
//	test/quality_inputs_test.lg sits in the `test` package's directory, and
//	`test` is a benchmark root, but it belongs to no build and
//	BenchmarkClojureTestSuite reads only clojure-test-suite/**/*.cljc. Matching
//	by directory would charge every .lg test edit a full benchmark run.
func matchedDepFile(p string, depFiles map[string]string) string {
	return depFiles[normalize(p)]
}

// normalize makes a path comparable: slash-separated, no leading "./".
func normalize(p string) string {
	p = filepath.ToSlash(strings.TrimSpace(p))
	return strings.TrimPrefix(p, "./")
}

// fileListFields are the `go list` fields naming files that end up in, or are
// embedded into, a package's build. IgnoredGoFiles is included deliberately: a
// file excluded by build constraints under one configuration is compiled under
// another, and a constraint we do not enumerate here would otherwise read as
// "not part of any package" — the unsafe direction.
var fileListFields = []string{
	"GoFiles", "CgoFiles", "TestGoFiles", "XTestGoFiles", "IgnoredGoFiles",
	"EmbedFiles", "TestEmbedFiles", "XTestEmbedFiles",
}

// dependencyFiles maps every repo-relative file compiled into or embedded by
// the benchmark binaries to the package that owns it, across all tag
// configurations.
func dependencyFiles() (map[string]string, error) {
	root, err := moduleRoot()
	if err != nil {
		return nil, err
	}
	files := map[string]string{}
	for _, tags := range tagConfigs {
		if err := listDepFiles(tags, root, files); err != nil {
			return nil, fmt.Errorf("tags %q: %w", tags, err)
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("dependency expansion produced no in-module files")
	}
	return files, nil
}

// listDepFiles runs `go list -deps -test` over the roots and records each
// package's files. -test is required: the benchmarks live in _test.go files,
// so without it the closure omits anything only the test binary pulls in.
func listDepFiles(tags, root string, into map[string]string) error {
	var tmpl strings.Builder
	tmpl.WriteString("{{.ImportPath}}\t{{.Dir}}")
	for _, f := range fileListFields {
		fmt.Fprintf(&tmpl, "\t{{join .%s \",\"}}", f)
	}
	args := []string{"list", "-deps", "-test", "-f", tmpl.String()}
	if tags != "" {
		args = append(args, "-tags", tags)
	}
	args = append(args, benchmarkRoots...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(out), "\n") {
		cols := strings.Split(strings.TrimRight(line, "\r"), "\t")
		if len(cols) < 3 {
			continue
		}
		pkg, dir := cols[0], cols[1]
		rel, err := filepath.Rel(root, dir)
		if err != nil || strings.HasPrefix(rel, "..") {
			// Outside the module (standard library, module cache). Those
			// cannot appear in a change set for this repository.
			continue
		}
		for _, group := range cols[2:] {
			for _, name := range strings.Split(group, ",") {
				if name = strings.TrimSpace(name); name != "" {
					into[filepath.ToSlash(filepath.Join(rel, name))] = pkg
				}
			}
		}
	}
	return nil
}

// moduleRoot returns the absolute path of the module being built.
func moduleRoot() (string, error) {
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// readPaths reads newline-separated paths from a file or stdin. With summary
// set, each line is a `jj diff --summary` line instead, expanded by
// summaryPaths.
func readPaths(from string, summary bool) ([]string, error) {
	f := os.Stdin
	if from != "-" {
		var err error
		if f, err = os.Open(from); err != nil {
			return nil, err
		}
		defer f.Close()
	}
	var paths []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
		case summary:
			ps, err := summaryPaths(line)
			if err != nil {
				return nil, err
			}
			paths = append(paths, ps...)
		default:
			paths = append(paths, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

// summaryPaths returns the paths one `jj diff --summary` line names. A rename
// or copy is spelled `R prefix{old => new}suffix` and yields both sides: the
// old side of a rename is a deletion, and a rename out of the closure is caught
// only by the deletion rule seeing it.
func summaryPaths(line string) ([]string, error) {
	if len(line) < 3 || line[1] != ' ' || !strings.ContainsRune("MADRC", rune(line[0])) {
		return nil, fmt.Errorf("unrecognised summary line %q", line)
	}
	p := line[2:]
	if line[0] != 'R' && line[0] != 'C' {
		return []string{p}, nil
	}
	open, arrow, end := strings.IndexByte(p, '{'), strings.Index(p, " => "), strings.LastIndexByte(p, '}')
	if open < 0 || arrow < open || end < arrow {
		return nil, fmt.Errorf("unrecognised rename %q", line)
	}
	pre, post := p[:open], p[end+1:]
	return []string{
		path.Clean(pre + p[open+1:arrow] + post),
		path.Clean(pre + p[arrow+4:end] + post),
	}, nil
}
