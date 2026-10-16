/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"embed"
	"io/fs"
	"strings"
	"sync"
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
// core/ tree: gogen (nooga/let-go#425), and lg layers that out-of-tree Go
// packages register through RegisterEmbeddedSource (nooga/let-go#1063). They
// resolve like embedded core but are deliberately excluded from the core
// lowering universe: EmbeddedNSNames walks core/ only, so an aux source is
// invisible to the self-hosting bootstrap (lgbgen) and always loads as
// interpreted source. For gogen that's the point — it is the AOT Go-emitter,
// not core content, and enrolling it in core/ perturbs the interpreted-vs-native
// lowering fixpoint. An out-of-tree package doesn't exist when let-go is built,
// so there would be nothing to lower anyway.
//
// Guarded by auxEmbeddedMu. Registration from init() alone would be race-free
// (inits run one at a time), but an exported function can't hold callers to
// that, and sealing the map at the first read would refuse legitimate ones:
// pkg/compiler's own init reads embedded sources, before the init of every
// package that imports it.
var (
	auxEmbeddedMu      sync.RWMutex
	auxEmbeddedSources = map[string]string{}
)

// RegisterEmbeddedSource makes the lg source src resolvable as namespace name,
// with no search path: (require 'name) finds it before any file on
// -source-paths, so a registered namespace shadows a file of the same name.
// It lets a Go package ship the lg layer over its natives inside the package,
// for any binary that links it — a custom lg built on cli.Main, or a -w/-b
// build that imports it.
//
// Call it from init(), next to the package's rt.RegisterNS calls, so the
// namespace exists before any program can require it. Give library namespaces
// a qualified name (mylib.core, not core).
//
// It panics if src is empty (the //go:embed directive that feeds it never
// populated, a broken build), if name is a core namespace or a native one
// already registered, or if name is already registered: each would otherwise
// surface later as an obscure resolution failure, as a source that never loads
// because the native namespace wins, or as one package's source silently
// replacing another's. A native namespace registered after the source isn't
// caught, so register the source under a name of its own, not one of the
// package's rt.RegisterNS names.
func RegisterEmbeddedSource(name, src string) {
	if src == "" {
		panic("rt: embedded source for namespace " + name + " is empty (//go:embed failed?)")
	}
	if _, err := coreFS.ReadFile("core/" + nsNameToPath(name)); err == nil {
		panic("rt: embedded source for namespace " + name + " would shadow the core namespace of that name")
	}
	if LookupNS(name) != nil {
		panic("rt: embedded source for namespace " + name + " would never load: a native namespace of that name is already registered")
	}
	auxEmbeddedMu.Lock()
	defer auxEmbeddedMu.Unlock()
	if _, dup := auxEmbeddedSources[name]; dup {
		panic("rt: embedded source for namespace " + name + " is registered twice")
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
	auxEmbeddedMu.RLock()
	src, ok := auxEmbeddedSources[name]
	auxEmbeddedMu.RUnlock()
	if ok {
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
