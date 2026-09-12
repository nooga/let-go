package vm

import (
	"math"
	"math/big"
	"testing"
)

// --- Primitive overflow checks -----------------------------------------------
//
// All cases use math.MaxInt64 / math.MinInt64: Int is 64 bits on every
// target, so the overflow primitives are too.

func TestAddIntChecked(t *testing.T) {
	cases := []struct {
		a, b     int64
		want     int64
		overflow bool
	}{
		{1, 2, 3, false},
		{-1, 1, 0, false},
		{math.MaxInt64, 0, math.MaxInt64, false},
		{math.MinInt64, 0, math.MinInt64, false},
		{math.MaxInt64, 1, 0, true},
		{1, math.MaxInt64, 0, true},
		{math.MinInt64, -1, 0, true},
		{-1, math.MinInt64, 0, true},
		{math.MaxInt64, math.MaxInt64, 0, true},
		{math.MinInt64, math.MinInt64, 0, true},
		{math.MaxInt64, math.MinInt64, -1, false},
	}
	for _, c := range cases {
		got, of := addIntChecked(c.a, c.b)
		if of != c.overflow {
			t.Errorf("addIntChecked(%d, %d) overflow=%v, want %v", c.a, c.b, of, c.overflow)
		}
		if !of && got != c.want {
			t.Errorf("addIntChecked(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestSubIntChecked(t *testing.T) {
	cases := []struct {
		a, b     int64
		want     int64
		overflow bool
	}{
		{3, 2, 1, false},
		{1, 1, 0, false},
		{math.MinInt64, 1, 0, true},
		{math.MinInt64, -1, math.MinInt64 + 1, false}, // MinInt - -1 = MinInt+1, in range
		{math.MinInt64, math.MaxInt64, 0, true},
		{math.MaxInt64, -1, 0, true},
		{math.MaxInt64, math.MinInt64, 0, true},
		{0, math.MinInt64, 0, true}, // -(MinInt) overflows
		{-1, math.MaxInt64, math.MinInt64, false},
	}
	for _, c := range cases {
		got, of := subIntChecked(c.a, c.b)
		if of != c.overflow {
			t.Errorf("subIntChecked(%d, %d) overflow=%v, want %v", c.a, c.b, of, c.overflow)
		}
		if !of && got != c.want {
			t.Errorf("subIntChecked(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestMulIntChecked(t *testing.T) {
	cases := []struct {
		a, b     int64
		want     int64
		overflow bool
	}{
		{0, 0, 0, false},
		{1, 1, 1, false},
		{-1, -1, 1, false},
		{math.MinInt64, 0, 0, false},
		{math.MinInt64, 1, math.MinInt64, false},
		{math.MinInt64, -1, 0, true},
		{-1, math.MinInt64, 0, true},
		{math.MaxInt64, 1, math.MaxInt64, false},
		{math.MaxInt64, -1, math.MinInt64 + 1, false},
		{math.MaxInt64, 2, 0, true},
		{2, math.MaxInt64, 0, true},
		{math.MaxInt64, math.MaxInt64, 0, true},
		{math.MinInt64 / 2, 3, 0, true},
		{3, math.MinInt64 / 2, 0, true},
	}
	for _, c := range cases {
		got, of := mulIntChecked(c.a, c.b)
		if of != c.overflow {
			t.Errorf("mulIntChecked(%d, %d) overflow=%v, want %v", c.a, c.b, of, c.overflow)
		}
		if !of && got != c.want {
			t.Errorf("mulIntChecked(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestNegIntChecked(t *testing.T) {
	cases := []struct {
		a        int64
		want     int64
		overflow bool
	}{
		{0, 0, false},
		{1, -1, false},
		{-1, 1, false},
		{math.MaxInt64, -math.MaxInt64, false},
		{math.MinInt64, 0, true},
	}
	for _, c := range cases {
		got, of := negIntChecked(c.a)
		if of != c.overflow {
			t.Errorf("negIntChecked(%d) overflow=%v, want %v", c.a, of, c.overflow)
		}
		if !of && got != c.want {
			t.Errorf("negIntChecked(%d) = %d, want %d", c.a, got, c.want)
		}
	}
}

// --- High-level NumXP behaviour ---------------------------------------------

func mustBigIntString(t *testing.T, v Value, want string) {
	t.Helper()
	bi, ok := v.(*BigInt)
	if !ok {
		t.Fatalf("expected *BigInt, got %T (%v)", v, v)
	}
	if bi.val.String() != want {
		t.Errorf("BigInt value = %s, want %s", bi.val.String(), want)
	}
}

func mustInt(t *testing.T, v Value, want int64) {
	t.Helper()
	i, ok := v.(Int)
	if !ok {
		t.Fatalf("expected Int, got %T (%v)", v, v)
	}
	if int64(i) != want {
		t.Errorf("Int value = %d, want %d", int64(i), want)
	}
}

// bigIntFromInt converts an Int-width value to a *big.Int. It took a
// platform-width int before Int became int64, when the distinction mattered;
// the two are now the same width on every host.
func bigIntFromInt(n int64) *big.Int {
	return new(big.Int).SetInt64(int64(n))
}

func TestNumAddP_Boundary(t *testing.T) {
	// In-range: stays Int
	r, err := NumAddP(Int(1), Int(2))
	if err != nil {
		t.Fatal(err)
	}
	mustInt(t, r, 3)

	// MaxInt+1 promotes
	r, err = NumAddP(Int(math.MaxInt64), Int(1))
	if err != nil {
		t.Fatal(err)
	}
	expected := new(big.Int).Add(bigIntFromInt(math.MaxInt64), big.NewInt(1))
	mustBigIntString(t, r, expected.String())

	// MinInt-1 (via add of -1) promotes
	r, err = NumAddP(Int(math.MinInt64), Int(-1))
	if err != nil {
		t.Fatal(err)
	}
	expected = new(big.Int).Add(bigIntFromInt(math.MinInt64), big.NewInt(-1))
	mustBigIntString(t, r, expected.String())
}

func TestNumSubP_Boundary(t *testing.T) {
	// In-range: stays Int
	r, err := NumSubP(Int(5), Int(3))
	if err != nil {
		t.Fatal(err)
	}
	mustInt(t, r, 2)

	r, err = NumSubP(Int(math.MinInt64), Int(1))
	if err != nil {
		t.Fatal(err)
	}
	expected := new(big.Int).Sub(bigIntFromInt(math.MinInt64), big.NewInt(1))
	mustBigIntString(t, r, expected.String())

	r, err = NumSubP(Int(0), Int(math.MinInt64))
	if err != nil {
		t.Fatal(err)
	}
	expected = new(big.Int).Neg(bigIntFromInt(math.MinInt64))
	mustBigIntString(t, r, expected.String())
}

func TestNumMulP_Boundary(t *testing.T) {
	// In-range: stays Int
	r, err := NumMulP(Int(6), Int(7))
	if err != nil {
		t.Fatal(err)
	}
	mustInt(t, r, 42)

	// MaxInt*2 promotes
	r, err = NumMulP(Int(math.MaxInt64), Int(2))
	if err != nil {
		t.Fatal(err)
	}
	expected := new(big.Int).Mul(bigIntFromInt(math.MaxInt64), big.NewInt(2))
	mustBigIntString(t, r, expected.String())

	// MinInt*-1 promotes
	r, err = NumMulP(Int(math.MinInt64), Int(-1))
	if err != nil {
		t.Fatal(err)
	}
	expected = new(big.Int).Neg(bigIntFromInt(math.MinInt64))
	mustBigIntString(t, r, expected.String())

	// (MinInt/2)*3 promotes
	r, err = NumMulP(Int(math.MinInt64/2), Int(3))
	if err != nil {
		t.Fatal(err)
	}
	expected = new(big.Int).Mul(bigIntFromInt(math.MinInt64/2), big.NewInt(3))
	mustBigIntString(t, r, expected.String())
}

func TestNumNegP_MinInt(t *testing.T) {
	r, err := NumNegP(Int(math.MinInt64))
	if err != nil {
		t.Fatal(err)
	}
	expected := new(big.Int).Neg(bigIntFromInt(math.MinInt64))
	mustBigIntString(t, r, expected.String())

	// In-range: stays Int
	r, err = NumNegP(Int(5))
	if err != nil {
		t.Fatal(err)
	}
	mustInt(t, r, -5)
}

// Clojure preserves BigInt identity through promoting ops: `(-' 1N)` returns
// `-1N`, not `Int -1`. NumNeg (the non-promoting version) downgrades via
// MaybeDowngrade, so NumNegP must handle BigInt directly to avoid that.
func TestNumNegP_BigIntStaysBigInt(t *testing.T) {
	r, err := NumNegP(NewBigIntFromInt64(1))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "-1")

	r, err = NumNegP(NewBigIntFromInt64(-1))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "1")

	// Large BigInt that wouldn't fit in int64 should also stay BigInt
	huge, _ := NewBigIntFromString("99999999999999999999999")
	r, err = NumNegP(huge)
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "-99999999999999999999999")
}

// --- Delegation paths -------------------------------------------------------
//
// Non-Int/Int combinations must behave exactly like NumAdd/Sub/Mul/Neg.

func TestNumAddP_Delegation(t *testing.T) {
	// Int + Float → Float (delegates to NumAdd)
	r, err := NumAddP(Int(math.MaxInt64), Float(1.0))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(Float); !ok {
		t.Errorf("Int + Float should yield Float, got %T", r)
	}

	// Float + Int → Float
	r, err = NumAddP(Float(1.0), Int(math.MaxInt64))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(Float); !ok {
		t.Errorf("Float + Int should yield Float, got %T", r)
	}

	// Int + BigInt → BigInt (NumAdd already promotes)
	r, err = NumAddP(Int(5), NewBigIntFromInt64(10))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "15")

	// BigInt + Int → BigInt
	r, err = NumAddP(NewBigIntFromInt64(10), Int(5))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "15")

	// BigInt + BigInt → BigInt
	r, err = NumAddP(NewBigIntFromInt64(10), NewBigIntFromInt64(20))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "30")

	// Float + BigInt → Float (lossy but consistent with NumAdd)
	r, err = NumAddP(Float(1.0), NewBigIntFromInt64(2))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(Float); !ok {
		t.Errorf("Float + BigInt should yield Float, got %T", r)
	}
}

func TestNumMulP_BigIntBoundary(t *testing.T) {
	// BigInt * Int big enough to stay BigInt — should not silently downgrade
	big1 := NewBigInt(new(big.Int).Lsh(big.NewInt(1), 70)) // 2^70
	r, err := NumMulP(big1, Int(2))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "2361183241434822606848")
}

func TestNumAddP_InvalidInputs(t *testing.T) {
	// (+' 1 nil) → error (let-go surfaces error from NumAdd's bottom branch)
	_, err := NumAddP(Int(1), NIL)
	if err == nil {
		t.Error("NumAddP(Int, NIL) should return error")
	}
}

// NumSubP and NumMulP delegate to NumSub/NumMul for any non-Int/Int pair.
// Cover the same delegation paths NumAddP_Delegation does, to lock in the
// invariant that BigInt identity is preserved through binary ops.
func TestNumSubP_Delegation(t *testing.T) {
	r, err := NumSubP(NewBigIntFromInt64(10), Int(3))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "7")

	r, err = NumSubP(Float(5.0), Int(2))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(Float); !ok {
		t.Errorf("Float - Int should yield Float, got %T", r)
	}
}

func TestNumMulP_Delegation(t *testing.T) {
	r, err := NumMulP(NewBigIntFromInt64(10), Int(3))
	if err != nil {
		t.Fatal(err)
	}
	mustBigIntString(t, r, "30")

	r, err = NumMulP(Float(5.0), Int(2))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(Float); !ok {
		t.Errorf("Float * Int should yield Float, got %T", r)
	}
}

// NumNegP delegates for any non-Int / non-BigInt operand (Float, Ratio,
// BigDecimal). Float is the simplest representative.
func TestNumNegP_Delegation(t *testing.T) {
	r, err := NumNegP(Float(3.5))
	if err != nil {
		t.Fatal(err)
	}
	f, ok := r.(Float)
	if !ok {
		t.Fatalf("NumNegP(Float) should yield Float, got %T", r)
	}
	if float64(f) != -3.5 {
		t.Errorf("NumNegP(3.5) = %v, want -3.5", f)
	}
}
