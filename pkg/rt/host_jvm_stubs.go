/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// This file provides compile-only stubs for JVM-only surfaces that a Clojure
// library references on its :clj branch but that let-go cannot meaningfully run.
// The forms must RESOLVE (so the namespace compiles) but fail loudly / degrade if
// actually executed. Motivated by metosin/malli's java.time date coercion,
// bounded-execution timeout helper, and borkdude.dynaload's LazyVar deftype.

// loudStub returns a fn that resolves at compile time but throws if ever called.
func loudStub(what string) vm.Value {
	return mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NIL, fmt.Errorf("%s is not supported under let-go", what)
	})
}

// theHostDateStubType is the ValueType for the java.time chainStub.
type theHostDateStubType struct{}

func (t *theHostDateStubType) String() string     { return t.Name() }
func (t *theHostDateStubType) Type() vm.ValueType { return vm.TypeType }
func (t *theHostDateStubType) Unbox() any         { return nil }
func (t *theHostDateStubType) Name() string       { return "java.time.stub" }
func (t *theHostDateStubType) Box(any) (vm.Value, error) {
	return vm.NIL, nil
}

var hostDateStubType = &theHostDateStubType{}

// chainStub backs the java.time formatter objects a library builds at LOAD time
// via (-> (DateTimeFormatterBuilder.) (.appendPattern ..) .. (.toFormatter)).
// Builder methods return self so the chain threads; every other method (runtime
// .parse/.format, a typo, or a new method) fails loudly.
type chainStub struct{}

func (c *chainStub) Type() vm.ValueType { return hostDateStubType }
func (c *chainStub) Unbox() any         { return c }
func (c *chainStub) String() string     { return "#<java.time.stub>" }

var chainStubBuilderMethods = map[string]bool{
	"appendPattern": true, "optionalStart": true, "optionalEnd": true,
	"appendFraction": true, "appendOffset": true, "parseDefaulting": true,
	"toFormatter": true, "withZone": true,
}

func (c *chainStub) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	if chainStubBuilderMethods[string(name)] {
		return c, nil
	}
	return vm.NIL, fmt.Errorf("java.time .%s is not supported under let-go (date coercion degraded)", name)
}

// installJVMStubs registers the compile-only stubs (depends on defStaticNS from
// host_jvm_statics.go).
func installJVMStubs(ns *vm.Namespace) {
	// clojure.lang.IDeref / IFn as Protocols so a deftype implementing them (e.g.
	// borkdude.dynaload's LazyVar) compiles. The generated methods are never
	// dispatched (let-go uses -deref/-invoke), so such a value fails loudly if used.
	ns.Def("clojure.lang.IDeref", vm.NewProtocol("clojure.lang.IDeref", []vm.Symbol{"deref"}))
	ns.Def("clojure.lang.IFn", vm.NewProtocol("clojure.lang.IFn", []vm.Symbol{"invoke", "applyTo"}))

	// Bounded-execution / threads: load-only loud stubs.
	futureTask := loudStub("java.util.concurrent.FutureTask")
	ns.Def("FutureTask.", futureTask)
	ns.Def("->FutureTask", futureTask)
	thread := loudStub("java.lang.Thread")
	ns.Def("Thread.", thread)
	ns.Def("->Thread", thread)
	for _, nm := range []string{"TimeUnit", "java.util.concurrent.TimeUnit"} {
		defStaticNS(nm).Def("MILLISECONDS", vm.Symbol("TimeUnit/MILLISECONDS"))
	}

	// java.time builder/formatter forms -> chainable stub so the load-time
	// formatter construction threads; .parse/.format throw at runtime (degraded).
	dtb := mustWrap(func(vs []vm.Value) (vm.Value, error) { return &chainStub{}, nil })
	ns.Def("DateTimeFormatterBuilder.", dtb)
	ns.Def("->DateTimeFormatterBuilder", dtb)
	for _, nm := range []string{"DateTimeFormatter", "ZoneId", "Date", "Instant"} {
		nsb := defStaticNS(nm)
		nsb.Def("ofPattern", dtb)
		nsb.Def("of", dtb)
		nsb.Def("from", loudStub(nm+" date coercion"))
		nsb.Def("ofEpochMilli", loudStub(nm+" date coercion"))
	}
	chrono := defStaticNS("ChronoField")
	for _, f := range []string{"MICRO_OF_SECOND", "HOUR_OF_DAY", "OFFSET_SECONDS"} {
		chrono.Def(f, vm.Symbol("ChronoField/"+f))
	}

	// Decimal / URI string coercion — loud stubs (a caller's try/catch usually
	// degrades these to a pass-through).
	bigDec := loudStub("java.math.BigDecimal string coercion")
	ns.Def("BigDecimal.", bigDec)
	ns.Def("->BigDecimal", bigDec)
	uri := loudStub("java.net.URI string coercion")
	ns.Def("URI.", uri)
	ns.Def("->URI", uri)
}
