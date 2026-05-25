/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import "strings"

// OpByKeyword looks up an Op by its keyword name. The match is
// case-insensitive AND dash-insensitive, so "load-arg", "loadarg",
// and "LoadArg" all match opTable entry "LoadArg".
//
// Used by Lisp-side IR construction primitives, where ops are passed
// as kebab-case keywords (the Clojure idiom).
func OpByKeyword(name string) (Op, bool) {
	normalized := strings.ReplaceAll(name, "-", "")
	for i, info := range opTable {
		if strings.EqualFold(info.name, normalized) {
			return Op(i), true
		}
	}
	return OpInvalid, false
}
