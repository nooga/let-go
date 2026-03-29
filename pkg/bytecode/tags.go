package bytecode

// Magic bytes identifying an LGB file.
var Magic = [4]byte{'L', 'G', 'B', 0x01}

// FormatVersion is the current serialization format version.
const FormatVersion uint16 = 1

// Module flags.
const (
	FlagConstsBase uint16 = 1 << 0 // ConstsBase field is present in consts section
)

// Type tags for const pool entries.
const (
	TagNil       byte = 0x00
	TagTrue      byte = 0x01
	TagFalse     byte = 0x02
	TagInt       byte = 0x03
	TagFloat     byte = 0x04
	TagString    byte = 0x05
	TagKeyword   byte = 0x06
	TagSymbol    byte = 0x07
	TagChar      byte = 0x08
	TagBigInt    byte = 0x09
	TagVoid      byte = 0x0A
	TagFunc      byte = 0x10
	TagVarRef    byte = 0x11
	TagEmptyList byte = 0x20
	TagList      byte = 0x21
	TagVector    byte = 0x22
	TagMap       byte = 0x23
	TagSet       byte = 0x24
	TagRecordType byte = 0x30
	TagRecord     byte = 0x31
	TagRegex      byte = 0x32
	TagAtom       byte = 0x33
)
