package bytecode

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func BenchmarkDecodeCore(b *testing.B) {
	b.ReportAllocs()
	corePath := filepath.Join("..", "rt", "core_compiled.lgb")
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		b.Skip("core_compiled.lgb not found at expected path")
	}
	data, err := os.ReadFile(corePath)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeModuleSmall(b *testing.B) {
	b.ReportAllocs()
	mb := NewModuleBuilder()
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)
	mb.AddChunk(chunk)
	for i := range 10 {
		mb.AddConst(vm.Int(i))
	}
	m := mb.Build()
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeModuleLargeCollections(b *testing.B) {
	b.ReportAllocs()
	mb := NewModuleBuilder()
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)
	mb.AddChunk(chunk)

	// Large map
	m := vm.EmptyPersistentMap
	for i := range 100 {
		m = m.Assoc(vm.Keyword(string(rune('a'+i%26))), vm.Int(i)).(*vm.PersistentMap)
	}
	mb.AddConst(m)

	// Large set
	s := vm.NewPersistentSet(nil)
	for i := range 100 {
		s = s.Conj(vm.Int(i)).(*vm.PersistentSet)
	}
	mb.AddConst(s)

	// Large vector
	vec := make(vm.ArrayVector, 100)
	for i := range 100 {
		vec[i] = vm.Int(i)
	}
	mb.AddConst(vec)

	module := mb.Build()
	var buf bytes.Buffer
	if err := Encode(&buf, module); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkV1vsV2Decode compares decode performance of identical content
// through the frozen v1 path vs the v2 path.
func BenchmarkV1vsV2Decode(b *testing.B) {
	mb := NewModuleBuilder()
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)
	mb.AddChunk(chunk)

	// Mix of scalar and collection values
	mb.AddConst(vm.NIL)
	mb.AddConst(vm.TRUE)
	mb.AddConst(vm.Int(42))
	mb.AddConst(vm.String("hello"))
	mb.AddConst(vm.Keyword("kw"))

	m := mb.Build()
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		b.Fatal(err)
	}
	dataV2 := buf.Bytes()

	// Clone and patch version byte to 1 for v1 path
	dataV1 := make([]byte, len(dataV2))
	copy(dataV1, dataV2)
	// Magic[4] + Version uint16 LE => byte 4 is low byte of version
	if len(dataV1) > 5 && dataV1[4] == 2 {
		dataV1[4] = 1
	}

	b.Run("v1", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := Decode(bytes.NewReader(dataV1))
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("v2", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := Decode(bytes.NewReader(dataV2))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
