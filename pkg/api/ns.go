package api

import (
	"reflect"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// NS is a handle for registering host values into one namespace. It is
// the namespace-targeted counterpart of LetGo.Def: same boxing, explicit
// destination. Obtain one from LetGo.NS.
type NS struct {
	ns *vm.Namespace
}

// NS resolves the named namespace, creating it if it does not exist yet,
// and returns a registration handle for it. Typical embedder use is
// grouping host functions into their own namespaces at startup:
//
//	patch := lg.NS("music.patch")
//	patch.Def("oscillator", oscillatorFn)
//	patch.DefShadowing("filter", filterFn) // collides with clojure.core/filter
func (l *LetGo) NS(name string) *NS {
	return &NS{ns: rt.NS(name)}
}

// Def boxes value exactly as LetGo.Def does and defines it in this
// namespace. Defining a name that is referred in from clojure.core works
// (the local var wins) but emits the Clojure-parity shadow warning; use
// DefShadowing when the collision is intentional.
func (n *NS) Def(name string, value any) error {
	val, err := vm.BoxValue(reflect.ValueOf(value))
	if err != nil {
		return err
	}
	n.ns.Def(name, val)
	return nil
}

// DefShadowing defines name in this namespace, first excluding it from
// the clojure.core auto-refer — the programmatic equivalent of
// (:refer-clojure :exclude [name]). The definition is silent and
// unqualified references in the namespace resolve to it instead of the
// core var.
//
// Boxing happens before any namespace mutation: if value is unboxable the
// namespace is left untouched, so a failed call is a true no-op rather than
// leaving name excluded-but-undefined (which would silently strip the
// unqualified clojure.core/<name> resolution).
func (n *NS) DefShadowing(name string, value any) error {
	val, err := vm.BoxValue(reflect.ValueOf(value))
	if err != nil {
		return err
	}
	n.ns.Exclude(name)
	n.ns.Def(name, val)
	return nil
}
