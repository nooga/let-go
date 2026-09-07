/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

// Package gomod scaffolds the throwaway Go module that `lg` generates around
// emitted Go sources before handing them to `go build`.
//
// It was extracted from the WASM build path, which was its only caller and had
// the module name baked in. The AOT native path (nooga/let-go#596) needs the
// same three behaviors — pin a released binary's require to its own version,
// prefer a local source tree in dev builds, fall back to the module proxy — so
// this exists to keep one implementation rather than two. A fourth behavior
// serves a custom host built on pkg/cli: reproduce the host's own `replace`
// directive for let-go verbatim (GenerateWithReplace), so its -w builds
// against the let-go it actually links.
package gomod

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ModulePath is the let-go module the generated module requires.
const ModulePath = "github.com/nooga/let-go"

// Files is a generated module's go.mod text and go.sum bytes. Sum is empty when
// the local-source path supplied no go.sum.
type Files struct {
	Mod string
	Sum []byte
}

// Generate produces go.mod/go.sum for a module named moduleName, to be written
// in dir.
//
// version is the lg binary's own version string ("dev" for an unreleased
// build). A released binary pins the require to its exact version so the
// generated program builds against the runtime it was emitted by; a dev build
// prefers the local source tree, and falls back to the module proxy at @latest
// when there isn't one.
// Replacement is the right-hand side of a go.mod `replace` directive for
// let-go, as a dependent binary's build info records it. Version is "" for a
// directory replacement; a module replacement carries the module path and
// its version, fork or not.
type Replacement struct {
	Path    string
	Version string
}

// IsDir reports whether the replacement names a directory rather than a
// module.
func (r Replacement) IsDir() bool { return r.Version == "" }

// GenerateWithReplace produces go.mod/go.sum for a module named moduleName in
// dir, requiring let-go through rep reproduced verbatim. It exists for a
// custom host whose go.mod replaces let-go: the host's module root is not
// let-go's, so FindLetGoSourceDir cannot recover a directory replacement, and
// a module replacement is not a version of github.com/nooga/let-go at all.
// Collapsing either to a version silently built the WASM module against a
// different let-go (or @latest); reproducing the directive builds it against
// the one the host links.
//
// A directory replacement must be absolute and must hold let-go's module.
// Both failures are errors rather than fall-throughs: a binary that recorded
// `../let-go` has genuinely lost the path (build info carries no module
// root), and guessing a root would only move the surprise. LETGO_SRC is the
// override for that case; callers honor it before consulting build info.
func GenerateWithReplace(dir, moduleName string, rep Replacement) (Files, error) {
	if !rep.IsDir() {
		return getAndRead(dir, replacedGoMod(moduleName, rep), ModulePath)
	}
	if !filepath.IsAbs(rep.Path) {
		return Files{}, fmt.Errorf("let-go replace path %q is relative and cannot be resolved from a built binary; set LETGO_SRC to the checkout or build the host with an absolute replace", rep.Path)
	}
	if !isLetGoSourceDir(rep.Path) {
		return Files{}, fmt.Errorf("let-go replace target %q is not a let-go checkout (set LETGO_SRC to one)", rep.Path)
	}
	return localFiles(rep.Path, moduleName)
}

// replacedGoMod is the go.mod a module replacement starts from: the require
// carries the v0.0.0 placeholder, and the replace directive is the one from
// the host, verbatim. `go get` on the require then resolves the replacement
// and fills in the rest.
func replacedGoMod(moduleName string, rep Replacement) string {
	return "module " + moduleName + "\n\ngo 1.26\n\n" +
		"require " + ModulePath + " v0.0.0\n\n" +
		"replace " + ModulePath + " => " + rep.Path + " " + rep.Version + "\n"
}

// isLetGoSourceDir reports whether dir holds let-go's own module.
func isLetGoSourceDir(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	return err == nil && strings.Contains(string(data), "module "+ModulePath)
}

func Generate(dir, moduleName, version string) (Files, error) {
	if ref, ok := ReleaseRef(version); ok {
		// Released binary: pin the require to this exact version. We can't
		// return a nil go.sum (go build rejects a module with unresolved
		// sums), so resolve it from the proxy the same way the no-source path
		// below does. `go get` on the single require writes go.sum without the
		// `go mod tidy` that the build step deliberately avoids.
		return resolveFiles(dir, moduleName, ref)
	}
	// Dev build — try local source first
	srcDir, err := FindLetGoSourceDir()
	if err == nil {
		return localFiles(srcDir, moduleName)
	}
	// No local source — resolve latest version from module proxy
	return resolveFiles(dir, moduleName, ModulePath+"@latest")
}

// ReleaseRef reports the module ref a released binary should pin to, and
// whether this version is a release at all. A version is a release when it
// starts with a digit: "dev" and "" are not, so they fall through to the
// local-source and proxy paths.
func ReleaseRef(version string) (string, bool) {
	if version == "dev" || version == "" || version[0] < '0' || version[0] > '9' {
		return "", false
	}
	return ModulePath + "@v" + version, true
}

// resolveFiles scaffolds a minimal module in dir and resolves the let-go
// require named by modRef via `go get`, returning the resulting go.mod and
// go.sum. Used for both released binaries (pinned version) and dev builds with
// no local source tree (@latest).
func resolveFiles(dir, moduleName, modRef string) (Files, error) {
	return getAndRead(dir, "module "+moduleName+"\n\ngo 1.26\n", modRef)
}

// getAndRead writes goMod into dir, runs `go get modRef` there, and reads
// back the go.mod/go.sum the toolchain produced. With a replace directive
// already in goMod, `go get` on the replaced path resolves the replacement
// and writes sums for it and its dependencies — the same mechanics whether
// the require is pinned, @latest, or replaced.
func getAndRead(dir, goMod, modRef string) (Files, error) {
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		return Files{}, err
	}
	get := exec.Command(GoToolPath(), "get", modRef)
	get.Dir = dir
	get.Stderr = os.Stderr
	if err := get.Run(); err != nil {
		return Files{}, fmt.Errorf("resolving let-go module: %w (set LETGO_SRC for local source)", err)
	}
	// go get wrote the go.mod with the resolved version — read it back
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return Files{}, err
	}
	sum, _ := os.ReadFile(filepath.Join(dir, "go.sum"))
	return Files{Mod: string(data), Sum: sum}, nil
}

// localFiles rewrites the let-go checkout's own go.mod into one for moduleName
// that requires let-go through a replace pointing back at srcDir, so a dev
// build compiles against the working tree rather than a published version.
func localFiles(srcDir, moduleName string) (Files, error) {
	modData, err := os.ReadFile(filepath.Join(srcDir, "go.mod"))
	if err != nil {
		return Files{}, err
	}
	modText := strings.Replace(string(modData), "module "+ModulePath, "module "+moduleName, 1)
	modText = strings.TrimRight(modText, "\n") + "\n\nrequire " + ModulePath + " v0.0.0\n"
	modText = strings.TrimRight(modText, "\n") + fmt.Sprintf("\n\nreplace %s => %s\n", ModulePath, srcDir)
	sumData, err := os.ReadFile(filepath.Join(srcDir, "go.sum"))
	if err != nil && !os.IsNotExist(err) {
		return Files{}, err
	}
	return Files{Mod: modText, Sum: sumData}, nil
}

// FindLetGoSourceDir locates a let-go checkout to build against: LETGO_SRC if
// set, else the module root above the current directory, else the one above the
// running executable.
func FindLetGoSourceDir() (string, error) {
	if src := os.Getenv("LETGO_SRC"); src != "" {
		return src, nil
	}
	if dir := findModuleRoot(mustGetwd()); dir != "" {
		return dir, nil
	}
	if exe, err := os.Executable(); err == nil {
		if dir := findModuleRoot(filepath.Dir(exe)); dir != "" {
			return dir, nil
		}
	}
	return "", fmt.Errorf("cannot find let-go source tree (dev build); set LETGO_SRC or run from source directory")
}

func findModuleRoot(start string) string {
	for d := start; d != "/" && d != "."; d = filepath.Dir(d) {
		if isLetGoSourceDir(d) {
			return d
		}
	}
	return ""
}

func mustGetwd() string {
	d, _ := os.Getwd()
	return d
}

// GoToolPath returns the `go` binary to invoke, preferring the one under the
// GOROOT this binary was built with so a build doesn't silently pick a
// different toolchain off PATH.
func GoToolPath() string {
	if goroot := runtime.GOROOT(); goroot != "" {
		if path := filepath.Join(goroot, "bin", "go"); fileExists(path) {
			return path
		}
	}
	return "go"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
