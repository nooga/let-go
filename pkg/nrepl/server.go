/*
 * Copyright (c) 2022-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package nrepl

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/go-uuid"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/bencode"
)

// session holds per-client state.
type session struct {
	id  string
	ctx *compiler.Context
}

// NreplServer implements the nREPL protocol over TCP.
type NreplServer struct {
	ctx      *compiler.Context
	listener net.Listener
	stop     chan struct{}
	wg       sync.WaitGroup
	mu       sync.Mutex
	sessions map[string]*session
	port     int
}

func NewNreplServer(ctx *compiler.Context) *NreplServer {
	return &NreplServer{
		ctx:      ctx,
		sessions: make(map[string]*session),
	}
}

func (n *NreplServer) newSession() *session {
	id, err := uuid.GenerateUUID()
	if err != nil {
		id = "fallback-session"
	}
	// Each session gets its own compiler context sharing the same namespace
	s := &session{
		id:  id,
		ctx: n.ctx,
	}
	n.mu.Lock()
	n.sessions[id] = s
	n.mu.Unlock()
	return s
}

func (n *NreplServer) closeSession(id string) {
	n.mu.Lock()
	delete(n.sessions, id)
	n.mu.Unlock()
}

func (n *NreplServer) sessionIDs() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	ids := make([]string, 0, len(n.sessions))
	for id := range n.sessions {
		ids = append(ids, id)
	}
	return ids
}

// Start starts the nREPL server on the given port.
func (n *NreplServer) Start(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	n.listener = l
	n.port = port
	n.stop = make(chan struct{})

	// Write .nrepl-port file
	os.WriteFile(".nrepl-port", []byte(fmt.Sprintf("%d", port)), 0644)

	fmt.Printf("nREPL server started on port %d on host 127.0.0.1 - nrepl://127.0.0.1:%d\n", port, port)

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		for {
			select {
			case <-n.stop:
				return
			default:
				conn, err := n.listener.Accept()
				if err != nil {
					select {
					case <-n.stop:
						return
					default:
						continue
					}
				}
				go n.handleConn(conn)
			}
		}
	}()
	return nil
}

// Stop shuts down the server and cleans up.
func (n *NreplServer) Stop() {
	close(n.stop)
	n.listener.Close()
	n.wg.Wait()
	os.Remove(".nrepl-port")
}

// handleConn processes a single client connection.
func (n *NreplServer) handleConn(conn net.Conn) {
	defer conn.Close()
	dec := bencode.NewDecoder(conn)

	for {
		var msg map[string]interface{}
		err := dec.Decode(&msg)
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
		n.handleMsg(conn, msg)
	}
}

// handleMsg dispatches a single nREPL message.
func (n *NreplServer) handleMsg(conn net.Conn, msg map[string]interface{}) {
	op, _ := msg["op"].(string)
	id := msgStr(msg, "id")
	sessID := msgStr(msg, "session")

	switch op {
	case "clone":
		s := n.newSession()
		respond(conn, map[string]interface{}{
			"id":          id,
			"status":      []string{"done"},
			"new-session": s.id,
		})

	case "close":
		n.closeSession(sessID)
		respond(conn, map[string]interface{}{
			"id":      id,
			"session": sessID,
			"status":  []string{"done", "session-closed"},
		})

	case "describe":
		respond(conn, map[string]interface{}{
			"id":      id,
			"session": sessID,
			"ops": map[string]interface{}{
				"clone":       map[string]interface{}{},
				"close":       map[string]interface{}{},
				"eval":        map[string]interface{}{},
				"load-file":   map[string]interface{}{},
				"describe":    map[string]interface{}{},
				"completions": map[string]interface{}{},
				"lookup":      map[string]interface{}{},
				"info":        map[string]interface{}{},
				"complete":    map[string]interface{}{},
				"ls-sessions": map[string]interface{}{},
			},
			"versions": map[string]interface{}{
				"let-go": map[string]interface{}{
					"major": "1", "minor": "0",
				},
				"nrepl": map[string]interface{}{
					"major": "1", "minor": "0",
				},
			},
			"aux":    map[string]interface{}{},
			"status": []string{"done"},
		})

	case "eval":
		n.handleEval(conn, msg)

	case "load-file":
		// Transform to eval
		code := msgStr(msg, "file")
		msg["code"] = code
		n.handleEval(conn, msg)

	case "completions", "complete":
		n.handleCompletions(conn, msg)

	case "info", "lookup":
		n.handleInfo(conn, msg)

	case "ls-sessions":
		respond(conn, map[string]interface{}{
			"id":       id,
			"session":  sessID,
			"sessions": n.sessionIDs(),
			"status":   []string{"done"},
		})

	case "interrupt":
		respond(conn, map[string]interface{}{
			"id":      id,
			"session": sessID,
			"status":  []string{"done", "session-idle"},
		})

	default:
		respond(conn, map[string]interface{}{
			"id":      id,
			"session": sessID,
			"status":  []string{"done", "error", "unknown-op"},
		})
	}
}

// handleEval evaluates code and streams out/err/value/done messages.
func (n *NreplServer) handleEval(conn net.Conn, msg map[string]interface{}) {
	id := msgStr(msg, "id")
	sessID := msgStr(msg, "session")
	code := msgStr(msg, "code")

	// Capture stdout
	var outBuf bytes.Buffer
	origStdout := os.Stdout
	pr, pw, _ := os.Pipe()
	os.Stdout = pw

	// Read pipe in background
	outDone := make(chan struct{})
	go func() {
		io.Copy(&outBuf, pr)
		close(outDone)
	}()

	// Eval
	_, val, err := n.ctx.CompileMultiple(strings.NewReader(code))

	// Restore stdout and flush
	pw.Close()
	os.Stdout = origStdout
	<-outDone

	// Send stdout if any
	if outBuf.Len() > 0 {
		respond(conn, map[string]interface{}{
			"id":      id,
			"session": sessID,
			"out":     outBuf.String(),
		})
	}

	if err != nil {
		errStr := vm.FormatError(err)
		respond(conn, map[string]interface{}{
			"id":      id,
			"session": sessID,
			"err":     errStr + "\n",
			"ex":      "let-go.lang.Error",
			"root-ex": "let-go.lang.Error",
		})
	} else {
		valStr := "nil"
		if val != nil {
			valStr = val.String()
		}
		respond(conn, map[string]interface{}{
			"id":      id,
			"session": sessID,
			"value":   valStr,
			"ns":      n.ctx.CurrentNS().Name(),
		})
	}

	// Always end with done
	respond(conn, map[string]interface{}{
		"id":      id,
		"session": sessID,
		"status":  []string{"done"},
	})
}

// handleCompletions returns completion candidates.
func (n *NreplServer) handleCompletions(conn net.Conn, msg map[string]interface{}) {
	id := msgStr(msg, "id")
	sessID := msgStr(msg, "session")

	prefix := msgStr(msg, "prefix")
	if prefix == "" {
		prefix = msgStr(msg, "symbol")
	}

	var completions []interface{}
	if prefix != "" {
		sym := vm.Symbol(prefix)
		matches := rt.FuzzyNamespacedSymbolLookup(n.ctx.CurrentNS(), sym)
		for _, m := range matches {
			completions = append(completions, map[string]interface{}{
				"candidate": string(m),
				"type":      "function",
			})
		}
	}

	if completions == nil {
		completions = []interface{}{}
	}

	respond(conn, map[string]interface{}{
		"id":          id,
		"session":     sessID,
		"completions": completions,
		"status":      []string{"done"},
	})
}

// handleInfo returns symbol info.
func (n *NreplServer) handleInfo(conn net.Conn, msg map[string]interface{}) {
	id := msgStr(msg, "id")
	sessID := msgStr(msg, "session")
	op := msgStr(msg, "op")

	sym := msgStr(msg, "sym")
	if sym == "" {
		sym = msgStr(msg, "symbol")
	}

	nsName := msgStr(msg, "ns")
	if nsName == "" {
		nsName = n.ctx.CurrentNS().Name()
	}

	resp := map[string]interface{}{
		"id":      id,
		"session": sessID,
		"status":  []string{"done"},
	}

	if sym != "" {
		// Try to look up the symbol
		ns := n.ctx.CurrentNS()
		if nsName != "" && nsName != ns.Name() {
			if found := rt.NS(nsName); found != nil {
				ns = found
			}
		}
		v := ns.Lookup(vm.Symbol(sym))
		if v != vm.NIL {
			info := map[string]interface{}{
				"name": sym,
				"ns":   nsName,
			}
			if op == "info" {
				// CIDER info returns flat
				resp["name"] = sym
				resp["ns"] = nsName
			} else {
				// lookup nests under "info"
				resp["info"] = info
			}
		} else {
			resp["status"] = []string{"done", "no-info"}
		}
	}

	respond(conn, resp)
}

// --- Helpers ---

func respond(conn net.Conn, msg map[string]interface{}) {
	bs, err := bencode.EncodeBytes(msg)
	if err != nil {
		return
	}
	conn.Write(bs)
}

func msgStr(msg map[string]interface{}, key string) string {
	v, ok := msg[key]
	if !ok || v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	case int64:
		return fmt.Sprintf("%d", s)
	default:
		return fmt.Sprintf("%v", v)
	}
}
