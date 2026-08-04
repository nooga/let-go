/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// registerReplSupport wires the one native primitive clojure.repl/source-fn
// needs but can't have without reader access: given a (file, 0-based line,
// 0-based column) position, read exactly one form starting there and return
// its raw source text. Registered directly into the "clojure.repl" namespace
// object (fetch-or-create via rt.NS, same as postCoreInit's read-string
// registration into clojure.core) so pkg/rt/core/clojure/repl.lg — loaded on
// demand, source-only (see its lgbgen:skip header) — can build source-fn/doc
// on top of it purely in .lg.
//
// This lives in pkg/compiler (not pkg/rt) because it needs NewLispReader;
// pkg/compiler already depends on pkg/rt, so the dependency only runs one way.
func registerReplSupport() {
	replNS := rt.NS("clojure.repl")
	// rt.NS runs before rt.SetNSLoader (postCoreInit fires from the compiler
	// package's init(), well before main wires up the resolver), so this
	// first-touch creates an empty placeholder namespace with no loader
	// attached to mark it "needs loading". Without this, a later
	// (require 'clojure.repl) sees the namespace already registered and
	// never runs clojure/repl.lg's body — leaving only raw-source defined
	// (see the same fix already needed for ir.data, MarkNSNeedsLoad's other
	// caller).
	rt.MarkNSNeedsLoad("clojure.repl")

	rawSourceFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("raw-source: wrong number of arguments %d (expected 3)", len(vs))
		}
		file, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("raw-source: expected String file, got %T", vs[0])
		}
		lineV, ok1 := vs[1].(vm.Int)
		colV, ok2 := vs[2].(vm.Int)
		if !ok1 || !ok2 {
			return vm.NIL, fmt.Errorf("raw-source: expected Int line/column")
		}
		text, ok := vm.SourceRegistry.Text(string(file))
		if !ok {
			return vm.NIL, nil
		}
		snippet, ok := extractFormText(text, int(lineV), int(colV))
		if !ok {
			return vm.NIL, nil
		}
		return vm.String(snippet), nil
	})
	replNS.Def("raw-source", rawSourceFn)
}

// byteOffsetForPos converts a 0-based (line, column) position — both counted
// in RUNES, matching LispReader's own line/column bookkeeping (see next() in
// reader.go) — into a byte offset into src. Returns false if src is shorter
// than the requested position (stale source text for a moved/edited file).
func byteOffsetForPos(src string, line, col int) (int, bool) {
	curLine, curCol := 0, 0
	for i, r := range src {
		if curLine == line && curCol == col {
			return i, true
		}
		if r == '\n' {
			curLine++
			curCol = 0
		} else {
			curCol++
		}
	}
	if curLine == line && curCol == col {
		return len(src), true
	}
	return 0, false
}

// runePrefixByteLen returns the byte length of the first n runes of s.
func runePrefixByteLen(s string, n int) int {
	i := 0
	for count := 0; count < n && i < len(s); count++ {
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
	}
	return i
}

// extractFormText slices exactly one form's source text out of src, starting
// at the given 0-based (line, column). It repositions a fresh LispReader at
// that byte offset and reads one form; the reader's own rune-count (r.pos)
// after Read tells us how much of the slice the form actually consumed, which
// we convert back to a byte length to cut the substring. This reuses the
// reader as-is rather than requiring EndLine/EndColumn tracking (which the
// reader does not populate anywhere today).
func extractFormText(src string, line, col int) (string, bool) {
	off, ok := byteOffsetForPos(src, line, col)
	if !ok {
		return "", false
	}
	sub := src[off:]
	reader := NewLispReader(strings.NewReader(sub), "<source-fn>")
	if _, err := reader.Read(); err != nil {
		return "", false
	}
	byteLen := runePrefixByteLen(sub, reader.pos)
	return sub[:byteLen], true
}
