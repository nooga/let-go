package vm

import (
	"strings"
	"unicode/utf8"
)

// Strings are character-indexed but stored as UTF-8. These helpers map rune
// indices to byte offsets by walking the string in place, so string natives
// don't convert the whole string to []rune on every call.

// RuneOffset returns the byte offset of rune index i in s. ok is false when i
// is negative or past the rune count; i equal to the rune count yields len(s).
func RuneOffset(s string, i int) (int, bool) {
	if i < 0 {
		return 0, false
	}
	ri, bo := 0, 0
	for ri < i && bo < len(s) {
		if s[bo] < utf8.RuneSelf {
			bo++
		} else {
			_, size := utf8.DecodeRuneInString(s[bo:])
			bo += size
		}
		ri++
	}
	return bo, ri == i
}

// RuneAt returns the rune at rune index i of s.
func RuneAt(s string, i int) (rune, bool) {
	bo, ok := RuneOffset(s, i)
	if !ok || bo >= len(s) {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(s[bo:])
	return r, true
}

// ByteSearchable reports whether searching s's bytes for needle finds the
// same matches as searching its runes. UTF-8 resynchronises, so it does for
// any valid needle — except one containing U+FFFD, which as a rune also
// matches each invalid byte of s.
func ByteSearchable(needle string) bool {
	return utf8.ValidString(needle) && !strings.ContainsRune(needle, utf8.RuneError)
}
