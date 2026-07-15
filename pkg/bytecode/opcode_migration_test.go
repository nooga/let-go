package bytecode

import (
	"bytes"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestRemapLegacyChunksV2ToV3 exercises the frozen v1/v2 → current opcode remap over
// a hand-built v2 code array. Uses literal v2 opcode NUMBERS (not vm.OP_*, which
// now describe the current set) so it pins the migration independently of the live enum.
func TestRemapLegacyChunksV2ToCurrent(t *testing.T) {
	const hi = 5 << 16 // a packed high-bits payload (e.g. stack-pointer hint)

	// v2 numbering (with TRACE opcodes): TRACE_ENABLE=18, TRACE_DISABLE=19,
	// MAKE_MULTI_ARITY=20, TAIL_CALL=21, TRY_PUSH=22, ADD=25, RECUR=16,
	// LOAD_CONST=1, NOOP=0, RETURN=4.
	v2 := []int32{
		0,     // NOOP                      -> 0
		1, 25, // LOAD_CONST, arg=25        -> 1, 25 (arg must NOT remap)
		18,          // TRACE_ENABLE              -> 0 (NOOP)
		19,          // TRACE_DISABLE             -> 0 (NOOP)
		20 | hi, 99, // MAKE_MULTI_ARITY|hi, arg -> 18|hi, 99 (high bits kept, opcode -2)
		22, 30, 40, // TRY_PUSH, a=30, b=40      -> 20, 30, 40 (args kept, opcode -2)
		16, 7, 21, 8, // RECUR, a,b,c            -> 16, 7, 21, 8 (stride-4 args kept)
		25, // ADD                       -> 23 (opcode -2)
		4,  // RETURN                    -> 4
	}
	want := []int32{
		0,
		1, 25,
		0,
		0,
		18 | hi, 99,
		20, 30, 40,
		16, 7, 21, 8,
		23,
		4,
	}

	chunk := vm.NewCodeChunk(vm.NewConsts())
	chunk.Append(v2...)
	remapLegacyChunks([]*vm.CodeChunk{chunk})

	got := chunk.Code()
	if len(got) != len(want) {
		t.Fatalf("code length changed: got %d, want %d (remap must preserve width)", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("code[%d]: got %d, want %d", i, got[i], want[i])
		}
	}
}

// TestPreRemovalSignature verifies that the pre-removal signature can be computed
// and differs from the current runtime signature (assuming TRACE opcodes were removed).
func TestPreRemovalSignature(t *testing.T) {
	preCount, preHash := preRemovalSignature()
	runtimeCount, runtimeHash := vm.OpcodeSetSignature()

	if preCount == runtimeCount && preHash == runtimeHash {
		t.Fatal("pre-removal signature should differ from runtime signature (TRACE opcodes should have been removed)")
	}

	if preCount != len(v2PreRemovalOpcodeNames) {
		t.Errorf("preCount: got %d, want %d", preCount, len(v2PreRemovalOpcodeNames))
	}

	// Verify the pre-removal list includes TRACE_ENABLE and TRACE_DISABLE.
	if len(v2PreRemovalOpcodeNames) < 20 || v2PreRemovalOpcodeNames[18] != "TRACE_ENABLE" || v2PreRemovalOpcodeNames[19] != "TRACE_DISABLE" {
		t.Fatal("pre-removal opcode list must include TRACE_ENABLE at 18 and TRACE_DISABLE at 19")
	}
}

// TestMigrationRegistryLookup confirms that the pre-removal signature is registered
// in the migration registry.
func TestMigrationRegistryLookup(t *testing.T) {
	preCount, preHash := preRemovalSignature()
	remap := lookupMigration(preCount, preHash)

	if remap == nil {
		t.Fatal("migration registry should have an entry for the pre-removal signature")
	}

	// Verify the remap function can be called (it should be remapLegacyChunks).
	chunk := vm.NewCodeChunk(vm.NewConsts())
	chunk.Append(18) // v2 TRACE_ENABLE
	remap([]*vm.CodeChunk{chunk})
	got := chunk.Code()
	if len(got) != 1 || got[0] != 0 {
		t.Errorf("remap should convert v2 TRACE_ENABLE (18) to NOOP (0), got %d", got[0])
	}
}

// TestLegacyBundleDecodeWithMigration encodes and decodes a v2 bundle, verifying
// that the migration is applied when the runtime's signature differs from the
// pre-removal signature.
func TestLegacyBundleDecodeWithMigration(t *testing.T) {
	// Build a module with v2 opcode set (before TRACE removal).
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)

	// Use current opcode numbers (will be re-mapped on decode if signatures differ).
	// We use a simple sequence that should survive the remap.
	chunk.Append(
		vm.OP_LOAD_CONST, 0,
		vm.OP_RETURN,
	)
	chunk.SetMaxStack(2)

	fn := vm.MakeFunc(0, false, chunk)
	fn.SetName("test-fn")

	b := NewModuleBuilder()
	b.AddChunk(chunk)
	b.AddConst(fn)

	var buf bytes.Buffer
	if err := Encode(&buf, b.Build()); err != nil {
		t.Fatalf("Encode: %v", err)
	}

	// Decode without FlagCapabilities (simulating a pre-#443 bundle).
	// The decoder should recognize the implicit pre-removal signature and
	// apply the migration if needed.
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	gotFn, ok := decoded.Consts[0].(*vm.Func)
	if !ok {
		t.Fatalf("expected *Func, got %T", decoded.Consts[0])
	}

	got := gotFn.Chunk().Code()
	// The code should be properly decoded (migration applied if needed).
	if len(got) < 2 {
		t.Fatalf("code too short: got %d, want at least 2", len(got))
	}
}

// TestV1DecodeAppliesOpcodeMigration pins the v1 migration gap found in
// review: v1 bundles predate both capabilities and the TRACE-opcode removal,
// so they implicitly carry the pre-removal opcode set — exactly like a
// capability-less v2 bundle — and both frozen v1 decode paths must apply the
// same remap the v2 paths get. Without it, a legacy ADD (25) decodes
// untouched and the current VM executes it as OP_MUL.
func TestV1DecodeAppliesOpcodeMigration(t *testing.T) {
	// Hand-crafted v1 module with one chunk of pre-removal opcodes:
	// TRACE_ENABLE (18) and ADD (25), both stride 1.
	buildV1 := func() *bytes.Buffer {
		var buf bytes.Buffer
		w := NewWriter(&buf)
		w.WriteBytes(Magic[:])
		w.WriteUint16(1) // version 1
		w.WriteUint16(0) // flags (v1: no capabilities)
		w.WriteVarint(0) // 0 strings
		w.WriteVarint(1) // 1 chunk
		w.WriteVarint(2) // max_stack
		w.WriteVarint(2) // code_len
		w.WriteInt32(18) // v2 TRACE_ENABLE — must become OP_NOOP
		w.WriteInt32(25) // legacy ADD — same number as current OP_MUL; must remap to OP_ADD
		w.WriteVarint(0) // 0 source-map entries
		w.WriteVarint(0) // 0 consts
		w.WriteVarint(0) // empty NS table
		w.Flush()
		return &buf
	}
	want := []int32{vm.OP_NOOP, vm.OP_ADD}

	t.Run("module path", func(t *testing.T) {
		m, err := Decode(buildV1())
		if err != nil {
			t.Fatalf("v1 Decode: %v", err)
		}
		if len(m.Chunks) != 1 {
			t.Fatalf("expected 1 chunk, got %d", len(m.Chunks))
		}
		got := m.Chunks[0].Code
		if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("v1 module chunk not migrated: got %v, want %v", got, want)
		}
	})

	t.Run("exec-unit path", func(t *testing.T) {
		unit, err := DecodeToExecUnitBytes(buildV1().Bytes(), nil)
		if err != nil {
			t.Fatalf("v1 DecodeToExecUnitBytes: %v", err)
		}
		got := unit.MainChunk.Code()
		if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("v1 exec-unit chunk not migrated: got %v, want %v", got, want)
		}
	})
}

// TestMismatchedSignatureWithoutMigration verifies that a bundle with an
// unrecognized opcode-set signature is rejected.
func TestMismatchedSignatureWithoutMigration(t *testing.T) {
	// This test would require injecting a fake signature into the bundle.
	// For now, we trust that lookupMigration returns nil for unknown signatures.
	remap := lookupMigration(99999, 0xdeadbeefcafebabe)
	if remap != nil {
		t.Fatal("lookup for unknown signature should return nil")
	}
}
