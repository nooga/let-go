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

// --- window.LetGoHost — public host API ---
// Single object on window for custom shells to talk to. Surface:
//   onOutput(cb)      register sink for VM stdout; cb(string)
//   onEmit(cb)        register sink for js/emit events; cb(name, parsedData)
//   sendInput(str)    inject UTF-8 keystrokes/bytes toward the VM
//   setSize(c, r)     advertise a new terminal size to the VM
//
// wake() is intentionally NOT in this slice. Unblocking a parked read-key
// without sending real input requires a SAB-level protocol change (a
// wake-epoch cell or tri-state ready flag) so the Go-side read-key can
// distinguish a real keystroke from an unblock. That lands in its own
// PR with a concrete contract.
//
// The internal _lg* globals (_lgKey, _lgSetSize, _lgFlush, _lgEmit) remain
// as compatibility hooks and are also the implementation backing for the
// LetGoHost methods. Callers using either shape keep working.
window.LetGoHost = (function() {
  let outputSink = (s) => console.log(s);
  let emitSink = (name, data) => {
    try {
      window.dispatchEvent(new CustomEvent(name, { detail: data }));
    } catch (err) { console.error('lg emit:', err); }
  };
  return {
    onOutput(cb) { outputSink = cb; },
    onEmit(cb) { emitSink = cb; },
    sendInput(s) { return window._lgKey ? window._lgKey(s) : false; },
    setSize(c, r) { return window._lgSetSize ? window._lgSetSize(c, r) : undefined; },
    // Internal — invoked by the runtime/relay code below.
    _output(s) { outputSink(s); },
    _emit(name, data) { emitSink(name, data); },
  };
})();

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

// xterm is the default output sink in this template. Custom shells can
// call LetGoHost.onOutput(...) before the WASM boots to redirect.
window.LetGoHost.onOutput((s) => term.write(s));

window.addEventListener('resize', () => fitAddon.fit());

// --- Worker mode (interactive, needs cross-origin isolation) ---
function startWorkerMode() {
  const sab = new SharedArrayBuffer(64);
  const keyInt32 = new Int32Array(sab);
  const keyUint8 = new Uint8Array(sab, 8, 16);

  showTerminal();

  // Public input hook — backs LetGoHost.sendInput. Returns true if accepted,
  // false if the input was too long (>16 bytes). xterm's onData calls into
  // this; custom shells can call it directly or via LetGoHost.sendInput.
  window._lgKey = function(data) {
    const bytes = new TextEncoder().encode(data);
    if (bytes.length === 0 || bytes.length > 16) return false;
    while (Atomics.load(keyInt32, 0) !== 0) { /* busy wait */ }
    keyUint8.set(bytes);
    Atomics.store(keyInt32, 1, bytes.length);
    Atomics.store(keyInt32, 0, 1);
    Atomics.notify(keyInt32, 0);
    return true;
  };

  // Public size hook — backs LetGoHost.setSize.
  window._lgSetSize = function(cols, rows) {
    Atomics.store(keyInt32, 6, cols);
    Atomics.store(keyInt32, 7, rows);
  };

  // Initial size + xterm resize hook route through the public API.
  window._lgSetSize(term.cols, term.rows);
  term.onResize(({cols, rows}) => window._lgSetSize(cols, rows));

  // xterm keystrokes feed into _lgKey (which is what the VM reads).
  term.onData((data) => window._lgKey(data));

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
    // Worker side of the js/emit bridge — forward to main thread, which
    // dispatches into LetGoHost (workers have no DOM, no LetGoHost).
    globalThis._lgEmit = function(name, dataJson) {
      postMessage({t:'emit', name, data: dataJson});
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
    if (e.data.t === 'out') window.LetGoHost._output(e.data.d);
    if (e.data.t === 'exit') window.LetGoHost._output('\r\n\x1b[90m[program exited]\x1b[0m\r\n');
    if (e.data.t === 'emit') {
      try { window.LetGoHost._emit(e.data.name, JSON.parse(e.data.data)); }
      catch (err) { console.error('lg emit relay:', err); }
    }
  };

  worker.postMessage({ t: 'init', sab, wasmGzB64: WASM_GZ_B64, wasmExecJS: WASM_EXEC_JS });
}

// --- Main-thread mode (output only, no input) ---
async function startMainThreadMode() {
  showTerminal();
  const out = (s) => window.LetGoHost._output(s);
  if (location.protocol === 'file:') {
    out('\x1b[33mInteractive input requires a local server. Run:\x1b[0m\r\n');
    out('\x1b[33m  python3 -m http.server\x1b[0m\r\n');
    out('\x1b[33mthen open http://localhost:8000\x1b[0m\r\n\r\n');
  } else {
    out('\x1b[33mInteractive input unavailable (no cross-origin isolation).\x1b[0m\r\n');
    out('\x1b[33mDeploy coi-serviceworker.js alongside this file.\x1b[0m\r\n\r\n');
  }

  const decoder = new TextDecoder('utf-8');
  const enosys = () => { const e = new Error("not implemented"); e.code = "ENOSYS"; return e; };
  globalThis.fs = {
    constants: { O_WRONLY:-1, O_RDWR:-1, O_CREAT:-1, O_TRUNC:-1, O_APPEND:-1, O_EXCL:-1, O_DIRECTORY:-1 },
    writeSync(fd, buf) {
      if (fd === 1 || fd === 2) { window.LetGoHost._output(decoder.decode(buf)); return buf.length; }
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
  // Main-thread side of the js/emit bridge — dispatch straight into
  // LetGoHost (no worker round-trip needed).
  globalThis._lgEmit = function(name, dataJson) {
    try { window.LetGoHost._emit(name, JSON.parse(dataJson)); }
    catch (err) { console.error('lg emit:', err); }
  };

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