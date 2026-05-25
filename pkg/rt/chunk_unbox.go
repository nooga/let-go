/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// chunkFromBoxed extracts a *vm.CodeChunk from a Lisp Value. The IR
// chunk primitives generated from examples/go-gen/ir_bridge.lg use
// this as their unbox helper when an arg is declared :Self (the
// receiver of a chunk method).
func chunkFromBoxed(v vm.Value) (*vm.CodeChunk, error) {
	b, ok := v.(*vm.Boxed)
	if !ok {
		return nil, fmt.Errorf("expected boxed *vm.CodeChunk, got %s", v.Type().Name())
	}
	c, ok := b.Unbox().(*vm.CodeChunk)
	if !ok {
		return nil, fmt.Errorf("expected boxed *vm.CodeChunk, got boxed %T", b.Unbox())
	}
	return c, nil
}
