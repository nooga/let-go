/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"github.com/zeebo/xxh3"

	"github.com/nooga/let-go/pkg/vm"
)

func mustBox(name string, fn any) vm.Value {
	v, err := vm.NativeFnType.Box(fn)
	if err != nil {
		panic("hash NS: " + name + ": " + err.Error())
	}
	return v
}

func makeXxh3Hasher(h *xxh3.Hasher) vm.Value {
	write := mustBox("xxh3-hasher/write", func(b []byte) {
		_, _ = h.Write(b)
	})
	writeStr := mustBox("xxh3-hasher/write-str", func(s string) {
		_, _ = h.WriteString(s)
	})
	sum := mustBox("xxh3-hasher/sum", h.Sum64)
	reset := mustBox("xxh3-hasher/reset", h.Reset)

	m := vm.EmptyPersistentMap.
		Assoc(vm.Keyword("write"), write).(*vm.PersistentMap).
		Assoc(vm.Keyword("write-str"), writeStr).(*vm.PersistentMap).
		Assoc(vm.Keyword("sum"), sum).(*vm.PersistentMap).
		Assoc(vm.Keyword("reset"), reset).(*vm.PersistentMap)
	return m
}

// nolint
func installHashNS() {
	ns := vm.NewNamespace("hash")

	ns.Def("xxh3-64", mustBox("xxh3-64", xxh3.Hash))
	ns.Def("xxh3-64-seed", mustBox("xxh3-64-seed", xxh3.HashSeed))
	ns.Def("xxh3-64-str", mustBox("xxh3-64-str", xxh3.HashString))
	ns.Def("xxh3-64-str-seed", mustBox("xxh3-64-str-seed", xxh3.HashStringSeed))

	ns.Def("xxh3-hasher", mustBox("xxh3-hasher", func() vm.Value {
		return makeXxh3Hasher(xxh3.New())
	}))
	ns.Def("xxh3-hasher-seed", mustBox("xxh3-hasher-seed", func(seed uint64) vm.Value {
		return makeXxh3Hasher(xxh3.NewSeed(seed))
	}))

	RegisterNS(ns)
}
