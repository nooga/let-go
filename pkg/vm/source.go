/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"sync"
)

// SourceInfo tracks the source location of a form.
//
// Symbol is an optional user-visible name for the value the form binds
// (e.g. a let-binding or parameter name). It is empty for SourceInfos
// attached purely for line/column tracking, and populated by the IR
// builder when a value flows through a `bind-local!` site so downstream
// passes (notably lower-go) can emit human-readable identifiers.
type SourceInfo struct {
	File      string
	Line      int // 0-based
	Column    int // 0-based
	EndLine   int
	EndColumn int
	Symbol    string
}

func (s *SourceInfo) String() string {
	if s == nil {
		return "<unknown>"
	}
	return fmt.Sprintf("%s:%d:%d", s.File, s.Line+1, s.Column+1)
}

// SourceMap maps bytecode IP offsets to source locations.
type SourceMap struct {
	entries []sourceMapEntry
}

type sourceMapEntry struct {
	startIP int
	info    SourceInfo
}

// NewSourceMap creates a new empty SourceMap.
func NewSourceMap() *SourceMap {
	return &SourceMap{}
}

// Add records a source location for the given instruction pointer offset.
func (sm *SourceMap) Add(ip int, info SourceInfo) {
	sm.entries = append(sm.entries, sourceMapEntry{startIP: ip, info: info})
}

// Reserve grows the backing array so that n subsequent Add calls do not
// reallocate. The decoder knows the entry count up front (it reads a counted
// section), so a single sized allocation replaces O(n) append-growth garbage —
// the dominant source of startup heap churn.
func (sm *SourceMap) Reserve(n int) {
	if n > cap(sm.entries) {
		grown := make([]sourceMapEntry, len(sm.entries), n)
		copy(grown, sm.entries)
		sm.entries = grown
	}
}

// SourceMapEntry is the exported version of sourceMapEntry.
type SourceMapEntry struct {
	StartIP int
	Info    SourceInfo
}

// Entries returns the source map entries.
func (sm *SourceMap) Entries() []SourceMapEntry {
	if sm == nil {
		return nil
	}
	out := make([]SourceMapEntry, len(sm.entries))
	for i, e := range sm.entries {
		out[i] = SourceMapEntry{StartIP: e.startIP, Info: e.info}
	}
	return out
}

// Lookup finds the SourceInfo for a given instruction pointer.
// Uses the last entry whose startIP <= ip.
func (sm *SourceMap) Lookup(ip int) *SourceInfo {
	if sm == nil || len(sm.entries) == 0 {
		return nil
	}
	var best *sourceMapEntry
	for i := range sm.entries {
		if sm.entries[i].startIP <= ip {
			best = &sm.entries[i]
		} else {
			break
		}
	}
	if best == nil {
		return nil
	}
	return &best.info
}

// SourceRegistry stores source text for error display.
// Maps file names to their full source text.
var SourceRegistry = &sourceRegistry{sources: map[string]string{}}

type sourceRegistry struct {
	mu      sync.RWMutex
	sources map[string]string
}

// Register stores source text for a given file name.
func (r *sourceRegistry) Register(name string, src string) {
	r.mu.Lock()
	r.sources[name] = src
	r.mu.Unlock()
}

// GetLine returns the line at the given 0-based index for the named file.
func (r *sourceRegistry) GetLine(file string, line int) string {
	r.mu.RLock()
	src, ok := r.sources[file]
	r.mu.RUnlock()
	if !ok {
		return ""
	}
	lines := splitLines(src)
	if line < 0 || line >= len(lines) {
		return ""
	}
	return lines[line]
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start <= len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// FormSource maps form values to their source locations.
// Used by the compiler to attach source info to bytecode.
// Only pointer-based types (like *List) can be tracked; slice/value types
// are not hashable and are silently ignored.
var FormSource = &formSourceMap{m: map[any]*SourceInfo{}}

type formSourceMap struct {
	mu sync.RWMutex
	m  map[any]*SourceInfo
}

// Set associates a source location with a form value.
// Only pointer-identity types can be used as keys; slice types are silently ignored.
func (f *formSourceMap) Set(form Value, info SourceInfo) {
	// Only store for pointer-identity types that can be used as map keys.
	// Slice types like ArrayVector and Map (backed by slices) will panic
	// if used as map keys, so we skip them.
	switch form.(type) {
	case *List, *Cons:
		// These are pointer types, safe to use as map keys
	default:
		return
	}
	f.mu.Lock()
	cp := info
	f.m[form] = &cp
	f.mu.Unlock()
}

// Get retrieves the source location for a form value.
func (f *formSourceMap) Get(form Value) *SourceInfo {
	switch form.(type) {
	case *List, *Cons:
		// pointer types we track
	default:
		return nil
	}
	f.mu.RLock()
	info := f.m[form]
	f.mu.RUnlock()
	return info
}

// Len reports how many form→SourceInfo entries are currently held.
func (f *formSourceMap) Len() int {
	f.mu.RLock()
	n := len(f.m)
	f.mu.RUnlock()
	return n
}

// Reset drops all form→SourceInfo entries. This map is keyed by form
// object identity and never evicts, so a process that compiles many
// distinct forms over time (e.g. BenchmarkClojureTestSuite re-running
// the whole suite per iteration) accumulates an entry — and a live
// reference pinning the *List/*Cons form — for every form ever read.
// Callers that re-compile in a loop should Reset between rounds to keep
// the map (and the forms it pins) from growing unboundedly. Safe to call
// between compiles; source info for the current compile is repopulated
// as forms are read.
func (f *formSourceMap) Reset() {
	f.mu.Lock()
	f.m = map[any]*SourceInfo{}
	f.mu.Unlock()
}
