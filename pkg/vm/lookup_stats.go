package vm

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

type LookupStats struct {
	NamespacedCalls uint64
	LookupCalls     uint64
	NamespacedBySym map[string]uint64
	LookupByNS      map[string]uint64
	LookupByNSSym   map[string]uint64
}

type lookupCount struct {
	key   string
	count uint64
}

func (s LookupStats) Summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[lookup] namespaced calls=%d\n", s.NamespacedCalls)
	fmt.Fprintf(&b, "[lookup] namespace lookup calls=%d\n", s.LookupCalls)
	writeLookupCounts(&b, "namespaced symbol", s.NamespacedBySym)
	writeLookupCounts(&b, "lookup namespace", s.LookupByNS)
	writeLookupCounts(&b, "lookup symbol", s.LookupByNSSym)
	return b.String()
}

func writeLookupCounts(b *strings.Builder, prefix string, m map[string]uint64) {
	counts := make([]lookupCount, 0, len(m))
	for k, v := range m {
		counts = append(counts, lookupCount{key: k, count: v})
	}
	slices.SortFunc(counts, func(a, c lookupCount) int {
		if a.count != c.count {
			if a.count > c.count {
				return -1
			}
			return 1
		}
		return strings.Compare(a.key, c.key)
	})
	for _, c := range counts {
		fmt.Fprintf(b, "[lookup] %s %s count=%d\n", prefix, c.key, c.count)
	}
}

var (
	lookupStatsMu sync.Mutex
	// lookupStatsEnabled is atomic so the noters can bail without taking
	// lookupStatsMu — they sit on Namespace.Lookup and Symbol.Namespaced,
	// where an always-taken global mutex is a cross-goroutine serialization
	// point paid even with stats off (the production case). The mutex still
	// guards the stats maps themselves.
	lookupStatsEnabled atomic.Bool
	lookupStatsGlobal  = LookupStats{
		NamespacedBySym: map[string]uint64{},
		LookupByNS:      map[string]uint64{},
		LookupByNSSym:   map[string]uint64{},
	}
)

func SetLookupStatsEnabled(enabled bool) {
	lookupStatsEnabled.Store(enabled)
}

func ResetLookupStats() {
	lookupStatsMu.Lock()
	lookupStatsGlobal = LookupStats{
		NamespacedBySym: map[string]uint64{},
		LookupByNS:      map[string]uint64{},
		LookupByNSSym:   map[string]uint64{},
	}
	lookupStatsMu.Unlock()
}

func SnapshotLookupStats() LookupStats {
	lookupStatsMu.Lock()
	defer lookupStatsMu.Unlock()
	out := LookupStats{
		NamespacedCalls: lookupStatsGlobal.NamespacedCalls,
		LookupCalls:     lookupStatsGlobal.LookupCalls,
		NamespacedBySym: make(map[string]uint64, len(lookupStatsGlobal.NamespacedBySym)),
		LookupByNS:      make(map[string]uint64, len(lookupStatsGlobal.LookupByNS)),
		LookupByNSSym:   make(map[string]uint64, len(lookupStatsGlobal.LookupByNSSym)),
	}
	for k, v := range lookupStatsGlobal.NamespacedBySym {
		out.NamespacedBySym[k] = v
	}
	for k, v := range lookupStatsGlobal.LookupByNS {
		out.LookupByNS[k] = v
	}
	for k, v := range lookupStatsGlobal.LookupByNSSym {
		out.LookupByNSSym[k] = v
	}
	return out
}

func noteNamespaced(sym string) {
	if !lookupStatsEnabled.Load() {
		return
	}
	lookupStatsMu.Lock()
	// Re-check under the lock so the maps mutate only while enabled is
	// observed true. Note SetLookupStatsEnabled(false) is no longer a
	// recording barrier: a noter already past the fast-path check may
	// record one final sample after Set returns. (The previous
	// always-locked version quiesced synchronously; no current caller
	// relies on that.)
	if lookupStatsEnabled.Load() {
		lookupStatsGlobal.NamespacedCalls++
		lookupStatsGlobal.NamespacedBySym[sym]++
	}
	lookupStatsMu.Unlock()
}

func noteLookup(ns, sym string) {
	if !lookupStatsEnabled.Load() {
		return
	}
	lookupStatsMu.Lock()
	if lookupStatsEnabled.Load() {
		lookupStatsGlobal.LookupCalls++
		lookupStatsGlobal.LookupByNS[ns]++
		lookupStatsGlobal.LookupByNSSym[ns+"::"+sym]++
	}
	lookupStatsMu.Unlock()
}
