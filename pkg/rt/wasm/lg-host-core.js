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

// --- wasm_exec.js + WASM payload ---
// WASM_MODE is 'inline' (gzip-base64 baked into this page, decoded in JS) or
// 'external' (a separate main.wasm fetched + stream-compiled). Both load paths
// ship below; the build picks one via -w-wasm.
const WASM_EXEC_JS = __WASM_EXEC_JS__;
const WASM_MODE = __WASM_MODE__;
const WASM_GZ_B64 = __WASM_GZ_B64__; // inline payload; "" in external mode
const HOST_EVAL = __LG_HOST_EVAL__;  // built with -w-host-eval? gates LetGoHost.eval
// External mode resolves the payload to an absolute URL so the Blob-URL worker
// (whose relative base is the blob, not the page) can fetch it. A client can
// override this to point at its own CDN.
const WASM_URL = WASM_MODE === 'external' ? new URL('main.wasm', location.href).href : null;

// --- Decompress gzipped base64 WASM (inline mode) ---
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
// Single object on window. This is the whole contract between the runtime
// glue (this file) and a shell — the default xterm shell that ships with
// `lg -w`, or a client's own under `-w-shell none`. Surface:
//   onReady(cb)       cb(mode) once glue is wired and the boot mode
//                     ('worker' | 'main') is known, before the VM runs.
//                     A shell opens its UI and binds output/input here.
//   onOutput(cb)      register sink for VM stdout; cb(string)
//   onEmit(cb)        register sink for js/emit events; cb(name, parsedData)
//   sendInput(str)    inject UTF-8 keystrokes/bytes toward the VM
//   setSize(c, r)     advertise a new terminal size to the VM
//   eval(code)        -w-host-eval bundles only: run code in the loaded image,
//                     -> Promise<string> (stringified value or error). Direct
//                     call on the main thread, worker round-trip when isolated.
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
  // Output emitted before a shell registers onOutput is buffered and
  // replayed on the first onOutput, not dropped to console. The boot path
  // runs synchronously to its first await, so in main-thread mode the
  // "interactive input unavailable" notice is written before the default
  // shell (concatenated after core) registers its sink; without the buffer
  // that text would leak to console instead of the terminal. Buffering also
  // covers a client shell that binds onOutput late (e.g. after an async CDN
  // load) under -w-shell none.
  let outputSink = null;
  const outputBuffer = [];
  let emitSink = (name, data) => {
    try {
      window.dispatchEvent(new CustomEvent(name, { detail: data }));
    } catch (err) { console.error('lg emit:', err); }
  };
  let readyCb = null;
  let readyMode = null; // set if _ready fires before onReady registers
  return {
    onReady(cb) {
      readyCb = cb;
      if (readyMode !== null) cb(readyMode); // late registration still fires
    },
    onOutput(cb) {
      outputSink = cb;
      if (outputBuffer.length) { for (const s of outputBuffer.splice(0)) cb(s); }
    },
    onEmit(cb) { emitSink = cb; },
    sendInput(s) { return window._lgKey ? window._lgKey(s) : false; },
    setSize(c, r) { return window._lgSetSize ? window._lgSetSize(c, r) : undefined; },
    // Internal — invoked by the runtime/relay code below.
    _output(s) { if (outputSink) outputSink(s); else outputBuffer.push(s); },
    _emit(name, data) { emitSink(name, data); },
    _ready(mode) { readyMode = mode; if (readyCb) readyCb(mode); },
  };
})();

// --- LetGoHost request/eval plumbing (-w-host-eval bundles only) ---
// runtimeReady resolves once the wasm runtime has installed its _lgEval/_lgRequest
// hooks and called _lgRuntimeReady (wired per boot mode below). requestImpl is
// installed by whichever boot mode runs: a direct _lgRequest call on the main
// thread, or a worker round-trip when cross-origin isolated. The vars stay
// declared in every bundle (the boot relays below reference them), but the public
// request/eval APIs are exposed only when the bundle was built with -w-host-eval:
// a default bundle never calls _lgRuntimeReady, so installed APIs would await a
// ready signal that never comes (hang) and would lie to feature detection.
let runtimeReadyResolve;
const runtimeReady = new Promise((r) => { runtimeReadyResolve = r; });
let requestImpl = null;
let requestSeq = 0;
const requestPending = new Map();
// Worker-path safety net: a worker crash or dropped message means the
// 'request-result' reply never arrives, so without this the pending entry leaks
// and the caller's Promise hangs forever. The timer rejects + removes the entry
// so the leak is bounded and the caller sees a failure. A late reply arriving
// after the timeout is harmless — its id is already gone from requestPending and
// the result is dropped. Generous by design: the worker runs host requests
// synchronously today, so this guards lost responses, not slow-but-valid work.
const REQUEST_TIMEOUT_MS = 30000;
if (HOST_EVAL) {
  window.LetGoHost.request = function(req) {
    return runtimeReady.then(() => requestImpl(req || {}));
  };
  window.LetGoHost.eval = function(code) {
    return window.LetGoHost.request({ op: 'eval', session: 'default', code }).then((resp) => {
      if (!resp || !resp.ok) {
        const msg = resp && resp.error && resp.error.message ? resp.error.message : 'request failed';
        throw new Error(msg);
      }
      return resp.value;
    });
  };
}

// Boot status text. May be absent when a client ships its own HTML without
// a #status element, so every use is guarded.
const status = document.getElementById('status');
function setStatus(t) { if (status) status.textContent = t; }

// --- Worker mode (interactive, needs cross-origin isolation) ---
async function startWorkerMode() {
  const sab = new SharedArrayBuffer(64);
  const keyInt32 = new Int32Array(sab);
  const keyUint8 = new Uint8Array(sab, 8, 16);

  // Input is unsafe until the worker has been told to start (init posted):
  // before that there is no consumer to drain the SAB slot. Drop pre-start
  // keys rather than write into a slot nothing will ever read.
  let workerReady = false;

  // Public input hook — backs LetGoHost.sendInput. Returns true if accepted,
  // false if dropped (worker not started yet, or input too long >16 bytes).
  // A shell's keystrokes feed through here via LetGoHost.sendInput.
  //
  // Non-blocking: never wait for the consumer to drain (that wait is what froze
  // the page on key bursts). We drop rather than overwrite while a key is still
  // pending (ready==1): overwriting concurrently would tear the consumer's
  // byte-by-byte copy and, worse, the consumer's unconditional clear after copy
  // would stomp our just-set ready flag and lose the accepted key.
  //
  // This makes the slot oldest-pending-wins, not newest-wins. In practice when
  // the program is waiting on input the consumer drains immediately and parks at
  // ready==0, so the freshest key still lands (the interactive common case).
  // Drops only happen while a prior key sits unconsumed — i.e. the program is
  // busy and not reading keys — where keeping the first pending key (the
  // interrupt) and dropping the rest is the behavior we want anyway.
  window._lgKey = function(data) {
    if (!workerReady) return false;
    const bytes = new TextEncoder().encode(data);
    if (bytes.length === 0 || bytes.length > 16) return false;
    if (Atomics.load(keyInt32, 0) !== 0) return false;
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

  // Glue is wired (_lgKey / _lgSetSize installed). The shell opens its UI,
  // binds output via onOutput, and pushes input/size via sendInput/setSize.
  window.LetGoHost._ready('worker');

  // Build worker code: fs shim + wasm_exec.js + bootstrap
  const workerCode = `
    const decoder = new TextDecoder('utf-8');
    const enosys = () => { const e = new Error("not implemented"); e.code = "ENOSYS"; return e; };
    globalThis.fs = {
      constants: { O_WRONLY:-1, O_RDWR:-1, O_CREAT:-1, O_TRUNC:-1, O_APPEND:-1, O_EXCL:-1, O_DIRECTORY:-1 },
      // App output now flows through _lgOutput (the Go-side HostWriter), not
      // fd 1/2 interception. This only catches Go-runtime direct writes
      // (e.g. panics) — surface those on the worker console.
      writeSync(fd, buf) {
        if (fd === 1 || fd === 2) { console.log(decoder.decode(buf)); return buf.length; }
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
    // Worker side of the output bridge — the Go HostWriter calls _lgOutput;
    // relay each chunk to the main thread, which feeds LetGoHost._output.
    globalThis._lgOutput = function(s) {
      postMessage({t:'out', d:s});
    };
    globalThis._lgFlush = function() {};
    // Worker side of the js/emit bridge — forward to main thread, which
    // dispatches into LetGoHost (workers have no DOM, no LetGoHost).
    globalThis._lgEmit = function(name, dataJson) {
      postMessage({t:'emit', name, data: dataJson});
    };
    // Host-eval (worker): the runtime sets globalThis._lgEval, then calls this to
    // announce readiness; relay to the main thread so it can resolve the gate.
    globalThis._lgRuntimeReady = function() {
      postMessage({t:'host-eval-ready'});
    };
    onmessage = async (e) => {
      if (e.data.t === 'request') {
        // Runs synchronously even while the program is parked on go.run: an
        // async onmessage doesn't block the worker, and the _lgRequest FuncOf
        // returns its value on the calling stack.
        const r = globalThis._lgRequest ? globalThis._lgRequest(JSON.stringify(e.data.req || {})) : '{"ok":false,"error":{"code":"not-ready","message":"runtime not ready"}}';
        postMessage({t:'request-result', id: e.data.id, result: r});
        return;
      }
      if (e.data.t !== 'init') return;
      const { sab, mode, wasmModule, wasmGzB64, wasmExecJS, urlSearch } = e.data;
      globalThis._lgKeyInt32 = new Int32Array(sab);
      globalThis._lgKeyUint8 = new Uint8Array(sab, 8, 16);
      // Page URL search string, forwarded from the main thread (the worker's
      // own location is the blob URL, not the page). Read by js/url-param.
      globalThis._lgUrlSearch = urlSearch || '';
      // Load wasm_exec.js in worker scope
      eval(wasmExecJS);
      const go = new Go();
      let instance;
      if (mode === 'external') {
        // External: the main thread compiled the payload via compileStreaming
        // and posted the WebAssembly.Module (structured-cloneable). The worker
        // only instantiates — no fetch, no recompile — so download+compile
        // overlapped worker spin-up. instantiate(Module, imports) resolves to
        // the Instance directly (the {instance,module} shape is BufferSource-only).
        instance = await WebAssembly.instantiate(wasmModule, go.importObject);
      } else {
        // Inline: decode the gzip-base64 payload posted from the main thread.
        const compressed = Uint8Array.from(atob(wasmGzB64), c => c.charCodeAt(0));
        const ds = new DecompressionStream('gzip');
        const w = ds.writable.getWriter(); w.write(compressed); w.close();
        const r = ds.readable.getReader();
        const chunks = []; let total = 0;
        while (true) { const {done,value} = await r.read(); if(done) break; chunks.push(value); total += value.length; }
        const wasmBytes = new Uint8Array(total);
        let off = 0; for(const c of chunks) { wasmBytes.set(c, off); off += c.length; }
        ({ instance } = await WebAssembly.instantiate(wasmBytes, go.importObject));
      }
      postMessage({t:'ready'});
      await go.run(instance);
      globalThis._lgFlush();
      postMessage({t:'exit'});
    };
  `;

  // External mode: kick off download + streaming compile on the main thread
  // *now*, so it overlaps worker creation. instantiateStreaming is MIME-strict,
  // so fall back to compile(arrayBuffer) if the host doesn't serve application/wasm.
  const modPromise = WASM_MODE === 'external' ? (async () => {
    try {
      return await WebAssembly.compileStreaming(fetch(WASM_URL));
    } catch (streamErr) {
      return await WebAssembly.compile(await (await fetch(WASM_URL)).arrayBuffer());
    }
  })() : null;

  const blob = new Blob([workerCode], { type: 'application/javascript' });
  const worker = new Worker(URL.createObjectURL(blob));

  worker.onmessage = (e) => {
    if (e.data.t === 'out') window.LetGoHost._output(e.data.d);
    if (e.data.t === 'exit') window.LetGoHost._output('\r\n\x1b[90m[program exited]\x1b[0m\r\n');
    if (e.data.t === 'emit') {
      try { window.LetGoHost._emit(e.data.name, JSON.parse(e.data.data)); }
      catch (err) { console.error('lg emit relay:', err); }
    }
    // Host-request relay: the worker runs the request and posts the JSON string
    // back, matched to its request id. requestImpl is installed here (not at
    // startup) so it only exists once the worker's runtime is actually ready.
    if (e.data.t === 'host-eval-ready') {
      requestImpl = (req) => new Promise((resolve, reject) => {
        const id = ++requestSeq;
        const timer = setTimeout(() => {
          if (requestPending.delete(id)) {
            reject(new Error('LetGoHost.request: no worker response within ' + REQUEST_TIMEOUT_MS + 'ms'));
          }
        }, REQUEST_TIMEOUT_MS);
        requestPending.set(id, { resolve, reject, timer });
        worker.postMessage({t:'request', id, req});
      });
      runtimeReadyResolve();
    }
    if (e.data.t === 'request-result') {
      const pending = requestPending.get(e.data.id);
      if (pending) {
        clearTimeout(pending.timer);
        requestPending.delete(e.data.id);
        try {
          pending.resolve(JSON.parse(e.data.result));
        } catch (err) {
          pending.reject(err);
        }
      }
    }
  };

  const initMsg = { t: 'init', sab, mode: WASM_MODE, wasmExecJS: WASM_EXEC_JS, urlSearch: location.search };
  if (WASM_MODE === 'external') {
    initMsg.wasmModule = await modPromise; // compiled on the main thread
  } else {
    initMsg.wasmGzB64 = WASM_GZ_B64;       // worker decodes the inline payload
  }
  worker.postMessage(initMsg);
  // A consumer now exists — accept input. Keys typed during the (possibly
  // multi-second, external-mode) load above were dropped, not queued.
  workerReady = true;
}

// --- Main-thread mode (output only, no input) ---
async function startMainThreadMode() {
  // Glue is wired (no SAB input in this mode). The shell opens its UI and
  // binds output via onOutput before the unavailable-input notice below.
  window.LetGoHost._ready('main');

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
    // App output now flows through _lgOutput (the Go-side HostWriter), not
    // fd 1/2 interception. This only catches Go-runtime direct writes
    // (e.g. panics) — surface those on the console.
    writeSync(fd, buf) {
      if (fd === 1 || fd === 2) { console.log(decoder.decode(buf)); return buf.length; }
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
  // Main-thread path runs on the page itself, so its location IS the page —
  // expose the search string for js/url-param (matches the worker path).
  globalThis._lgUrlSearch = location.search;
  // Main-thread side of the output bridge — the Go HostWriter calls
  // _lgOutput; feed it straight into LetGoHost (no worker round-trip).
  globalThis._lgOutput = function(s) {
    window.LetGoHost._output(s);
  };
  // Main-thread side of the js/emit bridge — dispatch straight into
  // LetGoHost (no worker round-trip needed).
  globalThis._lgEmit = function(name, dataJson) {
    try { window.LetGoHost._emit(name, JSON.parse(dataJson)); }
    catch (err) { console.error('lg emit:', err); }
  };
  // Host-request (main thread): the runtime sets window._lgRequest, then calls
  // this to announce readiness. Requests run in-page, so dispatch is a direct call.
  globalThis._lgRuntimeReady = function() {
    requestImpl = (req) => JSON.parse(window._lgRequest(JSON.stringify(req || {})));
    runtimeReadyResolve();
  };

  // Load wasm_exec.js
  eval(WASM_EXEC_JS);
  const go = new Go();
  let instance;
  if (WASM_MODE === 'external') {
    // Fetch the separate payload (streaming, with arrayBuffer fallback).
    try {
      ({ instance } = await WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject));
    } catch (streamErr) {
      const buf = await (await fetch(WASM_URL)).arrayBuffer();
      ({ instance } = await WebAssembly.instantiate(buf, go.importObject));
    }
  } else {
    // Decode the inline gzip-base64 payload.
    const wasmBytes = await decompressWasm(WASM_GZ_B64);
    ({ instance } = await WebAssembly.instantiate(wasmBytes, go.importObject));
  }
  go.run(instance);
}

// --- Entry point ---
(async () => {
  try {
    setStatus(WASM_MODE === 'external' ? 'fetching wasm...' : 'decompressing wasm...');
    if (typeof SharedArrayBuffer !== 'undefined' && crossOriginIsolated) {
      // await so external-mode fetch/compile rejections (after the first await
      // in startWorkerMode) surface through the catch below, not as an
      // unhandled rejection that strands the page at "fetching wasm...".
      await startWorkerMode();
    } else {
      await startMainThreadMode();
    }
  } catch(err) {
    setStatus('error: ' + err);
    console.error(err);
  }
})();
