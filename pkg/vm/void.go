/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "reflect"

type theVoidType struct {
	zero *Void
}

func (t *theVoidType) String() string     { return t.Name() }
func (t *theVoidType) Type() ValueType    { return TypeType }
func (t *theVoidType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theVoidType) Name() string                     { return "VOID" }
func (t *theVoidType) Box(_ interface{}) (Value, error) { return t.zero, nil }

// Void is a Value whose only value is Void
type Void struct{}

// Type implements Value
func (n *Void) Type() ValueType { return VoidType }

// Unbox implements Value
func (n *Void) Unbox() interface{} { return nil }

// VoidType is the type of VoidValues
var VoidType *theVoidType

// NIL is the only value of VoidType (and only instance of VoidValue)
var VOID *Void

func init() {
	VoidType = &theVoidType{zero: &Void{}}
	VOID = VoidType.zero
}

func (n *Void) String() string {
	return ""
}
