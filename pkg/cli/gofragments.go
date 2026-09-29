//go:build !lg_no_gofragments

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package cli

// lg reads #go{...} raw Go fragments unless built with -tags
// lg_no_gofragments; importing pkg/gofragments installs the reader.
import _ "github.com/nooga/let-go/pkg/gofragments"
