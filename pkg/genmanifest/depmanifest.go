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
	"bufio"
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
// skipping directories, _test.go, and — when skipGenerated is set — generated
// files (both input and generator sweeps set it, so transient generated .go
// artifacts never enter the manifest).
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
			// Skip generated .go too: a generator dir can hold a transient,
			// gitignored generated file (e.g. cmd/lgbgen/main_gogen_ir.go from a
			// prior `-tags gogen_ir` build) that is absent in a clean CI checkout.
			// Including it would make the front-gate staleness check fail there.
			files, err := sweepFiles(repoRoot, s, true)
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

// HashedEdge is an Edge plus the sha256 of its Input at manifest-write time.
type HashedEdge struct {
	Edge
	Sum string
}

const depManifestHeader = "# Auto-generated by `make generate`. DO NOT EDIT.\n" +
	"# <output> <input> <kind> <sha256>   kind in {input, generator}\n"

// WriteDepManifest resolves the edges, hashes each input, and writes the
// space-separated manifest deterministically.
func WriteDepManifest(repoRoot string) error {
	edges, err := Edges(repoRoot)
	if err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString(depManifestHeader)
	for _, e := range edges {
		sum, herr := hashFile(filepath.Join(repoRoot, e.Input))
		if herr != nil {
			return fmt.Errorf("hash %s: %w", e.Input, herr)
		}
		fmt.Fprintf(&b, "%s %s %s %s\n", e.Output, e.Input, e.Kind, sum)
	}
	return os.WriteFile(filepath.Join(repoRoot, DepManifestRelPath), []byte(b.String()), 0644)
}

// ReadDepManifest parses the committed manifest.
func ReadDepManifest(repoRoot string) ([]HashedEdge, error) {
	f, err := os.Open(filepath.Join(repoRoot, DepManifestRelPath))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []HashedEdge
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 4 {
			return nil, fmt.Errorf("malformed manifest line: %q", line)
		}
		out = append(out, HashedEdge{Edge{parts[0], parts[1], parts[2]}, parts[3]})
	}
	return out, sc.Err()
}

// StaleOutputs returns the outputs whose current input/generator file hashes no
// longer match the committed manifest (including outputs entirely absent from
// it). Sorted, deduplicated.
func StaleOutputs(repoRoot string) ([]string, error) {
	recorded, err := ReadDepManifest(repoRoot)
	if err != nil {
		return nil, err
	}
	recByOutput := map[string]map[string]string{} // output -> input -> sum
	for _, he := range recorded {
		if recByOutput[he.Output] == nil {
			recByOutput[he.Output] = map[string]string{}
		}
		recByOutput[he.Output][he.Input] = he.Sum
	}
	current, err := Edges(repoRoot)
	if err != nil {
		return nil, err
	}
	staleSet := map[string]bool{}
	for _, e := range current {
		rec, ok := recByOutput[e.Output][e.Input]
		if !ok {
			staleSet[e.Output] = true
			continue
		}
		sum, herr := hashFile(filepath.Join(repoRoot, e.Input))
		if herr != nil {
			return nil, herr
		}
		if sum != rec {
			staleSet[e.Output] = true
		}
	}
	// An input removed since the manifest was written also dirties the output.
	curByOutput := map[string]map[string]bool{}
	for _, e := range current {
		if curByOutput[e.Output] == nil {
			curByOutput[e.Output] = map[string]bool{}
		}
		curByOutput[e.Output][e.Input] = true
	}
	for _, he := range recorded {
		if !curByOutput[he.Output][he.Input] {
			staleSet[he.Output] = true
		}
	}
	var out []string
	for o := range staleSet {
		out = append(out, o)
	}
	sort.Strings(out)
	return out, nil
}

// CheckDepManifest returns nil when the manifest matches the current sources,
// else an error naming the first stale output and a changed input.
func CheckDepManifest(repoRoot string) error {
	stale, err := StaleOutputs(repoRoot)
	if err != nil {
		return err
	}
	if len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("dependency manifest stale for %d output(s): %s\n%s",
		len(stale), strings.Join(stale, ", "), Remediation)
}
