/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// theHostThreadType is the ValueType for the java.lang.Thread compat shim.
// let-go has no threads, but a library's cooperative-cancellation idiom
// ((.isInterrupted (Thread/currentThread)) -> (throw (InterruptedException.)))
// maps cleanly onto the scope tree, so the shim gives it real meaning instead
// of a load-only stub. Motivated by weavejester/ragtime, whose migrate-all
// checks for interruption between migrations.
type theHostThreadType struct{}

func (t *theHostThreadType) String() string     { return t.Name() }
func (t *theHostThreadType) Type() vm.ValueType { return vm.TypeType }
func (t *theHostThreadType) Unbox() any         { return nil }
func (t *theHostThreadType) Name() string       { return "java.lang.Thread" }
func (t *theHostThreadType) Box(any) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("java.lang.Thread cannot be boxed")
}

// hostThreadType is the singleton type for hostThread values.
var hostThreadType = &theHostThreadType{}

// hostThread is what Thread/currentThread returns: the calling scope,
// dressed as a java.lang.Thread. Interruption is scope cancellation.
type hostThread struct{ scope *vm.Scope }

func (h *hostThread) Type() vm.ValueType { return hostThreadType }
func (h *hostThread) Unbox() any         { return h }
func (h *hostThread) String() string {
	return fmt.Sprintf("#<java.lang.Thread scope=%d>", h.scope.Live())
}

func (h *hostThread) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	switch string(name) {
	case "isInterrupted":
		if len(args) == 0 {
			// The same test scope-cancelled? makes.
			return vm.Boolean(h.scope.Context().Err() != nil), nil
		}
	case "interrupt":
		if len(args) == 0 {
			// Cancel is terminal for the scope (it does not install a fresh
			// context generation the way CancelAll does), which is what keeps
			// isInterrupted true afterwards. Inside with-scope that is bounded;
			// at the root scope it cancels every tracked goroutine's context for
			// the rest of the process — what interrupting the main thread means.
			h.scope.Cancel()
			return vm.NIL, nil
		}
	}
	return vm.NIL, fmt.Errorf("java.lang.Thread .%s is not supported under let-go", name)
}

// installHostThread registers Thread/currentThread. Thread. itself stays a
// loud stub (host_jvm_stubs.go): let-go cannot construct threads.
func installHostThread() {
	defStaticNS("Thread").Def("currentThread", vm.NewCtxNativeFn("Thread/currentThread",
		func(ec *vm.ExecContext, args []vm.Value) (vm.Value, error) {
			if len(args) != 0 {
				return vm.NIL, fmt.Errorf("Thread/currentThread expects no arguments, got %d", len(args))
			}
			return &hostThread{scope: ec.Scope()}, nil
		}))
}
