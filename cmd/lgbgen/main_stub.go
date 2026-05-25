//go:build !bootstrap

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "lgbgen must be run with source bootstrap enabled: go run -tags bootstrap ./cmd/lgbgen [output-path]")
	os.Exit(2)
}
