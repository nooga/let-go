package api_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/vm"
)

// scriptedKeys is a deterministic rt.KeySource for tests: it hands out a
// fixed sequence of keys, then reports end-of-input.
type scriptedKeys struct {
	keys []string
	i    int
}

func (s *scriptedKeys) ReadKey() (string, error) {
	if s.i >= len(s.keys) {
		return "", nil // EOF → read-key nil
	}
	k := s.keys[s.i]
	s.i++
	return k, nil
}

func (s *scriptedKeys) KeyPending() bool { return s.i < len(s.keys) }

// TestWithKeySource proves the input dual of TestWithStdout: keys supplied
// via api.WithKeySource flow out of (term/read-key) in order, and exhaustion
// surfaces as the nil contract — no real terminal involved.
func TestWithKeySource(t *testing.T) {
	ks := &scriptedKeys{keys: []string{"a", "\x1b[A", "b"}}
	lg, err := api.NewLetGo("withkeys-test", api.WithKeySource(ks))
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"a", "\x1b[A", "b"} {
		v, err := lg.Run(`(term/read-key)`)
		if err != nil {
			t.Fatal(err)
		}
		s, ok := v.(vm.String)
		if !ok || string(s) != want {
			t.Fatalf("read-key = %#v, want %q", v, want)
		}
	}

	// Exhausted → nil.
	v, err := lg.Run(`(term/read-key)`)
	if err != nil {
		t.Fatal(err)
	}
	if v != vm.NIL {
		t.Errorf("read-key after exhaustion = %#v, want nil", v)
	}
}

// TestWithKeySourcePending proves key-pending? routes through the same bound
// source: true while keys remain, false once drained.
func TestWithKeySourcePending(t *testing.T) {
	ks := &scriptedKeys{keys: []string{"x"}}
	lg, err := api.NewLetGo("withkeys-pending-test", api.WithKeySource(ks))
	if err != nil {
		t.Fatal(err)
	}

	if v, err := lg.Run(`(term/key-pending?)`); err != nil || v != vm.TRUE {
		t.Fatalf("key-pending? with a queued key = %#v (err %v), want true", v, err)
	}
	if _, err := lg.Run(`(term/read-key)`); err != nil {
		t.Fatal(err)
	}
	if v, err := lg.Run(`(term/key-pending?)`); err != nil || v != vm.FALSE {
		t.Fatalf("key-pending? after draining = %#v (err %v), want false", v, err)
	}
}
