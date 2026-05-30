package vm

import (
	"runtime"
	"testing"
)

// BenchmarkRatchetAnchor is the calibration micro-benchmark used by
// cmd/bench-ratchet to normalize benchmark results across machines.
//
// Design constraints (deliberately frozen — DO NOT modify the loop
// body in this function without also rebasing the entire
// docs/perf/baseline.json):
//
//   - Pure CPU: no allocations, no syscalls, no map / slice / channel.
//   - No project-specific code. Stays valid even if every other package
//     in the repo is rewritten.
//   - Compiler-resistant: the result is escaped via runtime.KeepAlive so
//     dead-code elimination can't fold the loop away.
//   - Cheap per-iteration so b.N converges quickly: each iter is one
//     fused multiply-add over uint64. PCG-XSH-RR step constants chosen
//     because they're well-known and unlikely to ever be optimized
//     specially by the Go compiler.
//
// The ratchet tool divides every other benchmark's ns/op by this
// benchmark's ns/op to get a machine-portable "this is N anchors worth
// of work" number. The baseline file stores both the absolute ns/op
// (so eyeballing same-machine drift still works) and the ratio.
func BenchmarkRatchetAnchor(b *testing.B) {
	var x uint64 = 1
	for i := 0; i < b.N; i++ {
		x = x*6364136223846793005 + 1442695040888963407
	}
	runtime.KeepAlive(x)
}
