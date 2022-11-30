package main

import (
	"fmt"
	"syscall/js"

	"github.com/nooga/let-go/pkg/api"
)

var lg *api.LetGo

func Eval(this js.Value, args []js.Value) interface{} {
	x := args[0].String()
	v, err := lg.Run(x)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return v.String()
}

func main() {
	var err error
	lg, err = api.NewLetGo("user")
	if err != nil {
		panic("let-go runtime failed to init")
	}
	js.Global().Set("Eval", js.FuncOf(Eval))

	<-make(chan struct{})
}
