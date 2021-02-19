/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
 * persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/lego/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestReaderBasic(t *testing.T) {
	cases := map[string]vm.Value{
		"1":                    vm.Int(1),
		"+1":                   vm.Int(1),
		"-1":                   vm.Int(-1),
		"987654321":            vm.Int(987654321),
		"+987654321":           vm.Int(987654321),
		"-987654321":           vm.Int(-987654321),
		"true":                 vm.TRUE,
		"false":                vm.FALSE,
		"nil":                  vm.NIL,
		"foo":                  vm.Symbol("foo"),
		"()":                   vm.EmptyList,
		"(    )":               vm.EmptyList,
		"\"hello\"":            vm.String("hello"),
		"\"h\\\"el\\tl\\\\o\"": vm.String("h\"el\tl\\o"),
	}

	for p, e := range cases {
		r := NewLispReader(strings.NewReader(p), "<reader>")
		o, err := r.Read()
		assert.NoError(t, err)
		assert.Equal(t, e, o)
	}
}

func TestSimpleCall(t *testing.T) {
	p := "(+ 40 2)"
	r := NewLispReader(strings.NewReader(p), "<reader>")
	o, err := r.Read()
	assert.NoError(t, err)

	out, err := vm.ListType.Box([]vm.Value{
		vm.Symbol("+"),
		vm.Int(40),
		vm.Int(2),
	})

	assert.NoError(t, err)
	assert.Equal(t, out, o)
}
