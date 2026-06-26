package ir_test

import (
	"strings"
	"testing"
)

// TestBuildDoesNotDuplicateNestedCallArg guards a build-call defect: a nested
// compound argument was emitted TWICE — once as a dead copy (build-call built
// the args only to switch blocks, then discarded them) and once as the real
// arg (build-call-with-head rebuilt them). The dead copy is harmless for pure
// ops (DCE drops it) but a dead duplicate of an EFFECTFUL call — e.g.
// (chunk b), which consumes its buffer — still executes and corrupts state,
// so a lowered (filter pred (range n)) silently produced an empty result.
// build must emit each argument expression exactly once.
func TestBuildDoesNotDuplicateNestedCallArg(t *testing.T) {
	ensureLoader()
	// (inc a) is a single Inc inst; before the fix build emitted it twice.
	f := buildLispIR(t, `(defn tf [a b] (cons (inc a) b))`)
	dump := lispDump(t, f)
	if n := strings.Count(dump, "= Inc "); n != 1 {
		t.Fatalf("nested arg (inc a) built %d times, want exactly 1:\n--- IR ---\n%s", n, dump)
	}
}
