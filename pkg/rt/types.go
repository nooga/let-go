package rt

import (
	"github.com/nooga/let-go/pkg/vm"
)

// FileStat is the struct behind os/stat return values.
type FileStat struct {
	Name    string `letgo:"name"`
	Size    int64  `letgo:"size"`
	IsDir   bool   `letgo:"dir?"`
	ModTime string `letgo:"mod-time"`
}

// ShellResult is the struct behind os/sh return values.
type ShellResult struct {
	Exit int    `letgo:"exit"`
	Out  string `letgo:"out"`
	Err  string `letgo:"err"`
}

// HTTPRequest is the struct behind ring-style HTTP request maps.
type HTTPRequest struct {
	RequestMethod string    `letgo:"request-method"`
	Scheme        string    `letgo:"scheme"`
	URI           string    `letgo:"uri"`
	Path          string    `letgo:"path"`
	QueryString   string    `letgo:"query-string"`
	Body          string    `letgo:"body"`
	RemoteAddr    string    `letgo:"remote-addr"`
	ServerAddr    string    `letgo:"server-addr"`
	ServerPort    string    `letgo:"server-port"`
	ContentType   string    `letgo:"content-type"`
	Headers       vm.Value  `letgo:"headers"`
}

// HTTPResponse is the struct behind HTTP client response maps.
type HTTPResponse struct {
	Status  int      `letgo:"status"`
	Body    string   `letgo:"body"`
	Headers vm.Value `letgo:"headers"`
}

var (
	fileStatMapping    *vm.StructMapping
	shellResultMapping *vm.StructMapping
	httpRequestMapping *vm.StructMapping
	httpResponseMapping *vm.StructMapping
)

func init() {
	fileStatMapping = vm.RegisterStruct[FileStat]("os/FileStat")
	shellResultMapping = vm.RegisterStruct[ShellResult]("os/ShellResult")
	httpRequestMapping = vm.RegisterStruct[HTTPRequest]("http/Request")
	httpResponseMapping = vm.RegisterStruct[HTTPResponse]("http/Response")
}
