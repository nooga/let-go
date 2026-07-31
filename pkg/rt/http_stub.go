//go:build tinygo || lg_no_http

/*
 * Stub for the two lanes that build without net/http.
 *
 * TinyGo: net/http is broken in TinyGo's stdlib (roundtrip_js.go references a
 * method that doesn't exist on its own Transport), so the lane has never had a
 * working http namespace; browser callers use JS fetch via the js bridge.
 *
 * lg_no_http: an opt-in tag for stock-Go builds that embed the runtime and
 * never serve or fetch. net/http pulls crypto/tls and crypto/x509 with it, and
 * the linker cannot drop them — installHttpNS is reached from an init(), so
 * every symbol behind it is live no matter what the program does.
 */

package rt

import (
	"fmt"
	"io"
)

func installHttpNS() {}

func slurpURL(url string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("URLs need net/http, which this build omits (under tinygo, use JS fetch)")
}
