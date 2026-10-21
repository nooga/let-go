package rt

import (
	"fmt"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// prThroughToString renders vs through print-method into a single string,
// space-separated, at the given readability.
func prThroughToString(vs []vm.Value, readably bool) (vm.Value, error) {
	sb := &strings.Builder{}
	w := newLGWriter(sb, nil)
	w.readably = readably
	boxed := vm.NewBoxed(w)
	for i := range vs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if err := dispatchPrintMethod(vs[i], boxed, w, readably); err != nil {
			return vm.NIL, err
		}
	}
	return vm.String(sb.String()), nil
}

// renderLeaf renders a scalar/leaf value. readably=true → reader form
// (strings quoted, chars as \x); false → human form (raw string/char).
func renderLeaf(v vm.Value, readably bool) string {
	if !readably {
		switch v.Type() {
		case vm.StringType:
			return string(v.(vm.String))
		case vm.CharType:
			return string(rune(v.(vm.Char)))
		}
	} else if sb, ok := v.(*hostStringBuilder); ok {
		// pr shows a builder's contents as a string literal.
		return vm.String(sb.b.String()).String()
	}
	return v.String()
}

// writerFromValue unwraps a boxed *LGWriter passed as a print-method writer.
func writerFromValue(v vm.Value) (*LGWriter, error) {
	if b, ok := v.(*vm.Boxed); ok {
		if w, ok := b.Unbox().(*LGWriter); ok {
			return w, nil
		}
	}
	return nil, fmt.Errorf("print-method writer is not an LGWriter")
}

var printMethodVar *vm.Var

// dispatchPrintMethod writes obj to w by routing through the print-method
// multimethod. Falls back to native String() if print-method is unbound
// (bootstrap, before core.lg defines it).
func dispatchPrintMethod(obj vm.Value, boxedWriter vm.Value, w *LGWriter, readably bool) error {
	if printMethodVar == nil {
		printMethodVar = CoreNS.LookupLocal(vm.Symbol("print-method"))
	}
	if printMethodVar == nil || printMethodVar.Deref() == vm.NIL {
		_, err := w.WriteString(renderLeaf(obj, readably))
		return err
	}
	fn, ok := vm.AsFn(printMethodVar.Deref())
	if !ok {
		_, err := w.WriteString(renderLeaf(obj, readably))
		return err
	}
	_, err := fn.Invoke([]vm.Value{obj, boxedWriter})
	return err
}

// printMethodDefault is the :default method of print-method: [obj writer]
// writes obj's leaf form to the boxed *LGWriter using readability from the writer.
var printMethodDefault, _ = vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
	if len(vs) != 2 {
		return vm.NIL, fmt.Errorf("print-method default expects [obj writer]")
	}
	w, err := writerFromValue(vs[1])
	if err != nil {
		return vm.NIL, err
	}
	_, err = w.WriteString(renderLeaf(vs[0], w.readably))
	return vm.NIL, err
})
