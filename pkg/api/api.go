package api

import (
	"io"
	"reflect"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Option configures a LetGo runtime instance at construction. Options
// are applied in order; later options override earlier ones for the
// same configuration key.
type Option func(*config)

// config collects the resolved option values. Pure data, no behavior.
type config struct {
	stdout io.Writer
	stderr io.Writer
}

// WithStdout configures the runtime to route output written via *out*
// (i.e. (println ...), (print ...), (pr ...), (prn ...), and any user
// code that consults *out*) through w for this runtime instance.
//
// Implementation: each Run call pushes an IOHandle wrapping w as a
// dynamic binding on *out*, then pops it on return. Per-Run scope
// means two LetGo instances with different WithStdout values DO NOT
// interfere with each other's sequential output capture, and a
// LetGo constructed without WithStdout sees the unaltered os.Stdout
// default regardless of what other instances have done.
//
// Concurrency caveat: vm.Var's binding stack is process-global. If
// two goroutines call Run on different LetGo instances concurrently,
// their bindings push onto the same *out* stack — push/pop
// interleavings can scramble captures. For deterministic isolation,
// serialize Run calls or run different instances in different
// processes.
//
// Default: os.Stdout (no binding pushed; runtime root unchanged).
func WithStdout(w io.Writer) Option {
	return func(c *config) { c.stdout = w }
}

// WithStderr configures the runtime to route output written via *err*
// (HTTP server error logs, pod stderr, panic reports, etc.) through w
// for this runtime instance. Same per-Run binding semantics as
// WithStdout; same concurrency caveat.
//
// Default: os.Stderr.
func WithStderr(w io.Writer) Option {
	return func(c *config) { c.stderr = w }
}

// (Other options deliberately NOT exposed:
//
//   - WithStdin: stdin substitution is tied to the wake() / SAB protocol
//     deferred from nooga/let-go#174 and the readline-driven REPL path.
//     *in*'s root binding remains os.Stdin; embedders that need stdin
//     substitution today can rebind *in* manually before calling Run.
//
//   - WithEmit: a callback for (js/emit ...) would require either
//     a *emit-fn* var in the core ns that js/emit consults, or a
//     refactor of pkg/rt/js_wasm.go's _lgEmit lookup path. Neither is
//     done. Until the wiring is real, a public option that silently
//     does nothing would be misleading.)

type LetGo struct {
	cp     *vm.Consts
	c      *compiler.Context
	loader *resolver.NSResolver

	// stdoutHandle/stderrHandle are pre-constructed Boxed IOHandles, or
	// nil when no option was supplied. Run pushes them as dynamic
	// bindings on *out*/*err* and pops on return.
	stdoutHandle vm.Value
	stderrHandle vm.Value
}

// NewLetGo constructs a runtime. With no options, behavior is exactly
// as it was pre-option: *out* and *err* keep their global default root
// bindings (os.Stdout / os.Stderr). Each I/O option installed gets
// pushed as a dynamic binding around each Run call, so different
// LetGo instances don't share or overwrite each other's I/O routing.
func NewLetGo(ns string, opts ...Option) (*LetGo, error) {
	cfg := &config{}
	for _, o := range opts {
		o(cfg)
	}

	cp := vm.NewConsts()
	nso := rt.NS(ns)
	c := compiler.NewCompiler(cp, nso)
	loader := resolver.NewNSResolver(c, []string{"."})
	loader.DiscoverDepsEdn(".")
	ret := &LetGo{
		cp:     cp,
		c:      c,
		loader: loader,
	}
	rt.SetNSLoader(ret.loader)

	if cfg.stdout != nil {
		ret.stdoutHandle = vm.NewBoxed(rt.NewWriterHandle("api.WithStdout", cfg.stdout))
	}
	if cfg.stderr != nil {
		ret.stderrHandle = vm.NewBoxed(rt.NewWriterHandle("api.WithStderr", cfg.stderr))
	}

	return ret, nil
}

func (l *LetGo) SetLoadPath(path []string) {
	l.loader.SetPath(path)
}

func (l *LetGo) Def(name string, value any) error {
	val, err := vm.BoxValue(reflect.ValueOf(value))
	if err != nil {
		return err
	}
	l.c.CurrentNS().Def(name, val)

	return nil
}

// Run compiles and evaluates expr. If this LetGo was constructed with
// WithStdout / WithStderr, the corresponding IOHandle is pushed as a
// dynamic binding on *out* / *err* for the eval's scope and popped on
// return (including on error / panic).
//
// Two LetGo instances calling Run sequentially each get their own
// configured streams: A's Run pushes A's handle, evaluates, pops; then
// B's Run pushes B's handle. The root binding of *out*/*err* is
// never mutated.
func (l *LetGo) Run(expr string) (vm.Value, error) {
	if l.stdoutHandle != nil {
		if v := rt.LookupCoreVar("*out*"); v != nil {
			v.PushBinding(l.stdoutHandle)
			defer v.PopBinding()
		}
	}
	if l.stderrHandle != nil {
		if v := rt.LookupCoreVar("*err*"); v != nil {
			v.PushBinding(l.stderrHandle)
			defer v.PopBinding()
		}
	}

	c, err := l.c.Compile(expr)
	if err != nil {
		return vm.NIL, err
	}
	frame := vm.NewFrame(c, nil)
	result, err := frame.RunProtected()
	vm.ReleaseFrame(frame)
	return result, err
}
