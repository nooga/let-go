/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
)

// An invalid byte decodes as U+FFFD, so a cmap keyed on U+FFFD replaces it,
// and one that is not keyed on it leaves the original byte in place. A let-go
// string literal cannot hold an invalid byte, hence a Go test.
func TestStringEscapeInvalidUTF8(t *testing.T) {
	src := vm.String("a\xffb<")
	cases := []struct {
		name string
		cmap vm.Value
		want string
	}{
		{"keyed on U+FFFD", vm.NewPersistentMap([]vm.Value{vm.Char(utf8.RuneError), vm.String("?")}), "a?b<"},
		{"not keyed on U+FFFD", vm.NewPersistentMap([]vm.Value{vm.Char('<'), vm.String("&lt;")}), "a\xffb&lt;"},
	}
	for _, c := range cases {
		got, err := Escape(src, c.cmap)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if got != vm.String(c.want) {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
