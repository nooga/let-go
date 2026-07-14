/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

// MetaValue decorates any Value with metadata while delegating its ordinary
// Value behavior to the wrapped value.
type MetaValue[T Value] struct {
	wrapped T
	meta    Value
}

func NewMetaValue[T Value](wrapped T, meta Value) *MetaValue[T] {
	return &MetaValue[T]{wrapped: wrapped, meta: meta}
}

func (m *MetaValue[T]) Wrapped() T             { return m.wrapped }
func (m *MetaValue[T]) Type() ValueType        { return m.wrapped.Type() }
func (m *MetaValue[T]) Unbox() any             { return m.wrapped.Unbox() }
func (m *MetaValue[T]) String() string         { return m.wrapped.String() }
func (m *MetaValue[T]) WithMeta(v Value) Value { return NewMetaValue(m.wrapped, v) }

func (m *MetaValue[T]) Meta() Value {
	if m.meta == nil {
		return NIL
	}
	return m.meta
}

// MetaFn adds the Fn surface to the generic metadata decorator. ExecContext
// unwraps it before dispatch so closures and context-aware natives still
// receive the caller's dynamic context.
type MetaFn struct {
	*MetaValue[Fn]
}

func NewMetaFn(fn Fn, meta Value) *MetaFn {
	return &MetaFn{MetaValue: NewMetaValue[Fn](fn, meta)}
}

func (f *MetaFn) Invoke(args []Value) (Value, error) { return f.Wrapped().Invoke(args) }
func (f *MetaFn) Arity() int                         { return f.Wrapped().Arity() }
func (f *MetaFn) WithMeta(meta Value) Value          { return NewMetaFn(f.Wrapped(), meta) }
