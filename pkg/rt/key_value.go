package rt

import "github.com/nooga/let-go/pkg/vm"

// keyValue is what read-key returns for one token from the key source, shared
// by the native and wasm read-key so both return the same shape. Hybrid
// return: plain keys stay strings (no caller churn), an SGR mouse report
// becomes a tagged map, and "" is the end-of-input nil. Mouse reports only
// arrive once mouse reporting is enabled, so string-only callers are
// unaffected.
func keyValue(s string) vm.Value {
	if s == "" {
		return vm.NIL
	}
	if ev, ok := decodeSGRMouse(s); ok {
		return mouseEventMap(ev)
	}
	return vm.String(s)
}

// mouseEventMap renders a decoded SGR mouse report as the tagged map read-key
// hands to Clojure: {:type :mouse :action :press :button :left :x 10 :y 5
// :mods #{:shift}}. :mods is an empty set when no modifier is held.
func mouseEventMap(ev MouseEvent) vm.Value {
	var mods []vm.Value
	if ev.Shift {
		mods = append(mods, vm.Keyword("shift"))
	}
	if ev.Ctrl {
		mods = append(mods, vm.Keyword("ctrl"))
	}
	if ev.Meta {
		mods = append(mods, vm.Keyword("meta"))
	}
	m := vm.EmptyPersistentMap
	m = m.Assoc(vm.Keyword("type"), vm.Keyword("mouse")).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("action"), vm.Keyword(ev.Action)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("button"), vm.Keyword(ev.Button)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("x"), vm.MakeInt(ev.X)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("y"), vm.MakeInt(ev.Y)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("mods"), vm.NewPersistentSet(mods)).(*vm.PersistentMap)
	return m
}
