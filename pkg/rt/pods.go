/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/bencode"
)

// podMsg is a decoded bencode message from a pod.
type podMsg struct {
	raw map[string]interface{}
}

// Pod represents a running babashka pod subprocess.
type Pod struct {
	id         string // identifier derived from first namespace
	name       string
	format     string // "json", "edn", or "transit+json"
	cmd        *exec.Cmd
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	stderr     io.ReadCloser
	decoder    *bencode.Decoder
	writeMu    sync.Mutex // protects stdin writes
	namespaces []podNamespace
	ops        map[string]interface{}
	idCounter  atomic.Int64
	shutdown   bool

	// Response routing
	pending   map[string]chan podMsg // id -> response channel
	streaming map[string]bool       // ids that should not auto-close on done
	pendingMu sync.Mutex
	routerUp  bool
}

type podNamespace struct {
	name   string
	vars   []podVar
	defer_ bool
}

type podVar struct {
	name    string
	code    string // client-side code (macros, wrapper fns)
	meta    string
	argMeta bool
}

// --- Pod registry ---

var (
	podRegistry   map[string]*Pod // id -> pod
	podRegistryMu sync.Mutex
)

func init() {
	podRegistry = make(map[string]*Pod)
}

func registerPod(p *Pod) {
	podRegistryMu.Lock()
	podRegistry[p.id] = p
	podRegistryMu.Unlock()
}

func lookupPod(id string) *Pod {
	podRegistryMu.Lock()
	defer podRegistryMu.Unlock()
	return podRegistry[id]
}

// ShutdownAllPods gracefully shuts down all loaded pods.
func ShutdownAllPods() {
	podRegistryMu.Lock()
	defer podRegistryMu.Unlock()
	for _, p := range podRegistry {
		p.Shutdown()
	}
	podRegistry = make(map[string]*Pod)
}

// --- Pod protocol ---

func (p *Pod) nextID() string {
	return fmt.Sprintf("lg-%d", p.idCounter.Add(1))
}

func (p *Pod) send(msg map[string]interface{}) error {
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	bs, err := bencode.EncodeBytes(msg)
	if err != nil {
		return fmt.Errorf("pod %s: bencode encode: %w", p.name, err)
	}
	_, err = p.stdin.Write(bs)
	return err
}

// recvRaw reads a single message (only used during describe before router starts).
func (p *Pod) recvRaw() (map[string]interface{}, error) {
	var msg map[string]interface{}
	if err := p.decoder.Decode(&msg); err != nil {
		return nil, fmt.Errorf("pod %s: bencode decode: %w", p.name, err)
	}
	return msg, nil
}

// startRouter starts the background message router goroutine.
func (p *Pod) startRouter() {
	p.pending = make(map[string]chan podMsg)
	p.streaming = make(map[string]bool)
	p.routerUp = true
	go func() {
		for {
			var msg map[string]interface{}
			if err := p.decoder.Decode(&msg); err != nil {
				// Pod closed or errored - close all pending channels
				p.pendingMu.Lock()
				for _, ch := range p.pending {
					close(ch)
				}
				p.pending = make(map[string]chan podMsg)
				p.pendingMu.Unlock()
				return
			}

			// Handle out/err on any message
			if out, ok := msg["out"].(string); ok {
				fmt.Print(out)
			}
			if errStr, ok := msg["err"].(string); ok {
				fmt.Fprint(os.Stderr, errStr)
			}

			// Route by id
			id := bencStr(msg, "id")
if id == "" {
				continue
			}

			p.pendingMu.Lock()
			ch, ok := p.pending[id]
			p.pendingMu.Unlock()

			if ok {
				ch <- podMsg{raw: msg}
				// Clean up if done (but not for streaming channels)
				if isDone(msg) {
					p.pendingMu.Lock()
					isStreaming := p.streaming[id]
					if !isStreaming {
						delete(p.pending, id)
						p.pendingMu.Unlock()
						close(ch)
					} else {
						p.pendingMu.Unlock()
					}
				}
			}
		}
	}()
}

// registerPending creates a channel for receiving responses with the given id.
func (p *Pod) registerPending(id string) chan podMsg {
	ch := make(chan podMsg, 16)
	p.pendingMu.Lock()
	p.pending[id] = ch
	p.pendingMu.Unlock()
	return ch
}

// registerStreaming creates a channel that stays open after "done" (for async streaming).
func (p *Pod) registerStreaming(id string) chan podMsg {
	ch := make(chan podMsg, 64)
	p.pendingMu.Lock()
	p.pending[id] = ch
	p.streaming[id] = true
	p.pendingMu.Unlock()
	return ch
}

func (p *Pod) describe() error {
	if err := p.send(map[string]interface{}{"op": "describe"}); err != nil {
		return err
	}
	msg, err := p.recvRaw()
	if err != nil {
		return err
	}

	if f, ok := msg["format"].(string); ok {
		p.format = f
	} else {
		p.format = "json"
	}

	if ops, ok := msg["ops"].(map[string]interface{}); ok {
		p.ops = ops
	}

	if nsList, ok := msg["namespaces"].([]interface{}); ok {
		for _, nsRaw := range nsList {
			nsMap, ok := nsRaw.(map[string]interface{})
			if !ok {
				continue
			}
			pns := podNamespace{name: bencStr(nsMap, "name")}
			if d, _ := nsMap["defer"].(string); d == "true" {
				pns.defer_ = true
			}
			if varsList, ok := nsMap["vars"].([]interface{}); ok {
				for _, varRaw := range varsList {
					varMap, ok := varRaw.(map[string]interface{})
					if !ok {
						continue
					}
					pv := podVar{
						name: bencStr(varMap, "name"),
						code: bencStr(varMap, "code"),
						meta: bencStr(varMap, "meta"),
					}
					if am, _ := varMap["arg-meta"].(string); am == "true" {
						pv.argMeta = true
					}
					pns.vars = append(pns.vars, pv)
				}
			}
			p.namespaces = append(p.namespaces, pns)
		}
	}

	if len(p.namespaces) > 0 {
		p.id = p.namespaces[0].name
	}

	return nil
}

// Invoke calls a pod var and returns the result (synchronous).
func (p *Pod) Invoke(varName string, args []vm.Value) (vm.Value, error) {
	encoded, err := p.encodeArgs(args)
	if err != nil {
		return vm.NIL, fmt.Errorf("pod %s: encode args: %w", p.name, err)
	}

	id := p.nextID()
	ch := p.registerPending(id)

	if err := p.send(map[string]interface{}{
		"op":   "invoke",
		"id":   id,
		"var":  varName,
		"args": encoded,
	}); err != nil {
		return vm.NIL, err
	}

	// Collect responses from the router
	var result vm.Value = vm.NIL
	for msg := range ch {
		m := msg.raw
		if exMsg, ok := m["ex-message"].(string); ok {
			exData := bencStr(m, "ex-data")
			var dataVal vm.Value = vm.NIL
			if exData != "" {
				dataVal, _ = p.decodePayload(exData)
			}
			errMap := vm.EmptyPersistentMap
			errMap = errMap.Assoc(vm.Keyword("message"), vm.String(exMsg)).(*vm.PersistentMap)
			if dataVal != vm.NIL {
				errMap = errMap.Assoc(vm.Keyword("data"), dataVal).(*vm.PersistentMap)
			}
			return vm.NIL, vm.NewThrownError(errMap)
		}
		if valStr, ok := m["value"].(string); ok {
			result, err = p.decodePayload(valStr)
			if err != nil {
				return vm.NIL, fmt.Errorf("pod %s: decode value: %w", p.name, err)
			}
		}
	}
	return result, nil
}

// InvokeAsync calls a pod var and sends results to a channel.
// The channel stays open for streaming - the pod sends multiple values with the same id.
func (p *Pod) InvokeAsync(varName string, args []vm.Value) (vm.Chan, error) {
	encoded, err := p.encodeArgs(args)
	if err != nil {
		return nil, fmt.Errorf("pod %s: encode args: %w", p.name, err)
	}

	id := p.nextID()
	// Register as streaming (don't auto-close on done)
	routerCh := p.registerStreaming(id)

	if err := p.send(map[string]interface{}{
		"op":   "invoke",
		"id":   id,
		"var":  varName,
		"args": encoded,
	}); err != nil {
		return nil, err
	}

	valCh := make(vm.Chan, 32)
	go func() {
		defer close(valCh)
		for msg := range routerCh {
			m := msg.raw
			if exMsg, ok := m["ex-message"].(string); ok {
				errMap := vm.EmptyPersistentMap
				errMap = errMap.Assoc(vm.Keyword("error"), vm.String(exMsg)).(*vm.PersistentMap)
				if exData := bencStr(m, "ex-data"); exData != "" {
					if dv, err := p.decodePayload(exData); err == nil {
						errMap = errMap.Assoc(vm.Keyword("data"), dv).(*vm.PersistentMap)
					}
				}
				valCh <- errMap
				return
			}
			if valStr, ok := m["value"].(string); ok {
				if v, err := p.decodePayload(valStr); err == nil {
					valCh <- v
				}
			}
		}
	}()

	return valCh, nil
}

// Shutdown sends shutdown op if supported and waits for process to exit.
func (p *Pod) Shutdown() {
	if p.shutdown {
		return
	}
	p.shutdown = true
	if p.ops != nil {
		if _, ok := p.ops["shutdown"]; ok {
			_ = p.send(map[string]interface{}{"op": "shutdown"})
		}
	}
	p.stdin.Close()
	p.cmd.Wait() //nolint:errcheck
}

// --- Payload encoding/decoding ---

func (p *Pod) encodeArgs(args []vm.Value) (string, error) {
	switch p.format {
	case "json":
		return JSONEncodeArgs(args)
	case "transit+json":
		return TransitEncodeArgs(args)
	case "edn":
		return EDNEncodeArgs(args)
	default:
		return "", fmt.Errorf("unsupported pod format: %s", p.format)
	}
}

func (p *Pod) decodePayload(s string) (vm.Value, error) {
	switch p.format {
	case "json":
		return JSONDecodeValue(s)
	case "transit+json":
		return TransitDecodeValue(s)
	case "edn":
		if readEDN == nil {
			return vm.NIL, fmt.Errorf("EDN reader not initialized")
		}
		return readEDN(s)
	default:
		return vm.NIL, fmt.Errorf("unsupported pod format: %s", p.format)
	}
}

// --- Namespace proxy creation ---

func createProxyNamespaces(p *Pod) error {
	for _, pns := range p.namespaces {
		ns := LookupOrRegisterNSNoLoad(pns.name)

		for _, pv := range pns.vars {
			if pv.code != "" {
				// Client-side code: evaluate in the pod's namespace
				if err := evalPodCode(pv.code, ns); err != nil {
					// Non-fatal: log and continue (the var may still be usable)
					fmt.Fprintf(os.Stderr, "pod %s: client code eval error for %s/%s: %v\n",
						p.id, pns.name, pv.name, err)
				}
				continue
			}

			// Create proxy function
			qualifiedName := pns.name + "/" + pv.name
			ns.Def(pv.name, makePodProxy(p, qualifiedName))
		}
	}
	return nil
}

func evalPodCode(code string, ns *vm.Namespace) error {
	if evalInNS == nil {
		return fmt.Errorf("evalInNS not initialized")
	}
	_, err := evalInNS(code, ns)
	return err
}

func makePodProxy(p *Pod, qualifiedName string) vm.Value {
	fn, _ := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
		return p.Invoke(qualifiedName, args)
	})
	return fn
}

// --- Pod resolution ---

func resolvePodBinary(name string, version string) (string, error) {
	if version != "" {
		return findCachedPod(name, version)
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("pod binary %q not found on PATH", name)
	}
	return path, nil
}

func findCachedPod(name string, version string) (string, error) {
	podsDir := os.Getenv("BABASHKA_PODS_DIR")
	if podsDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		podsDir = filepath.Join(home, ".babashka", "pods")
	}

	osName := runtime.GOOS
	arch := runtime.GOARCH
	switch osName {
	case "darwin":
		osName = "mac_os_x"
	}
	switch arch {
	case "arm64":
		arch = "aarch64"
	case "amd64":
		arch = "x86_64"
	}

	repoDir := filepath.Join(podsDir, "repository", name, version, osName, arch)
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return "", fmt.Errorf("pod %s@%s not found in cache (%s) - install with: bb -e '(pods/load-pod (quote %s) \"%s\")'",
			name, version, repoDir, name, version)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return filepath.Join(repoDir, entry.Name()), nil
		}
	}
	return "", fmt.Errorf("pod %s@%s: no executable in %s", name, version, repoDir)
}

func startPod(binary string) (*Pod, error) {
	cmd := exec.Command(binary)
	cmd.Env = append(os.Environ(), "BABASHKA_POD=true")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start pod %s: %w", binary, err)
	}

	go func() {
		io.Copy(os.Stderr, stderr) //nolint:errcheck
	}()

	p := &Pod{
		name:    binary,
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		decoder: bencode.NewDecoder(stdout),
	}

	if err := p.describe(); err != nil {
		cmd.Process.Kill() //nolint:errcheck
		return nil, fmt.Errorf("pod describe failed: %w", err)
	}

	// Start background message router after describe handshake
	p.startRouter()

	return p, nil
}

// --- Helpers ---

func isDone(msg map[string]interface{}) bool {
	statusRaw, ok := msg["status"]
	if !ok {
		return false
	}
	switch s := statusRaw.(type) {
	case string:
		return s == "done" || s == `["done"]`
	case []interface{}:
		for _, v := range s {
			if str, ok := v.(string); ok && str == "done" {
				return true
			}
		}
	}
	return false
}

func bencStr(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// --- Namespace ---

func installPodsNS() {
	loadPod, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("load-pod: expected 1-2 arguments (name [version])")
		}

		var name, version string
		switch v := vs[0].(type) {
		case vm.String:
			name = string(v)
		case vm.Symbol:
			name = string(v)
		default:
			return vm.NIL, fmt.Errorf("load-pod: expected string or symbol, got %s", vs[0].Type().Name())
		}
		if len(vs) == 2 {
			v, ok := vs[1].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("load-pod: version must be a string")
			}
			version = string(v)
		}

		binary, err := resolvePodBinary(name, version)
		if err != nil {
			return vm.NIL, err
		}

		pod, err := startPod(binary)
		if err != nil {
			return vm.NIL, err
		}

		registerPod(pod)
		if err := createProxyNamespaces(pod); err != nil {
			return vm.NIL, err
		}

		if len(pod.namespaces) > 0 {
			return vm.String(pod.namespaces[0].name), nil
		}
		return vm.NIL, nil
	})

	// pods/invoke - low-level invoke for async use
	// (pods/invoke pod-id 'var-sym [args] {:handlers {:success fn :error fn :done fn}})
	invokePod, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 3 {
			return vm.NIL, fmt.Errorf("pods/invoke: expected (pod-id var-sym args [opts])")
		}
		podID, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("pods/invoke: pod-id must be a string, got %s (%s)", vs[0].Type().Name(), vs[0])
		}
		pod := lookupPod(string(podID))
		if pod == nil {
			return vm.NIL, fmt.Errorf("pods/invoke: no pod with id %q", podID)
		}
		varSym, ok := vs[1].(vm.Symbol)
		if !ok {
			return vm.NIL, fmt.Errorf("pods/invoke: var must be a symbol")
		}

		// Extract args
		var args []vm.Value
		if sq, ok := vs[2].(vm.Sequable); ok {
			for s := sq.Seq(); s != nil; s = s.Next() {
				args = append(args, s.First())
			}
		}

		// Check for opts with :handlers
		if len(vs) >= 4 {
			if opts, ok := vs[3].(vm.Lookup); ok {
				handlers := opts.ValueAt(vm.Keyword("handlers"))
				if handlers != vm.NIL {
					// Async mode: invoke and dispatch to handler callbacks
					ch, err := pod.InvokeAsync(string(varSym), args)
					if err != nil {
						return vm.NIL, err
					}

					var successFn, errorFn, doneFn vm.Value
					if hmap, ok := handlers.(vm.Lookup); ok {
						successFn = hmap.ValueAt(vm.Keyword("success"))
						errorFn = hmap.ValueAt(vm.Keyword("error"))
						doneFn = hmap.ValueAt(vm.Keyword("done"))
						}

					go func() {
						for val := range ch {
							// Check if it's an error map
							if m, ok := val.(*vm.PersistentMap); ok {
								if errMsg := m.ValueAt(vm.Keyword("error")); errMsg != vm.NIL {
									if errorFn != nil && errorFn != vm.NIL {
										callFn(errorFn, []vm.Value{val})
									}
									continue
								}
							}
							if successFn != nil && successFn != vm.NIL {
								callFn(successFn, []vm.Value{val})
							}
						}
						if doneFn != nil && doneFn != vm.NIL {
							callFn(doneFn, nil)
						}
					}()
					return vm.NIL, nil
				}
			}
		}

		// Synchronous mode
		return pod.Invoke(string(varSym), args)
	})

	ns := vm.NewNamespace("pods")
	ns.Def("load-pod", loadPod)
	ns.Def("invoke", invokePod)
	RegisterNS(ns)

	// Alias as babashka.pods for compatibility with pod client-side code
	bbns := vm.NewNamespace("babashka.pods")
	bbns.Def("load-pod", loadPod)
	bbns.Def("invoke", invokePod)
	RegisterNS(bbns)
}

// invokable is anything with an Invoke method (Func, NativeFn, MultiArityFn, Closure, etc).
type invokable interface {
	Invoke(args []vm.Value) (vm.Value, error)
}

// callFn safely calls a vm function value, recovering from panics.
func callFn(fn vm.Value, args []vm.Value) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "pod callback panic: %v\n", r)
		}
	}()
	if f, ok := fn.(invokable); ok {
		f.Invoke(args) //nolint:errcheck
	}
}
