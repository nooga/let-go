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
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	wasmassets "github.com/nooga/let-go/pkg/rt/wasm"
	"github.com/nooga/let-go/pkg/vm"
)

const wasmMainTmpl = `package main

import (
	_ "embed"
	"bytes"
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

//go:embed program.lgb
var lgbData []byte

func main() {
	consts := vm.NewConsts()
	ns := rt.NS("user")
	ctx := compiler.NewCompiler(consts, ns)
	nsResolver := resolver.NewNSResolver(ctx, []string{"."})
	rt.SetNSLoader(nsResolver)

	resolve := func(nsName, name string) *vm.Var {
		n := rt.DefNSBare(nsName)
		v := n.LookupLocal(vm.Symbol(name))
		if v == nil {
			return n.DefStub(name)
		}
		return v
	}

	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(lgbData), resolve)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %%v\n", err)
		return
	}

	for _, name := range unit.NSOrder {
		chunk := unit.NSChunks[name]
		if chunk == nil || chunk == unit.MainChunk {
			continue
		}
		f := vm.NewFrame(chunk, nil)
		_, err := f.RunProtected()
		vm.ReleaseFrame(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading %%s: %%v\n", name, err)
			return
		}
	}

	f := vm.NewFrame(unit.MainChunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	if err != nil {
		fmt.Fprint(os.Stderr, vm.FormatError(err))
	}
}
`

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

func buildWasm(ctx *compiler.Context, nsRes *resolver.NSResolver, src string, outDir string) error {
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

	// 3. Write generated source files
	if err := os.WriteFile(filepath.Join(tmpDir, "program.lgb"), lgbBuf.Bytes(), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(wasmMainTmpl), 0644); err != nil {
		return err
	}

	// 4. Write go.mod
	goMod, err := generateWasmGoMod(tmpDir)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}

	// 5. Resolve dependencies
	fmt.Println("resolving dependencies...")
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = tmpDir
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	// 6. Build WASM binary to temp dir
	fmt.Println("building wasm...")
	wasmPath := filepath.Join(tmpDir, "app.wasm")
	build := exec.Command("go", "build", "-o", wasmPath, ".")
	build.Dir = tmpDir
	build.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("go build wasm: %w", err)
	}

	// 7. Read and compress WASM binary
	fmt.Println("compressing...")
	wasmData, err := os.ReadFile(wasmPath)
	if err != nil {
		return err
	}
	var gzBuf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&gzBuf, gzip.BestCompression)
	gz.Write(wasmData)
	gz.Close()
	wasmB64 := base64.StdEncoding.EncodeToString(gzBuf.Bytes())

	// 8. Read wasm_exec.js
	wasmExecJS, err := readWasmExecJS()
	if err != nil {
		return err
	}

	// 9. Build single self-contained HTML
	html := wasmassets.AssembleHTML(string(wasmExecJS), wasmB64)

	// 10. Write output
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	outPath := filepath.Join(outDir, "index.html")
	if err := os.WriteFile(outPath, []byte(html), 0644); err != nil {
		return err
	}

	// 11. Write coi-serviceworker.js for cross-origin isolation on hosted servers
	if err := os.WriteFile(filepath.Join(outDir, "coi-serviceworker.js"), []byte(coiServiceWorkerJS), 0644); err != nil {
		return err
	}

	fi, _ := os.Stat(outPath)
	fmt.Printf("output: %s (%.1f MB)\n", outPath, float64(fi.Size())/(1024*1024))
	return nil
}

func generateWasmGoMod(tmpDir string) (string, error) {
	v := version
	if v != "dev" && v != "" && v[0] >= '0' && v[0] <= '9' {
		return fmt.Sprintf("module lg-wasm-app\n\ngo 1.26\n\nrequire github.com/nooga/let-go v%s\n", v), nil
	}
	// Dev build — try local source first
	srcDir, err := findLetGoModuleDir()
	if err == nil {
		return fmt.Sprintf("module lg-wasm-app\n\ngo 1.26\n\nrequire github.com/nooga/let-go v0.0.0\n\nreplace github.com/nooga/let-go => %s\n", srcDir), nil
	}
	// No local source — resolve latest version from module proxy
	goMod := "module lg-wasm-app\n\ngo 1.26\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return "", err
	}
	get := exec.Command("go", "get", "github.com/nooga/let-go@latest")
	get.Dir = tmpDir
	get.Stderr = os.Stderr
	if err := get.Run(); err != nil {
		return "", fmt.Errorf("resolving let-go module: %w (set LETGO_SRC for local source)", err)
	}
	// go get wrote the go.mod with the resolved version — read it back
	data, err := os.ReadFile(filepath.Join(tmpDir, "go.mod"))
	if err != nil {
		return "", err
	}
	return string(data), nil
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
		out, err := exec.Command("go", "env", "GOROOT").Output()
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
