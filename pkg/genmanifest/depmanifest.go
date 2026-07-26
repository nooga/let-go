/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Provenance dependency manifest: which generator INPUTS feed which generated
// OUTPUT, content-hashed, so `make generate` re-runs only the stages whose
// inputs changed. Refines the coarse combined digest (generated.sums) into a
// per-output graph. Content-based (sha256), never mtime — checkout-safe.
package genmanifest

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Edge is one (output, input, kind) dependency. Kind is "input" (a domain
// source: annotation .go, .lg, IR spec) or "generator" (the generator's own
// source — a change there means the generation LOGIC changed).
type Edge struct {
	Output string
	Input  string
	Kind   string
}

// DepManifestRelPath is the committed manifest, co-located with generated.sums.
const DepManifestRelPath = "pkg/rt/generated.manifest"

type sweep struct{ dir, ext string }

// outputSpec declares one generated output and how to find its inputs. A sweep
// walks dir for *.ext (skipping generated + _test.go for input sweeps); files
// are explicit paths.
type outputSpec struct {
	output      string
	inputSweeps []sweep
	inputFiles  []string
	genSweeps   []sweep
	genFiles    []string
}

// outputSpecs is the authoritative edge declaration (see the plan's table).
var outputSpecs = []outputSpec{
	{
		output:      "pkg/rt/zz_primitives_generated.go",
		inputSweeps: []sweep{{"pkg/rt", ".go"}},
		genSweeps:   []sweep{{"internal/primgen", ".go"}, {"cmd/lgprimgen", ".go"}},
	},
	{
		output:      "pkg/rt/corefns/zz_primitives_generated.go",
		inputSweeps: []sweep{{"pkg/rt/corefns", ".go"}},
		genSweeps:   []sweep{{"internal/primgen", ".go"}, {"cmd/lgprimgen", ".go"}},
	},
	{
		output:      "pkg/rt/core_compiled.lgb",
		inputSweeps: []sweep{{"pkg/rt/core", ".lg"}},
		genSweeps:   []sweep{{"cmd/lgbgen", ".go"}},
	},
	{
		output:      "pkg/rt/core_go_lowered/",
		inputSweeps: []sweep{{"pkg/rt/core", ".lg"}},
		genSweeps:   []sweep{{"cmd/lgbgen", ".go"}},
	},
	{
		output:     "pkg/ir/op_generated.go",
		inputFiles: []string{"pkg/ir/ir_ops.lg", "scripts/gen/op_generated.head"},
		genSweeps:  []sweep{{"cmd/lgbgen", ".go"}},
	},
	{
		output:     "pkg/rt/ir_bridge_generated.go",
		inputFiles: []string{"pkg/ir/ir_bridge.lg", "scripts/gen/ir_bridge_generated.head"},
		genSweeps:  []sweep{{"cmd/lgbgen", ".go"}},
	},
	{
		output:     "pkg/rt/core/ir/data/generated.lg",
		inputFiles: []string{"pkg/ir/ir_data.lg"},
		genSweeps:  []sweep{{"cmd/lgbgen", ".go"}},
	},
}

// sweepFiles returns repo-relative paths under root/s.dir with extension s.ext,
// skipping directories, _test.go, and (for input sweeps) generated files.
func sweepFiles(repoRoot string, s sweep, skipGenerated bool) ([]string, error) {
	var out []string
	base := filepath.Join(repoRoot, s.dir)
	err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, s.ext) {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if skipGenerated && s.ext == ".go" {
			gen, gerr := isGenerated(path)
			if gerr != nil {
				return gerr
			}
			if gen {
				return nil
			}
		}
		rel, rerr := filepath.Rel(repoRoot, path)
		if rerr != nil {
			return rerr
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	return out, err
}

// Edges resolves every output's sweeps/files into a sorted, deduplicated edge
// list (no hashes). Missing explicit input files are an error (a stale spec).
func Edges(repoRoot string) ([]Edge, error) {
	seen := map[Edge]bool{}
	var edges []Edge
	add := func(output, input, kind string) {
		e := Edge{output, input, kind}
		if !seen[e] {
			seen[e] = true
			edges = append(edges, e)
		}
	}
	for _, spec := range outputSpecs {
		for _, s := range spec.inputSweeps {
			files, err := sweepFiles(repoRoot, s, true)
			if err != nil {
				return nil, err
			}
			for _, f := range files {
				add(spec.output, f, "input")
			}
		}
		for _, f := range spec.inputFiles {
			if _, err := os.Stat(filepath.Join(repoRoot, f)); err != nil {
				return nil, fmt.Errorf("input file for %s missing: %s: %w", spec.output, f, err)
			}
			add(spec.output, f, "input")
		}
		for _, s := range spec.genSweeps {
			files, err := sweepFiles(repoRoot, s, false)
			if err != nil {
				return nil, err
			}
			for _, f := range files {
				add(spec.output, f, "generator")
			}
		}
		for _, f := range spec.genFiles {
			add(spec.output, f, "generator")
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].Output != edges[j].Output {
			return edges[i].Output < edges[j].Output
		}
		if edges[i].Kind != edges[j].Kind {
			return edges[i].Kind < edges[j].Kind
		}
		return edges[i].Input < edges[j].Input
	})
	return edges, nil
}
