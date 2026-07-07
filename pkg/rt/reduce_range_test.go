package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Behavior tests for the Range direct-arithmetic reduce fast path.
func TestReduceRangeFastPath(t *testing.T) {
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
			"(reduce + (range 10)) = 45",
			[]vm.Value{plus, vm.NewRange(vm.Int(0), vm.Int(10), vm.Int(1))},
			vm.Int(45),
		},
		{
			"(reduce + 100 (range 10)) = 145",
			[]vm.Value{plus, vm.Int(100), vm.NewRange(vm.Int(0), vm.Int(10), vm.Int(1))},
			vm.Int(145),
		},
		{
			"(reduce + (range 5 10)) = 35",
			[]vm.Value{plus, vm.NewRange(vm.Int(5), vm.Int(10), vm.Int(1))},
			vm.Int(35),
		},
		{
			"(reduce + (range 10 0 -2)) = 30",
			[]vm.Value{plus, vm.NewRange(vm.Int(10), vm.Int(0), vm.Int(-2))},
			vm.Int(30), // 10+8+6+4+2
		},
		{
			"(reduce + (range 0)) = 0",
			[]vm.Value{plus, vm.NewRange(vm.Int(0), vm.Int(0), vm.Int(1))},
			vm.Int(0), // empty, no init → (+) = 0
		},
		{
			"(reduce + 7 (range 0)) = 7",
			[]vm.Value{plus, vm.Int(7), vm.NewRange(vm.Int(0), vm.Int(0), vm.Int(1))},
			vm.Int(7), // empty with init
		},
		{
			"(reduce + (range 1)) = 0",
			[]vm.Value{plus, vm.NewRange(vm.Int(0), vm.Int(1), vm.Int(1))},
			vm.Int(0), // single element 0, no init — fn never called
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
		// (range 100) must stop at :stop without consuming the tail.
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
		rng := vm.NewRange(vm.Int(0), vm.Int(100), vm.Int(1))
		out, err := reduceFn.Invoke([]vm.Value{f, rng})
		if err != nil {
			t.Fatalf("reduce error: %v", err)
		}
		if out != vm.Keyword("stop") {
			t.Fatalf("got %v want :stop", out)
		}
	})
}
