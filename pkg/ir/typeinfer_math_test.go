/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/genmanifest"
)

func TestTypeInferFloatReturningMathFnsMatchRuntimeRegistrations(t *testing.T) {
	root, err := genmanifest.FindRepoRoot(".")
	if err != nil {
		t.Fatalf("locate repo root: %v", err)
	}

	mathGo, err := os.ReadFile(filepath.Join(root, "pkg/rt/math.go"))
	if err != nil {
		t.Fatalf("read math.go: %v", err)
	}
	typeinferLg, err := os.ReadFile(filepath.Join(root, "pkg/rt/core/ir/passes/typeinfer.lg"))
	if err != nil {
		t.Fatalf("read typeinfer.lg: %v", err)
	}

	want := stringSet(regexp.MustCompile(`ns\.Def\("([^"]+)",\s*mathFn[12]\(`).FindAllStringSubmatch(string(mathGo), -1))

	blockRe := regexp.MustCompile(`(?s)\(def \^:private float-returning-math-fns\s+#\{(.*?)\}\)`)
	block := blockRe.FindStringSubmatch(string(typeinferLg))
	if block == nil {
		t.Fatalf("could not find float-returning-math-fns set in typeinfer.lg")
	}
	got := stringSet(regexp.MustCompile(`"([^"]+)"`).FindAllStringSubmatch(block[1], -1))

	if diff := compareStringSets(want, got); diff != "" {
		t.Fatalf("float-returning-math-fns drifted from mathFn1/mathFn2 registrations:\n%s", diff)
	}
}

func stringSet(matches [][]string) map[string]struct{} {
	out := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		out[match[1]] = struct{}{}
	}
	return out
}

func compareStringSets(want, got map[string]struct{}) string {
	var missing, extra []string
	for name := range want {
		if _, ok := got[name]; !ok {
			missing = append(missing, name)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			extra = append(extra, name)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)

	var b strings.Builder
	if len(missing) > 0 {
		b.WriteString("missing from typeinfer: ")
		b.WriteString(strings.Join(missing, ", "))
		b.WriteByte('\n')
	}
	if len(extra) > 0 {
		b.WriteString("extra in typeinfer: ")
		b.WriteString(strings.Join(extra, ", "))
		b.WriteByte('\n')
	}
	return b.String()
}
