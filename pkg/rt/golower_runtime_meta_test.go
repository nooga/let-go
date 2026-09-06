/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestApplyVarMetaAcceptsLegacyMap(t *testing.T) {
	v := vm.NewVar(vm.NewNamespace("test.legacy-meta"), "test.legacy-meta", "x")
	meta := vm.NewPersistentMap([]vm.Value{
		vm.Keyword("dynamic"), vm.TRUE,
		vm.Keyword("private"), vm.TRUE,
		vm.Keyword("doc"), vm.String("legacy"),
	})
	ApplyVarMeta(v, meta)
	if !v.IsDynamic() || !v.IsPrivate() {
		t.Fatal("legacy metadata map did not set var flags")
	}
	if got := v.Meta(); got != meta {
		t.Fatalf("legacy metadata = %v, want original map %v", got, meta)
	}
}

func TestApplyVarMetaDefersCompactPairsButSetsFlags(t *testing.T) {
	v := vm.NewVar(vm.NewNamespace("test.compact-meta"), "test.compact-meta", "x")
	pairs := vm.DefMetaPairs{
		vm.Keyword("dynamic"), vm.TRUE,
		vm.Keyword("private"), vm.TRUE,
		vm.Keyword("doc"), vm.String("compact"),
	}
	ApplyVarMeta(v, pairs)
	if !v.IsDynamic() || !v.IsPrivate() {
		t.Fatal("compact metadata did not set var flags before Meta")
	}
	meta, ok := v.Meta().(*vm.PersistentMap)
	if !ok {
		t.Fatalf("Meta() = %T, want *vm.PersistentMap", v.Meta())
	}
	if got := meta.ValueAt(vm.Keyword("doc")); got != vm.String("compact") {
		t.Fatalf(":doc = %v, want compact", got)
	}
}

func TestApplyVarMetaCompactPairsUseLastDuplicateFlag(t *testing.T) {
	v := vm.NewVar(vm.NewNamespace("test.duplicate-meta"), "test.duplicate-meta", "x")
	pairs := vm.DefMetaPairs{
		vm.Keyword("dynamic"), vm.TRUE,
		vm.Keyword("dynamic"), vm.FALSE,
		vm.Keyword("private"), vm.TRUE,
		vm.Keyword("private"), vm.FALSE,
	}
	ApplyVarMeta(v, pairs)
	if v.IsDynamic() || v.IsPrivate() {
		t.Fatal("an earlier duplicate overrode the final false metadata flag")
	}
	meta := v.Meta().(*vm.PersistentMap)
	if meta.ValueAt(vm.Keyword("dynamic")) != vm.FALSE || meta.ValueAt(vm.Keyword("private")) != vm.FALSE {
		t.Fatalf("materialized duplicate metadata = %v, want final false values", meta)
	}
}
