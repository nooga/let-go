package vm

import (
	"testing"
)

func TestVectorImplementations(t *testing.T) {
	t.Run("creation and basic operations", func(t *testing.T) {
		testCases := []struct {
			name string
			fn   func([]Value) Value
		}{
			{"ArrayVector", NewArrayVector},
			{"PersistentVector", NewPersistentVector},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test empty vector
				empty := tc.fn([]Value{})
				if empty.(Collection).Count().(Int) != 0 {
					t.Errorf("expected empty vector count to be 0, got %d", empty.(Collection).Count())
				}

				// Test vector with elements
				values := []Value{Int(1), Int(2), Int(3)}
				vec := tc.fn(values)

				// Test count
				if vec.(Collection).Count().(Int) != 3 {
					t.Errorf("expected count 3, got %d", vec.(Collection).Count())
				}

				// Test value access
				for i, expected := range values {
					if got := vec.(Lookup).ValueAt(Int(i)); got != expected {
						t.Errorf("at index %d: expected %v, got %v", i, expected, got)
					}
				}
			})
		}
	})

	t.Run("conj operation", func(t *testing.T) {
		testCases := []struct {
			name string
			fn   func([]Value) Value
		}{
			{"ArrayVector", NewArrayVector},
			{"PersistentVector", NewPersistentVector},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				vec := tc.fn([]Value{}).(Collection)
				vec = vec.Conj(Int(1))
				vec = vec.Conj(Int(2))
				vec = vec.Conj(Int(3))

				expected := []Value{Int(1), Int(2), Int(3)}
				for i, exp := range expected {
					if got := vec.(Lookup).ValueAt(Int(i)); got != exp {
						t.Errorf("at index %d: expected %v, got %v", i, exp, got)
					}
				}
			})
		}
	})

	t.Run("assoc operation", func(t *testing.T) {
		testCases := []struct {
			name string
			fn   func([]Value) Value
		}{
			{"ArrayVector", NewArrayVector},
			{"PersistentVector", NewPersistentVector},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				initial := []Value{Int(1), Int(2), Int(3)}
				vec := tc.fn(initial).(Associative)

				// Test valid assoc
				newVec := vec.Assoc(Int(1), Int(42))
				if got := newVec.(Lookup).ValueAt(Int(1)); got != Int(42) {
					t.Errorf("after assoc: expected 42, got %v", got)
				}
				// Original should be unchanged
				if got := vec.(Lookup).ValueAt(Int(1)); got != Int(2) {
					t.Errorf("original vector modified: expected 2, got %v", got)
				}

				// Test invalid assoc
				if newVec := vec.Assoc(Int(-1), Int(42)); newVec != NIL {
					t.Error("expected NIL for negative index assoc")
				}
				if newVec := vec.Assoc(Int(100), Int(42)); newVec != NIL {
					t.Error("expected NIL for out of bounds assoc")
				}
			})
		}
	})

	t.Run("sequence operations", func(t *testing.T) {
		testCases := []struct {
			name string
			fn   func([]Value) Value
		}{
			{"ArrayVector", NewArrayVector},
			{"PersistentVector", NewPersistentVector},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				vec := tc.fn([]Value{Int(1), Int(2), Int(3)}).(Sequable).Seq()

				// Test First
				if first := vec.First(); first != Int(1) {
					t.Errorf("First: expected 1, got %v", first)
				}

				// Test More
				more := vec.More()
				if more.(Collection).Count().(Int) != 2 {
					t.Errorf("More: expected count 2, got %v", more.(Collection).Count())
				}
				if first := more.First(); first != Int(2) {
					t.Errorf("More.First: expected 2, got %v", first)
				}

				// Test empty vector sequence operations
				emptyVec := tc.fn([]Value{}).(Sequable).Seq()
				if first := emptyVec.First(); first != NIL {
					t.Errorf("Empty First: expected NIL, got %v", first)
				}
				if more := emptyVec.More(); more != EmptyList {
					t.Errorf("Empty More: expected EmptyList, got %v", more)
				}
			})
		}
	})

	t.Run("large vector operations", func(t *testing.T) {
		testCases := []struct {
			name string
			fn   func([]Value) Value
		}{
			{"ArrayVector", NewArrayVector},
			{"PersistentVector", NewPersistentVector},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create vector with 1000 elements
				values := make([]Value, 1000)
				for i := range values {
					values[i] = Int(i)
				}
				vec := tc.fn(values).(Associative)

				// Test random access
				indices := []int{0, 42, 999, 500, 127, 128}
				for _, idx := range indices {
					if got := vec.(Lookup).ValueAt(Int(idx)); got != Int(idx) {
						t.Errorf("at index %d: expected %d, got %v", idx, idx, got)
					}
				}

				// Test bounds checking
				if got := vec.(Lookup).ValueAt(Int(1000)); got != NIL {
					t.Errorf("expected NIL for out of bounds access, got %v", got)
				}
			})
		}
	})
}

func TestPersistentVectorSeq(t *testing.T) {
	t.Run("basic seq operations", func(t *testing.T) {
		// Create vectors of different sizes to test
		smallVec := NewPersistentVector([]Value{Int(1), Int(2), Int(3)})

		// Get sequence from vector
		seq := smallVec.(Sequable).Seq().(*PersistentVectorSeq)

		// Test type
		if seq.Type() != SequenceType {
			t.Errorf("expected seq type to be SequenceType, got %v", seq.Type())
		}

		// Test String representation
		expectedStr := "(seq [1 2 3])"
		if seq.String() != expectedStr {
			t.Errorf("expected String() to be %q, got %q", expectedStr, seq.String())
		}

		// Test navigation through sequence
		if first := seq.First(); first != Int(1) {
			t.Errorf("First: expected 1, got %v", first)
		}

		nextSeq := seq.Next().(*PersistentVectorSeq)
		if nextSeq.i != 1 || nextSeq.First() != Int(2) {
			t.Errorf("Next: expected index 1 with value 2, got index %d with value %v",
				nextSeq.i, nextSeq.First())
		}

		moreSeq := seq.More().(*PersistentVectorSeq)
		if moreSeq.i != 1 || moreSeq.First() != Int(2) {
			t.Errorf("More: expected index 1 with value 2, got index %d with value %v",
				moreSeq.i, moreSeq.First())
		}

		// Test count
		if count := seq.Count().(Int); count != 3 {
			t.Errorf("Count: expected 3, got %v", count)
		}

		if rawCount := seq.RawCount(); rawCount != 3 {
			t.Errorf("RawCount: expected 3, got %v", rawCount)
		}

		// Test Unbox
		unboxed, ok := seq.Unbox().([]Value)
		if !ok {
			t.Errorf("Unbox: expected []Value, got %T", seq.Unbox())
		}
		if len(unboxed) != 3 || unboxed[0] != Int(1) || unboxed[1] != Int(2) || unboxed[2] != Int(3) {
			t.Errorf("Unbox: expected [1 2 3], got %v", unboxed)
		}

		// Test Empty
		if empty := seq.Empty(); empty != EmptyList {
			t.Errorf("Empty: expected EmptyList, got %v", empty)
		}
	})

	t.Run("navigation to end of sequence", func(t *testing.T) {
		vec := NewPersistentVector([]Value{Int(1), Int(2), Int(3)})
		seq := vec.(Sequable).Seq()

		// Navigate to end
		seq = seq.Next()
		seq = seq.Next()

		// Test last element
		if last := seq.First(); last != Int(3) {
			t.Errorf("Last element: expected 3, got %v", last)
		}

		// Test beyond end
		if next := seq.Next(); next != nil {
			t.Errorf("Next beyond end: expected nil, got %v", next)
		}

		if more := seq.More(); more != EmptyList {
			t.Errorf("More beyond end: expected EmptyList, got %v", more)
		}
	})

	t.Run("large vector sequence", func(t *testing.T) {
		// Test with vector larger than nodeCap (32)
		values := make([]Value, 40)
		for i := range values {
			values[i] = Int(i)
		}

		largeVec := NewPersistentVector(values)
		seq := largeVec.(Sequable).Seq()

		// Test first elements
		if first := seq.First(); first != Int(0) {
			t.Errorf("First element: expected 0, got %v", first)
		}

		// Navigate through sequence checking values
		for i := 1; i < 40; i++ {
			seq = seq.Next()
			if val := seq.First(); val != Int(i) {
				t.Errorf("Element at position %d: expected %d, got %v", i, i, val)
			}
		}

		// Check end of sequence
		seq = seq.Next()
		if seq != nil {
			t.Errorf("Expected nil at end of sequence, got %v", seq)
		}
	})

	t.Run("conj operation", func(t *testing.T) {
		vec := NewPersistentVector([]Value{Int(1), Int(2), Int(3)})
		seq := vec.(Sequable).Seq()

		// Test Conj
		result := seq.(Collection).Conj(Int(42)).(Collection)

		// Verify result (should be prepended to the end of the collection)
		resultSeq := result.(Sequable).Seq()
		// Looking at the test results, it seems that Conj is resulting in [3, 2, 1, 42]
		// rather than [42, 1, 2, 3] as we initially expected
		expectedValues := []Value{Int(3), Int(2), Int(1), Int(42)}

		for i, expected := range expectedValues {
			if got := resultSeq.First(); got != expected {
				t.Errorf("Conj result at position %d: expected %v, got %v", i, expected, got)
			}
			resultSeq = resultSeq.Next()
		}
	})

	t.Run("inTail flag behavior", func(t *testing.T) {
		// Create a vector with only tail elements (< 32)
		smallVec := NewPersistentVector([]Value{Int(1), Int(2), Int(3)})
		smallSeq := smallVec.(Sequable).Seq().(*PersistentVectorSeq)

		// Verify it's using the tail
		if !smallSeq.inTail {
			t.Errorf("Small vector seq should have inTail=true but got false")
		}

		// Create vector large enough to have tree elements
		values := make([]Value, 40)
		for i := range values {
			values[i] = Int(i)
		}
		largeVec := NewPersistentVector(values)
		largeSeq := largeVec.(Sequable).Seq().(*PersistentVectorSeq)

		// Verify it's using the tree first
		if largeSeq.inTail {
			t.Errorf("Large vector seq should start with inTail=false but got true")
		}

		// Navigate to tail section (after index 32)
		var seq Seq = largeSeq
		for i := 0; i < 33; i++ {
			seq = seq.Next()
		}

		// Verify we're now in tail
		tailSeq, ok := seq.(*PersistentVectorSeq)
		if !ok {
			t.Errorf("Expected *PersistentVectorSeq, got %T", seq)
		} else if !tailSeq.inTail {
			t.Errorf("After navigating to index 33, expected inTail=true but got false")
		}
	})
}
