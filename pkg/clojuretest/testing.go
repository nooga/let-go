/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

// Package clojuretest runs AOT-lowered clojure.test bodies as Go subtests.
package clojuretest

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Test connects a Clojure var name to its lowered Go body.
type Test struct {
	Name string
	File string
	Line int
	Body func(*vm.ExecContext) (vm.Value, error)
}

// Namespace initialization and test metadata share the runtime's global registry.
var suiteMu sync.Mutex

// Run loads an AOT bundle and runs tests with clojure.test's once/each fixtures.
// Test bodies replace :test metadata, so composed tests also use lowered bodies.
// Subtests are sequential; callers must not mutate the runtime concurrently.
func Run(t *testing.T, bundle []byte, namespace string, tests []Test) {
	t.Helper()
	suiteMu.Lock()
	defer suiteMu.Unlock()
	if err := rt.LoadCore(); err != nil {
		t.Fatalf("load core: %v", err)
	}
	rt.UseBytecodeNSLoader()
	if _, err := rt.RequireNS("test"); err != nil {
		t.Fatalf("load clojure.test: %v", err)
	}
	previousNS := rt.CurrentNS.Deref()
	rt.CurrentNS.SetRoot(rt.LookupOrRegisterNS(namespace))
	defer rt.CurrentNS.SetRoot(previousNS)
	ec := vm.NewExecContext()
	unit, err := rt.DecodeExecUnit(bundle)
	if err != nil {
		t.Fatalf("decode test bundle: %v", err)
	}
	if err := rt.LoadProgramNamespaces(unit); err != nil {
		t.Fatalf("load test namespaces: %v", err)
	}
	if err := rt.RunProgramMainChunk(unit); err != nil {
		t.Fatalf("initialize tests: %v", err)
	}
	ns := rt.LookupNS(namespace)
	if ns == nil {
		t.Fatalf("test namespace %q is missing", namespace)
	}
	rt.ApplyGoOverrides(ns)
	bind := func(name string, value vm.Value) {
		v := rt.LookupVar("clojure.test", name)
		ec.PushBinding(v, value)
		t.Cleanup(func() { ec.PopBinding(v) })
	}
	bind("*testing-vars*", vm.EmptyList)
	bind("*testing-contexts*", vm.EmptyList)
	bind("*report-counters*", vm.NewAtom(vm.Map{
		vm.Keyword("test"): vm.Int(0), vm.Keyword("pass"): vm.Int(0),
		vm.Keyword("fail"): vm.Int(0), vm.Keyword("error"): vm.Int(0),
	}))
	nsVar := rt.LookupVar("clojure.core", "*ns*")
	ec.PushBinding(nsVar, ns)
	defer ec.PopBinding(nsVar)
	vars := make([]*vm.Var, len(tests))
	for i, test := range tests {
		v := ns.LookupLocal(vm.Symbol(test.Name))
		if v == nil || test.Body == nil {
			t.Fatalf("missing test var or lowered body: %s/%s", namespace, test.Name)
		}
		meta := v.Meta()
		m, ok := meta.(vm.Associative)
		if !ok || field(meta, "test") == vm.NIL {
			t.Fatalf("%s/%s has no :test metadata", namespace, test.Name)
		}
		body := vm.NewArityNativeFn(test.Name, 0, false, func(ec *vm.ExecContext, _ []vm.Value) (vm.Value, error) {
			return test.Body(ec)
		})
		v.SetMeta(m.Assoc(vm.Keyword("test"), body).
			Assoc(vm.Keyword("file"), vm.String(test.File)).
			Assoc(vm.Keyword("line"), vm.Int(test.Line)))
		defer v.SetMeta(meta)
		vars[i] = v
	}
	call := func(name string, args ...vm.Value) (vm.Value, error) {
		fn, ok := ec.Deref(rt.LookupVar("clojure.test", name)).(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("clojure.test/%s is not callable", name)
		}
		return ec.Invoke(fn, args)
	}
	once, err := call("join-fixtures", field(ns.Meta(), "test/once-fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	each, err := call("join-fixtures", field(ns.Meta(), "test/each-fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	_, restoreReports := captureReports(t, ec, Test{})
	defer restoreReports()
	run := vm.NewArityNativeFn("go-test-suite", 0, false, func(_ *vm.ExecContext, _ []vm.Value) (vm.Value, error) {
		for i, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				counts, restore := captureReports(t, ec, test)
				defer restore()
				body := vm.NewArityNativeFn(test.Name, 0, false, func(_ *vm.ExecContext, _ []vm.Value) (vm.Value, error) {
					return call("test-var", vars[i])
				})
				if _, err := ec.Invoke(each.(vm.Fn), []vm.Value{body}); err != nil {
					t.Errorf("test or each fixture: %v", err)
				}
				t.Logf("assertions: pass=%d fail=%d error=%d", counts["pass"], counts["fail"], counts["error"])
			})
		}
		return vm.NIL, nil
	})
	if _, err := ec.Invoke(once.(vm.Fn), []vm.Value{run}); err != nil {
		t.Errorf("once fixture: %v", err)
	}
}

// A suite-level reporter catches once-fixture assertions; each subtest shadows
// it while its each-fixture and body run, then restores it for teardown.
func captureReports(t *testing.T, ec *vm.ExecContext, location Test) (map[string]int, func()) {
	counts := map[string]int{}
	report := vm.NewArityNativeFn("go-test-report", 1, false, func(_ *vm.ExecContext, args []vm.Value) (vm.Value, error) {
		event := args[0]
		kind := field(event, "type")
		switch kind {
		case vm.Keyword("pass"), vm.Keyword("fail"), vm.Keyword("error"):
			counts[string(kind.(vm.Keyword))]++
			fn := ec.Deref(rt.LookupVar("clojure.test", "inc-report-counter")).(vm.Fn)
			if _, err := ec.Invoke(fn, []vm.Value{kind}); err != nil {
				return vm.NIL, err
			}
		}
		if kind == vm.Keyword("fail") || kind == vm.Keyword("error") {
			file, line := field(event, "file"), field(event, "line")
			if file == vm.NIL {
				file = vm.String(location.File)
			}
			if line == vm.NIL {
				line = vm.Int(location.Line)
			}
			contexts := ec.Deref(rt.LookupVar("clojure.test", "*testing-contexts*"))
			fileText := fmt.Sprint(file)
			if s, ok := file.(vm.String); ok {
				fileText = string(s)
			}
			t.Errorf("%s:%v: %v: %v\ncontext: %v\nexpected: %v\nactual: %v",
				fileText, line, kind, field(event, "message"), contexts, field(event, "expected"), field(event, "actual"))
		}
		return vm.NIL, nil
	})
	v := rt.LookupVar("clojure.test", "report")
	ec.PushBinding(v, report)
	return counts, func() { ec.PopBinding(v) }
}

func field(value vm.Value, key string) vm.Value {
	if m, ok := value.(vm.Lookup); ok {
		return m.ValueAt(vm.Keyword(key))
	}
	return vm.NIL
}
