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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContext_Compile(t *testing.T) {
	tests := map[string]interface{}{
		"(+ (* 2 20) 2)":          42,
		`(if true "big" "meh")`:   "big",
		`(if false "big" "meh")`:  "meh",
		`(if nil 1 2)`:            2,
		`(if true 101)`:           101,
		`(if false 101)`:          nil,
		`(do 1 2 3)`:              3,
		`(do (+ 1 2))`:            3,
		`(do)`:                    nil,
		`(do (def x 40) (+ x 2))`: 42,
	}
	for k, v := range tests {
		out, err := Eval(k)
		assert.NoError(t, err)
		assert.Equal(t, v, out.Unbox())
	}
}
