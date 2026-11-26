package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestInternTwoArgsLeavesVarUnbound(t *testing.T) {
	ns := vm.NewNamespace("intern-unbound-test")

	got, err := CoreInternf(ns, vm.Symbol("x"))
	if err != nil {
		t.Fatal(err)
	}
	v, ok := got.(*vm.Var)
	if !ok {
		t.Fatalf("got %T, want *vm.Var", got)
	}
	if v.IsBound() {
		t.Error("var is bound after the two-argument intern")
	}

	if _, err := CoreInternf(ns, vm.Symbol("y"), vm.Int(1)); err != nil {
		t.Fatal(err)
	}
	if !ns.LookupLocal(vm.Symbol("y")).IsBound() {
		t.Error("var is unbound after the three-argument intern")
	}
}
