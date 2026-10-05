package bytecode

import "github.com/nooga/let-go/pkg/vm"

// Module is the serializable unit — a complete compilation result.
type Module struct {
	Version uint16
	Flags   uint16
	// Capabilities is an optional feature mask. If FlagCapabilities is set in Flags,
	// a uint32 capability mask follows the header. Bits indicate optional features
	// the decoder must support; see KnownCapabilities / CapOpcodeSet.
	Capabilities uint32
	Strings      []string
	Chunks       []*ChunkData
	Consts       []vm.Value
	// ConstsBase is the starting global index for the consts in this module.
	// For layered pools, indices 0..ConstsBase-1 are in a parent pool.
	ConstsBase int
	// NSTable maps namespace names to their main chunk indices (for bundles).
	NSTable map[string]int
	// NSRanges lists, under FlagNSRanges, the namespaces whose main chunk and
	// function chunks no other namespace's code reaches, so a decoder can
	// defer them (see NSRange). Sorted by name.
	NSRanges []NSRange
}

// NSRange marks one bundle namespace as self-contained: MainChunk is its
// NSTable entry, ConstRuns are the const-pool index runs of the constants only
// its code reaches (functions, def metadata, literals, var references), and
// [ChunkLo, ChunkHi) are the chunks of the functions among them — everything a
// decoder skips while the namespace is deferred. ChunkLo == ChunkHi means the
// namespace has no function chunks of its own. The encoder emits an entry only
// after proving, by walking every namespace's code, that no other namespace's
// code loads any of those constants or a function whose chunk lies in the
// range. A constant two namespaces share (values dedupe across compilations)
// is never in a run.
type NSRange struct {
	Name      string
	MainChunk int
	ChunkLo   int
	ChunkHi   int
	ConstRuns []ConstRange
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
