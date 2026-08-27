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
			name: "empty dep version leaves no usable version",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom", Version: "v9.9.9"},
				Deps: []*runtimeDebug.Module{dep(gomod.ModulePath, "")},
			},
			stamped: "9.9.9",
			want:    "dev",
		},
		{
			// The documented local-replace setup: `require ... v0.0.0` plus a
			// directory replace. The require's placeholder names a release
			// that does not exist; the replacement (no version) must win and
			// send the build down the local-source path.
			name: "directory replace overrides the v0.0.0 require placeholder",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom", Version: "v9.9.9"},
				Deps: []*runtimeDebug.Module{{
					Path: gomod.ModulePath, Version: "v0.0.0",
					Replace: &runtimeDebug.Module{Path: gomod.ModulePath},
				}},
			},
			stamped: "9.9.9",
			want:    "dev",
		},
		{
			name: "versioned replace pins to the replacement's version",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom"},
				Deps: []*runtimeDebug.Module{{
					Path: gomod.ModulePath, Version: "v0.0.0",
					Replace: &runtimeDebug.Module{Path: gomod.ModulePath, Version: "v1.2.4"},
				}},
			},
			stamped: "dev",
			want:    "1.2.4",
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

// What rt exposes as let-go.version/let-go.commit must describe the runtime,
// not the host CLI. A custom binary calling cli.Main("host-v9", ...) may not
// leak "host-v9" into System/getProperty feature checks.
func TestRuntimeMetadataFrom(t *testing.T) {
	tests := []struct {
		name                string
		info                *runtimeDebug.BuildInfo
		hostVer, hostCommit string
		wantVer, wantCommit string
	}{
		{
			name:    "no build info falls back to the host stamp",
			info:    nil,
			hostVer: "1.2.3", hostCommit: "abc",
			wantVer: "1.2.3", wantCommit: "abc",
		},
		{
			name:    "stock lg: let-go is the main module, host stamp IS the runtime",
			info:    &runtimeDebug.BuildInfo{Main: runtimeDebug.Module{Path: gomod.ModulePath}},
			hostVer: "1.2.3", hostCommit: "abc",
			wantVer: "1.2.3", wantCommit: "abc",
		},
		{
			name: "custom host metadata does not overwrite the runtime identity",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom", Version: "v9.0.0"},
				Deps: []*runtimeDebug.Module{{Path: gomod.ModulePath, Version: "v1.2.3"}},
			},
			hostVer: "host-v9", hostCommit: "hostcommit000",
			wantVer: "1.2.3", wantCommit: "none",
		},
		{
			name: "pseudo-version dep yields its revision as the commit",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom"},
				Deps: []*runtimeDebug.Module{{Path: gomod.ModulePath, Version: "v0.0.0-20260101000000-abcdef123456"}},
			},
			hostVer: "host-v9", hostCommit: "hostcommit000",
			wantVer: "0.0.0-20260101000000-abcdef123456", wantCommit: "abcdef123456",
		},
		{
			name: "directory replace resolves to dev/none, not the host stamp",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom"},
				Deps: []*runtimeDebug.Module{{
					Path: gomod.ModulePath, Version: "v0.0.0",
					Replace: &runtimeDebug.Module{Path: gomod.ModulePath},
				}},
			},
			hostVer: "host-v9", hostCommit: "hostcommit000",
			wantVer: "dev", wantCommit: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver, com := runtimeMetadataFrom(tt.info, tt.hostVer, tt.hostCommit)
			if ver != tt.wantVer || com != tt.wantCommit {
				t.Errorf("runtimeMetadataFrom() = (%q, %q), want (%q, %q)",
					ver, com, tt.wantVer, tt.wantCommit)
			}
		})
	}
}
