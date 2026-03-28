package vm

import "fmt"

// Numeric type promotion and arithmetic dispatch.
// All functions use direct type assertions (no Unbox) to avoid allocation.

// NumAdd adds two numeric Values.
func NumAdd(a, b Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return MakeInt(int(av) + int(bv)), nil
		case Float:
			return Float(float64(av) + float64(bv)), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return Float(float64(av) + float64(bv)), nil
		case Float:
			return Float(float64(av) + float64(bv)), nil
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
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return Float(float64(av) - float64(bv)), nil
		case Float:
			return Float(float64(av) - float64(bv)), nil
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
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return Float(float64(av) * float64(bv)), nil
		case Float:
			return Float(float64(av) * float64(bv)), nil
		}
	}
	return NIL, fmt.Errorf("cannot multiply %s and %s", a.Type().Name(), b.Type().Name())
}

// NumDiv divides a by b. Int/Int returns Float when not exact.
// Use NumQuot for integer division.
func NumDiv(a, b Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			// Integer division when exact, otherwise promote to float
			if int(av)%int(bv) == 0 {
				return MakeInt(int(av) / int(bv)), nil
			}
			return Float(float64(av) / float64(bv)), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return Float(float64(av) / float64(bv)), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return Float(float64(av) / float64(bv)), nil
		case Float:
			if float64(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return Float(float64(av) / float64(bv)), nil
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
	}
	return NIL, fmt.Errorf("cannot quot %s and %s", a.Type().Name(), b.Type().Name())
}

// NumMod computes modulus.
func NumMod(a, b Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			if int(bv) == 0 {
				return NIL, fmt.Errorf("divide by zero")
			}
			return MakeInt(int(av) % int(bv)), nil
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
	}
	return NIL, fmt.Errorf("cannot abs %s", a.Type().Name())
}

// NumGt returns true if a > b.
func NumGt(a, b Value) (bool, error) {
	switch av := a.(type) {
	case Int:
		switch bv := b.(type) {
		case Int:
			return int(av) > int(bv), nil
		case Float:
			return float64(av) > float64(bv), nil
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) > float64(int(bv)), nil
		case Float:
			return float64(av) > float64(bv), nil
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
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) < float64(int(bv)), nil
		case Float:
			return float64(av) < float64(bv), nil
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
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) >= float64(int(bv)), nil
		case Float:
			return float64(av) >= float64(bv), nil
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
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) <= float64(int(bv)), nil
		case Float:
			return float64(av) <= float64(bv), nil
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
		}
	case Float:
		switch bv := b.(type) {
		case Int:
			return float64(av) == float64(int(bv))
		case Float:
			return float64(av) == float64(bv)
		}
	}
	return false
}

// IsNumber returns true if the value is Int or Float.
func IsNumber(v Value) bool {
	switch v.(type) {
	case Int, Float:
		return true
	}
	return false
}

// ToFloat converts an Int or Float to float64.
func ToFloat(v Value) (float64, bool) {
	switch n := v.(type) {
	case Int:
		return float64(n), true
	case Float:
		return float64(n), true
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
	}
	return 0, false
}
