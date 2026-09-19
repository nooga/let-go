package rt

import "github.com/nooga/let-go/pkg/vm"

// WithFile binds core/*file* to path for the duration of fn. Used by every
// file-loading entry point (CLI file runner, require/load, test harness).
func WithFile(path string, fn func() error) error {
	v := NS(NameCoreNS).LookupLocal(vm.Symbol("*file*"))
	if v == nil {
		return fn()
	}
	v.PushBinding(vm.String(path))
	defer v.PopBinding()
	return fn()
}
