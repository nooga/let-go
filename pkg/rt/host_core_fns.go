/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// installCoreCompatFns adds a few clojure.core fns/forms that let-go was missing.
// Motivated by metosin/malli (which registers `indexed?` as a default schema and
// uses `class` in a countability check; borkdude.dynaload expands `locking` to
// monitor-enter/exit), but each is a general, standalone clojure.core addition.
func installCoreCompatFns(ns *vm.Namespace) {
	// indexed? — true iff the value is positionally indexed (vectors). Matches
	// Clojure's (instance? clojure.lang.Indexed x): strings are positionally
	// accessible in let-go (vm.Indexed) but are NOT clojure.lang.Indexed, so
	// exclude them.
	ns.Def("indexed?", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("indexed? expects 1 arg")
		}
		if _, isStr := vs[0].(vm.String); isStr {
			return vm.FALSE, nil
		}
		_, ok := vs[0].(vm.Indexed)
		return vm.Boolean(ok), nil
	}))

	// class — alias of `type`. Clojure's class/type differ on primitives, but
	// let-go has no boxed primitives so they coincide.
	ns.Def("class", ns.Lookup("type").(*vm.Var).Deref())

	// uri? — let-go has no java.net.URI type, so nothing is a URI.
	ns.Def("uri?", mustWrap(func(vs []vm.Value) (vm.Value, error) { return vm.FALSE, nil }))

	// monitor-enter / monitor-exit — no-ops; let-go has no object monitors, and
	// `locking` needs them to expand.
	noop := mustWrap(func(vs []vm.Value) (vm.Value, error) { return vm.NIL, nil })
	ns.Def("monitor-enter", noop)
	ns.Def("monitor-exit", noop)
}
