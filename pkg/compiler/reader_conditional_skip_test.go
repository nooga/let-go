/*
 * Copyright (c) 2025 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

// A non-matching reader-conditional branch must be skipped in its entirety,
// whatever syntax it contains. skipReaderForm used to under-skip several form
// shapes (metadata prefixes, tagged literals, regex, quote prefixes, nested
// conditionals, char literals/comments inside a list), leaving tokens behind
// that desynced the surrounding conditional and swallowed following forms.
//
// Each case reads TWO forms: the conditional's :default value, then the form
// immediately after the conditional. If the gnarly :clj branch is mis-skipped,
// the trailing form is corrupted or never reached.
func TestReaderConditionalSkipsGnarlyBranches(t *testing.T) {
	type twoForms struct{ first, second vm.Value }
	cases := map[string]twoForms{
		"#?(:clj ^long size :default 1) 2":            {vm.Int(1), vm.Int(2)},   // metadata prefix
		"#?(:clj #js [1 2] :default 3) 4":             {vm.Int(3), vm.Int(4)},   // tagged literal
		"#?(:clj #\"a)b\" :default 5) 6":              {vm.Int(5), vm.Int(6)},   // regex containing )
		"#?(:clj 'foo :default 7) 8":                  {vm.Int(7), vm.Int(8)},   // quote prefix
		"#?(:clj #?(:bb 1 :default 2) :default 9) 10": {vm.Int(9), vm.Int(10)},  // nested conditional
		"#?(:clj #'foo :default 11) 12":               {vm.Int(11), vm.Int(12)}, // var-quote
		"#?(:clj #_ x :default 13) 14":                {vm.Int(13), vm.Int(14)}, // discard
		"#?(:clj (a \\) b) :default 15) 16":           {vm.Int(15), vm.Int(16)}, // char literal ) inside list
		"#?(:clj (a ; )\n b) :default 17) 18":         {vm.Int(17), vm.Int(18)}, // comment ) inside list
	}
	for p, e := range cases {
		r := NewLispReader(strings.NewReader(p), "<reader>")
		first, err := r.Read()
		assert.NoErrorf(t, err, "first form of %q", p)
		assert.Equalf(t, e.first, first, "first form of %q", p)
		second, err := r.Read()
		assert.NoErrorf(t, err, "second form of %q", p)
		assert.Equalf(t, e.second, second, "second form of %q", p)
	}
}
