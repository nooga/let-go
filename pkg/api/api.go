package api

import (
	"reflect"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type LetGo struct {
	cp *vm.Consts
	c  *compiler.Context
}

func NewLetGo(ns string) (*LetGo, error) {
	cp := vm.NewConsts()
	nso := rt.NS(ns)
	ret := &LetGo{
		cp: cp,
		c:  compiler.NewCompiler(cp, nso),
	}
	return ret, nil
}

func (l *LetGo) Def(name string, value interface{}) error {
	val, err := vm.BoxValue(reflect.ValueOf(value))
	if err != nil {
		return err
	}
	l.c.CurrentNS().Def(name, val)

	return nil
}

func (l *LetGo) Run(expr string) (vm.Value, error) {
	c, err := l.c.Compile(expr)
	if err != nil {
		return vm.NIL, err
	}
	frame := vm.NewFrame(c, nil)
	return frame.Run()
}
