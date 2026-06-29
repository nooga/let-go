package vm

import (
	"strings"
	"testing"
)

func TestLookupStatsSummarySortsHotEntries(t *testing.T) {
	stats := LookupStats{
		NamespacedCalls: 12,
		LookupCalls:     9,
		NamespacedBySym: map[string]uint64{
			"foo/bar": 7,
			"baz":     3,
		},
		LookupByNS: map[string]uint64{
			"core": 5,
			"user": 2,
		},
		LookupByNSSym: map[string]uint64{
			"core::map": 4,
			"user::x":   1,
		},
	}

	got := stats.Summary()
	wantParts := []string{
		"[lookup] namespaced calls=12",
		"[lookup] namespace lookup calls=9",
		"[lookup] namespaced symbol foo/bar count=7",
		"[lookup] namespaced symbol baz count=3",
		"[lookup] lookup namespace core count=5",
		"[lookup] lookup namespace user count=2",
		"[lookup] lookup symbol core::map count=4",
		"[lookup] lookup symbol user::x count=1",
	}
	for _, want := range wantParts {
		if !strings.Contains(got, want) {
			t.Fatalf("summary missing %q\n%s", want, got)
		}
	}
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) < 8 {
		t.Fatalf("expected >= 8 lines, got %d\n%s", len(lines), got)
	}
}
