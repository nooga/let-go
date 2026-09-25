/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func compileWarningSource(t *testing.T, source string) error {
	t.Helper()
	ctx := NewCompiler(vm.NewConsts(), rt.CoreNS).SetSource("reflection-warning.lg")
	_, err := ctx.Compile(source)
	return err
}

func TestReflectionWarningDefaultsOff(t *testing.T) {
	var output bytes.Buffer
	restore := rt.SetReflectionWarningWriter(&output)
	defer restore()
	rt.ResetReflectionWarnings()

	if err := compileWarningSource(t, "(fn [x]\n  (.missing x))"); err != nil {
		t.Fatal(err)
	}
	if output.Len() != 0 {
		t.Fatalf("default-off compilation warned: %s", output.String())
	}
}

func TestReflectionWarningUsesDynamicBindingAndDeduplicatesSourceSites(t *testing.T) {
	warnVar := rt.CoreNS.LookupLocal(vm.Symbol("*warn-on-reflection*"))
	if warnVar == nil {
		t.Fatal("*warn-on-reflection* var not defined")
	}
	old := warnVar.Root()
	warnVar.SetRoot(vm.FALSE)
	defer warnVar.SetRoot(old)
	vm.RootExecContext.PushBinding(warnVar, vm.TRUE)
	defer vm.RootExecContext.PopBinding(warnVar)

	var output bytes.Buffer
	restore := rt.SetReflectionWarningWriter(&output)
	defer restore()
	rt.ResetReflectionWarnings()
	const source = "(fn [x]\n  (.missing x)\n  (.other x))"
	if err := compileWarningSource(t, source); err != nil {
		t.Fatal(err)
	}
	if err := compileWarningSource(t, source); err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if count := strings.Count(got, "reflection warning"); count != 2 {
		t.Fatalf("warning count = %d, want once for each of two sites:\n%s", count, got)
	}
	for _, location := range []string{"reflection-warning.lg:2:", "reflection-warning.lg:3:"} {
		if !strings.Contains(got, location) {
			t.Errorf("warning missing precise source %s:\n%s", location, got)
		}
	}
}

func TestReflectionWarningIgnoresKnownHostTargetsAndGenericCalls(t *testing.T) {
	warnVar := rt.CoreNS.LookupLocal(vm.Symbol("*warn-on-reflection*"))
	if warnVar == nil {
		t.Fatal("*warn-on-reflection* var not defined")
	}
	old := warnVar.Root()
	warnVar.SetRoot(vm.TRUE)
	defer warnVar.SetRoot(old)

	var output bytes.Buffer
	restore := rt.SetReflectionWarningWriter(&output)
	defer restore()
	rt.ResetReflectionWarnings()
	if err := compileWarningSource(t, "(fn [f x]\n  (.missing \"known\")\n  (f x))"); err != nil {
		t.Fatal(err)
	}
	if output.Len() != 0 {
		t.Fatalf("known host target or ordinary IFn call warned: %s", output.String())
	}
}

func TestReflectionWarningHonoursHintedAndConstructorBoundLocals(t *testing.T) {
	warnVar := rt.CoreNS.LookupLocal(vm.Symbol("*warn-on-reflection*"))
	if warnVar == nil {
		t.Fatal("*warn-on-reflection* var not defined")
	}
	old := warnVar.Root()
	warnVar.SetRoot(vm.TRUE)
	defer warnVar.SetRoot(old)

	cases := []struct {
		name   string
		source string
		// want lists the line:column of each expected warning.
		want []string
	}{
		{"hinted param", "(fn [^String s]\n  (.length s))", nil},
		{"unhinted param", "(fn [s]\n  (.length s))", []string{"2:3"}},
		{"constructor-bound let", "(fn []\n  (let [sb (StringBuilder.)]\n    (.append sb \"x\")))", nil},
		{"hinted let", "(fn []\n  (let [^Object o (identity 1)]\n    (.toString o)))", nil},
		{"unhinted rebinding masks a hinted param", "(fn [^String s]\n  (let [s (identity s)]\n    (.length s)))", []string{"3:5"}},
		{"same-scope rebinding masks a known local", "(fn []\n  (let [s (StringBuilder.) s (identity s)]\n    (.length s)))", []string{"3:5"}},
		{"hinted param through a closure", "(fn [^String s]\n  (fn []\n    (.length s)))", nil},
		{"unhinted inner param shadows a hinted outer one", "(fn [^String s]\n  (fn [s]\n    (.length s)))", []string{"3:5"}},
		{"let bound to a known local", "(fn [^String s]\n  (let [t s]\n    (.length t)))", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			restore := rt.SetReflectionWarningWriter(&output)
			defer restore()
			rt.ResetReflectionWarnings()
			if err := compileWarningSource(t, tc.source); err != nil {
				t.Fatal(err)
			}
			got := output.String()
			if count := strings.Count(got, "reflection warning"); count != len(tc.want) {
				t.Fatalf("warning count = %d, want %d:\n%s", count, len(tc.want), got)
			}
			for _, location := range tc.want {
				if !strings.Contains(got, "reflection-warning.lg:"+location+":") {
					t.Errorf("warning missing at %s:\n%s", location, got)
				}
			}
		})
	}
}
