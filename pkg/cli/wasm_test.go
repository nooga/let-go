/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	runtimeDebug "runtime/debug"
	"strings"
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

// A custom host built against a `replace` must carry the whole directive into
// the generated WASM module. Reducing it to a version pins the wrong module
// for a fork; reducing it to a path loses the version. And a relative
// directory path is unrecoverable from a built binary — build info records
// it verbatim and carries no module root — so it is an error, not a guess.
func TestLetgoReplacementFrom(t *testing.T) {
	withReplace := func(replacePath, replaceVer string) *runtimeDebug.BuildInfo {
		return &runtimeDebug.BuildInfo{
			Main: runtimeDebug.Module{Path: "example.com/custom"},
			Deps: []*runtimeDebug.Module{{
				Path: gomod.ModulePath, Version: "v0.0.0",
				Replace: &runtimeDebug.Module{Path: replacePath, Version: replaceVer},
			}},
		}
	}
	tests := []struct {
		name    string
		info    *runtimeDebug.BuildInfo
		want    *gomod.Replacement
		wantErr string
	}{
		{name: "no build info", info: nil, want: nil},
		{
			name: "no replace",
			info: &runtimeDebug.BuildInfo{
				Main: runtimeDebug.Module{Path: "example.com/custom"},
				Deps: []*runtimeDebug.Module{{Path: gomod.ModulePath, Version: "v1.2.3"}},
			},
			want: nil,
		},
		{
			// Go writes "(devel)" as the version of a directory replacement.
			name: "absolute directory replace with (devel)",
			info: withReplace("/src/let-go", "(devel)"),
			want: &gomod.Replacement{Path: "/src/let-go"},
		},
		{
			name: "absolute directory replace with an empty version",
			info: withReplace("/src/let-go", ""),
			want: &gomod.Replacement{Path: "/src/let-go"},
		},
		{
			name:    "relative directory replace is an error",
			info:    withReplace("../let-go", "(devel)"),
			wantErr: "relative",
		},
		{
			name: "fork replace keeps its module path and version",
			info: withReplace("example.com/fork", "v1.2.3"),
			want: &gomod.Replacement{Path: "example.com/fork", Version: "v1.2.3"},
		},
		{
			name: "same-path pin is carried verbatim too",
			info: withReplace(gomod.ModulePath, "v1.2.4"),
			want: &gomod.Replacement{Path: gomod.ModulePath, Version: "v1.2.4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := letgoReplacementFrom(tt.info)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("letgoReplacementFrom() = %+v, want an error mentioning %q", got, tt.wantErr)
				}
				for _, want := range []string{tt.wantErr, "LETGO_SRC"} {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("error %q should mention %q", err, want)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("letgoReplacementFrom() error: %v", err)
			}
			if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
				t.Errorf("letgoReplacementFrom() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// LETGO_SRC is the explicit override and must win BEFORE build info is
// consulted: it is the only way out for a host whose replace is relative, so
// an error from the unrecoverable replacement may not pre-empt it.
func TestWasmLetgoSource(t *testing.T) {
	replaced := func(path, ver string) *runtimeDebug.BuildInfo {
		return &runtimeDebug.BuildInfo{
			Main: runtimeDebug.Module{Path: "example.com/custom"},
			Deps: []*runtimeDebug.Module{{
				Path: gomod.ModulePath, Version: "v0.0.0",
				Replace: &runtimeDebug.Module{Path: path, Version: ver},
			}},
		}
	}
	t.Run("no override, relative replace: the error surfaces", func(t *testing.T) {
		if _, err := wasmLetgoSource("", replaced("../let-go", "(devel)")); err == nil {
			t.Error("expected the relative-replace error")
		}
	})
	t.Run("override set, relative replace: override wins, no error", func(t *testing.T) {
		rep, err := wasmLetgoSource("/src/let-go", replaced("../let-go", "(devel)"))
		if err != nil || rep != nil {
			t.Errorf("wasmLetgoSource() = (%+v, %v), want (nil, nil)", rep, err)
		}
	})
	t.Run("override set, fork replace: override still wins", func(t *testing.T) {
		rep, err := wasmLetgoSource("/src/let-go", replaced("example.com/fork", "v1.2.3"))
		if err != nil || rep != nil {
			t.Errorf("wasmLetgoSource() = (%+v, %v), want (nil, nil)", rep, err)
		}
	})
	t.Run("no override, absolute replace: the replacement is used", func(t *testing.T) {
		rep, err := wasmLetgoSource("", replaced("/src/let-go", "(devel)"))
		if err != nil || rep == nil || rep.Path != "/src/let-go" || rep.Version != "" {
			t.Errorf("wasmLetgoSource() = (%+v, %v), want the directory replacement", rep, err)
		}
	})
}
