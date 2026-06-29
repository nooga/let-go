package vm

import "testing"

func TestNewCodeChunkWithCapacityPreallocatesCode(t *testing.T) {
	chunk := NewCodeChunkWithCapacity(NewConsts(), 17)
	if len(chunk.code) != 0 {
		t.Fatalf("len(code) = %d, want 0", len(chunk.code))
	}
	if cap(chunk.code) != 17 {
		t.Fatalf("cap(code) = %d, want 17", cap(chunk.code))
	}
}

func TestNewSourceMapWithCapacityPreallocatesEntries(t *testing.T) {
	sm := NewSourceMapWithCapacity(23)
	if len(sm.entries) != 0 {
		t.Fatalf("len(entries) = %d, want 0", len(sm.entries))
	}
	if cap(sm.entries) != 23 {
		t.Fatalf("cap(entries) = %d, want 23", cap(sm.entries))
	}
}

func TestCodeChunkReserveLocalVarsPreallocatesStorage(t *testing.T) {
	chunk := NewCodeChunk(NewConsts())
	chunk.ReserveLocalVars(11)
	if len(chunk.localVars) != 0 {
		t.Fatalf("len(localVars) = %d, want 0", len(chunk.localVars))
	}
	if cap(chunk.localVars) != 11 {
		t.Fatalf("cap(localVars) = %d, want 11", cap(chunk.localVars))
	}
}

func TestCodeChunkReserveLocalVarsPreservesExistingEntries(t *testing.T) {
	chunk := NewCodeChunk(NewConsts())
	chunk.AddLocalVar(3, "x")
	chunk.ReserveLocalVars(11)
	if len(chunk.localVars) != 1 {
		t.Fatalf("len(localVars) = %d, want 1", len(chunk.localVars))
	}
	if chunk.localVars[0].Slot != 3 || chunk.localVars[0].Name != "x" {
		t.Fatalf("localVars[0] = %#v, want slot=3 name=x", chunk.localVars[0])
	}
	if cap(chunk.localVars) != 11 {
		t.Fatalf("cap(localVars) = %d, want 11", cap(chunk.localVars))
	}
}

func TestCodeChunkAddSourceInfoAtUsesProvidedIP(t *testing.T) {
	chunk := NewCodeChunk(NewConsts())
	chunk.AddSourceInfoAt(7, SourceInfo{File: "test.lg", Line: 3})
	if chunk.sourceMap == nil {
		t.Fatal("sourceMap is nil")
	}
	if len(chunk.sourceMap.entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(chunk.sourceMap.entries))
	}
	if chunk.sourceMap.entries[0].startIP != 7 {
		t.Fatalf("startIP = %d, want 7", chunk.sourceMap.entries[0].startIP)
	}
}

func TestConstsReservePreallocatesConstsAndBuckets(t *testing.T) {
	consts := NewConsts()
	consts.Reserve(200)
	if cap(consts.consts) < 200 {
		t.Fatalf("cap(consts) = %d, want >= 200", cap(consts.consts))
	}
	if len(consts.buckets) < 267 {
		t.Fatalf("len(buckets) = %d, want >= 267", len(consts.buckets))
	}
}
