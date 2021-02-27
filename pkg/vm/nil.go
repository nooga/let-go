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

type theNilType struct {
	zero *Nil
}

func (t *theNilType) Name() string                     { return "nil" }
func (t *theNilType) Box(_ interface{}) (Value, error) { return t.zero, nil }

// Nil is a Value whose only value is Nil
type Nil struct{}

// Type implements Value
func (n *Nil) Type() ValueType { return NilType }

// Unbox implements Value
func (n *Nil) Unbox() interface{} { return nil }

// NilType is the type of NilValues
var NilType *theNilType

// NIL is the only value of NilType (and only instance of NilValue)
var NIL *Nil

func init() {
	NilType = &theNilType{zero: nil}
	NIL = NilType.zero
}

func (n *Nil) String() string {
	return "nil"
}
