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
	emit   func(name, dataJSON string)
	keys   rt.KeySource
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

// WithEmit configures the runtime to route (js/emit event-name data)
// through fn for this instance, the emit dual of WithStdout. fn receives
// the event name and the JSON-marshaled data — the same (name, dataJSON)
// pair the WASM bundle hands to LetGoHost.onEmit — so a Go embedder gets
// the events its guest dispatches without depending on the browser host.
//
// Implementation: each Run pushes a FuncEmitter wrapping fn as a dynamic
// binding on *emit*, popped on return. Same per-Run isolation and
// process-global-binding-stack concurrency caveat as WithStdout.
//
// Default: no-op (the *emit* root nopEmitter; (js/emit ...) is dropped).
func WithEmit(fn func(name, dataJSON string)) Option {
	return func(c *config) { c.emit = fn }
}

// WithKeySource routes term/read-key and key-pending? through ks for this
// instance, the input dual of WithStdout. An embedder or test supplies the
// keystrokes the guest reads — driving a TUI without a real terminal, or
// feeding a deterministic script. ks is an rt.KeySource (ReadKey +
// KeyPending); see pkg/rt/keysource.go.
//
// Implementation: each Run pushes ks as a dynamic binding on *keys*, popped
// on return. Same per-Run isolation and process-global-binding-stack
// concurrency caveat as WithStdout.
//
// Default: the platform key source (stdin+SIGWINCH on native, the SAB in
// WASM). Wake — unblocking a parked read-key without a real key — is not part
// of this seam yet; see the note in pkg/rt/wasm/lg-host.js.
func WithKeySource(ks rt.KeySource) Option {
	return func(c *config) { c.keys = ks }
}

// (Other options deliberately NOT exposed:
//
//   - WithStdin: line-stream stdin substitution is the *in* io.Reader dual
//     of WithStdout, separate from WithKeySource's interactive key input.
//     *in*'s root binding remains os.Stdin; embedders that need stdin
//     substitution today can rebind *in* manually before calling Run.)

type LetGo struct {
	cp     *vm.Consts
	c      *compiler.Context
	loader *resolver.NSResolver

	// stdoutHandle/stderrHandle are pre-constructed Boxed IOHandles, or
	// nil when no option was supplied. Run pushes them as dynamic
	// bindings on *out*/*err* and pops on return.
	stdoutHandle vm.Value
	stderrHandle vm.Value
	emitHandle   vm.Value
	keysHandle   vm.Value
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
	if cfg.emit != nil {
		ret.emitHandle = vm.NewBoxed(rt.FuncEmitter(cfg.emit))
	}
	if cfg.keys != nil {
		ret.keysHandle = vm.NewBoxed(cfg.keys)
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
	if l.emitHandle != nil {
		if v := rt.LookupCoreVar("*emit*"); v != nil {
			v.PushBinding(l.emitHandle)
			defer v.PopBinding()
		}
	}
	if l.keysHandle != nil {
		if v := rt.LookupCoreVar("*keys*"); v != nil {
			v.PushBinding(l.keysHandle)
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
