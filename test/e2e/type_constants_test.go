/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Each source quotes a defrecord or deftype instance into the compiled code,
// so it reaches the .lgb as a constant. The type is defined in the same file,
// so a loaded bundle decodes the constant before the defrecord or deftype form
// runs.
var typeConstantsCases = []struct{ name, src, want string }{
	{"defrecord", `(ns app)
(defprotocol Area (area [s]))
(defrecord Pt [x y] Area (area [_] (* x y)))
(defmacro pt-const [] (list 'quote (->Pt 2 3)))
(when-not *compiling-aot*
  (let [p (pt-const)]
    (prn (instance? Pt p) p (= p (->Pt 2 3)) (area p))))
`, "true #Pt{:x 2, :y 3} true 6\n"},
	{"deftype", `(ns app)
(defprotocol Area (area [s]))
(deftype Box [v] Area (area [_] v))
(defmacro box-const [] (list 'quote (Box. 7)))
(when-not *compiling-aot*
  (let [b (box-const)]
    (prn (instance? Box b) (= Box (type b)) (area b))))
`, "true true 7\n"},
}

// TestTypeConstantsSurviveBytecode checks that a defrecord or deftype value
// compiled in as a constant decodes with the live type, so instance?, =,
// printing and protocol dispatch match the uncompiled run under lg and
// lg-runtime.
func TestTypeConstantsSurviveBytecode(t *testing.T) {
	full := buildLG(t)
	runtime := buildLGRuntime(t)

	run := func(t *testing.T, bin string, args ...string) string {
		t.Helper()
		out, err := exec.Command(bin, args...).CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", filepath.Base(bin), args, err, out)
		}
		return string(out)
	}

	for _, tc := range typeConstantsCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "app.lg")
			lgb := filepath.Join(dir, "app.lgb")
			if err := os.WriteFile(src, []byte(tc.src), runtimeOnlyTestFilePerm); err != nil {
				t.Fatal(err)
			}
			if got := run(t, full, src); got != tc.want {
				t.Fatalf("source run: got %q, want %q", got, tc.want)
			}
			run(t, full, "-c", lgb, src)
			for _, bin := range []string{full, runtime} {
				if got := run(t, bin, lgb); got != tc.want {
					t.Errorf("%s app.lgb: got %q, want %q", filepath.Base(bin), got, tc.want)
				}
			}
		})
	}
}
