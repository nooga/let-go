/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// Package rt — profile.go installs Clojure-callable controls for the
// per-opcode dynamic profiler (see pkg/vm/profile.go).
//
// API:
//   (profile/enable!)    — start counting opcode dispatches
//   (profile/disable!)   — stop counting
//   (profile/reset!)     — zero all counters
//   (profile/snapshot)   — returns map with :ops and :pairs vectors

package rt

import (
	"fmt"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/nooga/let-go/pkg/vm"
)

// cpuProfFile holds the open file for an in-progress CPU profile, so
// stop! can close it. nil when no profile is active. Guarded by cpuProfMu
// because Lisp can call cpu-start!/cpu-stop! from goroutines via
// future / async / pmap.
var (
	cpuProfMu   sync.Mutex
	cpuProfFile *os.File
)

func installProfileNS() {
	enable, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("profile/enable!: expected 0 args, got %d", len(vs))
		}
		vm.ProfilingEnabled.Store(true)
		return vm.NIL, nil
	})

	disable, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("profile/disable!: expected 0 args, got %d", len(vs))
		}
		vm.ProfilingEnabled.Store(false)
		return vm.NIL, nil
	})

	reset, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("profile/reset!: expected 0 args, got %d", len(vs))
		}
		vm.ResetProfile()
		return vm.NIL, nil
	})

	snapshot, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("profile/snapshot: expected 0 args, got %d", len(vs))
		}
		// Build a PersistentMap with :ops and :pairs keys.
		// :ops is a vector of [opcode-name count] pairs (already sorted).
		ops := vm.ProfileSnapshot()
		opsVec := make([]vm.Value, 0, len(ops))
		for _, s := range ops {
			opsVec = append(opsVec, vm.NewArrayVector([]vm.Value{
				vm.Keyword(s.Name),
				vm.Int(int(s.Count)),
			}))
		}

		pairs := vm.PairSnapshot()
		pairsVec := make([]vm.Value, 0, len(pairs))
		for _, p := range pairs {
			pairsVec = append(pairsVec, vm.NewArrayVector([]vm.Value{
				vm.Keyword(p.PrevName),
				vm.Keyword(p.CurrName),
				vm.Int(int(p.Count)),
			}))
		}

		result := vm.EmptyPersistentMap
		result = result.Assoc(vm.Keyword("ops"), vm.NewArrayVector(opsVec)).(*vm.PersistentMap)
		result = result.Assoc(vm.Keyword("pairs"), vm.NewArrayVector(pairsVec)).(*vm.PersistentMap)
		return result, nil
	})

	cpuStart, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("profile/cpu-start!: expected 1 arg (path), got %d", len(vs))
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("profile/cpu-start!: arg must be String")
		}
		cpuProfMu.Lock()
		defer cpuProfMu.Unlock()
		if cpuProfFile != nil {
			return vm.NIL, fmt.Errorf("profile/cpu-start!: profile already in progress")
		}
		f, err := os.Create(string(path))
		if err != nil {
			return vm.NIL, fmt.Errorf("profile/cpu-start!: %w", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			return vm.NIL, fmt.Errorf("profile/cpu-start!: %w", err)
		}
		cpuProfFile = f
		return vm.NIL, nil
	})

	cpuStop, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("profile/cpu-stop!: expected 0 args, got %d", len(vs))
		}
		cpuProfMu.Lock()
		defer cpuProfMu.Unlock()
		if cpuProfFile == nil {
			return vm.NIL, fmt.Errorf("profile/cpu-stop!: no profile in progress")
		}
		pprof.StopCPUProfile()
		if err := cpuProfFile.Close(); err != nil {
			cpuProfFile = nil
			return vm.NIL, fmt.Errorf("profile/cpu-stop!: %w", err)
		}
		cpuProfFile = nil
		return vm.NIL, nil
	})

	ns := vm.NewNamespace("profile")
	ns.Def("enable!", enable)
	ns.Def("disable!", disable)
	ns.Def("reset!", reset)
	ns.Def("snapshot", snapshot)
	ns.Def("cpu-start!", cpuStart)
	ns.Def("cpu-stop!", cpuStop)
	RegisterNS(ns)
}
