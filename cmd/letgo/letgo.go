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
	message := ` █   ██▀ ▀█▀    ▄▀  ▄▀▄
 █▄▄ █▄▄  █  ▀▀ ▀▄█ ▀▄▀

`
	fmt.Println(message)
}

func repl(ctx *compiler.Context) {
	ctx.SetSource("REPL")
	scanner := bufio.NewScanner(os.Stdin)
	prompt := ctx.CurrentNS().Name() + "=> "
	fmt.Print(prompt)
	for scanner.Scan() {
		in := scanner.Text()
		chunk, err := ctx.Compile(in)
		if err != nil {
			fmt.Println(err)
			continue
		}

		val, err := vm.NewFrame(chunk, nil).Run()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(val.String())
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

func init() {
	flag.BoolVar(&runREPL, "repl", false, "attach REPL after running given files")
}

func initCompiler() *compiler.Context {
	ns := rt.NS("lang")
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

	if len(files) >= 1 {
		for i := range files {
			err := runFile(context, files[i])
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	} else {
		runREPL = true
	}

	if runREPL {
		motd()
		repl(context)
	}
}
