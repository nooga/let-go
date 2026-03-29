package bytecode

import (
	"fmt"
	"io"
	"math/big"

	"github.com/nooga/let-go/pkg/vm"
)

// Encode serializes a Module to binary format.
func Encode(w io.Writer, m *Module) error {
	enc := &encoder{
		w:        NewWriter(w),
		strings:  m.Strings,
		strIndex: make(map[string]int, len(m.Strings)),
		chunks:   m.Chunks,
	}
	for i, s := range m.Strings {
		enc.strIndex[s] = i
	}
	if err := enc.writeHeader(m); err != nil {
		return err
	}
	if err := enc.writeStringTable(); err != nil {
		return err
	}
	if err := enc.writeChunks(); err != nil {
		return err
	}
	if err := enc.writeConsts(m.Consts); err != nil {
		return err
	}
	if err := enc.writeNSTable(m.NSTable); err != nil {
		return err
	}
	return enc.w.Flush()
}

// EncodeModule builds a Module from live VM objects and serializes it.
func EncodeModule(w io.Writer, consts *vm.Consts, chunks []*vm.CodeChunk) error {
	b := NewModuleBuilder()
	vals := consts.Values()
	for _, c := range chunks {
		b.AddChunk(c)
	}
	for _, v := range vals {
		b.AddConst(v)
	}
	m := b.Build()
	return Encode(w, m)
}

// EncodeCompilation serializes a compilation result (main chunk + const pool).
// The main chunk is always chunk index 0. All CodeChunks referenced by Funcs
// in the const pool are collected automatically.
func EncodeCompilation(w io.Writer, consts *vm.Consts, mainChunk *vm.CodeChunk) error {
	b := NewModuleBuilder()
	// Main chunk must be index 0
	b.AddChunk(mainChunk)
	// Collect all func chunks from the const pool (AddConst interns their chunks)
	vals := consts.Values()
	for _, v := range vals {
		b.AddConst(v)
	}
	m := b.Build()
	return Encode(w, m)
}

// EncodeBundle serializes a multi-namespace compilation bundle.
// nsChunks maps namespace names to their main CodeChunks.
// The "core" entry (chunk index 0) is treated as the entry point.
func EncodeBundle(w io.Writer, consts *vm.Consts, nsChunks map[string]*vm.CodeChunk) error {
	b := NewModuleBuilder()
	// Register all namespace main chunks
	for name, chunk := range nsChunks {
		b.SetNSEntry(name, chunk)
	}
	// Collect all func chunks from the const pool
	vals := consts.Values()
	for _, v := range vals {
		b.AddConst(v)
	}
	m := b.Build()
	return Encode(w, m)
}

// ModuleBuilder collects strings, chunks, and consts for serialization.
type ModuleBuilder struct {
	strIndex   map[string]int
	strings    []string
	chunkIndex map[*vm.CodeChunk]int
	chunks     []*ChunkData
	consts     []vm.Value
	nsTable    map[string]int
}

// NewModuleBuilder creates a new builder.
func NewModuleBuilder() *ModuleBuilder {
	return &ModuleBuilder{
		strIndex:   make(map[string]int),
		chunkIndex: make(map[*vm.CodeChunk]int),
	}
}

func (b *ModuleBuilder) internString(s string) int {
	if idx, ok := b.strIndex[s]; ok {
		return idx
	}
	idx := len(b.strings)
	b.strings = append(b.strings, s)
	b.strIndex[s] = idx
	return idx
}

// AddChunk registers a CodeChunk and interns its strings.
func (b *ModuleBuilder) AddChunk(c *vm.CodeChunk) int {
	if idx, ok := b.chunkIndex[c]; ok {
		return idx
	}
	cd := &ChunkData{
		MaxStack: c.MaxStack(),
		Code:     c.Code(),
	}
	if sm := c.GetSourceMap(); sm != nil {
		for _, e := range sm.Entries() {
			b.internString(e.Info.File)
			cd.SourceMap = append(cd.SourceMap, SourceEntry{
				StartIP:   e.StartIP,
				File:      e.Info.File,
				Line:      e.Info.Line,
				Column:    e.Info.Column,
				EndLine:   e.Info.EndLine,
				EndColumn: e.Info.EndColumn,
			})
		}
	}
	idx := len(b.chunks)
	b.chunks = append(b.chunks, cd)
	b.chunkIndex[c] = idx
	return idx
}

// AddConst registers a const value and interns related strings/chunks.
func (b *ModuleBuilder) AddConst(v vm.Value) {
	b.internStringsForValue(v)
	b.consts = append(b.consts, v)
}

func (b *ModuleBuilder) internStringsForValue(v vm.Value) {
	switch val := v.(type) {
	case vm.String:
		b.internString(string(val))
	case vm.Keyword:
		b.internString(string(val))
	case vm.Symbol:
		b.internString(string(val))
	case *vm.Func:
		b.internString(val.FuncName())
		b.AddChunk(val.Chunk())
	case *vm.Var:
		b.internString(val.NS())
		b.internString(val.VarName())
	case *vm.RecordType:
		b.internString(val.TypeName())
		for _, f := range val.Fields() {
			b.internString(string(f))
		}
	case *vm.Regex:
		b.internString(val.Pattern())
	case *vm.List:
		var s vm.Seq = val
		for s != nil && s != vm.EmptyList {
			b.internStringsForValue(s.First())
			s = s.Next()
		}
	case vm.ArrayVector:
		for _, item := range val {
			b.internStringsForValue(item)
		}
	case *vm.PersistentMap:
		s := val.Seq()
		for s != nil && s != vm.EmptyList {
			entry := s.First().(vm.ArrayVector)
			b.internStringsForValue(entry[0])
			b.internStringsForValue(entry[1])
			s = s.Next()
		}
	case *vm.PersistentSet:
		s := val.Seq()
		for s != nil && s != vm.EmptyList {
			b.internStringsForValue(s.First())
			s = s.Next()
		}
	case *vm.Record:
		rt := val.RecordType()
		b.internString(rt.TypeName())
		for _, f := range rt.Fields() {
			b.internString(string(f))
		}
		for _, fv := range val.FixedFields() {
			if fv != nil {
				b.internStringsForValue(fv)
			}
		}
		b.internStringsForValue(val.Extra())
	case *vm.Atom:
		b.internStringsForValue(val.Deref())
	}
}

// SetNSEntry records a namespace name → main chunk mapping for bundle modules.
func (b *ModuleBuilder) SetNSEntry(name string, chunk *vm.CodeChunk) {
	if b.nsTable == nil {
		b.nsTable = make(map[string]int)
	}
	b.internString(name)
	idx := b.AddChunk(chunk)
	b.nsTable[name] = idx
}

// Build creates the Module.
func (b *ModuleBuilder) Build() *Module {
	return &Module{
		Version: FormatVersion,
		Flags:   0,
		Strings: b.strings,
		Chunks:  b.chunks,
		Consts:  b.consts,
		NSTable: b.nsTable,
	}
}

// ChunkIndex returns the index for a given CodeChunk pointer.
func (b *ModuleBuilder) ChunkIndex(c *vm.CodeChunk) (int, bool) {
	idx, ok := b.chunkIndex[c]
	return idx, ok
}

type encoder struct {
	w        *Writer
	strings  []string
	strIndex map[string]int
	chunks   []*ChunkData
	// chunkMap maps live CodeChunk pointers to chunk indices (populated by EncodeModule path)
}

func (e *encoder) writeHeader(m *Module) error {
	if err := e.w.WriteBytes(Magic[:]); err != nil {
		return err
	}
	if err := e.w.WriteUint16(m.Version); err != nil {
		return err
	}
	return e.w.WriteUint16(m.Flags)
}

func (e *encoder) writeStringTable() error {
	if err := e.w.WriteVarint(uint64(len(e.strings))); err != nil {
		return err
	}
	for _, s := range e.strings {
		b := []byte(s)
		if err := e.w.WriteVarint(uint64(len(b))); err != nil {
			return err
		}
		if err := e.w.WriteBytes(b); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) stringRef(s string) (uint64, error) {
	idx, ok := e.strIndex[s]
	if !ok {
		return 0, fmt.Errorf("string %q not in string table", s)
	}
	return uint64(idx), nil
}

func (e *encoder) writeStringRef(s string) error {
	ref, err := e.stringRef(s)
	if err != nil {
		return err
	}
	return e.w.WriteVarint(ref)
}

func (e *encoder) writeChunks() error {
	if err := e.w.WriteVarint(uint64(len(e.chunks))); err != nil {
		return err
	}
	for _, ch := range e.chunks {
		if err := e.w.WriteVarint(uint64(ch.MaxStack)); err != nil {
			return err
		}
		if err := e.w.WriteVarint(uint64(len(ch.Code))); err != nil {
			return err
		}
		for _, op := range ch.Code {
			if err := e.w.WriteInt32(op); err != nil {
				return err
			}
		}
		if err := e.w.WriteVarint(uint64(len(ch.SourceMap))); err != nil {
			return err
		}
		for _, sm := range ch.SourceMap {
			if err := e.w.WriteVarint(uint64(sm.StartIP)); err != nil {
				return err
			}
			if err := e.writeStringRef(sm.File); err != nil {
				return err
			}
			if err := e.w.WriteVarint(uint64(sm.Line)); err != nil {
				return err
			}
			if err := e.w.WriteVarint(uint64(sm.Column)); err != nil {
				return err
			}
			if err := e.w.WriteVarint(uint64(sm.EndLine)); err != nil {
				return err
			}
			if err := e.w.WriteVarint(uint64(sm.EndColumn)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *encoder) writeConsts(consts []vm.Value) error {
	if err := e.w.WriteVarint(uint64(len(consts))); err != nil {
		return err
	}
	for _, v := range consts {
		if err := e.writeValue(v); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) writeValue(v vm.Value) error {
	switch val := v.(type) {
	case *vm.Nil:
		return e.w.WriteByte(TagNil)
	case vm.Boolean:
		if bool(val) {
			return e.w.WriteByte(TagTrue)
		}
		return e.w.WriteByte(TagFalse)
	case vm.Int:
		if err := e.w.WriteByte(TagInt); err != nil {
			return err
		}
		return e.w.WriteSvarint(int64(val))
	case vm.Float:
		if err := e.w.WriteByte(TagFloat); err != nil {
			return err
		}
		return e.w.WriteFloat64(float64(val))
	case vm.String:
		if err := e.w.WriteByte(TagString); err != nil {
			return err
		}
		return e.writeStringRef(string(val))
	case vm.Keyword:
		if err := e.w.WriteByte(TagKeyword); err != nil {
			return err
		}
		return e.writeStringRef(string(val))
	case vm.Symbol:
		if err := e.w.WriteByte(TagSymbol); err != nil {
			return err
		}
		return e.writeStringRef(string(val))
	case vm.Char:
		if err := e.w.WriteByte(TagChar); err != nil {
			return err
		}
		return e.w.WriteInt32(int32(val))
	case *vm.BigInt:
		if err := e.w.WriteByte(TagBigInt); err != nil {
			return err
		}
		bi := val.Val()
		sign := byte(0)
		if bi.Sign() < 0 {
			sign = 1
		}
		if err := e.w.WriteByte(sign); err != nil {
			return err
		}
		mag := new(big.Int).Abs(bi).Bytes()
		if err := e.w.WriteVarint(uint64(len(mag))); err != nil {
			return err
		}
		return e.w.WriteBytes(mag)
	case *vm.Void:
		return e.w.WriteByte(TagVoid)
	case *vm.Func:
		if err := e.w.WriteByte(TagFunc); err != nil {
			return err
		}
		// Find chunk index - search through e.chunks by matching the code
		chunkIdx := e.findChunkIndex(val.Chunk())
		if err := e.w.WriteVarint(uint64(chunkIdx)); err != nil {
			return err
		}
		if err := e.w.WriteVarint(uint64(val.Arity())); err != nil {
			return err
		}
		variadic := byte(0)
		if val.IsVariadic() {
			variadic = 1
		}
		if err := e.w.WriteByte(variadic); err != nil {
			return err
		}
		return e.writeStringRef(val.FuncName())
	case *vm.Var:
		if err := e.w.WriteByte(TagVarRef); err != nil {
			return err
		}
		if err := e.writeStringRef(val.NS()); err != nil {
			return err
		}
		return e.writeStringRef(val.VarName())
	case *vm.List:
		if val == vm.EmptyList {
			return e.w.WriteByte(TagEmptyList)
		}
		if err := e.w.WriteByte(TagList); err != nil {
			return err
		}
		return e.writeSeqConsts(val)
	case vm.ArrayVector:
		if err := e.w.WriteByte(TagVector); err != nil {
			return err
		}
		if err := e.w.WriteVarint(uint64(len(val))); err != nil {
			return err
		}
		for _, item := range val {
			if err := e.writeValue(item); err != nil {
				return err
			}
		}
		return nil
	case *vm.PersistentMap:
		if err := e.w.WriteByte(TagMap); err != nil {
			return err
		}
		return e.writeMapConsts(val)
	case *vm.PersistentSet:
		if err := e.w.WriteByte(TagSet); err != nil {
			return err
		}
		return e.writeSetConsts(val)
	case *vm.RecordType:
		if err := e.w.WriteByte(TagRecordType); err != nil {
			return err
		}
		if err := e.writeStringRef(val.TypeName()); err != nil {
			return err
		}
		fields := val.Fields()
		if err := e.w.WriteVarint(uint64(len(fields))); err != nil {
			return err
		}
		for _, f := range fields {
			if err := e.writeStringRef(string(f)); err != nil {
				return err
			}
		}
		return nil
	case *vm.Record:
		if err := e.w.WriteByte(TagRecord); err != nil {
			return err
		}
		// Write the RecordType inline
		rt := val.RecordType()
		if err := e.writeStringRef(rt.TypeName()); err != nil {
			return err
		}
		fields := rt.Fields()
		if err := e.w.WriteVarint(uint64(len(fields))); err != nil {
			return err
		}
		for _, f := range fields {
			if err := e.writeStringRef(string(f)); err != nil {
				return err
			}
		}
		// Write fixed field values
		ff := val.FixedFields()
		for _, fv := range ff {
			if fv == nil {
				if err := e.writeValue(vm.NIL); err != nil {
					return err
				}
			} else {
				if err := e.writeValue(fv); err != nil {
					return err
				}
			}
		}
		// Write extra map
		return e.writeMapConsts(val.Extra())

	case *vm.Regex:
		if err := e.w.WriteByte(TagRegex); err != nil {
			return err
		}
		return e.writeStringRef(val.Pattern())
	case *vm.Atom:
		if err := e.w.WriteByte(TagAtom); err != nil {
			return err
		}
		return e.writeValue(val.Deref())
	default:
		return fmt.Errorf("unsupported value type for serialization: %T", v)
	}
}

func (e *encoder) findChunkIndex(c *vm.CodeChunk) int {
	code := c.Code()
	for i, ch := range e.chunks {
		if len(ch.Code) == len(code) {
			match := true
			for j := range code {
				if ch.Code[j] != code[j] {
					match = false
					break
				}
			}
			if match {
				return i
			}
		}
	}
	return -1
}

func (e *encoder) writeSeqConsts(l *vm.List) error {
	count := l.RawCount()
	if err := e.w.WriteVarint(uint64(count)); err != nil {
		return err
	}
	var s vm.Seq = l
	for s != nil && s != vm.EmptyList {
		if err := e.writeValue(s.First()); err != nil {
			return err
		}
		s = s.Next()
	}
	return nil
}

func (e *encoder) writeMapConsts(m *vm.PersistentMap) error {
	count := m.RawCount()
	if err := e.w.WriteVarint(uint64(count)); err != nil {
		return err
	}
	s := m.Seq()
	for s != nil && s != vm.EmptyList {
		entry := s.First().(vm.ArrayVector)
		if err := e.writeValue(entry[0]); err != nil {
			return err
		}
		if err := e.writeValue(entry[1]); err != nil {
			return err
		}
		s = s.Next()
	}
	return nil
}

func (e *encoder) writeSetConsts(set *vm.PersistentSet) error {
	count := set.RawCount()
	if err := e.w.WriteVarint(uint64(count)); err != nil {
		return err
	}
	s := set.Seq()
	for s != nil && s != vm.EmptyList {
		if err := e.writeValue(s.First()); err != nil {
			return err
		}
		s = s.Next()
	}
	return nil
}

func (e *encoder) writeNSTable(nsTable map[string]int) error {
	if err := e.w.WriteVarint(uint64(len(nsTable))); err != nil {
		return err
	}
	for name, chunkIdx := range nsTable {
		if err := e.writeStringRef(name); err != nil {
			return err
		}
		if err := e.w.WriteVarint(uint64(chunkIdx)); err != nil {
			return err
		}
	}
	return nil
}
