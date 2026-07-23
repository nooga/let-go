/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// lgbstat reports where a .lgb's bytes go. It decodes to the structural
// Module, tallies contents (chunks, consts by type, string table), attributes
// chunks to source files via their first source-map entry, and measures the
// debug sections' weight by ablation: re-encode with a section stripped and
// diff against the as-is re-encode. Complements scripts/lgbdump.lg, which
// dumps instructions; this accounts bytes.
//
// Usage:
//
//	lgbstat <bundle.lgb>
package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

func encodedLen(m *bytecode.Module) int {
	var buf bytes.Buffer
	if err := bytecode.Encode(&buf, m); err != nil {
		fmt.Fprintln(os.Stderr, "re-encode failed:", err)
		os.Exit(1)
	}
	return buf.Len()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: lgbstat <bundle.lgb>")
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	m, err := bytecode.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, "decode failed:", err)
		os.Exit(1)
	}

	codeInts, smEntries, lvEntries := 0, 0, 0
	for _, c := range m.Chunks {
		codeInts += len(c.Code)
		smEntries += len(c.SourceMap)
		lvEntries += len(c.LocalVars)
	}
	typeCount := map[string]int{}
	var strs []string
	for _, v := range m.Consts {
		typeCount[fmt.Sprintf("%T", v)]++
		if s, ok := v.(vm.String); ok {
			strs = append(strs, string(s))
		}
	}
	fmt.Printf("file=%d bytes  chunks=%d  code_int32s=%d  sourcemap_entries=%d  localvar_entries=%d  consts=%d  strings_table=%d\n",
		len(data), len(m.Chunks), codeInts, smEntries, lvEntries, len(m.Consts), len(m.Strings))

	type kv struct {
		k string
		n int
	}
	var kvs []kv
	for k, n := range typeCount {
		kvs = append(kvs, kv{k, n})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].n > kvs[j].n })
	fmt.Println("const types:")
	for _, e := range kvs {
		fmt.Printf("  %6d  %s\n", e.n, e.k)
	}

	sort.Slice(strs, func(i, j int) bool { return len(strs[i]) > len(strs[j]) })
	fmt.Println("largest string consts:")
	for i := 0; i < 8 && i < len(strs); i++ {
		s := strings.ReplaceAll(strs[i], "\n", "\\n")
		if len(s) > 90 {
			s = s[:90] + "..."
		}
		fmt.Printf("  %5d  %q\n", len(strs[i]), s)
	}

	// String-table accounting is direct: the table is interned against consts,
	// so ablation-by-blanking would desync the encoder.
	tableBytes, proseBytes, proseN := 0, 0, 0
	for _, s := range m.Strings {
		tableBytes += len(s)
		if len(s) >= 40 && strings.Contains(s, " ") {
			proseBytes += len(s)
			proseN++
		}
	}
	fmt.Printf("string table: %d entries, %d bytes; prose-like (docstrings): %d entries, %d bytes (%.1f%% of file)\n",
		len(m.Strings), tableBytes, proseN, proseBytes, float64(proseBytes)/float64(len(data))*100)

	base := encodedLen(m)
	fmt.Printf("re-encoded as-is: %d bytes (file %d)\n", base, len(data))

	saveSM := make([][]bytecode.SourceEntry, len(m.Chunks))
	saveLV := make([][]bytecode.LocalVarEntry, len(m.Chunks))
	for i, c := range m.Chunks {
		saveSM[i], saveLV[i] = c.SourceMap, c.LocalVars
	}
	for i := range m.Chunks {
		m.Chunks[i].SourceMap = nil
	}
	noSM := encodedLen(m)
	for i := range m.Chunks {
		m.Chunks[i].SourceMap = saveSM[i]
	}
	for i := range m.Chunks {
		m.Chunks[i].LocalVars = nil
	}
	// Clear the flag too, or the encoder still writes a zero-count varint per
	// chunk and the ablation under-reports the section by chunks bytes.
	savedFlags := m.Flags
	m.Flags &^= bytecode.FlagLocalVars
	noLV := encodedLen(m)
	m.Flags = savedFlags
	for i := range m.Chunks {
		m.Chunks[i].LocalVars = saveLV[i]
	}
	fmt.Printf("ablation: sourcemaps=%d bytes (%.1f%%)  localvars=%d bytes (%.1f%%)\n",
		base-noSM, float64(base-noSM)/float64(base)*100,
		base-noLV, float64(base-noLV)/float64(base)*100)

	type acc struct{ chunks, ints, sm int }
	byFile := map[string]*acc{}
	for _, c := range m.Chunks {
		file := "<no-sourcemap>"
		if len(c.SourceMap) > 0 {
			file = c.SourceMap[0].File
		}
		a := byFile[file]
		if a == nil {
			a = &acc{}
			byFile[file] = a
		}
		a.chunks++
		a.ints += len(c.Code)
		a.sm += len(c.SourceMap)
	}
	var files []string
	for f := range byFile {
		files = append(files, f)
	}
	sort.Slice(files, func(i, j int) bool { return byFile[files[i]].ints > byFile[files[j]].ints })
	fmt.Println("per-source-file (chunks / code int32s / sourcemap entries):")
	for _, f := range files {
		a := byFile[f]
		fmt.Printf("  %5d chunks %7d ints %6d sm  %s\n", a.chunks, a.ints, a.sm, f)
	}

	var nsNames []string
	for n := range m.NSTable {
		nsNames = append(nsNames, n)
	}
	sort.Strings(nsNames)
	fmt.Println("NS table:")
	fmt.Println(" ", strings.Join(nsNames, " "))
}
