package api_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	c, err := api.NewLetGo("interop")
	assert.NoError(t, err)

	// Scalar result: Unbox is the canonical conversion.
	v, err := c.Run("(* 11 11)")
	assert.NoError(t, err)
	assert.Equal(t, 121, v.Unbox())

	// Sequence result: `map` returns a lazy seq, so callers should either
	// iterate with vm.Iter or realize explicitly with vm.SeqToSlice.
	v, err = c.Run("(map inc [1 2 3])")
	assert.NoError(t, err)
	realized, err := vm.SeqToSlice(v)
	assert.NoError(t, err)
	assert.Equal(t, []vm.Value{vm.Int(2), vm.Int(3), vm.Int(4)}, realized)

	// And vm.Iter is the idiomatic range-loop shape.
	var collected []vm.Value
	for elem := range vm.Iter(v) {
		collected = append(collected, elem)
	}
	assert.Equal(t, []vm.Value{vm.Int(2), vm.Int(3), vm.Int(4)}, collected)
}

func TestDef(t *testing.T) {
	c, err := api.NewLetGo("interop")
	assert.NoError(t, err)
	err = c.Def("x", 2)
	assert.NoError(t, err)
	err = c.Def("f", func(a, b int) int {
		return a * b
	})
	assert.NoError(t, err)
	tests := []struct {
		code   string
		result interface{}
	}{
		{"(def y (+ x 11))", 13},
		{"(f x y)", 26},
	}
	for _, z := range tests {
		v, err := c.Run(z.code)
		assert.NoError(t, err, z.code)
		assert.Equal(t, z.result, v.Unbox(), z.code)
	}
}

func TestChannels(t *testing.T) {
	c, err := api.NewLetGo("interop")
	assert.NoError(t, err)

	inch := make(chan int)
	outch := make(vm.Chan)

	err = c.Def("in", inch)
	assert.NoError(t, err)
	err = c.Def("out", outch)
	assert.NoError(t, err)

	_, err = c.Run(`
		(go (loop [i (<! in)]
				(when i
					(>! out (inc i))
					(recur (<! in)))))
	`)
	assert.NoError(t, err)

	go func() {
		for i := 0; i < 10; i++ {
			inch <- i
		}
		close(inch)
	}()

	for i := 0; i < 10; i++ {
		j := <-outch
		assert.Equal(t, i+1, j.Unbox())
	}
}

type Item struct {
	Name  string
	Price float64
	Qty   int
}

func TestStructRecordRoundtrip(t *testing.T) {
	vm.RegisterStruct[Item]("interop/Item")

	c, err := api.NewLetGo("structInterop")
	assert.NoError(t, err)

	// Pass a Go struct to let-go — it becomes a Record
	err = c.Def("item", Item{Name: "Widget", Price: 9.99, Qty: 5})
	assert.NoError(t, err)

	// Access fields with keywords in let-go
	v, err := c.Run(`(:name item)`)
	assert.NoError(t, err)
	assert.Equal(t, "Widget", string(v.(vm.String)))

	v, err = c.Run(`(:price item)`)
	assert.NoError(t, err)
	assert.Equal(t, 9.99, float64(v.(vm.Float)))

	// Unmutated record roundtrips back via Unbox (fast path)
	v, err = c.Run(`item`)
	assert.NoError(t, err)
	got := v.Unbox().(Item)
	assert.Equal(t, Item{Name: "Widget", Price: 9.99, Qty: 5}, got)

	// Mutated record goes through slow path
	v, err = c.Run(`(assoc item :qty 10)`)
	assert.NoError(t, err)
	mutated, err := vm.ToStruct[Item](v.(*vm.Record))
	assert.NoError(t, err)
	assert.Equal(t, Item{Name: "Widget", Price: 9.99, Qty: 10}, mutated)

	// Record works with all map operations
	v, err = c.Run(`(count item)`)
	assert.NoError(t, err)
	assert.Equal(t, 3, v.Unbox())

	v, err = c.Run(`(contains? item :name)`)
	assert.NoError(t, err)
	assert.Equal(t, true, v.Unbox())
}

func TestStructPassedToLetGoFunction(t *testing.T) {
	vm.RegisterStruct[Item]("interop/Item")

	c, err := api.NewLetGo("structFn")
	assert.NoError(t, err)

	// Define a let-go function that processes structs
	_, err = c.Run(`(defn total [item] (* (:price item) (:qty item)))`)
	assert.NoError(t, err)

	// Call it with a Go struct
	err = c.Def("my-item", Item{Name: "Gadget", Price: 4.50, Qty: 3})
	assert.NoError(t, err)

	v, err := c.Run(`(total my-item)`)
	assert.NoError(t, err)
	assert.Equal(t, 13.5, float64(v.(vm.Float)))
}

// TestLazySeqToGoSliceFn: passing a let-go lazy sequence (from `map`,
// `filter`, etc.) to a Go function taking []T must auto-realize and
// per-element convert. Exercises the reflect bridge in boxArgForReflect.
func TestLazySeqToGoSliceFn(t *testing.T) {
	c, err := api.NewLetGo("lazyToSlice")
	assert.NoError(t, err)

	assert.NoError(t, c.Def("sum-ints", func(xs []int) int {
		s := 0
		for _, x := range xs {
			s += x
		}
		return s
	}))

	// (take 100 (iterate inc 1)) is an infinite-seq slice — must be
	// realized to []int by the bridge, not passed as a Seq value.
	v, err := c.Run(`(sum-ints (take 100 (iterate inc 1)))`)
	assert.NoError(t, err)
	assert.Equal(t, 5050, v.Unbox())

	// `(map inc [1 2 3])` returns a lazy seq — same path.
	v, err = c.Run(`(sum-ints (map inc [1 2 3]))`)
	assert.NoError(t, err)
	assert.Equal(t, 2+3+4, v.Unbox())
}

// TestEmptySeqToGoSliceFn: empty/nil sequences must produce empty Go
// slices, not panic or fall through to a one-element nil-filled slice.
// Regression test for the EmptyList sentinel handling in unboxSliceInto.
func TestEmptySeqToGoSliceFn(t *testing.T) {
	c, err := api.NewLetGo("emptySeq")
	assert.NoError(t, err)

	assert.NoError(t, c.Def("count-strs", func(xs []string) int {
		return len(xs)
	}))

	// (map identity []) → EmptyList. Must round-trip as empty []string.
	v, err := c.Run(`(count-strs (map identity []))`)
	assert.NoError(t, err)
	assert.Equal(t, 0, v.Unbox())

	// (filter pred …) producing no matches.
	v, err = c.Run(`(count-strs (filter (constantly false) ["a" "b"]))`)
	assert.NoError(t, err)
	assert.Equal(t, 0, v.Unbox())
}

// TestStructFieldWithLazySeq: assoc'ing a lazy seq into a struct field
// and round-tripping through ToStruct must work. Exercises the existing
// unboxSliceInto path via struct_mapping for the slice-field case.
type ItemList struct {
	Tags []string
}

func TestStructFieldWithLazySeq(t *testing.T) {
	vm.RegisterStruct[ItemList]("interop/ItemList")

	c, err := api.NewLetGo("structSeq")
	assert.NoError(t, err)

	assert.NoError(t, c.Def("base", ItemList{Tags: []string{"keep"}}))

	// `(map upper …)` is a lazy seq — assoc'd into a struct field
	// must realize to []string on ToStruct round-trip.
	v, err := c.Run(`(assoc base :tags (map (fn [s] (str s "!")) ["a" "b" "c"]))`)
	assert.NoError(t, err)

	got, err := vm.ToStruct[ItemList](v.(*vm.Record))
	assert.NoError(t, err)
	assert.Equal(t, []string{"a!", "b!", "c!"}, got.Tags)
}

// TestVmIterOverResult: vm.Iter is the documented way for Go-side code
// to consume a let-go sequence lazily. Confirms it works for both
// concrete (*List) and lazy (*LazySeq) result shapes.
func TestVmIterOverResult(t *testing.T) {
	c, err := api.NewLetGo("vmIter")
	assert.NoError(t, err)

	// Concrete list
	v, err := c.Run(`'(1 2 3)`)
	assert.NoError(t, err)
	var got []int64
	for x := range vm.Iter(v) {
		got = append(got, int64(x.(vm.Int)))
	}
	assert.Equal(t, []int64{1, 2, 3}, got)

	// Lazy seq, with early break — must stop without realizing further.
	v, err = c.Run(`(iterate inc 0)`)
	assert.NoError(t, err)
	got = got[:0]
	count := 0
	for x := range vm.Iter(v) {
		got = append(got, int64(x.(vm.Int)))
		count++
		if count >= 5 {
			break
		}
	}
	assert.Equal(t, []int64{0, 1, 2, 3, 4}, got)
}

func BenchmarkUse(b *testing.B) {
	c, err := api.NewLetGo("useBenchmark")
	if err != nil {
		b.Fatal(err)
	}
	c.SetLoadPath([]string{"../../", "."})
	for n := 0; n < b.N; n++ {
		_, err = c.Run("(use 'tns)")
	}
	if err != nil {
		b.Fatal(err)
	}
}
