/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

var consts *vm.Consts

func runFile(filename string) error {
	ns := rt.NS(rt.NameCoreNS)
	if ns == nil {
		fmt.Println("namespace not found")
		return nil
	}
	ctx := compiler.NewCompiler(consts, ns)
	ctx.SetSource(filename)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, _, err = ctx.CompileMultiple(f)
	errc := f.Close()
	if err != nil {
		return err
	}
	if errc != nil {
		return errc
	}
	return nil
}

func TestRunner(t *testing.T) {
	consts = vm.NewConsts()
	// Set up a loader so rt.NS can autoload namespaces from files during tests.
	loaderCtx := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	// Search paths for `require`: current dir for in-tree test helpers
	// (test/test.lg etc.), plus examples/go-gen so tests can exercise
	// example libraries like gogen.
	rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, []string{".", "../examples/go-gen"}))

	// Per-file isolation baseline: a snapshot of the (clean) dynamic-binding
	// state taken before any test file runs. Each file is executed within this
	// scope and *ns* is reset below, so a file that leaves a dynamic var dirty
	// — e.g. *ns* left pointing at a scratch namespace by an in-ns or a
	// throwing (binding [*ns* ...] ...) body — cannot corrupt the unqualified
	// symbol resolution of files that run after it in this shared runtime.
	cleanBindings := vm.SnapshotBindings()
	coreNS := rt.NS(rt.NameCoreNS)

	file, err := os.Open("./")
	assert.NoError(t, err)
	// removed unused names := file.Readdirnames(0)
	err = file.Close()
	assert.NoError(t, err)

	// compile all .lg files first so tests are defined and registered
	// Run per file to retain Go subtest reporting and accurate counters
	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip directories that contain non-test .lg files (e.g.
			// compat/ holds the corpus runner, which is invoked manually;
			// benches/ holds standalone benchmarks + bug-repro fixtures that
			// run for seconds and assert nothing — run by hand with ./lg;
			// gogen/ holds native-lowering harness fixtures driven by the
			// deftype_skeleton_lowering_e2e_test, not bytecode deftests).
			if info.Name() == "compat" || info.Name() == "clojure-test-suite" || info.Name() == "benches" || info.Name() == "gogen" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".lg" {
			return nil
		}
		name := info.Name()
		t.Run(name, func(t *testing.T) {
			// Reset *ns* to a clean baseline (clojure.core) so a prior file's
			// leaked current-namespace doesn't carry over, then run the file's
			// compile+test cycle within the clean dynamic-binding scope so any
			// bindings it leaks are dropped afterward. Together these isolate
			// each file's runtime state from the next.
			rt.CurrentNS.SetRoot(coreNS)
			_, _ = vm.RunWithBindings(cleanBindings, func() (vm.Value, error) {
				// reset registry for per-file isolation
				_, _, cerr := compiler.NewCompiler(consts, rt.NS("test")).CompileMultiple(strings.NewReader("(clear-registered-tests!)"))
				assert.NoError(t, cerr)

				// compile the file to define tests
				cerr = runFile(path)
				assert.NoError(t, cerr)

				// run only this file's tests
				outcomeVar := rt.NS("test").Lookup("*test-result*").(*vm.Var)
				_, _, cerr = compiler.NewCompiler(consts, rt.NS("test")).CompileMultiple(strings.NewReader("(run-tests)"))
				assert.NoError(t, cerr)
				outcome := bool(outcomeVar.Deref().(vm.Boolean))
				assert.True(t, outcome, "some tests failed in "+name)
				return vm.NIL, nil
			})
		})
		return nil
	})
	assert.NoError(t, err)
}
