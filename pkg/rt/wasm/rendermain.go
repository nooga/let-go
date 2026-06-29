/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"strconv"
	"strings"
)

// wasmMainTmpl is the generated main.go for a `lg -w` bundle: it decodes the
// embedded program.lgb, runs each namespace chunk, then the main chunk. The
// __LG_* markers are filled by RenderMain.
const wasmMainTmpl = `package main

import (
	_ "embed"
	"bytes"
	"fmt"
	"os"
__LG_HOST_EVAL_IMPORTS__

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

	// Route *out*/*err* to the JS host via _lgOutput (HostWriter), instead of
	// os.Stdout/Stderr + the bundle's fs.writeSync fd interception. SetRoot,
	// not a per-Run binding, because this generated main drives bytecode
	// directly rather than through pkg/api. Guarded: if the core I/O vars
	// aren't installed yet, output falls back to os.Stdout.
	hostWriter := rt.NewHostWriter()
	if v := rt.LookupCoreVar("*out*"); v != nil {
		v.SetRoot(vm.NewBoxed(rt.NewWriterHandle("host-stdout", hostWriter)))
	}
	if v := rt.LookupCoreVar("*err*"); v != nil {
		v.SetRoot(vm.NewBoxed(rt.NewWriterHandle("host-stderr", hostWriter)))
	}

	// Route (js/emit ...) to the JS host via _lgEmit (HostEmitter), the dual
	// of the HostWriter *out* routing above. Same SetRoot rationale.
	hostEmitter := rt.NewHostEmitter()
	if v := rt.LookupCoreVar("*emit*"); v != nil {
		v.SetRoot(vm.NewBoxed(hostEmitter))
	}

	// Route storage through browser localStorage, scoped by the bundle's
	// host-selected store id so guest keys remain app-local.
	hostStorage := rt.NewHostStorage(__LG_STORAGE_ID__)
	if v := rt.LookupCoreVar("*storage*"); v != nil {
		v.SetRoot(vm.NewBoxed(hostStorage))
	}

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
__LG_HOST_EVAL_BODY__
}
`

// wasmHostEvalSnippet is spliced into wasmMainTmpl at __LG_HOST_EVAL_BODY__ when
// -w-host-eval is set. After the program's main chunk runs, it installs an
// internal _lgEval hook — compile + run a string in the loaded image, returning
// a stringified value (FormatError on failure) — plus a structured _lgRequest
// hook that currently supports the `eval` op and normalizes an implicit
// `default` session. The JS host wraps _lgRequest as the public
// LetGoHost.request(req) API and keeps LetGoHost.eval(code) as a compatibility
// helper. After wiring the hooks, the runtime signals readiness and parks so it
// stays callable.
//
// _lgEval is the internal hook (like _lgKey / _lgEmit); lg-host-core.js wraps it
// as the public LetGoHost.eval(code), calling it directly on the main thread or
// relaying through the worker in cross-origin-isolated mode, so one API works in
// both boot modes. The host installs _lgRuntimeReady; calling it resolves the
// LetGoHost.eval ready gate (main thread) or posts the ready message (worker),
// closing the race where a client could call eval before the runtime is up.
// Structured request responses return as JSON strings; richer data still flows
// out-of-band via (js/emit ...).
const wasmHostEvalSnippet = `	type hostRequest struct {
		ID      string ` + "`json:\"id,omitempty\"`" + `
		Session string ` + "`json:\"session,omitempty\"`" + `
		Op      string ` + "`json:\"op\"`" + `
		NS      string ` + "`json:\"ns,omitempty\"`" + `
		Code    string ` + "`json:\"code,omitempty\"`" + `
	}
	type hostError struct {
		Code    string ` + "`json:\"code\"`" + `
		Message string ` + "`json:\"message\"`" + `
	}
	type hostResponse struct {
		ID       string      ` + "`json:\"id,omitempty\"`" + `
		OK       bool        ` + "`json:\"ok\"`" + `
		Op       string      ` + "`json:\"op,omitempty\"`" + `
		Session  string      ` + "`json:\"session,omitempty\"`" + `
		Value    any         ` + "`json:\"value,omitempty\"`" + `
		Error    *hostError  ` + "`json:\"error,omitempty\"`" + `
		Warnings []string    ` + "`json:\"warnings\"`" + `
		Output   []string    ` + "`json:\"output\"`" + `
		Metrics  interface{} ` + "`json:\"metrics\"`" + `
		Artifacts interface{} ` + "`json:\"artifacts\"`" + `
	}
	type hostSession struct {
		NS string
	}
	sessions := map[string]*hostSession{
		"default": {NS: "user"},
	}
	getSession := func(id string) *hostSession {
		if id == "" {
			id = "default"
		}
		if s := sessions[id]; s != nil {
			return s
		}
		s := &hostSession{NS: "user"}
		sessions[id] = s
		return s
	}
	opcodeStride := func(op int32) int {
		switch op & 0xff {
		case vm.OP_TRY_PUSH:
			return 3
		case vm.OP_RECUR:
			return 4
		case vm.OP_LOAD_ARG, vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE, vm.OP_JUMP,
			vm.OP_POP_N, vm.OP_DUP_NTH, vm.OP_INVOKE, vm.OP_LOAD_CLOSEDOVER,
			vm.OP_RECUR_FN, vm.OP_MAKE_MULTI_ARITY, vm.OP_TAIL_CALL,
			vm.OP_LOAD_CONST, vm.OP_LOAD_VAR:
			return 2
		default:
			return 1
		}
	}
	projectConst := func(v vm.Value) any {
		switch x := v.(type) {
		case vm.Symbol:
			return string(x)
		case vm.Keyword:
			return ":" + string(x)
		case vm.String:
			return string(x)
		case vm.Int:
			return int(x)
		case vm.Float:
			return float64(x)
		case vm.Boolean:
			return bool(x)
		default:
			return v.String()
		}
	}
	disassembleChunk := func(chunk *vm.CodeChunk) []map[string]any {
		code := chunk.Code()
		consts := chunk.Consts().AllValues()
		rows := make([]map[string]any, 0, len(code))
		for i := 0; i < len(code); {
			op := code[i]
			stride := opcodeStride(op)
			args := make([]int, 0, stride-1)
			row := map[string]any{
				"offset": i,
				"op":     vm.OpcodeToString(op),
				"args":   args,
			}
			for j := 1; j < stride && i+j < len(code); j++ {
				args = append(args, int(code[i+j]))
			}
			row["args"] = args
			switch op & 0xff {
			case vm.OP_LOAD_CONST, vm.OP_LOAD_VAR:
				if i+1 < len(code) {
					if idx := int(code[i+1]); idx >= 0 && idx < len(consts) {
						row["resolved"] = projectConst(consts[idx])
					}
				}
			}
			rows = append(rows, row)
			i += stride
		}
		return rows
	}
	compileForms := func(req hostRequest) (*vm.CodeChunk, vm.Value, *hostSession, error) {
		sess := getSession(req.Session)
		if req.NS != "" {
			sess.NS = req.NS
		}
		ns := rt.NS(sess.NS)
		c := compiler.NewCompiler(consts, ns)
		c.SetSource("<host-request>")
		chunk, result, err := c.CompileMultiple(bytes.NewReader([]byte(req.Code)))
		sess.NS = c.CurrentNS().Name()
		return chunk, result, sess, err
	}
	evalCode := func(code string) (string, error) {
		chunk, cerr := ctx.Compile(code)
		if cerr != nil {
			return "", cerr
		}
		frame := vm.NewFrame(chunk, nil)
		result, rerr := frame.RunProtected()
		vm.ReleaseFrame(frame)
		if rerr != nil {
			return "", rerr
		}
		return result.String(), nil
	}
	hostEval := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return "error: eval expects one string argument"
		}
		result, err := evalCode(args[0].String())
		if err != nil {
			return vm.FormatError(err)
		}
		return result
	})
	hostRequestFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		resp := hostResponse{
			Warnings: []string{},
			Output:   []string{},
		}
		if len(args) < 1 {
			resp.OK = false
			resp.Error = &hostError{Code: "bad-request", Message: "request expects one JSON string argument"}
			b, _ := json.Marshal(resp)
			return string(b)
		}
		var req hostRequest
		if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
			resp.OK = false
			resp.Error = &hostError{Code: "bad-request", Message: err.Error()}
			b, _ := json.Marshal(resp)
			return string(b)
		}
		resp.ID = req.ID
		resp.Op = req.Op
		if req.Session == "" {
			req.Session = "default"
		}
		resp.Session = req.Session
		switch req.Op {
		case "eval":
			if req.Code == "" {
				resp.OK = false
				resp.Error = &hostError{Code: "bad-request", Message: "eval requires code"}
			} else if _, result, sess, err := compileForms(req); err != nil {
				resp.OK = false
				resp.Error = &hostError{Code: "eval-error", Message: vm.FormatError(err)}
			} else {
				resp.OK = true
				resp.Session = req.Session
				resp.Value = result.String()
				resp.Artifacts = map[string]any{
					"namespace": sess.NS,
				}
			}
		case "compile":
			if req.Code == "" {
				resp.OK = false
				resp.Error = &hostError{Code: "bad-request", Message: "compile requires code"}
			} else if chunk, result, sess, err := compileForms(req); err != nil {
				resp.OK = false
				resp.Error = &hostError{Code: "compile-error", Message: vm.FormatError(err)}
			} else {
				resp.OK = true
				resp.Value = result.String()
				resp.Artifacts = map[string]any{
					"namespace":  sess.NS,
					"chunkWords": len(chunk.Code()),
					"maxStack":   chunk.MaxStack(),
				}
			}
		case "disassemble":
			if req.Code == "" {
				resp.OK = false
				resp.Error = &hostError{Code: "bad-request", Message: "disassemble requires code"}
			} else if chunk, result, sess, err := compileForms(req); err != nil {
				resp.OK = false
				resp.Error = &hostError{Code: "compile-error", Message: vm.FormatError(err)}
			} else {
				resp.OK = true
				resp.Value = result.String()
				resp.Artifacts = map[string]any{
					"namespace": sess.NS,
					"bytecode": map[string]any{
						"status": "ok",
						"rows":   disassembleChunk(chunk),
					},
				}
			}
		case "inspect-all":
			if req.Code == "" {
				resp.OK = false
				resp.Error = &hostError{Code: "bad-request", Message: "inspect-all requires code"}
			} else if chunk, result, sess, err := compileForms(req); err != nil {
				resp.OK = false
				resp.Error = &hostError{Code: "compile-error", Message: vm.FormatError(err)}
			} else {
				resp.OK = true
				resp.Value = result.String()
				resp.Artifacts = map[string]any{
					"namespace": sess.NS,
					"result": map[string]any{
						"status": "ok",
						"value":  result.String(),
					},
					"bytecode": map[string]any{
						"status": "ok",
						"rows":   disassembleChunk(chunk),
					},
					"ir": map[string]any{
						"status": "unavailable",
						"reason": "ir-build is not implemented on the browser bridge yet",
					},
					"optimizedBytecode": map[string]any{
						"status": "unavailable",
						"reason": "IR bytecode lowering is not implemented on the browser bridge yet",
					},
					"loweredGo": map[string]any{
						"status": "unavailable",
						"reason": "Go lowering is not implemented on the browser bridge yet",
					},
				}
			}
		default:
			resp.OK = false
			resp.Error = &hostError{Code: "unknown-op", Message: "unsupported op: " + req.Op}
		}
		b, _ := json.Marshal(resp)
		return string(b)
	})
	js.Global().Set("_lgEval", hostEval)
	js.Global().Set("_lgRequest", hostRequestFn)
	if ready := js.Global().Get("_lgRuntimeReady"); ready.Type() == js.TypeFunction {
		ready.Invoke()
	}
	select {}`

// RenderMain fills wasmMainTmpl's placeholders: the storage-id, and the
// -w-host-eval splice (the window.Eval bridge + park, plus its syscall/js
// import). With hostEval false the marker lines are removed whole, so the
// default bundle's generated main is byte-identical to the pre-flag output.
func RenderMain(storeID string, hostEval bool) string {
	s := strings.ReplaceAll(wasmMainTmpl, "__LG_STORAGE_ID__", strconv.Quote(storeID))
	if hostEval {
		s = strings.ReplaceAll(s, "__LG_HOST_EVAL_IMPORTS__", "\t\"encoding/json\"\n\t\"syscall/js\"")
		s = strings.ReplaceAll(s, "__LG_HOST_EVAL_BODY__", wasmHostEvalSnippet)
	} else {
		s = strings.ReplaceAll(s, "__LG_HOST_EVAL_IMPORTS__\n", "")
		s = strings.ReplaceAll(s, "__LG_HOST_EVAL_BODY__\n", "")
	}
	return s
}
