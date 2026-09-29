package bytecode

import (
	"fmt"
	"sort"

	"github.com/nooga/let-go/pkg/vm"
)

// Namespace ranges (FlagNSRanges) let a bundle decoder leave a namespace's
// code undecoded until the namespace is required. The encoder side proves
// which namespaces qualify; the decoder side (deferred.go) consumes the
// section. Both sides share the wire layout below:
//
//	count varint
//	per entry: name string-ref, main chunk index, chunk-range lo, chunk-range hi,
//	           run count, then per run: const lo, const hi
//
// Entries are sorted by name so the bytes are deterministic.

const (
	ownerNone     int32 = -1 // reached by no namespace's code
	ownerMultiple int32 = -2 // reached by more than one namespace's code
)

// forEachFunc calls fn for every function constant inside v, including one
// nested in a collection literal.
func forEachFunc(v vm.Value, fn func(*vm.Func)) {
	switch val := v.(type) {
	case *vm.Func:
		fn(val)
	case *vm.List:
		var s vm.Seq = val
		for s != nil && s != vm.EmptyList {
			forEachFunc(s.First(), fn)
			s = s.Next()
		}
	case vm.DefMetaPairs:
		for _, item := range val {
			forEachFunc(item, fn)
		}
	case vm.ArrayVector:
		for _, item := range val {
			forEachFunc(item, fn)
		}
	case *vm.PersistentMap:
		s := val.Seq()
		for s != nil && s != vm.EmptyList {
			if k, mv, ok := vm.MapEntryKV(s.First()); ok {
				forEachFunc(k, fn)
				forEachFunc(mv, fn)
			}
			s = s.Next()
		}
	case *vm.PersistentSet:
		s := val.Seq()
		for s != nil && s != vm.EmptyList {
			forEachFunc(s.First(), fn)
			s = s.Next()
		}
	case *vm.Record:
		for _, fv := range val.FixedFields() {
			if fv != nil {
				forEachFunc(fv, fn)
			}
		}
		forEachFunc(val.Extra(), fn)
	case *vm.Atom:
		forEachFunc(val.Deref(), fn)
	}
}

// reachability walks every namespace's code from its main chunk — through
// each function constant it loads into that function's chunk — and records,
// per constant and per chunk, which single namespace reaches it, ownerNone,
// or ownerMultiple. names[i] is owner i.
func (b *ModuleBuilder) reachability(names []string) (constOwner, chunkOwner []int32) {
	constOwner = make([]int32, len(b.consts))
	chunkOwner = make([]int32, len(b.chunks))
	for i := range constOwner {
		constOwner[i] = ownerNone
	}
	for i := range chunkOwner {
		chunkOwner[i] = ownerNone
	}
	mark := func(owners []int32, idx int, ns int32) {
		switch owners[idx] {
		case ownerNone:
			owners[idx] = ns
		case ns, ownerMultiple:
		default:
			owners[idx] = ownerMultiple
		}
	}
	for nsIdx, name := range names {
		ns := int32(nsIdx)
		seen := make(map[int]bool)
		stack := []int{b.nsTable[name]}
		for len(stack) > 0 {
			k := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if seen[k] {
				continue
			}
			seen[k] = true
			mark(chunkOwner, k, ns)
			code := b.chunks[k].Code
			for i := 0; i < len(code); i += vm.OpcodeStride(code[i]) {
				if !vm.ReadsConst(code[i]) || i+1 >= len(code) {
					continue
				}
				c := int(code[i+1])
				if c < b.constsBase || c-b.constsBase >= len(b.consts) {
					continue // a parent-pool index, or out of range: not this module's
				}
				c -= b.constsBase
				mark(constOwner, c, ns)
				forEachFunc(b.consts[c], func(f *vm.Func) {
					if idx, ok := b.chunkIndex[f.Chunk()]; ok && !seen[idx] {
						stack = append(stack, idx)
					}
				})
			}
		}
	}
	return constOwner, chunkOwner
}

// nsRanges derives the FlagNSRanges entries. For a namespace with a registered
// const range, the deferrable constants are those in the range that only its
// own code reaches (or that nothing reaches — dead after macro expansion); a
// constant another namespace also reaches stays eager. The namespace
// qualifies when every function constant in the range is deferrable, those
// functions' chunks form one contiguous run that holds no other chunk, no
// function constant outside the deferrable set points into that run, and
// nothing but the namespace itself reaches its main chunk. Anything less
// keeps the namespace eager — the decoder then decodes it like before, so a
// false negative costs startup bytes, never correctness.
func (b *ModuleBuilder) nsRanges() ([]NSRange, error) {
	if len(b.nsConstRange) == 0 || len(b.nsTable) == 0 {
		return nil, nil
	}
	names := make([]string, 0, len(b.nsTable))
	for name := range b.nsTable {
		names = append(names, name)
	}
	sort.Strings(names)
	constOwner, chunkOwner := b.reachability(names)

	mainChunks := make(map[int]bool, len(b.nsTable))
	for _, idx := range b.nsTable {
		mainChunks[idx] = true
	}
	// chunkRefs[k] lists the const indices whose function constant uses chunk k.
	chunkRefs := make(map[int][]int)
	for c, v := range b.consts {
		forEachFunc(v, func(f *vm.Func) {
			if k, ok := b.chunkIndex[f.Chunk()]; ok {
				chunkRefs[k] = append(chunkRefs[k], c)
			}
		})
	}

	var out []NSRange
	for nsIdx, name := range names {
		r, ok := b.nsConstRange[name]
		if !ok {
			continue
		}
		ns := int32(nsIdx)
		main := b.nsTable[name]
		if chunkOwner[main] != ns || len(chunkRefs[main]) > 0 {
			continue
		}
		lo, hi, count, selfContained := -1, -1, 0, true
		var runs []ConstRange
		deferrable := func(c int) bool {
			return c >= r.Lo && c < r.Hi && (constOwner[c] == ns || constOwner[c] == ownerNone)
		}
		for c := r.Lo; c < r.Hi && c < len(b.consts); c++ {
			if !deferrable(c) {
				continue
			}
			if n := len(runs); n > 0 && runs[n-1].Hi == c {
				runs[n-1].Hi = c + 1
			} else {
				runs = append(runs, ConstRange{Lo: c, Hi: c + 1})
			}
		}
		for c := r.Lo; c < r.Hi && c < len(b.consts); c++ {
			f, isFn := b.consts[c].(*vm.Func)
			if !isFn {
				continue
			}
			k, ok := b.chunkIndex[f.Chunk()]
			if !ok || mainChunks[k] || !deferrable(c) ||
				(chunkOwner[k] != ns && chunkOwner[k] != ownerNone) {
				selfContained = false
				break
			}
			for _, ref := range chunkRefs[k] {
				if !deferrable(ref) {
					selfContained = false
				}
			}
			if lo < 0 || k < lo {
				lo = k
			}
			if k > hi {
				hi = k
			}
			count++
		}
		if !selfContained {
			continue
		}
		if count == 0 {
			out = append(out, NSRange{Name: name, MainChunk: main, ConstRuns: runs})
			continue
		}
		if hi-lo+1 != count {
			continue // another chunk sits inside the run
		}
		out = append(out, NSRange{Name: name, MainChunk: main, ChunkLo: lo, ChunkHi: hi + 1, ConstRuns: runs})
	}
	if err := validateNSRanges(out, len(b.chunks), len(b.consts)); err != nil {
		return nil, err
	}
	return out, nil
}

// validateNSRanges checks the structural invariants a decoder relies on:
// in-bounds indices, ascending pairwise-disjoint const runs, and
// pairwise-disjoint chunk runs that never contain a main chunk. A negative
// constCount skips the const upper bound (the decoder learns it only at the
// consts section and checks it there).
func validateNSRanges(ranges []NSRange, chunkCount, constCount int) error {
	mains := make(map[int]string, len(ranges))
	var allRuns []ConstRange
	for _, r := range ranges {
		prev := -1
		for _, run := range r.ConstRuns {
			if run.Lo < 0 || run.Lo >= run.Hi || (constCount >= 0 && run.Hi > constCount) || run.Lo < prev {
				return fmt.Errorf("namespace range %s: const run [%d,%d) invalid for %d consts", r.Name, run.Lo, run.Hi, constCount)
			}
			prev = run.Hi
			allRuns = append(allRuns, run)
		}
		if r.MainChunk < 0 || r.MainChunk >= chunkCount {
			return fmt.Errorf("namespace range %s: main chunk %d out of range (have %d)", r.Name, r.MainChunk, chunkCount)
		}
		if r.ChunkLo < 0 || r.ChunkLo > r.ChunkHi || r.ChunkHi > chunkCount {
			return fmt.Errorf("namespace range %s: chunk range [%d,%d) invalid for %d chunks", r.Name, r.ChunkLo, r.ChunkHi, chunkCount)
		}
		if prev, dup := mains[r.MainChunk]; dup {
			return fmt.Errorf("namespace ranges %s and %s share main chunk %d", prev, r.Name, r.MainChunk)
		}
		mains[r.MainChunk] = r.Name
	}
	sort.Slice(allRuns, func(i, j int) bool { return allRuns[i].Lo < allRuns[j].Lo })
	for i := 1; i < len(allRuns); i++ {
		if allRuns[i].Lo < allRuns[i-1].Hi {
			return fmt.Errorf("namespace ranges: const runs [%d,%d) and [%d,%d) overlap", allRuns[i-1].Lo, allRuns[i-1].Hi, allRuns[i].Lo, allRuns[i].Hi)
		}
	}
	for i, a := range ranges {
		for _, r := range ranges {
			if r.MainChunk >= a.ChunkLo && r.MainChunk < a.ChunkHi {
				return fmt.Errorf("namespace range %s: chunk range [%d,%d) contains main chunk %d of %s", a.Name, a.ChunkLo, a.ChunkHi, r.MainChunk, r.Name)
			}
		}
		for _, c := range ranges[i+1:] {
			if a.ChunkLo < c.ChunkHi && c.ChunkLo < a.ChunkHi && a.ChunkLo != a.ChunkHi && c.ChunkLo != c.ChunkHi {
				return fmt.Errorf("namespace ranges %s [%d,%d) and %s [%d,%d) overlap", a.Name, a.ChunkLo, a.ChunkHi, c.Name, c.ChunkLo, c.ChunkHi)
			}
		}
	}
	return nil
}

func (e *encoder) writeNSRanges(ranges []NSRange) error {
	if err := e.w.WriteVarint(uint64(len(ranges))); err != nil {
		return err
	}
	for _, r := range ranges {
		if err := e.writeStringRef(r.Name); err != nil {
			return err
		}
		for _, n := range []int{r.MainChunk, r.ChunkLo, r.ChunkHi, len(r.ConstRuns)} {
			if err := e.w.WriteVarint(uint64(n)); err != nil {
				return err
			}
		}
		for _, run := range r.ConstRuns {
			if err := e.w.WriteVarint(uint64(run.Lo)); err != nil {
				return err
			}
			if err := e.w.WriteVarint(uint64(run.Hi)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *decoder) readNSRanges() ([]NSRange, error) {
	count, err := d.r.ReadVarint()
	if err != nil {
		return nil, fmt.Errorf("reading namespace range count: %w", err)
	}
	ranges := make([]NSRange, 0, count)
	for i := 0; i < int(count); i++ {
		name, err := d.readStringRef()
		if err != nil {
			return nil, fmt.Errorf("reading namespace range name[%d]: %w", i, err)
		}
		var n [4]int
		for j := range n {
			v, err := d.r.ReadVarint()
			if err != nil {
				return nil, fmt.Errorf("reading namespace range %s: %w", name, err)
			}
			n[j] = int(v)
		}
		runs := make([]ConstRange, n[3])
		for j := range runs {
			lo, err := d.r.ReadVarint()
			if err != nil {
				return nil, fmt.Errorf("reading namespace range %s run[%d]: %w", name, j, err)
			}
			hi, err := d.r.ReadVarint()
			if err != nil {
				return nil, fmt.Errorf("reading namespace range %s run[%d]: %w", name, j, err)
			}
			runs[j] = ConstRange{Lo: int(lo), Hi: int(hi)}
		}
		ranges = append(ranges, NSRange{Name: name, MainChunk: n[0], ChunkLo: n[1], ChunkHi: n[2], ConstRuns: runs})
	}
	return ranges, nil
}
