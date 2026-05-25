/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// Opcode execution profiler.
//
// When ProfilingEnabled is true, Frame.Run increments per-opcode and
// per-pair counters on every dispatch. The cost when disabled is a
// single branch on a global bool (well-predicted).
//
// Read counters via ProfileSnapshot(); reset via ResetProfile().
// Exposed to Clojure as profile/enable!, profile/disable!,
// profile/snapshot, profile/reset!.

package vm

import "sync/atomic"

// ProfilingEnabled gates the per-opcode increments in Frame.Run.
// Read on every opcode dispatch when true; a single atomic Load.
var ProfilingEnabled atomic.Bool

// Per-opcode execution counts. opcodeCounts[op] is incremented for
// every dispatch of opcode op (extracted via op & 0xff). Indexed by
// the raw opcode byte (0-255).
var opcodeCounts [256]atomic.Uint64

// Per-pair execution counts. pairCounts[prev][curr] is incremented
// whenever opcode `curr` is dispatched immediately after opcode `prev`.
// `prev` starts at 0 (NOOP) for the first opcode of each frame.
var pairCounts [256][256]atomic.Uint64

// RecordOpcode is called from Frame.Run on every dispatch when
// ProfilingEnabled is true. prevOp is the previous opcode byte (or 0
// at frame start).
func RecordOpcode(prevOp, currOp uint8) {
	opcodeCounts[currOp].Add(1)
	pairCounts[prevOp][currOp].Add(1)
}

// ProfileSample is one entry in a profile snapshot.
type ProfileSample struct {
	Opcode uint8
	Name   string
	Count  uint64
}

// PairSample is one entry in a pair-count snapshot.
type PairSample struct {
	Prev, Curr uint8
	PrevName   string
	CurrName   string
	Count      uint64
}

// ProfileSnapshot returns the current per-opcode counts, sorted desc.
// Only entries with non-zero counts are returned.
func ProfileSnapshot() []ProfileSample {
	var out []ProfileSample
	for op := range 256 {
		c := opcodeCounts[op].Load()
		if c == 0 {
			continue
		}
		out = append(out, ProfileSample{
			Opcode: uint8(op),
			Name:   opcodeMnemonicByByte(uint8(op)),
			Count:  c,
		})
	}
	// Sort descending by count.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Count > out[i].Count {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

// PairSnapshot returns the current per-pair counts, sorted desc.
// Only entries with non-zero counts are returned.
func PairSnapshot() []PairSample {
	var out []PairSample
	for p := range 256 {
		for c := range 256 {
			n := pairCounts[p][c].Load()
			if n == 0 {
				continue
			}
			out = append(out, PairSample{
				Prev:     uint8(p),
				Curr:     uint8(c),
				PrevName: opcodeMnemonicByByte(uint8(p)),
				CurrName: opcodeMnemonicByByte(uint8(c)),
				Count:    n,
			})
		}
	}
	// Sort descending by count.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Count > out[i].Count {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

// ResetProfile zeroes all counters.
func ResetProfile() {
	for op := range 256 {
		opcodeCounts[op].Store(0)
		for c := range 256 {
			pairCounts[op][c].Store(0)
		}
	}
}

// opcodeMnemonicByByte returns just the mnemonic for an opcode byte
// (without the sp prefix that OpcodeToString includes).
func opcodeMnemonicByByte(op uint8) string {
	s := OpcodeToString(int32(op))
	// Slice past the "%d/" prefix.
	for i, c := range s {
		if c == '/' {
			// Find first non-space after /.
			for j := i + 1; j < len(s); j++ {
				if s[j] != ' ' {
					return s[j:]
				}
			}
		}
	}
	return s
}
