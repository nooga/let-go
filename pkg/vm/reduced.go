package vm

import "fmt"

// Reduced wraps a value to signal early termination in reduce/transduce.
type Reduced struct {
	value Value
}

func NewReduced(v Value) *Reduced {
	return &Reduced{value: v}
}

func (r *Reduced) Deref() Value        { return r.value }
func (r *Reduced) Type() ValueType     { return ReducedType }
func (r *Reduced) Unbox() interface{}  { return r.value }
func (r *Reduced) String() string {
	return fmt.Sprintf("#reduced<%s>", r.value.String())
}

func IsReduced(v Value) bool {
	_, ok := v.(*Reduced)
	return ok
}

// ReducedType
type theReducedType struct{}

func (t *theReducedType) String() string     { return t.Name() }
func (t *theReducedType) Type() ValueType    { return TypeType }
func (t *theReducedType) Unbox() interface{} { return nil }
func (t *theReducedType) Name() string       { return "let-go.lang.Reduced" }
func (t *theReducedType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var ReducedType *theReducedType = &theReducedType{}
