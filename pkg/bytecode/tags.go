package bytecode

// Magic bytes identifying an LGB file.
var Magic = [4]byte{'L', 'G', 'B', 0x01}

// FormatVersion is the newest serialization format version. Version 3 adds
// the compressed-body framing. Plain bundles continue to use version 2 so
// compression remains opt-in and their encoding stays byte-identical.
const FormatVersion uint16 = 3

const uncompressedFormatVersion uint16 = 2

// Module flags. Bits are positional via iota — never write an explicit shift.
// Append new flags before flagsEnd; knownFlags is the derived full mask and must
// not be used as the admitted set for older format versions (see vNFlags below).
const (
	FlagConstsBase   uint16 = 1 << iota // ConstsBase field is present in consts section
	FlagCapabilities                    // Capability mask follows the header
	FlagLocalVars                       // per-chunk local-variable debug tables follow the NS table (v2+)
	// FlagCompressed: the module body (everything after the header + capability
	// section) is a single compressed stream, prefixed by its declared
	// uncompressed size and a codec byte. The magic, version, flags, and
	// capability payload — including the opcode-set signature — stay plaintext,
	// so a version/opcode mismatch is still rejected before any inflate.
	// Compression is opt-in at compile time
	// (lg -c -z / lg -b -z); a bundle without this bit decodes byte-identically
	// to before.
	FlagCompressed
	// FlagDebugSplit marks an artifact whose source maps and local-variable
	// tables live in an external, digest-bound debug companion.
	FlagDebugSplit

	flagsEnd // first unused bit; keep last
)

// knownFlags covers every bit assigned in the block above. It is for layout
// checks only — readHeader admits flags per format version via vNFlags.
const knownFlags = flagsEnd - 1

// Per-version admitted flag sets. Appending a flag forces an explicit decision
// about which versions accept it; using knownFlags as the admitted set would
// silently widen what older versions accept.
const (
	v1Flags uint16 = FlagConstsBase | FlagCapabilities
	v2Flags uint16 = v1Flags | FlagLocalVars
	v3Flags uint16 = v2Flags | FlagCompressed | FlagDebugSplit
)

// Compression codecs (the uncompressed byte after the declared body size).
// Kept as a value, not a second flag bit, so a future codec (e.g. zstd) slots
// in without spending flag space or breaking the framing.
const (
	compressionNone  byte = 0 // reserved; a compressed bundle never writes this
	compressionFlate byte = 1 // compress/flate (raw DEFLATE) — stdlib, TinyGo/wasip1-safe
)

// Keep a corrupt or adversarial bundle from making the byte-backed decode path
// allocate without bound. This is intentionally well above current production
// bundles while still providing a deterministic failure before an OOM.
const maxUncompressedBundleBodySize uint64 = 256 << 20

// Capability bits (valid when FlagCapabilities is set).
const (
	// CapOpcodeSet: the capability mask is followed by the producer's opcode-set
	// signature (varint opcode count + uint64 FNV-64a hash of the mnemonics in
	// enum order). The decoder rejects the bundle when the signature differs
	// from the running VM's, so an opcode-enum change between the tree that
	// compiled a bundle and the tree running it fails with a clear message
	// instead of decoding shifted opcodes into undefined behavior.
	CapOpcodeSet uint32 = 1 << 0
)

// Tag byte layout: 0bVV_TTTTTT
//
//	VV     = 2-bit tag version
//	TTTTTT = 6-bit tag ID
const (
	tagVersionShift = 6
	tagVersionMask  = 0b11000000
	tagIDMask       = 0b00111111
)

// Tag versions.
const (
	TagVer0 byte = iota << tagVersionShift
	TagVer1
	TagVer2
	TagVer3
)

// Tag IDs (6-bit). These are the semantic tag identifiers.
const (
	TagIDNil        byte = 0x00
	TagIDTrue       byte = 0x01
	TagIDFalse      byte = 0x02
	TagIDInt        byte = 0x03
	TagIDFloat      byte = 0x04
	TagIDString     byte = 0x05
	TagIDKeyword    byte = 0x06
	TagIDSymbol     byte = 0x07
	TagIDChar       byte = 0x08
	TagIDBigInt     byte = 0x09
	TagIDVoid       byte = 0x0A
	TagIDUUID       byte = 0x0B
	TagIDInstant    byte = 0x0C
	TagIDFunc       byte = 0x10
	TagIDVarRef     byte = 0x11
	TagIDEmptyList  byte = 0x20
	TagIDList       byte = 0x21
	TagIDVector     byte = 0x22
	TagIDMap        byte = 0x23
	TagIDSet        byte = 0x24
	TagIDRecordType byte = 0x30
	TagIDRecord     byte = 0x31
	TagIDRegex      byte = 0x32
	TagIDAtom       byte = 0x33
)

// Tag byte values (version 0). These are the actual bytes on the wire
// for the initial v2 release. They are byte-identical to v1 tags.
const (
	TagNil        byte = TagIDNil | TagVer0
	TagTrue       byte = TagIDTrue | TagVer0
	TagFalse      byte = TagIDFalse | TagVer0
	TagInt        byte = TagIDInt | TagVer0
	TagFloat      byte = TagIDFloat | TagVer0
	TagString     byte = TagIDString | TagVer0
	TagKeyword    byte = TagIDKeyword | TagVer0
	TagSymbol     byte = TagIDSymbol | TagVer0
	TagChar       byte = TagIDChar | TagVer0
	TagBigInt     byte = TagIDBigInt | TagVer0
	TagVoid       byte = TagIDVoid | TagVer0
	TagUUID       byte = TagIDUUID | TagVer0
	TagInstant    byte = TagIDInstant | TagVer0
	TagFunc       byte = TagIDFunc | TagVer0
	TagVarRef     byte = TagIDVarRef | TagVer0
	TagEmptyList  byte = TagIDEmptyList | TagVer0
	TagList       byte = TagIDList | TagVer0
	TagVector     byte = TagIDVector | TagVer0
	TagMap        byte = TagIDMap | TagVer0
	TagSet        byte = TagIDSet | TagVer0
	TagRecordType byte = TagIDRecordType | TagVer0
	TagRecord     byte = TagIDRecord | TagVer0
	TagRegex      byte = TagIDRegex | TagVer0
	TagAtom       byte = TagIDAtom | TagVer0
)

// Reserved tag IDs for future standard tags (0x34–0x3F).
const (
	TagIDReserved0  byte = 0x34
	TagIDReserved1  byte = 0x35
	TagIDReserved2  byte = 0x36
	TagIDReserved3  byte = 0x37
	TagIDReserved4  byte = 0x38
	TagIDReserved5  byte = 0x39
	TagIDReserved6  byte = 0x3A
	TagIDReserved7  byte = 0x3B
	TagIDReserved8  byte = 0x3C
	TagIDReserved9  byte = 0x3D
	TagIDReserved10 byte = 0x3E
	TagIDReserved11 byte = 0x3F
)
