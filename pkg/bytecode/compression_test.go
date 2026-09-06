package bytecode

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// buildCompressibleModule returns a module whose string table repeats, so a
// compressed encoding is meaningfully smaller than the plain one (real bundles
// are dominated by recurring symbol/const/opcode bytes).
func buildCompressibleModule() *Module {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)

	fn := vm.MakeFunc(0, false, chunk)
	fn.SetName("compressible-fn")

	b := NewModuleBuilder()
	b.AddChunk(chunk)
	b.AddConst(fn)
	// A repetitive string table is the compressible payload.
	for i := 0; i < 200; i++ {
		b.internString("a-recurring-namespace/symbol-name-that-repeats")
	}
	return b.Build()
}

func TestCompressedRoundtrip(t *testing.T) {
	m := buildCompressibleModule()

	var plain bytes.Buffer
	if err := Encode(&plain, m); err != nil {
		t.Fatalf("plain encode: %v", err)
	}

	mc := buildCompressibleModule()
	mc.Flags |= FlagCompressed
	var comp bytes.Buffer
	if err := Encode(&comp, mc); err != nil {
		t.Fatalf("compressed encode: %v", err)
	}

	if comp.Len() >= plain.Len() {
		t.Fatalf("compressed (%d) not smaller than plain (%d)", comp.Len(), plain.Len())
	}

	// Header stays plaintext: magic + version readable without inflating, so a
	// version/opcode-set mismatch is rejected before any decompression.
	got := comp.Bytes()
	if !bytes.Equal(got[:4], Magic[:]) {
		t.Errorf("compressed header magic = %x, want %x", got[:4], Magic[:])
	}
	if version := binary.LittleEndian.Uint16(got[4:6]); version != FormatVersion {
		t.Errorf("compressed version = %d, want %d", version, FormatVersion)
	}
	if version := binary.LittleEndian.Uint16(plain.Bytes()[4:6]); version != uncompressedFormatVersion {
		t.Errorf("plain version = %d, want %d", version, uncompressedFormatVersion)
	}

	// Decodes back to the same shape as the plain encoding.
	dp, err := Decode(bytes.NewReader(plain.Bytes()))
	if err != nil {
		t.Fatalf("decode plain: %v", err)
	}
	dc, err := Decode(bytes.NewReader(comp.Bytes()))
	if err != nil {
		t.Fatalf("decode compressed: %v", err)
	}
	if dc.Flags&FlagCompressed == 0 {
		t.Error("decoded compressed module lost FlagCompressed")
	}
	if dc.Version != FormatVersion {
		t.Errorf("decoded compressed version = %d, want %d", dc.Version, FormatVersion)
	}
	if len(dc.Chunks) != len(dp.Chunks) || len(dc.Consts) != len(dp.Consts) {
		t.Fatalf("shape mismatch: chunks %d/%d consts %d/%d",
			len(dc.Chunks), len(dp.Chunks), len(dc.Consts), len(dp.Consts))
	}
	fn, ok := dc.Consts[0].(*vm.Func)
	if !ok || fn.FuncName() != "compressible-fn" {
		t.Fatalf("const[0] = %#v, want func compressible-fn", dc.Consts[0])
	}
}

func TestV2RejectsCompressionFlagAtHeader(t *testing.T) {
	m := buildCompressibleModule()
	var plain bytes.Buffer
	if err := Encode(&plain, m); err != nil {
		t.Fatalf("plain encode: %v", err)
	}

	data := append([]byte(nil), plain.Bytes()...)
	flags := binary.LittleEndian.Uint16(data[6:8]) | FlagCompressed
	binary.LittleEndian.PutUint16(data[6:8], flags)
	_, err := Decode(bytes.NewReader(data))
	if err == nil || !strings.Contains(err.Error(), "unsupported LGB flags") {
		t.Fatalf("Decode() error = %v, want unsupported-flags error", err)
	}
}

func writeCompressedTestFrame(t *testing.T, declaredSize uint64, body []byte) []byte {
	t.Helper()
	var framed bytes.Buffer
	w := NewWriter(&framed)
	if err := w.WriteBytes(Magic[:]); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteUint16(FormatVersion); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteUint16(FlagCompressed); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteVarint(declaredSize); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteByte(compressionFlate); err != nil {
		t.Fatal(err)
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
	fw, err := flate.NewWriter(&framed, flate.BestSpeed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(body); err != nil {
		t.Fatal(err)
	}
	if err := fw.Close(); err != nil {
		t.Fatal(err)
	}
	return framed.Bytes()
}

func TestCompressedBodyRejectsOversizedDeclaration(t *testing.T) {
	data := writeCompressedTestFrame(t, maxUncompressedBundleBodySize+1, nil)
	_, err := DecodeToExecUnitBytes(data, nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds limit") {
		t.Fatalf("DecodeToExecUnitBytes() error = %v, want size-limit error", err)
	}
}

func TestCompressedBodyRejectsInflateBeyondDeclaration(t *testing.T) {
	data := writeCompressedTestFrame(t, 1, []byte{0, 1})
	_, err := DecodeToExecUnitBytes(data, nil)
	if err == nil || !strings.Contains(err.Error(), "does not match declared size") {
		t.Fatalf("DecodeToExecUnitBytes() error = %v, want declared-size mismatch", err)
	}
}

type failingWriteCloser struct {
	writeErr error
	closeErr error
	closed   bool
}

func (w *failingWriteCloser) Write([]byte) (int, error) { return 0, w.writeErr }

func (w *failingWriteCloser) Close() error {
	w.closed = true
	return w.closeErr
}

func TestCompressedBodyWriterClosesAfterWriteError(t *testing.T) {
	writeErr := errors.New("write failed")
	closeErr := errors.New("close failed")
	w := &failingWriteCloser{writeErr: writeErr, closeErr: closeErr}
	err := writeAndCloseCompressedBody(w, []byte("body"))
	if !w.closed {
		t.Fatal("compressed writer was not closed after write error")
	}
	if !errors.Is(err, writeErr) || !errors.Is(err, closeErr) {
		t.Fatalf("error = %v, want joined write and close errors", err)
	}
}

var _ io.WriteCloser = (*failingWriteCloser)(nil)

// TestCompressedExecUnitBytePath exercises the byte-backed inflate swap
// (DecodeToExecUnitBytes, the embedded-core path) so the zero-copy source-map
// reader is rebuilt over the inflated buffer rather than the compressed one.
func TestCompressedExecUnitBytePath(t *testing.T) {
	mc := buildCompressibleModule()
	mc.Flags |= FlagCompressed
	var comp bytes.Buffer
	if err := Encode(&comp, mc); err != nil {
		t.Fatalf("compressed encode: %v", err)
	}

	resolve := func(ns, name string) *vm.Var { return nil }

	// Streaming path.
	if _, err := DecodeToExecUnit(bytes.NewReader(comp.Bytes()), resolve); err != nil {
		t.Fatalf("streaming DecodeToExecUnit: %v", err)
	}
	// Byte-backed (zero-copy) path.
	if _, err := DecodeToExecUnitBytes(comp.Bytes(), resolve); err != nil {
		t.Fatalf("byte-backed DecodeToExecUnitBytes: %v", err)
	}
}

// TestEncodeNormalizesDownWithoutCompression covers the sharp edge where a
// decoded v3 module is re-encoded with FlagCompressed cleared: the wire
// version must drop back to the uncompressed write version so older decoders
// are not rejected for a plaintext body.
func TestEncodeNormalizesDownWithoutCompression(t *testing.T) {
	mc := buildCompressibleModule()
	mc.Flags |= FlagCompressed
	var comp bytes.Buffer
	if err := Encode(&comp, mc); err != nil {
		t.Fatalf("compressed encode: %v", err)
	}
	decoded, err := Decode(bytes.NewReader(comp.Bytes()))
	if err != nil {
		t.Fatalf("decode compressed: %v", err)
	}
	if decoded.Version != FormatVersion {
		t.Fatalf("decoded version = %d, want %d", decoded.Version, FormatVersion)
	}

	decoded.Flags &^= FlagCompressed
	var plain bytes.Buffer
	if err := Encode(&plain, decoded); err != nil {
		t.Fatalf("re-encode without compression: %v", err)
	}
	if version := binary.LittleEndian.Uint16(plain.Bytes()[4:6]); version != uncompressedFormatVersion {
		t.Fatalf("re-encoded version = %d, want %d", version, uncompressedFormatVersion)
	}
	if flags := binary.LittleEndian.Uint16(plain.Bytes()[6:8]); flags&FlagCompressed != 0 {
		t.Fatalf("re-encoded flags still have FlagCompressed: 0x%04x", flags)
	}
	if _, err := Decode(bytes.NewReader(plain.Bytes())); err != nil {
		t.Fatalf("decode re-encoded plain: %v", err)
	}
}
