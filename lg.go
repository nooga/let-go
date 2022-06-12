/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/alimpfard/line"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func motd() {
	message := "\u001B[1m\u001b[37;1mLET-GO\u001B[0m \u001B[36mdev\u001b[0m    \u001b[90m(Ctrl-C to quit)\u001b[0m\n"
	fmt.Print(message)
}

func runForm(ctx *compiler.Context, in string) (vm.Value, error) {
	chunk, err := ctx.Compile(in)
	if err != nil {
		return nil, err
	}
	var val vm.Value
	if debug {
		val, err = vm.NewDebugFrame(chunk, nil).Run()
	} else {
		val, err = vm.NewFrame(chunk, nil).Run()
	}
	if err != nil {
		return nil, err
	}
	return val, err
}

var completionTerminators map[byte]bool
var styles map[compiler.TokenKind]line.Style

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
		reader := compiler.NewLispReader(strings.NewReader(lin), "syntax")
		reader.Read()
		//fmt.Fprintln(os.Stderr, reader.Tokens)
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
	for {
		if interrupted {
			break
		}
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
			fmt.Println(err)
		} else {
			fmt.Println(val.String())
		}
		prompt = ctx.CurrentNS().Name() + "=> "
	}
}

func runFile(ctx *compiler.Context, filename string) error {
	ctx.SetSource(filename)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, _, err = ctx.CompileMultiple(f)
	errc := f.Close()
	if err != nil {
		return err
	}
	if errc != nil {
		return errc
	}
	//_, err = vm.NewFrame(chunk, nil).Run()
	//if err != nil {
	//	return err
	//}
	return nil
}

var runREPL bool
var expr string
var debug bool

func init() {
	flag.BoolVar(&runREPL, "r", false, "attach REPL after running given files")
	flag.StringVar(&expr, "e", "", "eval given expression")
	flag.BoolVar(&debug, "d", false, "enable VM debug mode")
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

func initCompiler(debug bool) *compiler.Context {
	ns := rt.NS("user")
	if ns == nil {
		fmt.Println("namespace not found")
		return nil
	}
	if debug {
		return compiler.NewDebugCompiler(ns)
	} else {
		return compiler.NewCompiler(ns)
	}
}

func main() {
	flag.Parse()
	files := flag.Args()

	context := initCompiler(debug)

	ranSomething := false
	if len(files) >= 1 {
		for i := range files {
			err := runFile(context, files[i])
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		ranSomething = true
	}

	if expr != "" {
		context.SetSource("EXPR")
		val, err := runForm(context, expr)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}
		ranSomething = true
	}

	if !ranSomething || runREPL {
		motd()
		repl(context)
	}

}
