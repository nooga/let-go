/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	runtimeDebug "runtime/debug"
	"testing"

	"github.com/nooga/let-go/pkg/gomod"
)

// The version handed to gomod.Generate must describe let-go, not the binary.
// For a custom lg the ldflags-stamped version is the HOST module's, and using
// it would pin the generated WASM module to a let-go release that does not
// exist.
func TestLetgoVersionFrom(t *testing.T) {
	dep := func(path, ver string) *runtimeDebug.Module {
		return &runtimeDebug.Module{Path: path, Version: ver}
	}
	tests := []struct {
		name    string
		info    *runtimeDebug.BuildInfo
		stamped string
		want    string
	}{
		{
			name:    "no build info falls back to the stamped version",
			info:    nil,
			stamped: "1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "stock lg: let-go is the main module, stamp wins",
			info:    &runtimeDebug.BuildInfo{Main: runtimeDebug.Module{Path: gomod.ModulePath, Version: "v1.2.3"}},
			stamped: "1.2.3",
			want:    "1.2.3",
		},
		{
			name: "custom lg: let-go's own dep version wins over the host stamp",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom", Version: "v9.9.9"},
				Deps: []*runtimeDebug.Module{dep("example.com/other", "v1.0.0"), dep(gomod.ModulePath, "v1.2.3")},
			},
			stamped: "9.9.9",
			want:    "1.2.3",
		},
		{
			name: "custom lg with a pseudo-version dep",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom"},
				Deps: []*runtimeDebug.Module{dep(gomod.ModulePath, "v0.0.0-20260101000000-abcdef123456")},
			},
			stamped: "dev",
			want:    "0.0.0-20260101000000-abcdef123456",
		},
		{
			name: "replace directive leaves no usable version",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom", Version: "v9.9.9"},
				Deps: []*runtimeDebug.Module{dep(gomod.ModulePath, "")},
			},
			stamped: "9.9.9",
			want:    "dev",
		},
		{
			name: "let-go absent from deps does not leak the host version",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom", Version: "v9.9.9"},
				Deps: []*runtimeDebug.Module{dep("example.com/other", "v1.0.0")},
			},
			stamped: "9.9.9",
			want:    "dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := letgoVersionFrom(tt.info, tt.stamped); got != tt.want {
				t.Errorf("letgoVersionFrom() = %q, want %q", got, tt.want)
			}
		})
	}
}
