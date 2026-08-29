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
// this exists to keep one implementation rather than two.
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
// GenerateFrom is Generate with an explicit let-go source directory, for a
// caller that already knows it. A custom host built against a local `replace`
// has the replacement path in its build info, but FindLetGoSourceDir cannot
// recover it: the host's own module root is not let-go's, so the lookup fails
// and Generate falls back to @latest, which breaks offline and silently builds
// against a different release online. An empty or non-let-go srcDir behaves
// exactly like Generate.
func GenerateFrom(dir, moduleName, version, srcDir string) (Files, error) {
	if srcDir != "" && IsLetGoSourceDir(srcDir) {
		return localFiles(srcDir, moduleName)
	}
	return Generate(dir, moduleName, version)
}

// IsLetGoSourceDir reports whether dir holds let-go's own module.
func IsLetGoSourceDir(dir string) bool {
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
	goMod := "module " + moduleName + "\n\ngo 1.26\n"
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
		if IsLetGoSourceDir(d) {
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
