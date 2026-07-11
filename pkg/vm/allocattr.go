/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

// Byte-level allocation attribution back to .lg call sites.
//
// Go's alloc profile bottoms out at builtins (every (assoc ...) is one
// rt closure), hiding which .lg caller drives the bytes. This tooling
// keeps a shadow stack of interpreter frames (Frame.Run pushes/pops when
// enabled) and, at instrumented allocation sites in the collection code,
// records bytes+count keyed by the top few .lg frames (chunk+ip), which
// resolve to file:line via the chunk source maps at dump time.
//
// Enable with LG_ALLOC_ATTR=1. Zero overhead when disabled beyond a
// predictable branch per instrumented site.
//
// Concurrency: the shadow stack is a single process-global guarded by
// attrMu, so concurrent access is memory-SAFE (no data race), but it is a
// single logical stack, so attribution is only ACCURATE when one goroutine
// interprets .lg at a time. The pmapv builtin drops to 1 worker while
// enabled to preserve that. NOTE: in-VM concurrency already exists beyond
// pmapv — `future*` (lang.go) and core.async `go`/`thread` blocks run .lg
// bodies on goroutines via vm.Goroutines.Go, and are NOT serialized under
// the flag. Under LG_ALLOC_ATTR those interleave their frames onto this one
// stack, garbling attribution for concurrent workloads (still safe, just
// inaccurate). Accurate concurrent attribution needs a goroutine-local
// stack; the tool's primary use — profiling the single-threaded lowering
// pipeline (lgbgen --target=go) plus serialized pmapv — is unaffected.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

var allocAttrEnabled = os.Getenv("LG_ALLOC_ATTR") != ""

// AllocAttrEnabled reports whether .lg allocation attribution is on.
// Used by rt to serialize pmapv workers so the shadow stack is coherent.
func AllocAttrEnabled() bool { return allocAttrEnabled }

type allocKind uint8

const (
	akMapCloneAndSet allocKind = iota
	akMapNodeSeq
	akMapEntries
	akNewArrayVector
	akArrayVectorSeq
	akLazySeq
	nAllocKinds
)

var allocKindNames = [nAllocKinds]string{
	"PersistentMap cloneAndSet (assoc path)",
	"PersistentMap nodeSeq (map iteration)",
	"PersistentMap entries (map->seq)",
	"NewArrayVector",
	"ArrayVector Seq/Next",
	"NewLazySeq",
}

type attrFrameRef struct {
	chunk *CodeChunk
	ip    int
}

const attrDepth = 3

type attrKey struct {
	kind   allocKind
	frames [attrDepth]attrFrameRef
}

type attrStat struct {
	bytes uint64
	count uint64
}

var (
	attrMu    sync.Mutex
	attrStack []*Frame
	attrStats = map[attrKey]*attrStat{}
)

// attrPushFrame / attrPopFrame maintain the shadow stack. Called from
// Frame.Run only when allocAttrEnabled. The mutex is cheap relative to
// interpretation and keeps the odd cross-goroutine handoff coherent.
func attrPushFrame(f *Frame) {
	attrMu.Lock()
	attrStack = append(attrStack, f)
	attrMu.Unlock()
}

func attrPopFrame() {
	attrMu.Lock()
	if n := len(attrStack); n > 0 {
		attrStack[n-1] = nil
		attrStack = attrStack[:n-1]
	}
	attrMu.Unlock()
}

// recordAllocAttr books bytes against the current .lg call stack.
// Callers must guard with `if allocAttrEnabled` so the disabled path
// costs one branch.
func recordAllocAttr(kind allocKind, bytes int) {
	attrMu.Lock()
	var key attrKey
	key.kind = kind
	for i, j := 0, len(attrStack)-1; i < attrDepth && j >= 0; i, j = i+1, j-1 {
		f := attrStack[j]
		key.frames[i] = attrFrameRef{chunk: f.code, ip: f.ip}
	}
	st := attrStats[key]
	if st == nil {
		st = &attrStat{}
		attrStats[key] = st
	}
	st.bytes += uint64(bytes)
	st.count++
	attrMu.Unlock()
}

func (r attrFrameRef) String() string {
	if r.chunk == nil {
		return "<host>"
	}
	si := r.chunk.LookupSource(r.ip)
	if si == nil {
		return "<no-source>"
	}
	if si.Symbol != "" {
		return fmt.Sprintf("%s:%d %s", si.File, si.Line+1, si.Symbol)
	}
	return fmt.Sprintf("%s:%d", si.File, si.Line+1)
}

// DumpAllocAttr writes the attribution report: per-kind totals, then the
// top call-stack keys by bytes.
func DumpAllocAttr(w io.Writer) {
	attrMu.Lock()
	defer attrMu.Unlock()
	if len(attrStats) == 0 {
		fmt.Fprintln(w, "alloc-attr: no samples recorded")
		return
	}

	var kindBytes, kindCount [nAllocKinds]uint64
	// Aggregate by the RESOLVED source locations — distinct ips on the same
	// .lg line would otherwise show as duplicate rows.
	type row struct {
		kind  allocKind
		sites [attrDepth]string
		st    attrStat
	}
	agg := map[string]*row{}
	for k, st := range attrStats {
		kindBytes[k.kind] += st.bytes
		kindCount[k.kind] += st.count
		var sites [attrDepth]string
		for d := 0; d < attrDepth; d++ {
			sites[d] = k.frames[d].String()
		}
		id := fmt.Sprintf("%d|%s|%s|%s", k.kind, sites[0], sites[1], sites[2])
		r := agg[id]
		if r == nil {
			r = &row{kind: k.kind, sites: sites}
			agg[id] = r
		}
		r.st.bytes += st.bytes
		r.st.count += st.count
	}
	rows := make([]*row, 0, len(agg))
	for _, r := range agg {
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].st.bytes > rows[j].st.bytes })

	fmt.Fprintln(w, "=== alloc-attr: bytes by instrumented kind ===")
	for k := allocKind(0); k < nAllocKinds; k++ {
		if kindCount[k] == 0 {
			continue
		}
		fmt.Fprintf(w, "%10.1f MB %12d calls  %s\n",
			float64(kindBytes[k])/(1<<20), kindCount[k], allocKindNames[k])
	}

	fmt.Fprintln(w, "\n=== alloc-attr: top .lg sites by bytes ===")
	limit := 40
	for i, r := range rows {
		if i >= limit {
			break
		}
		fmt.Fprintf(w, "%10.1f MB %12d calls  [%s]\n", float64(r.st.bytes)/(1<<20),
			r.st.count, allocKindNames[r.kind])
		for d := 0; d < attrDepth; d++ {
			fmt.Fprintf(w, "      %d: %s\n", d, r.sites[d])
		}
	}
}
