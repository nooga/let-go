package vm

import (
	"fmt"
	"math"
	"math/big"
)

// Numeric type promotion and arithmetic dispatch.
// All functions use direct type assertions (no Unbox) to avoid allocation.
// Promotion: Int op BigInt → BigInt, BigInt op Float → Float, Int op Float → Float.

// NumAdd adds two numeric Values.
func NumAdd(a, b Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return MakeInt(int(av) + int(bv)), nil
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
			bf := new(big.Float).SetInt(bv.val)
			af := new(big.Float).SetFloat64(float64(av))
			r, _ := new(big.Float).Add(af, bf).Float64()
			return Float(r), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			r := new(big.Int).Add(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case Float:
			af := new(big.Float).SetInt(av.val)
			bf := new(big.Float).SetFloat64(float64(bv))
			r, _ := new(big.Float).Add(af, bf).Float64()
			return Float(r), nil
		case *BigInt:
			r := new(big.Int).Add(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	return NIL, fmt.Errorf("cannot add %s and %s", a.Type().Name(), b.Type().Name())
}

// NumSub subtracts b from a.
func NumSub(a, b Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return MakeInt(int(av) - int(bv)), nil
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
			af := new(big.Float).SetFloat64(float64(av))
			bf := new(big.Float).SetInt(bv.val)
			r, _ := new(big.Float).Sub(af, bf).Float64()
			return Float(r), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			r := new(big.Int).Sub(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case Float:
			af := new(big.Float).SetInt(av.val)
			bf := new(big.Float).SetFloat64(float64(bv))
			r, _ := new(big.Float).Sub(af, bf).Float64()
			return Float(r), nil
		case *BigInt:
			r := new(big.Int).Sub(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	return NIL, fmt.Errorf("cannot subtract %s and %s", a.Type().Name(), b.Type().Name())
}

// NumMul multiplies two numeric Values.
func NumMul(a, b Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return MakeInt(int(av) * int(bv)), nil
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
			af := new(big.Float).SetFloat64(float64(av))
			bf := new(big.Float).SetInt(bv.val)
			r, _ := new(big.Float).Mul(af, bf).Float64()
			return Float(r), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			r := new(big.Int).Mul(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case Float:
			af := new(big.Float).SetInt(av.val)
			bf := new(big.Float).SetFloat64(float64(bv))
			r, _ := new(big.Float).Mul(af, bf).Float64()
			return Float(r), nil
		case *BigInt:
			r := new(big.Int).Mul(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	return NIL, fmt.Errorf("cannot multiply %s and %s", a.Type().Name(), b.Type().Name())
}

// NumDiv divides a by b. Int/Int returns Float when not exact.
func NumDiv(a, b Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			if int(av)%int(bv) == 0 {
				return MakeInt(int(av) / int(bv)), nil
			}
			return Float(float64(av) / float64(bv)), nil
		case Float:
			// Int/Float: IEEE 754 semantics (allows Inf)
			return Float(float64(av) / float64(bv)), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			ai := big.NewInt(int64(av))
			mod := new(big.Int).Mod(ai, bv.val)
			if mod.Sign() == 0 {
				r := new(big.Int).Div(ai, bv.val)
				return NewBigInt(r), nil
			}
			af := new(big.Float).SetInt(ai)
			bf := new(big.Float).SetInt(bv.val)
			r, _ := new(big.Float).Quo(af, bf).Float64()
			return Float(r), nil
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
			if bv.val.Sign() == 0 {
				// Float/BigInt(0): IEEE 754
				return Float(float64(av) / 0.0), nil
			}
			af := new(big.Float).SetFloat64(float64(av))
			bf := new(big.Float).SetInt(bv.val)
			r, _ := new(big.Float).Quo(af, bf).Float64()
			return Float(r), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			bi := big.NewInt(int64(bv))
			mod := new(big.Int).Mod(av.val, bi)
			if mod.Sign() == 0 {
				r := new(big.Int).Div(av.val, bi)
				return NewBigInt(r), nil
			}
			af := new(big.Float).SetInt(av.val)
			bf := new(big.Float).SetInt64(int64(bv))
			r, _ := new(big.Float).Quo(af, bf).Float64()
			return Float(r), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			af := new(big.Float).SetInt(av.val)
			bf := new(big.Float).SetFloat64(float64(bv))
			r, _ := new(big.Float).Quo(af, bf).Float64()
			return Float(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			mod := new(big.Int).Mod(av.val, bv.val)
			if mod.Sign() == 0 {
				r := new(big.Int).Div(av.val, bv.val)
				return NewBigInt(r), nil
			}
			af := new(big.Float).SetInt(av.val)
			bf := new(big.Float).SetInt(bv.val)
			r, _ := new(big.Float).Quo(af, bf).Float64()
			return Float(r), nil
		}
	}
	return NIL, fmt.Errorf("cannot divide %s and %s", a.Type().Name(), b.Type().Name())
}

// NumQuot performs integer division (quot in Clojure).
func NumQuot(a, b Value) (Value, error) {
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
			return Float(float64(int(float64(av) / float64(bv)))), nil
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
			return Float(float64(int(float64(av) / float64(bv)))), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return Float(float64(int(float64(av) / float64(bv)))), nil
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
	return NIL, fmt.Errorf("cannot quot %s and %s", a.Type().Name(), b.Type().Name())
}

// NumRem computes the remainder of truncated division (like Java's %).
// Sign follows the dividend: (rem -10 3) => -1
func NumRem(a, b Value) (Value, error) {
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
			return Float(math.Remainder(float64(av), float64(bv))), nil
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
			return Float(math.Remainder(float64(av), float64(bv))), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return Float(math.Remainder(float64(av), float64(bv))), nil
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
	return NIL, fmt.Errorf("cannot rem %s and %s", a.Type().Name(), b.Type().Name())
}

// NumMod computes the floored modulus (like Clojure's mod).
// Sign follows the divisor: (mod -10 3) => 2
func NumMod(a, b Value) (Value, error) {
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
			r := math.Mod(float64(av), float64(bv))
			if r != 0 && (r > 0) != (float64(bv) > 0) {
				r += float64(bv)
			}
			return Float(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Mod(big.NewInt(int64(av)), bv.val)
			return NewBigInt(r), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := math.Mod(float64(av), float64(bv))
			if r != 0 && (r > 0) != (float64(bv) > 0) {
				r += float64(bv)
			}
			return Float(r), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := math.Mod(float64(av), float64(bv))
			if r != 0 && (r > 0) != (float64(bv) > 0) {
				r += float64(bv)
			}
			return Float(r), nil
		}
	case *BigInt:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Mod(av.val, big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		case *BigInt:
			if bv.val.Sign() == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			r := new(big.Int).Mod(av.val, bv.val)
			return NewBigInt(r), nil
		}
	}
	return NIL, fmt.Errorf("cannot mod %s and %s", a.Type().Name(), b.Type().Name())
}

// NumNeg negates a numeric value.
func NumNeg(a Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		return MakeInt(-int(av)), nil
	case Float:
		return Float(-float64(av)), nil
	case *BigInt:
		return MaybeDowngrade(new(big.Int).Neg(av.val)), nil
	}
	return NIL, fmt.Errorf("cannot negate %s", a.Type().Name())
}

// NumAbs returns absolute value.
func NumAbs(a Value) (Value, error) {
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
	}
	return NIL, fmt.Errorf("cannot abs %s", a.Type().Name())
}

// --- Comparison helpers ---


// NumGt returns true if a > b.
func NumGt(a, b Value) (bool, error) {
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
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

// NumLt returns true if a < b.
func NumLt(a, b Value) (bool, error) {
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
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

func NumGe(a, b Value) (bool, error) {
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
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

func NumLe(a, b Value) (bool, error) {
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
	return false, fmt.Errorf("cannot compare %s and %s", a.Type().Name(), b.Type().Name())
}

// NumEq tests numeric equality (cross-type: 1 == 1.0 is true).
func NumEq(a, b Value) bool {
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
	return false
}

// IsNumber returns true if the value is Int, Float, or BigInt.
func IsNumber(v Value) bool {
	switch v.(type) {
	case Int, Float, *BigInt:
		return true
	}
	return false
}

// ToFloat converts an Int, Float, or BigInt to float64.
func ToFloat(v Value) (float64, bool) {
	switch n := v.(type) {
	case Int:
		return float64(n), true
	case Float:
		return float64(n), true
	case *BigInt:
		f, _ := new(big.Float).SetInt(n.val).Float64()
		return f, true
	}
	return 0, false
}

// ToInt converts a numeric Value to int if possible.
func ToInt(v Value) (int, bool) {
	switch n := v.(type) {
	case Int:
		return int(n), true
	case Float:
		return int(n), true
	case *BigInt:
		if n.val.IsInt64() {
			return int(n.val.Int64()), true
		}
		return 0, false
	}
	return 0, false
}
