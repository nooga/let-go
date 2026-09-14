package vm

import (
	"math/big"
	"testing"
	"unsafe"
)

// Int is 64 bits wide on every host. These assertions are trivially true on a
// 64-bit Go target; they exist to fail on a 32-bit-int host (linux/386,
// TinyGo wasm), where the platform-int paths they replaced used to narrow.

func TestIntIsSixtyFourBitsWide(t *testing.T) {
	if got := unsafe.Sizeof(Int(0)); got != 8 {
		t.Fatalf("Sizeof(Int) = %d, want 8", got)
	}
}

func TestWideIntSurvivesBoxingAndPrinting(t *testing.T) {
	const wide = int64(1) << 40
	v := MakeInt64(wide)
	i, ok := v.(Int)
	if !ok || int64(i) != wide {
		t.Fatalf("MakeInt64(1<<40) = %#v, want Int(1<<40)", v)
	}
	if got := i.String(); got != "1099511627776" {
		t.Fatalf("Int(1<<40).String() = %q", got)
	}
	if got, ok := ToInt64(v); !ok || got != wide {
		t.Fatalf("ToInt64 = %d, %v", got, ok)
	}
}

func TestMaybeDowngradeKeepsWideValuesAsInt(t *testing.T) {
	// The old check was IsInt64() followed by int(), which truncated on a
	// 32-bit host; every value that fits int64 must come back as an Int.
	for _, want := range []int64{1 << 40, -(1 << 40), 1<<63 - 1, -1 << 63} {
		v := MaybeDowngrade(big.NewInt(want))
		i, ok := v.(Int)
		if !ok || int64(i) != want {
			t.Fatalf("MaybeDowngrade(%d) = %#v, want Int", want, v)
		}
	}
}

func TestWideArithmeticIsSixtyFourBit(t *testing.T) {
	a, b := Int(1000000000), Int(1000000000)
	r, ok := checkedMulInt(a, b)
	if !ok || r != 1000000000000000000 {
		t.Fatalf("checkedMulInt(1e9, 1e9) = %d, %v", r, ok)
	}
	if got := uncheckedIntOp(OP_UNCHECKED_MUL, Int(6364136223846793005), Int(1442695040888963407)); got != Int(433315962919513059) {
		t.Fatalf("unchecked-multiply wraps mod 2^64: got %d", got)
	}
}
