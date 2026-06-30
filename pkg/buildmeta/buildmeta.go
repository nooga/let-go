/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package buildmeta

import (
	"runtime/debug"
	"strings"
)

// Resolve fills missing let-go version/commit metadata from Go build info.
// Explicit ldflag-stamped values always take precedence.
func Resolve(stampedVersion, stampedCommit string, info *debug.BuildInfo) (string, string) {
	resolvedVersion := stampedVersion
	resolvedCommit := stampedCommit
	if info == nil {
		return resolvedVersion, resolvedCommit
	}

	if resolvedVersion == "dev" && !modified(info) {
		if buildVersion := version(info.Main.Version); buildVersion != "" {
			resolvedVersion = buildVersion
		}
	}
	if resolvedCommit == "none" {
		if vcsRevision := vcsRevision(info); vcsRevision != "" {
			resolvedCommit = vcsRevision
		} else if pseudoRevision := pseudoVersionRevision(info.Main.Version); pseudoRevision != "" {
			resolvedCommit = pseudoRevision
		}
	}
	return resolvedVersion, resolvedCommit
}

func version(moduleVersion string) string {
	if moduleVersion == "" || moduleVersion == "(devel)" {
		return ""
	}
	return strings.TrimPrefix(moduleVersion, "v")
}

func vcsRevision(info *debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" && setting.Value != "" {
			return setting.Value
		}
	}
	return ""
}

func modified(info *debug.BuildInfo) bool {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.modified" && setting.Value == "true" {
			return true
		}
	}
	return false
}

func pseudoVersionRevision(moduleVersion string) string {
	parts := strings.Split(moduleVersion, "-")
	if len(parts) < 3 {
		return ""
	}
	revision := parts[len(parts)-1]
	timePart := parts[len(parts)-2]
	if dot := strings.LastIndex(timePart, "."); dot >= 0 {
		timePart = timePart[dot+1:]
	}
	if len(timePart) != 14 || !allDigits(timePart) {
		return ""
	}
	if len(revision) != 12 || !allHex(revision) {
		return ""
	}
	return revision
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return s != ""
}

func allHex(s string) bool {
	for _, r := range s {
		if !hexDigit(r) {
			return false
		}
	}
	return s != ""
}

func hexDigit(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}
