/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"encoding/json"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// js/emit encodes through internal/jsonenc so lg_no_json builds keep it. The
// payload bytes must stay what encoding/json produced before the switch.
func TestPrepareEmitMatchesEncodingJSON(t *testing.T) {
	kw := func(s string) vm.Value { return vm.Keyword(s) }
	m := func(kvs ...vm.Value) vm.Value {
		out := vm.EmptyPersistentMap
		for i := 0; i < len(kvs); i += 2 {
			out = out.Assoc(kvs[i], kvs[i+1]).(*vm.PersistentMap)
		}
		return out
	}
	payloads := []vm.Value{
		m(kw("hp"), vm.Int(12), kw("max-hp"), vm.Int(20), kw("depth"), vm.Int(3), kw("turn"), vm.Int(1234)),
		m(kw("title"), vm.String(`The "Crypt" <of> Ab&ol`), kw("quest"), vm.String("line\u2028sep\n")),
		m(kw("code"), vm.String("AbC-_09xyz"), kw("turn"), vm.Int(9007199254740993)),
		m(kw("ok"), vm.FALSE, kw("error"), vm.String("bad \\ input\t\x01")),
		m(kw("metrics"), m(kw("gen-ms"), vm.Float(62.302), kw("tiny"), vm.Float(1e-7), kw("big"), vm.Float(1e21))),
		vm.NewArrayVector([]vm.Value{vm.NIL, vm.TRUE, kw("xsofy/k"), vm.NewArrayVector(nil)}),
	}
	for _, p := range payloads {
		_, got, err := prepareEmit([]vm.Value{kw("stats"), p})
		if err != nil {
			t.Fatalf("prepareEmit(%s): %v", p, err)
		}
		goVal, err := fromValue(p)
		if err != nil {
			t.Fatal(err)
		}
		want, err := json.Marshal(goVal)
		if err != nil {
			t.Fatal(err)
		}
		if got != string(want) {
			t.Errorf("payload %s:\n got  %s\n want %s", p, got, want)
		}
	}
}
