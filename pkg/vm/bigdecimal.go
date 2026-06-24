package vm

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"
	"unsafe"
)

type theBigDecimalType struct{}

func (t *theBigDecimalType) String() string  { return t.Name() }
func (t *theBigDecimalType) Type() ValueType { return TypeType }
func (t *theBigDecimalType) Unbox() any      { return reflect.TypeFor[*theBigDecimalType]() }
func (t *theBigDecimalType) Name() string    { return "let-go.lang.BigDecimal" }

func (t *theBigDecimalType) Box(bare any) (Value, error) {
	switch v := bare.(type) {
	case *big.Float:
		return &BigDecimal{val: v}, nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var BigDecimalType *theBigDecimalType = &theBigDecimalType{}

// BigDecimal wraps *math/big.Float as a VM value with arbitrary precision.
type BigDecimal struct {
	val *big.Float
}

const bigDecimalPrec = 128 // bits of precision
const BigDecimalPrecConst = bigDecimalPrec

func NewBigDecimal(v *big.Float) *BigDecimal {
	return &BigDecimal{val: v}
}

func NewBigDecimalFromString(s string) (*BigDecimal, bool) {
	v, _, err := new(big.Float).SetPrec(bigDecimalPrec).Parse(s, 10)
	if err != nil {
		return nil, false
	}
	return &BigDecimal{val: v}, true
}

// MustBigDecimalFromString builds a BigDecimal from a string, panicking on an
// unparseable input. Intended for gogen-emitted constant literals, whose source
// strings were already validated at read time, so the panic is unreachable in
// generated code — it exists only to give a single-value expression form (the
// (value, ok) signature of NewBigDecimalFromString can't be used inline).
func MustBigDecimalFromString(s string) *BigDecimal {
	v, ok := NewBigDecimalFromString(s)
	if !ok {
		panic("MustBigDecimalFromString: invalid BigDecimal literal: " + s)
	}
	return v
}

func NewBigDecimalFromFloat64(f float64) *BigDecimal {
	return &BigDecimal{val: new(big.Float).SetPrec(bigDecimalPrec).SetFloat64(f)}
}

func NewBigDecimalFromInt64(n int64) *BigDecimal {
	return &BigDecimal{val: new(big.Float).SetPrec(bigDecimalPrec).SetInt64(n)}
}

func (b *BigDecimal) Val() *big.Float { return b.val }

func (b *BigDecimal) Type() ValueType { return BigDecimalType }
func (b *BigDecimal) Unbox() any      { return b.val }

func (b *BigDecimal) String() string {
	s := b.val.Text('f', -1)
	// Ensure there's a decimal point
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	// Remove trailing zeros after decimal, but keep at least one
	if idx := strings.IndexByte(s, '.'); idx >= 0 {
		end := len(s)
		for end > idx+2 && s[end-1] == '0' {
			end--
		}
		s = s[:end]
	}
	return s + "M"
}

// Hash implements Hashable.
func (b *BigDecimal) Hash() uint32 {
	f, _ := b.val.Float64()
	if math.IsNaN(f) {
		return 0
	}
	return hashUint64(*(*uint64)(unsafe.Pointer(&f)))
}

// Equals compares BigDecimals.
func (b *BigDecimal) Equals(other Value) bool {
	if o, ok := other.(*BigDecimal); ok {
		return b.val.Cmp(o.val) == 0
	}
	return false
}

// ToFloat64 converts to float64.
func (b *BigDecimal) ToFloat64() float64 {
	f, _ := b.val.Float64()
	return f
}

// IsBigDecimal returns true if v is a *BigDecimal.
func IsBigDecimal(v Value) bool {
	_, ok := v.(*BigDecimal)
	return ok
}

func (b *BigDecimal) GoString() string {
	return fmt.Sprintf("BigDecimal(%s)", b.val.Text('f', -1))
}
