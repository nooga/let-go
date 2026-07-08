/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bytes"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// The lazy (bytes) decode path defers source-map materialization to first
// Lookup. This must be observationally identical to the eager (io.Reader) path:
// for every chunk and every IP, LookupSource must return the same source
// location. Decoding the REAL core bundle both ways and comparing every lookup
// is un-gameable — a mis-captured byte span or a wrong string-ref would diverge.
func TestDecodeBytesLazySourceMapMatchesEager(t *testing.T) {
	if len(rt.CoreCompiledLGB) == 0 {
		t.Skip("no precompiled core_compiled.lgb")
	}
	resolve := func(ns, name string) *vm.Var {
		n := rt.NS(ns)
		v := n.Lookup(vm.Symbol(name))
		if v == vm.NIL {
			return n.Def(name, vm.NIL)
		}
		return v.(*vm.Var)
	}

	eager, err := bytecode.DecodeToExecUnit(bytes.NewReader(rt.CoreCompiledLGB), resolve)
	if err != nil {
		t.Fatalf("eager decode: %v", err)
	}
	lazy, err := bytecode.DecodeToExecUnitBytes(rt.CoreCompiledLGB, resolve)
	if err != nil {
		t.Fatalf("lazy decode: %v", err)
	}

	if len(eager.NSChunks) != len(lazy.NSChunks) {
		t.Fatalf("NSChunks count: eager=%d lazy=%d", len(eager.NSChunks), len(lazy.NSChunks))
	}

	sawSource := false
	for name, ec := range eager.NSChunks {
		lc := lazy.NSChunks[name]
		if lc == nil {
			t.Fatalf("lazy missing ns chunk %q", name)
		}
		n := ec.Length()
		if lc.Length() != n {
			t.Fatalf("ns %q length: eager=%d lazy=%d", name, n, lc.Length())
		}
		for ip := 0; ip <= n; ip++ {
			es, ls := ec.LookupSource(ip), lc.LookupSource(ip)
			if (es == nil) != (ls == nil) {
				t.Fatalf("ns %q ip=%d: eager=%v lazy=%v (nil mismatch)", name, ip, es, ls)
			}
			if es != nil {
				sawSource = true
				if *es != *ls {
					t.Fatalf("ns %q ip=%d: eager=%+v lazy=%+v", name, ip, *es, *ls)
				}
			}
		}
	}
	if !sawSource {
		t.Skip("core bundle carries no source maps; equivalence check vacuous")
	}
}
