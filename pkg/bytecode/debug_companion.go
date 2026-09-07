/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package bytecode

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/nooga/let-go/pkg/vm"
)

// DebugCompanionSuffix is appended to a stripped artifact unless the caller
// chooses an explicit debug-output path.
const DebugCompanionSuffix = ".debug"

var debugCompanionMagic = [4]byte{'L', 'G', 'D', 0x01}

const debugCompanionVersion uint16 = 1

// HasSplitDebug reports whether data has an LGB header marked for an external
// debug companion. It is intentionally a cheap header probe: full format and
// capability validation remains the decoder's responsibility.
func HasSplitDebug(data []byte) bool {
	return len(data) >= 8 && bytes.Equal(data[:4], Magic[:]) &&
		binary.LittleEndian.Uint16(data[6:8])&FlagDebugSplit != 0
}

type debugCompanion struct {
	chunks []debugChunk
}

type debugChunk struct {
	sourceMap []SourceEntry
	localVars []LocalVarEntry
}

// SplitDebug removes source maps and local-variable tables from an LGB and
// returns both the smaller runtime artifact and a digest-bound companion that
// can restore its debug information. Chunk code and ordering are checked after
// the strip so the companion can never be emitted against a structurally
// different program.
func SplitDebug(data []byte) (stripped, companion []byte, err error) {
	original, err := Decode(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	stripped, err = StripDebug(data)
	if err != nil {
		return nil, nil, err
	}
	slim, err := Decode(bytes.NewReader(stripped))
	if err != nil {
		return nil, nil, fmt.Errorf("decoding stripped LGB: %w", err)
	}
	if len(original.Chunks) != len(slim.Chunks) {
		return nil, nil, fmt.Errorf("strip changed chunk count: %d to %d", len(original.Chunks), len(slim.Chunks))
	}
	for i := range original.Chunks {
		if original.Chunks[i].MaxStack != slim.Chunks[i].MaxStack ||
			!slices.Equal(original.Chunks[i].Code, slim.Chunks[i].Code) {
			return nil, nil, fmt.Errorf("strip changed executable chunk %d", i)
		}
	}
	companion, err = encodeDebugCompanion(stripped, original.Chunks)
	if err != nil {
		return nil, nil, err
	}
	return stripped, companion, nil
}

func encodeDebugCompanion(stripped []byte, chunks []*ChunkData) ([]byte, error) {
	strings := make([]string, 0)
	stringIndex := make(map[string]int)
	intern := func(s string) {
		if _, ok := stringIndex[s]; ok {
			return
		}
		stringIndex[s] = len(strings)
		strings = append(strings, s)
	}
	for _, chunk := range chunks {
		for _, entry := range chunk.SourceMap {
			intern(entry.File)
		}
		for _, local := range chunk.LocalVars {
			intern(local.Name)
		}
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	if err := w.WriteBytes(debugCompanionMagic[:]); err != nil {
		return nil, err
	}
	if err := w.WriteUint16(debugCompanionVersion); err != nil {
		return nil, err
	}
	digest := sha256.Sum256(stripped)
	if err := w.WriteBytes(digest[:]); err != nil {
		return nil, err
	}
	if err := w.WriteVarint(uint64(len(strings))); err != nil {
		return nil, err
	}
	for _, s := range strings {
		if err := w.WriteVarint(uint64(len(s))); err != nil {
			return nil, err
		}
		if err := w.WriteBytes([]byte(s)); err != nil {
			return nil, err
		}
	}
	if err := w.WriteVarint(uint64(len(chunks))); err != nil {
		return nil, err
	}
	for _, chunk := range chunks {
		if err := w.WriteVarint(uint64(len(chunk.SourceMap))); err != nil {
			return nil, err
		}
		for _, entry := range chunk.SourceMap {
			values := []uint64{
				uint64(entry.StartIP), uint64(stringIndex[entry.File]),
				uint64(entry.Line), uint64(entry.Column),
				uint64(entry.EndLine), uint64(entry.EndColumn),
			}
			for _, value := range values {
				if err := w.WriteVarint(value); err != nil {
					return nil, err
				}
			}
		}
		if err := w.WriteVarint(uint64(len(chunk.LocalVars))); err != nil {
			return nil, err
		}
		for _, local := range chunk.LocalVars {
			if err := w.WriteVarint(uint64(local.Slot)); err != nil {
				return nil, err
			}
			if err := w.WriteVarint(uint64(stringIndex[local.Name])); err != nil {
				return nil, err
			}
		}
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeDebugCompanion(data, stripped []byte) (*debugCompanion, error) {
	r := NewReaderBytes(data)
	magic, err := r.ReadBytes(len(debugCompanionMagic))
	if err != nil {
		return nil, fmt.Errorf("reading debug companion magic: %w", err)
	}
	if !bytes.Equal(magic, debugCompanionMagic[:]) {
		return nil, fmt.Errorf("invalid debug companion magic: %x", magic)
	}
	version, err := r.ReadUint16()
	if err != nil {
		return nil, fmt.Errorf("reading debug companion version: %w", err)
	}
	if version != debugCompanionVersion {
		return nil, fmt.Errorf("unsupported debug companion version %d", version)
	}
	wantDigest, err := r.ReadBytes(sha256.Size)
	if err != nil {
		return nil, fmt.Errorf("reading debug companion digest: %w", err)
	}
	gotDigest := sha256.Sum256(stripped)
	if !bytes.Equal(wantDigest, gotDigest[:]) {
		return nil, fmt.Errorf("debug companion does not match stripped LGB (SHA-256 mismatch)")
	}

	stringCount, err := readDebugCount(r, "string")
	if err != nil {
		return nil, err
	}
	strings := make([]string, stringCount)
	for i := range strings {
		length, err := readDebugCount(r, fmt.Sprintf("string length[%d]", i))
		if err != nil {
			return nil, err
		}
		strings[i], err = r.ReadString(length)
		if err != nil {
			return nil, fmt.Errorf("reading debug string[%d]: %w", i, err)
		}
	}
	stringAt := func(raw uint64, field string) (string, error) {
		if raw >= uint64(len(strings)) {
			return "", fmt.Errorf("debug %s string index %d out of range", field, raw)
		}
		return strings[raw], nil
	}

	chunkCount, err := readDebugCount(r, "chunk")
	if err != nil {
		return nil, err
	}
	debug := &debugCompanion{chunks: make([]debugChunk, chunkCount)}
	for i := range debug.chunks {
		sourceCount, err := readDebugCount(r, fmt.Sprintf("source-map[%d]", i))
		if err != nil {
			return nil, err
		}
		debug.chunks[i].sourceMap = make([]SourceEntry, sourceCount)
		for j := range debug.chunks[i].sourceMap {
			values := make([]uint64, 6)
			for k := range values {
				values[k], err = r.ReadVarint()
				if err != nil {
					return nil, fmt.Errorf("reading source-map[%d][%d]: %w", i, j, err)
				}
			}
			file, err := stringAt(values[1], fmt.Sprintf("source-map[%d][%d]", i, j))
			if err != nil {
				return nil, err
			}
			debug.chunks[i].sourceMap[j] = SourceEntry{
				StartIP: int(values[0]), File: file,
				Line: int(values[2]), Column: int(values[3]),
				EndLine: int(values[4]), EndColumn: int(values[5]),
			}
		}
		localCount, err := readDebugCount(r, fmt.Sprintf("local-var[%d]", i))
		if err != nil {
			return nil, err
		}
		debug.chunks[i].localVars = make([]LocalVarEntry, localCount)
		for j := range debug.chunks[i].localVars {
			slot, err := r.ReadVarint()
			if err != nil {
				return nil, fmt.Errorf("reading local-var slot[%d][%d]: %w", i, j, err)
			}
			nameRef, err := r.ReadVarint()
			if err != nil {
				return nil, fmt.Errorf("reading local-var name[%d][%d]: %w", i, j, err)
			}
			name, err := stringAt(nameRef, fmt.Sprintf("local-var[%d][%d]", i, j))
			if err != nil {
				return nil, err
			}
			debug.chunks[i].localVars[j] = LocalVarEntry{Slot: int(slot), Name: name}
		}
	}
	return debug, nil
}

func readDebugCount(r *Reader, field string) (int, error) {
	value, err := r.ReadVarint()
	if err != nil {
		return 0, fmt.Errorf("reading debug %s count: %w", field, err)
	}
	maxInt := uint64(^uint(0) >> 1)
	if value > maxInt {
		return 0, fmt.Errorf("debug %s count %d overflows int", field, value)
	}
	return int(value), nil
}

// DecodeToExecUnitBytesWithDebug decodes an LGB and, when it is marked as
// split-debug, verifies and attaches the supplied companion before execution.
// A companion beside an ordinary unstripped LGB is ignored, which makes stale
// sidecars harmless after recompiling without -strip.
func DecodeToExecUnitBytesWithDebug(data, companion []byte, resolve VarResolver) (*ExecUnit, error) {
	d := &decoder{
		r:       NewReaderBytes(data),
		resolve: resolve,
		stats:   decoderStats(),
	}
	unit, err := d.decodeExec(nil)
	if err != nil {
		return nil, err
	}
	if d.flags&FlagDebugSplit == 0 || len(companion) == 0 {
		return unit, nil
	}
	debug, err := decodeDebugCompanion(companion, data)
	if err != nil {
		return nil, err
	}
	if len(debug.chunks) != len(d.chunks) {
		return nil, fmt.Errorf("debug companion chunk count %d does not match LGB chunk count %d", len(debug.chunks), len(d.chunks))
	}
	for i, info := range debug.chunks {
		chunk := d.chunks[i]
		if len(info.sourceMap) > 0 {
			chunk.ReserveSourceMap(len(info.sourceMap))
			for _, entry := range info.sourceMap {
				chunk.AddSourceInfoAt(entry.StartIP, vm.SourceInfo{
					File: entry.File, Line: entry.Line, Column: entry.Column,
					EndLine: entry.EndLine, EndColumn: entry.EndColumn,
				})
			}
		}
		if len(info.localVars) > 0 {
			chunk.ReserveLocalVars(len(info.localVars))
			for _, local := range info.localVars {
				chunk.AddLocalVar(local.Slot, local.Name)
			}
		}
	}
	return unit, nil
}
