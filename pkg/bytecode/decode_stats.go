package bytecode

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

type DecodeStats struct {
	Tags          [64]uint64
	StringEntries uint64
	StringBytes   uint64
	VarRefHits    uint64
	VarRefMisses  uint64
	VarRefCreates uint64
}

type decodeTagCount struct {
	id    byte
	count uint64
}

func (s *DecodeStats) addTag(id byte) {
	if int(id) < len(s.Tags) {
		s.Tags[id]++
	}
}

func (s *DecodeStats) addString(n int) {
	s.StringEntries++
	s.StringBytes += uint64(n)
}

func (s *DecodeStats) addVarRefHit() {
	s.VarRefHits++
}

func (s *DecodeStats) addVarRefMiss(created bool) {
	s.VarRefMisses++
	if created {
		s.VarRefCreates++
	}
}

func (s *DecodeStats) mergeFrom(other *DecodeStats) {
	if other == nil {
		return
	}
	for i, n := range other.Tags {
		s.Tags[i] += n
	}
	s.StringEntries += other.StringEntries
	s.StringBytes += other.StringBytes
	s.VarRefHits += other.VarRefHits
	s.VarRefMisses += other.VarRefMisses
	s.VarRefCreates += other.VarRefCreates
}

func (s DecodeStats) Summary() string {
	var counts []decodeTagCount
	for i, n := range s.Tags {
		if n > 0 {
			counts = append(counts, decodeTagCount{id: byte(i), count: n})
		}
	}
	slices.SortFunc(counts, func(a, b decodeTagCount) int {
		if a.count != b.count {
			if a.count > b.count {
				return -1
			}
			return 1
		}
		return strings.Compare(tagName(a.id), tagName(b.id))
	})
	var b strings.Builder
	fmt.Fprintf(&b, "[decode] strings entries=%d bytes=%d\n", s.StringEntries, s.StringBytes)
	fmt.Fprintf(&b, "[decode] var-ref hits=%d misses=%d creates=%d\n", s.VarRefHits, s.VarRefMisses, s.VarRefCreates)
	for _, c := range counts {
		fmt.Fprintf(&b, "[decode] tag %s count=%d\n", tagName(c.id), c.count)
	}
	return b.String()
}

func tagName(id byte) string {
	switch id {
	case TagIDNil:
		return "nil"
	case TagIDTrue:
		return "true"
	case TagIDFalse:
		return "false"
	case TagIDInt:
		return "int"
	case TagIDFloat:
		return "float"
	case TagIDString:
		return "string"
	case TagIDKeyword:
		return "keyword"
	case TagIDSymbol:
		return "symbol"
	case TagIDChar:
		return "char"
	case TagIDBigInt:
		return "bigint"
	case TagIDVoid:
		return "void"
	case TagIDUUID:
		return "uuid"
	case TagIDInstant:
		return "instant"
	case TagIDFunc:
		return "func"
	case TagIDVarRef:
		return "var-ref"
	case TagIDEmptyList:
		return "empty-list"
	case TagIDList:
		return "list"
	case TagIDVector:
		return "vector"
	case TagIDMap:
		return "map"
	case TagIDSet:
		return "set"
	case TagIDRecordType:
		return "record-type"
	case TagIDRecord:
		return "record"
	case TagIDRegex:
		return "regex"
	case TagIDAtom:
		return "atom"
	default:
		return fmt.Sprintf("0x%02x", id)
	}
}

var (
	decodeStatsMu      sync.Mutex
	decodeStatsEnabled bool
	decodeStatsGlobal  DecodeStats
)

func SetDecodeStatsEnabled(enabled bool) {
	decodeStatsMu.Lock()
	decodeStatsEnabled = enabled
	decodeStatsMu.Unlock()
}

func ResetDecodeStats() {
	decodeStatsMu.Lock()
	decodeStatsGlobal = DecodeStats{}
	decodeStatsMu.Unlock()
}

func SnapshotDecodeStats() DecodeStats {
	decodeStatsMu.Lock()
	defer decodeStatsMu.Unlock()
	return decodeStatsGlobal
}

func NoteDecodeVarRefHit() {
	decodeStatsMu.Lock()
	if decodeStatsEnabled {
		decodeStatsGlobal.addVarRefHit()
	}
	decodeStatsMu.Unlock()
}

func NoteDecodeVarRefMiss(created bool) {
	decodeStatsMu.Lock()
	if decodeStatsEnabled {
		decodeStatsGlobal.addVarRefMiss(created)
	}
	decodeStatsMu.Unlock()
}

func decoderStats() *DecodeStats {
	decodeStatsMu.Lock()
	defer decodeStatsMu.Unlock()
	if !decodeStatsEnabled {
		return nil
	}
	return &DecodeStats{}
}

func recordDecodeStats(stats *DecodeStats) {
	if stats == nil {
		return
	}
	decodeStatsMu.Lock()
	decodeStatsGlobal.mergeFrom(stats)
	decodeStatsMu.Unlock()
}
