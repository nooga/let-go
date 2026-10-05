//go:build !lg_no_gofragments

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package test

// The .lg suite reads #go{...} fragments, as lg does unless built with
// -tags lg_no_gofragments.
import _ "github.com/nooga/let-go/pkg/gofragments"
