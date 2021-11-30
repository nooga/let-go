/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"fmt"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func runFile(filename string) error {
	ns := rt.NS(rt.NameCoreNS)
	if ns == nil {
		fmt.Println("namespace not found")
		return nil
	}
	ctx := compiler.NewCompiler(ns)
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
	file, err := os.Open("./")
	assert.NoError(t, err)
	names, err := file.Readdirnames(0)
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)
	outcomeVar := rt.CoreNS.Lookup("*test-flag*").(*vm.Var)
	for f := range names {
		fn := "./" + names[f]
		if filepath.Ext(fn) != ".lg" {
			continue
		}
		t.Run(names[f], func(t *testing.T) {
			outcomeVar.SetRoot(vm.TRUE)
			err := runFile(fn)
			assert.NoError(t, err)
			outcome := bool(outcomeVar.Deref().(vm.Boolean))
			assert.True(t, outcome, "some tests failed")
		})
	}
}
