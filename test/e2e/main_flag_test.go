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
		{name: "missing entry", args: []string{"-m", "missing"}, want: "-m: missing/-main not found", code: 1, contains: true},
		{name: "not callable", args: []string{"-m", "notfn"}, want: "-m: notfn/-main is not callable", code: 1, contains: true},
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

// -m keys persistent storage by namespace, the way script mode keys it by file,
// so a program keeps its store whichever way it is launched.
func TestMainFlagStorageID(t *testing.T) {
	bin := buildLG(t)
	dir := t.TempDir()
	home := t.TempDir()
	cases := []struct{ file, ns string }{
		{"store.lg", "store"},
		{"my_app.clj", "my-app"},
		{"nested/cljc_app.cljc", "nested.cljc-app"},
		{"literal-app.lg", "literal-app"},
		{"choice_app.lg", "choice-app"},
	}
	for _, tc := range cases {
		file := filepath.Join(dir, tc.file)
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatal(err)
		}
		body := "(ns " + tc.ns + `) (def saved (storage/get "k")) (defn -main [] (prn saved))`
		if err := os.WriteFile(file, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// The .lg underscore candidate takes precedence over the .clj hyphen candidate.
	if err := os.WriteFile(filepath.Join(dir, "choice-app.clj"), []byte(`(ns choice-app) (defn -main [] (prn :wrong-file))`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "other.lg"), []byte(`(ns other) (defn -main [] (prn (storage/get "k")))`), 0644); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) string {
		t.Helper()
		cmd := exec.Command(bin, args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "HOME="+home, "XDG_CONFIG_HOME="+home, "APPDATA="+home)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%v: %v: %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	for _, tc := range cases {
		t.Run(tc.ns, func(t *testing.T) {
			run("-e", `(storage/set "k" "from-script")`, tc.file)
			if got := run("-source-paths", dir, "-m", tc.ns); got != `"from-script"` {
				t.Fatalf("-m %s read %s, want the script-mode value", tc.ns, got)
			}
		})
	}
	run("-storage-id", "custom", "-e", `(storage/set "k" "explicit")`)
	if got := run("-storage-id", "custom", "-source-paths", dir, "-m", "my-app/-main"); got != `"explicit"` {
		t.Fatalf("explicit store read %s, want explicit", got)
	}
	if got := run("-source-paths", dir, "-m", "other"); got != "nil" {
		t.Fatalf("-m other read %s, want its own empty store", got)
	}
}
