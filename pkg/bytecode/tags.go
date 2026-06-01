package bytecode

// Magic bytes identifying an LGB file.
var Magic = [4]byte{'L', 'G', 'B', 0x01}

// FormatVersion is the current serialization format version.
const FormatVersion uint16 = 2

// Module flags.
const (
	FlagConstsBase   uint16 = 1 << 0 // ConstsBase field is present in consts section
	FlagCapabilities uint16 = 1 << 1 // Capability mask follows the header
	FlagLocalVars    uint16 = 1 << 2 // per-chunk local-variable debug tables follow the NS table (v2+)
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
