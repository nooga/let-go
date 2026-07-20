/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	wasmassets "github.com/nooga/let-go/pkg/rt/wasm"
	"github.com/nooga/let-go/pkg/vm"
)

// The generated main.go template and its -w-host-eval splice live in
// pkg/rt/wasm (RenderMain), alongside the HTML/JS build assets.
//
// The HTML page and host JS live as embedded assets in pkg/rt/wasm.
// See pkg/rt/wasm.AssembleHTML for the assembly contract.

const coiServiceWorkerJS = `addEventListener('install', () => skipWaiting());
addEventListener('activate', e => e.waitUntil(clients.claim()));
addEventListener('fetch', e => {
  if (e.request.cache === 'only-if-cached' && e.request.mode !== 'same-origin') return;
  e.respondWith(fetch(e.request).then(r => {
    if (r.status === 0) return r;
    const h = new Headers(r.headers);
    // Pass server-set isolation headers through untouched. Overriding them
    // (the previous behavior) broke require-corp setups on dev servers
    // that already provide proper headers, by replacing them with the
    // credentialless variant Safari rejects — yielding pages that look
    // like they should be isolated but aren't.
    if (!h.has('Cross-Origin-Embedder-Policy')) {
      // require-corp is the broadest-compatible option: Safari, Firefox,
      // and Chrome all accept it; credentialless is Chrome-only.
      h.set('Cross-Origin-Embedder-Policy', 'require-corp');
    }
    if (!h.has('Cross-Origin-Opener-Policy')) {
      h.set('Cross-Origin-Opener-Policy', 'same-origin');
    }
    return new Response(r.body, {status: r.status, statusText: r.statusText, headers: h});
  }).catch(() => new Response(null, {status: 500})));
});
`

func buildWasm(ctx *compiler.Context, nsRes *resolver.NSResolver, src string, outDir string, shell bool, externalWasm bool, hostEval bool, storeID string, customShellTemplate string) error {
	// 1. Compile .lg → .lgb in memory
	ctx.SetSource(src)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	chunk, _, err := ctx.CompileMultiple(f)
	f.Close()
	if err != nil {
		return err
	}

	var lgbBuf bytes.Buffer
	if len(nsRes.LoadedChunks) > 0 {
		mainNS := ctx.CurrentNS().Name()
		nsChunks := make(map[string]*vm.CodeChunk, len(nsRes.LoadedChunks)+1)
		maps.Copy(nsChunks, nsRes.LoadedChunks)
		nsChunks[mainNS] = chunk
		nsOrder := append(nsRes.LoadOrder, mainNS)
		if err := bytecode.EncodeBundleOrdered(&lgbBuf, ctx.Consts(), nsChunks, nsOrder); err != nil {
			return fmt.Errorf("encoding lgb: %w", err)
		}
	} else {
		if err := bytecode.EncodeCompilation(&lgbBuf, ctx.Consts(), chunk); err != nil {
			return fmt.Errorf("encoding lgb: %w", err)
		}
	}

	// 2. Create temp build directory
	tmpDir, err := os.MkdirTemp("", "lg-wasm-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	if err := prepareWasmBuildDirs(tmpDir); err != nil {
		return err
	}
	goEnv := wasmBuildEnv(tmpDir)
	goTool := goToolPath()
	buildTags := strings.TrimSpace(os.Getenv("LG_WASM_BUILD_TAGS"))

	// 3. Write generated source files
	if err := os.WriteFile(filepath.Join(tmpDir, "program.lgb"), lgbBuf.Bytes(), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(wasmassets.RenderMain(storeID, hostEval)), 0644); err != nil {
		return err
	}
	if wasmassets.HasBuildTag(buildTags, "gogen_ir") {
		srcDir, err := findLetGoModuleDir()
		if err != nil {
			return fmt.Errorf("gogen_ir wasm build requires local let-go source for wireup: %w", err)
		}
		if err := wasmassets.WriteGogenIRWireup(tmpDir, srcDir); err != nil {
			return err
		}
	}

	// 4. Write go.mod
	goMod, goSum, err := generateWasmModuleFiles(tmpDir)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}
	if len(goSum) > 0 {
		if err := os.WriteFile(filepath.Join(tmpDir, "go.sum"), goSum, 0644); err != nil {
			return err
		}
	}

	// 5. Build WASM binary to temp dir. We intentionally skip `go mod tidy`:
	// the generated app imports only runtime packages, while tidy also walks
	// test-only deps from the replaced local module and can spuriously pull
	// network-only packages that the wasm build itself does not need.
	fmt.Println("building wasm...")
	wasmPath := filepath.Join(tmpDir, "app.wasm")
	build := exec.Command(goTool, wasmassets.GoBuildArgs(wasmPath, buildTags)...)
	build.Dir = tmpDir
	build.Env = append(goEnv, "GOOS=js", "GOARCH=wasm")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("go build wasm: %w", err)
	}

	// 6. Read the WASM binary. Inline mode gzip+base64s it into the page;
	// external mode ships it as a separate file the loader fetches + streams.
	wasmData, err := os.ReadFile(wasmPath)
	if err != nil {
		return err
	}
	var wasmB64 string
	if !externalWasm {
		fmt.Println("compressing...")
		var gzBuf bytes.Buffer
		gz, _ := gzip.NewWriterLevel(&gzBuf, gzip.BestCompression)
		gz.Write(wasmData)
		gz.Close()
		wasmB64 = base64.StdEncoding.EncodeToString(gzBuf.Bytes())
	}

	// 7. Read wasm_exec.js
	wasmExecJS, err := readWasmExecJS()
	if err != nil {
		return err
	}

	// 8. Build the HTML. shell=false emits the core glue only (no xterm shell
	// / CDN tags). externalWasm=true emits the streaming loader and an empty
	// inline payload (the wasm ships as main.wasm, written below).
	var html string
	if customShellTemplate != "" {
		tmpl, err := os.ReadFile(customShellTemplate)
		if err != nil {
			return fmt.Errorf("read -w-shell template %q: %w", customShellTemplate, err)
		}
		html = wasmassets.AssembleHTMLWithTemplate(string(tmpl), string(wasmExecJS), wasmB64, externalWasm, hostEval)
	} else {
		html = wasmassets.AssembleHTML(string(wasmExecJS), wasmB64, shell, externalWasm, hostEval)
	}

	// 9. Write output
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	outPath := filepath.Join(outDir, "index.html")
	if err := os.WriteFile(outPath, []byte(html), 0644); err != nil {
		return err
	}

	// External mode: emit the raw wasm as a separately-servable asset. The
	// loader fetches it (instantiateStreaming) instead of decoding an inline blob.
	if externalWasm {
		if err := os.WriteFile(filepath.Join(outDir, "main.wasm"), wasmData, 0644); err != nil {
			return err
		}
	}

	// 10. Write coi-serviceworker.js for cross-origin isolation on hosted servers
	if err := os.WriteFile(filepath.Join(outDir, "coi-serviceworker.js"), []byte(coiServiceWorkerJS), 0644); err != nil {
		return err
	}

	fi, _ := os.Stat(outPath)
	fmt.Printf("output: %s (%.1f MB)\n", outPath, float64(fi.Size())/(1024*1024))
	return nil
}

func prepareWasmBuildDirs(tmpDir string) error {
	for _, dir := range []string{
		filepath.Join(tmpDir, ".gocache"),
		filepath.Join(tmpDir, ".gotmp"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func wasmBuildEnv(tmpDir string) []string {
	return append(os.Environ(),
		"GOCACHE="+filepath.Join(tmpDir, ".gocache"),
		"GOTMPDIR="+filepath.Join(tmpDir, ".gotmp"),
	)
}

func goToolPath() string {
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

func generateWasmModuleFiles(tmpDir string) (string, []byte, error) {
	v := version
	if v != "dev" && v != "" && v[0] >= '0' && v[0] <= '9' {
		// Released binary: pin the require to this exact version. We can't
		// return a nil go.sum (go build rejects a module with unresolved
		// sums), so resolve it from the proxy the same way the no-source path
		// below does. `go get` on the single require writes go.sum without the
		// `go mod tidy` that the build step deliberately avoids.
		return resolveWasmModuleFiles(tmpDir, "github.com/nooga/let-go@v"+v)
	}
	// Dev build — try local source first
	srcDir, err := findLetGoModuleDir()
	if err == nil {
		goMod, goSum, err := localWasmModuleFiles(srcDir)
		if err != nil {
			return "", nil, err
		}
		return goMod, goSum, nil
	}
	// No local source — resolve latest version from module proxy
	return resolveWasmModuleFiles(tmpDir, "github.com/nooga/let-go@latest")
}

// resolveWasmModuleFiles scaffolds a minimal module in tmpDir and resolves the
// let-go require named by modRef via `go get`, returning the resulting go.mod
// and go.sum. Used for both released binaries (pinned version) and dev builds
// with no local source tree (@latest).
func resolveWasmModuleFiles(tmpDir, modRef string) (string, []byte, error) {
	goMod := "module lg-wasm-app\n\ngo 1.26\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return "", nil, err
	}
	get := exec.Command(goToolPath(), "get", modRef)
	get.Dir = tmpDir
	get.Stderr = os.Stderr
	if err := get.Run(); err != nil {
		return "", nil, fmt.Errorf("resolving let-go module: %w (set LETGO_SRC for local source)", err)
	}
	// go get wrote the go.mod with the resolved version — read it back
	data, err := os.ReadFile(filepath.Join(tmpDir, "go.mod"))
	if err != nil {
		return "", nil, err
	}
	sum, _ := os.ReadFile(filepath.Join(tmpDir, "go.sum"))
	return string(data), sum, nil
}

func localWasmModuleFiles(srcDir string) (string, []byte, error) {
	modPath := filepath.Join(srcDir, "go.mod")
	modData, err := os.ReadFile(modPath)
	if err != nil {
		return "", nil, err
	}
	modText := string(modData)
	modText = strings.Replace(modText, "module github.com/nooga/let-go", "module lg-wasm-app", 1)
	modText = strings.TrimRight(modText, "\n") + "\n\nrequire github.com/nooga/let-go v0.0.0\n"
	modText = strings.TrimRight(modText, "\n") + fmt.Sprintf("\n\nreplace github.com/nooga/let-go => %s\n", srcDir)
	sumData, err := os.ReadFile(filepath.Join(srcDir, "go.sum"))
	if err != nil && !os.IsNotExist(err) {
		return "", nil, err
	}
	return modText, sumData, nil
}

func findLetGoModuleDir() (string, error) {
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
		data, err := os.ReadFile(filepath.Join(d, "go.mod"))
		if err == nil && strings.Contains(string(data), "module github.com/nooga/let-go") {
			return d
		}
	}
	return ""
}

func mustGetwd() string {
	d, _ := os.Getwd()
	return d
}

func readWasmExecJS() ([]byte, error) {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		goroot = runtime.GOROOT()
	}
	if goroot == "" {
		out, err := exec.Command(goToolPath(), "env", "GOROOT").Output()
		if err != nil {
			return nil, fmt.Errorf("cannot find GOROOT: %w", err)
		}
		goroot = strings.TrimSpace(string(out))
	}
	candidates := []string{
		filepath.Join(goroot, "lib", "wasm", "wasm_exec.js"),
		filepath.Join(goroot, "misc", "wasm", "wasm_exec.js"),
	}
	for _, src := range candidates {
		data, err := os.ReadFile(src)
		if err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("wasm_exec.js not found in GOROOT (%s)", goroot)
}
