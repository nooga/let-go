package rt

import (
	"errors"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

type testNSLoader struct {
	calls int
	load  func(string, int) (*vm.Namespace, error)
}

func (l *testNSLoader) Load(name string) *vm.Namespace {
	ns, _ := l.LoadWithError(name)
	return ns
}

func (l *testNSLoader) LoadWithError(name string) (*vm.Namespace, error) {
	l.calls++
	return l.load(name, l.calls)
}

func TestLookupOrRegisterNS_RetriesAfterLoaderFailure(t *testing.T) {
	name := "rt.test.retry.loader.failure"
	origLoader := GetNSLoader()
	defer SetNSLoader(origLoader)

	delete(nsRegistry, name)
	delete(nsNeedsLoad, name)
	delete(nsLoadErrors, name)
	defer delete(nsRegistry, name)
	defer delete(nsNeedsLoad, name)
	defer delete(nsLoadErrors, name)

	ldr := &testNSLoader{
		load: func(_ string, call int) (*vm.Namespace, error) {
			if call < 2 {
				return nil, nil
			}
			return vm.NewNamespace(name), nil
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
	delete(nsLoadErrors, name)
	defer delete(nsRegistry, name)
	defer delete(nsNeedsLoad, name)
	defer delete(nsLoadErrors, name)

	ldr := &testNSLoader{
		load: func(_ string, _ int) (*vm.Namespace, error) { return nil, nil },
	}
	SetNSLoader(ldr)

	ns, err := RequireNS(name)
	if err == nil {
		t.Fatalf("expected RequireNS to fail, got ns=%v", ns)
	}
}

func TestRequireNSPreservesLoaderError(t *testing.T) {
	name := "rt.test.require.loader.original-error"
	boom := errors.New("original compile failure")
	origLoader := GetNSLoader()
	defer SetNSLoader(origLoader)

	delete(nsRegistry, name)
	delete(nsNeedsLoad, name)
	delete(nsLoadErrors, name)
	defer delete(nsRegistry, name)
	defer delete(nsNeedsLoad, name)
	defer delete(nsLoadErrors, name)

	SetNSLoader(&testNSLoader{
		load: func(_ string, _ int) (*vm.Namespace, error) { return nil, boom },
	})

	_, err := RequireNS(name)
	if !errors.Is(err, boom) {
		t.Fatalf("RequireNS error = %v, want original %v", err, boom)
	}
}
