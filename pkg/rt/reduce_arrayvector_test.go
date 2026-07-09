package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Behavior tests for the ArrayVector direct-index reduce fast path.
func TestReduceArrayVectorFastPath(t *testing.T) {
	reduceVar := LookupCoreVar("reduce")
	if reduceVar == nil {
		t.Fatal("core/reduce not found")
	}
	reduceFn, ok := reduceVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("reduce is not an Fn")
	}

	plusVar := LookupCoreVar("+")
	if plusVar == nil {
		t.Fatal("core/+ not found")
	}
	plus, ok := plusVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("+ is not an Fn")
	}

	cases := []struct {
		name string
		args []vm.Value
		want vm.Value
	}{
		{
			"two-arg sum",
			[]vm.Value{plus, vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3), vm.Int(4), vm.Int(5)})},
			vm.Int(15),
		},
		{
			"three-arg sum with init",
			[]vm.Value{plus, vm.Int(100), vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3), vm.Int(4), vm.Int(5)})},
			vm.Int(115),
		},
		{
			"single element no init",
			[]vm.Value{plus, vm.NewArrayVector([]vm.Value{vm.Int(42)})},
			vm.Int(42),
		},
		{
			"empty with init",
			[]vm.Value{plus, vm.Int(7), vm.NewArrayVector([]vm.Value{})},
			vm.Int(7),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := reduceFn.Invoke(c.args)
			if err != nil {
				t.Fatalf("reduce error: %v", err)
			}
			if out != c.want {
				t.Fatalf("got %v want %v", out, c.want)
			}
		})
	}

	t.Run("reduced short-circuit", func(t *testing.T) {
		// f = (fn [a x] (if (> a 5) (reduced :stop) (+ a x))); reducing
		// [1 2 3 4 5 6 7] must stop at :stop without consuming the tail.
		f, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
			a := int64(vs[0].(vm.Int))
			if a > 5 {
				return vm.NewReduced(vm.Keyword("stop")), nil
			}
			return vm.Int(a + int64(vs[1].(vm.Int))), nil
		})
		if err != nil {
			t.Fatalf("wrap error: %v", err)
		}
		vec := vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3), vm.Int(4), vm.Int(5), vm.Int(6), vm.Int(7)})
		out, err := reduceFn.Invoke([]vm.Value{f, vec})
		if err != nil {
			t.Fatalf("reduce error: %v", err)
		}
		if out != vm.Keyword("stop") {
			t.Fatalf("got %v want :stop", out)
		}
	})
}
