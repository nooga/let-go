package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// testGreaterThan5Pred builds a unary bytecode fn computing (< 5 x), so Some
// exercises the PreparedCall path with a real dispatch per element.
func testGreaterThan5Pred() *vm.Func {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST)
	chunk.Append32(consts.Intern(vm.Int(5)))
	chunk.Append(vm.OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(vm.OP_LT)
	chunk.Append(vm.OP_RETURN)
	chunk.SetMaxStack(4)
	return vm.MakeFunc(1, false, chunk)
}

func TestSomeBytecodePredOverChunkedAndLinearSeqs(t *testing.T) {
	pred := testGreaterThan5Pred()

	// Chunked source (Range) drives the chunk-walking arm.
	v, err := Some(vm.RootExecContext, pred, vm.NewRange(vm.Int(0), vm.Int(10), vm.Int(1)))
	if err != nil {
		t.Fatal(err)
	}
	if v != vm.TRUE {
		t.Fatalf("chunked walk: got %v", v)
	}

	// No match returns nil.
	v, err = Some(vm.RootExecContext, pred, vm.NewRange(vm.Int(0), vm.Int(5), vm.Int(1)))
	if err != nil {
		t.Fatal(err)
	}
	if v != vm.NIL {
		t.Fatalf("chunked no-match: got %v", v)
	}

	// Linear source (List) drives the First/Next arm.
	list, err := vm.ListType.Box([]vm.Value{vm.Int(1), vm.Int(9), vm.Int(2)})
	if err != nil {
		t.Fatal(err)
	}
	v, err = Some(vm.RootExecContext, pred, list)
	if err != nil {
		t.Fatal(err)
	}
	if v != vm.TRUE {
		t.Fatalf("linear walk: got %v", v)
	}
}

func TestSomeNativePredStillFallsBack(t *testing.T) {
	native, err := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
		return vm.Boolean(args[0] == vm.Int(2)), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	v, err := Some(vm.RootExecContext, native, vm.NewRange(vm.Int(0), vm.Int(4), vm.Int(1)))
	if err != nil {
		t.Fatal(err)
	}
	if v != vm.TRUE {
		t.Fatalf("native pred fallback: got %v", v)
	}
}
