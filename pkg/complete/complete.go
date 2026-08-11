/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

// Package complete supplies REPL symbol completion. It lives outside package
// main so it can be tested — the repo keeps test files out of the repo root.
package complete

import (
	"strings"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Completer implements readline's AutoCompleter interface.
// readline passes line as []rune and pos as a rune index; we stay in runes
// throughout so non-ASCII input and mid-line cursors are handled correctly.
type Completer struct {
	Ctx *compiler.Context
}

func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if pos > len(line) {
		pos = len(line)
	}
	head := line[:pos]

	start := pos
	for start > 0 && !compiler.IsTokenBoundary(head[start-1]) {
		start--
	}
	prefix := string(head[start:])

	symbols, namePrefix := rt.FuzzyNamespacedSymbolLookup(c.Ctx.CurrentNS(), vm.Symbol(prefix))

	// readline INSERTS the candidate at the cursor, so a candidate has to be the
	// part not typed yet — returning whole symbols turns "ns-unm<tab>" into
	// "ns-unmns-unmap". Candidates come back unqualified, because the lookup
	// resolves the ns segment of the prefix and matches on the name, so the typed
	// part to strip is namePrefix, the same name segment the lookup itself
	// matched against. length is that segment's rune count: readline echoes
	// that many runes ahead of each candidate when it lists them.
	for _, s := range symbols {
		newLine = append(newLine, []rune(strings.TrimPrefix(string(s), string(namePrefix))+" "))
	}
	length = len([]rune(namePrefix))
	return
}
