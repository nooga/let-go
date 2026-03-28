package vm

import (
	"fmt"
	"sync"
)

// Delay wraps a thunk (zero-arg function) that is evaluated at most once.
// The result is cached after the first call to Force().
type Delay struct {
	fn       Fn
	value    Value
	realized bool
	err      error
	mu       sync.Mutex
}

func NewDelay(fn Fn) *Delay {
	return &Delay{fn: fn}
}

func (d *Delay) Force() (Value, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.realized {
		if d.err != nil {
			return NIL, d.err
		}
		return d.value, nil
	}
	d.realized = true
	v, err := d.fn.Invoke(nil)
	if err != nil {
		d.err = err
		return NIL, err
	}
	d.value = v
	d.fn = nil // release the thunk for GC
	return v, nil
}

func (d *Delay) IsRealized() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.realized
}

func (d *Delay) Deref() Value {
	v, err := d.Force()
	if err != nil {
		panic(&thrownPanic{err: err})
	}
	return v
}

func (d *Delay) Type() ValueType    { return DelayType }
func (d *Delay) Unbox() interface{} { return d }
func (d *Delay) String() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.realized {
		return fmt.Sprintf("#<Delay@%p: %s>", d, d.value.String())
	}
	return fmt.Sprintf("#<Delay@%p: :pending>", d)
}

// DelayType
type theDelayType struct{}

func (t *theDelayType) String() string     { return t.Name() }
func (t *theDelayType) Type() ValueType    { return TypeType }
func (t *theDelayType) Unbox() interface{} { return nil }
func (t *theDelayType) Name() string       { return "let-go.lang.Delay" }
func (t *theDelayType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var DelayType *theDelayType = &theDelayType{}
