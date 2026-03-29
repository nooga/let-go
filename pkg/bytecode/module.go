package bytecode

import "github.com/nooga/let-go/pkg/vm"

// Module is the serializable unit — a complete compilation result.
type Module struct {
	Version uint16
	Flags   uint16
	Strings []string
	Chunks  []*ChunkData
	Consts  []vm.Value
	// NSTable maps namespace names to their main chunk indices (for bundles).
	NSTable map[string]int
}

// ChunkData holds the data for a single code chunk.
type ChunkData struct {
	MaxStack  int
	Code      []int32
	SourceMap []SourceEntry
}

// SourceEntry is a source map entry for serialization.
type SourceEntry struct {
	StartIP                          int
	File                             string
	Line, Column, EndLine, EndColumn int
}
