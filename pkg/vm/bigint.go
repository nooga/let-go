package vm

import (
	"fmt"
	"math/big"
	"reflect"
)

type theBigIntType struct{}

func (t *theBigIntType) String() string     { return t.Name() }
func (t *theBigIntType) Type() ValueType    { return TypeType }
func (t *theBigIntType) Unbox() interface{} { return reflect.TypeOf(t) }
func (t *theBigIntType) Name() string       { return "let-go.lang.BigInt" }

func (t *theBigIntType) Box(bare interface{}) (Value, error) {
	switch v := bare.(type) {
	case *big.Int:
		return &BigInt{val: v}, nil
	case int64:
		return &BigInt{val: big.NewInt(v)}, nil
	case int:
		return &BigInt{val: big.NewInt(int64(v))}, nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var BigIntType *theBigIntType = &theBigIntType{}

// BigInt wraps *math/big.Int as a VM value.
type BigInt struct {
	val *big.Int
}

func NewBigInt(v *big.Int) *BigInt {
	return &BigInt{val: v}
}

func NewBigIntFromString(s string) (*BigInt, bool) {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, false
	}
	return &BigInt{val: v}, true
}

func NewBigIntFromInt64(n int64) *BigInt {
	return &BigInt{val: big.NewInt(n)}
}

func (b *BigInt) Val() *big.Int { return b.val }

func (b *BigInt) Type() ValueType    { return BigIntType }
func (b *BigInt) Unbox() interface{} { return b.val }

func (b *BigInt) String() string {
	return b.val.String() + "N"
}

// Hash implements Hashable — consistent with Int.Hash() for values that fit in int64.
func (b *BigInt) Hash() uint32 {
	if b.val.IsInt64() {
		return hashUint64(uint64(b.val.Int64()))
	}
	// For values that don't fit in int64, hash the bytes
	bytes := b.val.Bytes()
	h := uint32(0)
	for _, by := range bytes {
		h = h*31 + uint32(by)
	}
	if b.val.Sign() < 0 {
		h = ^h
	}
	return h
}

// ToInt64 returns the int64 value if it fits, or false.
func (b *BigInt) ToInt64() (int64, bool) {
	if b.val.IsInt64() {
		return b.val.Int64(), true
	}
	return 0, false
}

// Equals compares two BigInts.
func (b *BigInt) Equals(other *BigInt) bool {
	return b.val.Cmp(other.val) == 0
}

// --- Helper: convert Value to *big.Int for arithmetic ---

func ToBigInt(v Value) (*big.Int, bool) {
	switch n := v.(type) {
	case *BigInt:
		return n.val, true
	case Int:
		return big.NewInt(int64(n)), true
	}
	return nil, false
}

// MaybeDowngrade returns an Int if the BigInt fits in int64, otherwise BigInt.
func MaybeDowngrade(b *big.Int) Value {
	if b.IsInt64() {
		return MakeInt(int(b.Int64()))
	}
	return &BigInt{val: b}
}

// IsBigInt returns true if v is a *BigInt.
func IsBigInt(v Value) bool {
	_, ok := v.(*BigInt)
	return ok
}

// BigIntFromValue returns a string representation suitable for printing.
func (b *BigInt) GoString() string {
	return fmt.Sprintf("BigInt(%s)", b.val.String())
}
