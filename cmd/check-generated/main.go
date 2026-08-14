/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// check-generated exits non-zero when the committed generated artifacts
// (pkg/rt/core_compiled.lgb + the lowered Go tree) are stale relative to
// their .lg / lgbgen sources. Content-based, so it is reliable across
// git/jj checkouts. Used by the Makefile `check-generated` target, by
// CI, and by the git pre-commit hook (scripts/pre-commit).
//
// With -write it instead recomputes the digest and rewrites the canonical
// manifest (a light refresh — the digest only, no bundle/tree rebuild).
// With -o PATH it writes the recomputed digest to PATH; that is how the
// generated.sums git merge driver (scripts/git-merge-sums.sh) produces the
// merged manifest at git's output path. Recomputing beats keep-current so
// the signature never lands stale after a merge.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/genmanifest"
)

func main() {
	writeManifest := flag.Bool("write-manifest", false, "write pkg/rt/generated.manifest and exit")
	staleList := flag.Bool("stale", false, "print outputs with stale input hashes (one per line) and exit")
	needsGeneration := flag.Bool("needs-generation", false,
		"print outputs with stale inputs or missing/incomplete generated files and exit")
	writeCanonical := flag.Bool("write", false,
		"recompute the source digest and rewrite the canonical manifest, then exit")
	outPath := flag.String("o", "",
		"recompute the source digest and write the manifest to this path (for the merge driver)")
	goTool := flag.String("go", "", "Go executable used to resolve generator dependencies")
	flag.Parse()
	if *goTool != "" {
		if err := os.Setenv("LETGO_GO", *goTool); err != nil {
			fmt.Fprintf(os.Stderr, "check-generated: set Go tool: %v\n", err)
			os.Exit(2)
		}
	}

	root, err := genmanifest.FindRepoRoot(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "check-generated: %v\n", err)
		os.Exit(2)
	}

	// Manifest write mode: write the dependency manifest and exit.
	if *writeManifest {
		if err := genmanifest.WriteDepManifest(root); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	// Stale outputs mode: print stale generated outputs and exit.
	if *staleList || *needsGeneration {
		var stale []string
		if *needsGeneration {
			stale, err = genmanifest.StaleOutputs(root)
		} else {
			stale, err = genmanifest.StaleInputOutputs(root)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, s := range stale {
			fmt.Println(s)
		}
		// -stale is a gate: stale inputs fail it. -needs-generation is a query
		// consumed by generate.lg: a stale list is a successful result, while a
		// non-zero exit unambiguously means the query itself failed.
		if *staleList && len(stale) > 0 {
			os.Exit(1)
		}
		return
	}

	// Regenerate modes: recompute the digest from the sources on disk and
	// write it (canonical for -write, an arbitrary path for -o). No bundle
	// or lowered-tree rebuild — the digest is a pure function of the sources.
	if *writeCanonical || *outPath != "" {
		digest, err := genmanifest.Compute(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "check-generated: %v\n", err)
			os.Exit(2)
		}
		target := *outPath
		if target == "" {
			if err := genmanifest.Write(root, digest); err != nil {
				fmt.Fprintf(os.Stderr, "check-generated: %v\n", err)
				os.Exit(2)
			}
		} else if err := genmanifest.WriteTo(target, digest); err != nil {
			fmt.Fprintf(os.Stderr, "check-generated: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("check-generated: wrote manifest digest %s\n", digest[:12])
		return
	}

	res, err := genmanifest.Check(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "check-generated: %v\n", err)
		os.Exit(1)
	}

	if res.Recorded == "" {
		fmt.Fprintf(os.Stderr, "check-generated: manifest %s missing.\n%s\n",
			genmanifest.ManifestRelPath, genmanifest.Remediation)
		os.Exit(1)
	}
	if !res.Fresh {
		fmt.Fprintf(os.Stderr,
			"check-generated: STALE — .lg/lgbgen sources changed without regeneration.\n"+
				"  recorded: %s\n  computed: %s\n%s\n",
			res.Recorded, res.Computed, genmanifest.Remediation)
		os.Exit(1)
	}

	fmt.Println("check-generated: OK — bundle + lowered tree are in sync with sources.")
}
