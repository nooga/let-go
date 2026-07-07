package rt

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// asBytes backs the binary file/stream sinks (spit, write!). The byte-array case
// is the one that used to be impossible: a byte-array handed to spit/write! would
// stringify to its #byte-array[…] repr, so bytes >127 could never be written.
func TestAsBytes(t *testing.T) {
	// String → its bytes verbatim, including a high byte.
	if b, ok := asBytes(vm.String("hi\xff")); !ok || string(b) != "hi\xff" {
		t.Fatalf("String: got %q ok=%v", b, ok)
	}

	// byte-array → raw bytes, high bytes preserved (PNG signature + 0xFF/0xC8).
	want := []byte{137, 80, 78, 71, 0, 255, 200}
	ba := vm.NewByteArrayFrom(append([]byte(nil), want...))
	if b, ok := asBytes(ba); !ok || !bytes.Equal(b, want) {
		t.Fatalf("byte-array: got %v ok=%v want %v", b, ok, want)
	}

	// A non-byte typed array is not byte-coercible.
	if _, ok := asBytes(vm.NewIntArrayFrom([]int64{1, 2})); ok {
		t.Fatal("int-array should not be byte-coercible")
	}
}

func TestSpitWritesByteArrayVerbatim(t *testing.T) {
	want := []byte{0x89, 0x50, 0x4e, 0x47, 0x00, 0xff, 0xc8}
	path := filepath.Join(t.TempDir(), "out.bin")

	spit := LookupCoreVar("spit")
	if spit == nil {
		t.Fatal("core/spit not found")
	}
	if _, err := spit.Invoke([]vm.Value{
		vm.String(path),
		vm.NewByteArrayFrom(append([]byte(nil), want...)),
	}); err != nil {
		t.Fatalf("spit: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("spit wrote %v, want %v", got, want)
	}
}

func TestWriteWritesByteArrayVerbatim(t *testing.T) {
	want := []byte{0x89, 0x50, 0x4e, 0x47, 0x00, 0xff, 0xc8}
	var dst bytes.Buffer

	write := LookupCoreVar("write!")
	if write == nil {
		t.Fatal("core/write! not found")
	}
	if _, err := write.Invoke([]vm.Value{
		vm.NewBoxed(NewWriterHandle("test", &dst)),
		vm.NewByteArrayFrom(append([]byte(nil), want...)),
	}); err != nil {
		t.Fatalf("write!: %v", err)
	}

	if got := dst.Bytes(); !bytes.Equal(got, want) {
		t.Fatalf("write! wrote %v, want %v", got, want)
	}
}
