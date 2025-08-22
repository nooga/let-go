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
	rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, []string{"."}))

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
			return nil
		}
		if filepath.Ext(path) != ".lg" {
			return nil
		}
		name := info.Name()
		t.Run(name, func(t *testing.T) {
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
		})
		return nil
	})
	assert.NoError(t, err)
}
