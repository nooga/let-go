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

import (
	"unicode/utf8"
)

type theCharType struct {
	zero Char
}

func (lt *theCharType) Name() string { return "Char" }

func (lt *theCharType) Box(bare interface{}) (Value, error) {
	raw, ok := bare.(rune)
	if !ok {
		return CharType.zero, NewTypeError(bare, "can't be boxed as", lt)
	}
	return Char(raw), nil
}

// CharType is the type of CharValues
var CharType *theCharType

func init() {
	CharType = &theCharType{zero: utf8.RuneError}
}

// Char is boxed rune
type Char rune

// Type implements Value
func (l Char) Type() ValueType { return CharType }

// Unbox implements Unbox
func (l Char) Unbox() interface{} {
	return rune(l)
}

func (l Char) String() string {
	return "\\" + string(l)
}
