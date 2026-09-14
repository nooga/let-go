package vm

import (
	"math"
	"math/big"
)

// Overflow-checked integer primitives at the host's `int` width.
//
// Go's signed integer overflow is undefined-style silent wrap. These helpers
// detect overflow without allocating, returning a flag the caller uses to
// fall back to *big.Int. They power the apostrophe arithmetic operations
// (`+'`, `-'`, `*'`) which promote to BigInt on overflow rather than
// silently wrapping like let-go's current default `+`/`-`/`*`.
//
// The XOR-sign and divide-back tricks below are width-agnostic, so on a
// 32-bit host they detect overflow at the 32-bit boundary; on a 64-bit
// host, at the 64-bit boundary. Either way matches `vm.Int`'s representable
// range exactly, which is what `+'` should promote at.

// addIntChecked returns (a+b, overflow).
// Overflow iff a and b share a sign but the sum's sign differs from theirs.
func addIntChecked(a, b int64) (int64, bool) {
	sum := a + b
	return sum, (a^sum)&(b^sum) < 0
}

// subIntChecked returns (a-b, overflow).
// Overflow iff a and b have different signs and the result's sign differs from a.
func subIntChecked(a, b int64) (int64, bool) {
	diff := a - b
	return diff, (a^b)&(a^diff) < 0
}

// mulIntChecked returns (a*b, overflow).
// Uses divide-back to detect overflow, with explicit MinInt*-1 special case
// (which would otherwise wrap silently during the divide-back check).
func mulIntChecked(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, false
	}
	if (a == math.MinInt64 && b == -1) || (b == math.MinInt64 && a == -1) {
		return 0, true
	}
	p := a * b
	if p/b != a {
		return 0, true
	}
	return p, false
}

// negIntChecked returns (-a, overflow). Only math.MinInt64 overflows.
func negIntChecked(a int64) (int64, bool) {
	if a == math.MinInt64 {
		return 0, true
	}
	return -a, false
}

// NumAddP is `+` with overflow promotion for Int+Int.
// All other type combinations match NumAdd exactly. Delegating to NumAdd
// for the non-Int/Int paths is safe because NumAdd's BigInt branches
// return NewBigInt(...) directly (no MaybeDowngrade), preserving BigInt
// identity through `+'`. The same invariant is relied on by NumSubP and
// NumMulP below — if you change MaybeDowngrade usage in NumAdd/Sub/Mul,
// re-check the eval tests `(bigint? (+' 1 1N))` etc. NumNeg DOES
// downgrade *BigInt, hence the explicit branch in NumNegP.
// Naming mirrors Clojure JVM's `addP`/`multiplyP` precedent.
func NumAddP(a, b Value) (Value, error) {
	if av, ok := a.(Int); ok {
		if bv, ok := b.(Int); ok {
			sum, overflow := addIntChecked(int64(av), int64(bv))
			if !overflow {
				return MakeInt64(sum), nil
			}
			r := new(big.Int).Add(big.NewInt(int64(av)), big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		}
	}
	return NumAdd(a, b)
}

// NumSubP is `-` with overflow promotion for Int-Int.
func NumSubP(a, b Value) (Value, error) {
	if av, ok := a.(Int); ok {
		if bv, ok := b.(Int); ok {
			diff, overflow := subIntChecked(int64(av), int64(bv))
			if !overflow {
				return MakeInt64(diff), nil
			}
			r := new(big.Int).Sub(big.NewInt(int64(av)), big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		}
	}
	return NumSub(a, b)
}

// NumMulP is `*` with overflow promotion for Int*Int.
func NumMulP(a, b Value) (Value, error) {
	if av, ok := a.(Int); ok {
		if bv, ok := b.(Int); ok {
			prod, overflow := mulIntChecked(int64(av), int64(bv))
			if !overflow {
				return MakeInt64(prod), nil
			}
			r := new(big.Int).Mul(big.NewInt(int64(av)), big.NewInt(int64(bv)))
			return NewBigInt(r), nil
		}
	}
	return NumMul(a, b)
}

// NumNegP is unary `-` with overflow promotion at math.MinInt64.
// For *BigInt we negate directly without downgrading — `(-' 1N)` must
// return `-1N`, not `Int -1` (Clojure preserves BigInt identity through
// promoting ops). NumNeg in pkg/vm/numbers.go does downgrade, hence the
// explicit BigInt branch here.
func NumNegP(a Value) (Value, error) {
	switch av := a.(type) {
	case Int:
		neg, overflow := negIntChecked(int64(av))
		if !overflow {
			return MakeInt64(neg), nil
		}
		return NewBigInt(new(big.Int).Neg(big.NewInt(int64(av)))), nil
	case *BigInt:
		return NewBigInt(new(big.Int).Neg(av.val)), nil
	}
	return NumNeg(a)
}
