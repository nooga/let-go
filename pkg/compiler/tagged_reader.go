/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"sync"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TaggedDataReader transforms the next ordinary Lisp form after a registered
// tag, such as #probe [1 2].
type TaggedDataReader func(vm.Value) (vm.Value, error)

type taggedDataReaderResolver func(tag string) (TaggedDataReader, bool, error)

func execContextDataReaderResolver(ec *vm.ExecContext) taggedDataReaderResolver {
	return func(tag string) (TaggedDataReader, bool, error) {
		dataReaders := rt.NS(rt.NameCoreNS).Lookup(vm.Symbol("*data-readers*"))
		dataReadersVar, ok := dataReaders.(*vm.Var)
		if !ok {
			return nil, false, nil
		}
		readers := ec.Deref(dataReadersVar)
		if readers == vm.NIL {
			return nil, false, nil
		}
		lookup, ok := readers.(vm.Lookup)
		if !ok || (readers.Type() != vm.MapType && readers.Type() != vm.SortedMapType) {
			return nil, false, fmt.Errorf("*data-readers* must be a map, got %s", readers.Type().Name())
		}
		handler := lookup.ValueAt(vm.Symbol(tag))
		if handler == nil || handler == vm.NIL {
			return nil, false, nil
		}
		if handlerVar, ok := handler.(*vm.Var); ok {
			handler = ec.Deref(handlerVar)
		}
		fn, ok := handler.(vm.Fn)
		if !ok {
			return nil, false, fmt.Errorf("*data-readers* entry for %s must be callable, got %s", tag, handler.Type().Name())
		}
		return func(value vm.Value) (vm.Value, error) {
			return ec.Invoke(fn, []vm.Value{value})
		}, true, nil
	}
}

var rootDataReaderResolver = execContextDataReaderResolver(vm.RootExecContext)

// taggedReaderRegistryVar is a private dynamic binding key. Registry values are
// installed only on child execution contexts, never on RootExecContext, so
// concurrent compiler evaluations cannot observe or overwrite one another.
var taggedReaderRegistryVar = vm.NewVar(nil, "compiler", "*tagged-reader-registry*")

type taggedReaderRegistryBinding struct {
	registry *TaggedReaderRegistry
}

func (b *taggedReaderRegistryBinding) String() string     { return "#<tagged-reader-registry>" }
func (b *taggedReaderRegistryBinding) Type() vm.ValueType { return vm.AnyType }
func (b *taggedReaderRegistryBinding) Unbox() any         { return b.registry }

func taggedReadersFromExecContext(ec *vm.ExecContext) *TaggedReaderRegistry {
	binding, ok := ec.Deref(taggedReaderRegistryVar).(*taggedReaderRegistryBinding)
	if !ok {
		return nil
	}
	return binding.registry
}

func execContextWithTaggedReaders(ec *vm.ExecContext, registry *TaggedReaderRegistry) *vm.ExecContext {
	if registry == nil || taggedReadersFromExecContext(ec) == registry {
		return ec
	}
	child := ec.Child()
	child.PushBinding(taggedReaderRegistryVar, &taggedReaderRegistryBinding{registry: registry})
	return child
}

// TaggedReaderRegistry is a concurrency-safe, per-reader registry of custom
// tagged literal handlers. The zero value is ready to use.
type TaggedReaderRegistry struct {
	mu      sync.RWMutex
	readers map[string]TaggedDataReader
}

// NewTaggedReaderRegistry returns an empty tagged-reader registry.
func NewTaggedReaderRegistry() *TaggedReaderRegistry {
	return &TaggedReaderRegistry{readers: make(map[string]TaggedDataReader)}
}

// RegisterData registers a handler that receives the next normally-read Lisp
// form. Registering the same tag more than once is an error.
func (r *TaggedReaderRegistry) RegisterData(tag string, reader TaggedDataReader) error {
	if reader == nil {
		return fmt.Errorf("tagged reader #%s is nil", tag)
	}
	if r == nil {
		return fmt.Errorf("tagged reader registry is nil")
	}
	if !validTaggedReaderTag(tag) {
		return fmt.Errorf("tagged reader tag %q cannot be read", tag)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.readers == nil {
		r.readers = make(map[string]TaggedDataReader)
	}
	if _, exists := r.readers[tag]; exists {
		return fmt.Errorf("tagged reader #%s is already registered", tag)
	}
	r.readers[tag] = reader
	return nil
}

func validTaggedReaderTag(tag string) bool {
	first := true
	for _, ch := range tag {
		if first {
			if !isLetter(ch) {
				return false
			}
			first = false
			continue
		}
		if isWhitespace(ch) || isTerminatingMacro(ch) {
			return false
		}
	}
	return !first
}

func (r *TaggedReaderRegistry) lookup(tag string) (TaggedDataReader, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	reader, ok := r.readers[tag]
	return reader, ok
}
