package genmanifest

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// goModuleBuildInputs returns a sorted, deduplicated, slash-separated stable
// superset of main-module sources that can build the generator on any host.
// The canonical Go query seeds package discovery; all local build variants and
// their local imports are then included.

func TestGoModuleBuildInputsIncludesTransitiveLocalBuildInputs(t *testing.T) {
	root := t.TempDir()
	t.Setenv("GOOS", "darwin")
	t.Setenv("GOARCH", "arm64")
	t.Setenv("CGO_ENABLED", "0")
	t.Setenv("GOWORK", filepath.Join(root, "missing.work"))
	t.Setenv("GOFLAGS", "-tags=gogen_ir")
	t.Setenv("GOENV", filepath.Join(root, "missing-goenv"))
	t.Setenv("GOEXPERIMENT", "not-a-real-experiment")
	t.Setenv("GOTOOLCHAIN", "go0.0")
	writeFixtureFile(t, root, "go.mod", "module example.com/fixture\n\ngo 1.23\n")
	writeFixtureFile(t, root, "cmd/tool/main.go", `package main

import (
	"embed"
	"example.com/fixture/internal/direct"
)

//go:embed all:assets
var messageFS embed.FS

func main() { _, _ = direct.Value, messageFS }
`)
	writeFixtureFile(t, root, "cmd/tool/main_bootstrap.go", `//go:build bootstrap

package main

const bootstrapSelected = true
`)
	writeFixtureFile(t, root, "cmd/tool/main_not_bootstrap.go", `//go:build !bootstrap

package main

const bootstrapSelected = false
`)
	writeFixtureFile(t, root, "cmd/tool/assets/message.txt", "embedded fixture\n")
	writeFixtureFile(t, root, "cmd/tool/assets/message.txt.backup", "ignored editor backup\n")
	writeFixtureFile(t, root, "internal/direct/direct.go", `package direct

import "example.com/fixture/internal/transitive"

var Value = transitive.Value
`)
	writeFixtureFile(t, root, "internal/transitive/transitive.go", "package transitive\n\nconst Value = 42\n")
	writeFixtureFile(t, root, "internal/direct/cgo.go", `package direct

/*
#include "native.h"
*/
import "C"

var Native = C.fixture_answer
`)
	writeFixtureFile(t, root, "internal/direct/native.c", "#include \"native.h\"\nint fixture_answer = 7;\n")
	writeFixtureFile(t, root, "internal/direct/native.h", "extern int fixture_answer;\n")

	want := []string{
		"cmd/tool/assets/message.txt",
		"cmd/tool/main.go",
		"cmd/tool/main_bootstrap.go",
		"cmd/tool/main_not_bootstrap.go",
		"internal/direct/cgo.go",
		"internal/direct/direct.go",
		"internal/direct/native.c",
		"internal/direct/native.h",
		"internal/transitive/transitive.go",
	}

	got, err := goModuleBuildInputs(root, "./cmd/tool", "bootstrap")
	if err != nil {
		t.Fatalf("goModuleBuildInputs: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("module-local selected inputs mismatch\n got: %q\nwant: %q", got, want)
	}
}

func TestGoModuleBuildInputsReportsSelectedGoToolFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake Go tool requires a POSIX host")
	}
	root := t.TempDir()
	writeFixtureFile(t, root, "go.mod", "module example.com/fixture\n\ngo 1.23\n")
	fake := filepath.Join(root, "fake-go")
	// Delay stderr so the assertion catches callers that stop reading and kill the
	// tool as soon as stdout fails to decode; diagnostics must include both streams.
	if err := os.WriteFile(fake, []byte("#!/bin/sh\necho fake-go-stdout\nsleep 0.1\necho fake-go-stderr >&2\nexit 17\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LETGO_GO", fake)
	_, err := goModuleBuildInputs(root, "./cmd/tool", "bootstrap")
	if err == nil {
		t.Fatal("fake Go tool failure was accepted")
	}
	for _, want := range []string{"fake-go", "fake-go-stdout", "fake-go-stderr"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q lacks %q", err, want)
		}
	}
}

func TestLGBGenUsesGoToolModuleBuildInputs(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	want := goListModuleInputsOracle(t, root, "./cmd/lgbgen", "bootstrap")

	got, err := goModuleBuildInputs(root, "./cmd/lgbgen", "bootstrap")
	if err != nil {
		t.Fatalf("goModuleBuildInputs: %v", err)
	}
	if !sort.StringsAreSorted(got) {
		t.Fatalf("paths are not sorted: %q", got)
	}
	for i, path := range got {
		if filepath.IsAbs(path) || filepath.ToSlash(path) != path {
			t.Errorf("path must be repo-relative and slash-separated: %q", path)
		}
		if i > 0 && got[i-1] == path {
			t.Errorf("duplicate path: %q", path)
		}
	}
	gotSet := make(map[string]bool, len(got))
	for _, path := range got {
		gotSet[path] = true
	}
	for _, selected := range want {
		if isEphemeralBackup(selected) {
			continue
		}
		if !gotSet[selected] {
			t.Errorf("Go-selected dependency missing from stable superset: %s", selected)
		}
	}
	for _, variant := range []string{"pkg/rt/types_linux.go", "pkg/rt/types_other.go"} {
		if !gotSet[variant] {
			t.Errorf("platform build variant missing from stable superset: %s", variant)
		}
	}

	edges, err := Edges(root)
	if err != nil {
		t.Fatalf("Edges: %v", err)
	}
	for _, output := range []string{"pkg/rt/core_compiled.lgb", "pkg/rt/core_go_lowered/"} {
		var generators []string
		for _, edge := range edges {
			if edge.Output == output && edge.Kind == "generator" {
				generators = append(generators, edge.Input)
			}
		}
		sort.Strings(generators)
		if !reflect.DeepEqual(generators, got) {
			t.Errorf("%s generator edges do not match stable module closure\n got: %q\nwant: %q", output, generators, got)
		}
	}
}

type listedGoPackage struct {
	Dir          string
	Module       *struct{ Main bool }
	GoFiles      []string
	CgoFiles     []string
	CFiles       []string
	CXXFiles     []string
	MFiles       []string
	HFiles       []string
	FFiles       []string
	SFiles       []string
	SwigFiles    []string
	SwigCXXFiles []string
	SysoFiles    []string
	EmbedFiles   []string
}

func goListModuleInputsOracle(t *testing.T, root, packagePattern string, tags ...string) []string {
	t.Helper()
	args := []string{"list", "-deps", "-json"}
	if len(tags) != 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	args = append(args, packagePattern)
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	ignoredEnv := map[string]bool{
		"CGO_ENABLED": true, "GO111MODULE": true, "GOAMD64": true,
		"GOARCH": true, "GOENV": true, "GOEXPERIMENT": true,
		"GOFLAGS": true, "GOOS": true, "GOTOOLCHAIN": true, "GOWORK": true,
	}
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if !ignoredEnv[strings.ToUpper(key)] {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env,
		"GOOS=linux", "GOARCH=amd64", "GOAMD64=v1", "CGO_ENABLED=1",
		"GO111MODULE=on", "GOWORK=off", "GOFLAGS=", "GOENV=off",
		"GOEXPERIMENT=", "GOTOOLCHAIN=local",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", strings.Join(cmd.Args, " "), err, out)
	}

	seen := make(map[string]bool)
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var pkg listedGoPackage
		if err := dec.Decode(&pkg); err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("decode go list JSON: %v", err)
		}
		if pkg.Module == nil || !pkg.Module.Main {
			continue
		}
		groups := [][]string{
			pkg.GoFiles, pkg.CgoFiles, pkg.CFiles, pkg.CXXFiles, pkg.MFiles,
			pkg.HFiles, pkg.FFiles, pkg.SFiles, pkg.SwigFiles,
			pkg.SwigCXXFiles, pkg.SysoFiles, pkg.EmbedFiles,
		}
		for _, files := range groups {
			for _, file := range files {
				absolute := filepath.Join(pkg.Dir, file)
				rel, err := filepath.Rel(root, absolute)
				if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					t.Fatalf("module-local file %q is outside root %q", absolute, root)
				}
				seen[filepath.ToSlash(rel)] = true
			}
		}
	}
	paths := make([]string, 0, len(seen))
	for path := range seen {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func writeFixtureFile(t *testing.T, root, rel, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
