package bytecode

import (
	"fmt"
	"io"
	"math/big"
	"regexp"

	"github.com/nooga/let-go/pkg/vm"
)

// VarResolver resolves a var reference by namespace and name.
type VarResolver func(ns, name string) *vm.Var

// ExecUnit is a decoded compilation unit ready for execution.
type ExecUnit struct {
	Consts    *vm.Consts
	MainChunk *vm.CodeChunk
	// NSChunks maps namespace names to their main chunks (for bundles).
	NSChunks map[string]*vm.CodeChunk
}

// Decode reads a binary module from r.
func Decode(r io.Reader) (*Module, error) {
	return DecodeWithResolver(r, nil)
}

// DecodeToExecUnit decodes an LGB module and returns a ready-to-execute unit.
// The main chunk is chunk index 0. All decoded consts are populated into a
// shared Consts pool that all chunks reference.
func DecodeToExecUnit(r io.Reader, resolve VarResolver) (*ExecUnit, error) {
	return DecodeToExecUnitWithParent(r, resolve, nil)
}

// DecodeToExecUnitWithParent decodes an LGB module with an optional parent const pool.
// If parent is non-nil and the module has a ConstsBase, the decoded consts are layered
// on top of the parent pool — indices < base resolve from the parent.
func DecodeToExecUnitWithParent(r io.Reader, resolve VarResolver, parent *vm.Consts) (*ExecUnit, error) {
	d := &decoder{
		r:       NewReader(r),
		resolve: resolve,
	}

	_, flags, err := d.readHeader()
	if err != nil {
		return nil, err
	}
	d.flags = flags
	strings, err := d.readStringTable()
	if err != nil {
		return nil, err
	}
	d.strings = strings

	chunkDatas, err := d.readChunks()
	if err != nil {
		return nil, err
	}

	// Build the const pool — layered if parent provided
	var sharedConsts *vm.Consts
	if parent != nil {
		sharedConsts = vm.NewChildConsts(parent)
	} else {
		sharedConsts = vm.NewConsts()
	}

	// Build live CodeChunks first — readConsts needs d.chunks for Func resolution
	d.chunks = make([]*vm.CodeChunk, len(chunkDatas))
	for i, cd := range chunkDatas {
		chunk := vm.NewCodeChunk(sharedConsts)
		chunk.Append(cd.Code...)
		chunk.SetMaxStack(cd.MaxStack)
		if len(cd.SourceMap) > 0 {
			for _, e := range cd.SourceMap {
				chunk.AddSourceInfo(vm.SourceInfo{
					File:      e.File,
					Line:      e.Line,
					Column:    e.Column,
					EndLine:   e.EndLine,
					EndColumn: e.EndColumn,
				})
			}
		}
		d.chunks[i] = chunk
	}

	// Now decode consts (Func entries reference d.chunks)
	consts, err := d.readConsts()
	if err != nil {
		return nil, err
	}
	for _, v := range consts {
		sharedConsts.Append(v)
	}

	// Read NS table
	nsTable, err := d.readNSTable()
	if err != nil {
		return nil, err
	}

	if len(d.chunks) == 0 {
		return nil, fmt.Errorf("no chunks in module")
	}

	unit := &ExecUnit{
		Consts:    sharedConsts,
		MainChunk: d.chunks[0],
	}

	// If NS table is present, resolve chunk indices to live CodeChunks
	if len(nsTable) > 0 {
		unit.NSChunks = make(map[string]*vm.CodeChunk, len(nsTable))
		for name, idx := range nsTable {
			if idx >= len(d.chunks) {
				return nil, fmt.Errorf("NS table chunk index %d out of range for %q", idx, name)
			}
			unit.NSChunks[name] = d.chunks[idx]
		}
		// MainChunk is core's chunk if present
		if coreChunk, ok := unit.NSChunks["core"]; ok {
			unit.MainChunk = coreChunk
		}
	}

	return unit, nil
}

// DecodeWithResolver reads a binary module, resolving var references with the given function.
func DecodeWithResolver(r io.Reader, resolve VarResolver) (*Module, error) {
	d := &decoder{
		r:       NewReader(r),
		resolve: resolve,
	}
	m, err := d.readModule()
	if err != nil {
		return nil, err
	}
	return m, nil
}

type decoder struct {
	r          *Reader
	resolve    VarResolver
	flags      uint16
	constsBase int
	strings    []string
	chunks     []*vm.CodeChunk
}

func (d *decoder) readModule() (*Module, error) {
	version, flags, err := d.readHeader()
	if err != nil {
		return nil, err
	}
	d.flags = flags
	strings, err := d.readStringTable()
	if err != nil {
		return nil, err
	}
	d.strings = strings

	chunkDatas, err := d.readChunks()
	if err != nil {
		return nil, err
	}

	// Build live CodeChunk objects for func resolution
	sharedConsts := vm.NewConsts()
	d.chunks = make([]*vm.CodeChunk, len(chunkDatas))
	for i, cd := range chunkDatas {
		chunk := vm.NewCodeChunk(sharedConsts)
		chunk.Append(cd.Code...)
		chunk.SetMaxStack(cd.MaxStack)
		if len(cd.SourceMap) > 0 {
			for _, e := range cd.SourceMap {
				chunk.AddSourceInfo(vm.SourceInfo{
					File:      e.File,
					Line:      e.Line,
					Column:    e.Column,
					EndLine:   e.EndLine,
					EndColumn: e.EndColumn,
				})
			}
		}
		d.chunks[i] = chunk
	}

	consts, err := d.readConsts()
	if err != nil {
		return nil, err
	}

	nsTable, err := d.readNSTable()
	if err != nil {
		return nil, err
	}

	return &Module{
		Version:    version,
		Flags:      flags,
		Strings:    strings,
		Chunks:     chunkDatas,
		Consts:     consts,
		ConstsBase: d.constsBase,
		NSTable:    nsTable,
	}, nil
}

func (d *decoder) readHeader() (version, flags uint16, err error) {
	magic, err := d.r.ReadBytes(4)
	if err != nil {
		return 0, 0, fmt.Errorf("reading magic: %w", err)
	}
	if magic[0] != Magic[0] || magic[1] != Magic[1] || magic[2] != Magic[2] || magic[3] != Magic[3] {
		return 0, 0, fmt.Errorf("invalid magic bytes: %x", magic)
	}
	version, err = d.r.ReadUint16()
	if err != nil {
		return 0, 0, fmt.Errorf("reading version: %w", err)
	}
	flags, err = d.r.ReadUint16()
	if err != nil {
		return 0, 0, fmt.Errorf("reading flags: %w", err)
	}
	return version, flags, nil
}

func (d *decoder) readStringTable() ([]string, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading string count: %w", err)
	}
	strings := make([]string, count)
	for i := range strings {
		slen, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading string length: %w", err)
		}
		b, err := d.r.ReadBytes(int(slen))
		if err != nil {
			return nil, fmt.Errorf("reading string data: %w", err)
		}
		strings[i] = string(b)
	}
	return strings, nil
}

func (d *decoder) readStringRef() (string, error) {
	idx, err := d.r.ReadVarint()
	if err != nil {
		return "", err
	}
	if int(idx) >= len(d.strings) {
		return "", fmt.Errorf("string ref %d out of range (have %d)", idx, len(d.strings))
	}
	return d.strings[idx], nil
}

func (d *decoder) readChunks() ([]*ChunkData, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading chunk count: %w", err)
	}
	chunks := make([]*ChunkData, count)
	for i := range chunks {
		ch := &ChunkData{}
		ms, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading max_stack: %w", err)
		}
		ch.MaxStack = int(ms)

		codeLen, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading code_len: %w", err)
		}
		ch.Code = make([]int32, codeLen)
		for j := range ch.Code {
			ch.Code[j], err = d.r.ReadInt32()
			if err != nil {
				return nil, fmt.Errorf("reading code[%d]: %w", j, err)
			}
		}

		smCount, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading source_map count: %w", err)
		}
		ch.SourceMap = make([]SourceEntry, smCount)
		for j := range ch.SourceMap {
			startIP, err := d.r.ReadVarint()
			if err != nil {
				return nil, err
			}
			file, err := d.readStringRef()
			if err != nil {
				return nil, err
			}
			line, err := d.r.ReadVarint()
			if err != nil {
				return nil, err
			}
			col, err := d.r.ReadVarint()
			if err != nil {
				return nil, err
			}
			eline, err := d.r.ReadVarint()
			if err != nil {
				return nil, err
			}
			ecol, err := d.r.ReadVarint()
			if err != nil {
				return nil, err
			}
			ch.SourceMap[j] = SourceEntry{
				StartIP:   int(startIP),
				File:      file,
				Line:      int(line),
				Column:    int(col),
				EndLine:   int(eline),
				EndColumn: int(ecol),
			}
		}
		chunks[i] = ch
	}
	return chunks, nil
}

func (d *decoder) readConsts() ([]vm.Value, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading const count: %w", err)
	}
	// Read base offset if flag is set
	if d.flags&FlagConstsBase != 0 {
		base, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading consts base: %w", err)
		}
		d.constsBase = int(base)
	}
	consts := make([]vm.Value, count)
	for i := range consts {
		v, err := d.readValue()
		if err != nil {
			return nil, fmt.Errorf("reading const[%d]: %w", i, err)
		}
		consts[i] = v
	}
	return consts, nil
}

func (d *decoder) readNSTable() (map[string]int, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		// EOF is OK — old format modules don't have NS tables
		return nil, nil
	}
	if count == 0 {
		return nil, nil
	}
	table := make(map[string]int, count)
	for i := 0; i < int(count); i++ {
		name, err := d.readStringRef()
		if err != nil {
			return nil, fmt.Errorf("reading NS table name[%d]: %w", i, err)
		}
		chunkIdx, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading NS table chunk index[%d]: %w", i, err)
		}
		table[name] = int(chunkIdx)
	}
	return table, nil
}

func (d *decoder) readValue() (vm.Value, error) {
	tag, err := d.r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("reading tag: %w", err)
	}
	switch tag {
	case TagNil:
		return vm.NIL, nil
	case TagTrue:
		return vm.TRUE, nil
	case TagFalse:
		return vm.FALSE, nil
	case TagInt:
		v, err := d.r.ReadSvarint()
		if err != nil {
			return nil, err
		}
		return vm.Int(v), nil
	case TagFloat:
		v, err := d.r.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return vm.Float(v), nil
	case TagString:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		return vm.String(s), nil
	case TagKeyword:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		return vm.Keyword(s), nil
	case TagSymbol:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		return vm.Symbol(s), nil
	case TagChar:
		v, err := d.r.ReadInt32()
		if err != nil {
			return nil, err
		}
		return vm.Char(v), nil
	case TagBigInt:
		sign, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		magLen, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		mag, err := d.r.ReadBytes(int(magLen))
		if err != nil {
			return nil, err
		}
		bi := new(big.Int).SetBytes(mag)
		if sign != 0 {
			bi.Neg(bi)
		}
		return vm.NewBigInt(bi), nil
	case TagVoid:
		return vm.VOID, nil
	case TagFunc:
		chunkIdx, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		if int(chunkIdx) >= len(d.chunks) {
			return nil, fmt.Errorf("chunk index %d out of range (have %d)", chunkIdx, len(d.chunks))
		}
		arity, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		variadic, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		name, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		fn := vm.MakeFunc(int(arity), variadic != 0, d.chunks[chunkIdx])
		fn.SetName(name)
		return fn, nil
	case TagVarRef:
		ns, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		name, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		if d.resolve != nil {
			v := d.resolve(ns, name)
			if v != nil {
				return v, nil
			}
		}
		// Return a placeholder var if no resolver
		return vm.NewVar(nil, ns, name), nil
	case TagEmptyList:
		return vm.EmptyList, nil
	case TagList:
		count, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		items := make([]vm.Value, count)
		for i := range items {
			items[i], err = d.readValue()
			if err != nil {
				return nil, err
			}
		}
		result, _ := vm.ListType.Box(items)
		return result, nil
	case TagVector:
		count, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		items := make(vm.ArrayVector, count)
		for i := range items {
			items[i], err = d.readValue()
			if err != nil {
				return nil, err
			}
		}
		return items, nil
	case TagMap:
		return d.readMapValue()
	case TagSet:
		count, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		items := make([]vm.Value, count)
		for i := range items {
			items[i], err = d.readValue()
			if err != nil {
				return nil, err
			}
		}
		return vm.NewPersistentSet(items), nil
	case TagRecordType:
		name, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		fieldCount, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		fields := make([]vm.Keyword, fieldCount)
		for i := range fields {
			s, err := d.readStringRef()
			if err != nil {
				return nil, err
			}
			fields[i] = vm.Keyword(s)
		}
		return vm.NewRecordType(name, fields), nil
	case TagRecord:
		// Read the record type inline
		typeName, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		fieldCount, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		fieldKws := make([]vm.Keyword, fieldCount)
		for i := range fieldKws {
			s, err := d.readStringRef()
			if err != nil {
				return nil, err
			}
			fieldKws[i] = vm.Keyword(s)
		}
		rt := vm.NewRecordType(typeName, fieldKws)
		// Read fixed field values
		fixedFields := make([]vm.Value, fieldCount)
		for i := range fixedFields {
			fixedFields[i], err = d.readValue()
			if err != nil {
				return nil, err
			}
		}
		// Read extra map
		extraMap, err := d.readMapValue()
		if err != nil {
			return nil, err
		}
		// Build the data map from fields + extra
		data := extraMap.(*vm.PersistentMap)
		for i, kw := range fieldKws {
			if fixedFields[i] != vm.NIL {
				data = data.Assoc(kw, fixedFields[i]).(*vm.PersistentMap)
			}
		}
		return vm.NewRecord(rt, data), nil
	case TagRegex:
		pattern, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("recompiling regex %q: %w", pattern, err)
		}
		v, _ := vm.RegexType.Box(re)
		return v, nil
	case TagAtom:
		val, err := d.readValue()
		if err != nil {
			return nil, err
		}
		return vm.NewAtom(val), nil
	default:
		return nil, fmt.Errorf("unknown tag 0x%02x", tag)
	}
}

func (d *decoder) readMapValue() (vm.Value, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, err
	}
	m := vm.EmptyPersistentMap
	for i := 0; i < int(count); i++ {
		k, err := d.readValue()
		if err != nil {
			return nil, err
		}
		v, err := d.readValue()
		if err != nil {
			return nil, err
		}
		m = m.Assoc(k, v).(*vm.PersistentMap)
	}
	return m, nil
}
