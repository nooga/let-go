package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestLGBGenUsesSourceBootstrap(t *testing.T) {
	root := repoRoot(t)
	out := filepath.Join(t.TempDir(), "core_compiled.lgb")

	cmd := exec.Command("go", "run", "-tags", "bootstrap", "./cmd/lgbgen", out)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lgbgen source bootstrap failed: %v\n%s", err, output)
	}

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat generated lgb: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("generated lgb is empty")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	// go test runs a test in its package directory, cmd/lgbgen; the source
	// path from runtime.Caller is not usable, since -trimpath rewrites it.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

// buildLGBGen builds the bootstrap lgbgen once for a test that runs it
// several times.
func buildLGBGen(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "lgbgen")
	cmd := exec.Command("go", "build", "-tags", "bootstrap", "-o", bin, "./cmd/lgbgen")
	cmd.Dir = repoRoot(t)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build lgbgen: %v\n%s", err, output)
	}
	return bin
}

// A tree lowered for another module holds the namespaces asked for and nothing
// else, and its packages import each other under the caller's prefix, so the
// module needs no replace onto a let-go checkout.
func TestGoTargetLowersNamedNamespacesUnderImportPrefix(t *testing.T) {
	if testing.Short() {
		t.Skip("lowers clojure.core; a few seconds")
	}
	bin := buildLGBGen(t)
	out := filepath.Join(t.TempDir(), "lgcore")

	cmd := exec.Command(bin, "--target=go", "--ns", "core,walk",
		"--import-prefix", "example.com/app/lgcore", out)
	cmd.Dir = repoRoot(t)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lgbgen: %v\n%s", err, output)
	}

	var pkgs []string
	err := filepath.WalkDir(out, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			rel, _ := filepath.Rel(out, path)
			pkgs = append(pkgs, filepath.ToSlash(rel))
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(src), "core_go_lowered") {
				t.Errorf("%s imports let-go's own lowered tree", rel)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(pkgs)
	if got, want := strings.Join(pkgs, " "), "core/core.go walk/walk.go"; got != want {
		t.Fatalf("lowered files = %q, want %q", got, want)
	}
	walk, err := os.ReadFile(filepath.Join(out, "walk", "walk.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(walk), `"example.com/app/lgcore/core"`) {
		t.Fatal("walk.go does not import core under the caller's prefix")
	}
}

// The caller links exactly the packages it names, so a name that cannot be
// lowered stops the run instead of yielding a smaller tree.
func TestGoTargetRejectsNamespacesItCannotLower(t *testing.T) {
	bin := buildLGBGen(t)
	root := repoRoot(t)
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"unknown namespace", []string{"--target=go", "--ns", "core,no-such-ns", "OUT"}, "no embedded namespace named no-such-ns"},
		{"skip-go namespace", []string{"--target=go", "--ns", "string", "OUT"}, "lgbgen:skip-go"},
		{"no out dir", []string{"--target=go", "--ns", "core"}, "explicit <out-dir>"},
		{"prefix without target", []string{"--import-prefix", "example.com/x"}, "explicit <out-dir>"},
		{"the canonical tree as out dir", []string{"--target=go", "--ns", "walk", "pkg/rt/core_go_lowered"}, "let-go's own lowered tree"},
		{"the canonical tree by another spelling", []string{"--target=go", "--import-prefix", "example.com/x", "./pkg/rt/../rt/core_go_lowered/"}, "let-go's own lowered tree"},
		{"a directory above the canonical tree", []string{"--target=go", "--ns", "walk", "pkg/rt"}, "let-go's own lowered tree"},
		{"a directory inside the canonical tree", []string{"--target=go", "--ns", "walk", "pkg/rt/core_go_lowered/sub/tree"}, "let-go's own lowered tree"},
		{"a namespace with nothing to lower", []string{"--target=go", "--ns", "core,async", "OUT"}, "no functions to lower"},
		{"an empty --ns", []string{"--target=go", "--ns", " , ", "OUT"}, "names no namespace"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "lgcore")
			before := treeFiles(t, filepath.Join(root, "pkg", "rt"))
			args := append([]string(nil), tc.args...)
			for i, a := range args {
				if a == "OUT" {
					args[i] = out
				}
			}
			cmd := exec.Command(bin, args...)
			cmd.Dir = root
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("lgbgen %v succeeded\n%s", tc.args, output)
			}
			if !strings.Contains(string(output), tc.want) {
				t.Fatalf("lgbgen %v: output lacks %q\n%s", tc.args, tc.want, output)
			}
			if entries, _ := os.ReadDir(out); len(entries) > 0 {
				t.Fatalf("lgbgen %v left %d entries in the out dir", tc.args, len(entries))
			}
			if after := treeFiles(t, filepath.Join(root, "pkg", "rt")); after != before {
				t.Fatalf("lgbgen %v changed the files under pkg/rt:\nbefore %s\nafter  %s", tc.args, before, after)
			}
		})
	}
}

// treeFiles fingerprints the regular files under dir by path and size, enough
// to notice a rejected run that still removed or wrote something there.
func treeFiles(t *testing.T, dir string) string {
	t.Helper()
	var n int
	var size int64
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		n++
		size += info.Size()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%d files, %d bytes", n, size)
}
