/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// Package gofragments installs the #go{...} raw Go fragment reader: importing
// it adds go -> the native raw reader to the root of
// clojure.core/*data-readers*. A program that reads #go literals imports it,
// usually for its side effect alone:
//
//	import _ "github.com/nooga/let-go/pkg/gofragments"
package gofragments

import "github.com/nooga/let-go/pkg/compiler"

func init() {
	if err := compiler.InstallDataReader("go", compiler.GoRawDataReader()); err != nil {
		panic(err)
	}
}
