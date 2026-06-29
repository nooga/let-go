// Package perfdata defines the JSON schema shared by the perf tools.
package perfdata

// Baseline is the on-disk format. It is partitioned BY MACHINE PROFILE because
// metrics are mixed-hardware in practice (e.g. Apple M1/M2/M3 — all arm64 but
// different microarchitectures — plus amd64 CI) and the anchor-relative ratio
// is only stable within a single CPU model. Each profile is a self-contained
// sub-baseline (its own machine fingerprint, calibration anchor, and benchmark
// set); the ratchet only ever merges/compares WITHIN one profile, so an M1
// number never contaminates an M3 number. Keep field tags stable.
type Baseline struct {
	Version int `json:"version"`
	// Machines maps a machine-profile key (see MachineKey, e.g. "arm64/Apple M3")
	// to that profile's sub-baseline.
	Machines map[string]MachineBaseline `json:"machines"`
}

// MachineBaseline is one machine profile's self-contained baseline: the host
// that captured it, its calibration anchor, and the per-benchmark metrics. All
// pairing happens within a single MachineBaseline.
type MachineBaseline struct {
	CapturedAt    string                    `json:"captured_at"`
	CapturedAtSHA string                    `json:"captured_at_sha"`
	Machine       Machine                   `json:"machine"`
	Anchor        Anchor                    `json:"anchor"`
	Benchmarks    map[string]BenchmarkEntry `json:"benchmarks"`
}

// Machine fingerprints the host that captured a baseline.
type Machine struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	NumCPU    int    `json:"num_cpu"`
	CPUModel  string `json:"cpu_model"`
	GoVersion string `json:"go_version"`
}

// MachineKey is the partition key for a machine profile: architecture plus CPU
// model, the two attributes that determine relative benchmark cost. GoVersion
// and NumCPU are deliberately NOT in the key (a Go upgrade or a differently-
// provisioned host shouldn't fragment the baseline); they are surfaced as a
// warning on mismatch instead.
func MachineKey(m Machine) string {
	return m.Arch + "/" + m.CPUModel
}

// Anchor captures the absolute speed of the calibration benchmark.
type Anchor struct {
	Name       string            `json:"name"`
	Package    string            `json:"package"`
	NSPerOp    float64           `json:"ns_per_op"`
	Iterations int64             `json:"iterations,omitempty"`
	Samples    []BenchmarkSample `json:"samples,omitempty"`
}

// BenchmarkEntry is one benchmark's current summary plus optional raw samples.
type BenchmarkEntry struct {
	NSPerOp       float64           `json:"ns_per_op"`
	AllocsPerOp   int64             `json:"allocs_per_op"`
	BytesPerOp    int64             `json:"bytes_per_op"`
	RatioToAnchor float64           `json:"ratio_to_anchor"`
	BestSinceSHA  string            `json:"best_since_sha,omitempty"`
	BestSinceAt   string            `json:"best_since_at,omitempty"`
	Samples       []BenchmarkSample `json:"samples,omitempty"`
}

// BenchmarkSample is one raw benchmark measurement retained for statistics.
type BenchmarkSample struct {
	Iterations    int64   `json:"iterations"`
	NSPerOp       float64 `json:"ns_per_op"`
	BytesPerOp    int64   `json:"bytes_per_op"`
	AllocsPerOp   int64   `json:"allocs_per_op"`
	RatioToAnchor float64 `json:"ratio_to_anchor,omitempty"`
	CapturedAt    string  `json:"captured_at,omitempty"`
}

// LegacyBaselineV1 is the pre-multi-machine (version 1) on-disk shape: a single
// machine + anchor + flat benchmark map. readBaseline migrates it into the
// profile-partitioned form.
type LegacyBaselineV1 struct {
	Version       int                       `json:"version"`
	CapturedAt    string                    `json:"captured_at"`
	CapturedAtSHA string                    `json:"captured_at_sha"`
	Machine       Machine                   `json:"machine"`
	Anchor        Anchor                    `json:"anchor"`
	Benchmarks    map[string]BenchmarkEntry `json:"benchmarks"`
}

// ToMachineBaseline returns the v1 baseline as a single MachineBaseline (the
// caller keys it by MachineKey(v.Machine)).
func (v LegacyBaselineV1) ToMachineBaseline() MachineBaseline {
	return MachineBaseline{
		CapturedAt:    v.CapturedAt,
		CapturedAtSHA: v.CapturedAtSHA,
		Machine:       v.Machine,
		Anchor:        v.Anchor,
		Benchmarks:    v.Benchmarks,
	}
}

// StreamRecord is one .jsonl line emitted by bench-ratchet capture.
type StreamRecord struct {
	Package     string  `json:"package"`
	Name        string  `json:"name"`
	Iterations  int64   `json:"iterations"`
	NSPerOp     float64 `json:"ns_per_op"`
	BytesPerOp  int64   `json:"bytes_per_op"`
	AllocsPerOp int64   `json:"allocs_per_op"`
	CapturedAt  string  `json:"captured_at"`
}

// Sample returns the record's benchmark measurement without identity fields.
func (r StreamRecord) Sample() BenchmarkSample {
	return BenchmarkSample{
		Iterations:  r.Iterations,
		NSPerOp:     r.NSPerOp,
		BytesPerOp:  r.BytesPerOp,
		AllocsPerOp: r.AllocsPerOp,
		CapturedAt:  r.CapturedAt,
	}
}
