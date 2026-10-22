package vm

import (
	"slices"
	"testing"
)

// split is called directly so an error from the matcher fails the test;
// through .split, the method fallback turns it into a quiet wrong answer.
// Expected values were taken from the JVM.
func TestSplitKeepsAdjacentEmptyMatchesForAnyPattern(t *testing.T) {
	cases := []struct {
		s, pat string
		want   []string
	}{
		{"abc", "a*|^", []string{"", "", "b", "c", ""}},
		{"abc", `a*|\b`, []string{"", "", "b", "c", ""}},
		{"abc", `\Qa`, []string{"", "bc"}},
		{"a)bc", `\Qa)`, []string{"", "bc"}},
		{"abc", `a*|\Q`, []string{"", "", "b", "c", ""}},
		{"abc", `\Qa\E*|^`, []string{"", "", "b", "c", ""}},
		{"aBc", `(?i)b*|\b`, []string{"a", "", "c", ""}},
	}
	for _, c := range cases {
		v, err := String(c.s).split([]Value{String(c.pat), Int(-1)})
		if err != nil {
			t.Errorf("split %q on %q: %v", c.s, c.pat, err)
			continue
		}
		var got []string
		for _, p := range v.(*TypedArray).Unbox().([]Value) {
			got = append(got, string(p.(String)))
		}
		if !slices.Equal(got, c.want) {
			t.Errorf("split %q on %q = %q, want %q", c.s, c.pat, got, c.want)
		}
	}
}
