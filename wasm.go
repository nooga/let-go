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
	"regexp"
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
	useTinyGo := os.Getenv("LETGO_USE_TINYGO") == "1"
	wasmPath := filepath.Join(tmpDir, "app.wasm")
	var build *exec.Cmd
	if useTinyGo {
		fmt.Println("building wasm with tinygo...")
		// -stack-size: TinyGo's wasm target auto-sizes goroutine stacks and falls
		// back to a 64KB default it can't bound against let-go's unbounded
		// interpreter recursion (Frame.Run→Func.Invoke→Frame.Run). Deep recursion
		// overflows the stack into adjacent static memory (the reader's macros map),
		// which manifested as the "3-minute hang". Override via LETGO_TINYGO_STACK
		// for sweeping. See docs-xsofy/tiny-let-go-the-stack-overflow-hunt.md.
		stack := os.Getenv("LETGO_TINYGO_STACK")
		if stack == "" {
			stack = "1MB"
		}
		// -panic=trap keeps the bundle small (a panic becomes a bare wasm
		// unreachable). Override with LETGO_TINYGO_PANIC=print while debugging so
		// the panic message reaches the console instead of a silent trap.
		panicMode := os.Getenv("LETGO_TINYGO_PANIC")
		if panicMode == "" {
			panicMode = "trap"
		}
		tgArgs := []string{"build",
			"-target=wasm", "-no-debug", "-opt=z", "-panic=" + panicMode,
			"-stack-size=" + stack}
		if gc := os.Getenv("LETGO_TINYGO_GC"); gc != "" {
			tgArgs = append(tgArgs, "-gc="+gc)
		}
		tgArgs = append(tgArgs, "-o", wasmPath, ".")
		build = exec.Command("tinygo", tgArgs...)
		build.Dir = tmpDir
		build.Env = os.Environ()
	} else {
		fmt.Println("building wasm...")
		build = exec.Command(goTool, wasmassets.GoBuildArgs(wasmPath, buildTags)...)
		build.Dir = tmpDir
		build.Env = append(goEnv, "GOOS=js", "GOARCH=wasm")
	}
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("wasm build: %w", err)
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
		return fmt.Sprintf("module lg-wasm-app\n\ngo 1.26\n\nrequire github.com/nooga/let-go v%s\n", v), nil, nil
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
	goMod := "module lg-wasm-app\n\ngo 1.26\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return "", nil, err
	}
	get := exec.Command(goToolPath(), "get", "github.com/nooga/let-go@latest")
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

// tinygoFdWriteRe matches TinyGo's WASI fd_write import in its wasm_exec.js.
// TinyGo routes os.Stdout through this (NOT through globalThis.fs), and it
// line-buffers into logLine, emitting to console.log only on LF — so under our
// browser host the terminal output never reaches xterm, and xsofy's newline-less
// ANSI would never flush at all. Non-greedy up to the first `return 0;` (the only
// one in the function) captures the whole body.
var tinygoFdWriteRe = regexp.MustCompile(`(?s)fd_write: function\(fd, iovs_ptr, iovs_len, nwritten_ptr\) \{.*?return 0;\s*\},`)

// patchTinyGoStdout rewrites TinyGo's fd_write so fd 1/2 feed the host terminal
// sink (globalThis.fs.writeSync accumulates into outputBuf; term/flush -> _lgFlush
// posts it to xterm), writing bytes immediately rather than line-buffering. This
// is the seam stock Go's wasm_exec.js already honors via globalThis.fs; TinyGo's
// WASI path bypassed it. See docs-xsofy/tiny-let-go-the-stack-overflow-hunt.md
// "Sequel (2026-07-05)".
func patchTinyGoStdout(data []byte) []byte {
	const replacement = `fd_write: function(fd, iovs_ptr, iovs_len, nwritten_ptr) {
						let nwritten = 0;
						if (fd == 1 || fd == 2) {
							for (let iovs_i = 0; iovs_i < iovs_len; iovs_i++) {
								let iov_ptr = iovs_ptr + iovs_i*8; // wasm32
								let ptr = mem().getUint32(iov_ptr + 0, true);
								let len = mem().getUint32(iov_ptr + 4, true);
								nwritten += len;
								if (globalThis.fs && globalThis.fs.writeSync) {
									globalThis.fs.writeSync(fd, new Uint8Array(mem().buffer, ptr, len));
								}
							}
						}
						mem().setUint32(nwritten_ptr, nwritten, true);
						return 0;
					},`
	if !tinygoFdWriteRe.Match(data) {
		fmt.Fprintln(os.Stderr, "warning: could not patch TinyGo fd_write for terminal output (wasm_exec.js layout changed?)")
		return data
	}
	return tinygoFdWriteRe.ReplaceAll(data, []byte(replacement))
}

func readWasmExecJS() ([]byte, error) {
	if os.Getenv("LETGO_USE_TINYGO") == "1" {
		out, err := exec.Command("tinygo", "env", "TINYGOROOT").Output()
		if err != nil {
			return nil, fmt.Errorf("cannot find TINYGOROOT: %w", err)
		}
		root := strings.TrimSpace(string(out))
		src := filepath.Join(root, "targets", "wasm_exec.js")
		data, err := os.ReadFile(src)
		if err != nil {
			return nil, fmt.Errorf("tinygo wasm_exec.js not found at %s: %w", src, err)
		}
		return patchTinyGoStdout(data), nil
	}
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
