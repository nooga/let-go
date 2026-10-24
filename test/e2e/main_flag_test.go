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

func TestMainFlag(t *testing.T) {
	bin := buildLG(t)
	dir := t.TempDir()
	fixtures := map[string]string{
		"app.clj": `(ns app) (prn :load *command-line-args*)
(defn -main [& args] (prn :main args *command-line-args*) :returned)
(defn other [& args] (prn :other args) :returned)`,
		"missing.lg": `(ns missing)`,
		"notfn.lg":   `(ns notfn) (def -main 42)`,
		"boom.lg":    `(ns boom) (defn -main [& args] (throw (ex-info "entry failed" {})))`,
		"badload.lg": `(ns badload) (throw (ex-info "load failed" {}))`,
		"once.lg":    `(ns once) (println "loaded") (defn -main [] (println "called"))`,
	}
	for name, body := range fixtures {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	tests := []struct {
		name     string
		args     []string
		want     string
		code     int
		contains bool
	}{
		{name: "args", args: []string{"-m", "app", "a", "b"}, want: ":load (\"a\" \"b\")\n:main (\"a\" \"b\") (\"a\" \"b\")"},
		{name: "no args", args: []string{"-m", "app"}, want: ":load nil\n:main nil nil"},
		{name: "alias", args: []string{"--main", "app", "x"}, want: ":load (\"x\")\n:main (\"x\") (\"x\")"},
		{name: "equals", args: []string{"-m=app", "x"}, want: ":load (\"x\")\n:main (\"x\") (\"x\")"},
		{name: "verbatim flags", args: []string{"-m", "app", "-v", "--", "-m", "", "a b"}, want: ":load (\"-v\" \"--\" \"-m\" \"\" \"a b\")\n:main (\"-v\" \"--\" \"-m\" \"\" \"a b\") (\"-v\" \"--\" \"-m\" \"\" \"a b\")"},
		{name: "qualified", args: []string{"-m", "app/other", "x"}, want: ":load (\"x\")\n:other (\"x\")"},
		{name: "init before require", args: []string{"-e", `(prn :init *command-line-args*)`, "-m", "app", "x"}, want: ":init (\"x\")\nnil\n:load (\"x\")\n:main (\"x\") (\"x\")"},
		{name: "already required", args: []string{"-e", `(require 'once)`, "-m", "once"}, want: "loaded\nnil\ncalled"},
		{name: "init failure", args: []string{"-e", `(throw (ex-info "init failed" {}))`, "-m", "app"}, want: "init failed", code: 1, contains: true},
		{name: "missing namespace", args: []string{"-m", "no-such-app"}, want: "unable to load namespace", code: 1, contains: true},
		{name: "missing entry", args: []string{"-m", "missing"}, want: "missing/-main not found", code: 1, contains: true},
		{name: "not callable", args: []string{"-m", "notfn"}, want: "notfn/-main is not callable", code: 1, contains: true},
		{name: "throws", args: []string{"-m", "boom"}, want: "entry failed", code: 1, contains: true},
		{name: "load throws", args: []string{"-m", "badload"}, want: "load failed", code: 1, contains: true},
		{name: "missing value", args: []string{"-m"}, want: "flag needs an argument", code: 2, contains: true},
		{name: "empty value", args: []string{"-m="}, want: "-m requires", code: 2, contains: true},
		{name: "malformed entry", args: []string{"-m", "app/"}, want: "-m requires", code: 2, contains: true},
		{name: "compile conflict", args: []string{"-c", filepath.Join(dir, "out.lgb"), "-m", "app"}, want: "-m cannot be combined", code: 2, contains: true},
		{name: "bundle conflict", args: []string{"-b", filepath.Join(dir, "out"), "-m", "app"}, want: "-m cannot be combined", code: 2, contains: true},
		{name: "wasm conflict", args: []string{"-w", filepath.Join(dir, "site"), "-m", "app"}, want: "-m cannot be combined", code: 2, contains: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"-source-paths", dir}, tt.args...)
			cmd := exec.Command(bin, args...)
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			code := 0
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok {
					code = exit.ExitCode()
				} else {
					t.Fatal(err)
				}
			}
			got := strings.TrimSpace(string(out))
			if code != tt.code {
				t.Fatalf("exit %d, want %d: %s", code, tt.code, got)
			}
			if tt.contains {
				if !strings.Contains(got, tt.want) {
					t.Fatalf("got %q, want containing %q", got, tt.want)
				}
				if tt.name == "init failure" && strings.Contains(got, ":load") {
					t.Fatalf("namespace loaded after init failure: %s", got)
				}
			} else if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
