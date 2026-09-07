package rt

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

func debugFileTestLGB(t *testing.T) []byte {
	t.Helper()
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append32(int(vm.OP_LOAD_CONST))
	chunk.Append32(consts.Intern(vm.Int(42)))
	chunk.Append32(int(vm.OP_RETURN))
	chunk.SetMaxStack(1)
	chunk.AddSourceInfoAt(0, vm.SourceInfo{File: "debug-file-test.lg", Line: 1, Column: 1})

	var encoded bytes.Buffer
	if err := bytecode.EncodeCompilation(&encoded, consts, chunk); err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}

func TestDecodeExecUnitWithDebugFileSkipsSidecarForPlainLGB(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.debug")
	t.Setenv("LG_DEBUG_FILE", missing)
	if _, err := DecodeExecUnitWithDebugFile(debugFileTestLGB(t), "plain.lgb"); err != nil {
		t.Fatalf("plain LGB consulted explicit debug path: %v", err)
	}
}

func TestDecodeExecUnitWithDebugFileLoadsSidecarOnlyForSplitLGB(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.debug")
	t.Setenv("LG_DEBUG_FILE", missing)
	slim, _, err := bytecode.SplitDebug(debugFileTestLGB(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeExecUnitWithDebugFile(slim, "split.lgb"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("split LGB missing-sidecar error = %v, want os.ErrNotExist", err)
	}
}
