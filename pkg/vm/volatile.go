package vm

import "fmt"

// Volatile is a mutable container with no synchronization.
// Useful for transducers and other single-threaded mutation contexts.
type Volatile struct {
	value Value
}

func NewVolatile(v Value) *Volatile {
	return &Volatile{value: v}
}

func (v *Volatile) Deref() Value {
	return v.value
}

func (v *Volatile) Reset(val Value) Value {
	old := v.value
	v.value = val
	_ = old
	return val
}

func (v *Volatile) Type() ValueType    { return VolatileType }
func (v *Volatile) Unbox() interface{} { return v }
func (v *Volatile) String() string {
	return fmt.Sprintf("#<Volatile@%p: %s>", v, v.value.String())
}

// VolatileType
type theVolatileType struct{}

func (t *theVolatileType) String() string     { return t.Name() }
func (t *theVolatileType) Type() ValueType    { return TypeType }
func (t *theVolatileType) Unbox() interface{} { return nil }
func (t *theVolatileType) Name() string       { return "let-go.lang.Volatile" }
func (t *theVolatileType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var VolatileType *theVolatileType = &theVolatileType{}
