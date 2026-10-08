package rt

import "github.com/nooga/let-go/pkg/vm"

// loadScopedVars are the compile-time flags Clojure's load binds per file, so
// a file's top-level set! of one of them ends with that file's load instead of
// reaching every file loaded after it.
var loadScopedVars = []vm.Symbol{"*warn-on-reflection*", "*unchecked-math*"}

// WithFile establishes the per-load bindings Clojure's load establishes for the
// duration of fn: core/*file* bound to path, and each of loadScopedVars bound
// to its current value. Used by every file-loading entry point (CLI file
// runner, require/load, test harness).
func WithFile(path string, fn func() error) error {
	core := NS(NameCoreNS)
	if v := core.LookupLocal(vm.Symbol("*file*")); v != nil {
		v.PushBinding(vm.String(path))
		defer v.PopBinding()
	}
	for _, name := range loadScopedVars {
		if v := core.LookupLocal(name); v != nil {
			v.PushBinding(v.Deref())
			defer v.PopBinding()
		}
	}
	return fn()
}
