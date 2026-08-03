package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestBindGeneratedPrimitiveBindsIntoNamespace(t *testing.T) {
	adapter, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) { return vm.NIL, nil })
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}
	BindGeneratedPrimitive("bind.probe.ns", "probe-fn", adapter)
	ns := LookupNS("bind.probe.ns")
	if ns == nil {
		t.Fatal("namespace not created")
	}
	v := ns.LookupLocal(vm.Symbol("probe-fn"))
	if v == nil || !v.IsBound() {
		t.Fatal("primitive not bound in its namespace")
	}
}
