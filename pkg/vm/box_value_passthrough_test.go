package vm_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestBoxValueArgPassthrough(t *testing.T) {
	// Box a Go fn whose param is vm.Value. Call it with an Int.
	// Before fix: reflect.Call panicked because Int.Unbox() → int64.
	fn, err := vm.NativeFnType.Box(func(v vm.Value) vm.Value {
		return v
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}
	got, err := fn.(*vm.NativeFn).Invoke([]vm.Value{vm.Int(42)})
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if i, ok := got.(vm.Int); !ok || int64(i) != 42 {
		t.Fatalf("expected Int(42), got %v (%T)", got, got)
	}
}
