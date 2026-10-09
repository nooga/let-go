/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLgCompileNonLiteralDefs pins #1049: lg-compile loads and lowers a
// definition that is not a literal top-level form, one a macro emits or one
// wrapped in do, the same as a literal one. Before the fix each case failed
// with "Can't resolve" when lowering the code that refers to it.
func TestLgCompileNonLiteralDefs(t *testing.T) {
	bin := buildLG(t)
	root := repoRoot(t)

	// compile writes each named file into one temp dir, runs lg-compile on
	// them in order, and returns the output dir, the combined log and the error.
	compile := func(t *testing.T, files ...[2]string) (string, string, error) {
		t.Helper()
		dir := t.TempDir()
		args := []string{"scripts/lg-compile", t.TempDir(), "tmpmod"}
		for _, f := range files {
			p := filepath.Join(dir, f[0])
			if err := os.WriteFile(p, []byte(f[1]), 0644); err != nil {
				t.Fatal(err)
			}
			args = append(args, p)
		}
		cmd := exec.Command(bin, args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "LG_SOURCE_PATHS=pkg/rt/core")
		out, err := cmd.CombinedOutput()
		return args[1], string(out), err
	}
	mustCompile := func(t *testing.T, files ...[2]string) (string, string) {
		t.Helper()
		outDir, log, err := compile(t, files...)
		if err != nil {
			t.Fatalf("lg-compile: %v\n%s", err, log)
		}
		return outDir, log
	}
	goSrc := func(t *testing.T, outDir, pkg string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(outDir, pkg, pkg+".go"))
		if err != nil {
			t.Fatalf("lowered package %s missing: %v", pkg, err)
		}
		return string(b)
	}
	mustLower := func(t *testing.T, src string, funcs ...string) {
		t.Helper()
		for _, f := range funcs {
			if !strings.Contains(src, "\nfunc "+f+"(") {
				t.Fatalf("want lowered func %s; got:\n%s", f, src)
			}
		}
	}

	t.Run("defns a macro emits in a do", func(t *testing.T) {
		out, _ := mustCompile(t, [2]string{"app.lg", "(ns app)\n" +
			"(defmacro both [& ds] `(do ~@ds))\n" +
			"(both (defn f [x] (inc x)) (defn g [x] (* 2 x)))\n" +
			"(defn -main [& _] (println (f 1) (g 3)))\n"})
		mustLower(t, goSrc(t, out, "app"), "F", "G", "Main")
	})

	t.Run("a single defn a macro returns", func(t *testing.T) {
		out, _ := mustCompile(t, [2]string{"app.lg", "(ns app)\n" +
			"(defmacro one [d] d)\n" +
			"(one (defn f [x] (inc x)))\n" +
			"(defn -main [& _] (println (f 1)))\n"})
		mustLower(t, goSrc(t, out, "app"), "F", "Main")
	})

	t.Run("defn and def wrapped in do", func(t *testing.T) {
		out, _ := mustCompile(t, [2]string{"app.lg", "(ns app)\n" +
			"(do (defn f [x] (inc x)) (def k 41))\n" +
			"(defn -main [& _] (println (f k)))\n"})
		mustLower(t, goSrc(t, out, "app"), "F", "Main")
	})

	t.Run("a macro-emitted defn- stays private", func(t *testing.T) {
		out, _ := mustCompile(t, [2]string{"app.lg", "(ns app)\n" +
			"(defmacro private-fn [n] `(defn- ~n [x#] (+ x# 100)))\n" +
			"(private-fn hidden)\n" +
			"(defn -main [& _] (println (hidden 1)))\n"})
		mustLower(t, goSrc(t, out, "app"), "hidden", "Main")
	})

	t.Run("macro-emitted defmulti and defmethod", func(t *testing.T) {
		mustCompile(t, [2]string{"app.lg", "(ns app)\n" +
			"(defmacro shapes [] `(do (defmulti ~'area :shape)\n" +
			"                         (defmethod ~'area :sq [s#] (* (:w s#) (:w s#)))))\n" +
			"(shapes)\n" +
			"(defn -main [& _] (println (area {:shape :sq :w 3})))\n"})
	})

	// The macro is expanded while its own ns is current, not the last one loaded.
	t.Run("two files: the macro's ns loads before the caller's", func(t *testing.T) {
		out, _ := mustCompile(t,
			[2]string{"twoa.lg", "(ns twoa)\n(defmacro one [d] d)\n(one (defn f [x] (inc x)))\n"},
			[2]string{"twob.lg", "(ns twob (:require [twoa]))\n(defn -main [& _] (println (twoa/f 41)))\n"})
		mustLower(t, goSrc(t, out, "twoa"), "F")
		mustLower(t, goSrc(t, out, "twob"), "Main")
	})

	t.Run("a guarded entry call builds without a warning", func(t *testing.T) {
		_, log := mustCompile(t, [2]string{"app.lg", "(ns app)\n" +
			"(defn -main [& _] (println 1))\n" +
			"(comment (defn unused [] 1))\n" +
			"(when-not *compiling-aot* (-main))\n"})
		if strings.Contains(log, "warning") {
			t.Fatalf("want no warning; got:\n%s", log)
		}
	})

	// A definition under a conditional still isn't loaded; the warning names the form.
	t.Run("a conditional definition is named", func(t *testing.T) {
		_, log, err := compile(t, [2]string{"app.lg", "(ns app)\n" +
			"(when true (defn f [x] (inc x)))\n" +
			"(defn -main [& _] (println (f 1)))\n"})
		if err == nil {
			t.Fatalf("want a failure: f is not loaded at compile time; got:\n%s", log)
		}
		if !strings.Contains(log, "warning") || !strings.Contains(log, "(when true (defn f") {
			t.Fatalf("want a warning naming the when form; got:\n%s", log)
		}
	})
}
