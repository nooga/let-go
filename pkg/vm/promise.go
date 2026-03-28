package vm

import (
	"fmt"
	"sync"
)

// Promise is a write-once, read-many synchronization primitive.
// Created with (promise), written with (deliver p val), read with (deref p).
// Deref blocks until delivered.
type Promise struct {
	value     Value
	delivered bool
	mu        sync.Mutex
	ch        chan struct{} // closed on deliver to wake all waiters
}

func NewPromise() *Promise {
	return &Promise{ch: make(chan struct{})}
}

func (p *Promise) Deliver(val Value) Value {
	p.mu.Lock()
	if p.delivered {
		p.mu.Unlock()
		return NIL
	}
	p.value = val
	p.delivered = true
	close(p.ch)
	p.mu.Unlock()
	return val
}

func (p *Promise) Deref() Value {
	p.mu.Lock()
	if p.delivered {
		v := p.value
		p.mu.Unlock()
		return v
	}
	p.mu.Unlock()
	<-p.ch // block until delivered
	return p.value
}

func (p *Promise) IsRealized() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.delivered
}

func (p *Promise) Type() ValueType    { return PromiseType }
func (p *Promise) Unbox() interface{} { return p }
func (p *Promise) String() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.delivered {
		return fmt.Sprintf("#<Promise@%p: %s>", p, p.value.String())
	}
	return fmt.Sprintf("#<Promise@%p: :pending>", p)
}

type thePromiseType struct{}

func (t *thePromiseType) String() string     { return t.Name() }
func (t *thePromiseType) Type() ValueType    { return TypeType }
func (t *thePromiseType) Unbox() interface{} { return nil }
func (t *thePromiseType) Name() string       { return "let-go.lang.Promise" }
func (t *thePromiseType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var PromiseType *thePromiseType = &thePromiseType{}
