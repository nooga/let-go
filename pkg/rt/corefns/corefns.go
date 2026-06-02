/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package corefns is a placeholder while the lginterop-driven direct-call
// path is being wired up. It currently exports just `Seq` so the dynamic
// native-direct registry plumbing (registry, pipeline-seeding,
// lower-go emit) has at least one entry to exercise.
//
// The long-term shape is: corefns hosts real Go top-level funcs with
// typed Go signatures (not the vm.Value-only shape used here). The
// metadata the registry holds is whatever lginterop's go/types pass
// extracts — typed params, typed results, variadic, &c. lower-go reads
// the typed signature and emits per-arg coercion + direct call, mirroring
// what vm.NativeFnType.Box already does reflectively at the boundary.
package corefns

import (
	"fmt"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Seq returns a seq of coll, or nil for empty/nil colls. Mirrors the
// closure body previously inlined in pkg/rt/lang.go around the `seq`
// registration.
func Seq(v vm.Value) (vm.Value, error) {
	if v == vm.NIL {
		return vm.NIL, nil
	}
	switch v.(type) {
	case *vm.Cons, *vm.LazySeq:
		// fall through
	default:
		if coll, ok := v.(vm.Collection); ok {
			if coll.RawCount() == 0 {
				return vm.NIL, nil
			}
		}
	}
	var n vm.Seq
	if sqbl, ok := v.(vm.Sequable); ok {
		n = sqbl.Seq()
	} else if s, ok := v.(vm.Seq); ok {
		n = s
	} else {
		return vm.NIL, fmt.Errorf("seq expected Seqable")
	}
	if n == nil || n == vm.EmptyList {
		return vm.NIL, nil
	}
	return n, nil
}

func init() {
	rt.RegisterNativeModule(&rt.NativeModule{
		GoPkg:     "github.com/nooga/let-go/pkg/rt/corefns",
		Namespace: "clojure.core",
		Fns: map[string]rt.NativeDirectFn{
			"seq": {
				GoIdent:    "Seq",
				Arity:      1,
				ParamSpecs: []string{"vm.Value"},
				ResultSpec: "vm.Value",
				NeedsError: true,
			},
		},
	})
}
