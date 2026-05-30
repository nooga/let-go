package vm

import (
	"strings"
	"testing"
)

func TestFormatErrorSurfacesGoPanicOrigin(t *testing.T) {
	fakeStack := `goroutine 1 [running]:
runtime/debug.Stack()
	/usr/lib/go/src/runtime/debug/stack.go:26 +0x5e
github.com/nooga/let-go/pkg/vm.recoverThrownPanic(0x1)
	/repo/pkg/vm/errors.go:181 +0x44
panic({0x1, 0x2})
	/usr/lib/go/src/runtime/panic.go:770 +0x132
github.com/nooga/let-go/pkg/rt.installTermNS.func7()
	/repo/pkg/rt/term_wasm.go:91 +0x88
github.com/nooga/let-go/pkg/vm.(*NativeFn).Invoke(...)
	/repo/pkg/vm/native_func.go:200 +0x40`
	gpe := &GoPanicError{value: "syscall/js: call of Value.Int on undefined", stack: []byte(fakeStack)}
	wrapped := &ExecutionError{message: "calling require", cause: gpe}
	out := FormatError(wrapped)
	if !strings.Contains(out, "go panic origin") {
		t.Fatalf("expected 'go panic origin' section, got:\n%s", out)
	}
	if !strings.Contains(out, "pkg/rt/term_wasm.go:91") {
		t.Fatalf("expected the crash site term_wasm.go:91, got:\n%s", out)
	}
	if strings.Contains(out, "recoverThrownPanic") {
		t.Fatalf("recover machinery should be filtered, got:\n%s", out)
	}
	t.Logf("\n%s", out)
}
