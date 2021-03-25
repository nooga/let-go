/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
 * persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"log"
	"os"
)

func motd() {
	message := "\u001B[1m\u001b[37;1mLET-GO\u001B[0m \u001B[36mvery-⍺ much-λ\u001b[0m    \u001b[90m(Ctrl-C to quit)\u001b[0m\n\n"
	fmt.Print(message)
}

func runForm(ctx *compiler.Context, in string) (vm.Value, error) {
	chunk, err := ctx.Compile(in)
	if err != nil {
		return nil, err
	}

	val, err := vm.NewFrame(chunk, nil).Run()
	if err != nil {
		return nil, err
	}
	return val, err
}

func repl(ctx *compiler.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	prompt := ctx.CurrentNS().Name() + "=> "
	fmt.Print(prompt)
	for scanner.Scan() {
		in := scanner.Text()
		ctx.SetSource("REPL")
		val, err := runForm(ctx, in)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val.String())
		}
		prompt = ctx.CurrentNS().Name() + "=> "
		fmt.Print(prompt)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
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

func init() {
	flag.BoolVar(&runREPL, "r", false, "attach REPL after running given files")
	flag.StringVar(&expr, "e", "", "eval given expression")
}

func initCompiler() *compiler.Context {
	ns := rt.NS("user")
	if ns == nil {
		fmt.Println("namespace not found")
		return nil
	}
	return compiler.NewCompiler(ns)
}

func main() {
	flag.Parse()
	files := flag.Args()

	context := initCompiler()

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
