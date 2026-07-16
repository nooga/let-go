package wasmhost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type InspectOptions struct {
	Bytecode          bool `json:"bytecode,omitempty"`
	IR                bool `json:"ir,omitempty"`
	OptimizedBytecode bool `json:"optimized_bytecode,omitempty"`
	LoweredGo         bool `json:"lowered_go,omitempty"`
}

type Request struct {
	ID      string         `json:"id,omitempty"`
	Session string         `json:"session,omitempty"`
	Op      string         `json:"op"`
	NS      string         `json:"ns,omitempty"`
	Code    string         `json:"code,omitempty"`
	Inspect InspectOptions `json:"inspect,omitempty"`
}

type Frame struct {
	ID          string   `json:"id,omitempty"`
	Session     string   `json:"session,omitempty"`
	Op          string   `json:"op,omitempty"`
	Kind        string   `json:"kind"`
	FormIndex   int      `json:"form_index,omitempty"`
	Stage       string   `json:"stage,omitempty"`
	Artifact    string   `json:"artifact,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	Content     any      `json:"content,omitempty"`
	Text        string   `json:"text,omitempty"`
	Value       string   `json:"value,omitempty"`
	Status      []string `json:"status,omitempty"`
}

type session struct {
	NS string
}

type Host struct {
	consts   *vm.Consts
	sessions map[string]*session
	mu       sync.Mutex
}

func New(consts *vm.Consts) *Host {
	return &Host{
		consts: consts,
		sessions: map[string]*session{
			"default": {NS: "user"},
		},
	}
}

func DecodeRequest(bs []byte) (Request, error) {
	var req Request
	err := json.Unmarshal(bs, &req)
	return req, err
}

func EncodeFrame(frame Frame) ([]byte, error) {
	return json.Marshal(frame)
}

func (h *Host) HandleJSON(reqJSON string, emit func(string)) string {
	req, err := DecodeRequest([]byte(reqJSON))
	if err != nil {
		h.emit(emit, Frame{
			Kind:   "err",
			Text:   err.Error(),
			Status: nil,
		})
		h.emit(emit, Frame{Kind: "status", Status: []string{"done", "error"}})
		return `{"ok":false}`
	}
	h.Handle(req, emit)
	return `{"ok":true}`
}

func (h *Host) Handle(req Request, emit func(string)) {
	if req.Session == "" {
		req.Session = "default"
	}
	switch req.Op {
	case "eval":
		h.handleEval(req, emit)
	default:
		h.emit(emit, Frame{
			ID:      req.ID,
			Session: req.Session,
			Op:      req.Op,
			Kind:    "err",
			Text:    "unsupported op: " + req.Op,
		})
		h.emit(emit, Frame{
			ID:      req.ID,
			Session: req.Session,
			Op:      req.Op,
			Kind:    "status",
			Status:  []string{"done", "error", "unknown-op"},
		})
	}
}

func (h *Host) handleEval(req Request, emit func(string)) {
	if strings.TrimSpace(req.Code) == "" {
		h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "err", Text: "eval requires code"})
		h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "status", Status: []string{"done", "error", "bad-request"}})
		return
	}
	forms, err := compiler.SplitTopLevelForms(req.Code, "<host-request>")
	if err != nil {
		h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "err", Text: vm.FormatError(err)})
		h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "status", Status: []string{"done", "error"}})
		return
	}
	sess := h.getSession(req.Session)
	for i, form := range forms {
		formIndex := i + 1
		h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "progress", FormIndex: formIndex, Stage: "compile"})

		var (
			chunk  *vm.CodeChunk
			result vm.Value
			outBuf bytes.Buffer
			errBuf bytes.Buffer
			runErr error
		)
		func() {
			ns := rt.NS(sess.NS)
			// Per-request transient compiler: the host is long-lived, request
			// evals must not root its pool (see compiler.NewTransientCompiler).
			c := compiler.NewTransientCompiler(h.consts, ns)
			c.SetSource("<host-request>")
			outVar := rt.LookupCoreVar("*out*")
			if outVar != nil {
				outVar.PushBinding(vm.NewBoxed(rt.NewWriterHandle("wasmhost-out", &outBuf)))
				defer outVar.PopBinding()
			}
			errVar := rt.LookupCoreVar("*err*")
			if errVar != nil {
				errVar.PushBinding(vm.NewBoxed(rt.NewWriterHandle("wasmhost-err", &errBuf)))
				defer errVar.PopBinding()
			}
			chunk, result, runErr = c.CompileMultiple(bytes.NewReader([]byte(form.Source)))
			sess.NS = c.CurrentNS().Name()
		}()

		inspectable := isInspectableIRForm(form.Form)
		if req.Inspect.Bytecode {
			// Prefer the fn-body chunk (defn -> var -> *Func -> chunk); fall
			// back to the top-level form chunk for non-fn forms.
			bcChunk := chunk
			if fnChunk := codeChunkFromValue(result); fnChunk != nil {
				bcChunk = fnChunk
			}
			if bcChunk != nil {
				h.emit(emit, Frame{
					ID:          req.ID,
					Session:     req.Session,
					Op:          req.Op,
					Kind:        "artifact",
					FormIndex:   formIndex,
					Stage:       "bytecode",
					Artifact:    "bytecode",
					ContentType: "application/json",
					Content:     disassembleChunk(bcChunk),
				})
			}
		}
		if !inspectable && req.Inspect.IR {
			h.emit(emit, Frame{
				ID:        req.ID,
				Session:   req.Session,
				Op:        req.Op,
				Kind:      "note",
				FormIndex: formIndex,
				Stage:     "ir",
				Artifact:  "ir",
				Text:      "IR extraction unsupported for this form",
			})
		}
		if inspectable && req.Inspect.IR {
			if irDump, err := h.inspectIR(sess.NS, form.Source); err == nil && irDump != "" {
				h.emit(emit, Frame{
					ID:          req.ID,
					Session:     req.Session,
					Op:          req.Op,
					Kind:        "artifact",
					FormIndex:   formIndex,
					Stage:       "ir",
					Artifact:    "ir",
					ContentType: "text/plain",
					Content:     irDump,
				})
			} else if err != nil {
				h.emit(emit, Frame{
					ID:        req.ID,
					Session:   req.Session,
					Op:        req.Op,
					Kind:      "err",
					FormIndex: formIndex,
					Stage:     "ir",
					Artifact:  "ir",
					Text:      err.Error(),
				})
			}
		}
		if !inspectable && req.Inspect.OptimizedBytecode {
			h.emit(emit, Frame{
				ID:        req.ID,
				Session:   req.Session,
				Op:        req.Op,
				Kind:      "note",
				FormIndex: formIndex,
				Stage:     "optimized-bytecode",
				Artifact:  "optimized_bytecode",
				Text:      "Optimized bytecode extraction unsupported for this form",
			})
		}
		if inspectable && req.Inspect.OptimizedBytecode {
			if optimized, err := h.inspectOptimizedBytecode(sess.NS, form.Source); err == nil && optimized != nil {
				h.emit(emit, Frame{
					ID:          req.ID,
					Session:     req.Session,
					Op:          req.Op,
					Kind:        "artifact",
					FormIndex:   formIndex,
					Stage:       "optimized-bytecode",
					Artifact:    "optimized_bytecode",
					ContentType: "application/json",
					Content:     optimized,
				})
			} else if err != nil {
				h.emit(emit, Frame{
					ID:        req.ID,
					Session:   req.Session,
					Op:        req.Op,
					Kind:      "err",
					FormIndex: formIndex,
					Stage:     "optimized-bytecode",
					Artifact:  "optimized_bytecode",
					Text:      err.Error(),
				})
			}
		}
		if !inspectable && req.Inspect.LoweredGo {
			h.emit(emit, Frame{
				ID:        req.ID,
				Session:   req.Session,
				Op:        req.Op,
				Kind:      "note",
				FormIndex: formIndex,
				Stage:     "lowered-go",
				Artifact:  "lowered_go",
				Text:      "Lowered Go extraction unsupported for this form",
			})
		}
		if inspectable && req.Inspect.LoweredGo {
			if lowered, err := h.inspectLoweredGo(sess.NS, form.Source); err == nil && lowered != "" {
				h.emit(emit, Frame{
					ID:          req.ID,
					Session:     req.Session,
					Op:          req.Op,
					Kind:        "artifact",
					FormIndex:   formIndex,
					Stage:       "lowered-go",
					Artifact:    "lowered_go",
					ContentType: "text/plain",
					Content:     lowered,
				})
			} else if err != nil {
				h.emit(emit, Frame{
					ID:        req.ID,
					Session:   req.Session,
					Op:        req.Op,
					Kind:      "err",
					FormIndex: formIndex,
					Stage:     "lowered-go",
					Artifact:  "lowered_go",
					Text:      err.Error(),
				})
			}
		}
		if outBuf.Len() > 0 {
			h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "out", FormIndex: formIndex, Text: outBuf.String()})
		}
		if errBuf.Len() > 0 {
			h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "err", FormIndex: formIndex, Stage: "runtime", Text: errBuf.String()})
		}
		if runErr != nil {
			h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "err", FormIndex: formIndex, Stage: "compile", Text: vm.FormatError(runErr)})
			h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "status", Status: []string{"done", "error"}})
			return
		}
		h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "value", FormIndex: formIndex, Value: result.String()})
	}
	h.emit(emit, Frame{ID: req.ID, Session: req.Session, Op: req.Op, Kind: "status", Status: []string{"done"}})
}

func (h *Host) getSession(id string) *session {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s := h.sessions[id]; s != nil {
		return s
	}
	s := &session{NS: "user"}
	h.sessions[id] = s
	return s
}

func (h *Host) emit(emit func(string), frame Frame) {
	if emit == nil {
		return
	}
	bs, err := EncodeFrame(frame)
	if err != nil {
		return
	}
	emit(string(bs))
}

func disassembleChunk(chunk *vm.CodeChunk) []map[string]any {
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

func projectConst(v vm.Value) any {
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
		return fmt.Sprintf("[%s] %s", v.Type().Name(), v.String())
	}
}

func opcodeStride(op int32) int {
	switch op & 0xff {
	case vm.OP_TRY_PUSH:
		return 3
	case vm.OP_RECUR:
		return 4
	case vm.OP_LOAD_ARG, vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE, vm.OP_JUMP,
		vm.OP_POP_N, vm.OP_DUP_NTH, vm.OP_INVOKE, vm.OP_LOAD_CLOSEDOVER,
		vm.OP_RECUR_FN, vm.OP_MAKE_MULTI_ARITY, vm.OP_TAIL_CALL,
		vm.OP_LOAD_CONST, vm.OP_LOAD_VAR, vm.OP_FINALLY_END:
		return 2
	default:
		return 1
	}
}

func NewResolver(ctx *compiler.Context) *resolver.NSResolver {
	r := resolver.NewNSResolver(ctx, []string{"."})
	rt.SetNSLoader(r)
	return r
}

func (h *Host) inspectIR(nsName, formSource string) (string, error) {
	v, err := h.evalInNS(nsName, fmt.Sprintf(
		"(require 'ir.build)\n(require 'ir.dump)\n(require 'ir.passes.pipeline)\n(ir.dump/dump (ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote %s))))",
		formSource,
	))
	if err != nil {
		return "", err
	}
	if s, ok := v.(vm.String); ok {
		return string(s), nil
	}
	return v.String(), nil
}

func (h *Host) inspectOptimizedBytecode(nsName, formSource string) ([]map[string]any, error) {
	v, err := h.evalInNS(nsName, fmt.Sprintf(
		"(require 'ir.passes.pipeline)\n(binding [ir.passes.pipeline/*target* :bytecode] (ir.passes.pipeline/compile-form (quote %s)))",
		formSource,
	))
	if err != nil {
		return nil, err
	}
	chunk := codeChunkFromValue(v)
	if chunk == nil {
		return nil, fmt.Errorf("optimized bytecode unavailable for %T", v)
	}
	return disassembleChunk(chunk), nil
}

func (h *Host) inspectLoweredGo(nsName, formSource string) (string, error) {
	source := fmt.Sprintf(`(require 'clojure.string)
(require 'gogen)
(require 'ir.passes.pipeline)
(let [res (binding [ir.passes.pipeline/*target* :go]
            (ir.passes.pipeline/compile-form (quote %s)))]
  (cond
    (:decl res) (gogen/render (:decl res))
    (:decls res) (clojure.string/join "\n\n" (map gogen/render (:decls res)))
    :else (pr-str res)))`, formSource)
	v, err := h.evalInNS(nsName, source)
	if err != nil {
		return "", err
	}
	if s, ok := v.(vm.String); ok {
		return string(s), nil
	}
	return v.String(), nil
}

func (h *Host) evalInNS(nsName, source string) (vm.Value, error) {
	ns := rt.NS(nsName)
	if ns == nil {
		ns = rt.NS("user")
	}
	// Per-inspect transient compiler — same rationale as <host-request>.
	c := compiler.NewTransientCompiler(h.consts, ns)
	c.SetSource("<host-inspect>")
	_, result, err := c.CompileMultiple(strings.NewReader(source))
	return result, err
}

func codeChunkFromValue(v vm.Value) *vm.CodeChunk {
	switch x := v.(type) {
	case *vm.Boxed:
		if chunk, ok := x.Unbox().(*vm.CodeChunk); ok {
			return chunk
		}
	case *vm.Func:
		return x.Chunk()
	case *vm.Var:
		// A defn evaluates to its Var; reach through to the fn body chunk so
		// the Bytecode pane shows the function's code, not the (near-empty)
		// top-level def chunk.
		return codeChunkFromValue(x.Root())
	}
	return nil
}

func isInspectableIRForm(v vm.Value) bool {
	seqable, ok := v.(vm.Sequable)
	if !ok {
		return false
	}
	seq := seqable.Seq()
	if seq == nil {
		return false
	}
	sym, ok := seq.First().(vm.Symbol)
	if !ok {
		return false
	}
	switch string(sym) {
	case "defn", "defn-", "defprotocol", "deftype", "defmulti", "defmethod":
		return true
	default:
		return false
	}
}
