package bytecode

import (
	"fmt"
	"sync"

	"github.com/nooga/let-go/pkg/vm"
)

// DecodeBundle decodes an in-memory bundle like DecodeToExecUnitBytes, except
// that a namespace the bundle marks self-contained (FlagNSRanges) and deferNS
// accepts is left undecoded: its main chunk, the constants only it reaches
// and the chunks of the functions among them are skipped, the constants' pool
// slots are reserved empty (vm.Consts.AppendDeferred), and ExecUnit.NSChunk
// materializes all of it — once, under a lock — when the namespace is first
// required. Everything else (constants shared with another namespace, the
// other namespaces' chunks) decodes exactly as before, so startup pays for
// the code it runs, not for the code a later require may run. A namespace
// without a range entry decodes eagerly whatever deferNS says; the bundle's
// main chunk is never deferred.
func DecodeBundle(data []byte, resolve VarResolver, deferNS func(name string) bool) (*ExecUnit, error) {
	d := &decoder{
		r:       NewReaderBytes(data),
		resolve: resolve,
		stats:   decoderStats(),
		deferNS: deferNS,
	}
	return d.decodeExec(nil)
}

// deferredNS is one namespace DecodeBundle left undecoded: the byte offsets of
// the records the eager pass skipped, and the pool slots of its function
// constants, until materialize fills them in.
type deferredNS struct {
	NSRange
	mainOff   int           // main chunk record, in the chunk section
	chunksOff int           // chunk record ChunkLo; the rest of the run follows it
	mainLVOff int           // main chunk's local-var table, or -1
	lvOff     int           // chunk ChunkLo's local-var table, or -1
	runOffs   []int         // consts-section offset of each ConstRuns entry's first record
	main      *vm.CodeChunk // set last: the namespace is materialized
}

// deferredRun is one ConstRuns entry with its owner, for the const pass's
// ascending cursor.
type deferredRun struct {
	ConstRange
	ns  *deferredNS
	idx int // position in ns.ConstRuns
}

// deferredBundle is the decoder state an ExecUnit keeps so a deferred
// namespace can be decoded later against the same body, string table, pool
// and chunk table. body stays resident (source maps slice it zero-copy).
type deferredBundle struct {
	mu      sync.Mutex
	body    []byte
	strings []string
	flags   uint16
	remap   func([]*vm.CodeChunk)
	consts  *vm.Consts
	chunks  []*vm.CodeChunk
	resolve VarResolver
	ns      map[string]*deferredNS
	byChunk []*deferredNS // sorted by ChunkLo; the runs are disjoint
	runs    []deferredRun // every deferred const run, ascending
	runCur  int           // the const pass's cursor into runs
}

// deferredRunAt returns the deferred run holding const index i, or nil. The
// const pass visits indices in ascending order, so a cursor suffices.
func (d *decoder) deferredRunAt(i int) *deferredRun {
	b := d.deferred
	if b == nil {
		return nil
	}
	for b.runCur < len(b.runs) && i >= b.runs[b.runCur].Hi {
		b.runCur++
	}
	if b.runCur < len(b.runs) && i >= b.runs[b.runCur].Lo {
		return &b.runs[b.runCur]
	}
	return nil
}

// deferredOwnerOfChunk returns the deferred namespace whose main chunk or
// chunk run holds index k, or nil.
func (d *decoder) deferredOwnerOfChunk(k int) *deferredNS {
	if d.deferred == nil {
		return nil
	}
	for _, ns := range d.deferred.byChunk {
		if k == ns.MainChunk || (k >= ns.ChunkLo && k < ns.ChunkHi) {
			return ns
		}
	}
	return nil
}

// selectDeferred turns the bundle's range entries into decoder state for the
// namespaces deferNS accepts. Chunk 0 (the unit's main chunk) is never
// deferred.
func (d *decoder) selectDeferred(ranges []NSRange, chunkCount, constCount int) error {
	if d.deferNS == nil || len(ranges) == 0 || !d.r.HasBackingData() {
		return nil
	}
	if err := validateNSRanges(ranges, chunkCount, constCount); err != nil {
		return err
	}
	var picked []*deferredNS
	for _, r := range ranges {
		if r.MainChunk == 0 || !d.deferNS(r.Name) {
			continue
		}
		picked = append(picked, &deferredNS{NSRange: r, mainLVOff: -1, lvOff: -1, runOffs: make([]int, len(r.ConstRuns))})
	}
	if len(picked) == 0 {
		return nil
	}
	// Insertion-sort by ChunkLo (a bundle has a dozen namespaces).
	for i := 1; i < len(picked); i++ {
		for j := i; j > 0 && picked[j].ChunkLo < picked[j-1].ChunkLo; j-- {
			picked[j], picked[j-1] = picked[j-1], picked[j]
		}
	}
	d.deferred = &deferredBundle{
		strings: d.strings,
		flags:   d.flags,
		resolve: d.resolve,
		ns:      make(map[string]*deferredNS, len(picked)),
		byChunk: picked,
	}
	for _, ns := range picked {
		d.deferred.ns[ns.Name] = ns
		for i, run := range ns.ConstRuns {
			d.deferred.runs = append(d.deferred.runs, deferredRun{ConstRange: run, ns: ns, idx: i})
		}
	}
	runs := d.deferred.runs
	for i := 1; i < len(runs); i++ {
		for j := i; j > 0 && runs[j].Lo < runs[j-1].Lo; j-- {
			runs[j], runs[j-1] = runs[j-1], runs[j]
		}
	}
	return nil
}

// NSChunk returns the main chunk of a bundle namespace, decoding the
// namespace first if DecodeBundle deferred it. It returns (nil, nil) for a
// name the bundle does not define. Safe for concurrent callers: a deferred
// namespace is decoded exactly once.
func (u *ExecUnit) NSChunk(name string) (*vm.CodeChunk, error) {
	if ch := u.NSChunks[name]; ch != nil {
		return ch, nil
	}
	if u.deferred == nil {
		return nil, nil
	}
	return u.deferred.materialize(name)
}

// DeferredNS reports whether name is a namespace DecodeBundle deferred that
// has not been materialized yet.
func (u *ExecUnit) DeferredNS(name string) bool {
	if u.deferred == nil {
		return false
	}
	u.deferred.mu.Lock()
	defer u.deferred.mu.Unlock()
	ns := u.deferred.ns[name]
	return ns != nil && ns.main == nil
}

func (b *deferredBundle) materialize(name string) (*vm.CodeChunk, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	ns := b.ns[name]
	if ns == nil {
		return nil, nil
	}
	if ns.main != nil {
		return ns.main, nil
	}
	d := &decoder{
		strings: b.strings,
		flags:   b.flags,
		chunks:  b.chunks,
		resolve: b.resolve,
		stats:   decoderStats(),
	}
	defer recordDecodeStats(d.stats)

	// The chunk run, then the main chunk. A slot an earlier attempt filled
	// before failing is kept, and skipped below where the tables are per chunk.
	fresh := make(map[*vm.CodeChunk]bool, ns.ChunkHi-ns.ChunkLo+1)
	if ns.ChunkHi > ns.ChunkLo {
		d.r = NewReaderBytes(b.body[ns.chunksOff:])
		for k := ns.ChunkLo; k < ns.ChunkHi; k++ {
			ch, err := d.readChunkRecord(b.consts)
			if err != nil {
				return nil, fmt.Errorf("decoding deferred namespace %s chunk %d: %w", name, k, err)
			}
			if b.chunks[k] == nil {
				b.chunks[k] = ch
				fresh[ch] = true
			}
		}
	}
	d.r = NewReaderBytes(b.body[ns.mainOff:])
	main, err := d.readChunkRecord(b.consts)
	if err != nil {
		return nil, fmt.Errorf("decoding deferred namespace %s main chunk: %w", name, err)
	}
	fresh[main] = true

	if b.flags&FlagLocalVars != 0 {
		if ns.lvOff >= 0 {
			d.r = NewReaderBytes(b.body[ns.lvOff:])
			for k := ns.ChunkLo; k < ns.ChunkHi; k++ {
				if !fresh[b.chunks[k]] {
					if err := d.skipLocalVarTable(); err != nil {
						return nil, fmt.Errorf("decoding deferred namespace %s: %w", name, err)
					}
					continue
				}
				if err := d.readLocalVarTableInto(b.chunks[k], k); err != nil {
					return nil, fmt.Errorf("decoding deferred namespace %s: %w", name, err)
				}
			}
		}
		if ns.mainLVOff >= 0 {
			d.r = NewReaderBytes(b.body[ns.mainLVOff:])
			if err := d.readLocalVarTableInto(main, ns.MainChunk); err != nil {
				return nil, fmt.Errorf("decoding deferred namespace %s: %w", name, err)
			}
		}
	}
	if b.remap != nil {
		chunks := make([]*vm.CodeChunk, 0, len(fresh))
		for ch := range fresh {
			chunks = append(chunks, ch)
		}
		b.remap(chunks)
	}

	// The constants, now that the functions' chunks exist. A run is
	// consecutive records, so one reader per run.
	for i, run := range ns.ConstRuns {
		d.r = NewReaderBytes(b.body[ns.runOffs[i]:])
		for c := run.Lo; c < run.Hi; c++ {
			v, err := d.readValueV2()
			if err != nil {
				return nil, fmt.Errorf("decoding deferred namespace %s const[%d]: %w", name, c, err)
			}
			b.consts.SetDeferred(c, v)
		}
	}
	ns.main = main
	return main, nil
}
