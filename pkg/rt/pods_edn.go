/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"strings"

	"github.com/nooga/let-go/pkg/rt/internal/jsonenc"
	"github.com/nooga/let-go/pkg/vm"
)

// EDN payload helpers for pods, and the compiler hooks pods use. Kept apart
// from transit.go so they survive the lg_no_json build tag.

// prStr formats a value as EDN (used for EDN payload encoding).
func prStr(v vm.Value) string {
	switch v.Type() {
	case vm.StringType:
		return string(jsonenc.AppendString(nil, string(v.(vm.String))))
	case vm.NilType:
		return "nil"
	case vm.BooleanType:
		if bool(v.(vm.Boolean)) {
			return "true"
		}
		return "false"
	default:
		return v.String()
	}
}

// evalInNS is set by the compiler package to evaluate code in a namespace.
var evalInNS func(code string, ns *vm.Namespace) (vm.Value, error)

// SetEvalInNS sets the namespace-aware eval function.
func SetEvalInNS(fn func(string, *vm.Namespace) (vm.Value, error)) {
	evalInNS = fn
}

// EDNEncodeArgs encodes args as an EDN vector string.
func EDNEncodeArgs(args []vm.Value) (string, error) {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = prStr(a)
	}
	return "[" + strings.Join(parts, " ") + "]", nil
}
