package vm

import (
	"math/rand/v2"
	"strings"
	"testing"
	"unicode/utf8"
)

// The helpers must agree with indexing []rune(s), including on invalid UTF-8,
// where each bad byte decodes as one U+FFFD rune.
func TestRuneOffsetMatchesRuneSlice(t *testing.T) {
	corpus := []string{
		"", "a", "héllo", strings.Repeat("é—λ😀x", 20),
		"ab\xe2\x82cd\xffé", "\xc3", "a�b",
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for _, s := range corpus {
		runes := []rune(s)
		for k := 0; k < 500; k++ {
			i := rng.IntN(len(runes)+3) - 1
			bo, ok := RuneOffset(s, i)
			if wantOK := i >= 0 && i <= len(runes); ok != wantOK {
				t.Fatalf("RuneOffset(%q, %d) ok=%v want %v", s, i, ok, wantOK)
			}
			if ok && utf8.RuneCountInString(s[:bo]) != i {
				t.Fatalf("RuneOffset(%q, %d) = %d, prefix has %d runes", s, i, bo, utf8.RuneCountInString(s[:bo]))
			}
			r, ok := RuneAt(s, i)
			if wantOK := i >= 0 && i < len(runes); ok != wantOK || (ok && r != runes[i]) {
				t.Fatalf("RuneAt(%q, %d) = %q,%v", s, i, r, ok)
			}
		}
	}
}

func TestByteSearchable(t *testing.T) {
	for needle, want := range map[string]bool{
		"": true, "a": true, "é—": true,
		"�": false, "a�": false, "\xff": false,
	} {
		if got := ByteSearchable(needle); got != want {
			t.Errorf("ByteSearchable(%q) = %v want %v", needle, got, want)
		}
	}
}
