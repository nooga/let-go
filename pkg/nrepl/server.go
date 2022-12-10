package nrepl

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/go-uuid"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/bencode"
)

type session struct {
	ctx       *compiler.Context
	sessionID string
	lastID    int
	conn      net.Conn
	dec       *bencode.Decoder
	describe  map[string]interface{}
}

func newSession(ctx *compiler.Context, conn net.Conn) *session {
	empty := map[string]interface{}{}
	return &session{
		ctx:       ctx,
		sessionID: "none",
		lastID:    0,
		conn:      conn,
		dec:       bencode.NewDecoder(conn),
		describe: map[string]interface{}{
			"aux":    empty,
			"status": []string{"done"},
			"ops": map[string]interface{}{
				"clone":    empty,
				"eval":     empty,
				"describe": empty,
			},
			"versions": map[string]interface{}{
				"let-go":       "dev",
				"let-go.nrepl": "dev",
			},
		},
	}
}

func (s *session) handle() error {
	var m map[string]interface{}
	err := s.dec.Decode(&m)
	if err != nil {
		return err
	}
	fmt.Println("handle", m)
	op := m["op"].(string)
	id := m["id"].(string)
	s.lastID, err = strconv.Atoi(id)
	if err != nil {
		return err
	}
	switch op {
	case "eval":
		code := m["code"].(string)
		var v vm.Value
		_, v, err = s.ctx.CompileMultiple(strings.NewReader(code))
		if err != nil {
			fmt.Printf("Error evaluating: %s, %#v\n", code, err)
			err = s.respond(m, map[string]interface{}{
				"err":    fmt.Sprintf("%s", err),
				"status": []string{"eval-error"},
			})
			if err != nil {
				break
			}
		}
		fmt.Println("eval", code, v)
		err = s.respond(m, map[string]interface{}{
			"value":  v.String(),
			"ns":     s.ctx.CurrentNS().Name(),
			"status": []string{"done"},
		})
	case "describe":
		err = s.respond(m, s.describe)
	}
	if err != nil {
		fmt.Printf("Error responding: %#v\n", err)
		return err
	}
	if m["session"] == nil {
		err = s.send(map[string]interface{}{
			"new-session": s.sessionID,
			"status":      []string{"done"},
		})
		if err == nil {
			err = s.send(s.describe)
		}
	}
	return err
}

func (s *session) respond(m map[string]interface{}, o map[string]interface{}) error {
	id, err := strconv.Atoi(m["id"].(string))
	if err != nil {
		return err
	}
	session := "none"
	ses, ok := m["session"]
	if ok && ses != nil {
		session = ses.(string)
	}
	o["id"] = fmt.Sprintf("%d", id)
	o["session"] = session
	fmt.Println("respond", o)
	bs, err := bencode.EncodeBytes(o)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (s *session) send(o map[string]interface{}) error {
	s.lastID += 1
	id := s.lastID
	session := s.sessionID
	o["id"] = fmt.Sprintf("%d", id)
	o["session"] = session
	fmt.Println("send", o)
	bs, err := bencode.EncodeBytes(o)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

type NreplServer struct {
	ctx      *compiler.Context
	listener net.Listener
	stop     chan struct{}
	wg       sync.WaitGroup
	sessions map[string]*session
}

func NewNreplServer(compiler *compiler.Context) *NreplServer {
	return &NreplServer{ctx: compiler, sessions: map[string]*session{}}
}

func (n *NreplServer) adoptSession(s *session) string {
	newID, err := uuid.GenerateUUID()
	if err != nil {
		panic("wtf uuid gen failed")
	}
	n.sessions[newID] = s
	s.sessionID = newID
	return newID
}

func (n *NreplServer) Start(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	n.listener = l
	n.stop = make(chan struct{})

	n.wg.Add(1)
	go func() {
		for {
			select {
			case <-n.stop:
				goto done
			default:
				conn, err := n.listener.Accept()
				if err != nil {
					fmt.Println("error when accepting", err)
					continue
				}
				go func(conn net.Conn) {
					s := newSession(n.ctx, conn)
					n.adoptSession(s)
					for {
						err = s.handle()
						if err != nil {
							fmt.Println("handle failed", err)
							if err == io.EOF {
								break
							}
						}
					}
					conn.Close()
				}(conn)
			}
		}
	done:
		n.wg.Done()
	}()
	return nil
}

func (n *NreplServer) Stop() {
	close(n.stop)
	n.wg.Wait()
}
