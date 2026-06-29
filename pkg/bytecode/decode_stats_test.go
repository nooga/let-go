package bytecode

import (
	"strings"
	"testing"
)

func TestDecodeStatsSummarySortsByCount(t *testing.T) {
	var stats DecodeStats
	stats.StringEntries = 3
	stats.StringBytes = 42
	stats.VarRefHits = 11
	stats.VarRefMisses = 4
	stats.VarRefCreates = 2
	stats.Tags[TagIDVector] = 2
	stats.Tags[TagIDInstant] = 7
	stats.Tags[TagIDString] = 5

	got := stats.Summary()
	wantParts := []string{
		"[decode] strings entries=3 bytes=42",
		"[decode] var-ref hits=11 misses=4 creates=2",
		"[decode] tag instant count=7",
		"[decode] tag string count=5",
		"[decode] tag vector count=2",
	}
	for _, want := range wantParts {
		if !strings.Contains(got, want) {
			t.Fatalf("summary missing %q\n%s", want, got)
		}
	}
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) < 5 {
		t.Fatalf("expected at least 5 lines, got %d\n%s", len(lines), got)
	}
	if lines[1] != "[decode] var-ref hits=11 misses=4 creates=2" {
		t.Fatalf("line 2 = %q, want var-ref summary\n%s", lines[1], got)
	}
	if lines[2] != "[decode] tag instant count=7" {
		t.Fatalf("line 3 = %q, want instant first tag\n%s", lines[2], got)
	}
	if lines[3] != "[decode] tag string count=5" {
		t.Fatalf("line 4 = %q, want string second tag\n%s", lines[3], got)
	}
	if lines[4] != "[decode] tag vector count=2" {
		t.Fatalf("line 5 = %q, want vector third tag\n%s", lines[4], got)
	}
}
