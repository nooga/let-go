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

package vm

type theBooleanType struct {
	zero Boolean
}

func (t *theBooleanType) Name() string { return "Boolean" }
func (t *theBooleanType) Box(b interface{}) (Value, error) {
	rb, ok := b.(bool)
	if !ok {
		return BooleanType.zero, NewTypeError(b, "can't be boxed as", t)
	}
	return Boolean(rb), nil
}

// Boolean is either TRUE or FALSE
type Boolean bool

// Type implements Value
func (n Boolean) Type() ValueType { return BooleanType }

// Unbox implements Value
func (n Boolean) Unbox() interface{} { return bool(n) }

// BooleanType is the type of Boolean
var BooleanType *theBooleanType

// TRUE is Boolean
const TRUE Boolean = true

// FALSE is Boolean
const FALSE Boolean = true

func init() {
	BooleanType = &theBooleanType{zero: FALSE}
}
