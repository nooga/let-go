package bytecode

import "github.com/nooga/let-go/pkg/vm"

// Module is the serializable unit — a complete compilation result.
type Module struct {
	Version uint16
	Flags   uint16
	// Capabilities is an optional feature mask. If FlagCapabilities is set in Flags,
	// a uint32 capability mask follows the header. Bits indicate optional features
	// the decoder must support. Currently no capability bits are defined (all reserved).
	Capabilities uint32
	Strings      []string
	Chunks       []*ChunkData
	Consts       []vm.Value
	// ConstsBase is the starting global index for the consts in this module.
	// For layered pools, indices 0..ConstsBase-1 are in a parent pool.
	ConstsBase int
	// NSTable maps namespace names to their main chunk indices (for bundles).
	NSTable map[string]int
}

// ChunkData holds the data for a single code chunk.
type ChunkData struct {
	MaxStack  int
	Code      []int32
	SourceMap []SourceEntry
	// LocalVars is the chunk's local-variable debug table (slot -> name),
	// serialized in an optional section under FlagLocalVars.
	LocalVars []LocalVarEntry
}

// LocalVarEntry is a local-variable debug entry for serialization.
type LocalVarEntry struct {
	Slot int
	Name string
}

// SourceEntry is a source map entry for serialization.
type SourceEntry struct {
	StartIP                          int
	File                             string
	Line, Column, EndLine, EndColumn int
}
