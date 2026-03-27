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
	tests := []struct {
		code   string
		result interface{}
	}{
		{"(* 11 11)", 121},
		{"(map inc [1 2 3])", []vm.Value{vm.Int(2), vm.Int(3), vm.Int(4)}},
	}
	for _, z := range tests {
		v, err := c.Run(z.code)
		assert.NoError(t, err, z.code)
		assert.Equal(t, z.result, v.Unbox(), z.code)
	}
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
