//go:build gogen_ir

package wasmhost

import (
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestHostEvalLoadsNativeIRPassesUnderTag(t *testing.T) {
	consts := vm.NewConsts()
	ctx := compiler.NewCompiler(consts, rt.NS("user"))
	NewResolver(ctx)
	host := New(consts)

	host.Handle(Request{
		ID:      "probe-native",
		Session: "test",
		Op:      "eval",
		Code:    "(require 'ir.passes.pipeline)\n:ok",
	}, func(string) {})

	v := rt.LookupVar("ir.passes.dce", "dce")
	if v == nil {
		t.Fatal("ir.passes.dce/dce var not found after require")
	}
	got := v.Deref()
	if name := got.Type().Name(); name != "let-go.lang.NativeFn" {
		t.Fatalf("ir.passes.dce/dce should be a NativeFn override under -tags gogen_ir; got %s", name)
	}
}
