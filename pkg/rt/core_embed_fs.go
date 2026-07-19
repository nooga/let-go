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
// embed-cache-bust: run! Reduced handling
//
//go:embed all:core
var coreFS embed.FS

// auxEmbeddedSources holds embedded namespace sources that live OUTSIDE the
// core/ tree. They resolve like embedded core but are deliberately excluded
// from the core lowering universe: EmbeddedNSNames walks core/ only, so an aux
// source is invisible to the self-hosting bootstrap (lgbgen). gogen is the case
// (nooga/let-go#425) — it self-contains as embedded source but is the AOT
// Go-emitter, not core content, and enrolling it in core/ perturbs the
// interpreted-vs-native lowering fixpoint.
//
// Concurrency: a bare map with no lock, safe ONLY because it is written
// exclusively from init() (via registerEmbeddedSource), which happens-before
// every resolver read through EmbeddedSource. Do NOT add a runtime
// RegisterEmbeddedSource entry point — a write racing a resolver read would be a
// data race. Keep all registration in init().
var auxEmbeddedSources = map[string]string{}

// registerEmbeddedSource records an auxiliary embedded namespace source. Call it
// ONLY from init() (see auxEmbeddedSources for the race invariant). An empty
// source means the //go:embed directive that feeds it never populated — a broken
// build — so fail loudly and name the namespace here, rather than let it surface
// downstream as an obscure "Can't resolve <ns>/..." through the resolver's
// fallback chain.
func registerEmbeddedSource(name, src string) {
	if src == "" {
		panic("rt: embedded source for namespace " + name + " is empty (//go:embed failed?)")
	}
	auxEmbeddedSources[name] = src
}

// EmbeddedSource returns the source of an embedded namespace by its
// dotted ns name. Auxiliary sources (auxEmbeddedSources) are checked first,
// then the core/ tree.
//
// Naming rule for the core/ tree (the source-side analog of
// `cmd/lgbgen.nsToGoRelDir`, which nests the lowered Go tree the same way):
//   - dots are path separators: `ir.passes.dce` → `ir/passes/dce.lg`
//   - in the *leaf* segment, hyphens map to underscores so the file
//     name is a legal Go-style identifier: `ir.lower-go` → `ir/lower_go.lg`
//
// Returns ("", false) if no matching file exists.
func EmbeddedSource(name string) (string, bool) {
	if src, ok := auxEmbeddedSources[name]; ok {
		return src, true
	}
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
