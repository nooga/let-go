package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestKeyValue(t *testing.T) {
	if got := keyValue(""); got != vm.NIL {
		t.Errorf("keyValue(\"\") = %v, want nil", got)
	}
	for _, s := range []string{"l", "\x1b[A", "\x07"} {
		if got := keyValue(s); got != vm.String(s) {
			t.Errorf("keyValue(%q) = %v, want the string unchanged", s, got)
		}
	}

	m, ok := keyValue("\x1b[<20;12;7m").(*vm.PersistentMap)
	if !ok {
		t.Fatalf("keyValue(mouse report) is not a map")
	}
	want := map[string]vm.Value{
		"type": vm.Keyword("mouse"), "action": vm.Keyword("release"), "button": vm.Keyword("left"),
		"x": vm.MakeInt(12), "y": vm.MakeInt(7),
	}
	for k, v := range want {
		if got := m.ValueAt(vm.Keyword(k)); got != v {
			t.Errorf("%s = %v, want %v", k, got, v)
		}
	}
	mods, ok := m.ValueAt(vm.Keyword("mods")).(*vm.PersistentSet)
	if !ok || mods.Count() != vm.Int(2) || mods.Contains(vm.Keyword("shift")) != vm.TRUE || mods.Contains(vm.Keyword("ctrl")) != vm.TRUE {
		t.Errorf("mods = %v, want #{:shift :ctrl}", m.ValueAt(vm.Keyword("mods")))
	}
}
