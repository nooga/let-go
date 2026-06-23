package api_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestWithStorage(t *testing.T) {
	store := rt.NewMemoryStorage()
	lg, err := api.NewLetGo("withstorage-test", api.WithStorage(store))
	if err != nil {
		t.Fatal(err)
	}

	if v, err := lg.Run(`(storage/get "save:slot-1")`); err != nil || v != vm.NIL {
		t.Fatalf("missing storage/get = %#v (err %v), want nil", v, err)
	}
	if v, err := lg.Run(`(storage/set "save:slot-1" "{:turn 1}")`); err != nil || v != vm.NIL {
		t.Fatalf("storage/set = %#v (err %v), want nil", v, err)
	}
	if v, err := lg.Run(`(storage/get "save:slot-1")`); err != nil {
		t.Fatal(err)
	} else if s, ok := v.(vm.String); !ok || string(s) != "{:turn 1}" {
		t.Fatalf("storage/get = %#v, want save payload", v)
	}
	if v, err := lg.Run(`(storage/keys "save:")`); err != nil {
		t.Fatal(err)
	} else if got := stringSlice(t, v); len(got) != 1 || got[0] != "save:slot-1" {
		t.Fatalf("storage/keys = %#v, want [save:slot-1]", got)
	}
	if v, err := lg.Run(`(storage/remove "save:slot-1")`); err != nil || v != vm.NIL {
		t.Fatalf("storage/remove = %#v (err %v), want nil", v, err)
	}
	if v, err := lg.Run(`(storage/get "save:slot-1")`); err != nil || v != vm.NIL {
		t.Fatalf("storage/get after remove = %#v (err %v), want nil", v, err)
	}
}

func TestWithStorageIsolation(t *testing.T) {
	a := rt.NewMemoryStorage()
	b := rt.NewMemoryStorage()
	if err := a.Set("k", "a"); err != nil {
		t.Fatal(err)
	}
	if err := b.Set("k", "b"); err != nil {
		t.Fatal(err)
	}
	lgA, err := api.NewLetGo("withstorage-a", api.WithStorage(a))
	if err != nil {
		t.Fatal(err)
	}
	lgB, err := api.NewLetGo("withstorage-b", api.WithStorage(b))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		lg   *api.LetGo
		want string
	}{{lgA, "a"}, {lgB, "b"}} {
		v, err := tc.lg.Run(`(storage/get "k")`)
		if err != nil {
			t.Fatal(err)
		}
		if s, ok := v.(vm.String); !ok || string(s) != tc.want {
			t.Fatalf("storage/get = %#v, want %q", v, tc.want)
		}
	}
}

func stringSlice(t *testing.T, v vm.Value) []string {
	t.Helper()
	seqable, ok := v.(vm.Sequable)
	if !ok {
		t.Fatalf("value is not seqable: %T", v)
	}
	out := []string{}
	for s := seqable.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
		item, ok := s.First().(vm.String)
		if !ok {
			t.Fatalf("item is not string: %#v", s.First())
		}
		out = append(out, string(item))
	}
	return out
}
