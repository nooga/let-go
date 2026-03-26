package api

import (
	"reflect"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type LetGo struct {
	cp     *vm.Consts
	c      *compiler.Context
	loader *resolver.NSResolver
}

func NewLetGo(ns string) (*LetGo, error) {
	cp := vm.NewConsts()
	nso := rt.NS(ns)
	c := compiler.NewCompiler(cp, nso)
	ret := &LetGo{
		cp:     cp,
		c:      c,
		loader: resolver.NewNSResolver(c, []string{"."}),
	}
	rt.SetNSLoader(ret.loader)
	return ret, nil
}

func (l *LetGo) SetLoadPath(path []string) {
	l.loader.SetPath(path)
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
	result, err := frame.Run()
	vm.ReleaseFrame(frame)
	return result, err
}
