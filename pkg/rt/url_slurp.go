//go:build !tinygo

package rt

import (
	"io"
	"net/http"
)

func slurpURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
