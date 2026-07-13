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
