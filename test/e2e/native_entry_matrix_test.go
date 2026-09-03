/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

// nativeEntryCase is one row of the #425 / #628 standalone native-entry matrix.
// Cases either fail at lg-compile (genFail) or produce a runnable binary.
type nativeEntryCase struct {
	name string
	// files: path relative to the fixture root → source contents.
	files map[string]string
	// compileArgs: extra files after out-dir/prefix, as paths relative to fixture root.
	// Defaults to sorted keys of files when nil.
	compileArgs []string

	genFail        bool
	genErrSubstr   string
	skipRun        bool // generation (+ optional main.go checks) only
	args           []string
	wantStdout     string
	wantExit       int
	mainContains   []string
	mainNotContain []string
	timeout        time.Duration
}

func TestNativeEntryMatrix(t *testing.T) {
	if testing.Short() {
		t.Skip("native-entry matrix builds lg + Go binaries")
	}
	bin := buildLG(t)
	root := repoRoot(t)

	cases := []nativeEntryCase{
		{
			name: "public -main []",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main [] (println \"PUB\"))\n",
			},
			wantStdout: "PUB",
		},
		{
			name: "private -main []",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn- -main [] (println \"PRIV\"))\n",
			},
			wantStdout:   "PRIV",
			mainContains: []string{"prog.NativeMain("},
		},
		{
			name: "public main []",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn main [] (println \"MAIN\"))\n",
			},
			wantStdout: "MAIN",
		},
		{
			name: "argv adaptation",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main [argv] (println (count argv)))\n",
			},
			args:       []string{"a", "b", "c"},
			wantStdout: "3",
		},
		{
			name: "variadic & args adaptation",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main [& args] (println (count args)))\n",
			},
			args:       []string{"x", "y"},
			wantStdout: "2",
		},
		{
			name: "fixed & rest rejection",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main [first & rest] (println first))\n",
			},
			genFail:      true,
			genErrSubstr: "unsupported arity",
		},
		{
			name: "main + -main collision calls Main__main",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn main [] (println \"WRONG\"))\n(defn -main [] (println \"RIGHT\"))\n",
			},
			wantStdout:     "RIGHT",
			mainContains:   []string{"prog.Main__main("},
			mainNotContain: []string{"prog.Main(ec)"},
		},
		{
			name: "competing -main in two namespaces",
			// Selection is compile-time (first -main in file order). lg -c takes
			// one input, so this case asserts the frame targets namespace a's
			// package rather than running a multi-ns binary.
			files: map[string]string{
				"a.lg": "(ns a)\n(defn -main [] (println \"A\"))\n",
				"b.lg": "(ns b)\n(defn -main [] (println \"B\"))\n",
			},
			compileArgs:  []string{"a.lg", "b.lg"},
			skipRun:      true,
			mainContains: []string{`prog "tmpmod/a"`, "prog.Main("},
		},
		{
			name: "*command-line-args* visible to -main",
			// Frame publishes SetCommandLineArgs after BootCore and before the
			// entry call. (Fully-lowered top-level forms run in Go init()
			// before main, so argv visibility is asserted from -main here;
			// emit order vs LoadProgramNamespaces is pinned in entry_frame_test.)
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main [] (println (pr-str *command-line-args*)))\n",
			},
			args:       []string{"p", "q"},
			wantStdout: `("p" "q")`,
		},
		{
			name: "direct-entry error propagation",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main [] (throw (ex-info \"BOOM\" {})))\n",
			},
			wantExit: 1,
			// FormatError text includes the message; accept substring.
			wantStdout: "", // may be empty; check via combined output in runner
		},
		{
			name: "return-hint on arity vector",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main ^long [] (println \"RH-ARITY\") 0)\n",
			},
			wantStdout: "RH-ARITY",
		},
		{
			name: "return-hint on name",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn ^long -main [] (println \"RH-NAME\") 0)\n",
			},
			wantStdout: "RH-NAME",
		},
		{
			name: "typed no-error direct call shape",
			files: map[string]string{
				// Constant long return: lowerer sets :needs-error? false, so the
				// frame must call prog.Main(ec) without an err binding and
				// without importing unused pkg/vm.
				"app.lg": "(ns app)\n(defn -main ^long [] 0)\n",
			},
			mainContains:   []string{"prog.Main(ec)"},
			mainNotContain: []string{"if _, err := prog.Main(", "github.com/nooga/let-go/pkg/vm"},
		},
		{
			name: "return-hinted unsupported arity rejected",
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main ^long [first & rest] 1)\n",
			},
			genFail:      true,
			genErrSubstr: "unsupported arity",
		},
		{
			name: "fib example spike inputs",
			files: map[string]string{
				"fib.lg": `(ns fib)
(defn fib [n]
  (if (< n 2) n (+ (fib (- n 1)) (fib (- n 2)))))
(defn -main [& args]
  (let [n (if (seq args) 34 10)]
    (println (fib n))))
`,
			},
			args:       nil,
			wantStdout: "55",
			timeout:    30 * time.Second,
		},
		{
			name: "vm-fallback frame is namespace-qualified",
			// Generation-only: emit a not-lowered frame via the entry-frame API
			// and assert InvokeProgramEntry carries the owning namespace. Runtime
			// resolution is covered by TestInvokeProgramEntryUsesSelectedNamespace.
			skipRun: true,
			files: map[string]string{
				"app.lg": "(ns app)\n(defn -main [] 1)\n",
			},
			// Filled specially in the runner — see below.
		},
		{
			name: "duplicate blank imports deduped",
			// Multi-file non-entry ns would otherwise emit _ "tmpmod/lib" twice
			// (one go-pkg per input file). Entry package is aliased, not blank.
			files: map[string]string{
				"lib/a.lg": "(ns lib)\n(defn helper [] 1)\n",
				"lib/b.lg": "(ns lib)\n(defn helper [] 1)\n",
				"app.lg":   "(ns app)\n(defn -main [] (println \"ok\"))\n",
			},
			compileArgs:  []string{"lib/a.lg", "lib/b.lg", "app.lg"},
			skipRun:      true,
			mainContains: []string{`_ "tmpmod/lib"`, `prog "tmpmod/app"`},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			timeout := tc.timeout
			if timeout == 0 {
				timeout = 60 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if tc.name == "vm-fallback frame is namespace-qualified" {
				runVMFallbackFrameSource(ctx, t, bin, root)
				return
			}

			fix := t.TempDir()
			for rel, body := range tc.files {
				p := filepath.Join(fix, rel)
				if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(body), 0644); err != nil {
					t.Fatal(err)
				}
			}
			outDir := t.TempDir()
			prefix := "tmpmod"

			args := []string{"scripts/lg-compile", "--entry-frame", outDir, prefix}
			compFiles := tc.compileArgs
			if compFiles == nil {
				for rel := range tc.files {
					compFiles = append(compFiles, rel)
				}
				// Stable order for determinism.
				compFiles = sortedStrings(compFiles)
			}
			for _, rel := range compFiles {
				args = append(args, filepath.Join(fix, rel))
			}

			compileOut, compileErr := runCmd(ctx, t, bin, root, args, "LG_SOURCE_PATHS=pkg/rt/core")
			if tc.genFail {
				if compileErr == nil {
					t.Fatalf("expected lg-compile to fail; output:\n%s", compileOut)
				}
				if !strings.Contains(compileOut, tc.genErrSubstr) {
					t.Fatalf("want stderr/stdout to contain %q; got:\n%s", tc.genErrSubstr, compileOut)
				}
				if _, err := os.Stat(filepath.Join(outDir, "main.go")); !os.IsNotExist(err) {
					t.Fatalf("main.go must not be written on gen failure")
				}
				return
			}
			if compileErr != nil {
				t.Fatalf("lg-compile: %v\n%s", compileErr, compileOut)
			}

			mainPath := filepath.Join(outDir, "main.go")
			mainBody, err := os.ReadFile(mainPath)
			if err != nil {
				t.Fatalf("main.go missing: %v\n%s", err, compileOut)
			}
			mainStr := string(mainBody)
			for _, s := range tc.mainContains {
				if !strings.Contains(mainStr, s) {
					t.Fatalf("main.go missing %q:\n%s", s, mainStr)
				}
			}
			for _, s := range tc.mainNotContain {
				if strings.Contains(mainStr, s) {
					t.Fatalf("main.go unexpectedly contains %q:\n%s", s, mainStr)
				}
			}
			if tc.name == "duplicate blank imports deduped" {
				if c := strings.Count(mainStr, `_ "tmpmod/lib"`); c != 1 {
					t.Fatalf("want exactly one blank import of tmpmod/lib, got %d:\n%s", c, mainStr)
				}
			}
			if tc.skipRun {
				return
			}

			// Compile program.lgb from the same sources.
			lgb := filepath.Join(outDir, "program.lgb")
			lgcArgs := []string{"-c", lgb}
			for _, rel := range compFiles {
				lgcArgs = append(lgcArgs, filepath.Join(fix, rel))
			}
			if out, err := runCmd(ctx, t, bin, root, lgcArgs); err != nil {
				t.Fatalf("lg -c: %v\n%s", err, out)
			}

			mod := "module " + prefix + "\n\ngo 1.23\n\nrequire github.com/nooga/let-go v0.0.0\nreplace github.com/nooga/let-go => " + root + "\n"
			if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(mod), 0644); err != nil {
				t.Fatal(err)
			}
			if out, err := runCmd(ctx, t, "go", outDir, []string{"mod", "tidy"}); err != nil {
				t.Fatalf("go mod tidy: %v\n%s", err, out)
			}
			exe := filepath.Join(outDir, "app.native")
			if out, err := runCmd(ctx, t, "go", outDir, []string{"build", "-o", exe, "."}); err != nil {
				t.Fatalf("go build: %v\n%s", err, out)
			}

			runArgs := append([]string{}, tc.args...)
			cmd := exec.CommandContext(ctx, exe, runArgs...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			gotExit := 0
			if err != nil {
				if ee, ok := err.(*exec.ExitError); ok {
					gotExit = ee.ExitCode()
				} else {
					t.Fatalf("run: %v\nstdout=%s\nstderr=%s", err, stdout.String(), stderr.String())
				}
			}
			wantExit := tc.wantExit
			if gotExit != wantExit {
				t.Fatalf("exit %d, want %d\nstdout=%q\nstderr=%q", gotExit, wantExit, stdout.String(), stderr.String())
			}
			if tc.name == "direct-entry error propagation" {
				combined := stdout.String() + stderr.String()
				if !strings.Contains(combined, "BOOM") {
					t.Fatalf("want BOOM in output; stdout=%q stderr=%q", stdout.String(), stderr.String())
				}
				return
			}
			got := strings.TrimSpace(stdout.String())
			if got != tc.wantStdout {
				t.Fatalf("stdout %q, want %q\nstderr=%q", got, tc.wantStdout, stderr.String())
			}
		})
	}

	// Extra fib(34) check — same fixture, spike input.
	t.Run("fib 34", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		fix := t.TempDir()
		src := filepath.Join(fix, "fib.lg")
		body := `(ns fib)
(defn fib [n]
  (if (< n 2) n (+ (fib (- n 1)) (fib (- n 2)))))
(defn -main [& args]
  (let [n (if (seq args) 34 10)]
    (println (fib n))))
`
		if err := os.WriteFile(src, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
		outDir := t.TempDir()
		if out, err := runCmd(ctx, t, bin, root, []string{"scripts/lg-compile", "--entry-frame", outDir, "tmpmod", src}, "LG_SOURCE_PATHS=pkg/rt/core"); err != nil {
			t.Fatalf("compile: %v\n%s", err, out)
		}
		lgb := filepath.Join(outDir, "program.lgb")
		if out, err := runCmd(ctx, t, bin, root, []string{"-c", lgb, src}); err != nil {
			t.Fatalf("lg -c: %v\n%s", err, out)
		}
		mod := "module tmpmod\n\ngo 1.23\n\nrequire github.com/nooga/let-go v0.0.0\nreplace github.com/nooga/let-go => " + root + "\n"
		if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(mod), 0644); err != nil {
			t.Fatal(err)
		}
		if out, err := runCmd(ctx, t, "go", outDir, []string{"mod", "tidy"}); err != nil {
			t.Fatal(err, out)
		}
		exe := filepath.Join(outDir, "fib.native")
		if out, err := runCmd(ctx, t, "go", outDir, []string{"build", "-o", exe, "."}); err != nil {
			t.Fatal(err, out)
		}
		cmd := exec.CommandContext(ctx, exe, "x")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("run: %v\n%s", err, out)
		}
		if got := strings.TrimSpace(string(out)); got != "5702887" {
			t.Fatalf("got %q, want 5702887", got)
		}
	})
}

func runVMFallbackFrameSource(ctx context.Context, t *testing.T, bin, root string) {
	t.Helper()
	expr := `
(require 'ir.passes.entry-frame)
(let [raw (ir.passes.entry-frame/find-program-entry '[(defn -main [] 1)])
      e   (assoc (ir.passes.entry-frame/analyze-entry raw {:fn nil}) :ns "winner")
      src (ir.passes.entry-frame/emit-entry-frame-go {:entry e})]
  (print src))
`
	out, err := runCmd(ctx, t, bin, root, []string{"-e", expr}, "LG_SOURCE_PATHS=pkg/rt/core")
	if err != nil {
		t.Fatalf("lg -e: %v\n%s", err, out)
	}
	if !strings.Contains(out, `rt.InvokeProgramEntry(ec, "winner", "-main", argv)`) {
		t.Fatalf("fallback frame missing ns-qualified InvokeProgramEntry; got:\n%s", out)
	}
	if strings.Contains(out, `InvokeProgramEntry(ec, "-main"`) {
		t.Fatalf("legacy 3-arg InvokeProgramEntry must not appear:\n%s", out)
	}
}

func runCmd(ctx context.Context, t *testing.T, bin, dir string, args []string, extraEnv ...string) (string, error) {
	t.Helper()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), extraEnv...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
