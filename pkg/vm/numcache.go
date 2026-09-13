package vm

// Cached integer values to avoid interface boxing allocations.
// Covers the common range for loop counters, small arithmetic, etc.
const (
	intCacheMin = -128
	intCacheMax = 255
)

var intCache [intCacheMax - intCacheMin + 1]Value

func init() {
	for i := intCacheMin; i <= intCacheMax; i++ {
		intCache[i-intCacheMin] = Int(i)
	}
}

// MakeInt returns a cached Value for small ints, avoiding heap allocation.
func MakeInt(v int) Value {
	if v >= intCacheMin && v <= intCacheMax {
		return intCache[v-intCacheMin]
	}
	return Int(v)
}

// MakeFloat returns a Float Value.
func MakeFloat(v float64) Value {
	return Float(v)
}

// MakeInt64 is MakeInt for arithmetic results: takes the full 64-bit value so
// a 32-bit-int host (TinyGo wasm) does not truncate it on the way in.
func MakeInt64(v int64) Value {
	if v >= int64(intCacheMin) && v <= int64(intCacheMax) {
		return intCache[v-int64(intCacheMin)]
	}
	return Int(v)
}
