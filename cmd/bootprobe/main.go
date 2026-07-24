/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// bootprobe reports where core boot time goes: bundle decode vs definition
// replay vs everything else, using LoadCoreBundle's OnPhase hooks. It boots
// compiler-free (same surface as cmd/lg-runtime), so the numbers isolate the
// runtime boot path.
//
// Usage:
//
//	bootprobe                              one boot, phase breakdown + counts
//	bootprobe -decodes 300                 also re-decode the bundle N times
//	bootprobe -replays 300                 also re-run the core chunk N times
//	bootprobe -decodes 300 -cpuprofile p   CPU-profile the loop
//
// Loops report min/median/max so a quiet-host median is easy to extract;
// run the binary itself repeatedly to sample process-level variance.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"

	_ "github.com/nooga/let-go/pkg/rt/corefns"
)

func medianOf(xs []float64) float64 {
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	mid := len(s) / 2
	if len(s)%2 == 0 {
		return (s[mid-1] + s[mid]) / 2
	}
	return s[mid]
}

func startProfile(path string) func() {
	if path == "" {
		return func() {}
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return pprof.StopCPUProfile
}

func main() {
	replays := flag.Int("replays", 0, "re-run the core main chunk N times after boot")
	decodes := flag.Int("decodes", 0, "re-decode the core bundle N times after boot")
	cpuprofile := flag.String("cpuprofile", "", "write a CPU profile of the replay/decode loop")
	flag.Parse()

	var decodeMS, coreMS float64
	t0 := time.Now()
	unit, err := rt.LoadCoreBundle(rt.CoreLoadOptions{
		EagerHybrids: false,
		OnPhase: func(phase string, since time.Time) {
			d := float64(time.Since(since).Microseconds()) / 1000
			switch phase {
			case "decode-bundle":
				decodeMS = d
			case "run-core-chunk":
				coreMS = d
			}
		},
	})
	totalMS := float64(time.Since(t0).Microseconds()) / 1000
	if err != nil {
		fmt.Fprintln(os.Stderr, "boot failed:", err)
		os.Exit(1)
	}

	nses := rt.AllNSes()
	coreVars, totalVars := 0, 0
	for name, ns := range nses {
		totalVars += ns.RegistrySize()
		if name == rt.NameCoreNS {
			coreVars = ns.RegistrySize()
		}
	}

	fmt.Printf("phases_ms decode=%.3f core_replay=%.3f rest=%.3f total=%.3f\n",
		decodeMS, coreMS, totalMS-decodeMS-coreMS, totalMS)
	fmt.Printf("bundle_bytes=%d ns_chunks=%d core_vars=%d total_vars=%d nses=%d\n",
		len(rt.CoreCompiledLGB), len(unit.NSChunks), coreVars, totalVars, len(nses))

	if *decodes > 0 || *replays > 0 {
		stop := startProfile(*cpuprofile)
		defer stop()
	}

	if *decodes > 0 {
		ts := make([]float64, 0, *decodes)
		for i := 0; i < *decodes; i++ {
			t := time.Now()
			if _, derr := rt.DecodeExecUnit(rt.CoreCompiledLGB); derr != nil {
				fmt.Fprintln(os.Stderr, "re-decode failed:", derr)
				os.Exit(1)
			}
			ts = append(ts, float64(time.Since(t).Microseconds())/1000)
		}
		sort.Float64s(ts)
		fmt.Printf("re_decode_ms n=%d min=%.3f median=%.3f max=%.3f\n",
			*decodes, ts[0], medianOf(ts), ts[len(ts)-1])
	}

	if *replays > 0 {
		ts := make([]float64, 0, *replays)
		for i := 0; i < *replays; i++ {
			t := time.Now()
			fr := vm.NewFrame(unit.MainChunk, nil)
			_, rerr := fr.RunProtected()
			vm.ReleaseFrame(fr)
			if rerr != nil {
				fmt.Fprintln(os.Stderr, "re-replay failed:", rerr)
				os.Exit(1)
			}
			ts = append(ts, float64(time.Since(t).Microseconds())/1000)
		}
		sort.Float64s(ts)
		fmt.Printf("re_replay_ms n=%d min=%.3f median=%.3f max=%.3f\n",
			*replays, ts[0], medianOf(ts), ts[len(ts)-1])
	}
}
