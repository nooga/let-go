package vm

import (
	"fmt"
	"math/big"
	"reflect"
)

type theRatioType struct{}

func (t *theRatioType) String() string  { return t.Name() }
func (t *theRatioType) Type() ValueType { return TypeType }
func (t *theRatioType) Unbox() any      { return reflect.TypeFor[*theRatioType]() }
func (t *theRatioType) Name() string    { return "let-go.lang.Ratio" }

func (t *theRatioType) Box(bare any) (Value, error) {
	switch v := bare.(type) {
	case *big.Rat:
		return NewRatio(v), nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var RatioType *theRatioType = &theRatioType{}

// Ratio wraps *math/big.Rat as a VM value.
type Ratio struct {
	val *big.Rat
}

func NewRatio(v *big.Rat) *Ratio {
	return &Ratio{val: v}
}

func NewRatioFromInts(num, den int64) *Ratio {
	return &Ratio{val: new(big.Rat).SetFrac64(num, den)}
}

// MustRatioFromString builds a Ratio from its "num/den" string form, panicking
// on an unparseable input. Intended for gogen-emitted constant literals: the
// source was already validated at read time (and reduced past integer form, so
// the value is genuinely a Ratio), making the panic unreachable in generated
// code. Gives a single-value expression form keyed off the value's own String().
func MustRatioFromString(s string) *Ratio {
	v, ok := new(big.Rat).SetString(s)
	if !ok {
		panic("MustRatioFromString: invalid Ratio literal: " + s)
	}
	return &Ratio{val: v}
}

func (r *Ratio) Val() *big.Rat { return r.val }

func (r *Ratio) Type() ValueType { return RatioType }
func (r *Ratio) Unbox() any      { return r.val }

func (r *Ratio) String() string {
	return r.val.RatString()
}

// Hash implements Hashable.
func (r *Ratio) Hash() uint32 {
	num := r.val.Num()
	den := r.val.Denom()
	h := hashUint64(uint64(num.Int64()))
	h = h*31 + hashUint64(uint64(den.Int64()))
	return h
}

// Equals compares two Ratios.
func (r *Ratio) Equals(other Value) bool {
	if o, ok := other.(*Ratio); ok {
		return r.val.Cmp(o.val) == 0
	}
	return false
}

// ToFloat64 converts to float64.
func (r *Ratio) ToFloat64() float64 {
	f, _ := r.val.Float64()
	return f
}

// MaybeSimplifyRatio returns an Int/BigInt if denominator is 1, otherwise Ratio.
func MaybeSimplifyRatio(r *big.Rat) Value {
	if r.IsInt() {
		num := r.Num()
		if num.IsInt64() {
			return MakeInt(int(num.Int64()))
		}
		return NewBigInt(new(big.Int).Set(num))
	}
	return &Ratio{val: r}
}

// IsRatio returns true if v is a *Ratio.
func IsRatio(v Value) bool {
	_, ok := v.(*Ratio)
	return ok
}

// ToRat converts a Value to *big.Rat if possible.
func ToRat(v Value) (*big.Rat, bool) {
	switch n := v.(type) {
	case *Ratio:
		return n.val, true
	case Int:
		return new(big.Rat).SetInt64(int64(n)), true
	case *BigInt:
		return new(big.Rat).SetInt(n.val), true
	}
	return nil, false
}

func (r *Ratio) GoString() string {
	return fmt.Sprintf("Ratio(%s)", r.val.RatString())
}
