//go:build !tinygo || (!js && !wasm_unknown)

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package bundle

import "os"

// sameFile reports whether fi1 and fi2 describe the same file. Thin wrapper
// so the TinyGo js/wasm seam can stub the call — TinyGo's os omits SameFile
// under those GOOSes (types_anyos.go is //go:build !js && !wasm_unknown).
func sameFile(fi1, fi2 os.FileInfo) bool {
	return os.SameFile(fi1, fi2)
}
