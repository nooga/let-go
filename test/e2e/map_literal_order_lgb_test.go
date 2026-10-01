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

// A compiled .lgb keeps each collection literal its own constant: two map
// literals of equal content in different key orders print their own order
// after the bundle's constant pool is written and read back.
func TestMapLiteralKeyOrderSurvivesLgb(t *testing.T) {
	lg := buildLG(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "order.lg")
	prog := `(ns order)
(defn ab [] {:a 1 :b 2})
(defn ba [] {:b 2 :a 1})
(println (vec (keys (ab))) (vec (keys (ba))) (identical? [1 2] [1 2]))
`
	if err := os.WriteFile(src, []byte(prog), 0o644); err != nil {
		t.Fatal(err)
	}
	want := "[:a :b] [:b :a] false"
	direct, err := exec.Command(lg, src).CombinedOutput()
	if err != nil || strings.TrimSpace(string(direct)) != want {
		t.Fatalf("source run: %q, %v; want %q", direct, err, want)
	}
	lgb := filepath.Join(dir, "order.lgb")
	if out, err := exec.Command(lg, "-c", lgb, src).CombinedOutput(); err != nil {
		t.Fatalf("lg -c: %v\n%s", err, out)
	}
	got, err := exec.Command(lg, lgb).CombinedOutput()
	if err != nil || strings.TrimSpace(string(got)) != want {
		t.Fatalf("bundle run: %q, %v; want %q", got, err, want)
	}
}
