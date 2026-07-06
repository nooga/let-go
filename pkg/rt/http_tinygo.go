//go:build tinygo

/*
 * TinyGo stub: net/http is broken in TinyGo's stdlib (roundtrip_js.go
 * references a method that doesn't exist on its own Transport). Under
 * TinyGo, the http namespace is unavailable; browser callers should use
 * JS fetch via the js bridge instead.
 */

package rt

import (
	"fmt"
	"io"
)

func installHttpNS() {}

func slurpURL(url string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("http not available under tinygo (use JS fetch)")
}
