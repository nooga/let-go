//go:build !plan9

/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/alimpfard/line"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

var completionTerminators map[byte]bool
var styles map[compiler.TokenKind]line.Style

func init() {
	completionTerminators = map[byte]bool{
		'(':  true,
		')':  true,
		'[':  true,
		']':  true,
		'{':  true,
		'}':  true,
		'"':  true,
		'\\': true,
		'\'': true,
		'@':  true,
		'`':  true,
		'~':  true,
		';':  true,
		'#':  true,
	}
	styles = map[compiler.TokenKind]line.Style{
		compiler.TokenNumber:      {ForegroundColor: line.MakeXtermColor(line.XtermColorMagenta)},
		compiler.TokenPunctuation: {ForegroundColor: line.MakeXtermColor(line.XtermColorYellow)},
		compiler.TokenKeyword:     {ForegroundColor: line.MakeXtermColor(line.XtermColorBlue)},
		compiler.TokenString:      {ForegroundColor: line.MakeXtermColor(line.XtermColorCyan)},
		compiler.TokenSpecial:     {ForegroundColor: line.MakeXtermColor(line.XtermColorUnchanged), Bold: true},
	}
}

func repl(ctx *compiler.Context) {
	interrupted := false
	editor := line.NewEditor()
	prompt := ctx.CurrentNS().Name() + "=> "
	editor.SetInterruptHandler(func() {
		interrupted = true
		editor.Finish()
	})
	editor.SetTabCompletionHandler(func(editor line.Editor) []line.Completion {
		lin := editor.Line()
		prefix := ""
		for i := len(lin) - 1; i >= -1; i-- {
			if (i < 0 || completionTerminators[lin[i]] || unicode.IsSpace(rune(lin[i]))) && i+1 < len(lin) {
				prefix = lin[i+1:]
				break
			}
		}
		cur := ctx.CurrentNS()
		symbols := rt.FuzzyNamespacedSymbolLookup(cur, vm.Symbol(prefix))
		completions := []line.Completion{}
		for _, s := range symbols {
			completions = append(completions, line.Completion{
				Text:                      string(s) + " ",
				InvariantOffset:           uint32(len(prefix)),
				AllowCommitWithoutListing: true,
			})
		}
		return completions
	})
	editor.SetRefreshHandler(func(editor line.Editor) {
		lin := editor.Line()
		reader := compiler.NewLispReaderTokenizing(strings.NewReader(lin), "syntax")
		reader.Read() //nolint:errcheck // We really don't care, just need partial parse
		editor.StripStyles()
		for _, t := range reader.Tokens {
			if t.End == -1 {
				continue
			}
			style, ok := styles[t.Kind]
			if !ok {
				continue
			}
			editor.Stylize(line.Span{Start: uint32(t.Start), End: uint32(t.End), Mode: line.SpanModeByte}, style)
		}
	})
	for !interrupted {
		in, err := editor.GetLine(prompt)
		if err != nil {
			fmt.Println("prompt failed: ", err)
			continue
		}
		if in == "" {
			continue
		}
		editor.AddToHistory(in)
		ctx.SetSource("REPL")
		val, err := runForm(ctx, in)
		if err != nil {
			fmt.Print(vm.FormatError(err))
		} else {
			fmt.Println(val.String())
		}
		prompt = ctx.CurrentNS().Name() + "=> "
	}
}
