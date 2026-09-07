package vm_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// A conversion rejected by unboxMapInto used to be discarded by
// boxArgForReflect, which fell through to the Unbox fallback and handed
// reflect.Call a *vm.PersistentMap. The rejection still happened, but the
// caller saw
//
//	reflect: Call using *vm.PersistentMap as type map[string]interface {}
//
// rather than the reason. The careful diagnostics unboxMapInto produces were
// invisible at the one call site a wrapper author actually hits.
func TestBoxArgSurfacesMapConversionError(t *testing.T) {
	fn, err := vm.NativeFnType.Box(func(m map[string]any) vm.Value {
		return vm.Int(int64(len(m)))
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	// A numeric key into a string key target: rejected so it is not silently
	// rune-converted into "A".
	in := vm.NewPersistentMap([]vm.Value{vm.Int(65), vm.String("Ada")})

	var invokeErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panicked instead of returning an error: %v", r)
			}
		}()
		_, invokeErr = fn.(*vm.NativeFn).Invoke([]vm.Value{in})
	}()

	if invokeErr == nil {
		t.Fatal("want an error, got nil")
	}
	msg := invokeErr.Error()
	if strings.Contains(msg, "reflect: Call using") {
		t.Fatalf("still surfacing the reflect panic: %s", msg)
	}
	for _, want := range []string{"arg 0", "map key", "65"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error %q does not mention %q", msg, want)
		}
	}
}

// An unhashable key is the other error unboxMapInto produces; it must reach
// the caller the same way rather than as a recovered panic.
func TestBoxArgSurfacesUnhashableKeyError(t *testing.T) {
	fn, err := vm.NativeFnType.Box(func(m map[any]any) vm.Value {
		return vm.Int(int64(len(m)))
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	key := vm.ArrayVector{vm.Int(1), vm.Int(2)}
	in := vm.NewPersistentMap([]vm.Value{key, vm.String("v")})

	_, invokeErr := fn.(*vm.NativeFn).Invoke([]vm.Value{in})
	if invokeErr == nil {
		t.Fatal("want an error, got nil")
	}
	if !strings.Contains(invokeErr.Error(), "hashable") {
		t.Errorf("error %q does not explain the key is unhashable", invokeErr.Error())
	}
}

// The Unbox fallback is legitimate: a *Boxed already holding a Go map has no
// let-go map to convert, so unboxMapInto fails and the fallback serves it.
// Reporting the conversion error eagerly would break calls that work today.
func TestBoxArgFallbackStillServesABoxedGoMap(t *testing.T) {
	fn, err := vm.NativeFnType.Box(func(m map[string]string) vm.Value {
		return vm.String(m["k"])
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	boxed := vm.NewBoxed(map[string]string{"k": "v"})
	got, invokeErr := fn.(*vm.NativeFn).Invoke([]vm.Value{boxed})
	if invokeErr != nil {
		t.Fatalf("fallback broke: %v", invokeErr)
	}
	if got != vm.String("v") {
		t.Fatalf("got %v, want \"v\"", got)
	}
}
