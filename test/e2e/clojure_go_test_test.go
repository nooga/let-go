/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClojureGoTestAdapter(t *testing.T) {
	if testing.Short() {
		t.Skip("generates and executes an AOT Go test package")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := repoRoot(t)
	lg := buildLG(t)
	dir := generateClojureGoTests(ctx, t, lg, root, "test/fixtures/go-test-adapter/framework.lg")
	for _, tc := range []struct {
		name, counts, diagnostic string
		fail                     bool
		line                     int
	}{
		{"passing", "pass=4 fail=0 error=0", "", false, 0},
		{"failing-first", "pass=1 fail=1 error=0", "intentional mismatch", true, 18},
		{"assertion-error", "pass=0 fail=0 error=1", "assertion exploded", true, 23},
		{"uncaught-error", "pass=0 fail=0 error=1", "body exploded", true, 26},
		{"empty-test", "pass=0 fail=0 error=0", "", false, 0},
		{"composed", "pass=4 fail=0 error=0", "", false, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := runCmd(ctx, t, "go", dir, []string{"test", "-count=1", "-v", "-run", "^TestClojure$/^" + tc.name + "$"})
			if (err != nil) != tc.fail {
				t.Fatalf("go test error = %v, want failure %v\n%s", err, tc.fail, out)
			}
			for _, want := range []string{"=== RUN   TestClojure/" + tc.name, tc.counts, tc.diagnostic} {
				if !strings.Contains(string(out), want) {
					t.Errorf("missing %q in output:\n%s", want, out)
				}
			}
			if strings.Count(string(out), "=== RUN   TestClojure/") != 1 {
				t.Errorf("-run did not select exactly one Clojure test:\n%s", out)
			}
			if tc.fail && !strings.Contains(string(out), "framework.lg:"+strconv.Itoa(tc.line)+":") {
				t.Errorf("failure lost source location:\n%s", out)
			}
		})
	}
	for _, kind := range []string{"once", "each"} {
		t.Run(kind+"-fixture-failure", func(t *testing.T) {
			fixture := filepath.Join(t.TempDir(), "fixture.lg")
			source := `(ns adapter.fixture (:require [clojure.test :refer [deftest is use-fixtures]]))
(use-fixtures :` + kind + ` (fn [f] (f) (is false "fixture teardown failed")))
(deftest passing (is true))`
			if err := os.WriteFile(fixture, []byte(source), 0o644); err != nil {
				t.Fatal(err)
			}
			dir := generateClojureGoTests(ctx, t, lg, root, fixture)
			out, err := runCmd(ctx, t, "go", dir, []string{"test", "-count=1", "-v"})
			if err == nil || !strings.Contains(string(out), "fixture teardown failed") || !strings.Contains(string(out), "--- FAIL: TestClojure") {
				t.Fatalf("fixture assertion must fail go test:\n%s\n%v", out, err)
			}
		})
	}
	t.Run("generated-body-executes", func(t *testing.T) {
		path := filepath.Join(dir, "bodies.go")
		source, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		mutant := wrapLastReturnResult(t, source, "LgTestBody0",
			`func() vm.Value { _ = __EXPR__; panic("lowered body sentinel") }()`)
		if err := os.WriteFile(path, mutant, 0o644); err != nil {
			t.Fatal(err)
		}
		if out, err := runCmd(ctx, t, "go", dir, []string{"test", "-c", "-o", filepath.Join(dir, "adapter.test")}); err != nil {
			t.Fatalf("mutant must compile: %v\n%s", err, out)
		}
		out, err := runCmd(ctx, t, "go", dir, []string{"test", "-v", "-count=1", "-run", "^TestClojure$/^passing$"})
		if err == nil || !strings.Contains(string(out), "lowered body sentinel") || !strings.Contains(string(out), "--- FAIL: TestClojure/passing") {
			t.Fatalf("generated body was not executed: %v\n%s", err, out)
		}
	})
	t.Run("reject-name-collision", func(t *testing.T) {
		fixture := filepath.Join(t.TempDir(), "collision.lg")
		source := `(ns adapter.collision (:require [clojure.test :refer [deftest is]]))
(defn lg-test-body-0 [] 42)
(deftest passing (is (= 42 (lg-test-body-0))))`
		if err := os.WriteFile(fixture, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
		out, err := runCmd(ctx, t, lg, root, []string{"scripts/lg-test-go", fixture, t.TempDir()})
		if err == nil || !strings.Contains(string(out), "generated test name collides") {
			t.Fatalf("expected collision rejection: %v\n%s", err, out)
		}
	})
	// Auto-resolved keywords must resolve in the fixture's namespace, as the
	// bundle compiled by `lg -c` resolves them; reading them as :user/key made
	// the lowered body disagree with the bundle's definitions.
	t.Run("namespace-keywords", func(t *testing.T) {
		fixture := filepath.Join(t.TempDir(), "keywords.lg")
		source := `(ns adapter.keywords (:require [clojure.test :refer [deftest is]] [clojure.string :as str]))
(def expected ::key)
(def aliased ::str/key)
(deftest keyword-test (is (= expected ::key)) (is (= aliased ::str/key)))`
		if err := os.WriteFile(fixture, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
		dir := generateClojureGoTests(ctx, t, lg, root, fixture)
		out, err := runCmd(ctx, t, "go", dir, []string{"test", "-count=1", "-v"})
		if err != nil || !strings.Contains(string(out), "pass=2 fail=0 error=0") {
			t.Fatalf("namespace keywords must resolve in the fixture namespace: %v\n%s", err, out)
		}
	})
	// An alias from a later top-level require must resolve too: the file is
	// read form by form, as `lg -c` reads it, not only after the ns form.
	t.Run("late-require-alias", func(t *testing.T) {
		fixture := filepath.Join(t.TempDir(), "late_alias.lg")
		source := `(ns adapter.late-alias (:require [clojure.test :refer [deftest is]]))
(require '[clojure.string :as str])
(def expected ::str/key)
(deftest keyword-test (is (= expected ::str/key)))`
		if err := os.WriteFile(fixture, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
		dir := generateClojureGoTests(ctx, t, lg, root, fixture)
		out, err := runCmd(ctx, t, "go", dir, []string{"test", "-count=1", "-v"})
		if err != nil || !strings.Contains(string(out), "pass=1 fail=0 error=0") {
			t.Fatalf("an alias from a later require must resolve: %v\n%s", err, out)
		}
	})
	// A test redefined under another form would run the stale extracted body
	// and report PASS for a test the ordinary runner fails.
	t.Run("reject-redefined-test", func(t *testing.T) {
		fixture := filepath.Join(t.TempDir(), "redefined.lg")
		source := `(ns adapter.redefined (:require [clojure.test :refer [deftest is]]))
(deftest check (is true))
(when true (deftest check (is false "replacement must fail")))`
		if err := os.WriteFile(fixture, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
		out, err := runCmd(ctx, t, lg, root, []string{"scripts/lg-test-go", fixture, t.TempDir()})
		if err == nil || !strings.Contains(string(out), "defined more than once: check") {
			t.Fatalf("expected generation to reject the redefined test `check`: %v\n%s", err, out)
		}
	})
	// A deftest the discovery walk cannot extract must fail generation: a
	// silently omitted test would let a failing suite report PASS.
	t.Run("reject-unextracted-test", func(t *testing.T) {
		fixture := filepath.Join(t.TempDir(), "discovery.lg")
		source := `(ns adapter.discovery (:require [clojure.test :refer [deftest is]]))
(deftest passing (is true))
(when true (deftest missed (is false "this must fail")))`
		if err := os.WriteFile(fixture, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
		out, err := runCmd(ctx, t, lg, root, []string{"scripts/lg-test-go", fixture, t.TempDir()})
		if err == nil || !strings.Contains(string(out), "missed") {
			t.Fatalf("expected generation to reject the unextracted test `missed`: %v\n%s", err, out)
		}
	})
}

func generateClojureGoTests(ctx context.Context, t *testing.T, lg, root, fixture string) string {
	t.Helper()
	dir := t.TempDir()
	if out, err := runCmd(ctx, t, lg, root, []string{
		"scripts/lg-test-go", fixture, dir,
	}); err != nil {
		t.Fatalf("generate adapter: %v\n%s", err, out)
	}
	mod := "module adaptertest\n\ngo " + rootGoDirective(t, root) + "\n\n" +
		"require github.com/nooga/let-go v0.0.0\n" +
		"replace github.com/nooga/let-go => " + strconv.Quote(root) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCmd(ctx, t, "go", dir, []string{"mod", "tidy"}); err != nil {
		t.Fatalf("tidy adapter: %v\n%s", err, out)
	}
	// Build separately: a compilation failure must never count as a detected assertion.
	if out, err := runCmd(ctx, t, "go", dir, []string{"test", "-c", "-o", filepath.Join(dir, "adapter.test")}); err != nil {
		t.Fatalf("compile adapter: %v\n%s", err, out)
	}
	return dir
}
