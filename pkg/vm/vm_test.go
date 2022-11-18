/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilType(t *testing.T) {
	myNil, err := NilType.Box(nil)
	assert.NoError(t, err)
	assert.Equal(t, NIL, myNil)
	assert.Equal(t, "nil", myNil.Type().Name())
	assert.Nil(t, myNil.Unbox())
}

func TestIntType(t *testing.T) {
	for i := 0; i < 100; i++ {
		v := rand.Int()
		bv, err := IntType.Box(v)
		assert.NoError(t, err)
		assert.Equal(t, bv.Type(), IntType)
		assert.True(t, bv.Type() == IntType)
		assert.Equal(t, "let-go.lang.Int", bv.Type().Name())
		assert.Equal(t, v, bv.Unbox())
		assert.Equal(t, v, bv.Unbox().(int))
	}

	a, err := IntType.Box(9)
	assert.NoError(t, err)

	b, err := IntType.Box(1)
	assert.NoError(t, err)

	c, err := IntType.Box(1)
	assert.NoError(t, err)

	assert.NotEqual(t, a, b)
	assert.NotEqual(t, a, c)
	assert.Equal(t, c, b)
	assert.False(t, a == b)
	assert.True(t, c == b)

	bad := "bad"
	badInt, err := IntType.Box(bad)
	assert.Error(t, err)
	assert.Zero(t, badInt)
}

func TestListType(t *testing.T) {

	l := EmptyList
	assert.Zero(t, l.Count())

	rawList := l.Unbox()
	assert.Equal(t, []Value{}, rawList)

	emptyList, err := ListType.Box(rawList)
	assert.NoError(t, err)

	assert.Equal(t, EmptyList, emptyList)
	assert.Equal(t, l, emptyList)

	a, err := IntType.Box(9)
	assert.NoError(t, err)

	b, err := IntType.Box(1)
	assert.NoError(t, err)

	l2 := l.Cons(a).(*List)
	assert.Equal(t, 1, l2.Count().Unbox())

	l3 := l.Cons(b).(*List)
	assert.Equal(t, 1, l3.Count().Unbox())

	ltype := "let-go.lang.PersistentList"
	assert.Equal(t, ltype, l.Type().Name())
	assert.Equal(t, ltype, l2.Type().Name())
	assert.Equal(t, ltype, l3.Type().Name())
	assert.Same(t, l.Type(), l2.Type())

	assert.NotEqual(t, l2, l3)

	l4 := l2.Cons(b).(*List)
	assert.Equal(t, 2, l4.Count().Unbox())

	v2 := []Value{a}
	v3 := []Value{b}
	v4 := []Value{b, a}

	assert.Equal(t, v2, l2.Unbox())
	assert.Equal(t, v3, l3.Unbox())
	assert.Equal(t, v4, l4.Unbox())

	assert.Equal(t, l4.More(), l4.Next())

	assert.Equal(t, NIL, l.First())
	assert.Equal(t, l, l.Next())
	assert.Equal(t, l, l.More())
	assert.Equal(t, EmptyList, l.Next())
	assert.Equal(t, EmptyList, l.More())

	assert.Equal(t, l2, l4.Next())
	assert.True(t, l2 == l4.Next())

	assert.Equal(t, l2.First(), l4.Next().First())

	assert.Equal(t, l3.First(), l4.First())

	n := 100
	values := make([]Value, n)
	for i := 0; i < n; i++ {
		v, err := IntType.Box(rand.Int())
		assert.NoError(t, err)
		if err != nil {
			return
		}
		values[i] = v
	}

	list, err := ListType.Box(values)
	assert.NoError(t, err)
	assert.Equal(t, n, list.(*List).count)

	randomValues := list.Unbox().([]Value)
	for i := 0; i < n; i++ {
		assert.Equal(t, values[i], randomValues[i])
	}

	assert.Equal(t, l, l.Empty())
	assert.Equal(t, EmptyList, l.Empty())
	assert.Equal(t, EmptyList, list.(*List).Empty())
	assert.Equal(t, l, list.(*List).Empty())

	bad := "bad"
	badList, err := ListType.Box(bad)
	assert.Error(t, err)
	assert.Equal(t, EmptyList, badList)
}

func TestSimpleCall(t *testing.T) {

	forty, err := IntType.Box(40)
	assert.NoError(t, err)
	two, err := IntType.Box(2)
	assert.NoError(t, err)

	plus, err := NativeFnType.Box(func(a int, b int) int { return b + a })
	assert.NoError(t, err)

	consts := NewConsts()
	consts.Intern(forty)
	consts.Intern(two)
	consts.Intern(plus)

	c := NewCodeChunk(consts)
	c.maxStack = 4
	c.Append(OP_LOAD_CONST)
	c.Append32(2)

	c.Append(OP_LOAD_CONST)
	c.Append32(0)
	c.Append(OP_LOAD_CONST)
	c.Append32(1)

	c.Append(OP_INVOKE)
	c.Append32(2)

	c.Append(OP_RETURN)

	var out Value

	frame := NewFrame(c, nil)
	out, err = frame.Run()
	assert.NoError(t, err)

	assert.Equal(t, 42, out.Unbox())
}
