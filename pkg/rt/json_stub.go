//go:build lg_no_json

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * Stub for builds tagged lg_no_json, which leave out the json and transit
 * namespaces so encoding/json is not linked. Its reflection-based codec is
 * reached from the namespace installers' init(), so the linker keeps all of
 * it even in a program that never reads or writes JSON. js/emit still works:
 * it encodes through internal/jsonenc.
 *
 * Pods lose their json and transit+json payload formats here; EDN still works.
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

var errNoJSON = fmt.Errorf("JSON and transit need encoding/json, which this build omits (lg_no_json)")

func JSONEncodeArgs(args []vm.Value) (string, error)    { return "", errNoJSON }
func JSONDecodeValue(s string) (vm.Value, error)        { return vm.NIL, errNoJSON }
func TransitEncodeArgs(args []vm.Value) (string, error) { return "", errNoJSON }
func TransitDecodeValue(s string) (vm.Value, error)     { return vm.NIL, errNoJSON }
