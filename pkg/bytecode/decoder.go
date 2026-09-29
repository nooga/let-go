package bytecode

import (
	"bytes"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"math/big"
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
	// A namespace DecodeBundle deferred is listed here but absent from
	// NSChunks until NSChunk materializes it.
	NSOrder []string

	deferred *deferredBundle
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
	d.version = version

	if version == 1 {
		// v1 predates capabilities, so like a capability-less v2 bundle it
		// implicitly carries the pre-removal opcode set; resolve the migration
		// here and remap after decode, keeping the frozen v1 path untouched.
		preCount, preHash := preRemovalSignature()
		if err := d.resolveOpcodeSet(preCount, preHash); err != nil {
			return nil, err
		}
		unit, err := d.decodeToExecUnitV1(parent)
		if err != nil {
			return nil, err
		}
		if d.remapFunc != nil {
			d.remapFunc(d.chunks)
		}
		return unit, nil
	}
	if version == 2 || version == FormatVersion {
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

	if err := d.readLiveChunks(sharedConsts, nil); err != nil {
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
// For opcode-set signature changes, it may record a migration function to
// apply after chunks are decoded. Shared by the module and exec-unit v2
// decode paths.
// resolveOpcodeSet compares a bundle's opcode-set signature against the
// runtime's and, on mismatch, records the registered migration remap for the
// chunk-decode pass — or errors when no migration can bridge the gap.
func (d *decoder) resolveOpcodeSet(bundleCount int, bundleHash uint64) error {
	runtimeCount, runtimeHash := vm.OpcodeSetSignature()
	if bundleCount == runtimeCount && bundleHash == runtimeHash {
		return nil
	}
	d.remapFunc = lookupMigration(bundleCount, bundleHash)
	if d.remapFunc == nil {
		return fmt.Errorf(
			"opcode set mismatch: bundle compiled with %d opcodes (signature %016x), runtime has %d (%016x) — recompile the bundle with a matching lg",
			bundleCount, bundleHash, runtimeCount, runtimeHash)
	}
	return nil
}

func (d *decoder) readCapabilities() error {
	if d.flags&FlagCapabilities == 0 {
		// Pre-capability bundles (v1/v2 before #443) implicitly have the
		// pre-removal signature. Check if migration is needed.
		preCount, preHash := preRemovalSignature()
		return d.resolveOpcodeSet(preCount, preHash)
	}
	caps, err := d.r.ReadUint32()
	if err != nil {
		return fmt.Errorf("reading capability mask: %w", err)
	}
	if caps&^SupportedCapabilities != 0 {
		return UnsupportedCapabilityError(caps)
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
		if err := d.resolveOpcodeSet(int(bundleCount), bundleHash); err != nil {
			return err
		}
	}
	d.moduleCaps = caps
	return nil
}

// beginCompressedBody swaps d.r for a reader over the inflated body when
// FlagCompressed is set. It is called after the plaintext header + capability
// section (so a version/opcode mismatch is rejected before any inflate) and
// before the string table — every body section then reads through the new
// reader transparently.
//
// For a byte-backed reader (NewReaderBytes, the embedded-core path) the whole
// body is inflated into a resident buffer and re-wrapped with NewReaderBytes, so
// the decoder's zero-copy deferred source-map slicing keeps working — off the
// inflated buffer instead of the compressed one. For a streaming reader the
// remaining input is wrapped in a flate reader directly; source maps fall back
// to eager decode, which is already the streaming path's behavior.
func (d *decoder) beginCompressedBody() error {
	declaredSize, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading uncompressed bundle body size: %w", err)
	}
	if declaredSize > maxUncompressedBundleBodySize {
		return fmt.Errorf("declared uncompressed bundle body size %d exceeds limit %d", declaredSize, maxUncompressedBundleBodySize)
	}
	codec, err := d.r.ReadByte()
	if err != nil {
		return fmt.Errorf("reading compression codec: %w", err)
	}
	if codec != compressionFlate {
		return fmt.Errorf("unsupported bundle compression codec %d", codec)
	}
	if d.r.HasBackingData() {
		// data[pos:] is the not-yet-consumed remainder, regardless of bufio
		// readahead: pos counts logically-consumed bytes and data is the full
		// original slice.
		rest := d.r.data[d.r.pos:]
		zr := flate.NewReader(bytes.NewReader(rest))
		inflated, readErr := io.ReadAll(io.LimitReader(zr, int64(declaredSize)+1))
		closeErr := zr.Close()
		if readErr != nil {
			readErr = fmt.Errorf("inflating bundle body: %w", readErr)
		}
		if closeErr != nil {
			closeErr = fmt.Errorf("closing compressed bundle body: %w", closeErr)
		}
		if err := errors.Join(readErr, closeErr); err != nil {
			return err
		}
		if uint64(len(inflated)) != declaredSize {
			return fmt.Errorf("inflated bundle body size %d does not match declared size %d", len(inflated), declaredSize)
		}
		d.r = NewReaderBytes(inflated)
		return nil
	}
	zr := flate.NewReader(d.r.r)
	d.r = NewReader(io.LimitReader(zr, int64(declaredSize)+1))
	d.compressedBodySize = declaredSize
	d.compressedBodyCloser = zr
	return nil
}

// finishCompressedBody verifies that a streaming decode consumed exactly the
// declared body size and probes the DEFLATE reader once more so a missing end
// marker or extra inflated byte is reported. Byte-backed decodes are verified
// eagerly in beginCompressedBody and have no stored closer.
func (d *decoder) finishCompressedBody() error {
	if d.compressedBodyCloser == nil {
		return nil
	}
	var validationErr error
	if got := uint64(d.r.Offset()); got != d.compressedBodySize {
		validationErr = fmt.Errorf("decoded bundle body size %d does not match declared size %d", got, d.compressedBodySize)
	} else if _, err := d.r.ReadByte(); err == nil {
		validationErr = fmt.Errorf("inflated bundle body exceeds declared size %d", d.compressedBodySize)
	} else if !errors.Is(err, io.EOF) {
		validationErr = fmt.Errorf("validating compressed bundle body: %w", err)
	}
	closeErr := d.closeCompressedBody()
	return errors.Join(validationErr, closeErr)
}

func (d *decoder) closeCompressedBody() error {
	if d.compressedBodyCloser == nil {
		return nil
	}
	closer := d.compressedBodyCloser
	d.compressedBodyCloser = nil
	if err := closer.Close(); err != nil {
		return fmt.Errorf("closing compressed bundle body: %w", err)
	}
	return nil
}

func (d *decoder) decodeToExecUnitV2(parent *vm.Consts) (*ExecUnit, error) {
	if err := d.readCapabilities(); err != nil {
		return nil, err
	}
	if d.flags&FlagCompressed != 0 {
		if err := d.beginCompressedBody(); err != nil {
			return nil, err
		}
		defer d.closeCompressedBody()
	}

	strings, err := d.readStringTable()
	if err != nil {
		return nil, err
	}
	d.strings = strings

	var ranges []NSRange
	if d.flags&FlagNSRanges != 0 {
		if ranges, err = d.readNSRanges(); err != nil {
			return nil, err
		}
	}

	var sharedConsts *vm.Consts
	if parent != nil {
		sharedConsts = vm.NewChildConsts(parent)
	} else {
		sharedConsts = vm.NewConsts()
	}

	if err := d.readLiveChunks(sharedConsts, ranges); err != nil {
		return nil, err
	}

	// Apply opcode migration if the bundle's signature doesn't match the runtime.
	if d.remapFunc != nil {
		d.remapFunc(d.decodedChunks())
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
			if ns := d.deferred.lookup(name); ns != nil {
				if ns.MainChunk != idx {
					return nil, fmt.Errorf("namespace range %s names main chunk %d, NS table says %d", name, ns.MainChunk, idx)
				}
			} else {
				unit.NSChunks[name] = d.chunks[idx]
			}
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
			if last.idx != 0 && d.chunks[last.idx] == nil {
				return nil, fmt.Errorf("the unit's main chunk (namespace %s) cannot be deferred", last.name)
			}
			unit.MainChunk = d.chunks[last.idx]
		}
	}
	if d.deferred != nil {
		for _, ns := range d.deferred.byChunk {
			if _, ok := nsTable[ns.Name]; !ok {
				return nil, fmt.Errorf("namespace range %s has no NS table entry", ns.Name)
			}
		}
		d.deferred.body = d.r.data
		d.deferred.remap = d.remapFunc
		d.deferred.consts = sharedConsts
		d.deferred.chunks = d.chunks
		unit.deferred = d.deferred
	}

	if err := d.finishCompressedBody(); err != nil {
		return nil, err
	}
	return unit, nil
}

// lookup is nil-receiver safe: a decode with nothing deferred has no state.
func (b *deferredBundle) lookup(name string) *deferredNS {
	if b == nil {
		return nil
	}
	return b.ns[name]
}

// decodedChunks is d.chunks without the deferred (nil) entries, for a
// migration remap that walks each chunk's code. It is d.chunks itself when
// nothing was deferred.
func (d *decoder) decodedChunks() []*vm.CodeChunk {
	if d.deferred == nil {
		return d.chunks
	}
	out := make([]*vm.CodeChunk, 0, len(d.chunks))
	for _, ch := range d.chunks {
		if ch != nil {
			out = append(out, ch)
		}
	}
	return out
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
	d.version = version
	if version == 1 {
		// Same implicit pre-removal signature handling as the exec-unit v1
		// path above. The raw ChunkData code is re-synced from the remapped
		// live chunks so a re-encode of a decoded v1 module cannot smuggle
		// unmigrated opcodes under a current signature.
		preCount, preHash := preRemovalSignature()
		if err := d.resolveOpcodeSet(preCount, preHash); err != nil {
			return nil, err
		}
		m, err := d.readModuleV1()
		if err != nil {
			return nil, err
		}
		if d.remapFunc != nil {
			d.remapFunc(d.chunks)
			for i, cd := range m.Chunks {
				cd.Code = d.chunks[i].Code()
			}
		}
		return m, nil
	}
	if version == 2 || version == FormatVersion {
		return d.readModuleV2()
	}
	return nil, fmt.Errorf("unsupported LGB version %d", version)
}

type decoder struct {
	r                    *Reader
	resolve              VarResolver
	version              uint16
	flags                uint16
	constsBase           int
	strings              []string
	chunks               []*vm.CodeChunk
	moduleCaps           uint32 // populated when FlagCapabilities is set in v2+
	stats                *DecodeStats
	remapFunc            func([]*vm.CodeChunk) // migration to apply after chunks are decoded, or nil
	compressedBodySize   uint64
	compressedBodyCloser io.Closer
	// deferNS (DecodeBundle) picks the namespaces to leave undecoded;
	// deferred is the state selectDeferred built for them.
	deferNS  func(name string) bool
	deferred *deferredBundle
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
	if d.flags&FlagCompressed != 0 {
		if err := d.beginCompressedBody(); err != nil {
			return nil, err
		}
		defer d.closeCompressedBody()
	}

	strings, err := d.readStringTable()
	if err != nil {
		return nil, err
	}
	d.strings = strings

	var ranges []NSRange
	if d.flags&FlagNSRanges != 0 {
		if ranges, err = d.readNSRanges(); err != nil {
			return nil, err
		}
	}

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

	// Apply opcode migration if the bundle's signature doesn't match the runtime.
	if d.remapFunc != nil {
		d.remapFunc(d.chunks)
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
		Version:    d.version,
		Flags:      d.flags,
		Strings:    strings,
		Chunks:     chunkDatas,
		Consts:     consts,
		ConstsBase: d.constsBase,
		NSTable:    nsTable,
		NSRanges:   ranges,
	}
	if d.flags&FlagCapabilities != 0 {
		m.Capabilities = d.moduleCaps
	}
	if err := d.finishCompressedBody(); err != nil {
		return nil, err
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
	var supportedFlags uint16
	switch version {
	case 1:
		supportedFlags = v1Flags
	case 2:
		supportedFlags = v2Flags
	case FormatVersion:
		supportedFlags = v3Flags
	}
	if supportedFlags != 0 && flags&^supportedFlags != 0 {
		return 0, 0, fmt.Errorf("unsupported LGB flags 0x%04x for version %d (supported: 0x%04x)", flags&^supportedFlags, version, supportedFlags)
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

// readLiveChunks decodes the chunk section into live chunks. With ranges
// (FlagNSRanges) and a DecodeBundle deferral, the records of a deferred
// namespace are skipped instead: their offsets are noted for materialize and
// their slots in d.chunks stay nil.
func (d *decoder) readLiveChunks(sharedConsts *vm.Consts, ranges []NSRange) error {
	count, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading chunk count: %w", err)
	}
	if err := d.selectDeferred(ranges, int(count), -1); err != nil {
		return err
	}
	d.chunks = make([]*vm.CodeChunk, count)
	for i := range d.chunks {
		if ns := d.deferredOwnerOfChunk(i); ns != nil {
			off := d.r.Offset()
			switch i {
			case ns.MainChunk:
				ns.mainOff = off
			case ns.ChunkLo:
				ns.chunksOff = off
			}
			if err := d.skipChunkRecord(); err != nil {
				return fmt.Errorf("skipping deferred chunk %d: %w", i, err)
			}
			continue
		}
		chunk, err := d.readChunkRecord(sharedConsts)
		if err != nil {
			return err
		}
		d.chunks[i] = chunk
	}
	return nil
}

// skipChunkRecord advances past one chunk record without decoding it.
func (d *decoder) skipChunkRecord() error {
	if _, err := d.r.ReadVarint(); err != nil {
		return fmt.Errorf("reading max_stack: %w", err)
	}
	codeLen, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading code_len: %w", err)
	}
	if err := d.r.Skip(int(codeLen) * 4); err != nil {
		return fmt.Errorf("skipping code: %w", err)
	}
	smCount, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading source_map count: %w", err)
	}
	for j := 0; j < int(smCount)*6; j++ {
		if _, err := d.r.ReadVarint(); err != nil {
			return fmt.Errorf("skipping source_map entry: %w", err)
		}
	}
	return nil
}

// readChunkRecord decodes one chunk record: max_stack, code, source map.
func (d *decoder) readChunkRecord(sharedConsts *vm.Consts) (*vm.CodeChunk, error) {
	ms, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading max_stack: %w", err)
	}

	codeLen, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading code_len: %w", err)
	}
	chunk := vm.NewCodeChunkWithCapacity(sharedConsts, int(codeLen))
	for j := 0; j < int(codeLen); j++ {
		op, err := d.r.ReadInt32()
		if err != nil {
			return nil, fmt.Errorf("reading code[%d]: %w", j, err)
		}
		chunk.Append(op)
	}
	chunk.SetMaxStack(int(ms))
	if d.stats != nil {
		d.stats.addChunk(int(codeLen))
	}

	smCount, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading source_map count: %w", err)
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
						return nil, fmt.Errorf("skipping source_map entry: %w", err)
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
	return chunk, nil
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
// allocation used by the generic Module decode path. A deferred chunk's table
// is skipped and its offset noted for materialize.
func (d *decoder) readLocalVarTablesInto(chunks []*vm.CodeChunk) error {
	for i, chunk := range chunks {
		if chunk == nil {
			if ns := d.deferredOwnerOfChunk(i); ns != nil {
				off := d.r.Offset()
				switch i {
				case ns.MainChunk:
					ns.mainLVOff = off
				case ns.ChunkLo:
					ns.lvOff = off
				}
			}
			if err := d.skipLocalVarTable(); err != nil {
				return fmt.Errorf("skipping local var table[%d]: %w", i, err)
			}
			continue
		}
		if err := d.readLocalVarTableInto(chunk, i); err != nil {
			return err
		}
	}
	return nil
}

// skipLocalVarTable advances past one chunk's local-variable table.
func (d *decoder) skipLocalVarTable() error {
	count, err := d.r.ReadVarint()
	if err != nil {
		return err
	}
	for j := 0; j < int(count)*2; j++ {
		if _, err := d.r.ReadVarint(); err != nil {
			return err
		}
	}
	return nil
}

// readLocalVarTableInto reads one chunk's local-variable table (index i, for
// error messages) into chunk.
func (d *decoder) readLocalVarTableInto(chunk *vm.CodeChunk, i int) error {
	count, err := d.r.ReadVarint()
	if err != nil {
		return fmt.Errorf("reading local var count[%d]: %w", i, err)
	}
	if count == 0 {
		return nil
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
		// Reconstruct via vm.NewRegex, not raw regexp.Compile: NewRegex carries
		// the terminal-lookahead compatibility fallback, so a pattern that
		// compiled at read time round-trips through the bundle instead of
		// failing with "invalid or unsupported Perl syntax".
		v, err := vm.NewRegex(pattern)
		if err != nil {
			return nil, fmt.Errorf("recompiling regex %q: %w", pattern, err)
		}
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

func (d *decoder) readDefMetaPairs() (vm.Value, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading def metadata count: %w", err)
	}
	pairs := make(vm.DefMetaPairs, int(count)*2)
	for i := range pairs {
		pairs[i], err = d.readValueV2()
		if err != nil {
			return nil, fmt.Errorf("reading def metadata value[%d]: %w", i, err)
		}
	}
	return pairs, nil
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
	if d.deferred != nil {
		if n := len(d.deferred.runs); n > 0 && d.deferred.runs[n-1].Hi > int(count) {
			return fmt.Errorf("namespace range const run ends at %d, past %d consts", d.deferred.runs[n-1].Hi, count)
		}
	}
	for i := 0; i < int(count); i++ {
		if run := d.deferredRunAt(i); run != nil {
			// A constant only a deferred namespace reaches keeps an empty
			// slot; materialize re-reads the run from its first record.
			if i == run.Lo {
				run.ns.runOffs[run.idx] = d.r.Offset()
			}
			if err := d.skipValueV2(); err != nil {
				return fmt.Errorf("skipping deferred const[%d]: %w", i, err)
			}
			shared.AppendDeferred()
			continue
		}
		v, err := d.readValueV2()
		if err != nil {
			return fmt.Errorf("reading const[%d]: %w", i, err)
		}
		shared.Append(v)
	}
	return nil
}

// skipValueV2 advances past one encoded value without building it: no
// allocation and no var resolution. It mirrors readValueV2Tagged's layout per
// tag; TestSkipValueV2MatchesRead pins the two together over every constant
// of the core bundle.
func (d *decoder) skipValueV2() error {
	tagByte, err := d.r.ReadByte()
	if err != nil {
		return fmt.Errorf("reading tag: %w", err)
	}
	tagID := tagByte & tagIDMask
	tagVer := tagByte >> tagVersionShift
	if tagVer != 0 && (tagID != TagIDMap || tagVer != 1) && isKnownTagID(tagID) {
		return fmt.Errorf("unsupported tag version %d for tag ID 0x%02x", tagVer, tagID)
	}
	varints := func(n int) error {
		for j := 0; j < n; j++ {
			if _, err := d.r.ReadVarint(); err != nil {
				return err
			}
		}
		return nil
	}
	values := func(n uint64) error {
		for j := uint64(0); j < n; j++ {
			if err := d.skipValueV2(); err != nil {
				return err
			}
		}
		return nil
	}
	switch tagID {
	case TagIDNil, TagIDTrue, TagIDFalse, TagIDVoid, TagIDEmptyList:
		return nil
	case TagIDInt:
		_, err := d.r.ReadSvarint()
		return err
	case TagIDFloat:
		return d.r.Skip(8)
	case TagIDString, TagIDKeyword, TagIDSymbol, TagIDUUID, TagIDInstant, TagIDRegex:
		return varints(1)
	case TagIDChar:
		return d.r.Skip(4)
	case TagIDBigInt:
		if _, err := d.r.ReadByte(); err != nil {
			return err
		}
		magLen, err := d.r.ReadVarint()
		if err != nil {
			return err
		}
		return d.r.Skip(int(magLen))
	case TagIDFunc:
		if err := varints(2); err != nil { // chunk index, arity
			return err
		}
		if _, err := d.r.ReadByte(); err != nil { // variadic
			return err
		}
		return varints(1) // name
	case TagIDVarRef:
		return varints(2)
	case TagIDList, TagIDVector, TagIDSet:
		count, err := d.r.ReadVarint()
		if err != nil {
			return err
		}
		return values(count)
	case TagIDMap:
		count, err := d.r.ReadVarint()
		if err != nil {
			return err
		}
		return values(count * 2) // pairs: a map or def metadata
	case TagIDRecordType:
		if err := varints(1); err != nil { // name
			return err
		}
		fieldCount, err := d.r.ReadVarint()
		if err != nil {
			return err
		}
		return varints(int(fieldCount))
	case TagIDRecord:
		if err := varints(1); err != nil { // type name
			return err
		}
		fieldCount, err := d.r.ReadVarint()
		if err != nil {
			return err
		}
		if err := varints(int(fieldCount)); err != nil { // field keywords
			return err
		}
		if err := values(fieldCount); err != nil { // fixed field values
			return err
		}
		extra, err := d.r.ReadVarint() // extra map
		if err != nil {
			return err
		}
		return values(extra * 2)
	case TagIDAtom:
		return d.skipValueV2()
	default:
		return fmt.Errorf("unknown tag ID 0x%02x", tagID)
	}
}

// readFuncRest reads the fields of a function constant that follow its chunk
// index and builds the function over the chunk, which must be decoded.
func (d *decoder) readFuncRest(chunkIdx int) (vm.Value, error) {
	if chunkIdx >= len(d.chunks) {
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
	chunk := d.chunks[chunkIdx]
	if chunk == nil {
		return nil, fmt.Errorf("function %q uses chunk %d of a deferred namespace from outside it", name, chunkIdx)
	}
	fn := vm.MakeFunc(int(arity), variadic != 0, chunk)
	fn.SetName(name)
	return fn, nil
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
	return d.readValueV2Tagged(tagByte)
}

// readValueV2Tagged decodes the value whose tag byte was already read.
func (d *decoder) readValueV2Tagged(tagByte byte) (vm.Value, error) {
	tagID := tagByte & tagIDMask
	tagVer := tagByte >> tagVersionShift
	if d.stats != nil {
		d.stats.addTag(tagID)
	}

	if tagVer != 0 && (tagID != TagIDMap || tagVer != 1) && isKnownTagID(tagID) {
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
		return d.readFuncRest(int(chunkIdx))
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
		if tagVer == 1 {
			return d.readDefMetaPairs()
		}
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
		// Reconstruct via vm.NewRegex, not raw regexp.Compile: NewRegex carries
		// the terminal-lookahead compatibility fallback, so a pattern that
		// compiled at read time round-trips through the bundle instead of
		// failing with "invalid or unsupported Perl syntax".
		v, err := vm.NewRegex(pattern)
		if err != nil {
			return nil, fmt.Errorf("recompiling regex %q: %w", pattern, err)
		}
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
