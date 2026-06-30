/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package buildmeta

import (
	"runtime/debug"
	"testing"
)

func TestResolvePreservesStampedValues(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "v1.11.1"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "25ccf5505588cfd6fd0c1ed7e870446ae1ec4a2a"},
		},
	}

	gotVersion, gotCommit := Resolve("1.99.0", "feedface", info)

	if gotVersion != "1.99.0" {
		t.Fatalf("version = %q, want stamped value", gotVersion)
	}
	if gotCommit != "feedface" {
		t.Fatalf("commit = %q, want stamped value", gotCommit)
	}
}

func TestResolveUsesTaggedModuleVersion(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "v1.11.1"},
	}

	gotVersion, gotCommit := Resolve("dev", "none", info)

	if gotVersion != "1.11.1" {
		t.Fatalf("version = %q, want tag without v prefix", gotVersion)
	}
	if gotCommit != "none" {
		t.Fatalf("commit = %q, want unchanged unknown commit", gotCommit)
	}
}

func TestResolveKeepsDirtyBuildVersionUnknown(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "v1.11.1+dirty"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.modified", Value: "true"},
			{Key: "vcs.revision", Value: "25ccf5505588cfd6fd0c1ed7e870446ae1ec4a2a"},
		},
	}

	gotVersion, gotCommit := Resolve("dev", "none", info)

	if gotVersion != "dev" {
		t.Fatalf("version = %q, want unchanged dev version for dirty build", gotVersion)
	}
	if gotCommit != "25ccf5505588cfd6fd0c1ed7e870446ae1ec4a2a" {
		t.Fatalf("commit = %q, want vcs.revision", gotCommit)
	}
}

func TestResolveUsesVCSRevision(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "(devel)"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "25ccf5505588cfd6fd0c1ed7e870446ae1ec4a2a"},
		},
	}

	gotVersion, gotCommit := Resolve("dev", "none", info)

	if gotVersion != "dev" {
		t.Fatalf("version = %q, want unchanged dev version", gotVersion)
	}
	if gotCommit != "25ccf5505588cfd6fd0c1ed7e870446ae1ec4a2a" {
		t.Fatalf("commit = %q, want vcs.revision", gotCommit)
	}
}

func TestResolveUsesPseudoVersionRevision(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "v1.10.1-0.20260628062305-25ccf5505588"},
	}

	gotVersion, gotCommit := Resolve("dev", "none", info)

	if gotVersion != "1.10.1-0.20260628062305-25ccf5505588" {
		t.Fatalf("version = %q, want pseudo-version without v prefix", gotVersion)
	}
	if gotCommit != "25ccf5505588" {
		t.Fatalf("commit = %q, want pseudo-version revision", gotCommit)
	}
}

func TestPseudoVersionRevisionRejectsNonPseudoVersions(t *testing.T) {
	cases := []string{
		"",
		"(devel)",
		"v1.11.1",
		"v1.11.1-alpha.1",
		"v1.10.1-0.20260628062305-nothexvalue",
		"v1.10.1-0.notatime-25ccf5505588",
	}

	for _, tc := range cases {
		if got := pseudoVersionRevision(tc); got != "" {
			t.Fatalf("pseudoVersionRevision(%q) = %q, want empty", tc, got)
		}
	}
}
