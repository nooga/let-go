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
	"runtime"
	"testing"
	"time"
)

// TestLgCompileCommand drives `lg compile` (#596) end to end: lowering, the
// program.lgb child process, go.mod scaffolding, and go build, all from one
// invocation. The fixture is TestNativeEntryNamespaceReplay's, which runs the
// same pipeline step by step, so the two tests pin the same observable output.
func TestLgCompileCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("builds lg and a native Go binary")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := repoRoot(t)
	// -buildvcs=false keeps this a dev build. From a clean checkout, which is
	// what CI has, build info would stamp a pseudo-version for the commit, and
	// the generated go.mod would try to fetch that commit from the module
	// proxy instead of using LETGO_SRC below.
	bin := filepath.Join(t.TempDir(), "lg")
	if out, err := runCmd(ctx, t, "go", root, []string{"build", "-buildvcs=false", "-o", bin, "."}); err != nil {
		t.Fatalf("build lg: %v\n%s", err, out)
	}
	fix := t.TempDir()
	lib := filepath.Join(fix, "libx", "core.lg")
	if err := os.MkdirAll(filepath.Dir(lib), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lib, []byte(`(ns libx.core)
(defn helper [a] (+ a 41))
(when-not *compiling-aot* (println "libx: runtime init ran"))
`), 0644); err != nil {
		t.Fatal(err)
	}
	app := filepath.Join(fix, "app.lg")
	if err := os.WriteFile(app, []byte(`(ns app (:require [libx.core :as lib]))
(defn -main [] (println (lib/helper 1)))
(when-not *compiling-aot*
  (println "app: runtime init ran")
  (-main))
`), 0644); err != nil {
		t.Fatal(err)
	}
	// LETGO_SRC points the dev build's generated module at this checkout.
	env := []string{
		"LG_SOURCE_PATHS=" + fix + string(os.PathListSeparator) + filepath.Join(root, "pkg", "rt", "core"),
		"LETGO_SRC=" + root,
	}

	t.Run("builds and runs", func(t *testing.T) {
		work := t.TempDir()
		exe := filepath.Join(t.TempDir(), "app.native")
		if out, err := runCmd(ctx, t, bin, fix,
			[]string{"compile", "-o", exe, "-work", work, lib, app}, env...); err != nil {
			t.Fatalf("lg compile: %v\n%s", err, out)
		}
		stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
		if err != nil || exit != 0 {
			t.Fatalf("run: err=%v exit=%d stdout=%q stderr=%q", err, exit, stdout, stderr)
		}
		want := []byte("libx: runtime init ran\napp: runtime init ran\n42\n")
		if !bytes.Equal(stdout, want) {
			t.Fatalf("stdout %q, want %q; stderr=%q", stdout, want, stderr)
		}
		// -work keeps the generated module where the caller asked for it.
		for _, f := range []string{"main.go", "program.lgb", "go.mod"} {
			if _, err := os.Stat(filepath.Join(work, f)); err != nil {
				t.Errorf("-work dir is missing %s: %v", f, err)
			}
		}
	})

	// -lowered-core lowers let-go's own core into the generated module and
	// links it. The program reports how a core function prints: a lowered one
	// is a native fn, and it has to give the VM's answers.
	t.Run("links the lowered core", func(t *testing.T) {
		probe := filepath.Join(fix, "coreprobe.lg")
		if err := os.WriteFile(probe, []byte(`(ns coreprobe (:require [clojure.walk :as walk]))
(defn -main []
  (println (boolean (re-find #"native" (str filterv))))
  (println (filterv odd? [1 2 3]))
  (println (filterv (reify IFn (-invoke [_ args] (pos? (first args)))) [1 -1 2]))
  (println (walk/postwalk (fn [x] (if (number? x) (inc x) x)) {:a [1 2]})))
(when-not *compiling-aot* (-main))
`), 0644); err != nil {
			t.Fatal(err)
		}
		const want = "[1 3]\n[1 2]\n{:a [2 3]}\n"
		for _, tc := range []struct {
			name   string
			flags  []string
			native string
		}{
			{"without the flag", nil, "false\n"},
			{"with the flag", []string{"-lowered-core"}, "true\n"},
		} {
			work := t.TempDir()
			exe := filepath.Join(t.TempDir(), "coreprobe.native")
			args := append([]string{"compile", "-o", exe, "-work", work}, tc.flags...)
			if out, err := runCmd(ctx, t, bin, fix, append(args, probe), env...); err != nil {
				t.Fatalf("%s: lg compile: %v\n%s", tc.name, err, out)
			}
			stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
			if err != nil || exit != 0 {
				t.Fatalf("%s: run: err=%v exit=%d stdout=%q stderr=%q", tc.name, err, exit, stdout, stderr)
			}
			if got := string(stdout); got != tc.native+want {
				t.Fatalf("%s: stdout %q, want %q", tc.name, got, tc.native+want)
			}
			_, statErr := os.Stat(filepath.Join(work, "lgcore", "core", "core.go"))
			if lowered := tc.flags != nil; lowered != (statErr == nil) {
				t.Fatalf("%s: lgcore/core/core.go present=%v", tc.name, statErr == nil)
			}
		}

		// One kept work directory, named relative to the cwd, built with the
		// flag and then without it: the second binary must be back on the
		// bytecode core, not still linking the tree the first build left.
		exe := filepath.Join(t.TempDir(), "coreprobe.native")
		for _, tc := range []struct {
			flags  []string
			native string
		}{{[]string{"-lowered-core"}, "true\n"}, {nil, "false\n"}} {
			args := append([]string{"compile", "-o", exe, "-work", "kept-work"}, tc.flags...)
			if out, err := runCmd(ctx, t, bin, fix, append(args, probe), env...); err != nil {
				t.Fatalf("relative -work, flags %v: lg compile: %v\n%s", tc.flags, err, out)
			}
			stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
			if err != nil || exit != 0 {
				t.Fatalf("relative -work, flags %v: run: err=%v exit=%d stderr=%q", tc.flags, err, exit, stderr)
			}
			if got := string(stdout); got != tc.native+want {
				t.Fatalf("relative -work, flags %v: stdout %q, want %q", tc.flags, got, tc.native+want)
			}
		}
		if _, err := os.Stat(filepath.Join(fix, "kept-work", "kept-work")); err == nil {
			t.Fatal("the core was lowered into kept-work/kept-work: the relative -work path was resolved twice")
		}

		// Cross-compiling: lgbgen still has to be built for, and run on, this
		// machine. Only the build is checked; the binary is for another OS.
		cross := filepath.Join(t.TempDir(), "coreprobe.cross")
		target := "linux"
		if runtime.GOOS == "linux" {
			target = "darwin"
		}
		crossEnv := append(append([]string(nil), env...), "GOOS="+target, "GOARCH=arm64", "CGO_ENABLED=0")
		if out, err := runCmd(ctx, t, bin, fix, []string{"compile", "-lowered-core", "-o", cross, probe}, crossEnv...); err != nil {
			t.Fatalf("cross-compiling to %s with -lowered-core: %v\n%s", target, err, out)
		}

		// A program namespace lgcore.* lowers to the directory the core needs.
		clash := filepath.Join(fix, "lgcore", "core.lg")
		if err := os.MkdirAll(filepath.Dir(clash), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(clash, []byte("(ns lgcore.core)\n(defn -main [] (println 42))\n(when-not *compiling-aot* (-main))\n"), 0644); err != nil {
			t.Fatal(err)
		}
		out, err := runCmd(ctx, t, bin, fix, []string{"compile", "-lowered-core", "-o", exe, clash}, env...)
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
			t.Fatalf("lgcore.* program with -lowered-core: err=%v, want exit 1\n%s", err, out)
		}
		if !bytes.Contains([]byte(out), []byte("a program package is already there")) {
			t.Fatalf("missing the lgcore clash diagnostic:\n%s", out)
		}
	})

	t.Run("skips data-only input without importing its package", func(t *testing.T) {
		consts := filepath.Join(fix, "demo", "consts.lg")
		if err := os.MkdirAll(filepath.Dir(consts), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(consts, []byte(`(ns demo.consts)
(def greeting "hello")
(when-not *compiling-aot* (println "consts: runtime init ran"))
`), 0644); err != nil {
			t.Fatal(err)
		}
		app := filepath.Join(fix, "demo", "app.lg")
		if err := os.WriteFile(app, []byte(`(ns demo.app (:require [demo.consts :as c] [libx.core :as lib]))
(defn -main [] (println c/greeting) (println (lib/helper 1)))
`), 0644); err != nil {
			t.Fatal(err)
		}
		work := t.TempDir()
		exe := filepath.Join(t.TempDir(), "app.native")
		out, err := runCmd(ctx, t, bin, fix,
			[]string{"compile", "-o", exe, "-work", work, consts, lib, app}, env...)
		if err != nil {
			t.Fatalf("lg compile: %v\n%s", err, out)
		}
		if !bytes.Contains([]byte(out), []byte(consts+": no lowerable defns — skipped")) {
			t.Fatalf("missing skipped-input diagnostic:\n%s", out)
		}
		main, err := os.ReadFile(filepath.Join(work, "main.go"))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(main, []byte(`"lgprogram/demo/consts"`)) {
			t.Fatalf("main.go imports a skipped package:\n%s", main)
		}
		stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
		if err != nil || exit != 0 {
			t.Fatalf("run: err=%v exit=%d stdout=%q stderr=%q", err, exit, stdout, stderr)
		}
		want := []byte("consts: runtime init ran\nlibx: runtime init ran\nhello\n42\n")
		if !bytes.Equal(stdout, want) {
			t.Fatalf("stdout %q, want %q; stderr=%q", stdout, want, stderr)
		}
	})

	// A namespace's last segment names its Go package. main would make it a
	// program, which the entry frame cannot import, init cannot be imported
	// under its own name, and go and type are not identifiers; each takes a
	// trailing underscore, on the directory too.
	t.Run("namespaces named main, init, go and type", func(t *testing.T) {
		files := map[string]string{
			"named/init.lg": "(ns named.init)\n(defn setup [] :set)\n",
			"named/type.lg": "(ns named.type (:require [named.init :as i]))\n(defn kind [] [:typed (i/setup)])\n",
			"named/go.lg":   "(ns named.go (:require [named.type :as t]))\n(defn run [] [:ran (t/kind)])\n",
			"named/main.lg": `(ns named.main (:require [named.go :as g]))
(defn -main [] (println (g/run)))
(when-not *compiling-aot* (-main))
`,
		}
		var paths []string
		for _, rel := range []string{"named/init.lg", "named/type.lg", "named/go.lg", "named/main.lg"} {
			path := filepath.Join(fix, rel)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(files[rel]), 0644); err != nil {
				t.Fatal(err)
			}
			paths = append(paths, path)
		}
		work := t.TempDir()
		exe := filepath.Join(t.TempDir(), "named.native")
		if out, err := runCmd(ctx, t, bin, fix,
			append([]string{"compile", "-o", exe, "-work", work}, paths...), env...); err != nil {
			t.Fatalf("lg compile: %v\n%s", err, out)
		}
		if _, err := os.Stat(filepath.Join(work, "named", "main_", "main_.go")); err != nil {
			t.Errorf("named.main was not lowered into named/main_: %v", err)
		}
		stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
		if err != nil || exit != 0 || string(stdout) != "[:ran [:typed :set]]\n" {
			t.Fatalf("run: err=%v exit=%d stdout=%q stderr=%q", err, exit, stdout, stderr)
		}
	})

	// named.main and named.main- would both be written to named/main_, the
	// second over the first, and calls into one would reach the other.
	t.Run("two namespaces that lower to one directory are refused", func(t *testing.T) {
		var paths []string
		for rel, src := range map[string]string{
			"twin/main.lg":  "(ns twin.main)\n(defn v [] 11)\n",
			"twin/main_.lg": "(ns twin.main-)\n(defn v [] 22)\n",
		} {
			path := filepath.Join(fix, rel)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(src), 0644); err != nil {
				t.Fatal(err)
			}
			paths = append(paths, path)
		}
		out, err := runCmd(ctx, t, bin, fix,
			append([]string{"compile", "-o", filepath.Join(t.TempDir(), "twin")}, paths...), env...)
		if err == nil {
			t.Fatalf("lg compile built two namespaces into one directory:\n%s", out)
		}
		if !bytes.Contains([]byte(out), []byte("twin.main and twin.main- both lower to the Go package directory twin/main_")) {
			t.Fatalf("missing the diagnostic naming both namespaces:\n%s", out)
		}
	})

	t.Run("usage errors exit 2", func(t *testing.T) {
		for _, args := range [][]string{{"compile"}, {"compile", "-bogus", app}, {"compile", "-o"}} {
			cmd := exec.CommandContext(ctx, bin, args...)
			cmd.Dir = fix
			out, err := cmd.CombinedOutput()
			if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 2 {
				t.Errorf("lg %v: err=%v, want exit 2\n%s", args, err, out)
			}
		}
	})

	t.Run("a default output name that is a directory fails", func(t *testing.T) {
		// app.lg defaults to ./app; go build -o would write app/lgprogram.
		dir := filepath.Join(fix, "app")
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)
		out, err := runCmd(ctx, t, bin, fix, []string{"compile", lib, app}, env...)
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
			t.Fatalf("lg compile into a directory: err=%v, want exit 1\n%s", err, out)
		}
		if !bytes.Contains([]byte(out), []byte("is a directory")) {
			t.Fatalf("missing the directory diagnostic:\n%s", out)
		}
		if entries, _ := os.ReadDir(dir); len(entries) != 0 {
			t.Fatalf("lg compile wrote into %s: %v", dir, entries)
		}
	})

	// A def is an expression whose value is its var. A lowered function that
	// uses that value (here by returning it) has to build and give the VM's
	// answer; with-out-str is the same case, the def being the last form of
	// the closure the macro builds. The def whose value is dropped and the def
	// with no value pin the two neighbouring paths through the same code.
	t.Run("a def whose var is used as a value", func(t *testing.T) {
		prog := filepath.Join(fix, "defvalue.lg")
		if err := os.WriteFile(prog, []byte(`(ns defvalue)
(defn make [] (def answer 42))
(defn tagged [] (def ^:private hidden 7))
(defn unused [] (def ^:private side 8) nil)
(defn bare [] (def unbound))
(defn -main []
  (println (make) answer)
  (println (tagged) hidden)
  (println (unused) side)
  (println (bare))
  (println (count (with-out-str (def quiet 1))) quiet))
(when-not *compiling-aot* (-main))
`), 0644); err != nil {
			t.Fatal(err)
		}
		exe := filepath.Join(t.TempDir(), "defvalue.native")
		if out, err := runCmd(ctx, t, bin, fix, []string{"compile", "-o", exe, prog}, env...); err != nil {
			t.Fatalf("lg compile: %v\n%s", err, out)
		}
		stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
		if err != nil || exit != 0 {
			t.Fatalf("run: err=%v exit=%d stdout=%q stderr=%q", err, exit, stdout, stderr)
		}
		want := "#'defvalue/answer 42\n#'defvalue/hidden 7\nnil 8\n#'defvalue/unbound\n0 1\n"
		if string(stdout) != want {
			t.Fatalf("stdout %q, want %q; stderr=%q", stdout, want, stderr)
		}
	})

	t.Run("a program without an entry fails", func(t *testing.T) {
		out, err := runCmd(ctx, t, bin, fix, []string{"compile", lib}, env...)
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
			t.Fatalf("lg compile on an entry-less file: err=%v, want exit 1\n%s", err, out)
		}
		if !bytes.Contains([]byte(out), []byte("no program entry")) {
			t.Fatalf("missing the no-entry diagnostic:\n%s", out)
		}
	})
}
