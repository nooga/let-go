/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"fmt"
	"os"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
)

// SplitProgramDebug strips a wasm bundle's program of its debug sections and
// returns the companion together with the path to write it to.
//
// The path is derived from the bundle directory rather than a file inside it:
// every file in the output directory is served, and a debug sidecar is build
// output rather than something to hand to each visitor. A trailing separator on
// outDir would otherwise produce a companion named ".debug" inside the
// directory, which is what the trim guards against.
//
// override, when non-empty, selects the companion path explicitly and must not
// collide with the bundle directory itself.
func SplitProgramDebug(lgb []byte, outDir, override string) (stripped, companion []byte, path string, err error) {
	stripped, companion, err = bytecode.SplitDebug(lgb)
	if err != nil {
		return nil, nil, "", fmt.Errorf("splitting debug sections: %w", err)
	}
	base := strings.TrimRight(outDir, string(os.PathSeparator))
	if base == "" {
		return nil, nil, "", fmt.Errorf("cannot derive a debug companion path from output directory %q", outDir)
	}
	path = override
	if path == "" {
		path = base + bytecode.DebugCompanionSuffix
	}
	if path == base {
		return nil, nil, "", fmt.Errorf("debug output must differ from bundle output %s", base)
	}
	return stripped, companion, path, nil
}
