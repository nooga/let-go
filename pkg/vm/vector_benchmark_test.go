package vm

import (
	"strconv"
	"testing"
)

// Benchmark vector creation with different sizes
func BenchmarkVectorCreation(b *testing.B) {
	benchSizes := []int{10, 100, 1000, 10000}

	for _, size := range benchSizes {
		// Prepare test data
		values := make([]Value, size)
		for i := range values {
			values[i] = Int(i)
		}

		b.Run("ArrayVector/"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewArrayVector(values)
			}
		})

		b.Run("PersistentVector/"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewPersistentVector(values)
			}
		})
	}
}

// benchVecSink prevents the compiler from dead-code-eliminating the
// benchmarked ValueAt calls: storing the result to a package-level var is an
// observable side effect it must keep. Without a sink the whole access loop
// can be optimized away (giving a fictional sub-nanosecond time that then
// reads as a huge "regression" the moment devirtualization stops firing).
var benchVecSink Value

// Benchmark random access performance
func BenchmarkVectorAccess(b *testing.B) {
	benchSizes := []int{10, 100, 1000, 10000}

	for _, size := range benchSizes {
		// Prepare test data
		values := make([]Value, size)
		for i := range values {
			values[i] = Int(i)
		}

		arrayVec := NewArrayVector(values).(Lookup)
		persistentVec := NewPersistentVector(values).(Lookup)

		// Box the keys ONCE, outside the timed loop. Passing Int(n) at the call
		// site boxes it into the Value interface argument every iteration, and
		// for n > 255 that escapes and allocates — which measures Int boxing,
		// not vector access. Hoisting isolates the access cost.
		k0, kMid, kLast := Value(Int(0)), Value(Int(size/2)), Value(Int(size-1))

		b.Run("ArrayVector/"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchVecSink = arrayVec.ValueAt(k0)    // First
				benchVecSink = arrayVec.ValueAt(kMid)  // Middle
				benchVecSink = arrayVec.ValueAt(kLast) // Last
			}
		})

		b.Run("PersistentVector/"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchVecSink = persistentVec.ValueAt(k0)
				benchVecSink = persistentVec.ValueAt(kMid)
				benchVecSink = persistentVec.ValueAt(kLast)
			}
		})
	}
}

// Benchmark sequential appending with Conj
func BenchmarkVectorConj(b *testing.B) {
	benchSizes := []int{10, 100, 1000}

	// Conj never mutates its receiver, so one empty base serves every
	// iteration. Building it inside the loop under StopTimer/StartTimer
	// cost a stop-the-world ReadMemStats per iteration: tens of
	// microseconds of wall time around a few microseconds of measured work.
	// b.Loop keeps the setup above out of the timed region (the timer is
	// already running when this function is entered) and keeps the loop
	// body's results alive.
	for _, size := range benchSizes {
		b.Run("ArrayVector/"+strconv.Itoa(size), func(b *testing.B) {
			base := NewArrayVector([]Value{}).(Collection)
			for b.Loop() {
				vec := base
				// Add elements one by one
				for j := range size {
					vec = vec.Conj(Int(j))
				}
			}
		})

		b.Run("PersistentVector/"+strconv.Itoa(size), func(b *testing.B) {
			base := NewPersistentVector([]Value{}).(Collection)
			for b.Loop() {
				vec := base
				// Add elements one by one
				for j := range size {
					vec = vec.Conj(Int(j))
				}
			}
		})
	}
}

// Benchmark updating elements with Assoc
func BenchmarkVectorAssoc(b *testing.B) {
	benchSizes := []int{10, 100, 1000, 10000}

	for _, size := range benchSizes {
		// Prepare test data
		values := make([]Value, size)
		for i := range values {
			values[i] = Int(i)
		}

		// Assoc copies on write for both vector kinds, so the base built
		// once here is what every iteration used to rebuild under StopTimer.
		b.Run("ArrayVector/"+strconv.Itoa(size), func(b *testing.B) {
			vec := NewArrayVector(values).(Associative)
			i := 0
			for b.Loop() {
				i++
				// Update elements at different positions
				vec.Assoc(Int(0), Int(i))      // First
				vec.Assoc(Int(size/2), Int(i)) // Middle
				vec.Assoc(Int(size-1), Int(i)) // Last
			}
		})

		b.Run("PersistentVector/"+strconv.Itoa(size), func(b *testing.B) {
			vec := NewPersistentVector(values).(Associative)
			i := 0
			for b.Loop() {
				i++
				// Update elements at different positions
				vec.Assoc(Int(0), Int(i))      // First
				vec.Assoc(Int(size/2), Int(i)) // Middle
				vec.Assoc(Int(size-1), Int(i)) // Last
			}
		})
	}
}

// Benchmark sequential iteration using Seq
func BenchmarkVectorSeq(b *testing.B) {
	benchSizes := []int{10, 100, 1000}

	for _, size := range benchSizes {
		// Prepare test data
		values := make([]Value, size)
		for i := range values {
			values[i] = Int(i)
		}

		arrayVec := NewArrayVector(values).(Sequable)
		persistentVec := NewPersistentVector(values).(Sequable)

		// Seqs are immutable (Next returns a fresh node), so one head built
		// here is walked from the start on every iteration, and the timed
		// region still covers exactly the walk, as it did under StopTimer.
		b.Run("ArrayVector/"+strconv.Itoa(size), func(b *testing.B) {
			head := arrayVec.Seq()
			for b.Loop() {
				// Iterate through the entire sequence
				for seq := head; seq != nil; seq = seq.Next() {
					seq.First()
				}
			}
		})

		b.Run("PersistentVector/"+strconv.Itoa(size), func(b *testing.B) {
			head := persistentVec.Seq()
			for b.Loop() {
				// Iterate through the entire sequence
				for seq := head; seq != nil; seq = seq.Next() {
					seq.First()
				}
			}
		})
	}
}

// Benchmark large vector append operations (append 10000 items)
func BenchmarkLargeVectorAppend(b *testing.B) {
	const size = 10000

	b.Run("ArrayVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vec := NewArrayVector([]Value{}).(Collection)
			for j := range size {
				vec = vec.Conj(Int(j))
			}
		}
	})

	b.Run("PersistentVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vec := NewPersistentVector([]Value{}).(Collection)
			for j := range size {
				vec = vec.Conj(Int(j))
			}
		}
	})
}

// Benchmark creating multiple versions through successive updates
// This should highlight the advantage of persistent data structures
// when many versions need to be maintained
func BenchmarkVectorVersioning(b *testing.B) {
	const initialSize = 1000
	const numUpdates = 100

	// Create initial data
	values := make([]Value, initialSize)
	for i := range values {
		values[i] = Int(i)
	}

	b.Run("ArrayVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Create initial vector
			versions := make([]Lookup, numUpdates+1)
			original := NewArrayVector(values).(Associative)
			versions[0] = original.(Lookup)

			// Create multiple versions through updates
			for j := range numUpdates {
				// Create a copy
				originalValues := original.Unbox().([]Value)
				copyValues := make([]Value, len(originalValues))
				copy(copyValues, originalValues)

				nextVersion := NewArrayVector(copyValues).(Associative)
				// Update at different positions
				nextVersion = nextVersion.Assoc(Int(j%initialSize), Int(-j))
				versions[j+1] = nextVersion.(Lookup)
			}

			// Use the versions to prevent optimization
			for _, v := range versions {
				_ = v.ValueAt(Int(0))
			}
		}
	})

	b.Run("PersistentVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Create initial vector
			versions := make([]Lookup, numUpdates+1)
			original := NewPersistentVector(values).(Associative)
			versions[0] = original.(Lookup)

			// Create multiple versions through updates
			current := original
			for j := range numUpdates {
				// The beauty of persistent vectors - no need to copy
				nextVersion := current.Assoc(Int(j%initialSize), Int(-j))
				versions[j+1] = nextVersion.(Lookup)
				current = nextVersion
			}

			// Use the versions to prevent optimization
			for _, v := range versions {
				_ = v.ValueAt(Int(0))
			}
		}
	})
}
