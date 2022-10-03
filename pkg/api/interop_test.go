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
