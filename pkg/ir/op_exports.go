/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

// Public wrappers for helpers generated into op_generated.go. The
// generated file is rewritten by examples/go-gen/ir_ops.lg and must
// not be hand-edited; keep these exports here so the generator stays
// pristine.

// IROpToBytecode returns the bytecode OP_* constant corresponding to
// the given IR op. Returns vm.OP_NOOP for ops without a direct
// bytecode counterpart.
func IROpToBytecode(op Op) int32 { return irOpToBytecode(op) }

// IsCheapMaterializeOp reports whether op is one of the cheap-to-
// re-emit loads (Const, LoadArg, LoadVar, LoadClosed). Used by the
// lowerer's body-emit policy and use-site materialization fallback.
func IsCheapMaterializeOp(op Op) bool { return isCheapMaterializeOp(op) }
