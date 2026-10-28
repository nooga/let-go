/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// parseRunFlags makes -m a terminal option: everything after its value belongs
// to the program, including flags and a literal --. Other options retain Go's
// flag parsing rules, including flags registered by an embedding main.
func parseRunFlags(fs *flag.FlagSet, args []string) error {
	for i := 0; i < len(args); i++ {
		f, hasValue := lookupRunFlag(fs, args[i])
		if f == nil {
			break
		}
		if f.Name == "m" || f.Name == "main" {
			end := i + 1
			if !hasValue {
				end++
			}
			if end > len(args) {
				break
			}
			terminal := append([]string(nil), args[:end]...)
			terminal = append(terminal, "--")
			terminal = append(terminal, args[end:]...)
			return fs.Parse(terminal)
		}
		if !hasValue {
			if b, ok := f.Value.(interface{ IsBoolFlag() bool }); !ok || !b.IsBoolFlag() {
				i++
			}
		}
	}
	return fs.Parse(args)
}

func parseMainTarget(target string) (string, string, error) {
	ns, name, qualified := strings.Cut(target, "/")
	if !qualified {
		name = "-main"
	}
	if ns == "" || name == "" || strings.Contains(name, "/") {
		return "", "", fmt.Errorf("-m requires a namespace or namespace/fn, got %q", target)
	}
	return ns, name, nil
}

func runNamespaceMain(target string, args []string) error {
	ns, name, err := parseMainTarget(target)
	if err != nil {
		return err
	}
	namespace, err := rt.RequireNS(ns)
	if err != nil {
		return err
	}
	v := namespace.LookupLocal(vm.Symbol(name))
	if v == nil {
		return fmt.Errorf("-m: %s/%s not found", ns, name)
	}
	fn, ok := v.Deref().(vm.Fn)
	if !ok {
		return fmt.Errorf("-m: %s/%s is not callable", ns, name)
	}
	values := make([]vm.Value, len(args))
	for i, arg := range args {
		values[i] = vm.String(arg)
	}
	_, err = vm.NewExecContext().Invoke(fn, values)
	return err
}

func lookupRunFlag(fs *flag.FlagSet, arg string) (*flag.Flag, bool) {
	if arg == "--" || !strings.HasPrefix(arg, "-") || arg == "-" {
		return nil, false
	}
	name, _, hasValue := strings.Cut(strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-"), "=")
	return fs.Lookup(name), hasValue
}

func validateMainMode(fs *flag.FlagSet) (bool, error) {
	mainMode := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "m" || f.Name == "main" {
			mainMode = true
		}
	})
	if !mainMode {
		return false, nil
	}
	if _, _, err := parseMainTarget(mainTarget); err != nil {
		return true, err
	}
	if compileOutput != "" || bundleOutput != "" || wasmOutput != "" {
		return true, fmt.Errorf("-m cannot be combined with -c, -b or -w")
	}
	return true, nil
}

func programStorageSource(nsResolver *resolver.NSResolver, files []string, mainMode bool) string {
	if mainMode {
		ns, _, _ := parseMainTarget(mainTarget)
		if source := nsResolver.SourcePath(ns); source != "" {
			return source
		}
		return strings.ReplaceAll(ns, ".", "/") + ".lg"
	}
	if len(files) > 0 {
		return files[0]
	}
	return ""
}

func evalExpression(context *compiler.Context, expression string) error {
	context.SetSource("EXPR")
	val, err := runForm(context, expression)
	if err != nil {
		return err
	}
	return printResult(val)
}

// runProgram preserves script-then-expression ordering; namespace main runs
// after the expression, and only if initialization succeeded.
func runProgram(context *compiler.Context, files, userArgs []string, mainMode bool) (ranSomething, runFailed bool) {
	if !mainMode && len(files) >= 1 {
		if err := runScript(context, files[0]); err != nil {
			fmt.Print(vm.FormatError(err))
			runFailed = true
		}
		ranSomething = true
	}

	if expr != "" {
		err := evalExpression(context, expr)
		if err != nil {
			fmt.Print(vm.FormatError(err))
			runFailed = true
		}
		ranSomething = true
	}

	if !mainMode {
		return ranSomething, runFailed
	}
	if runFailed {
		return true, true
	}
	if err := runNamespaceMain(mainTarget, userArgs); err != nil {
		fmt.Fprint(os.Stderr, vm.FormatError(err))
		return true, true
	}
	return true, false
}
