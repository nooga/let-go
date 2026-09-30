/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// installExceptionClasses registers the java.lang.* exception hierarchy:
// each class under its qualified and bare name (Clojure auto-imports
// java.lang), plus the interop constructor spellings X., java.lang.X.,
// ->X, and ->java.lang.X. Constructors accept [], [msg], and [msg cause]
// and return a class-tagged exception value. ExceptionInfo aliases the
// existing ExInfoType, whose ancestry chains through RuntimeException in
// directTypeParents.
func installExceptionClasses(ns *vm.Namespace) {
	for _, class := range vm.ExceptionClasses {
		qualified := class.Name()
		bare := qualified[strings.LastIndex(qualified, ".")+1:]

		ctor := exceptionConstructor(class)
		ns.Def(qualified, class)
		ns.Def(bare, class)
		for _, spelling := range []string{
			bare + ".", qualified + ".", "->" + bare, "->" + qualified,
		} {
			ns.Def(spelling, ctor)
		}
	}

	ns.Def("ExceptionInfo", vm.ExInfoType)
	ns.Def("clojure.lang.ExceptionInfo", vm.ExInfoType)

	// Cancelled resolves for (catch Cancelled e ...) dispatch (see
	// catch-matches? below and vm.ClassCancelled's doc comment), but — unlike
	// every class in the loop above — gets no exceptionConstructor: nothing
	// in Lisp can construct a value of this class, which is what makes a
	// caught Cancelled trustworthy as a real scope-cancellation signal.
	ns.Def("Cancelled", vm.ClassCancelled)

	// catch-matches? [class-symbol caught] backs typed catch dispatch. Both
	// compilers desugar (catch SomeClass e ...) into a test through it, so
	// the semantics live in one place:
	//   - the symbol resolves at dispatch time; a JVM-only class let-go does
	//     not model yields a clause that never matches, instead of failing
	//     the whole namespace at compile time;
	//   - Throwable matches ANY thrown value: let-go permits throwing plain
	//     values (strings), and (catch Throwable e ...) is the conventional
	//     catch-everything. instance? itself stays honest — a thrown string
	//     is not (instance? Throwable s);
	//   - EXCEPT the scope-cancellation condition (vm.ClassCancelled, see
	//     pkg/vm/cancelled.go): Throwable and Exception do not match it, so
	//     an idiomatic (catch Throwable e (log e) (recur)) worker loop can't
	//     swallow its own termination. Only (catch Cancelled e ...) —
	//     naming it explicitly — matches, the same way catching
	//     java.lang.Exception in Java never catches an Error. Both compilers
	//     desugar a bare (catch e ...) into a Throwable-classed clause (see
	//     pkg/compiler/compiler.go tryCompiler and
	//     pkg/rt/core/ir/build.lg desugar-catch-clauses), so this single
	//     check also covers the bare form;
	//   - any other class matches by type identity or registered ancestry,
	//     like instance?. Cancelled's own ancestry (Throwable, Any) never
	//     includes Exception, so (catch Exception e ...) already misses it
	//     without a special case.
	catchMatches, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("catch-matches? expects 2 args")
		}
		sym, ok := vs[0].(vm.Symbol)
		if !ok {
			return vm.FALSE, nil
		}
		v, lok := ns.Lookup(sym).(*vm.Var)
		if !lok {
			return vm.FALSE, nil
		}
		class, cok := v.Deref().(vm.ValueType)
		if !cok {
			return vm.FALSE, nil
		}
		if class == vm.ValueType(vm.ClassThrowable) {
			if vm.IsCancelled(vs[1]) {
				return vm.FALSE, nil
			}
			return vm.TRUE, nil
		}
		if vs[1].Type() == class {
			return vm.TRUE, nil
		}
		if anc := directTypeAncestors(vs[1].Type()); anc != nil {
			return anc.Contains(class), nil
		}
		return vm.FALSE, nil
	})
	if err != nil {
		panic(err)
	}
	ns.Def("catch-matches?", catchMatches)
}

func exceptionConstructor(class *vm.ExceptionClass) vm.Value {
	ctor, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) > 2 {
			return vm.NIL, fmt.Errorf("%s. expects 0 to 2 args, got %d", class.Name(), len(vs))
		}
		message := ""
		if len(vs) >= 1 {
			s, ok := vs[0].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("%s. expected a String message, got %s", class.Name(), vs[0].Type().Name())
			}
			message = string(s)
		}
		// Mirror ex-info's cause handling: keep exception causes so
		// ex-cause/.getCause work, silently drop anything else.
		var cause error
		if len(vs) == 2 {
			if ei, ok := vs[1].(*vm.ExInfo); ok {
				cause = ei
			}
		}
		return vm.NewExInfoWithClass(class, message, cause), nil
	})
	if err != nil {
		panic(err)
	}
	return ctor
}
