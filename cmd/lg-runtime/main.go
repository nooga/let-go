/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 *
 * lg-runtime is the runtime-only let-go binary: it executes precompiled
 * bytecode (.lgb) with no reader, compiler, or resolver linked in. The
 * guarantee is structural — this package simply never imports pkg/compiler
 * or pkg/resolver (enforced by TestRuntimeOnlyDepGraph) — so there is no
 * eval / load-string / read-string and no dynamic source require: the
 * deployed artifact runs only bytecode compiled ahead of time by a trusted
 * toolchain.
 *
 *   lg-runtime <program.lgb> [args...]   run a precompiled program
 *   lg-runtime -v                        print version
 *
 * A binary produced by `lg -b app.lg -bundle-base lg-runtime <out>` carries
 * an appended payload; lg-runtime detects it at startup and runs it, giving
 * a standalone, compiler-free bundle.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	runtimeDebug "runtime/debug"
	"strings"

	"github.com/nooga/let-go/pkg/buildmeta"
	"github.com/nooga/let-go/pkg/bundle"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"

	_ "github.com/nooga/let-go/pkg/rt/corefns"
)

// Set by goreleaser via ldflags — same contract as the full lg build.
var (
	version = "dev"
	commit  = "none"
)

var showVersion bool

func init() {
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s <program.lgb> [args...]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func versionString() string {
	if commit != "none" && len(commit) > 7 {
		return fmt.Sprintf("%s (%s)", version, commit[:7])
	}
	return version
}

// resourcePathsFromEnv parses LG_RESOURCE_PATHS into io/resource roots. It
// mirrors resolver.ParseSearchPaths, which this binary cannot import:
// pkg/resolver links the compiler.
func resourcePathsFromEnv() []string {
	raw := os.Getenv("LG_RESOURCE_PATHS")
	if raw == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(raw, string(os.PathListSeparator)) {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func storageIDDefault(script string) string {
	cwd, _ := os.Getwd()
	exe, _ := os.Executable()
	return rt.StorageIDFrom("", script, cwd, exe)
}

func runUnit(data []byte) int {
	unit, err := rt.DecodeExecUnit(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if err := rt.RunExecUnit(unit); err != nil {
		fmt.Fprint(os.Stderr, vm.FormatError(err))
		return 1
	}
	return 0
}

// runBundle executes an appended payload: this binary is a standalone bundle,
// so every arg after the program name is a user arg. Mirrors the bundle path
// in the full lg's runMain, minus the source-loading resolver.
func runBundle(lgbData, resData []byte, bakedStoreID string) int {
	rt.SetCommandLineArgs(os.Args[1:])
	storeID := bakedStoreID
	if storeID == "" {
		storeID = storageIDDefault("")
	}
	rt.InstallPersistentStorage(storeID)

	// Resources are self-contained in a bundle: serve io/resource from the
	// embedded archive only, ignoring the filesystem.
	if resData != nil {
		files, err := rt.DecodeResourceArchive(resData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: decoding embedded resources: %v\n", err)
			return 1
		}
		rt.SetResourceProvider(rt.NewEmbeddedResourceProvider(files))
	}
	return runUnit(lgbData)
}

func runMain() int {
	if info, ok := runtimeDebug.ReadBuildInfo(); ok {
		version, commit = buildmeta.Resolve(version, commit, info)
	}
	// Propagate version metadata to the runtime so System/getProperty exposes it.
	rt.Version = version
	rt.Commit = commit

	if err := rt.LoadCore(); err != nil {
		fmt.Fprintln(os.Stderr, "boot:", err)
		return 1
	}
	rt.UseBytecodeNSLoader()
	defer rt.ShutdownAllPods()

	// Appended payload? Then we're a standalone bundle — run it directly,
	// before flag parsing (its args belong to the program).
	if lgbData, resData, bakedStoreID := bundle.ReadBundledSelf(); lgbData != nil {
		return runBundle(lgbData, resData, bakedStoreID)
	}

	flag.Parse()
	if showVersion {
		fmt.Printf("lg-runtime %s\n", versionString())
		return 0
	}
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return 2
	}

	program := args[0]
	rt.SetCommandLineArgs(args[1:])
	rt.InstallPersistentStorage(storageIDDefault(program))
	if rp := resourcePathsFromEnv(); len(rp) > 0 {
		rt.SetResourceProvider(rt.NewFSResourceProvider(rp))
	}

	data, err := os.ReadFile(program)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	return runUnit(data)
}

func main() {
	if code := runMain(); code != 0 {
		os.Exit(code)
	}
}
