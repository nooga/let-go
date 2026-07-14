/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

// TypeUnion is a reusable class/interface descriptor backed by several
// concrete let-go value types. It is a Value but deliberately not a ValueType:
// no runtime value has the union itself as its concrete type.
type TypeUnion struct {
	name    string
	types   []ValueType
	members map[ValueType]struct{}
}

func NewTypeUnion(name string, types ...ValueType) *TypeUnion {
	members := make(map[ValueType]struct{}, len(types))
	unique := make([]ValueType, 0, len(types))
	for _, valueType := range types {
		if _, exists := members[valueType]; exists {
			continue
		}
		members[valueType] = struct{}{}
		unique = append(unique, valueType)
	}
	return &TypeUnion{name: name, types: unique, members: members}
}

func (u *TypeUnion) Type() ValueType { return TypeType }
func (u *TypeUnion) Unbox() any      { return u.Types() }
func (u *TypeUnion) String() string  { return u.name }

func (u *TypeUnion) Contains(valueType ValueType) bool {
	_, ok := u.members[valueType]
	return ok
}

func (u *TypeUnion) Types() []ValueType {
	return append([]ValueType(nil), u.types...)
}
