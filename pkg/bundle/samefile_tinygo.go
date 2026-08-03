//go:build tinygo && (js || wasm_unknown)

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package bundle

import "os"

// TinyGo's js/wasm os has no SameFile. The exclude check only matters for
// host-side `lg -b` resource collection, which is not the TinyGo wasm
// payload path — returning false skips the hardlink-robust identity
// exclude (self-embed of dst under a resource root remains possible only
// if someone bundles under TinyGo wasm, which we don't).
func sameFile(_, _ os.FileInfo) bool {
	return false
}
