/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/alimpfard/line"
	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/nrepl"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func versionString() string {
	if commit != "none" && len(commit) > 7 {
		return fmt.Sprintf("%s (%s)", version, commit[:7])
	}
	return version
}

func motd() {
	banner := "" +
		" \u001b[1m λ\u001b[0m   \u001b[1mlet-go\u001b[0m %s\n" +
		" \u001b[1;36mGO\u001b[0m   \u001b[90mCtrl-C to quit\u001b[0m\n"
	fmt.Printf(banner, versionString())
}

func runForm(ctx *compiler.Context, in string) (vm.Value, error) {
	_, val, err := ctx.CompileMultiple(strings.NewReader(in))
	if err != nil {
		return nil, err
	}
	// if debug {
	// 	val, err = vm.NewDebugFrame(chunk, nil).Run()
	// } else {
	// 	val, err = vm.NewFrame(chunk, nil).Run()
	// }
	// if err != nil {
	// 	return nil, err
	// }
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
			fmt.Print(vm.FormatError(err))
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
	return nil
}

func runLGB(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	resolve := func(nsName, name string) *vm.Var {
		n := rt.DefNSBare(nsName)
		v := n.LookupLocal(vm.Symbol(name))
		if v == nil {
			return n.Def(name, vm.NIL)
		}
		return v
	}
	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(data), resolve)
	if err != nil {
		return fmt.Errorf("decoding %s: %w", filename, err)
	}
	f := vm.NewFrame(unit.MainChunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	return err
}

func compileLG(ctx *compiler.Context, src string, dst string) error {
	ctx.SetSource(src)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	chunk, _, err := ctx.CompileMultiple(f)
	f.Close()
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	return bytecode.EncodeCompilation(out, ctx.Consts(), chunk)
}

var nreplServer *nrepl.NreplServer

func nreplServe(ctx *compiler.Context, port int) error {
	nreplServer = nrepl.NewNreplServer(ctx)
	err := nreplServer.Start(port)
	if err != nil {
		return err
	}
	return nil
}

// Set by goreleaser via ldflags
var (
	version = "dev"
	commit  = "none"
)

var nreplPort int
var runNREPL bool
var runREPL bool
var expr string
var debug bool
var showVersion bool
var compileOutput string

func init() {
	flag.BoolVar(&runREPL, "r", false, "attach REPL after running given files")
	flag.StringVar(&expr, "e", "", "eval given expression")
	flag.BoolVar(&debug, "d", false, "enable VM debug mode")
	flag.BoolVar(&runNREPL, "n", false, "enable nREPL server")
	flag.IntVar(&nreplPort, "p", 2137, "set nREPL port, default is 2137")
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.StringVar(&compileOutput, "c", "", "compile .lg file to .lgb bytecode (specify output path)")

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
	consts := vm.NewConsts()
	ns := rt.NS("user")
	if ns == nil {
		fmt.Println("namespace not found")
		return nil
	}
	if debug {
		return compiler.NewDebugCompiler(consts, ns)
	} else {
		return compiler.NewCompiler(consts, ns)
	}
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("lg %s\n", versionString())
		os.Exit(0)
	}

	files := flag.Args()

	// Ensure all pods are shut down on exit
	defer rt.ShutdownAllPods()

	context := initCompiler(debug)
	nsResolver := resolver.NewNSResolver(context, []string{"."})
	rt.SetNSLoader(nsResolver)

	// Compile mode: compile .lg → .lgb
	if compileOutput != "" {
		if len(files) != 1 {
			fmt.Fprintln(os.Stderr, "error: -c requires exactly one input file")
			os.Exit(1)
		}
		if err := compileLG(context, files[0], compileOutput); err != nil {
			fmt.Fprint(os.Stderr, vm.FormatError(err))
			os.Exit(1)
		}
		return
	}

	ranSomething := false
	if len(files) >= 1 {
		for i := range files {
			if filepath.Ext(files[i]) == ".lgb" {
				// Run precompiled bytecode directly
				if err := runLGB(files[i]); err != nil {
					fmt.Print(vm.FormatError(err))
				}
			} else {
				if err := runFile(context, files[i]); err != nil {
					fmt.Print(vm.FormatError(err))
				}
			}
		}
		ranSomething = true
	}

	if expr != "" {
		context.SetSource("EXPR")
		val, err := runForm(context, expr)
		if err != nil {
			fmt.Print(vm.FormatError(err))
		} else {
			fmt.Println(val)
		}
		ranSomething = true
	}

	if !ranSomething || runREPL {
		motd()
		if runNREPL {
			err := nreplServe(context, nreplPort)
			if err != nil {
				fmt.Println("failed to run nREPL server on port", nreplPort, err)
			}
			fmt.Printf("nREPL server running at tcp://127.0.0.1:%d\n", nreplPort)
		}
		repl(context)
	}

}
