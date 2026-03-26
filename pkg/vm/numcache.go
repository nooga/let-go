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
