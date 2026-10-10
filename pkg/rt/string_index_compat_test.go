/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
)

// oldIndexOf and oldLastIndexOf are string/index-of and string/last-index-of
// as they were before the byte search: both sides converted to []rune, so
// each invalid UTF-8 byte of the string was a U+FFFD rune. from is nil when
// the call has no from argument.
func oldIndexOf(s, needle string, from *int) vm.Value {
	rs := []rune(s)
	start := 0
	if from != nil {
		start = *from
	}
	if start < 0 || start > len(rs) {
		return vm.NIL
	}
	if idx := runeIndex(rs[start:], []rune(needle)); idx != -1 {
		return vm.MakeInt(idx + start)
	}
	return vm.NIL
}

func oldLastIndexOf(s, needle string, from *int) (v vm.Value, panicked bool) {
	defer func() {
		if recover() != nil { // from < -1 sliced out of range
			v, panicked = vm.NIL, true
		}
	}()
	rs := []rune(s)
	end := len(rs)
	if from != nil {
		end = *from + 1
	}
	if end > len(rs) {
		end = len(rs)
	}
	if idx := runeLastIndex(rs[:end], []rune(needle)); idx != -1 {
		return vm.MakeInt(idx), false
	}
	return vm.NIL, false
}

func TestStringIndexOfMatchesRuneSemantics(t *testing.T) {
	corpus := []string{
		"", "a", "abcab", "héllo", "日本語日", "a\xffb", "\xff", "ab\xe2\x82cd", "€\x82€",
		"a�b\xff", "\xc3", "x\xe2\x82\xacy\xf0\x9f", strings.Repeat("é—", 5) + "\xfe",
	}
	needles := []string{"", "a", "b", "ab", "é", "€", "�", "\xff", "\x82", "\xe2\x82", "日"}
	indexOf, lastIndexOf := stringFn(t, "index-of"), stringFn(t, "last-index-of")
	call := func(fn vm.Fn, s, needle string, from *int) vm.Value {
		args := []vm.Value{vm.String(s), vm.String(needle)}
		if from != nil {
			args = append(args, vm.MakeInt(*from))
		}
		v, err := fn.Invoke(args)
		if err != nil {
			t.Fatalf("%q %q %v: %v", s, needle, from, err)
		}
		return v
	}
	show := func(from *int) string {
		if from == nil {
			return "-"
		}
		return fmt.Sprint(*from)
	}
	for _, s := range corpus {
		n := utf8.RuneCountInString(s)
		froms := []*int{nil}
		for i := -3; i <= n+2; i++ {
			froms = append(froms, &i)
		}
		for _, nd := range needles {
			for _, from := range froms {
				if got, want := call(indexOf, s, nd, from), oldIndexOf(s, nd, from); fmt.Sprint(got) != fmt.Sprint(want) {
					t.Errorf("index-of %q %q %s = %v, want %v", s, nd, show(from), got, want)
				}
				want, panicked := oldLastIndexOf(s, nd, from)
				if panicked {
					want = vm.NIL // the old panic; now nil, as Clojure
				}
				if got := call(lastIndexOf, s, nd, from); fmt.Sprint(got) != fmt.Sprint(want) {
					t.Errorf("last-index-of %q %q %s = %v, want %v", s, nd, show(from), got, want)
				}
			}
		}
	}
}
