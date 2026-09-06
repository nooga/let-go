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
	"github.com/nooga/let-go/pkg/gomod"
	"github.com/nooga/let-go/pkg/resolver"
	wasmassets "github.com/nooga/let-go/pkg/rt/wasm"
	"github.com/nooga/let-go/pkg/vm"
)

// wasmModuleName is the module name of the throwaway Go module the WASM build
// scaffolds around its generated sources. It used to be baked into the module
// scaffolder; it is passed in now that the AOT native path shares that code.
const wasmModuleName = "lg-wasm-app"

// tinyGoStackSizeRe matches a tinygo -stack-size value (a byte count or a
// K/M/G-suffixed size, e.g. 1MB, 64KB, 16384). Used only to warn on a typo in
// LETGO_TINYGO_STACK so it's attributed to the env var, not an opaque tinygo error.
var tinyGoStackSizeRe = regexp.MustCompile(`(?i)^\d+(\.\d+)?[kmg]?b?$`)

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

	// -strip splits source maps and local-variable tables out of the embedded
	// program before it is baked into the wasm binary. The companion is written
	// beside the bundle directory rather than inside it: everything in outDir is
	// served, and a debug sidecar is build output, not something to ship to every
	// visitor.
	var debugData []byte
	var debugPath string
	if stripDebug {
		stripped, companion, path, err := wasmassets.SplitProgramDebug(lgbBuf.Bytes(), outDir, debugOutput)
		if err != nil {
			return err
		}
		lgbBuf.Reset()
		lgbBuf.Write(stripped)
		debugData, debugPath = companion, path
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
	goTool := gomod.GoToolPath()
	buildTags := strings.TrimSpace(os.Getenv("LG_WASM_BUILD_TAGS"))

	// 3. Write generated source files
	if err := os.WriteFile(filepath.Join(tmpDir, "program.lgb"), lgbBuf.Bytes(), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(wasmassets.RenderMain(storeID, hostEval)), 0644); err != nil {
		return err
	}
	if wasmassets.HasBuildTag(buildTags, "gogen_ir") {
		srcDir, err := gomod.FindLetGoSourceDir()
		if err != nil {
			return fmt.Errorf("gogen_ir wasm build requires local let-go source for wireup: %w", err)
		}
		if err := wasmassets.WriteGogenIRWireup(tmpDir, srcDir); err != nil {
			return err
		}
	}

	// 4. Write go.mod
	mod, err := gomod.Generate(tmpDir, wasmModuleName, version)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(mod.Mod), 0644); err != nil {
		return err
	}
	if len(mod.Sum) > 0 {
		if err := os.WriteFile(filepath.Join(tmpDir, "go.sum"), mod.Sum, 0644); err != nil {
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
		if stack != "" && !tinyGoStackSizeRe.MatchString(stack) {
			fmt.Fprintf(os.Stderr, "warning: LETGO_TINYGO_STACK=%q is not a valid size (e.g. 1MB, 64KB); passing to tinygo as-is\n", stack)
		}
		if stack == "" {
			stack = "1MB"
		}
		// -panic=trap keeps the bundle small (a panic becomes a bare wasm
		// unreachable). Override with LETGO_TINYGO_PANIC=print while debugging so
		// the panic message reaches the console instead of a silent trap.
		panicMode := os.Getenv("LETGO_TINYGO_PANIC")
		if panicMode != "" && panicMode != "trap" && panicMode != "print" {
			fmt.Fprintf(os.Stderr, "warning: LETGO_TINYGO_PANIC=%q is not trap|print; passing to tinygo as-is\n", panicMode)
		}
		if panicMode == "" {
			panicMode = "trap"
		}
		// -opt=z is the size-lean default; LETGO_TINYGO_OPT=2 trades bundle
		// size for speed (tinygo accepts 0|1|2|s|z).
		optLevel := os.Getenv("LETGO_TINYGO_OPT")
		switch optLevel {
		case "":
			optLevel = "z"
		case "0", "1", "2", "s", "z":
		default:
			fmt.Fprintf(os.Stderr, "warning: LETGO_TINYGO_OPT=%q is not a known tinygo -opt level (0|1|2|s|z); passing as-is\n", optLevel)
		}
		tgArgs := []string{"build",
			"-target=wasm", "-no-debug", "-opt=" + optLevel, "-panic=" + panicMode,
			"-stack-size=" + stack}
		if gc := os.Getenv("LETGO_TINYGO_GC"); gc != "" {
			switch gc {
			case "none", "leaking", "conservative", "precise":
			default:
				fmt.Fprintf(os.Stderr, "warning: LETGO_TINYGO_GC=%q is not a known tinygo gc (none|leaking|conservative|precise); passing as-is\n", gc)
			}
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

	// The debug companion lands last, so a failed build leaves no sidecar
	// claiming to describe a bundle that was never written.
	if debugData != nil {
		if err := os.WriteFile(debugPath, debugData, 0644); err != nil {
			return fmt.Errorf("writing debug companion: %w", err)
		}
	}

	fi, _ := os.Stat(outPath)
	fmt.Printf("output: %s (%.1f MB)\n", outPath, float64(fi.Size())/(1024*1024))
	if debugData != nil {
		fmt.Printf("debug companion: %s (%d bytes)\n", debugPath, len(debugData))
	}
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
		out, err := exec.Command(gomod.GoToolPath(), "env", "GOROOT").Output()
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
