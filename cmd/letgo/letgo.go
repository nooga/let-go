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
	"fmt"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"log"
	"os"
)

func main() {
	ns := rt.NS("lang")
	if ns == nil {
		fmt.Println("namespace not found")
		return
	}
	comp := compiler.NewCompiler(ns)

	scanner := bufio.NewScanner(os.Stdin)
	prompt := ns.Name() + "=> "
	fmt.Print(prompt)
	for scanner.Scan() {
		in := scanner.Text()
		chunk, err := comp.Compile(in)
		val, err := vm.NewFrame(chunk, nil).Run()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(val.String())
		fmt.Print(prompt)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
