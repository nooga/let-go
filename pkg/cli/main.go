/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// parseRunFlags makes -m a terminal option: everything after its value belongs
// to the program, including flags and a literal --. Other options retain Go's
// flag parsing rules, including flags registered by an embedding main.
func parseRunFlags(fs *flag.FlagSet, args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" || !strings.HasPrefix(arg, "-") || arg == "-" {
			break
		}
		name, _, hasValue := strings.Cut(strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-"), "=")
		f := fs.Lookup(name)
		if f == nil {
			break
		}
		if name == "m" || name == "main" {
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
