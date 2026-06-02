// --- Cross-origin isolation: prefer server headers, fall back to SW ---
// When the dev/host server sends the COI headers itself (dev/serve.json
// does this), crossOriginIsolated is already true and we don't need the
// SW. Any leftover SW from a prior visit (e.g. earlier GitHub Pages load)
// would intercept future fetches with stale content — unregister it now
// so the headers path stays clean.
if (crossOriginIsolated && 'serviceWorker' in navigator) {
  navigator.serviceWorker.getRegistrations().then(rs => rs.forEach(r => r.unregister())).catch(()=>{});
}
// No isolation? Register the SW shim — but only once per tab. Without a
// loop guard, a SW that fails to provide isolation (Safari rejects
// credentialless, or activation races a tab close) reloads forever.
if (!crossOriginIsolated && window.isSecureContext && 'serviceWorker' in navigator
    && !sessionStorage.getItem('_lgCoiTried')) {
  sessionStorage.setItem('_lgCoiTried', '1');
  navigator.serviceWorker.register('coi-serviceworker.js').then(() => location.reload()).catch(()=>{});
}

// --- Inline wasm_exec.js and WASM data ---
const WASM_EXEC_JS = __WASM_EXEC_JS__;
const WASM_GZ_B64 = __WASM_GZ_B64__;

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
  const workerCode = `
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
  `;

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