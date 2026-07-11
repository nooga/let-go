package bytecode

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// buildTestModule returns a minimal builder-produced module (carries the
// opcode-set capability like every newly encoded bundle).
func buildTestModule() *Module {
	b := NewModuleBuilder()
	c := vm.NewCodeChunk(nil)
	c.Append(vm.OP_RETURN)
	b.AddChunk(c)
	return b.Build()
}

func TestOpcodeSetRoundtrip(t *testing.T) {
	m := buildTestModule()
	if m.Flags&FlagCapabilities == 0 || m.Capabilities&CapOpcodeSet == 0 {
		t.Fatal("builder should set FlagCapabilities|CapOpcodeSet on new modules")
	}
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("encode: %v", err)
	}
	dm, err := Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("decode of same-VM bundle should pass the opcode-set check: %v", err)
	}
	if dm.Capabilities&CapOpcodeSet == 0 {
		t.Fatal("decoded module lost CapOpcodeSet")
	}
}

// opcodeSigOffset is where the signature payload starts: magic(4) + version(2)
// + flags(2) + caps mask(4). The varint count comes first, then the u64 hash.
const opcodeSigOffset = 12

func TestOpcodeSetHashMismatch(t *testing.T) {
	m := buildTestModule()
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("encode: %v", err)
	}
	data := buf.Bytes()
	count, _ := vm.OpcodeSetSignature()
	if count >= 0x80 {
		t.Fatalf("test assumes a 1-byte varint opcode count, got %d", count)
	}
	data[opcodeSigOffset+1+3] ^= 0xff // flip a byte inside the hash
	_, err := Decode(bytes.NewReader(data))
	if err == nil {
		t.Fatal("expected opcode-set mismatch error")
	}
	if !strings.Contains(err.Error(), "opcode set mismatch") {
		t.Fatalf("expected 'opcode set mismatch' in error, got: %v", err)
	}
}

func TestOpcodeSetCountMismatch(t *testing.T) {
	m := buildTestModule()
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("encode: %v", err)
	}
	data := buf.Bytes()
	count, _ := vm.OpcodeSetSignature()
	if count >= 0x80 {
		t.Fatalf("test assumes a 1-byte varint opcode count, got %d", count)
	}
	data[opcodeSigOffset]++ // one more opcode than the runtime has
	_, err := Decode(bytes.NewReader(data))
	if err == nil {
		t.Fatal("expected opcode-set mismatch error")
	}
	if !strings.Contains(err.Error(), "opcode set mismatch") {
		t.Fatalf("expected 'opcode set mismatch' in error, got: %v", err)
	}
}

func TestLegacyBundleWithoutOpcodeSetDecodes(t *testing.T) {
	// A pre-guard bundle (no FlagCapabilities at all) must keep decoding.
	m := buildTestModule()
	m.Flags &^= FlagCapabilities
	m.Capabilities = 0
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if _, err := Decode(bytes.NewReader(buf.Bytes())); err != nil {
		t.Fatalf("legacy bundle without capability mask should decode: %v", err)
	}
}
