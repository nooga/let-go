/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"fmt"
	"os"
	"strings"
)

// ResolveShell interprets the -w-shell flag: "xterm"/"none" pick the built-in
// shell; any other value is a custom HTML template path, which must exist and
// contain exactly one HostBodyMarker. Returns the custom template path ("" for
// built-in) and whether the built-in xterm shell is selected.
func ResolveShell(wShell string) (customTemplate string, xtermShell bool, err error) {
	switch wShell {
	case "xterm":
		return "", true, nil
	case "none":
		return "", false, nil
	default:
		data, rerr := os.ReadFile(wShell)
		if rerr != nil {
			return "", false, fmt.Errorf("-w-shell %q is not 'xterm'/'none' and can't be read as a template: %w", wShell, rerr)
		}
		if strings.Count(string(data), HostBodyMarker) != 1 {
			return "", false, fmt.Errorf("-w-shell template %q must contain exactly one %s marker", wShell, HostBodyMarker)
		}
		return wShell, false, nil
	}
}
