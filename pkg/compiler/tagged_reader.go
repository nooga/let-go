/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TaggedDataReader transforms the next ordinary Lisp form after a registered
// tag, such as #probe [1 2].
type TaggedDataReader func(vm.Value) (vm.Value, error)

type taggedDataReaderResolver func(tag string, validate bool) (TaggedDataReader, bool, error)

func execContextDataReaderResolver(ec *vm.ExecContext) taggedDataReaderResolver {
	return func(tag string, validate bool) (TaggedDataReader, bool, error) {
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
		tagSymbol := vm.Symbol(tag)
		handler := lookup.ValueAt(tagSymbol)
		if handler == nil || handler == vm.NIL {
			return nil, false, nil
		}
		if !validate {
			return nil, true, nil
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

// TaggedRawInput is the narrow cursor exposed to a raw tagged reader. Reads
// update the enclosing Lisp reader's source position.
type TaggedRawInput interface {
	ReadRune() (rune, error)
	UnreadRune() error
}

// TaggedRawReader consumes a tag-specific payload directly from its input. Raw
// readers must consume exactly their payload and be side-effect free because
// reader-conditionals may invoke them while skipping a branch.
type TaggedRawReader func(TaggedRawInput) (vm.Value, error)

type taggedRawInput struct {
	reader *LispReader
}

func (i taggedRawInput) ReadRune() (rune, error) {
	return i.reader.next()
}

func (i taggedRawInput) UnreadRune() error {
	return i.reader.unread()
}

type taggedReader struct {
	data TaggedDataReader
	raw  TaggedRawReader
}

func defaultTaggedReader(tag string) (taggedReader, bool) {
	if tag == "go" {
		return taggedReader{raw: ReadRawGoFragment}, true
	}
	return taggedReader{}, false
}

// TaggedReaderRegistry is a concurrency-safe, per-reader registry of custom
// tagged literal handlers. A tag can have either a data reader or a raw reader;
// the zero value is ready to use.
type TaggedReaderRegistry struct {
	mu      sync.RWMutex
	readers map[string]taggedReader
}

// NewTaggedReaderRegistry returns an empty tagged-reader registry.
func NewTaggedReaderRegistry() *TaggedReaderRegistry {
	return &TaggedReaderRegistry{readers: make(map[string]taggedReader)}
}

// RegisterData registers a handler that receives the next normally-read Lisp
// form. Registering the same tag more than once is an error.
func (r *TaggedReaderRegistry) RegisterData(tag string, reader TaggedDataReader) error {
	if reader == nil {
		return fmt.Errorf("tagged reader #%s is nil", tag)
	}
	return r.register(tag, taggedReader{data: reader})
}

// RegisterRaw registers a handler that consumes its payload directly from the
// Lisp reader. Registering the same tag more than once is an error.
func (r *TaggedReaderRegistry) RegisterRaw(tag string, reader TaggedRawReader) error {
	if reader == nil {
		return fmt.Errorf("tagged reader #%s is nil", tag)
	}
	return r.register(tag, taggedReader{raw: reader})
}

func (r *TaggedReaderRegistry) register(tag string, reader taggedReader) error {
	if r == nil {
		return fmt.Errorf("tagged reader registry is nil")
	}
	if !validTaggedReaderTag(tag) {
		return fmt.Errorf("tagged reader tag %q cannot be read", tag)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.readers == nil {
		r.readers = make(map[string]taggedReader)
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

func (r *TaggedReaderRegistry) lookup(tag string) (taggedReader, bool) {
	if r == nil {
		return taggedReader{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	reader, ok := r.readers[tag]
	return reader, ok
}

type rawGoLexState uint8

const (
	rawGoCode rawGoLexState = iota
	rawGoString
	rawGoRawString
	rawGoRune
	rawGoLineComment
	rawGoBlockComment
)

// ReadRawGoFragment reads the balanced body of a #go{...} tagged literal and
// returns it unchanged as vm.String, excluding the outer braces. Go braces in
// strings, rune literals, raw strings, and comments do not affect balancing.
func ReadRawGoFragment(input TaggedRawInput) (vm.Value, error) {
	opener, err := input.ReadRune()
	for err == nil && isWhitespace(opener) {
		opener, err = input.ReadRune()
	}
	if err != nil {
		if err == io.EOF {
			return vm.NIL, fmt.Errorf("#go requires an opening {")
		}
		return vm.NIL, fmt.Errorf("reading #go fragment: %w", err)
	}
	if opener != '{' {
		return vm.NIL, fmt.Errorf("#go requires an opening {")
	}

	var body strings.Builder
	state := rawGoCode
	depth := 1
	escaped := false
	blockStar := false

	for {
		ch, err := input.ReadRune()
		if err != nil {
			if err == io.EOF {
				return vm.NIL, fmt.Errorf("unterminated #go fragment")
			}
			return vm.NIL, fmt.Errorf("reading #go fragment: %w", err)
		}

		switch state {
		case rawGoCode:
			switch ch {
			case '{':
				depth++
				body.WriteRune(ch)
			case '}':
				depth--
				if depth == 0 {
					return vm.String(body.String()), nil
				}
				body.WriteRune(ch)
			case '"':
				body.WriteRune(ch)
				state = rawGoString
			case '`':
				body.WriteRune(ch)
				state = rawGoRawString
			case '\'':
				body.WriteRune(ch)
				state = rawGoRune
			case '/':
				next, nextErr := input.ReadRune()
				if nextErr != nil {
					if nextErr == io.EOF {
						return vm.NIL, fmt.Errorf("unterminated #go fragment")
					}
					return vm.NIL, fmt.Errorf("reading #go fragment: %w", nextErr)
				}
				switch next {
				case '/':
					body.WriteString("//")
					state = rawGoLineComment
				case '*':
					body.WriteString("/*")
					state = rawGoBlockComment
					blockStar = false
				default:
					body.WriteRune(ch)
					if err := input.UnreadRune(); err != nil {
						return vm.NIL, fmt.Errorf("reading #go fragment: %w", err)
					}
				}
			default:
				body.WriteRune(ch)
			}

		case rawGoString:
			body.WriteRune(ch)
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				state = rawGoCode
			}

		case rawGoRawString:
			body.WriteRune(ch)
			if ch == '`' {
				state = rawGoCode
			}

		case rawGoRune:
			body.WriteRune(ch)
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '\'' {
				state = rawGoCode
			}

		case rawGoLineComment:
			body.WriteRune(ch)
			if ch == '\n' {
				state = rawGoCode
			}

		case rawGoBlockComment:
			body.WriteRune(ch)
			if blockStar && ch == '/' {
				state = rawGoCode
				blockStar = false
			} else {
				blockStar = ch == '*'
			}
		}
	}
}
