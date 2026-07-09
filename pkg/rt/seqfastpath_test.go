package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Behavior tests for ArrayVector direct-slice fast paths (first, second, next, rest, concat).
func TestSeqFastPathFirst(t *testing.T) {
	firstVar := LookupCoreVar("first")
	if firstVar == nil {
		t.Fatal("core/first not found")
	}
	firstFn, ok := firstVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("first is not an Fn")
	}

	cases := []struct {
		name string
		arg  vm.Value
		want vm.Value
	}{
		{
			"vector of 3",
			vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)}),
			vm.Int(1),
		},
		{
			"empty vector",
			vm.NewArrayVector([]vm.Value{}),
			vm.NIL,
		},
		{
			"nil",
			vm.NIL,
			vm.NIL,
		},
		{
			"single element",
			vm.NewArrayVector([]vm.Value{vm.Int(42)}),
			vm.Int(42),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := firstFn.Invoke([]vm.Value{c.arg})
			if err != nil {
				t.Fatalf("first error: %v", err)
			}
			if out != c.want {
				t.Fatalf("got %v want %v", out, c.want)
			}
		})
	}
}

func TestSeqFastPathSecond(t *testing.T) {
	secondVar := LookupCoreVar("second")
	if secondVar == nil {
		t.Fatal("core/second not found")
	}
	secondFn, ok := secondVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("second is not an Fn")
	}

	cases := []struct {
		name string
		arg  vm.Value
		want vm.Value
	}{
		{
			"vector of 3",
			vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)}),
			vm.Int(2),
		},
		{
			"single element",
			vm.NewArrayVector([]vm.Value{vm.Int(1)}),
			vm.NIL,
		},
		{
			"empty vector",
			vm.NewArrayVector([]vm.Value{}),
			vm.NIL,
		},
		{
			"nil",
			vm.NIL,
			vm.NIL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := secondFn.Invoke([]vm.Value{c.arg})
			if err != nil {
				t.Fatalf("second error: %v", err)
			}
			if out != c.want {
				t.Fatalf("got %v want %v", out, c.want)
			}
		})
	}
}

func TestSeqFastPathNext(t *testing.T) {
	nextVar := LookupCoreVar("next")
	if nextVar == nil {
		t.Fatal("core/next not found")
	}
	nextFn, ok := nextVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("next is not an Fn")
	}

	cases := []struct {
		name    string
		arg     vm.Value
		wantNil bool
		wantStr string
	}{
		{
			"vector of 3",
			vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)}),
			false,
			"(2 3)",
		},
		{
			"single element",
			vm.NewArrayVector([]vm.Value{vm.Int(1)}),
			true,
			"",
		},
		{
			"empty vector",
			vm.NewArrayVector([]vm.Value{}),
			true,
			"",
		},
		{
			"nil",
			vm.NIL,
			true,
			"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := nextFn.Invoke([]vm.Value{c.arg})
			if err != nil {
				t.Fatalf("next error: %v", err)
			}
			if c.wantNil {
				if out != vm.NIL && out != nil {
					t.Fatalf("got non-nil %v (type %T) want nil or vm.NIL", out, out)
				}
			} else {
				if out == nil || out == vm.NIL {
					t.Fatalf("got nil want %s", c.wantStr)
				}
				if out.String() != c.wantStr {
					t.Fatalf("got %s want %s", out.String(), c.wantStr)
				}
			}
		})
	}
}

func TestSeqFastPathRest(t *testing.T) {
	restVar := LookupCoreVar("rest")
	if restVar == nil {
		t.Fatal("core/rest not found")
	}
	restFn, ok := restVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("rest is not an Fn")
	}

	cases := []struct {
		name    string
		arg     vm.Value
		wantStr string
	}{
		{
			"vector of 3",
			vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)}),
			"(2 3)",
		},
		{
			"single element",
			vm.NewArrayVector([]vm.Value{vm.Int(1)}),
			"()",
		},
		{
			"empty vector",
			vm.NewArrayVector([]vm.Value{}),
			"()",
		},
		{
			"nil",
			vm.NIL,
			"()",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := restFn.Invoke([]vm.Value{c.arg})
			if err != nil {
				t.Fatalf("rest error: %v", err)
			}
			if out == nil {
				t.Fatalf("got nil want %s", c.wantStr)
			}
			if out.String() != c.wantStr {
				t.Fatalf("got %s want %s", out.String(), c.wantStr)
			}
		})
	}
}

func TestSeqFastPathConcat(t *testing.T) {
	concatVar := LookupCoreVar("concat*")
	if concatVar == nil {
		t.Fatal("core/concat* not found")
	}
	concatFn, ok := concatVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("concat* is not an Fn")
	}

	cases := []struct {
		name    string
		args    []vm.Value
		wantStr string
	}{
		{
			"two vectors",
			[]vm.Value{
				vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2)}),
				vm.NewArrayVector([]vm.Value{vm.Int(3)}),
			},
			"(1 2 3)",
		},
		{
			"no args",
			[]vm.Value{},
			"()",
		},
		{
			"empty vectors and nil",
			[]vm.Value{
				vm.NewArrayVector([]vm.Value{}),
				vm.NIL,
				vm.NewArrayVector([]vm.Value{}),
			},
			"()",
		},
		{
			"mixed vectors and nil",
			[]vm.Value{
				vm.NewArrayVector([]vm.Value{vm.Int(1)}),
				vm.NIL,
				vm.NewArrayVector([]vm.Value{vm.Int(2)}),
			},
			"(1 2)",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := concatFn.Invoke(c.args)
			if err != nil {
				t.Fatalf("concat error: %v", err)
			}
			if out == nil {
				t.Fatalf("got nil want %s", c.wantStr)
			}
			if out.String() != c.wantStr {
				t.Fatalf("got %s want %s", out.String(), c.wantStr)
			}
		})
	}
}
