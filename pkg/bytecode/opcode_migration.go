package bytecode

import "github.com/nooga/let-go/pkg/vm"

// Opcode migration: signature-keyed fallback for bundles compiled with older
// opcode sets.
//
// When a bundle's opcode-set signature doesn't match the running VM's, the
// decoder checks the migration registry. If a remap function is registered for
// the bundle's signature, it is applied to the decoded chunks before execution.
// Otherwise, the bundle is rejected with a clear error.
//
// Each migration is frozen at the time the opcode set changed. Do not update
// existing migrations; add a new one when opcodes change again.

// v2PreRemovalOpcodeNames is the frozen opcode-name list from before
// OP_TRACE_ENABLE and OP_TRACE_DISABLE were removed. It matches the state in
// main at the time of the removal. Used to compute the pre-removal signature
// for v1/v2 → current migration.
var v2PreRemovalOpcodeNames = []string{
	"NOOP",
	"LOAD_CONST",
	"LOAD_ARG",
	"INVOKE",
	"RETURN",
	"BRANCH_T",
	"BRANCH_F",
	"JUMP",
	"POP",
	"POP_N",
	"DUP_NTH",
	"SET_VAR",
	"LOAD_VAR",
	"MAKE_CLOSURE",
	"LOAD_CLOSEDOVER",
	"PUSH_CLOSEDOVER",
	"RECUR",
	"RECUR_FN",
	"TRACE_ENABLE",
	"TRACE_DISABLE",
	"MAKE_MULTI_ARITY",
	"TAIL_CALL",
	"TRY_PUSH",
	"TRY_POP",
	"THROW",
	"ADD",
	"SUB",
	"MUL",
	"BIT_AND",
	"BIT_OR",
	"BIT_XOR",
	"BIT_AND_NOT",
	"BIT_SHIFT_LEFT",
	"BIT_SHIFT_RIGHT",
	"UNSIGNED_BIT_SHIFT_RIGHT",
	"LT",
	"LTE",
	"GT",
	"GTE",
	"EQ",
	"INC",
	"DEC",
	"BIT_NOT",
	"QUOT",
	"DIV",
}

// preRemovalSignature computes the opcode-set signature that would have been
// produced by the VM when TRACE_ENABLE and TRACE_DISABLE were still in the
// opcode list. Used to identify old bundles that need migration.
func preRemovalSignature() (count int, hash uint64) {
	return vm.ComputeSignatureForNames(v2PreRemovalOpcodeNames)
}

const (
	v2opTraceEnable  int32 = 18
	v2opTraceDisable int32 = 19
	v2opFirstShifted int32 = 20 // OP_MAKE_MULTI_ARITY in v2 numbering
)

// v2Stride returns the width, in int32 words, of a v2 opcode's instruction.
// Mirrors pkg/rt/disasm.go opcodeStride as of FormatVersion 2.
func v2Stride(op int32) int {
	switch op & 0xff {
	case 16: // OP_RECUR (offset, argc)
		return 4
	case 22: // OP_TRY_PUSH (catchOffset, finallyOffset)
		return 3
	case 1, 2, 3, 5, 6, 7, 9, 10, 12, 14, 17, 20, 21:
		// LOAD_CONST, LOAD_ARG, INVOKE, BRANCH_TRUE, BRANCH_FALSE, JUMP,
		// POP_N, DUP_NTH, LOAD_VAR, LOAD_CLOSEDOVER, RECUR_FN,
		// MAKE_MULTI_ARITY, TAIL_CALL — one inline argument word.
		return 2
	default:
		return 1
	}
}

// remapV2Opcode maps a v2 opcode byte to its current value. TRACE_ENABLE/DISABLE
// become NOOP (also stride 1, so instruction width is preserved).
//
// This remap encodes the CURRENT enum numbering on its output side: it is
// only valid while later opcode changes are pure suffix appends. An insertion,
// deletion, or reordering anywhere in the enum invalidates every registered
// remap's output values — such a change must freeze the then-current name
// list (like v2PreRemovalOpcodeNames), register a new migration for the old
// signature, and REWRITE the older remaps to target the new numbering (or
// compose them through the frozen name tables).
func remapV2Opcode(op int32) int32 {
	switch {
	case op == v2opTraceEnable || op == v2opTraceDisable:
		return 0 // OP_NOOP
	case op >= v2opFirstShifted:
		return op - 2
	default:
		return op
	}
}

// remapLegacyChunks rewrites the opcode word of every instruction in each chunk
// from v2 numbering to current, in place. It walks by v2 stride so argument words
// are left untouched, and only the low byte (the opcode) of an opcode word is
// changed — packed high bits (e.g. the stack-pointer hint) are preserved. Every
// remap preserves instruction width (TRACE→NOOP is stride 1), so IPs, jump
// targets, and source-map offsets remain valid.
func remapLegacyChunks(chunks []*vm.CodeChunk) {
	for _, ch := range chunks {
		code := ch.Code()
		for i := 0; i < len(code); {
			op := code[i] & 0xff
			stride := v2Stride(op)
			if stride < 1 {
				stride = 1
			}
			code[i] = (code[i] &^ 0xff) | remapV2Opcode(op)
			i += stride
		}
	}
}

// A remapFunc applies a migration to chunks when a bundle's opcode-set signature
// doesn't match the current runtime signature.
type remapFunc func([]*vm.CodeChunk)

// migrationRegistry maps from an opcode-set signature (count, hash) to its
// remap function. When a bundle's effective signature (from CapOpcodeSet or
// inferred as pre-removal) doesn't match the current runtime, the decoder
// looks it up here. If found, the remap is applied; otherwise, the bundle
// is rejected.
var migrationRegistry = map[[2]interface{}]remapFunc{}

func init() {
	// Register the v2-with-TRACE → current migration.
	preCount, preHash := preRemovalSignature()
	// The registry key is [count, hash]. We use a generic key type
	// ([2]interface{}) to support the signature tuple.
	migrationRegistry[[2]interface{}{preCount, preHash}] = remapLegacyChunks
}

// lookupMigration returns the remap function for a signature, or nil if not found.
func lookupMigration(count int, hash uint64) remapFunc {
	key := [2]interface{}{count, hash}
	return migrationRegistry[key]
}
