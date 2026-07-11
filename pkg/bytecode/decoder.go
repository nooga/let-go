package bytecode

import (
	"fmt"
	"io"
	"math/big"
	"regexp"
	"sort"

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
	// NSOrder lists namespace names in chunk index order (load/dependency order).
	NSOrder []string
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
		stats:   decoderStats(),
	}
	return d.decodeExec(parent)
}

// DecodeToExecUnitBytes is like DecodeToExecUnit but decodes from an in-memory
// buffer. The buffer stays resident, so per-chunk source maps are captured
// zero-copy and decoded lazily (on first error/stack-trace lookup) instead of
// eagerly at load — removing the dominant startup heap churn. Prefer this for the
// embedded core bundle, which is already a []byte.
func DecodeToExecUnitBytes(data []byte, resolve VarResolver) (*ExecUnit, error) {
	return DecodeToExecUnitBytesWithParent(data, resolve, nil)
}

// DecodeToExecUnitBytesWithParent is DecodeToExecUnitBytes with an optional
// parent const pool (see DecodeToExecUnitWithParent).
func DecodeToExecUnitBytesWithParent(data []byte, resolve VarResolver, parent *vm.Consts) (*ExecUnit, error) {
	d := &decoder{
		r:       NewReaderBytes(data),
		resolve: resolve,
		stats:   decoderStats(),
	}
	return d.decodeExec(parent)
}

func (d *decoder) decodeExec(parent *vm.Consts) (*ExecUnit, error) {
	defer recordDecodeStats(d.stats)

	version, flags, err := d.readHeader()
	if err != nil {
		return nil, err
	}
	d.flags = flags

	if version == 1 {
		return d.decodeToExecUnitV1(parent)
	}
	if version == 2 {
		return d.decodeToExecUnitV2(parent)
	}
	return nil, fmt.Errorf("unsupported LGB version %d", version)
}

// decodeToExecUnitV1 is the frozen v1 decode path. Do not modify.
func (d *decoder) decodeToExecUnitV1(parent *vm.Consts) (*ExecUnit, error) {
	strings, err := d.readStringTable()
	if err != nil {
		return nil, err
	}
	d.strings = strings

	var sharedConsts *vm.Consts
	if parent != nil {
		sharedConsts = vm.NewChildConsts(parent)
	} else {
		sharedConsts = vm.NewConsts()
	}

	if err := d.readLiveChunks(sharedConsts); err != nil {
		return nil, err
	}

	if err := d.readConstsInto(sharedConsts); err != nil {
		return nil, err
	}

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

	if len(nsTable) > 0 {
		unit.NSChunks = make(map[string]*vm.CodeChunk, len(nsTable))
		type nsEntry struct {
			name string
			idx  int
		}
		entries := make([]nsEntry, 0, len(nsTable))
		for name, idx := range nsTable {
			if idx >= len(d.chunks) {
				return nil, fmt.Errorf("NS table chunk index %d out of range for %q", idx, name)
			}
			unit.NSChunks[name] = d.chunks[idx]
			entries = append(entries, nsEntry{name, idx})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].idx < entries[j].idx })
		unit.NSOrder = make([]string, len(entries))
		for i, e := range entries {
			unit.NSOrder[i] = e.name
		}
		if coreChunk, ok := unit.NSChunks["core"]; ok {
			unit.MainChunk = coreChunk
		} else if len(entries) > 0 {
			last := entries[len(entries)-1]
			unit.MainChunk = d.chunks[last.idx]
		}
	}

	return unit, nil
}

// readCapabilities reads and validates the capability mask (and each set
// capability's payload) that follows the header when FlagCapabilities is set.
// Shared by the module and exec-unit v2 decode paths.
func (d *decoder) readCapabilities() error {
	if d.flags&FlagCapabilities == 0 {
		return nil
	}
	caps, err := d.r.ReadUint32()
	if err != nil {
		return fmt.Errorf("reading capability mask: %w", err)
	}
	const supportedCaps = CapOpcodeSet
	if caps&^supportedCaps != 0 {
		return fmt.Errorf("unsupported capability mask 0x%08x (supported: 0x%08x)", caps, supportedCaps)
	}
	if caps&CapOpcodeSet != 0 {
		bundleCount, err := d.r.ReadVarint()
		if err != nil {
			return fmt.Errorf("reading opcode-set count: %w", err)
		}
		bundleHash, err := d.r.ReadUint64()
		if err != nil {
			return fmt.Errorf("reading opcode-set hash: %w", err)
		}
		count, hash := vm.OpcodeSetSignature()
		if int(bundleCount) != count || bundleHash != hash {
			return fmt.Errorf(
				"opcode set mismatch: bundle compiled with %d opcodes (signature %016x), runtime has %d (%016x) — recompile the bundle with a matching lg",
				bundleCount, bundleHash, count, hash)
		}
	}
	d.moduleCaps = caps
	return nil
}

func (d *decoder) decodeToExecUnitV2(parent *vm.Consts) (*ExecUnit, error) {
	if err := d.readCapabilities(); err != nil {
		return nil, err
	}

	strings, err := d.readStringTable()
	if err != nil {
		return nil, err
	}
	d.strings = strings

	var sharedConsts *vm.Consts
	if parent != nil {
		sharedConsts = vm.NewChildConsts(parent)
	} else {
		sharedConsts = vm.NewConsts()
	}

	if err := d.readLiveChunks(sharedConsts); err != nil {
		return nil, err
	}

	if err := d.readConstsV2Into(sharedConsts); err != nil {
		return nil, err
	}

	nsTable, err := d.readNSTable()
	if err != nil {
		return nil, err
	}
	if d.flags&FlagLocalVars != 0 {
		if err := d.readLocalVarTablesInto(d.chunks); err != nil {
			return nil, err
		}
	}

	if len(d.chunks) == 0 {
		return nil, fmt.Errorf("no chunks in module")
	}

	unit := &ExecUnit{
		Consts:    sharedConsts,
		MainChunk: d.chunks[0],
	}

	if len(nsTable) > 0 {
		unit.NSChunks = make(map[string]*vm.CodeChunk, len(nsTable))
		type nsEntry struct {
			name string
			idx  int
		}
		entries := make([]nsEntry, 0, len(nsTable))
		for name, idx := range nsTable {
			if idx >= len(d.chunks) {
				return nil, fmt.Errorf("NS table chunk index %d out of range for %q", idx, name)
			}
			unit.NSChunks[name] = d.chunks[idx]
			entries = append(entries, nsEntry{name, idx})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].idx < entries[j].idx })
		unit.NSOrder = make([]string, len(entries))
		for i, e := range entries {
			unit.NSOrder[i] = e.name
		}
		if coreChunk, ok := unit.NSChunks["core"]; ok {
			unit.MainChunk = coreChunk
		} else if len(entries) > 0 {
			last := entries[len(entries)-1]
			unit.MainChunk = d.chunks[last.idx]
		}
	}

	return unit, nil
}

// DecodeWithResolver reads a binary module, resolving var references with the given function.
func DecodeWithResolver(r io.Reader, resolve VarResolver) (*Module, error) {
	d := &decoder{
		r:       NewReader(r),
		resolve: resolve,
		stats:   decoderStats(),
	}
	defer recordDecodeStats(d.stats)
	version, flags, err := d.readHeader()
	if err != nil {
		return nil, err
	}
	d.flags = flags
	if version == 1 {
		return d.readModuleV1()
	}
	if version == 2 {
		return d.readModuleV2()
	}
	return nil, fmt.Errorf("unsupported LGB version %d", version)
}

type decoder struct {
	r          *Reader
	resolve    VarResolver
	flags      uint16
	constsBase int
	strings    []string
	chunks     []*vm.CodeChunk
	moduleCaps uint32 // populated when FlagCapabilities is set in v2
	stats      *DecodeStats
}

// readModuleV1 is the frozen v1 decode path. Do not modify.
func (d *decoder) readModuleV1() (*Module, error) {
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
		chunk := vm.NewCodeChunkWithCapacity(sharedConsts, len(cd.Code))
		chunk.Append(cd.Code...)
		chunk.SetMaxStack(cd.MaxStack)
		if len(cd.SourceMap) > 0 {
			chunk.ReserveSourceMap(len(cd.SourceMap))
			for _, e := range cd.SourceMap {
				chunk.AddSourceInfoAt(e.StartIP, vm.SourceInfo{
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
		Version:    1,
		Flags:      d.flags,
		Strings:    strings,
		Chunks:     chunkDatas,
		Consts:     consts,
		ConstsBase: d.constsBase,
		NSTable:    nsTable,
	}, nil
}

func (d *decoder) readModuleV2() (*Module, error) {
	if err := d.readCapabilities(); err != nil {
		return nil, err
	}

	strings, err := d.readStringTable()
	if err != nil {
		return nil, err
	}
	d.strings = strings

	chunkDatas, err := d.readChunks()
	if err != nil {
		return nil, err
	}

	sharedConsts := vm.NewConsts()
	d.chunks = make([]*vm.CodeChunk, len(chunkDatas))
	for i, cd := range chunkDatas {
		chunk := vm.NewCodeChunkWithCapacity(sharedConsts, len(cd.Code))
		chunk.Append(cd.Code...)
		chunk.SetMaxStack(cd.MaxStack)
		if len(cd.SourceMap) > 0 {
			chunk.ReserveSourceMap(len(cd.SourceMap))
			for _, e := range cd.SourceMap {
				chunk.AddSourceInfoAt(e.StartIP, vm.SourceInfo{
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

	consts, err := d.readConstsV2()
	if err != nil {
		return nil, err
	}

	nsTable, err := d.readNSTable()
	if err != nil {
		return nil, err
	}
	if d.flags&FlagLocalVars != 0 {
		tables, err := d.readLocalVarTables(len(chunkDatas))
		if err != nil {
			return nil, err
		}
		for i := range chunkDatas {
			chunkDatas[i].LocalVars = tables[i]
		}
	}

	m := &Module{
		Version:    2,
		Flags:      d.flags,
		Strings:    strings,
		Chunks:     chunkDatas,
		Consts:     consts,
		ConstsBase: d.constsBase,
		NSTable:    nsTable,
	}
	if d.flags&FlagCapabilities != 0 {
		m.Capabilities = d.moduleCaps
	}
	return m, nil
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
		s, err := d.r.ReadString(int(slen))
		if err != nil {
			return nil, fmt.Errorf("reading string data: %w", err)
		}
		strings[i] = s
		if d.stats != nil {
			d.stats.addString(len(s))
		}
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

func (d *decoder) readLiveChunks(sharedConsts *vm.Consts) error {
	count, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading chunk count: %w", err)
	}
	d.chunks = make([]*vm.CodeChunk, count)
	for i := range d.chunks {
		ms, err := d.r.ReadVarint()
		if err != nil {
			return fmt.Errorf("reading max_stack: %w", err)
		}

		codeLen, err := d.r.ReadVarint()
		if err != nil {
			return fmt.Errorf("reading code_len: %w", err)
		}
		chunk := vm.NewCodeChunkWithCapacity(sharedConsts, int(codeLen))
		for j := 0; j < int(codeLen); j++ {
			op, err := d.r.ReadInt32()
			if err != nil {
				return fmt.Errorf("reading code[%d]: %w", j, err)
			}
			chunk.Append(op)
		}
		chunk.SetMaxStack(int(ms))

		smCount, err := d.r.ReadVarint()
		if err != nil {
			return fmt.Errorf("reading source_map count: %w", err)
		}
		if smCount > 0 {
			if d.r.HasBackingData() {
				// Deferred path: capture the source-map section's raw bytes
				// (zero-copy — the backing buffer stays resident) and decode them
				// on first Lookup. Skips per-chunk entries allocation at load.
				// Each entry is 6 varints: startIP, file(string ref), line, col,
				// eline, ecol.
				start := d.r.Offset()
				for j := 0; j < int(smCount); j++ {
					for k := 0; k < 6; k++ {
						if _, err := d.r.ReadVarint(); err != nil {
							return fmt.Errorf("skipping source_map entry: %w", err)
						}
					}
				}
				// Closure-free lazy map: allocates only the SourceMap struct at
				// load (raw is a zero-copy slice of the resident bundle, strings
				// is shared) — decoding is deferred to first Lookup.
				raw := d.r.Slice(start, d.r.Offset())
				chunk.SetSourceMap(vm.NewLazySourceMapRaw(raw, d.strings, int(smCount)))
			} else {
				chunk.ReserveSourceMap(int(smCount))
				for j := 0; j < int(smCount); j++ {
					startIP, err := d.r.ReadVarint()
					if err != nil {
						return err
					}
					file, err := d.readStringRef()
					if err != nil {
						return err
					}
					line, err := d.r.ReadVarint()
					if err != nil {
						return err
					}
					col, err := d.r.ReadVarint()
					if err != nil {
						return err
					}
					eline, err := d.r.ReadVarint()
					if err != nil {
						return err
					}
					ecol, err := d.r.ReadVarint()
					if err != nil {
						return err
					}
					chunk.AddSourceInfoAt(int(startIP), vm.SourceInfo{
						File:      file,
						Line:      int(line),
						Column:    int(col),
						EndLine:   int(eline),
						EndColumn: int(ecol),
					})
				}
			}
		}
		d.chunks[i] = chunk
	}
	return nil
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

func (d *decoder) readConstsInto(shared *vm.Consts) error {
	count, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading const count: %w", err)
	}
	if d.flags&FlagConstsBase != 0 {
		base, err := d.r.ReadVarint()
		if err != nil {
			return fmt.Errorf("reading consts base: %w", err)
		}
		d.constsBase = int(base)
	}
	shared.Reserve(int(count))
	for i := 0; i < int(count); i++ {
		v, err := d.readValue()
		if err != nil {
			return fmt.Errorf("reading const[%d]: %w", i, err)
		}
		shared.Append(v)
	}
	return nil
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

// readLocalVarTables reads the optional per-chunk local-variable debug section
// (written under FlagLocalVars, after the NS table). Returns one slice per chunk
// in index order. Mirrors encoder.writeLocalVarTables.
func (d *decoder) readLocalVarTables(numChunks int) ([][]LocalVarEntry, error) {
	out := make([][]LocalVarEntry, numChunks)
	for i := 0; i < numChunks; i++ {
		count, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading local var count[%d]: %w", i, err)
		}
		if count == 0 {
			continue
		}
		lvs := make([]LocalVarEntry, count)
		for j := range lvs {
			slot, err := d.r.ReadVarint()
			if err != nil {
				return nil, fmt.Errorf("reading local var slot[%d][%d]: %w", i, j, err)
			}
			name, err := d.readStringRef()
			if err != nil {
				return nil, fmt.Errorf("reading local var name[%d][%d]: %w", i, j, err)
			}
			lvs[j] = LocalVarEntry{Slot: int(slot), Name: name}
		}
		out[i] = lvs
	}
	return out, nil
}

// readLocalVarTablesInto reads the optional per-chunk local-variable debug
// section directly into the live chunks, avoiding the temporary [][]LocalVarEntry
// allocation used by the generic Module decode path.
func (d *decoder) readLocalVarTablesInto(chunks []*vm.CodeChunk) error {
	for i, chunk := range chunks {
		count, err := d.r.ReadVarint()
		if err != nil {
			return fmt.Errorf("reading local var count[%d]: %w", i, err)
		}
		if count == 0 {
			continue
		}
		chunk.ReserveLocalVars(int(count))
		for j := 0; j < int(count); j++ {
			slot, err := d.r.ReadVarint()
			if err != nil {
				return fmt.Errorf("reading local var slot[%d][%d]: %w", i, j, err)
			}
			name, err := d.readStringRef()
			if err != nil {
				return fmt.Errorf("reading local var name[%d][%d]: %w", i, j, err)
			}
			chunk.AddLocalVar(int(slot), name)
		}
	}
	return nil
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
	case TagUUID:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		u := vm.ParseUUID(s)
		if u == nil {
			return nil, fmt.Errorf("invalid UUID in bytecode: %q", s)
		}
		return u, nil
	case TagInstant:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		i := vm.ParseInstant(s)
		if i == nil {
			return nil, fmt.Errorf("invalid #inst in bytecode: %q", s)
		}
		return i, nil
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

func (d *decoder) readVectorBatch() (vm.Value, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading vector count: %w", err)
	}
	items := make(vm.ArrayVector, count)
	for i := range items {
		items[i], err = d.readValueV2()
		if err != nil {
			return nil, fmt.Errorf("reading vector item[%d]: %w", i, err)
		}
	}
	return items, nil
}

func (d *decoder) readMapBatch() (vm.Value, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading map count: %w", err)
	}
	kvs := make([]vm.Value, 0, count*2)
	for i := 0; i < int(count); i++ {
		k, err := d.readValueV2()
		if err != nil {
			return nil, fmt.Errorf("reading map key[%d]: %w", i, err)
		}
		v, err := d.readValueV2()
		if err != nil {
			return nil, fmt.Errorf("reading map value[%d]: %w", i, err)
		}
		kvs = append(kvs, k, v)
	}
	return vm.NewPersistentMap(kvs), nil
}

func (d *decoder) readSetBatch() (vm.Value, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading set count: %w", err)
	}
	items := make([]vm.Value, 0, count)
	for i := 0; i < int(count); i++ {
		item, err := d.readValueV2()
		if err != nil {
			return nil, fmt.Errorf("reading set item[%d]: %w", i, err)
		}
		items = append(items, item)
	}
	return vm.NewPersistentSet(items), nil
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

func (d *decoder) readConstsV2() ([]vm.Value, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading const count: %w", err)
	}
	if d.flags&FlagConstsBase != 0 {
		base, err := d.r.ReadVarint()
		if err != nil {
			return nil, fmt.Errorf("reading consts base: %w", err)
		}
		d.constsBase = int(base)
	}
	consts := make([]vm.Value, count)
	for i := range consts {
		v, err := d.readValueV2()
		if err != nil {
			return nil, fmt.Errorf("reading const[%d]: %w", i, err)
		}
		consts[i] = v
	}
	return consts, nil
}

func (d *decoder) readConstsV2Into(shared *vm.Consts) error {
	count, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading const count: %w", err)
	}
	if d.flags&FlagConstsBase != 0 {
		base, err := d.r.ReadVarint()
		if err != nil {
			return fmt.Errorf("reading consts base: %w", err)
		}
		d.constsBase = int(base)
	}
	shared.Reserve(int(count))
	for i := 0; i < int(count); i++ {
		v, err := d.readValueV2()
		if err != nil {
			return fmt.Errorf("reading const[%d]: %w", i, err)
		}
		shared.Append(v)
	}
	return nil
}

func isKnownTagID(id byte) bool {
	switch id {
	case TagIDNil, TagIDTrue, TagIDFalse, TagIDInt, TagIDFloat, TagIDString,
		TagIDKeyword, TagIDSymbol, TagIDChar, TagIDBigInt, TagIDVoid, TagIDUUID,
		TagIDInstant, TagIDFunc, TagIDVarRef, TagIDEmptyList, TagIDList,
		TagIDVector, TagIDMap, TagIDSet, TagIDRecordType, TagIDRecord,
		TagIDRegex, TagIDAtom:
		return true
	}
	return false
}

func (d *decoder) readValueV2() (vm.Value, error) {
	tagByte, err := d.r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("reading tag: %w", err)
	}

	tagID := tagByte & tagIDMask
	tagVer := tagByte >> tagVersionShift
	if d.stats != nil {
		d.stats.addTag(tagID)
	}

	if tagVer != 0 && isKnownTagID(tagID) {
		return nil, fmt.Errorf("unsupported tag version %d for tag ID 0x%02x", tagVer, tagID)
	}

	switch tagID {
	case TagIDNil:
		return vm.NIL, nil
	case TagIDTrue:
		return vm.TRUE, nil
	case TagIDFalse:
		return vm.FALSE, nil
	case TagIDInt:
		v, err := d.r.ReadSvarint()
		if err != nil {
			return nil, err
		}
		return vm.Int(v), nil
	case TagIDFloat:
		v, err := d.r.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return vm.Float(v), nil
	case TagIDString:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		return vm.String(s), nil
	case TagIDKeyword:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		return vm.Keyword(s), nil
	case TagIDSymbol:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		return vm.Symbol(s), nil
	case TagIDChar:
		v, err := d.r.ReadInt32()
		if err != nil {
			return nil, err
		}
		return vm.Char(v), nil
	case TagIDBigInt:
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
	case TagIDVoid:
		return vm.VOID, nil
	case TagIDUUID:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		u := vm.ParseUUID(s)
		if u == nil {
			return nil, fmt.Errorf("invalid UUID in bytecode: %q", s)
		}
		return u, nil
	case TagIDInstant:
		s, err := d.readStringRef()
		if err != nil {
			return nil, err
		}
		i := vm.ParseInstant(s)
		if i == nil {
			return nil, fmt.Errorf("invalid #inst in bytecode: %q", s)
		}
		return i, nil
	case TagIDFunc:
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
	case TagIDVarRef:
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
		return vm.NewVar(nil, ns, name), nil
	case TagIDEmptyList:
		return vm.EmptyList, nil
	case TagIDList:
		count, err := d.r.ReadVarint()
		if err != nil {
			return nil, err
		}
		items := make([]vm.Value, count)
		for i := range items {
			items[i], err = d.readValueV2()
			if err != nil {
				return nil, err
			}
		}
		result, _ := vm.ListType.Box(items)
		return result, nil
	case TagIDVector:
		return d.readVectorBatch()
	case TagIDMap:
		return d.readMapBatch()
	case TagIDSet:
		return d.readSetBatch()
	case TagIDRecordType:
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
	case TagIDRecord:
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
		fixedFields := make([]vm.Value, fieldCount)
		for i := range fixedFields {
			fixedFields[i], err = d.readValueV2()
			if err != nil {
				return nil, err
			}
		}
		extraMap, err := d.readMapBatch()
		if err != nil {
			return nil, err
		}
		data := extraMap.(*vm.PersistentMap)
		for i, kw := range fieldKws {
			if fixedFields[i] != vm.NIL {
				data = data.Assoc(kw, fixedFields[i]).(*vm.PersistentMap)
			}
		}
		return vm.NewRecord(rt, data), nil
	case TagIDRegex:
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
	case TagIDAtom:
		val, err := d.readValueV2()
		if err != nil {
			return nil, err
		}
		return vm.NewAtom(val), nil
	default:
		return nil, fmt.Errorf("unknown tag ID 0x%02x", tagID)
	}
}
