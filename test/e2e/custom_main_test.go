/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCustomMain builds the thing this whole change exists for: a third-party
// binary that IS lg — full CLI via pkg/cli — wrapped around interop generated
// for a package let-go has never heard of.
//
// It has to be the real lg and not an lg-runtime-style bytecode host, because
// the custom binary serves two roles at once: it runs scripts (so the native
// namespace must resolve), and it is the compiling side of `lg -b` (top-level
// forms execute at AOT time, so the natives must be linked THERE too). Leg (b)
// covers that second role by using the custom binary as its own bundle base.
//
// Hermetic: a replace directive points at this worktree and the target is a
// stdlib package, so nothing is fetched. The -w subtests go further and run
// with a proxy that cannot serve github.com/nooga/let-go at all, so any fall
// through to @latest fails loudly instead of quietly building a different
// let-go.
func TestCustomMain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping custom-main e2e in -short mode (builds a Go binary)")
	}
	root := repoRoot(t)
	// Everything this test exercises runs in child processes, so the outer
	// test cache cannot see it: an edit to the emitter or to pkg/cli would
	// otherwise return a stale cached pass. Read the inputs here to key the
	// cache on their content — same blind spot trackLginteropInputs covers.
	trackLginteropInputs(t, root)
	trackInputs(t,
		filepath.Join(root, "pkg", "cli", "*.go"),
		filepath.Join(root, "pkg", "gomod", "*.go"),
		filepath.Join(root, "lg.go"))
	tmp := t.TempDir()

	// 1. Generate out-of-tree interop for a stdlib package.
	gen := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "hash/crc32", "-out-pkg", "interop",
		"-out", filepath.Join(tmp, "interop"))
	gen.Dir = root
	if out, err := gen.CombinedOutput(); err != nil {
		t.Fatalf("lginterop failed: %v\n%s", err, out)
	}

	// 2. A custom main: pkg/cli plus a blank import of the generated package.
	//    This is verbatim the pattern docs/guide/custom-lg.md documents.
	customLG := buildCustomHost(t, root, tmp, root)

	// crc32.ChecksumIEEE("hello") — the value is fixed, so a wrong result is a
	// real failure rather than a tautology against the same code path.
	const script = `(require '[crc32])
(println (crc32/ChecksumIEEE (.getBytes "hello")))
`
	scriptPath := filepath.Join(tmp, "script.lg")
	writeFile(t, scriptPath, script)
	const want = "907060870"

	verScript := filepath.Join(tmp, "version.lg")
	writeFile(t, verScript, `(println (System/getProperty "let-go.version"))`+"\n")

	t.Run("runs a script against the generated namespace", func(t *testing.T) {
		cmd := exec.Command(customLG, scriptPath)
		cmd.Dir = tmp
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("custom lg failed to run script: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), want) {
			t.Errorf("got %q, want it to contain %q", out, want)
		}
	})

	// The host stamp ("host-v9") is CLI metadata: -v may show it, but the
	// runtime's let-go.version must describe let-go itself — here resolved
	// through the directory replace, so "dev". A custom host leaking its own
	// version into System/getProperty would break runtime feature checks.
	t.Run("host metadata does not overwrite the runtime identity", func(t *testing.T) {
		cmd := exec.Command(customLG, verScript)
		cmd.Dir = tmp
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("custom lg failed to run version script: %v\n%s", err, out)
		}
		if strings.Contains(string(out), "host-v9") {
			t.Errorf("let-go.version leaked the host stamp: %q", out)
		}
		if !strings.Contains(string(out), "dev") {
			t.Errorf("let-go.version = %q, want the replace-resolved \"dev\"", out)
		}
	})

	// -w must resolve let-go through the local `replace`, not the proxy. The
	// replacement path lives in the custom binary's build info; without
	// carrying it into the generated module, gomod falls back to @latest,
	// which fails offline and silently builds a different let-go online.
	// GOPROXY=off is the assertion: any proxy access fails the test.
	t.Run("-w resolves let-go through the local replace, offline", func(t *testing.T) {
		outDir := filepath.Join(tmp, "wasmout")
		cmd := exec.Command(customLG, "-w", outDir, "-w-shell", "none", scriptPath)
		cmd.Dir = tmp
		cmd.Env = append(envWithout(os.Environ(), "LETGO_SRC"), "GOPROXY=off", "GOFLAGS=-mod=mod")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("custom lg -w failed with GOPROXY=off: %v\n%s", err, out)
		}
		if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
			t.Errorf("no wasm output produced: %v", err)
		}
	})

	// A relative directory replace (`=> ../let-go`) is recorded verbatim in
	// build info, which carries no module root to anchor it. Resolving it
	// against the working directory would honor whatever the user happens to
	// be standing in, so -w must refuse and name LETGO_SRC — and LETGO_SRC
	// must then actually rescue the build. Run from a foreign directory so a
	// cwd-relative guess cannot accidentally land on the checkout.
	t.Run("relative replace is rejected until LETGO_SRC names the checkout", func(t *testing.T) {
		// On macOS t.TempDir() sits under /var/folders, a symlink to
		// /private/var that the go tool resolves; a lexical Rel across that
		// boundary yields a path that misses. Resolve both sides first.
		tmpReal, err := filepath.EvalSymlinks(tmp)
		if err != nil {
			t.Fatal(err)
		}
		rootReal, err := filepath.EvalSymlinks(root)
		if err != nil {
			t.Fatal(err)
		}
		relDir := filepath.Join(tmpReal, "relhost")
		rel, err := filepath.Rel(relDir, rootReal)
		if err != nil {
			t.Fatal(err)
		}
		if filepath.IsAbs(rel) {
			t.Fatalf("expected a relative replace path, got %q", rel)
		}
		relLG := buildCustomHost(t, root, relDir, rel)

		elsewhere := t.TempDir()
		outDir := filepath.Join(tmp, "wasmout-rel")
		base := append(envWithout(os.Environ(), "LETGO_SRC"), "GOPROXY=off", "GOFLAGS=-mod=mod")

		cmd := exec.Command(relLG, "-w", outDir, "-w-shell", "none", scriptPath)
		cmd.Dir = elsewhere
		cmd.Env = base
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("-w with a relative replace and no LETGO_SRC succeeded; it cannot know the checkout:\n%s", out)
		}
		for _, want := range []string{"relative", "LETGO_SRC"} {
			if !strings.Contains(string(out), want) {
				t.Errorf("error output should mention %q:\n%s", want, out)
			}
		}
		if strings.Contains(string(out), "module lookup disabled") {
			t.Errorf("-w fell through to the proxy instead of rejecting the relative replace:\n%s", out)
		}

		cmd = exec.Command(relLG, "-w", outDir, "-w-shell", "none", scriptPath)
		cmd.Dir = elsewhere
		cmd.Env = append(base, "LETGO_SRC="+root)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("-w with LETGO_SRC set failed: %v\n%s", err, out)
		}
		if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
			t.Errorf("no wasm output produced with LETGO_SRC: %v", err)
		}
	})

	// A versioned replace onto a fork (`=> example.com/fork v1.2.3`) must be
	// reproduced whole: reducing it to the version would pin
	// github.com/nooga/let-go@v1.2.3, a different module. The proxy below
	// serves example.com/fork v1.2.3 (this worktree, zipped) and nothing
	// else; let-go's own deps come from the module cache. The only way -w can
	// succeed against it is by requiring the fork, which is the assertion.
	t.Run("versioned fork replace is reproduced, offline", func(t *testing.T) {
		proxy := writeFileProxy(t, root, "example.com/fork", "v1.2.3")
		proxyEnv := []string{"GOPROXY=" + proxy, "GOSUMDB=off", "GOFLAGS=-mod=mod"}
		forkLG := buildCustomHost(t, root, filepath.Join(tmp, "forkhost"), "example.com/fork v1.2.3", proxyEnv...)

		outDir := filepath.Join(tmp, "wasmout-fork")
		cmd := exec.Command(forkLG, "-w", outDir, "-w-shell", "none", scriptPath)
		cmd.Dir = tmp
		cmd.Env = append(envWithout(os.Environ(), "LETGO_SRC"), proxyEnv...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("custom lg -w against a fork replace failed: %v\n%s", err, out)
		}
		if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
			t.Errorf("no wasm output produced: %v", err)
		}

		// Runtime identity follows the fork too: the linked let-go is the
		// replacement at its version, not the v0.0.0 placeholder require.
		ver := exec.Command(forkLG, verScript)
		ver.Dir = tmp
		out, err := ver.CombinedOutput()
		if err != nil {
			t.Fatalf("fork host failed to run version script: %v\n%s", err, out)
		}
		if strings.TrimSpace(string(out)) != "1.2.3" {
			t.Errorf("let-go.version = %q, want the fork's 1.2.3", out)
		}
	})

	// The custom binary is its own bundle base: -b copies the running
	// executable and appends the bytecode, so the standalone binary inherits
	// the generated natives. A stock lg as base would produce a binary that
	// cannot resolve crc32.
	t.Run("bundles with itself as bundle base", func(t *testing.T) {
		bundled := filepath.Join(tmp, "bundled")
		bundle := exec.Command(customLG, "-b", bundled, scriptPath)
		bundle.Dir = tmp
		if out, err := bundle.CombinedOutput(); err != nil {
			t.Fatalf("custom lg -b failed: %v\n%s", err, out)
		}
		run := exec.Command(bundled)
		run.Dir = tmp
		out, err := run.CombinedOutput()
		if err != nil {
			t.Fatalf("bundled binary failed: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), want) {
			t.Errorf("bundled binary got %q, want it to contain %q", out, want)
		}
	})
}

// buildCustomHost writes the documented custom-main module into dir — a
// blank import of the interop package generated under <dir>/../interop's
// module, pkg/cli, and `replace github.com/nooga/let-go => <replaceRHS>` —
// and builds it. replaceRHS is spelled exactly as it would be in go.mod: an
// absolute or relative directory, or `<module> <version>`. Extra env entries
// apply to the build only (a fork replace needs a proxy). Returns the binary.
func buildCustomHost(t *testing.T, root, dir, replaceRHS string, env ...string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	// The interop package is generated once, next to the primary host; the
	// other hosts reach it through a relative directory replace of their own
	// module so a single generation serves all three.
	interopMod := "example.com/customlg"
	replaceInterop := ""
	if _, err := os.Stat(filepath.Join(dir, "interop")); err != nil {
		rel, err := filepath.Rel(dir, filepath.Join(filepath.Dir(dir)))
		if err != nil {
			t.Fatal(err)
		}
		replaceInterop = "\nreplace " + interopMod + " => " + rel + "\n"
	}
	goMod := `module example.com/customlg`
	if replaceInterop != "" {
		goMod = `module example.com/customhost

require example.com/customlg v0.0.0
` + replaceInterop
	}
	writeFile(t, filepath.Join(dir, "go.mod"), goMod+`
go 1.26

require github.com/nooga/let-go v0.0.0

replace github.com/nooga/let-go => `+replaceRHS+`
`)
	writeFile(t, filepath.Join(dir, "main.go"), `package main

import (
	"os"

	_ "example.com/customlg/interop"
	"github.com/nooga/let-go/pkg/cli"
)

func main() { os.Exit(cli.Main("host-v9", "none")) }
`)
	copyFile(t, filepath.Join(root, "go.sum"), filepath.Join(dir, "go.sum"))

	bin := filepath.Join(dir, "customlg")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = dir
	build.Env = append(append(os.Environ(), "GOFLAGS=-mod=mod"), env...)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("custom main (replace => %s) does not build: %v\n%s", replaceRHS, err, out)
	}
	return bin
}

// writeFileProxy lays out a GOPROXY-protocol directory serving exactly one
// module: modPath@version, whose content is the tracked files of the let-go
// checkout at root. Returns the file:// URL to put in GOPROXY. Together with
// GOSUMDB=off it lets a test build against "a fork of let-go" without the
// network, while still being unable to serve github.com/nooga/let-go itself.
func writeFileProxy(t *testing.T, root, modPath, version string) string {
	t.Helper()
	// Not t.TempDir(): its name embeds the subtest name, and a comma or pipe
	// in there would be read by GOPROXY as a list separator.
	base, err := os.MkdirTemp("", "letgo-proxy-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(base) })
	if strings.ContainsAny(base, ",|") {
		t.Fatalf("proxy dir %q would be split by GOPROXY", base)
	}
	dir := filepath.Join(base, modPath, "@v")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "list"), version+"\n")
	writeFile(t, filepath.Join(dir, version+".info"), `{"Version":"`+version+`"}`+"\n")
	copyFile(t, filepath.Join(root, "go.mod"), filepath.Join(dir, version+".mod"))

	// Only tracked files: that is what a real module zip of this repo holds,
	// and it keeps gitignored build output (bin/, .cache/) out of the archive.
	ls := exec.Command("git", "-C", root, "ls-files", "-z")
	listed, err := ls.Output()
	if err != nil {
		t.Fatalf("git ls-files: %v", err)
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	prefix := modPath + "@" + version + "/"
	for _, rel := range bytes.Split(bytes.TrimRight(listed, "\x00"), []byte{0}) {
		p := string(rel)
		fi, err := os.Lstat(filepath.Join(root, p))
		if err != nil || !fi.Mode().IsRegular() {
			continue // deleted in the worktree, or a symlink: not part of a module zip
		}
		w, err := zw.Create(prefix + filepath.ToSlash(p))
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.Open(filepath.Join(root, p))
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(w, f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, version+".zip"), buf.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}
	return "file://" + filepath.ToSlash(base)
}

// envWithout returns env minus every KEY=... entry for key.
func envWithout(env []string, key string) []string {
	out := make([]string, 0, len(env))
	for _, kv := range env {
		if !strings.HasPrefix(kv, key+"=") {
			out = append(out, kv)
		}
	}
	return out
}
