package bytecode

import (
	"bytes"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestEncodeCompilationPreservesFuncChunkIdentity(t *testing.T) {
	for _, compress := range []bool{false, true} {
		t.Run(map[bool]string{false: "plain", true: "compressed"}[compress], func(t *testing.T) {
			consts := vm.NewConsts()
			mainChunk := vm.NewCodeChunk(consts)

			firstChunk := vm.NewCodeChunk(consts)
			firstChunk.Append(vm.OP_RETURN)
			firstChunk.SetMaxStack(1)
			firstChunk.AddSourceInfoAt(0, vm.SourceInfo{File: "first.lg", Line: 1})
			first := vm.MakeFunc(0, false, firstChunk)
			first.SetName("first")
			consts.Intern(first)

			secondChunk := vm.NewCodeChunk(consts)
			secondChunk.Append(vm.OP_RETURN)
			secondChunk.SetMaxStack(1)
			secondChunk.AddSourceInfoAt(0, vm.SourceInfo{File: "second.lg", Line: 2})
			second := vm.MakeFunc(0, false, secondChunk)
			second.SetName("second")
			consts.Intern(second)

			var buf bytes.Buffer
			if err := EncodeCompilationCompressed(&buf, consts, mainChunk, compress); err != nil {
				t.Fatalf("EncodeCompilationCompressed: %v", err)
			}
			unit, err := DecodeToExecUnitBytes(buf.Bytes(), nil)
			if err != nil {
				t.Fatalf("DecodeToExecUnitBytes: %v", err)
			}

			got := make(map[string]string)
			for _, value := range unit.Consts.Values() {
				fn, ok := value.(*vm.Func)
				if !ok {
					continue
				}
				source := fn.Chunk().LookupSource(0)
				if source == nil {
					t.Fatalf("function %q has no source information", fn.FuncName())
				}
				got[fn.FuncName()] = source.File
			}

			if got["first"] != "first.lg" {
				t.Errorf("first function source = %q, want first.lg", got["first"])
			}
			if got["second"] != "second.lg" {
				t.Errorf("second function source = %q, want second.lg", got["second"])
			}
		})
	}
}
