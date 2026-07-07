package wasm

import (
	"strings"
	"testing"
)

func TestAssembleHTMLWithTemplate(t *testing.T) {
	tmpl := "<html><body><script>\n" + HostBodyMarker + "\n</script></body></html>"
	out := AssembleHTMLWithTemplate(tmpl, "// wasm_exec", "GZB64==", false, true)

	if strings.Contains(out, HostBodyMarker) {
		t.Fatal("marker was not consumed")
	}
	if !strings.Contains(out, "// wasm_exec") {
		t.Fatal("wasm_exec JS not injected")
	}
	if strings.Contains(out, "new Terminal") {
		t.Fatal("custom-template output must not include the xterm shell")
	}
	if !strings.HasPrefix(out, "<html><body>") {
		t.Fatal("custom template body not preserved")
	}
}
