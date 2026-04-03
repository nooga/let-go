/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
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
			return n.Def(name, vm.NIL)
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

// wasmHTMLTemplate has two %%s placeholders: wasm_exec.js source, base64-gzipped WASM.
// Uses xterm.js for terminal rendering. Runs Go WASM in a Web Worker with
// SharedArrayBuffer for blocking read-key. Falls back to output-only mode
// when cross-origin isolation is unavailable.
const wasmHTMLTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>let-go app</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/css/xterm.min.css">
  <style>
    *{margin:0;padding:0;box-sizing:border-box}
    html,body{height:100%%;background:#0c0c0c;color:#e8e6df;overflow:hidden}
    body{display:flex;flex-direction:column}
    #terminal{flex:1}
    #status{color:#5a584f;font-family:monospace;font-size:13px;padding:1rem;position:absolute;
            top:50%%;left:50%%;transform:translate(-50%%,-50%%)}
    .xterm{height:100%%}
  </style>
</head>
<body>
  <div id="status">loading...</div>
  <div id="terminal" style="display:none"></div>
  <script src="https://cdn.jsdelivr.net/npm/@xterm/xterm@5.5.0/lib/xterm.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/@xterm/addon-fit@0.10.0/lib/addon-fit.min.js"></script>
  <script>
// --- Cross-origin isolation via Service Worker ---
if (!crossOriginIsolated && window.isSecureContext && 'serviceWorker' in navigator) {
  navigator.serviceWorker.register('coi-serviceworker.js').then(() => location.reload()).catch(()=>{});
}

// --- Inline wasm_exec.js and WASM data ---
const WASM_EXEC_JS = %s;
const WASM_GZ_B64 = %s;

// --- Decompress gzipped base64 WASM ---
async function decompressWasm(b64) {
  const compressed = Uint8Array.from(atob(b64), c => c.charCodeAt(0));
  const ds = new DecompressionStream('gzip');
  const w = ds.writable.getWriter();
  w.write(compressed); w.close();
  const r = ds.readable.getReader();
  const chunks = [];
  while (true) { const {done,value} = await r.read(); if(done) break; chunks.push(value); }
  let total = 0; for(const c of chunks) total += c.length;
  const out = new Uint8Array(total);
  let off = 0; for(const c of chunks) { out.set(c, off); off += c.length; }
  return out;
}

const status = document.getElementById('status');
const termEl = document.getElementById('terminal');

// --- Initialize xterm.js ---
const term = new Terminal({
  fontFamily: '"IBM Plex Mono", "Menlo", "Consolas", monospace',
  fontSize: 14,
  theme: { background: '#0c0c0c', foreground: '#e8e6df', cursor: '#5ec4b6' },
  allowProposedApi: true,
  convertEol: true,
});
const fitAddon = new FitAddon.FitAddon();
term.loadAddon(fitAddon);

function showTerminal() {
  status.style.display = 'none';
  termEl.style.display = 'block';
  term.open(termEl);
  fitAddon.fit();
  term.focus();
}

window.addEventListener('resize', () => fitAddon.fit());

// --- Worker mode (interactive, needs cross-origin isolation) ---
function startWorkerMode() {
  const sab = new SharedArrayBuffer(64);
  const keyInt32 = new Int32Array(sab);
  const keyUint8 = new Uint8Array(sab, 8, 16);

  showTerminal();

  // Store terminal size in SAB
  Atomics.store(keyInt32, 6, term.cols);
  Atomics.store(keyInt32, 7, term.rows);
  term.onResize(({cols, rows}) => {
    Atomics.store(keyInt32, 6, cols);
    Atomics.store(keyInt32, 7, rows);
  });

  // Send keystrokes to worker via SAB
  term.onData(data => {
    const bytes = new TextEncoder().encode(data);
    if (bytes.length > 16) return;
    // Spin-wait if previous key hasn't been consumed yet
    while (Atomics.load(keyInt32, 0) !== 0) { /* busy wait */ }
    keyUint8.set(bytes);
    Atomics.store(keyInt32, 1, bytes.length);
    Atomics.store(keyInt32, 0, 1);
    Atomics.notify(keyInt32, 0);
  });

  // Build worker code: fs shim + wasm_exec.js + bootstrap
  const workerCode = ` + "`" + `
    let outputBuf = '';
    const decoder = new TextDecoder('utf-8');
    const enosys = () => { const e = new Error("not implemented"); e.code = "ENOSYS"; return e; };
    globalThis.fs = {
      constants: { O_WRONLY:-1, O_RDWR:-1, O_CREAT:-1, O_TRUNC:-1, O_APPEND:-1, O_EXCL:-1, O_DIRECTORY:-1 },
      writeSync(fd, buf) {
        if (fd === 1 || fd === 2) { outputBuf += decoder.decode(buf); return buf.length; }
        return 0;
      },
      write(fd, buf, offset, length, position, callback) {
        if (offset !== 0 || length !== buf.length || position !== null) { callback(enosys()); return; }
        callback(null, this.writeSync(fd, buf));
      },
      chmod(p,m,cb){cb(null);}, chown(p,u,g,cb){cb(null);}, close(fd,cb){cb(null);},
      fchmod(fd,m,cb){cb(null);}, fchown(fd,u,g,cb){cb(null);},
      fstat(fd,cb){cb(null,{isDirectory(){return false;},isFile(){return true;}});},
      fsync(fd,cb){cb(null);}, ftruncate(fd,l,cb){cb(null);},
      lchown(p,u,g,cb){cb(null);}, link(p,l,cb){cb(null);}, lstat(p,cb){cb(null);},
      mkdir(p,m,cb){cb(null);}, open(p,f,m,cb){cb(enosys());},
      read(fd,buf,off,len,pos,cb){cb(null,0);},
      readdir(p,cb){cb(null,[]);}, readlink(p,cb){cb(null,"");},
      rename(o,n,cb){cb(null);}, rmdir(p,cb){cb(null);},
      stat(p,cb){cb(null,{isDirectory(){return false;},isFile(){return true;}});},
      symlink(p,l,cb){cb(null);}, truncate(p,l,cb){cb(null);},
      unlink(p,cb){cb(null);}, utimes(p,a,m,cb){cb(null);},
    };
    globalThis._lgFlush = function() {
      if (outputBuf.length > 0) { postMessage({t:'out', d:outputBuf}); outputBuf = ''; }
    };
    onmessage = async (e) => {
      if (e.data.t !== 'init') return;
      const { sab, wasmGzB64, wasmExecJS } = e.data;
      globalThis._lgKeyInt32 = new Int32Array(sab);
      globalThis._lgKeyUint8 = new Uint8Array(sab, 8, 16);
      // Load wasm_exec.js in worker scope
      eval(wasmExecJS);
      // Decompress WASM
      const compressed = Uint8Array.from(atob(wasmGzB64), c => c.charCodeAt(0));
      const ds = new DecompressionStream('gzip');
      const w = ds.writable.getWriter(); w.write(compressed); w.close();
      const r = ds.readable.getReader();
      const chunks = []; let total = 0;
      while (true) { const {done,value} = await r.read(); if(done) break; chunks.push(value); total += value.length; }
      const wasmBytes = new Uint8Array(total);
      let off = 0; for(const c of chunks) { wasmBytes.set(c, off); off += c.length; }
      // Run Go WASM
      const go = new Go();
      const { instance } = await WebAssembly.instantiate(wasmBytes, go.importObject);
      postMessage({t:'ready'});
      await go.run(instance);
      globalThis._lgFlush();
      postMessage({t:'exit'});
    };
  ` + "`" + `;

  const blob = new Blob([workerCode], { type: 'application/javascript' });
  const worker = new Worker(URL.createObjectURL(blob));

  worker.onmessage = (e) => {
    if (e.data.t === 'out') term.write(e.data.d);
    if (e.data.t === 'exit') term.write('\r\n\x1b[90m[program exited]\x1b[0m\r\n');
  };

  worker.postMessage({ t: 'init', sab, wasmGzB64: WASM_GZ_B64, wasmExecJS: WASM_EXEC_JS });
}

// --- Main-thread mode (output only, no input) ---
async function startMainThreadMode() {
  showTerminal();
  if (location.protocol === 'file:') {
    term.write('\x1b[33mInteractive input requires a local server. Run:\x1b[0m\r\n');
    term.write('\x1b[33m  python3 -m http.server\x1b[0m\r\n');
    term.write('\x1b[33mthen open http://localhost:8000\x1b[0m\r\n\r\n');
  } else {
    term.write('\x1b[33mInteractive input unavailable (no cross-origin isolation).\x1b[0m\r\n');
    term.write('\x1b[33mDeploy coi-serviceworker.js alongside this file.\x1b[0m\r\n\r\n');
  }

  const decoder = new TextDecoder('utf-8');
  const enosys = () => { const e = new Error("not implemented"); e.code = "ENOSYS"; return e; };
  globalThis.fs = {
    constants: { O_WRONLY:-1, O_RDWR:-1, O_CREAT:-1, O_TRUNC:-1, O_APPEND:-1, O_EXCL:-1, O_DIRECTORY:-1 },
    writeSync(fd, buf) {
      if (fd === 1 || fd === 2) { term.write(decoder.decode(buf)); return buf.length; }
      return 0;
    },
    write(fd, buf, offset, length, position, callback) {
      if (offset !== 0 || length !== buf.length || position !== null) {
        callback(enosys()); return;
      }
      callback(null, this.writeSync(fd, buf));
    },
    chmod(p,m,cb){cb(null);}, chown(p,u,g,cb){cb(null);}, close(fd,cb){cb(null);},
    fchmod(fd,m,cb){cb(null);}, fchown(fd,u,g,cb){cb(null);},
    fstat(fd,cb){cb(null,{isDirectory(){return false;},isFile(){return true;}});},
    fsync(fd,cb){cb(null);}, ftruncate(fd,l,cb){cb(null);},
    lchown(p,u,g,cb){cb(null);}, link(p,l,cb){cb(null);}, lstat(p,cb){cb(null);},
    mkdir(p,m,cb){cb(null);}, open(p,f,m,cb){cb(enosys());},
    read(fd,buf,off,len,pos,cb){cb(null,0);},
    readdir(p,cb){cb(null,[]);}, readlink(p,cb){cb(null,"");},
    rename(o,n,cb){cb(null);}, rmdir(p,cb){cb(null);},
    stat(p,cb){cb(null,{isDirectory(){return false;},isFile(){return true;}});},
    symlink(p,l,cb){cb(null);}, truncate(p,l,cb){cb(null);},
    unlink(p,cb){cb(null);}, utimes(p,a,m,cb){cb(null);},
  };
  globalThis._lgFlush = function(){};

  // Load wasm_exec.js
  eval(WASM_EXEC_JS);
  const wasmBytes = await decompressWasm(WASM_GZ_B64);
  const go = new Go();
  const { instance } = await WebAssembly.instantiate(wasmBytes, go.importObject);
  go.run(instance);
}

// --- Entry point ---
(async () => {
  try {
    status.textContent = 'decompressing wasm...';
    if (typeof SharedArrayBuffer !== 'undefined' && crossOriginIsolated) {
      startWorkerMode();
    } else {
      await startMainThreadMode();
    }
  } catch(err) {
    status.textContent = 'error: ' + err;
    console.error(err);
  }
})();
  </script>
</body>
</html>
`

const coiServiceWorkerJS = `addEventListener('install', () => skipWaiting());
addEventListener('activate', e => e.waitUntil(clients.claim()));
addEventListener('fetch', e => {
  if (e.request.cache === 'only-if-cached' && e.request.mode !== 'same-origin') return;
  e.respondWith(fetch(e.request).then(r => {
    if (r.status === 0) return r;
    const h = new Headers(r.headers);
    h.set('Cross-Origin-Embedder-Policy', 'credentialless');
    h.set('Cross-Origin-Opener-Policy', 'same-origin');
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
		for k, v := range nsRes.LoadedChunks {
			nsChunks[k] = v
		}
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
	// JSON-encode JS strings so they're safe to embed in template
	wasmExecJSJSON, _ := json.Marshal(string(wasmExecJS))
	wasmB64JSON, _ := json.Marshal(wasmB64)
	html := fmt.Sprintf(wasmHTMLTemplate, string(wasmExecJSJSON), string(wasmB64JSON))

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
