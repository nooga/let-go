/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"embed"
	"io/fs"
	"strings"
)

// coreFS is the single embed.FS covering every `.lg` source that ships
// inside pkg/rt/core. Each new `.lg` is picked up automatically — no
// per-file `//go:embed` stub and no entry in pkg/resolver/resolver.go's
// embeddedSources map are required.
//
//go:embed all:core
var coreFS embed.FS

// EmbeddedSource returns the source of an embedded namespace by its
// dotted ns name.
//
// Naming rule (mirrors `cmd/lgbgen.nsToGoPkgName` direction-flipped):
//   - dots are path separators: `ir.passes.dce` → `ir/passes/dce.lg`
//   - in the *leaf* segment, hyphens map to underscores so the file
//     name is a legal Go-style identifier: `ir.lower-go` → `ir/lower_go.lg`
//
// Returns ("", false) if no matching file exists.
func EmbeddedSource(name string) (string, bool) {
	path := "core/" + nsNameToPath(name)
	data, err := coreFS.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

// nsNameToPath applies the EmbeddedSource naming rule.
func nsNameToPath(name string) string {
	parts := strings.Split(name, ".")
	// Mangle only the leaf — interior parts are directories whose names
	// stay verbatim. Hyphens in directory components would be unusual
	// (we don't have any) and treating them as `-` keeps `find`/`grep`
	// usable.
	parts[len(parts)-1] = strings.ReplaceAll(parts[len(parts)-1], "-", "_")
	return strings.Join(parts, "/") + ".lg"
}

// pathToNSName is the inverse of nsNameToPath. Used by EmbeddedNSNames
// to recover ns names from the embedded filesystem walk.
func pathToNSName(path string) string {
	// strip "core/" prefix and ".lg" suffix
	rel := strings.TrimSuffix(strings.TrimPrefix(path, "core/"), ".lg")
	parts := strings.Split(rel, "/")
	parts[len(parts)-1] = strings.ReplaceAll(parts[len(parts)-1], "_", "-")
	return strings.Join(parts, ".")
}

// EmbeddedNSNames returns every embedded namespace name (dotted form),
// derived from the embedded file tree. Used by lgbgen to discover what
// to precompile.
func EmbeddedNSNames() []string {
	var out []string
	_ = fs.WalkDir(coreFS, "core", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".lg") {
			return nil
		}
		out = append(out, pathToNSName(path))
		return nil
	})
	return out
}
