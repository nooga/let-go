/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "reflect"

type theDefMetaPairsType struct{}

func (t *theDefMetaPairsType) String() string  { return t.Name() }
func (t *theDefMetaPairsType) Type() ValueType { return TypeType }
func (t *theDefMetaPairsType) Unbox() any      { return reflect.TypeFor[*theDefMetaPairsType]() }
func (t *theDefMetaPairsType) Name() string    { return "let-go.internal.DefMetaPairs" }
func (t *theDefMetaPairsType) Box(bare any) (Value, error) {
	pairs, ok := bare.([]Value)
	if !ok {
		return NIL, NewTypeError(bare, "can't be boxed as", t)
	}
	return DefMetaPairs(pairs), nil
}

var defMetaPairsType = &theDefMetaPairsType{}

// DefMetaPairs is the compact on-wire and deferred in-memory representation of
// metadata attached to a def. Entries alternate key, value. It is internal to
// bundle replay: ApplyVarMeta stores the pairs on the Var, and Var.Meta lazily
// materializes the public PersistentMap only when metadata is observed.
type DefMetaPairs []Value

func (p DefMetaPairs) String() string  { return ArrayVector(p).String() }
func (p DefMetaPairs) Type() ValueType { return defMetaPairsType }
func (p DefMetaPairs) Unbox() any      { return []Value(p) }
func (p DefMetaPairs) Hash() uint32    { return ArrayVector(p).Hash() }

func (p DefMetaPairs) Equals(other Value) bool {
	q, ok := other.(DefMetaPairs)
	return ok && ArrayVector(p).Equals(ArrayVector(q))
}
