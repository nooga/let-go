package vm

import (
	"fmt"
	"math"
	"math/big"
)

const maxIntValue = Int(int(^uint(0) >> 1))
const minIntValue = -maxIntValue - 1

func checkedAddInt(a, b Int) (Int, bool) {
	if (b > 0 && a > maxIntValue-b) || (b < 0 && a < minIntValue-b) {
		return 0, false
	}
	return a + b, true
}

func checkedSubInt(a, b Int) (Int, bool) {
	if (b < 0 && a > maxIntValue+b) || (b > 0 && a < minIntValue+b) {
		return 0, false
	}
	return a - b, true
}

func checkedMulInt(a, b Int) (Int, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	if (a == minIntValue && b == -1) || (b == minIntValue && a == -1) {
		return 0, false
	}
	if a > 0 {
		if b > 0 && a > maxIntValue/b {
			return 0, false
		}
		if b < 0 && b < minIntValue/a {
			return 0, false
		}
	} else {
		if b > 0 && a < minIntValue/b {
			return 0, false
		}
		if b < 0 && a < maxIntValue/b {
			return 0, false
		}
	}
	return a * b, true
}

// Numeric type promotion and arithmetic dispatch.
// All functions use direct type assertions (no Unbox) to avoid allocation.
// Promotion: Int op BigInt → BigInt, BigInt op Float → Float, Int op Float → Float.

func normalizeFloat32(v Value) Value {
	if f, ok := v.(Float32); ok {
		return Float(f)
	}
	return v
}

// NumAdd adds two numeric Values.
func NumAdd(a, b Value) (Value, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			r, ok := checkedAddInt(av, bv)
			if !ok {
				return NIL, fmt.Errorf("integer overflow")
			}
			return MakeInt(int(r)), nil
		case Float:
			return Float(float64(av) + float64(bv)), nil
		case *BigInt:
			r := new(big.Int).Add(big.NewInt(int64(av)), bv.val)
			return NewBigInt(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return Float(float64(av) + float64(bv)), nil
		case Float:
			return Float(float64(av) + float64(bv)), nil
		case *BigInt:
			bf, _ := bv.val.Float64()
			return Float(float64(av) + float64(bf)), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			r := new(big.Int).Add(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case Float:
			af, _ := av.val.Float64()
			return Float(float64(af) + float64(bv)), nil
		case *BigInt:
			r := new(big.Int).Add(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	if r, err := numBinOpFallback(a, b, ratAdd, bdAdd, floatAdd); err == nil {
		return r, nil
	}
	return NIL, fmt.Errorf("cannot add %s and %s", a.Type().Name(), b.Type().Name())
}

// NumSub subtracts b from a.
func NumSub(a, b Value) (Value, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			r, ok := checkedSubInt(av, bv)
			if !ok {
				return NIL, fmt.Errorf("integer overflow")
			}
			return MakeInt(int(r)), nil
		case Float:
			return Float(float64(av) - float64(bv)), nil
		case *BigInt:
			r := new(big.Int).Sub(big.NewInt(int64(av)), bv.val)
			return NewBigInt(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return Float(float64(av) - float64(bv)), nil
		case Float:
			return Float(float64(av) - float64(bv)), nil
		case *BigInt:
			bf, _ := bv.val.Float64()
			return Float(float64(av) - bf), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			r := new(big.Int).Sub(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case Float:
			af, _ := av.val.Float64()
			return Float(af - float64(bv)), nil
		case *BigInt:
			r := new(big.Int).Sub(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	if r, err := numBinOpFallback(a, b, ratSub, bdSub, floatSub); err == nil {
		return r, nil
	}
	return NIL, fmt.Errorf("cannot subtract %s and %s", a.Type().Name(), b.Type().Name())
}

// NumMul multiplies two numeric Values.
func NumMul(a, b Value) (Value, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			r, ok := checkedMulInt(av, bv)
			if !ok {
				return NIL, fmt.Errorf("integer overflow")
			}
			return MakeInt(int(r)), nil
		case Float:
			return Float(float64(av) * float64(bv)), nil
		case *BigInt:
			r := new(big.Int).Mul(big.NewInt(int64(av)), bv.val)
			return NewBigInt(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return Float(float64(av) * float64(bv)), nil
		case Float:
			return Float(float64(av) * float64(bv)), nil
		case *BigInt:
			bf, _ := bv.val.Float64()
			return Float(float64(av) * bf), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			r := new(big.Int).Mul(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case Float:
			af, _ := av.val.Float64()
			return Float(af * float64(bv)), nil
		case *BigInt:
			r := new(big.Int).Mul(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	if r, err := numBinOpFallback(a, b, ratMul, bdMul, floatMul); err == nil {
		return r, nil
	}
	return NIL, fmt.Errorf("cannot multiply %s and %s", a.Type().Name(), b.Type().Name())
}

// NumDiv divides a by b. Int/Int returns Float when not exact.
func NumDiv(a, b Value) (Value, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	if _, ok := a.(*BigDecimal); ok {
		if bf, ok := ToFloat(b); ok && (math.IsNaN(bf) || math.IsInf(bf, 0)) {
			af, _ := ToFloat(a)
			return Float(floatDiv(af, bf)), nil
		}
	}
	if _, ok := b.(*BigDecimal); ok {
		if af, ok := ToFloat(a); ok && (math.IsNaN(af) || math.IsInf(af, 0)) {
			bf, _ := ToFloat(b)
			return Float(floatDiv(af, bf)), nil
		}
	}
	if _, ok := a.(*BigDecimal); ok {
		if bf, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			return Float(floatDiv(af, float64(bf))), nil
		}
	}
	if _, ok := b.(*BigDecimal); ok {
		if af, ok := a.(Float); ok {
			bf, _ := ToFloat(b)
			return Float(floatDiv(float64(af), bf)), nil
		}
	}
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Rat).SetFrac64(int64(av), int64(bv))
			return MaybeSimplifyRatio(r), nil
		case Float:
			// Int/Float: IEEE 754 semantics (allows Inf)
			return Float(float64(av) / float64(bv)), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Rat).SetFrac(big.NewInt(int64(av)), bv.val)
			return MaybeSimplifyRatio(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			// Float/Int: IEEE 754 semantics (allows Inf for /0)
			return Float(float64(av) / float64(bv)), nil
		case Float:
			// Float/Float: IEEE 754 semantics (allows NaN, Inf)
			return Float(float64(av) / float64(bv)), nil
		case *BigInt:
			bf, _ := bv.val.Float64()
			return Float(float64(av) / bf), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Rat).SetFrac(av.val, big.NewInt(int64(bv)))
			return MaybeSimplifyRatio(r), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			af, _ := av.val.Float64()
			r := af / float64(bv)
			return Float(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Rat).SetFrac(new(big.Int).Set(av.val), new(big.Int).Set(bv.val))
			return MaybeSimplifyRatio(r), nil
		}
	}
	if r, err := numBinOpFallback(a, b, ratDiv, bdDiv, floatDiv); err == nil {
		return r, nil
	}
	return NIL, fmt.Errorf("cannot divide %s and %s", a.Type().Name(), b.Type().Name())
}

// NumQuot performs integer division (quot in Clojure).
func NumQuot(a, b Value) (Value, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return MakeInt(int(av) / int(bv)), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			q, err := numQuotFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(q), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Quo(big.NewInt(int64(av)), bv.val)
			return NewBigInt(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			q, err := numQuotFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(q), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			q, err := numQuotFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(q), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Quo(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Quo(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	if IsRatio(a) || IsRatio(b) {
		if _, ok := a.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			q, err := numQuotFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(q), nil
		}
		if _, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			q, err := numQuotFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(q), nil
		}
		ar, aok := ToRat(a)
		br, bok := ToRat(b)
		if aok && bok {
			if br.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return NewBigInt(numTruncatedRatQuot(ar, br)), nil
		}
	}
	if _, ok := a.(*BigDecimal); ok {
		if _, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			q, err := numQuotFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(q), nil
		}
	}
	if _, ok := b.(*BigDecimal); ok {
		if _, ok := a.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			q, err := numQuotFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(q), nil
		}
	}
	// BigDecimal fallback: convert to float, truncate, return BigDecimal
	if _, ok := a.(*BigDecimal); ok {
		af, _ := ToFloat(a)
		bf, _ := ToFloat(b)
		if bf == 0 {
			return NIL, fmt.Errorf("divide by zero")
		}
		q, err := numQuotFloat(af, bf)
		if err != nil {
			return NIL, err
		}
		return NewBigDecimalFromFloat64(q), nil
	}
	if _, ok := b.(*BigDecimal); ok {
		af, _ := ToFloat(a)
		bf, _ := ToFloat(b)
		if bf == 0 {
			return NIL, fmt.Errorf("divide by zero")
		}
		q, err := numQuotFloat(af, bf)
		if err != nil {
			return NIL, err
		}
		return NewBigDecimalFromFloat64(q), nil
	}
	return NIL, fmt.Errorf("cannot quot %s and %s", a.Type().Name(), b.Type().Name())
}

func numQuotFloat(a, b float64) (float64, error) {
	if b == 0 || math.IsInf(a, 0) || math.IsNaN(a) || math.IsNaN(b) {
		return 0, fmt.Errorf("divide by zero")
	}
	if math.IsInf(b, 0) {
		return 0, nil
	}
	return math.Trunc(a / b), nil
}

// NumRem computes the remainder of truncated division (like Java's %).
// Sign follows the dividend: (rem -10 3) => -1
func NumRem(a, b Value) (Value, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return MakeInt(int(av) % int(bv)), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numRemFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Rem(big.NewInt(int64(av)), bv.val)
			return NewBigInt(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numRemFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numRemFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Rem(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Rem(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	if IsRatio(a) || IsRatio(b) {
		if _, ok := a.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numRemFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
		if _, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numRemFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
		ar, aok := ToRat(a)
		br, bok := ToRat(b)
		if aok && bok {
			if br.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return numRatRemainderValue(numRatRem(ar, br)), nil
		}
	}
	// Clojure returns a double for mixed BigDecimal/double remainders.
	if _, ok := a.(*BigDecimal); ok {
		if _, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numRemFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
	}
	if _, ok := b.(*BigDecimal); ok {
		if _, ok := a.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numRemFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
	}
	// BigDecimal fallback
	if _, ok := a.(*BigDecimal); ok {
		af, _ := ToFloat(a)
		bf, _ := ToFloat(b)
		if bf == 0 {
			return NIL, fmt.Errorf("divide by zero")
		}
		return NewBigDecimalFromFloat64(math.Mod(af, bf)), nil
	}
	if _, ok := b.(*BigDecimal); ok {
		af, _ := ToFloat(a)
		bf, _ := ToFloat(b)
		if bf == 0 {
			return NIL, fmt.Errorf("divide by zero")
		}
		return NewBigDecimalFromFloat64(math.Mod(af, bf)), nil
	}
	return NIL, fmt.Errorf("cannot rem %s and %s", a.Type().Name(), b.Type().Name())
}

func numRemFloat(a, b float64) (float64, error) {
	if b == 0 || math.IsInf(a, 0) || math.IsNaN(a) || math.IsNaN(b) {
		return 0, fmt.Errorf("divide by zero")
	}
	if math.IsInf(b, 0) {
		return math.NaN(), nil
	}
	return math.Mod(a, b), nil
}

func numTruncatedRatQuot(a, b *big.Rat) *big.Int {
	quot := new(big.Rat).Quo(a, b)
	return new(big.Int).Quo(quot.Num(), quot.Denom())
}

func numRatRem(a, b *big.Rat) *big.Rat {
	qRat := new(big.Rat).SetInt(numTruncatedRatQuot(a, b))
	return new(big.Rat).Sub(a, new(big.Rat).Mul(qRat, b))
}

func numRatRemainderValue(r *big.Rat) Value {
	if r.IsInt() {
		return NewBigInt(new(big.Int).Set(r.Num()))
	}
	return NewRatio(r)
}

// NumMod computes the floored modulus (like Clojure's mod).
// Sign follows the divisor: (mod -10 3) => 2
func NumMod(a, b Value) (Value, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := int(av) % int(bv)
			if r != 0 && (r > 0) != (int(bv) > 0) {
				r += int(bv)
			}
			return MakeInt(r), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numModFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			ai := big.NewInt(int64(av))
			r := new(big.Int).Rem(ai, bv.val)
			if r.Sign() != 0 && (r.Sign() > 0) != (bv.val.Sign() > 0) {
				r.Add(r, bv.val)
			}
			return NewBigInt(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numModFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numModFloat(float64(av), float64(bv))
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			bi := big.NewInt(int64(bv))
			r := new(big.Int).Rem(av.val, bi)
			if r.Sign() != 0 && (r.Sign() > 0) != (bi.Sign() > 0) {
				r.Add(r, bi)
			}
			return NewBigInt(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Rem(av.val, bv.val)
			if r.Sign() != 0 && (r.Sign() > 0) != (bv.val.Sign() > 0) {
				r.Add(r, bv.val)
			}
			return NewBigInt(r), nil
		}
	}
	if IsRatio(a) || IsRatio(b) {
		if _, ok := a.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numModFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
		if _, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numModFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
		ar, aok := ToRat(a)
		br, bok := ToRat(b)
		if aok && bok {
			if br.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := numRatRem(ar, br)
			if r.Sign() != 0 && (r.Sign() > 0) != (br.Sign() > 0) {
				r.Add(r, br)
			}
			return numRatRemainderValue(r), nil
		}
	}
	if _, ok := a.(*BigDecimal); ok {
		if _, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numModFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
	}
	if _, ok := b.(*BigDecimal); ok {
		if _, ok := a.(Float); ok {
			af, _ := ToFloat(a)
			bf, _ := ToFloat(b)
			if bf == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r, err := numModFloat(af, bf)
			if err != nil {
				return NIL, err
			}
			return Float(r), nil
		}
	}
	// BigDecimal fallback
	if _, ok := a.(*BigDecimal); ok {
		af, _ := ToFloat(a)
		bf, _ := ToFloat(b)
		if bf == 0 {
			return NIL, fmt.Errorf("divide by zero")
		}
		r, err := numModFloat(af, bf)
		if err != nil {
			return NIL, err
		}
		return NewBigDecimalFromFloat64(r), nil
	}
	if _, ok := b.(*BigDecimal); ok {
		af, _ := ToFloat(a)
		bf, _ := ToFloat(b)
		if bf == 0 {
			return NIL, fmt.Errorf("divide by zero")
		}
		r, err := numModFloat(af, bf)
		if err != nil {
			return NIL, err
		}
		return NewBigDecimalFromFloat64(r), nil
	}
	return NIL, fmt.Errorf("cannot mod %s and %s", a.Type().Name(), b.Type().Name())
}

func numModFloat(a, b float64) (float64, error) {
	r, err := numRemFloat(a, b)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(r) {
		return r, nil
	}
	if r != 0 && (r > 0) != (b > 0) {
		r += b
	}
	return r, nil
}

// NumNeg negates a numeric value.
func NumNeg(a Value) (Value, error) {
	a = normalizeFloat32(a)
	switch av := a.(type) {
	case Int:
		return MakeInt(-int(av)), nil
	case Float:
		return Float(-float64(av)), nil
	case *BigInt:
		return MaybeDowngrade(new(big.Int).Neg(av.val)), nil
	case *Ratio:
		return NewRatio(new(big.Rat).Neg(av.val)), nil
	case *BigDecimal:
		return NewBigDecimal(new(big.Float).SetPrec(bigDecimalPrec).Neg(av.val)), nil
	}
	return NIL, fmt.Errorf("cannot negate %s", a.Type().Name())
}

// NumAbs returns absolute value.
func NumAbs(a Value) (Value, error) {
	a = normalizeFloat32(a)
	switch av := a.(type) {
	case Int:
		v := int(av)
		if v < 0 {
			v = -v
		}
		return MakeInt(v), nil
	case Float:
		v := float64(av)
		if v < 0 {
			v = -v
		}
		return Float(v), nil
	case *BigInt:
		return MaybeDowngrade(new(big.Int).Abs(av.val)), nil
	case *Ratio:
		return NewRatio(new(big.Rat).Abs(av.val)), nil
	case *BigDecimal:
		return NewBigDecimal(new(big.Float).SetPrec(bigDecimalPrec).Abs(av.val)), nil
	}
	return NIL, fmt.Errorf("cannot abs %s", a.Type().Name())
}

// --- Comparison helpers ---

// NumGt returns true if a > b.
func NumGt(a, b Value) (bool, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return int(av) > int(bv), nil
		case Float:
			return float64(av) > float64(bv), nil
		case *BigInt:
			return big.NewInt(int64(av)).Cmp(bv.val) > 0, nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) > float64(int(bv)), nil
		case Float:
			return float64(av) > float64(bv), nil
		case *BigInt:
			bf, _ := new(big.Float).SetInt(bv.val).Float64()
			return float64(av) > bf, nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			return av.val.Cmp(big.NewInt(int64(bv))) > 0, nil
		case Float:
			af, _ := new(big.Float).SetInt(av.val).Float64()
			return af > float64(bv), nil
		case *BigInt:
			return av.val.Cmp(bv.val) > 0, nil
		}
	}
	if c, err := numCmpFallback(a, b); err == nil {
		return c > 0, nil
	}
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

// NumLt returns true if a < b.
func NumLt(a, b Value) (bool, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return int(av) < int(bv), nil
		case Float:
			return float64(av) < float64(bv), nil
		case *BigInt:
			return big.NewInt(int64(av)).Cmp(bv.val) < 0, nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) < float64(int(bv)), nil
		case Float:
			return float64(av) < float64(bv), nil
		case *BigInt:
			bf, _ := new(big.Float).SetInt(bv.val).Float64()
			return float64(av) < bf, nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			return av.val.Cmp(big.NewInt(int64(bv))) < 0, nil
		case Float:
			af, _ := new(big.Float).SetInt(av.val).Float64()
			return af < float64(bv), nil
		case *BigInt:
			return av.val.Cmp(bv.val) < 0, nil
		}
	}
	if c, err := numCmpFallback(a, b); err == nil {
		return c < 0, nil
	}
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

func NumGe(a, b Value) (bool, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return int(av) >= int(bv), nil
		case Float:
			return float64(av) >= float64(bv), nil
		case *BigInt:
			return big.NewInt(int64(av)).Cmp(bv.val) >= 0, nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) >= float64(int(bv)), nil
		case Float:
			return float64(av) >= float64(bv), nil
		case *BigInt:
			bf, _ := new(big.Float).SetInt(bv.val).Float64()
			return float64(av) >= bf, nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			return av.val.Cmp(big.NewInt(int64(bv))) >= 0, nil
		case Float:
			af, _ := new(big.Float).SetInt(av.val).Float64()
			return af >= float64(bv), nil
		case *BigInt:
			return av.val.Cmp(bv.val) >= 0, nil
		}
	}
	if c, err := numCmpFallback(a, b); err == nil {
		return c >= 0, nil
	}
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

func NumLe(a, b Value) (bool, error) {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return int(av) <= int(bv), nil
		case Float:
			return float64(av) <= float64(bv), nil
		case *BigInt:
			return big.NewInt(int64(av)).Cmp(bv.val) <= 0, nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) <= float64(int(bv)), nil
		case Float:
			return float64(av) <= float64(bv), nil
		case *BigInt:
			bf, _ := new(big.Float).SetInt(bv.val).Float64()
			return float64(av) <= bf, nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			return av.val.Cmp(big.NewInt(int64(bv))) <= 0, nil
		case Float:
			af, _ := new(big.Float).SetInt(av.val).Float64()
			return af <= float64(bv), nil
		case *BigInt:
			return av.val.Cmp(bv.val) <= 0, nil
		}
	}
	if c, err := numCmpFallback(a, b); err == nil {
		return c <= 0, nil
	}
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

// NumEq tests Clojure = numeric equality. Exact integers compare across
// fixed/big integer representations, floating values compare by value, and
// BigDecimal only compares with BigDecimal.
func NumEq(a, b Value) bool {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return int(av) == int(bv)
		case *BigInt:
			return big.NewInt(int64(av)).Cmp(bv.val) == 0
		}
	case Float:
		switch bv := b.(type) {
		case Float:
			return float64(av) == float64(bv)
		case Float32:
			return float64(av) == float64(bv)
		}
	case Float32:
		switch bv := b.(type) {
		case Float:
			return float64(av) == float64(bv)
		case Float32:
			return float64(av) == float64(bv)
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			return av.val.Cmp(big.NewInt(int64(bv))) == 0
		case *BigInt:
			return av.val.Cmp(bv.val) == 0
		}
	case *Ratio:
		if bv, ok := b.(*Ratio); ok {
			return av.val.Cmp(bv.val) == 0
		}
	case *BigDecimal:
		if bv, ok := b.(*BigDecimal); ok {
			return av.val.Cmp(bv.val) == 0
		}
	}
	return false
}

// NumEquivalent tests broad numeric equivalence for == and numeric predicates.
func NumEquivalent(a, b Value) bool {
	a = normalizeFloat32(a)
	b = normalizeFloat32(b)
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return int(av) == int(bv)
		case Float:
			return float64(av) == float64(bv)
		case *BigInt:
			return big.NewInt(int64(av)).Cmp(bv.val) == 0
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) == float64(int(bv))
		case Float:
			return float64(av) == float64(bv)
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			return av.val.Cmp(big.NewInt(int64(bv))) == 0
		case *BigInt:
			return av.val.Cmp(bv.val) == 0
		}
	}
	if c, err := numCmpFallback(a, b); err == nil {
		return c == 0
	}
	return false
}

// IsNumber returns true if the value is Int, Float, or BigInt.
func IsNumber(v Value) bool {
	switch v.(type) {
	case Int, Float, Float32, *BigInt, *Ratio, *BigDecimal:
		return true
	}
	return false
}

// ToFloat converts any numeric Value to float64.
func ToFloat(v Value) (float64, bool) {
	switch n := v.(type) {
	case Int:
		return float64(n), true
	case Float:
		return float64(n), true
	case Float32:
		return float64(n), true
	case *BigInt:
		f, _ := new(big.Float).SetInt(n.val).Float64()
		return f, true
	case *Ratio:
		f, _ := n.val.Float64()
		return f, true
	case *BigDecimal:
		f, _ := n.val.Float64()
		return f, true
	}
	return 0, false
}

func intFromFloat64(f float64) int {
	switch {
	case math.IsNaN(f):
		return 0
	case math.IsInf(f, 1) || f > float64(maxIntValue):
		return int(maxIntValue)
	case math.IsInf(f, -1) || f < float64(minIntValue):
		return int(minIntValue)
	default:
		return int(f)
	}
}

// ToInt converts a numeric Value to int if possible.
func ToInt(v Value) (int, bool) {
	switch n := v.(type) {
	case Int:
		return int(n), true
	case Float:
		return intFromFloat64(float64(n)), true
	case Float32:
		return intFromFloat64(float64(n)), true
	case *BigInt:
		if n.val.IsInt64() {
			return int(n.val.Int64()), true
		}
		return 0, false
	case Boolean:
		if bool(n) {
			return 1, true
		}
		return 0, true
	case *BigDecimal:
		f, _ := n.Val().Float64()
		return intFromFloat64(f), true
	case *Ratio:
		f, _ := n.Val().Float64()
		return intFromFloat64(f), true
	}
	return 0, false
}

// --- Ratio/BigDecimal fallback helpers ---

type ratOp func(a, b *big.Rat) *big.Rat
type bdOp func(a, b *big.Float) *big.Float
type floatOp func(a, b float64) float64

var (
	ratAdd   ratOp   = func(a, b *big.Rat) *big.Rat { return new(big.Rat).Add(a, b) }
	ratSub   ratOp   = func(a, b *big.Rat) *big.Rat { return new(big.Rat).Sub(a, b) }
	ratMul   ratOp   = func(a, b *big.Rat) *big.Rat { return new(big.Rat).Mul(a, b) }
	ratDiv   ratOp   = func(a, b *big.Rat) *big.Rat { return new(big.Rat).Quo(a, b) }
	bdAdd    bdOp    = func(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(bigDecimalPrec).Add(a, b) }
	bdSub    bdOp    = func(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(bigDecimalPrec).Sub(a, b) }
	bdMul    bdOp    = func(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(bigDecimalPrec).Mul(a, b) }
	bdDiv    bdOp    = func(a, b *big.Float) *big.Float { return new(big.Float).SetPrec(bigDecimalPrec).Quo(a, b) }
	floatAdd floatOp = func(a, b float64) float64 { return a + b }
	floatSub floatOp = func(a, b float64) float64 { return a - b }
	floatMul floatOp = func(a, b float64) float64 { return a * b }
	floatDiv floatOp = func(a, b float64) float64 { return a / b }
)

// numBinOpFallback handles Ratio and BigDecimal arithmetic.
// Promotion rules:
//   - Ratio ⊕ (Int|BigInt|Ratio) → Ratio (may simplify to Int)
//   - Ratio ⊕ Float → Float
//   - Ratio ⊕ BigDecimal → BigDecimal
//   - BigDecimal ⊕ (Int|BigInt|BigDecimal) → BigDecimal
//   - BigDecimal ⊕ Float → BigDecimal
func numBinOpFallback(a, b Value, rop ratOp, bop bdOp, fop floatOp) (Value, error) {
	// If either is BigDecimal, promote to BigDecimal
	if ad, ok := a.(*BigDecimal); ok {
		if bd, ok := b.(*BigDecimal); ok {
			return NewBigDecimal(bop(ad.val, bd.val)), nil
		}
		if bf, ok := ToFloat(b); ok {
			bv := new(big.Float).SetPrec(bigDecimalPrec).SetFloat64(bf)
			return NewBigDecimal(bop(ad.val, bv)), nil
		}
	}
	if bd, ok := b.(*BigDecimal); ok {
		if af, ok := ToFloat(a); ok {
			av := new(big.Float).SetPrec(bigDecimalPrec).SetFloat64(af)
			return NewBigDecimal(bop(av, bd.val)), nil
		}
	}
	// If either is Ratio, try Ratio path
	if _, ok := a.(*Ratio); ok {
		if _, ok := b.(Float); ok {
			af, _ := ToFloat(a)
			return Float(fop(af, float64(b.(Float)))), nil
		}
		ar, aok := ToRat(a)
		br, bok := ToRat(b)
		if aok && bok && rop != nil {
			return MaybeSimplifyRatio(rop(ar, br)), nil
		}
	}
	if _, ok := b.(*Ratio); ok {
		if _, ok := a.(Float); ok {
			bf, _ := ToFloat(b)
			return Float(fop(float64(a.(Float)), bf)), nil
		}
		ar, aok := ToRat(a)
		br, bok := ToRat(b)
		if aok && bok && rop != nil {
			return MaybeSimplifyRatio(rop(ar, br)), nil
		}
	}
	return NIL, fmt.Errorf("unsupported types")
}

// numCmpFallback handles Ratio and BigDecimal comparisons.
// Returns -1, 0, 1 like big.Rat.Cmp / big.Float.Cmp.
func numCmpFallback(a, b Value) (int, error) {
	// BigDecimal
	if ad, ok := a.(*BigDecimal); ok {
		if bd, ok := b.(*BigDecimal); ok {
			return ad.val.Cmp(bd.val), nil
		}
		bf, ok := ToFloat(b)
		if ok {
			bv := new(big.Float).SetPrec(bigDecimalPrec).SetFloat64(bf)
			return ad.val.Cmp(bv), nil
		}
	}
	if bd, ok := b.(*BigDecimal); ok {
		af, ok := ToFloat(a)
		if ok {
			av := new(big.Float).SetPrec(bigDecimalPrec).SetFloat64(af)
			return av.Cmp(bd.val), nil
		}
	}
	// Ratio
	if ar, ok := a.(*Ratio); ok {
		if br, ok := b.(*Ratio); ok {
			return ar.val.Cmp(br.val), nil
		}
		if _, ok := b.(Float); ok {
			af := ar.ToFloat64()
			bf := float64(b.(Float))
			if af < bf {
				return -1, nil
			}
			if af > bf {
				return 1, nil
			}
			return 0, nil
		}
		if br, ok := ToRat(b); ok {
			return ar.val.Cmp(br), nil
		}
	}
	if br, ok := b.(*Ratio); ok {
		if _, ok := a.(Float); ok {
			af := float64(a.(Float))
			bf := br.ToFloat64()
			if af < bf {
				return -1, nil
			}
			if af > bf {
				return 1, nil
			}
			return 0, nil
		}
		if ar, ok := ToRat(a); ok {
			return ar.Cmp(br.val), nil
		}
	}
	return 0, fmt.Errorf("unsupported types")
}
