package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

type testNSLoader struct {
	calls int
	load  func(string, int) *vm.Namespace
}

func (l *testNSLoader) Load(name string) *vm.Namespace {
	l.calls++
	return l.load(name, l.calls)
}

func TestLookupOrRegisterNS_RetriesAfterLoaderFailure(t *testing.T) {
	name := "rt.test.retry.loader.failure"
	origLoader := GetNSLoader()
	defer SetNSLoader(origLoader)

	delete(nsRegistry, name)
	delete(nsNeedsLoad, name)
	defer delete(nsRegistry, name)
	defer delete(nsNeedsLoad, name)

	ldr := &testNSLoader{
		load: func(_ string, call int) *vm.Namespace {
			if call < 2 {
				return nil
			}
			return vm.NewNamespace(name)
		},
	}
	SetNSLoader(ldr)

	first := LookupOrRegisterNS(name)
	if first == nil {
		t.Fatalf("first lookup returned nil namespace")
	}
	if !nsNeedsLoad[name] {
		t.Fatalf("expected nsNeedsLoad[%q] to remain true after load failure", name)
	}

	second := LookupOrRegisterNS(name)
	if second == nil {
		t.Fatalf("second lookup returned nil namespace")
	}
	if ldr.calls < 2 {
		t.Fatalf("expected loader retry on second lookup, got calls=%d", ldr.calls)
	}
	if nsNeedsLoad[name] {
		t.Fatalf("expected nsNeedsLoad[%q] to clear after successful load", name)
	}
}

func TestRequireNS_ReturnsErrorWhenLoaderCannotLoad(t *testing.T) {
	name := "rt.test.require.loader.failure"
	origLoader := GetNSLoader()
	defer SetNSLoader(origLoader)

	delete(nsRegistry, name)
	delete(nsNeedsLoad, name)
	defer delete(nsRegistry, name)
	defer delete(nsNeedsLoad, name)

	ldr := &testNSLoader{
		load: func(_ string, _ int) *vm.Namespace { return nil },
	}
	SetNSLoader(ldr)

	ns, err := RequireNS(name)
	if err == nil {
		t.Fatalf("expected RequireNS to fail, got ns=%v", ns)
	}
}
