package main

import (
	"testing"

	"github.com/nooga/let-go/pkg/perfdata"
)

const gbBench = "github.com/nooga/let-go/pkg/compiler.BenchmarkInitFromLGB [bytecode]"

// TestForceRebaselineGatesDeterministicOnGlobalBar pins that a timing-only
// -force update cannot silently RAISE the machine-independent deterministic
// bar. A value that is an improvement over THIS machine's stale row can still
// be a regression against the newest row any profile carries; adopting and
// stamping it makes it the newest provenance, so the gate then reports zero
// regressions and the real one disappears (mparrett, #891 review).
func TestForceRebaselineGatesDeterministicOnGlobalBar(t *testing.T) {
	machine := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3", GoVersion: "go1.27.1"}
	key := perfdata.MachineKey(machine)

	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		// The newest measured row anywhere: the global bar.
		"other": {
			CapturedAt: "2026-09-07T00:00:00Z", CapturedAtSHA: "477a5d36",
			Benchmarks: map[string]BenchmarkEntry{gbBench: {
				NSPerOp: 100, AllocsPerOp: 13995, BytesPerOp: 1016128,
				AllocsSinceSHA: "477a5d36", AllocsSinceAt: "2026-09-07T00:00:00Z",
				BytesSinceSHA: "477a5d36", BytesSinceAt: "2026-09-07T00:00:00Z",
			}},
		},
		// This machine's own row is much older and much worse.
		key: {
			CapturedAt: "2026-08-24T00:00:00Z", CapturedAtSHA: "dadb8e0b", Machine: machine,
			Benchmarks: map[string]BenchmarkEntry{gbBench: {
				NSPerOp: 150, AllocsPerOp: 26294, BytesPerOp: 1599971,
			}},
		},
	}}
	// A timing recapture: faster, and better than the LOCAL row on allocs and
	// bytes, but worse than the GLOBAL bar on both.
	march := MachineBaseline{
		CapturedAt: "2026-09-19T00:00:00Z", CapturedAtSHA: "827ef18b", Machine: machine,
		Benchmarks: map[string]BenchmarkEntry{gbBench: {
			NSPerOp: 110, AllocsPerOp: 14353, BytesPerOp: 1051381,
		}},
	}

	rejected := forceRebaseline(&baseline, key, march, false)

	got := baseline.Machines[key].Benchmarks[gbBench]
	if got.NSPerOp != 110 {
		t.Fatalf("forced timing was not written: ns_per_op = %v, want 110", got.NSPerOp)
	}
	// The stored local values are kept, not the bar's: writing another
	// machine's number here would claim this machine measured it. What matters
	// is that the measured value is not adopted and not stamped as newest.
	if got.AllocsPerOp != 26294 || got.BytesPerOp != 1599971 {
		t.Fatalf("worse-than-global values were adopted: %d allocs/%d bytes, want the stored 26294/1599971",
			got.AllocsPerOp, got.BytesPerOp)
	}
	if got.AllocsSinceSHA == "827ef18b" || got.BytesSinceSHA == "827ef18b" {
		t.Fatalf("worse values were stamped with this run: %q/%q", got.AllocsSinceSHA, got.BytesSinceSHA)
	}
	if len(rejected) != 2 {
		t.Fatalf("rejected = %+v, want the allocs and bytes regressions against the global bar", rejected)
	}

	// The regression the gate was reporting before this update must still be
	// reported after it.
	cur := MachineBaseline{Benchmarks: map[string]BenchmarkEntry{gbBench: {AllocsPerOp: 14353, BytesPerOp: 1051381}}}
	if n := compareDeterministic(machineIndependentBar(baseline), cur, 0.02); n != 2 {
		t.Fatalf("compareDeterministic reports %d regressions after the forced update, want 2", n)
	}
}
